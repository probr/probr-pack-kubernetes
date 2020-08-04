package podsecuritypolicy

import (
	"fmt"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"citihub.com/probr/internal/coreengine"
	"citihub.com/probr/test/features"
	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"
)

type probState struct {
	podName       string
	creationError *kubernetes.PodCreationError
}

func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.PodSecurityPolicies, Name: "pod_security_policy"}

	coreengine.TestHandleFunc(td, &coreengine.GoDogTestTuple{
		Handler: features.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
		},
	})
}

// general statements.  Cluster exists, etc. Also result/outcome

func (p *probState) creationWillWithAMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

func (p *probState) aKubernetesClusterExistsWhichWeCanDeployInto() error {
	c, err := kubernetes.GetClient()
	if err != nil {
		return err
	}

	if c == nil {
		return fmt.Errorf("client is nil")
	}

	//else we're good ...
	return nil
}

func (p *probState) aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster() error {
	//TODO: not sure this step is adding value ... return "pass" for now ...
	return nil
}

func (p *probState) theOperationWillWithAnError(res, msg string) error {
	if res == "Fail" {
		//expect pod creation error to be non-null
		if p.creationError == nil {
			//it's a fail:
			return fmt.Errorf("pod %v was created - test failed", p.podName)
		}
		//should also check code:
		if p.creationError.ReasonCode != kubernetes.PSPNoPrivilege {
			//also a fail:
			return fmt.Errorf("pod not was created but reason (%v) was not as expected (%v)- test failed",
				p.creationError.ReasonCode, kubernetes.PSPNoPrivilege)
		}

		//we're good
		return nil
	}

	if res == "Succeed" {
		// then expect the pod creation error to be nil
		if p.creationError != nil {
			//it's a fail:
			return fmt.Errorf("pod was not created - test failed: %v", p.creationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	return fmt.Errorf("desired result %v is not recognised", res)

}

func (p *probState) processCreationResult(pd *apiv1.Pod, err error) error {
	if err != nil {
		//check for expected error
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			p.creationError = e
			return nil
		}
		//unexpected error
		return fmt.Errorf("error attempting to create POD: %v", err)
	}

	if pd == nil {
		// valid if the request was for privileged (i.e. pod creation should fail)
		return nil
	}

	//hold on to the pod name
	p.podName = pd.GetObjectMeta().GetName()

	//we're good
	return nil
}

func (p *probState) runControlTest(cf func() (*bool, error), c string) error {
	yesNo, err := cf()

	if err != nil {
		return fmt.Errorf("error determining Pod Security Policy: %v error: %v", c, err)
	}
	if yesNo == nil {
		return fmt.Errorf("result of %v is nil despite no error being raised from the call", c)
	}

	if !*yesNo {
		return fmt.Errorf("%v is NOT restricted (result: %t)", c, *yesNo)
	}

	return nil
}

// CIS-5.2.1
// privileged access
func (p *probState) privilegedAccessRequestIsMarkedForTheKubernetesDeployment(privilegedAccessRequested string) error {
	var pa bool
	if privilegedAccessRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := kubernetes.CreatePODSettingSecurityContext(&pa,nil,nil)

	return p.processCreationResult(pd, err)	
}

func (p *probState) someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster() error {
	return p.runControlTest(kubernetes.PrivilegedAccessIsRestricted, "PrivilegedAccessIsRestricted")
}

// CIS-5.2.2
// hostPID
func (p *probState) hostPIDRequestIsMarkedForTheKubernetesDeployment(hostPIDRequested string) error {
	var hostPID bool
	if hostPIDRequested == "True" {
		hostPID = true
	} else {
		hostPID = false
	}

	pd, err := kubernetes.CreatePODSettingAttributes(&hostPID, nil, nil)

	return p.processCreationResult(pd, err)
}

func (p *probState) someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster() error {	
	return p.runControlTest(kubernetes.HostPIDIsRestricted, "HostPIDIsRestricted")
}

//CIS-5.2.3
func (p *probState) hostIPCRequestIsMarkedForTheKubernetesDeployment(hostIPCRequested string) error {
	var hostIPC bool
	if hostIPCRequested == "True" {
		hostIPC = true
	} else {
		hostIPC = false
	}

	pd, err := kubernetes.CreatePODSettingAttributes(nil, &hostIPC, nil)

	return p.processCreationResult(pd, err)
}

func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.HostIPCIsRestricted, "HostIPCIsRestricted")
}

//CIS-5.2.4
func (p *probState) hostNetworkRequestIsMarkedForTheKubernetesDeployment(hostNetworkRequested string) error {
	var hostNetwork bool
	if hostNetworkRequested == "True" {
		hostNetwork = true
	} else {
		hostNetwork = false
	}

	pd, err := kubernetes.CreatePODSettingAttributes(nil, nil, &hostNetwork)

	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.HostNetworkIsRestricted, "HostNetworkIsRestricted")
}

//CIS-5.2.5
func (p *probState) privilegedEscalationIsMarkedForTheKubernetesDeployment(privilegedEscalationRequested string) error {
	var pa bool
	if privilegedEscalationRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := kubernetes.CreatePODSettingSecurityContext(nil,&pa,nil)

	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.PrivilegedEscalationIsRestricted, "PrivilegedEscalationIsRestricted")
}

//CIS-5.2.6
func (p *probState) theUserRequestedIsForTheKubernetesDeployment(requestedUser string) error {
	var runAsUser int64
	if requestedUser == "Root" {
		runAsUser = 0
	} else {
		runAsUser = 1000
	}

	pd, err := kubernetes.CreatePODSettingSecurityContext(nil, nil, &runAsUser)
	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.RootUserIsRestricted, "RootUserIsRestricted")		
}

//CIS-5.2.7
func (p *probState) nETRAWIsMarkedForTheKubernetesDeployment(netRawRequested string) error {
	var c = make([]string, 1)
	if netRawRequested == "True" {
		c[0] = "NET_RAW"
	}

	pd, err := kubernetes.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningWithNETRAWCapabilityInAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.NETRawIsRestricted, "NETRAWIsRestricted")
}

//CIS-5.2.8
func (p *probState) additionalCapabilitiesForTheKubernetesDeployment(addCapabilities string) error {
	var c = make([]string, 1)
	if addCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		c[0] = "NET_ADMIN"
	}

	pd, err := kubernetes.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.AllowedCapabilitiesAreRestricted, "AllowedCapabilitiesAreRestricted")
}

//CIS-5.2.9
func (p *probState) assignedCapabilitiesForTheKubernetesDeployment(assignCapabilities string) error {
	var c = make([]string, 1)
	if assignCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		//what's the difference with 5.2.8???
		c[0] = "NET_ADMIN"
	}

	pd, err := kubernetes.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, err)
}
func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(kubernetes.AssignedCapabilitiesAreRestricted, "AssignedCapabilitiesAreRestricted")
}

func (p *probState) setup() {
	//just make sure this is reset
	p.podName = ""
	p.creationError = nil
}

func (p *probState) tearDown() {
	kubernetes.TeardownPodSecurityTestPod(&p.podName)
}

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probState{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		ps.setup()
	})

	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterExistsWhichWeCanDeployInto)

	ctx.Step(`^a Kubernetes deployment is applied to an existing Kubernetes cluster$`, ps.aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster)

	//CIS-5.2.1
	ctx.Step(`^privileged access request is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedAccessRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment running with privileged access in an existing Kubernetes cluster$`, ps.someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster)

	//CIS-5.2.2
	ctx.Step(`^hostPID request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostPIDRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostPID in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster)

	//CIS-5.2.3
	ctx.Step(`^hostIPC request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostIPCRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostIPC in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster)

	//CIS-5.2.4
	ctx.Step(`^hostNetwork request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostNetworkRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostNetwork in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster)

	//CIS-5.2.5
	ctx.Step(`^privileged escalation is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedEscalationIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the allowPrivilegeEscalation in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster)

	//CIS-5.2.6
	ctx.Step(`^the user requested is "([^"]*)" for the Kubernetes deployment$`, ps.theUserRequestedIsForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running as the root user in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster)

	//CIS-5.2.7
	ctx.Step(`^NET_RAW is marked "([^"]*)" for the Kubernetes deployment$`, ps.nETRAWIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running with NET_RAW capability in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningWithNETRAWCapabilityInAnExistingKubernetesCluster)

	//CIS-5.2.8
	ctx.Step(`^additional capabilities "([^"]*)" for the Kubernetes deployment$`, ps.additionalCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with capabilities beyond the default set from being deployed to an existing kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster)

	//CIS-5.2.9
	ctx.Step(`^assigned capabilities "([^"]*)" for the Kubernetes deployment$`, ps.assignedCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with assigned capabilities from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster)

	//general - outcome
	ctx.Step(`^the operation will "([^"]*)" with an error "([^"]*)"$`, ps.theOperationWillWithAnError)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		ps.tearDown()
	})

}
