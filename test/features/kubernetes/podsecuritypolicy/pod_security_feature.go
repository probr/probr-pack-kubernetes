package podsecuritypolicy

//go:generate go-bindata.exe -pkg $GOPACKAGE -o assets/assets.go assets/yaml

import (
	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
	"gitlab.com/citihub/probr/test/features"
	podsecuritypolicy "gitlab.com/citihub/probr/test/features/kubernetes/podsecuritypolicy/assets"
	apiv1 "k8s.io/api/core/v1"
)

type probState struct {
	podName        string
	creationError  *kubernetes.PodCreationError
	expectedReason *kubernetes.PodCreationErrorReason
}

//TODO: revise when interface this bit up ...
var psp kubernetes.PodSecurityPolicy

// SetPodSecurityPolicy ...
func SetPodSecurityPolicy(p kubernetes.PodSecurityPolicy) {
	psp = p
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
	b := psp.ClusterIsDeployed()

	if b == nil || !*b {
		return features.LogAndReturnError("kubernetes cluster is NOT deployed")
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
			return features.LogAndReturnError("pod %v was created - test failed", p.podName)
		}
		//should also check code:
		_, exists := p.creationError.ReasonCodes[*p.expectedReason]
		if !exists {
			//also a fail:
			return features.LogAndReturnError("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				p.creationError.ReasonCodes, p.expectedReason)
		}

		//we're good
		return nil
	}

	if res == "Succeed" {
		// then expect the pod creation error to be nil
		if p.creationError != nil {
			//it's a fail:
			return features.LogAndReturnError("pod was not created - test failed: %v", p.creationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	return features.LogAndReturnError("desired result %v is not recognised", res)

}

func (p *probState) performAllowedCommand() error {
	//check for lack of creation error, i.e. pod was created successfully
	if p.creationError == nil {
		// execute 'general allowed' command.  This is a simple inocuous cmd, e.g. ls
		c := kubernetes.Ls
		ex, err := psp.ExecPSPTestCmd(&p.podName, c)

		//want this to succeed - no errors & 0 exit code:
		if err != nil && ex == 0 {
			return features.LogAndReturnError("allowed command (%v) failed. Exit code: %v Error: %v", c, ex, err)
		}
		return nil
	}

	//pod wasn't created so nothing to test
	return nil
}

// common helper funcs
func (p *probState) processCreationResult(pd *apiv1.Pod, expected kubernetes.PodCreationErrorReason, err error) error {
	//first check for errors:
	if err != nil {
		//check if we've got a partial pod creation
		//e.g. pod was created but did't get to "running" state
		//in this case we need to hold onto the name so it can be deleted
		if pd != nil {
			p.podName = pd.GetObjectMeta().GetName()
		}

		//check for known error type
		//this means the pod has not been created for an expected reason and
		//is a valid result if the test is addressing prevention of insecure pod creation
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			p.creationError = e
			p.expectedReason = &expected
			return nil
		}
		//unexpected error
		//in this case something unexpected has happened, return an error to cucumber
		return features.LogAndReturnError("error attempting to create POD: %v", err)
	}

	//No errors: pod creation may or may not have been expected.  This will be determined
	//by the specific test case
	if pd == nil {
		// pod not created, which could be valid for some tests
		return nil
	}

	//if we've got this far, a pod was successfully created which could be
	//valid for some tests
	p.podName = pd.GetObjectMeta().GetName()

	//we're good
	return nil
}

func (p *probState) runControlTest(cf func() (*bool, error), c string) error {
	yesNo, err := cf()

	if err != nil {
		return features.LogAndReturnError("error determining Pod Security Policy: %v error: %v", c, err)
	}
	if yesNo == nil {
		return features.LogAndReturnError("result of %v is nil despite no error being raised from the call", c)
	}

	if !*yesNo {
		return features.LogAndReturnError("%v is NOT restricted (result: %t)", c, *yesNo)
	}

	return nil
}

func (p *probState) runVerificationTest(c kubernetes.PSPTestCommand) error {
	//check for lack of creation error, i.e. pod was created successfully
	if p.creationError == nil {
		ex, err := psp.ExecPSPTestCmd(&p.podName, c)
		//want this to fail as execution of a command requiring root should be blocked
		if err != nil {
			return nil
		}
		if err == nil {
			return features.LogAndReturnError("verification command (%v) succeeded with exit code %v", c, ex)
		}
	}

	//pod wasn't created so nothing to test
	return nil
}

// TEST STEPS:

// CIS-5.2.1
// privileged access
func (p *probState) privilegedAccessRequestIsMarkedForTheKubernetesDeployment(privilegedAccessRequested string) error {
	var pa bool
	if privilegedAccessRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := psp.CreatePODSettingSecurityContext(&pa, nil, nil)

	return p.processCreationResult(pd, kubernetes.PSPNoPrivilege, err)
}

func (p *probState) someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster() error {
	return p.runControlTest(psp.PrivilegedAccessIsRestricted, "PrivilegedAccessIsRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatRequiresPrivilegedAccess() error {
	return p.runVerificationTest(kubernetes.Chroot)
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

	pd, err := psp.CreatePODSettingAttributes(&hostPID, nil, nil)

	return p.processCreationResult(pd, kubernetes.PSPHostNamespace, err)
}

func (p *probState) someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster() error {
	return p.runControlTest(psp.HostPIDIsRestricted, "HostPIDIsRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostPIDNamespace() error {
	return p.runVerificationTest(kubernetes.EnterHostPIDNS)
}

//CIS-5.2.3
func (p *probState) hostIPCRequestIsMarkedForTheKubernetesDeployment(hostIPCRequested string) error {
	var hostIPC bool
	if hostIPCRequested == "True" {
		hostIPC = true
	} else {
		hostIPC = false
	}

	pd, err := psp.CreatePODSettingAttributes(nil, &hostIPC, nil)

	return p.processCreationResult(pd, kubernetes.PSPHostNamespace, err)
}

func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.HostIPCIsRestricted, "HostIPCIsRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostIPCNamespace() error {
	return p.runVerificationTest(kubernetes.EnterHostIPCNS)
}

//CIS-5.2.4
func (p *probState) hostNetworkRequestIsMarkedForTheKubernetesDeployment(hostNetworkRequested string) error {
	var hostNetwork bool
	if hostNetworkRequested == "True" {
		hostNetwork = true
	} else {
		hostNetwork = false
	}

	pd, err := psp.CreatePODSettingAttributes(nil, nil, &hostNetwork)

	return p.processCreationResult(pd, kubernetes.PSPHostNetwork, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.HostNetworkIsRestricted, "HostNetworkIsRestricted")
}
func (p *probState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostNetworkNamespace() error {
	return p.runVerificationTest(kubernetes.EnterHostNetworkNS)
}

//CIS-5.2.5
func (p *probState) privilegedEscalationIsMarkedForTheKubernetesDeployment(privilegedEscalationRequested string) error {
	var pa bool
	if privilegedEscalationRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := psp.CreatePODSettingSecurityContext(nil, &pa, nil)

	return p.processCreationResult(pd, kubernetes.PSPNoPrivilegeEscalation, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.PrivilegedEscalationIsRestricted, "PrivilegedEscalationIsRestricted")
}

//"but" same as 5.2.1

//CIS-5.2.6
func (p *probState) theUserRequestedIsForTheKubernetesDeployment(requestedUser string) error {
	var runAsUser int64
	if requestedUser == "Root" {
		runAsUser = 0
	} else {
		runAsUser = 1000
	}

	pd, err := psp.CreatePODSettingSecurityContext(nil, nil, &runAsUser)
	return p.processCreationResult(pd, kubernetes.PSPAllowedUsersGroups, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.RootUserIsRestricted, "RootUserIsRestricted")
}
func (p *probState) theKubernetesDeploymentShouldRunWithANonrootUID() error {
	return p.runVerificationTest(kubernetes.VerifyNonRootUID)
}

//CIS-5.2.7
func (p *probState) nETRAWIsMarkedForTheKubernetesDeployment(netRawRequested string) error {
	var c []string
	if netRawRequested == "True" {
		c = make([]string, 1)
		c[0] = "NET_RAW"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, kubernetes.PSPAllowedCapabilities, err)
}
func (p *probState) someSystemExistsToPreventAKubernetesDeploymentFromRunningWithNETRAWCapabilityInAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.NETRawIsRestricted, "NETRAWIsRestricted")
}
func (p *probState) iShouldNotBeAbleToPerformACommandThatRequiresNETRAWCapability() error {
	return p.runVerificationTest(kubernetes.NetRawTest)
}

//CIS-5.2.8
func (p *probState) additionalCapabilitiesForTheKubernetesDeployment(addCapabilities string) error {
	var c []string
	if addCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		c = make([]string, 1)
		c[0] = "NET_ADMIN"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, kubernetes.PSPAllowedCapabilities, err)
}
func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.AllowedCapabilitiesAreRestricted, "AllowedCapabilitiesAreRestricted")
}
func (p *probState) iShouldNotBeAbleToPerformACommandThatRequiresCapabilitiesOutsideOfTheDefaultSet() error {
	return p.runVerificationTest(kubernetes.SpecialCapTest)
}

//CIS-5.2.9
func (p *probState) assignedCapabilitiesForTheKubernetesDeployment(assignCapabilities string) error {
	var c []string
	if assignCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		//what's the difference with 5.2.8???
		c = make([]string, 1)
		c[0] = "NET_ADMIN"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return p.processCreationResult(pd, kubernetes.PSPAllowedCapabilities, err)
}
func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.AssignedCapabilitiesAreRestricted, "AssignedCapabilitiesAreRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatRequiresAnyCapabilities() error {
	return p.runVerificationTest(kubernetes.SpecialCapTest)
}

//AZ Policy - port range
func (p *probState) anPortRangeIsRequestedForTheKubernetesDeployment(portRange string) error {
	var y []byte
	var err error

	if portRange == "unapproved" {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-hostport-unapproved.yaml")
	} else {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-hostport-approved.yaml")
	}

	if err != nil {
		return err
	}

	//create from yaml
	pd, err := psp.CreatePodFromYaml(y)

	return p.processCreationResult(pd, kubernetes.PSPAllowedPortRange, err)
}

func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedPortRangeFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.HostPortsAreRestricted, "HostPortsAreRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatAccessAnUnapprovedPortRange() error {
	return p.runVerificationTest(kubernetes.NetCat)
}

//AZ Policy - volume type
func (p *probState) anVolumeTypeIsRequestedForTheKubernetesDeployment(volumeType string) error {
	var y []byte
	var err error

	if volumeType == "unapproved" {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-volumetypes-unapproved.yaml")
	} else {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-volumetypes-approved.yaml")
	}

	if err != nil {
		return err
	}

	//create from yaml
	pd, err := psp.CreatePodFromYaml(y)

	return p.processCreationResult(pd, kubernetes.PSPAllowedVolumeTypes, err)
}

func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedVolumeTypesFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.VolumeTypesAreRestricted, "VolumeTypesAreRestricted")
}

func (p *probState) iShouldNotBeAbleToPerformACommandThatAccessesAnUnapprovedVolumeType() error {
	//TODO: Not sure what the test is here - if any
	return nil
}

//AZ Policy - seccomp profile
func (p *probState) anSeccompProfileIsRequestedForTheKubernetesDeployment(seccompProfile string) error {
	var y []byte
	var err error

	if seccompProfile == "unapproved" {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-seccomp-unapproved.yaml")
	} else if seccompProfile == "undefined" {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-seccomp-undefined.yaml")
	} else if seccompProfile == "approved" {
		y, err = podsecuritypolicy.Asset("assets/yaml/psp-azp-seccomp-approved.yaml")
	}

	if err != nil {
		return err
	}

	//create from yaml
	pd, err := psp.CreatePodFromYaml(y)

	return p.processCreationResult(pd, kubernetes.PSPSeccompProfile, err)
}
func (p *probState) someSystemExistsToPreventKubernetesDeploymentsWithoutApprovedSeccompProfilesFromBeingDeployedToAnExistingKubernetesCluster() error {
	return p.runControlTest(psp.SeccompProfilesAreRestricted, "SeccompProfilesAreRestricted")
}
func (p *probState) iShouldNotBeAbleToPerformASystemCallThatIsBlockedByTheSeccompProfile() error {
	return p.runVerificationTest(kubernetes.Unshare)
}

func (p *probState) setup() {
	//just make sure this is reset
	p.podName = ""
	p.creationError = nil
}

func (p *probState) tearDown() {
	psp.TeardownPodSecurityTest(&p.podName)
	p.podName = ""
	p.creationError = nil
}

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		//check dependancies ...
		if psp == nil {
			// not been given one so set default
			psp = kubernetes.NewDefaultPSP()
		}
		psp.CreateConfigMap()
	})

	ctx.AfterSuite(func() {
		psp.DeleteConfigMap()
	})
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.setup()
		features.LogScenarioStart(s)
	})

	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterExistsWhichWeCanDeployInto)

	ctx.Step(`^a Kubernetes deployment is applied to an existing Kubernetes cluster$`, ps.aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster)

	//CIS-5.2.1
	ctx.Step(`^privileged access request is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedAccessRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment running with privileged access in an existing Kubernetes cluster$`, ps.someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires privileged access$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresPrivilegedAccess)

	//CIS-5.2.2
	ctx.Step(`^hostPID request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostPIDRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostPID in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host PID namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostPIDNamespace)

	//CIS-5.2.3
	ctx.Step(`^hostIPC request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostIPCRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostIPC in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host IPC namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostIPCNamespace)

	//CIS-5.2.4
	ctx.Step(`^hostNetwork request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostNetworkRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostNetwork in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host network namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostNetworkNamespace)

	//CIS-5.2.5
	ctx.Step(`^privileged escalation is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedEscalationIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the allowPrivilegeEscalation in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster)

	//CIS-5.2.6
	ctx.Step(`^the user requested is "([^"]*)" for the Kubernetes deployment$`, ps.theUserRequestedIsForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running as the root user in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster)
	ctx.Step(`^the Kubernetes deployment should run with a non-root UID$`, ps.theKubernetesDeploymentShouldRunWithANonrootUID)

	//CIS-5.2.7
	ctx.Step(`^NET_RAW is marked "([^"]*)" for the Kubernetes deployment$`, ps.nETRAWIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running with NET_RAW capability in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningWithNETRAWCapabilityInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires NET_RAW capability$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresNETRAWCapability)

	//CIS-5.2.8
	ctx.Step(`^additional capabilities "([^"]*)" requested for the Kubernetes deployment$`, ps.additionalCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with capabilities beyond the default set from being deployed to an existing kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires capabilities outside of the default set$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresCapabilitiesOutsideOfTheDefaultSet)

	//CIS-5.2.9
	ctx.Step(`^assigned capabilities "([^"]*)" requested for the Kubernetes deployment$`, ps.assignedCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with assigned capabilities from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires any capabilities$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresAnyCapabilities)

	//AZPolicy - port range
	ctx.Step(`^an "([^"]*)" port range is requested for the Kubernetes deployment$`, ps.anPortRangeIsRequestedForTheKubernetesDeployment)
	ctx.Step(`^I should not be able to perform a command that access an unapproved port range$`, ps.iShouldNotBeAbleToPerformACommandThatAccessAnUnapprovedPortRange)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with unapproved port range from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedPortRangeFromBeingDeployedToAnExistingKubernetesCluster)

	//AZPolicy - volume types
	ctx.Step(`^an "([^"]*)" volume type is requested for the Kubernetes deployment$`, ps.anVolumeTypeIsRequestedForTheKubernetesDeployment)
	ctx.Step(`^I should not be able to perform a command that accesses an unapproved volume type$`, ps.iShouldNotBeAbleToPerformACommandThatAccessesAnUnapprovedVolumeType)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with unapproved volume types from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedVolumeTypesFromBeingDeployedToAnExistingKubernetesCluster)

	//AZPolicy - seccomp profile
	ctx.Step(`^an "([^"]*)" seccomp profile is requested for the Kubernetes deployment$`, ps.anSeccompProfileIsRequestedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments without approved seccomp profiles from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithoutApprovedSeccompProfilesFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a system call that is blocked by the seccomp profile$`, ps.iShouldNotBeAbleToPerformASystemCallThatIsBlockedByTheSeccompProfile)

	//general - outcome
	ctx.Step(`^the operation will "([^"]*)" with an error "([^"]*)"$`, ps.theOperationWillWithAnError)
	ctx.Step(`^I should be able to perform an allowed command$`, ps.performAllowedCommand)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ps.tearDown()
		features.LogScenarioEnd(s)
	})

}
