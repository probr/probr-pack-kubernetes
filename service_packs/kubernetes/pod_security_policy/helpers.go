package pod_security_policy

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/api/policy/v1beta1"
)

// PrivilegedAccess type enumerating Privileged Access
type PrivilegedAccess int

// PrivilegedAccess enum
const (
	WithPrivilegedAccess PrivilegedAccess = iota
	WithoutPrivilegedAccess
)

// PSPProbeCommand type enumerating the commands that can be used to test pods for compliance with Pod Security Policies
type PSPProbeCommand int

type scenarioState struct {
	name     string
	audit    *audit.ScenarioAudit
	probe    *audit.Probe
	podState kubernetes.PodState
}

// PSPVerificationProbe encapsulates the command and expected result to be used in a Pod Security Policy probe.
type VerificationProbe struct {
	Cmd              PSPProbeCommand
	ExpectedExitCode int
}

// enumn supporting PSPProbeCommand type
const (
	Chroot PSPProbeCommand = iota
	EnterHostPIDNS
	EnterHostIPCNS
	EnterHostNetworkNS
	VerifyNonRootUID
	NetRawProbe
	SpecialCapProbe
	NetCat
	Unshare
	Ls
)

func (c PSPProbeCommand) String() string {
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
	//NOTE: either the above namespace needs to be added to the exclusion list on the
	//container registry image needs to be available in the allowed (probably internal) registry
	defaultPSPProbeContainer = "psp-test"
	defaultPSPProbePodName   = "psp-test-pod"
)

// PodSecurityPolicy interface defines a set of methods to support the testing of Pod Security Policies.
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
	CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64, probe *audit.Probe) (*apiv1.Pod, error)
	CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool, probe *audit.Probe) (*apiv1.Pod, error)
	CreatePODSettingCapabilities(c *[]string, probe *audit.Probe) (*apiv1.Pod, error)
	CreatePodFromYaml(y []byte, probe *audit.Probe) (*apiv1.Pod, error)
	ExecPSPProbeCmd(pName *string, cmd PSPProbeCommand, probe *audit.Probe) (*kubernetes.CmdExecutionResult, error)
	TeardownPodSecurityProbe(p string, e string) error
	CreateConfigMap() error
	DeleteConfigMap() error
}

// PSP implements PodSecurityPolicy.
type PSP struct {
	k                       kubernetes.Kubernetes
	securityPolicyProviders *[]SecurityPolicyProvider

	probeImage     string
	probeContainer string
	probePodName   string
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}

// NewPSP creates a new PSP using the supplied kubernetes instance and collection of SecurityPolicyProviders.
func NewPSP(k kubernetes.Kubernetes, sp *[]SecurityPolicyProvider) *PSP {
	p := &PSP{}
	p.k = k
	p.securityPolicyProviders = sp

	p.setenv()
	return p
}

// NewDefaultPSP creates a new PSP using the default kubernetes instance and the pre-defined SecurityPolicyProviders.
func NewDefaultPSP() *PSP {
	p := &PSP{}
	p.k = kubernetes.GetKubeInstance()

	//standard security providers
	p.securityPolicyProviders = &[]SecurityPolicyProvider{
		NewKubePodSecurityPolicyProvider(p.k),
		NewAzK8sConstraintTemplate(p.k)}

	p.setenv()
	return p

}

func (psp *PSP) setenv() {

	//just default these for now (not sure we'll ever want to supply):
	psp.probeContainer = defaultPSPProbeContainer
	psp.probePodName = defaultPSPProbePodName

	// Extract registry and image info from config
	psp.probeImage = config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry + "/" + config.Vars.ServicePacks.Kubernetes.ProbeImage
}

// ClusterIsDeployed verifies that a suitable kubernetes cluster is deployed.
func (psp *PSP) ClusterIsDeployed() *bool {
	return psp.k.ClusterIsDeployed()
}

// ClusterHasPSP determines if the cluster has any SecurityPolicyProvider's set.
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

// PrivilegedAccessIsRestricted looks for a SecurityPolicyProvider with 'Privileged' set to false (ie. NOT privileged).
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

	return logAndReturn("PrivilegedAccessIsRestricted", success, ret, err)
}

// HostPIDIsRestricted looks for a SecurityPolicyProvider with 'HostPID' set to false (i.e. NO Access to HostPID ).
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

	return logAndReturn("HostPIDIsRestricted", success, ret, err)
}

// HostIPCIsRestricted looks for a SecurityPolicyProvider with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
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

	return logAndReturn("HostIPCIsRestricted", success, ret, err)
}

// HostNetworkIsRestricted looks for a SecurityPolicyProvider with 'HostIPC' set to false (i.e. NO Access to HostIPC ).
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

	return logAndReturn("HostNetworkIsRestricted", success, ret, err)
}

// PrivilegedEscalationIsRestricted looks for a SecurityPolicyProvider with 'Privileged' set to false (ie. NOT privileged).
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

	return logAndReturn("PrivilegedEscalationIsRestricted", success, ret, err)
}

// RootUserIsRestricted looks for a SecurityPolicyProvider which prevents root user access.
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

	return logAndReturn("RootUserIsRestricted", success, ret, err)
}

// NETRawIsRestricted looks for a SecurityPolicyProvider where the NET_RAW capability is restricted.
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

	return logAndReturn("NETRawIsRestricted", success, ret, err)
}

// AllowedCapabilitiesAreRestricted looks for a SecurityPolicyProvider where allowed capabilities are restricted.
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

	return logAndReturn("AllowedCapabilitiesAreRestricted", success, ret, err)
}

// AssignedCapabilitiesAreRestricted looks for a SecurityPolicyProvider where assigned capabilities are restricted.
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

	return logAndReturn("AssignedCapabilitiesAreRestricted", success, ret, err)
}

// HostPortsAreRestricted looks for a SecurityPolicyProvider which has a HostPort restriction.
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

	return logAndReturn("HostPortsAreRestricted", success, ret, err)
}

// VolumeTypesAreRestricted looks for a SecurityPolicyProvider which has a VolumeType restriction.
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

	return logAndReturn("VolumeTypesAreRestricted", success, ret, err)
}

// SeccompProfilesAreRestricted looks for a SecurityPolicyProvider which restricts seccomp profiles.
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

	return logAndReturn("SeccompProfilesAreRestricted", success, ret, err)
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

//logAndReturn is a convenience function for processing the results of a security policy call.
//It logs the test name and result.
//Parameters:
//t - string: test name
//s - bool: success / fail
//r - bool: overall result
//e - error
func logAndReturn(t string, s bool, r bool, e error) (*bool, error) {
	log.Printf("[INFO] Security Policy check: %q.  Overall result: %t. (error: %v)", t, r, e)

	//if we 've had a success, ignore the error ...
	if s {
		//then we've made at least one successful call - nil out err, for client simplification
		return &r, nil
	}

	//otherwise just return
	return &r, e
}

// CreatePODSettingSecurityContext creates POD with a SecurityContext conforming to the parameters:
// pr *bool - set the Privileged flag.  Defaults to false.
// pe *bool - set the Allow Privileged Escalation flag.  Defaults to false.
// runAsUser *int64 - set RunAsUser.  Defaults to 1000.
func (psp *PSP) CreatePODSettingSecurityContext(pr *bool, pe *bool, runAsUser *int64, probe *audit.Probe) (*apiv1.Pod, error) {
	//default sensibly if not provided
	//this needs to take account of rules around allowedPrivilegedEscalation and Privileged:
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

	pname, ns, cname, image := kubernetes.GenerateUniquePodName(psp.probePodName), kubernetes.Namespace, psp.probeContainer, psp.probeImage

	//let caller handle ...
	pod, _, err := psp.k.CreatePod(pname, ns, cname, image, true, &sc, probe)
	return pod, err
}

// CreatePODSettingAttributes creates a POD with attributes:
// hostPID *bool - set the hostPID flag, defaults to false
// hostIPC *bool - set the hostIPC flag, defaults to false
// hostNetwork *bool - set the hostNetwork flag, defaults to false
func (psp *PSP) CreatePODSettingAttributes(hostPID *bool, hostIPC *bool, hostNetwork *bool, probe *audit.Probe) (*apiv1.Pod, error) {
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

	pname, ns, cname, image := kubernetes.GenerateUniquePodName(psp.probePodName), kubernetes.Namespace, psp.probeContainer, psp.probeImage

	// get the pod object and manipulate:
	po := psp.k.GetPodObject(pname, ns, cname, image, nil)
	po.Spec.HostPID = *hostPID
	po.Spec.HostIPC = *hostIPC
	po.Spec.HostNetwork = *hostNetwork

	// create from PO (and let caller handle ...)
	return psp.k.CreatePodFromObject(po, pname, ns, true, probe)
}

// CreatePODSettingCapabilities creates a pod with the supplied capabilities.
func (psp *PSP) CreatePODSettingCapabilities(c *[]string, probe *audit.Probe) (*apiv1.Pod, error) {
	pname, ns, cname, image := kubernetes.GenerateUniquePodName(psp.probePodName), kubernetes.Namespace, psp.probeContainer, psp.probeImage

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
	return psp.k.CreatePodFromObject(po, pname, ns, true, probe)
}

// CreatePodFromYaml creates a pod from the supplied yaml.
func (psp *PSP) CreatePodFromYaml(y []byte, probe *audit.Probe) (*apiv1.Pod, error) {
	pname := kubernetes.GenerateUniquePodName(psp.probePodName)

	return psp.k.CreatePodFromYaml(y, pname, kubernetes.Namespace, psp.probeImage, "", true, probe)
}

// ExecPSPProbeCmd executes the given PSPProbeCommand against the supplied pod name.
func (psp *PSP) ExecPSPProbeCmd(pName *string, cmd PSPProbeCommand, probe *audit.Probe) (*kubernetes.CmdExecutionResult, error) {
	var pn string
	//if we've not been given a pod name, assume one needs to be created:
	if pName == nil {
		//want one without privileged access or escalation
		f := false
		p, err := psp.CreatePODSettingSecurityContext(&f, &f, nil, probe)

		if err != nil {
			return nil, err
		}
		//grab the pod name:
		pn = p.GetObjectMeta().GetName()
	} else {
		pn = *pName
	}

	c := cmd.String()
	res := psp.k.ExecCommand(c, kubernetes.Namespace, &pn)

	log.Printf("[INFO] ExecPSPProbeCmd: %v stdout: %v exit code: %v (error: %v)", cmd, res.Stdout, res.Code, res.Err)

	return res, nil
}

// CreateConfigMap creates a config map to support PSP testing.
func (psp *PSP) CreateConfigMap() error {
	//set up anything that may be required for testing
	//1. config map
	_, err := psp.k.CreateConfigMap(utils.StringPtr("test-config-map"), kubernetes.Namespace)

	if err != nil {
		log.Printf("[NOTICE] Failed to create test config map: %v", err)
	}

	return err
}

// DeleteConfigMap deletes the config map supporting the PSP testing.
func (psp *PSP) DeleteConfigMap() error {
	return psp.k.DeleteConfigMap("test-config-map")
}

// TeardownPodSecurityProbe deletes the given pod name in the PSP test namespace.
func (psp *PSP) TeardownPodSecurityProbe(p string, e string) error {
	err := psp.k.DeletePod(p, kubernetes.Namespace, e) //don't worry about waiting
	return err
}

// KubePodSecurityPolicyProvider implements SecurityPolicyProvider and looks for kubernetes PodSecurityPolices.
type KubePodSecurityPolicyProvider struct {
	k        kubernetes.Kubernetes
	psps     *v1beta1.PodSecurityPolicyList
	pspMutex sync.Mutex
}

// NewKubePodSecurityPolicyProvider creates a new KubePodSecurityPolicyProvider with the supplied kubernetes instance.
func NewKubePodSecurityPolicyProvider(k kubernetes.Kubernetes) *KubePodSecurityPolicyProvider {
	return &KubePodSecurityPolicyProvider{k: k}
}

func (p *KubePodSecurityPolicyProvider) getPolicies() (*v1beta1.PodSecurityPolicyList, error) {
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

// HasSecurityPolicies provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasSecurityPolicies() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	b := len(psps.Items) > 0
	return &b, nil
}

// HasPrivilegedAccessRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasPrivilegedAccessRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have Privileged set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.Privileged {
			log.Printf("[DEBUG] PodSecurityPolicy: Privileged is set to %v on Policy: %v", e.Spec.Privileged, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasHostPIDRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasHostPIDRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostPID set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostPID {
			log.Printf("[DEBUG] PodSecurityPolicy: HostPID is set to %v on Policy: %v\n", e.Spec.HostPID, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasHostIPCRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasHostIPCRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostIPC set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostIPC {
			log.Printf("[DEBUG] PodSecurityPolicy: HostIPC is set to %v on Policy: %v\n", e.Spec.HostIPC, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasHostNetworkRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasHostNetworkRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have HostNetwork set to false
	var res bool
	for _, e := range psps.Items {
		if !e.Spec.HostNetwork {
			log.Printf("[DEBUG] PodSecurityPolicy: HostNetwork is set to %v on Policy: %v\n", e.Spec.HostNetwork, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasAllowPrivilegeEscalationRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have AllowPrivilegeEscalation set to false
	var res bool
	for _, e := range psps.Items {
		if !*e.Spec.AllowPrivilegeEscalation {
			log.Printf("[DEBUG] PodSecurityPolicy: AllowPrivilegeEscalation is set to %v on Policy: %v", e.Spec.AllowPrivilegeEscalation, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasRootUserRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasRootUserRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least on of the PSPs should have AllowPrivilegeEscalation set to false
	var res bool
	for _, e := range psps.Items {
		if e.Spec.RunAsUser.Rule == v1beta1.RunAsUserStrategyMustRunAsNonRoot {
			log.Printf("[DEBUG] PodSecurityPolicy: RunAsUserStrategyMustRunAsNonRoot is set on Policy: %v", e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

// HasNETRAWRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasNETRAWRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//at least one of the PSPs should have a RequiredDropCapability of "NET_RAW"
	var res bool
	for _, e := range psps.Items {
		for _, c := range e.Spec.RequiredDropCapabilities {
			if c == "NET_RAW" || c == "ALL" {
				log.Printf("[DEBUG] PodSecurityPolicy: HasNETRAWRestriction: RequiredDropCapability of %v is set on Policy: %v", c, e.GetName())
				res = true
				break
			}
		}
	}

	return &res, nil
}

// HasAllowedCapabilitiesRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasAllowedCapabilitiesRestriction() (*bool, error) {
	psps, err := p.getPolicies()
	if err != nil {
		return nil, err
	}

	//in this case we don't want "allowedCapabilities" on any PSP (default to true)
	res := true
	for _, e := range psps.Items {
		if e.Spec.AllowedCapabilities != nil && len(e.Spec.AllowedCapabilities) > 0 {
			log.Printf("[DEBUG] PodSecurityPolicy: HasAllowedCapabilitiesRestriction: at least one AllowedCapability is set on Policy: %v", e.GetName())
			res = false
			break
		}
	}

	return &res, nil
}

// HasAssignedCapabilitiesRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasAssignedCapabilitiesRestriction() (*bool, error) {
	//TODO: review - doesn't appear to be a PSP to enforce this
	b := false
	log.Printf("[INFO] PodSecurityPolicy: HasAssignedCapabilitiesRestriction defaulting to %t", b)
	return &b, nil
}

// HasHostPortRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasHostPortRestriction() (*bool, error) {
	//TODO: review this. From one view, this is always true as ports are locked down by
	//default and only opened via the hostport range on a PSP.  Which ports are allowed
	//to be open will be a case by case basis.
	//For now return 'true' as, theoretically, there is a 'host port restriction'
	//(see: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#hostportrange-v1beta1-policy)

	b := true
	log.Printf("[INFO] PodSecurityPolicy: HasHostPortRestriction defaulting to %t", b)
	return &b, nil
}

// HasVolumeTypeRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasVolumeTypeRestriction() (*bool, error) {
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
				log.Printf("[DEBUG] PodSecurityPolicy: HasVolumeTypeRestriction: at least one unapproved volume type (%v) is set on Policy: %v",
					v, e.GetName())
				res = false
				break
			}
		}
	}

	return &res, nil
}

// HasSeccompProfileRestriction provides the KubePodSecurityPolicyProvider implementation of SecurityPolicyProvider.
func (p *KubePodSecurityPolicyProvider) HasSeccompProfileRestriction() (*bool, error) {
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
			log.Printf("[DEBUG] PodSecurityPolicy: HasSeccompProfileRestriction: annotation of %v with value %v is set on Policy: %v",
				a, v, e.GetName())
			res = true
			break
		}
	}

	return &res, nil
}

func (p *KubePodSecurityPolicyProvider) getPodSecurityPolicies() (*v1beta1.PodSecurityPolicyList, error) {
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

	log.Printf("[NOTICE] PodSecurityPolicy: There are %d psp policies in the cluster\n", len(pspList.Items))

	for _, e := range pspList.Items {
		log.Printf("[INFO] PSP: %v \n", e.GetName())
		log.Printf("[INFO] Spec: %+v \n", e.Spec)
	}

	return pspList, nil
}
