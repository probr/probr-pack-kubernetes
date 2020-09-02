package kubernetes

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.com/citihub/probr/internal/clouddriver/azure"
	"gitlab.com/citihub/probr/internal/config"
	"gitlab.com/citihub/probr/internal/utils"
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
	NetCat
	Unshare
	Ls
)

func (c PSPTestCommand) String() string {
	return [...]string{"chroot .",
		"nsenter -t 1 -p ps",
		"nsenter -t 1 -i ps",
		"nsenter -t 1 -n ps",
		"id -u > 0 ",
		"ping google.com",
		"ip link add dummy0 type dummy",
		"nc -l 1234",
		"unshare",
		"ls"}[c]
}

const (
	//default values.  Overrides can be supplied via the environment.
	defaultPSPTestNamespace = "probr-pod-security-test-ns"
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry rule or busybox need to be available in the allowed (probably internal) registry
	defaultPSPImageRepository = "docker.io"
	defaultPSPTestImage       = "busybox"
	defaultPSPTestContainer   = "psp-test"
	defaultPSPTestPodName     = "psp-test-pod"
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
	HostPortsAreRestricted() (*bool, error)
	VolumeTypesAreRestricted() (*bool, error)
	SeccompProfilesAreRestricted() (*bool, error)
	CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64) (*apiv1.Pod, error)
	CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool) (*apiv1.Pod, error)
	CreatePODSettingCapabilities(c *[]string) (*apiv1.Pod, error)
	CreatePodFromYaml(y []byte) (*apiv1.Pod, error)
	ExecPSPTestCmd(pName *string, cmd PSPTestCommand) (int, error)
	TeardownPodSecurityTest(p *string) error
	CreateConfigMap() error
	DeleteConfigMap() error
}

// PSP ...
type PSP struct {
	k                       Kubernetes
	securityPolicyProviders *[]SecurityPolicyProvider

	testNamespace string
	testImage     string
	testContainer string
	testPodName   string
}

// NewPSP ...
func NewPSP(k Kubernetes, sp *[]SecurityPolicyProvider) *PSP {
	p := &PSP{}
	p.k = k
	p.securityPolicyProviders = sp

	p.setenv()
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

	p.setenv()
	return p

}

func (psp *PSP) setenv() {

	//just default these for now (not sure we'll ever want to supply):
	psp.testNamespace = defaultPSPTestNamespace
	psp.testContainer = defaultPSPTestContainer
	psp.testPodName = defaultPSPTestPodName

	// image repository + busy box from config
	// but default if not supplied
	i := config.Vars.Images.Repository
	if len(i) < 1 {
		i = defaultPSPImageRepository
	}
	b := config.Vars.Images.BusyBox
	if len(b) < 1 {
		b = defaultPSPTestImage
	}

	psp.testImage = i + "/" + b
}

// ClusterIsDeployed ...
func (psp *PSP) ClusterIsDeployed() *bool {
	return psp.k.ClusterIsDeployed()
}

//ClusterHasPSP determines if the cluster has any Pod Security Policies set.
func (psp *PSP) ClusterHasPSP() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasSecurityPolicies, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//PrivilegedAccessIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) PrivilegedAccessIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasPrivilegedAccessRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//HostPIDIsRestricted looks for a PodSecurityPolicy with 'HostPID' set to false (i.e. NO Access to HostPID ).
func (psp *PSP) HostPIDIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasHostPIDRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//HostIPCIsRestricted looks for a PodSecurityPolicy with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
func (psp *PSP) HostIPCIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasHostIPCRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//HostNetworkIsRestricted looks for a PodSecurityPolicy with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
func (psp *PSP) HostNetworkIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasHostNetworkRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//PrivilegedEscalationIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) PrivilegedEscalationIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasAllowPrivilegeEscalationRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

// RootUserIsRestricted ...
func (psp *PSP) RootUserIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasRootUserRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//NETRawIsRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) NETRawIsRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasNETRAWRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//AllowedCapabilitiesAreRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) AllowedCapabilitiesAreRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasAllowedCapabilitiesRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//AssignedCapabilitiesAreRestricted looks for a PodSecurityPolicy with 'Privileged' set to false (ie. NOT privileged).
func (psp *PSP) AssignedCapabilitiesAreRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasAssignedCapabilitiesRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

// HostPortsAreRestricted ...
func (psp *PSP) HostPortsAreRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasHostPortRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

// VolumeTypesAreRestricted ...
func (psp *PSP) VolumeTypesAreRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasVolumeTypeRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

// SeccompProfilesAreRestricted ...
func (psp *PSP) SeccompProfilesAreRestricted() (*bool, error) {
	var err error
	var ret, success bool

	// iterate over providers ...
	for _, p := range *psp.securityPolicyProviders {
		if p == nil {
			continue
		}
		if makeSecurityPolicyCall(p.HasSeccompProfileRestriction, &ret, &success, &err) {
			break
		}
	}

	//if we 've had a success, ignore the error ...
	if success {
		//then we've made at least one successful call - nil out err, for client simplification
		return &ret, nil
	}

	//otherwise just return
	return &ret, err
}

//convenience func to call the supplied 'SecurityPolicy' call and manage the results.
//Expects two bool pointers, one to track overall result and one to capture if any successful call has been made.
//Also requires an error pointer which will be updated based on supplied func call.
//The returned bool informs the caller on whether or not to break the loop, "true" indicating the loop can be broken.
func makeSecurityPolicyCall(f func() (*bool, error), b *bool, s *bool, e *error) bool {

	res, err := f()

	if err != nil {
		//hold onto the error
		*e = err
		//return false to the caller so loop will be continued
		return false
	}
	if res != nil {
		//set the overall result
		*b = *res

		//if we've got a result (irrespective of true/false), then we've had a successful, i.e. non-error call
		//update the success flag to capture
		//(note: this shouldn't be updated on futher errors - at least one success is all that's needed)
		*s = true

		if *res {
			//if true, we've got a successful result and loop can break
			return true
		}
	}

	return false
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

	pname, ns, cname, image := GenerateUniquePodName(psp.testPodName), psp.testNamespace, psp.testContainer, psp.testImage

	//let caller handle ...
	return psp.k.CreatePod(&pname, &ns, &cname, &image, true, &sc)
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

	pname, ns, cname, image := GenerateUniquePodName(psp.testPodName), psp.testNamespace, psp.testContainer, psp.testImage

	// get the pod object and manipulate:
	po := psp.k.GetPodObject(pname, ns, cname, image, nil)
	po.Spec.HostPID = *hostPID
	po.Spec.HostIPC = *hostIPC
	po.Spec.HostNetwork = *hostNetwork

	// create from PO (and let caller handle ...)
	return psp.k.CreatePodFromObject(po, &pname, &ns, true)
}

//CreatePODSettingCapabilities ...
func (psp *PSP) CreatePODSettingCapabilities(c *[]string) (*apiv1.Pod, error) {
	pname, ns, cname, image := GenerateUniquePodName(psp.testPodName), psp.testNamespace, psp.testContainer, psp.testImage

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

	// create from PO (and let caller handle ...)
	return psp.k.CreatePodFromObject(po, &pname, &ns, true)
}

// CreatePodFromYaml ...
func (psp *PSP) CreatePodFromYaml(y []byte) (*apiv1.Pod, error) {
	pname := GenerateUniquePodName(psp.testPodName)

	return psp.k.CreatePodFromYaml(y, &pname, utils.StringPtr(psp.testNamespace),
		utils.StringPtr(psp.testImage), true)
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
	ns := psp.testNamespace
	stdout, _, ex, err := psp.k.ExecCommand(&c, &ns, &pn)

	log.Printf("[NOTICE] ExecPSPTestCmd: %v stdout: %v exit code: %v (error: %v)", cmd, stdout, ex, err)

	if err != nil {
		return ex, err
	}

	return ex, nil
}

// CreateConfigMap ...
func (psp *PSP) CreateConfigMap() error {
	//set up anything that may be required for testing
	//1. config map
	_, err := psp.k.CreateConfigMap(utils.StringPtr("test-config-map"), &psp.testNamespace)

	if err != nil {
		log.Printf("[NOTICE] Failed to create test config map: %v", err)
	}

	return err
}

// DeleteConfigMap ...
func (psp *PSP) DeleteConfigMap() error {
	return psp.k.DeleteConfigMap(utils.StringPtr("test-config-map"), &psp.testNamespace)
}

//TeardownPodSecurityTest ...
func (psp *PSP) TeardownPodSecurityTest(p *string) error {
	ns := psp.testNamespace
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
			log.Printf("[NOTICE] Privileged is set to %v on Policy: %v", e.Spec.Privileged, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with Privileged set.\n")
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
			log.Printf("[NOTICE] HostPID is set to %v on Policy: %v\n", e.Spec.HostPID, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with HostPID set.\n")
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
			log.Printf("[NOTICE] HostIPC is set to %v on Policy: %v\n", e.Spec.HostIPC, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with HostIPC set.\n")
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
			log.Printf("[NOTICE] HostNetwork is set to %v on Policy: %v\n", e.Spec.HostNetwork, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with HostNetwork set.\n")
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
			log.Printf("[NOTICE] AllowPrivilegeEscalation is set to %v on Policy: %v", e.Spec.AllowPrivilegeEscalation, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with AllowPrivilegeEscalation set.\n")
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
			log.Printf("[NOTICE] RunAsUserStrategyMustRunAsNonRoot is set on Policy: %v", e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] NO Policies found with AllowPrivilegeEscalation set.\n")
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
				log.Printf("[NOTICE] HasNETRAWRestriction: RequiredDropCapability of %v is set on Policy: %v", c, e.GetName())
				res = true
				break
			}
		}
	}

	if !res {
		log.Printf("[NOTICE] HasNETRAWRestriction: NO Policies found with RequiredDropCapability of NET_RAW set.")
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
			log.Printf("[NOTICE] HasAllowedCapabilitiesRestriction: at least one AllowedCapability is set on Policy: %v", e.GetName())
			res = false
			break
		}
	}

	if res {
		log.Printf("[NOTICE] HasNETRAWRestriction: NO Policies found with AllowedCapabilities")
	}

	return &res, nil
}

// HasAssignedCapabilitiesRestriction ...
func (p *KubeSecurityPolicyProvider) HasAssignedCapabilitiesRestriction() (*bool, error) {
	return utils.BoolPtr(false), nil
}

// HasHostPortRestriction ...
func (p *KubeSecurityPolicyProvider) HasHostPortRestriction() (*bool, error) {
	//TODO: review this. From one view, this is always true as ports are locked down by
	//default and only opened via the hostport range on a PSP.  Which ports are allowed
	//to be open will be a case by case basis.
	//For now return 'true' as, theoretically, there is a 'host port restriction'
	//(see: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#hostportrange-v1beta1-policy)

	return utils.BoolPtr(true), nil
}

// HasVolumeTypeRestriction ...
func (p *KubeSecurityPolicyProvider) HasVolumeTypeRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//only want allowed volumes
	//TODO: this could be use case specific, so we may have to inject the 'good' volumes
	// allowed volume types are configMap, emptyDir, projected, downwardAPI, persistentVolumeClaim
	g := make(map[string]*string, 10)
	g["configMap"] = nil
	g["emptyDir"] = nil
	g["projected"] = nil
	g["downwardAPI"] = nil
	g["persistentVolumeClaim"] = nil

	res := true
	for _, e := range psps.Items {
		for _, v := range e.Spec.Volumes {
			_, exists := g[string(v)]
			if !exists {
				log.Printf("[NOTICE] HasVolumeTypeRestriction: at least one unapproved volume type (%v) is set on Policy: %v",
					v, e.GetName())
				res = false
				break
			}
		}
	}

	if res {
		log.Printf("[NOTICE] HasVolumeTypeRestriction: NO Policies found with unapproved volume types")
	}

	return &res, nil
}

// HasSeccompProfileRestriction ...
func (p *KubeSecurityPolicyProvider) HasSeccompProfileRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least one of the PSPs should have an annotation of "seccomp.security.alpha.kubernetes.io/allowedProfileNames"
	a := "seccomp.security.alpha.kubernetes.io/allowedProfileNames"
	var res bool
	for _, e := range psps.Items {
		v, exists := e.Annotations[a]
		if exists {
			log.Printf("[NOTICE] HasSeccompProfileRestriction: annotation of %v with value %v is set on Policy: %v",
				a, v, e.GetName())
			res = true
			break
		}
	}

	if !res {
		log.Printf("[NOTICE] HasSeccompProfileRestriction: NO Policies found with annotation %v set.", a)
	}

	return &res, nil
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
