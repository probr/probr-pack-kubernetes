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

const (
	//TODO: default to these values for MVP - need to expose in future
	//TODO: also, using network-access-test-ns here as there's an exclusion on the
	//container registry - needs to be tidied up ...
	pspTestNamespace = "probr-pod-security-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or busybox need to be available in the allowed (probably internal) registry
	pspTestImage     = "docker.io/busybox"
	pspTestContainer = "psp-test"
	pspTestPodName   = "psp-test-pod"
)

var securityPolicyProviders = []SecurityPolicyProvider{
	&KubeSecurityPolicyProvider{},
	azure.NewAzPolicyProvider()}

//ClusterHasPSP determines if the cluster has any Pod Security Policies set.
func ClusterHasPSP() (*bool, error) {
	var err error = nil
	var ret *bool = nil

	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func PrivilegedAccessIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func HostPIDIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func HostIPCIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func HostNetworkIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func PrivilegedEscalationIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func RootUserIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func NETRawIsRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func AllowedCapabilitiesAreRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func AssignedCapabilitiesAreRestricted() (*bool, error) {
	// iterate over providers ...
	for _, p := range securityPolicyProviders {
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
func CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64) (*apiv1.Pod, error) {
	//default sensibly if not provided
	f := false
	if pr == nil {
		pr = &f
	}
	if pe == nil {
		pe = &f		
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

	pname, ns, cname, image := pspTestPodName, pspTestNamespace, pspTestContainer, pspTestImage

	p, err := CreatePod(&pname, &ns, &cname, &image, true, &sc)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// CreatePODSettingAttributes creates a POD with attributes:
// hostPID *bool - set the hostPID flag, defaults to false
// hostIPC *bool - set the hostIPC flag, defaults to false
// hostNetwork *bool - set the hostNetwork flag, defaults to false
func CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool) (*apiv1.Pod, error) {
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

	pname, ns, cname, image := pspTestPodName, pspTestNamespace, pspTestContainer, pspTestImage

	// get the pod object and manipulate:
	po := GetPodObject(pname, ns, cname, image, nil)
	po.Spec.HostPID = *hostPID
	po.Spec.HostIPC = *hostIPC
	po.Spec.HostNetwork = *hostNetwork

	// create from PO
	p, err := CreatePodFromObject(po, &pname, &ns, true)

	if err != nil {
		return nil, err
	}

	return p, nil
}

//CreatePODSettingCapabilities ...
func CreatePODSettingCapabilities(c *[]string) (*apiv1.Pod, error) {	
	pname, ns, cname, image := pspTestPodName, pspTestNamespace, pspTestContainer, pspTestImage

	// get the pod object and manipulate:
	po := GetPodObject(pname, ns, cname, image, nil)
	
	if c != nil {
		for _, cap := range *c {
			for _, con := range po.Spec.Containers {
				if con.SecurityContext == nil {
					con.SecurityContext = &apiv1.SecurityContext{}					
				}
				if con.SecurityContext.Capabilities == nil {
					con.SecurityContext.Capabilities = &apiv1.Capabilities{}
					con.SecurityContext.Capabilities.Add = make([]apiv1.Capability,0)
				}
				con.SecurityContext.Capabilities.Add = 
					append(con.SecurityContext.Capabilities.Add, apiv1.Capability(cap))
			}
		} 
	}

	// create from PO
	p, err := CreatePodFromObject(po, &pname, &ns, true)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// ExecRootAccessCmd ...
func ExecRootAccessCmd() (int, error) {
	//make sure the pod is there (want one without privileged access or escalation)
	f := false
	_, err := CreatePODSettingSecurityContext(&f, &f, nil)

	if err != nil {
		return -1, err
	}

	//create a command that requires root access
	//try chroot
	cmd := "chroot ."
	ns, pn := pspTestNamespace, pspTestPodName
	stdout, _, ex, err := ExecCommand(&cmd, &ns, &pn)

	log.Printf("[NOTICE] ExecRootAccessCmd: %v stdout: %v exit code: %v", cmd, stdout, ex)

	if err != nil {
		return ex, err
	}

	return ex, nil
}

//TeardownPodSecurityTestPod ...
func TeardownPodSecurityTestPod(p *string) error {
	ns := pspTestNamespace
	err := DeletePod(p, &ns, true)
	return err
}

func getPodSecurityPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	c, err := GetClient()
	if err != nil {
		return nil, err
	}

	psp := c.PolicyV1beta1().PodSecurityPolicies()
	if psp == nil {
		return nil, fmt.Errorf("Pod Security Polices could not be obtained (nil returned)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pspList, err := psp.List(ctx, metav1.ListOptions{})
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

// KubeSecurityPolicyProvider ...
type KubeSecurityPolicyProvider struct {
	psps     *v1beta1.PodSecurityPolicyList
	pspMutex sync.Mutex
}

func (p *KubeSecurityPolicyProvider) getPolicies() (*v1beta1.PodSecurityPolicyList, error) {
	p.pspMutex.Lock()
	defer p.pspMutex.Unlock()

	//already got them?
	if p.psps == nil {
		ps, err := getPodSecurityPolicies()
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
