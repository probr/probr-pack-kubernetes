package kubernetes

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"citihub.com/probr/internal/clouddriver/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
)

//PrivilegedAccess ...
type PrivilegedAccess int

//PrivilegedAccess enum
const (
	WithPrivilegedAccess PrivilegedAccess = iota
	WithoutPrivilegedAccess
)

//PSPTestCommand ...
type PSPTestCommand int

//PSPTestCommand ...
const (
	Chroot PSPTestCommand = iota
	EnterHostPIDNS
	EnterHostIPCNS
	EnterHostNetworkNS
	VerifyNonRootUID
	NetRawTest
	SpecialCapTest
)

func (c PSPTestCommand) String() string {
	return [...]string{"chroot .",
		"nsenter -t 1 -p ps",
		"nsenter -t 1 -i ps",
		"nsenter -t 1 -n ps",
		"id -u > 0 ",
		"ping google.com",
		"ip link add dummy0 type dummy"}[c]
}

const (
	//TODO: default to these values for MVP - need to expose in future
	pspTestNamespace = "probr-pod-security-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or busybox need to be available in the allowed (probably internal) registry
	pspTestImage     = "docker.io/busybox"
	pspTestContainer = "psp-test"
	pspTestPodName   = "psp-test-pod"
)

// PodSecurityPolicy ...
type PodSecurityPolicy interface {
	ClusterIsDeployed() *bool
	ClusterHasPSP() (*bool, error)
	PrivilegedAccessIsRestricted() (*bool, error)
	HostPIDIsRestricted() (*bool, error)
	HostIPCIsRestricted() (*bool, error)
	HostNetworkIsRestricted() (*bool, error)
	PrivilegedEscalationIsRestricted() (*bool, error)
	RootUserIsRestricted() (*bool, error)
	NETRawIsRestricted() (*bool, error)
	AllowedCapabilitiesAreRestricted() (*bool, error)
	AssignedCapabilitiesAreRestricted() (*bool, error)
	CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64) (*apiv1.Pod, error)
	CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool) (*apiv1.Pod, error)
	CreatePODSettingCapabilities(c *[]string) (*apiv1.Pod, error)
	ExecPSPTestCmd(pName *string, cmd PSPTestCommand) (int, error)
	TeardownPodSecurityTestPod(p *string) error
}

// PSP ...
type PSP struct {
	k                       Kubernetes
	securityPolicyProviders *[]SecurityPolicyProvider
}

// NewPSP ...
func NewPSP(k Kubernetes, sp *[]SecurityPolicyProvider) *PSP {
	p := &PSP{}
	p.k = k
	p.securityPolicyProviders = sp

	return p
}

// NewDefaultPSP ...
func NewDefaultPSP() *PSP {
	p := &PSP{}
	p.k = GetKubeInstance()

	//standard security providers
	p.securityPolicyProviders = &[]SecurityPolicyProvider{
		NewKubeSecurityPolicyProvider(p.k),
		azure.NewAzPolicyProvider()}

	return p

}

// ClusterIsDeployed ...
func (psp *PSP) ClusterIsDeployed() *bool {
	return psp.k.ClusterIsDeployed()
}

//ClusterHasPSP determines if the cluster has any Pod Security Policies set.
func (psp *PSP) ClusterHasPSP() (*bool, error) {
	var err error = nil
	var ret *bool = nil

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, e := p.HasSecurityPolicies()
		if e != nil {
			//hold onto the error and continue
			err = e
			continue
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			//nil out any errors
			return b, nil
		}
		//if no policies (but a successful call), then make sure the ret value is set
		ret = b
	}

	//if we get to here, we haven't got any, but that could have been because of
	//errors from all providers, in which case "ret" will be nil
	if ret != nil {
		//then we've made at least one successful call - nil out err, for client simplification
		return ret, nil
	}

	//otherwise just return
	return ret, err
}

//PrivilegedAccessIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) PrivilegedAccessIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasPrivilegedAccessRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//HostPIDIsRestricted looks for a PodSecurityPolicy with 'HostPID' set to false (i.e. NO Access to HostPID ).
func (psp *PSP) HostPIDIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasHostPIDRestriction()
		if err != nil {
			return nil, err
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil

}

//HostIPCIsRestricted looks for a PodSecurityPolicy with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
func (psp *PSP) HostIPCIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasHostIPCRestriction()
		if err != nil {
			return nil, err
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil

}

//HostNetworkIsRestricted looks for a PodSecurityPolicy with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
func (psp *PSP) HostNetworkIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasHostNetworkRestriction()
		if err != nil {
			return nil, err
		}
		if *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil

}

//PrivilegedEscalationIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) PrivilegedEscalationIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasAllowPrivilegeEscalationRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

// RootUserIsRestricted ...
func (psp *PSP) RootUserIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasRootUserRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//NETRawIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) NETRawIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasNETRAWRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//AllowedCapabilitiesAreRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) AllowedCapabilitiesAreRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasAllowedCapabilitiesRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

//AssignedCapabilitiesAreRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) AssignedCapabilitiesAreRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		b, err := p.HasAssignedCapabilitiesRestriction()
		if err != nil {
			return nil, err
		}
		if b != nil && *b {
			//return on first 'true' - only trying to establish if we have any ...
			return b, nil
		}
	}

	//if we get to here, we haven't got any ...
	b := false
	return &b, nil
}

// CreatePODSettingSecurityContext creates POD with a SecurityContext conforming to the parameters:
// pr *bool - set the Privileged flag.  Defaults to false.
// pe *bool - set the Allow Privileged Escalation flag.  Defaults to false.
// runAsUser *int64 - set RunAsUser.  Defaults to 1000.
func (psp *PSP) CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64) (*apiv1.Pod, error) {
	//default sensibly if not provided
	//this needs to take account of rules around allowedPrivilegdEscalation and Privileged:
	// cannot set `allowPrivilegeEscalation` to false and `privileged` to true
	f := false
	if pr == nil {
		if pe != nil {
			pr = pe //set pr to pe value
		}
		//if pe is also nil then just set them both to false
		if pe == nil {
			pe, pr = &f, &f
		}
	}
	if pe == nil {
		if pr != nil {
			pe = pr //set pe to pr value
		}
		//if pr is also nil then just set them both to false
		if pr == nil {
			pe, pr = &f, &f
		}
	}
	if runAsUser == nil {
		i := int64(1000)
		runAsUser = &i
	}

	sc := apiv1.SecurityContext{
		Privileged:               pr,
		AllowPrivilegeEscalation: pe,
		RunAsUser:                runAsUser,
	}

	pname, ns, cname, image := GenerateUniquePodName(pspTestPodName), pspTestNamespace, pspTestContainer, pspTestImage

	p, err := psp.k.CreatePod(&pname, &ns, &cname, &image, true, &sc)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// CreatePODSettingAttributes creates a POD with attributes:
// hostPID *bool - set the hostPID flag, defaults to false
// hostIPC *bool - set the hostIPC flag, defaults to false
// hostNetwork *bool - set the hostNetwork flag, defaults to false
func (psp *PSP) CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool) (*apiv1.Pod, error) {
	//default sensibly if not provided
	f := false
	if hostPID == nil {
		hostPID = &f
	}
	if hostIPC == nil {
		hostIPC = &f
	}
	if hostNetwork == nil {
		hostNetwork = &f
	}

	pname, ns, cname, image := GenerateUniquePodName(pspTestPodName), pspTestNamespace, pspTestContainer, pspTestImage

	// get the pod object and manipulate:
	po := psp.k.GetPodObject(pname, ns, cname, image, nil)
	po.Spec.HostPID = *hostPID
	po.Spec.HostIPC = *hostIPC
	po.Spec.HostNetwork = *hostNetwork

	// create from PO
	p, err := psp.k.CreatePodFromObject(po, &pname, &ns, true)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//CreatePODSettingCapabilities ...
func (psp *PSP) CreatePODSettingCapabilities(c *[]string) (*apiv1.Pod, error) {
	pname, ns, cname, image := GenerateUniquePodName(pspTestPodName), pspTestNamespace, pspTestContainer, pspTestImage

	// get the pod object and manipulate:
	po := psp.k.GetPodObject(pname, ns, cname, image, nil)

	if c != nil {
		for _, cap := range *c {
			for _, con := range po.Spec.Containers {
				if con.SecurityContext == nil {
					con.SecurityContext = &apiv1.SecurityContext{}
				}
				if con.SecurityContext.Capabilities == nil {
					con.SecurityContext.Capabilities = &apiv1.Capabilities{}
					con.SecurityContext.Capabilities.Add = make([]apiv1.Capability, 0)
				}
				con.SecurityContext.Capabilities.Add =
					append(con.SecurityContext.Capabilities.Add, apiv1.Capability(cap))
			}
		}
	}

	// create from PO
	p, err := psp.k.CreatePodFromObject(po, &pname, &ns, true)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// ExecPSPTestCmd ...
func (psp *PSP) ExecPSPTestCmd(pName *string, cmd PSPTestCommand) (int, error) {
	var pn string
	//if we've not been given a pod name, assume one needs to be created:
	if pName == nil {
		//want one without privileged access or escalation
		f := false
		p, err := psp.CreatePODSettingSecurityContext(&f, &f, nil)

		if err != nil {
			return -1, err
		}
		//grab the pod name:
		pn = p.GetObjectMeta().GetName()
	} else {
		pn = *pName
	}

	c := cmd.String()
	ns := pspTestNamespace
	stdout, _, ex, err := psp.k.ExecCommand(&c, &ns, &pn)

	log.Printf("[NOTICE] ExecPSPTestCmd: %v stdout: %v exit code: %v", cmd, stdout, ex)

	if err != nil {
		return ex, err
	}

	return ex, nil
}

//TeardownPodSecurityTestPod ...
func (psp *PSP) TeardownPodSecurityTestPod(p *string) error {
	ns := pspTestNamespace
	err := psp.k.DeletePod(p, &ns, false) //don't worry about waiting
	return err
}

// KubeSecurityPolicyProvider ...
type KubeSecurityPolicyProvider struct {
	k        Kubernetes
	psps     *v1beta1.PodSecurityPolicyList
	pspMutex sync.Mutex
}

// NewKubeSecurityPolicyProvider ...
func NewKubeSecurityPolicyProvider(k Kubernetes) *KubeSecurityPolicyProvider {
	return &KubeSecurityPolicyProvider{k: k}
}

func (p *KubeSecurityPolicyProvider) getPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	p.pspMutex.Lock()
	defer p.pspMutex.Unlock()

	//already got them?
	if p.psps == nil {
		ps, err := p.getPodSecurityPolicies()
		if err != nil {
			return nil, err
		}
		p.psps = ps
	}
	return p.psps, nil
}

// HasSecurityPolicies ...
func (p *KubeSecurityPolicyProvider) HasSecurityPolicies() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	b := len(psps.Items) > 0
	return &b, nil
}

// HasPrivilegedAccessRestriction ...
func (p *KubeSecurityPolicyProvider) HasPrivilegedAccessRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have Privileged set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.Privileged {
			log.Printf("[NOTICE] PASS: Privileged is set to %v on Policy: %v", e.Spec.Privileged, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with Privileged set.\n")
	}

	return &res, nil

}

// HasHostPIDRestriction ...
func (p *KubeSecurityPolicyProvider) HasHostPIDRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostPID set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostPID {
			log.Printf("[NOTICE] PASS: HostPID is set to %v on Policy: %v\n", e.Spec.HostPID, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with HostPID set.\n")
	}

	return &res, nil

}

// HasHostIPCRestriction ...
func (p *KubeSecurityPolicyProvider) HasHostIPCRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostIPC set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostIPC {
			log.Printf("[NOTICE] PASS: HostIPC is set to %v on Policy: %v\n", e.Spec.HostIPC, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with HostIPC set.\n")
	}

	return &res, nil

}

// HasHostNetworkRestriction ...
func (p *KubeSecurityPolicyProvider) HasHostNetworkRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostNetwork set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostNetwork {
			log.Printf("[NOTICE] PASS: HostNetwork is set to %v on Policy: %v\n", e.Spec.HostNetwork, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with HostNetwork set.\n")
	}

	return &res, nil

}

// HasAllowPrivilegeEscalationRestriction ...
func (p *KubeSecurityPolicyProvider) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have AllowPrivilegeEscalation set to false
	var res bool
	for _, e := range psps.Items {
		if !*e.Spec.AllowPrivilegeEscalation {
			log.Printf("[NOTICE] PASS: AllowPrivilegeEscalation is set to %v on Policy: %v", e.Spec.AllowPrivilegeEscalation, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with AllowPrivilegeEscalation set.\n")
	}

	return &res, nil

}

// HasRootUserRestriction ...
func (p *KubeSecurityPolicyProvider) HasRootUserRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have AllowPrivilegeEscalation set to false
	var res bool
	for _, e := range psps.Items {
		if e.Spec.RunAsUser.Rule == v1beta1.RunAsUserStrategyMustRunAsNonRoot {
			log.Printf("[NOTICE] PASS: RunAsUserStrategyMustRunAsNonRoot is set on Policy: %v", e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] FAIL: NO Policies found with AllowPrivilegeEscalation set.\n")
	}

	return &res, nil
}

// HasNETRAWRestriction ...
func (p *KubeSecurityPolicyProvider) HasNETRAWRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least one of the PSPs should have a RequiredDropCapability of "NET_RAW"
	var res bool
	for _, e := range psps.Items {
		for _, c := range e.Spec.RequiredDropCapabilities {
			if c == "NET_RAW" || c == "ALL" {
				log.Printf("[NOTICE] HasNETRAWRestriction PASS: RequiredDropCapability of %v is set on Policy: %v", c, e.GetName())
				res = true
				break
			}
		}
	}

	if !res {
		log.Printf("[NOTICE] HasNETRAWRestriction FAIL: NO Policies found with RequiredDropCapability of NET_RAW set.")
	}

	return &res, nil
}

// HasAllowedCapabilitiesRestriction ...
func (p *KubeSecurityPolicyProvider) HasAllowedCapabilitiesRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//in this case we don't want "allowedCapabilities" on any PSP (default to true)
	res := true
	for _, e := range psps.Items {
		if e.Spec.AllowedCapabilities != nil && len(e.Spec.AllowedCapabilities) > 0 {
			log.Printf("[NOTICE] HasAllowedCapabilitiesRestriction FAIL: at least one AllowedCapability is set on Policy: %v", e.GetName())
			res = false
			break
		}
	}

	if res {
		log.Printf("[NOTICE] HasNETRAWRestriction PASS: NO Policies found with AllowedCapabilities")
	}

	return &res, nil
}

// HasAssignedCapabilitiesRestriction ...
func (p *KubeSecurityPolicyProvider) HasAssignedCapabilitiesRestriction() (*bool, error) {
	f := false
	return &f, nil
}

func (p *KubeSecurityPolicyProvider) getPodSecurityPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	c, err := p.k.GetClient()
	if err != nil {
		return nil, err
	}

	ps := c.PolicyV1beta1().PodSecurityPolicies()
	if ps == nil {
		return nil, fmt.Errorf("Pod Security Polices could not be obtained (nil returned)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pspList, err := ps.List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if pspList == nil {
		return nil, fmt.Errorf("Pod Security Polices list returned a nil list")
	}

	log.Printf("[NOTICE] There are %d psp policies in the cluster\n", len(pspList.Items))

	for _, e := range pspList.Items {
		log.Printf("[INFO] PSP: %v \n", e.GetName())
		log.Printf("[INFO] Spec: %+v \n", e.Spec)
	}

	return pspList, nil
}
