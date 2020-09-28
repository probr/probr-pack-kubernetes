// Package podsecuritypolicy provides the implementation required to execute the feature based test cases described in the
// the 'features' directory.  The 'assets' directory holds any assets required for the test cases.   Assets are 'embedded'
// via the 'go-bindata.exe' tool which is invoked via the 'go generate' tool.  It is important, therefore, that the
//'go:generate' comment is present in order to include this package in the scope of the 'go generate' tool.  This can be
// invoked directly on the command line of via the Makefile (e.g. make clean-build).
package podsecuritypolicy

//go:generate go-bindata.exe -pkg $GOPACKAGE -o assets/assets.go assets/yaml

import (
	"log"

	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/audit"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
	"gitlab.com/citihub/probr/probes"
	podsecuritypolicy "gitlab.com/citihub/probr/probes/kubernetes/podsecuritypolicy/assets"
	"gitlab.com/citihub/probr/probes/kubernetes/probe"
)

type probeState struct {
	state probe.State
}

const NAME = "pod_security_policy"

// PodSecurityPolicy is the section of the kubernetes package which provides the kubernetes interactions required to support
// pod security policy probes.
var psp kubernetes.PodSecurityPolicy

// SetPodSecurityPolicy allows injection of a specific PodSecurityPolicy helper.
func SetPodSecurityPolicy(p kubernetes.PodSecurityPolicy) {
	psp = p
}

// init() registers the feature tests descibed in this package with the test runner (coreengine.TestRunner) via the call
// to coreengine.AddTestHandler.  This links the test - described by the TestDescriptor - with the handler to invoke.  In
// this case, the general test handler is being used (probes.GodogTestHandler) and the GodogTest data provides the data
// require to execute the test.  Specifically, the data includes the Test Suite and Scenario Initializers from this package
// which will be called from probes.GodogTestHandler.  Note: a blank import at probr library level should be done to
// invoke this function automatically on initial load.
func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.PodSecurityPolicies, Name: NAME}

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: probes.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
		},
	})
}

// general statements.  Cluster exists, etc. Also result/outcome

func (p *probeState) creationWillWithAMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

func (p *probeState) aKubernetesClusterExistsWhichWeCanDeployInto() error {

	b := psp.ClusterIsDeployed()

	if b == nil || !*b {
		log.Fatalf("[ERROR] Kubernetes cluster is not deployed")
	}

	//else we're good ...
	return nil
}

func (p *probeState) aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster() error {

	//TODO: not sure this step is adding value ... return "pass" for now ...
	return nil
}

func (p *probeState) theOperationWillWithAnError(res, msg string) error {
	return probe.AssertResult(&p.state, res, msg)
}

func (p *probeState) performAllowedCommand() error {
	//'0' exit code as we expect this to succeed
	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.Ls, ExpectedExitCode: 0})
}

// common helper funcs
func (p *probeState) runControlTest(cf func() (*bool, error), c string) error {

	yesNo, err := cf()

	if err != nil {
		return probes.LogAndReturnError("error determining Pod Security Policy: %v error: %v", c, err)
	}
	if yesNo == nil {
		return probes.LogAndReturnError("result of %v is nil despite no error being raised from the call", c)
	}

	if !*yesNo {
		return probes.LogAndReturnError("%v is NOT restricted (result: %t)", c, *yesNo)
	}

	return nil
}

func (p *probeState) runVerificationTest(c kubernetes.PSPVerificationProbe) error {

	//check for lack of creation error, i.e. pod was created successfully
	if p.state.CreationError == nil {
		res, err := psp.ExecPSPTestCmd(&p.state.PodName, c.Cmd)

		//analyse the results
		if err != nil {
			//this is an error from trying to execute the command as opposed to
			//the command itself returning an error
			return probes.LogAndReturnError("error raised trying to execute verification command (%v) - %v", c.Cmd, err)
		}
		if res == nil {
			return probes.LogAndReturnError("<nil> result received when trying to execute verification command (%v)", c.Cmd)
		}
		if res.Err != nil && res.Internal {
			//we have an error which was raised before reaching the cluster (i.e. it's "internal")
			//this indicates that the command was not successfully executed
			return probes.LogAndReturnError("error raised trying to execute verification command (%v)", c.Cmd)
		}

		//we've managed to execution against the cluster.  This may have failed due to pod security, but this
		//is still a 'sucessful' execution.  The exit code of the command needs to be verified against expected
		//check the result against expected:
		if res.Code == c.ExpectedExitCode {
			//then as expected, test passes
			return nil
		}
		//else it's a fail:
		return probes.LogAndReturnError("exit code %d from verification commnad %q did not match expected %d",
			res.Code, c.Cmd, c.ExpectedExitCode)
	}

	//pod wasn't created so nothing to test
	//TODO: really, we don't want to 'pass' this.  Is there an alternative?
	return nil
}

// TEST STEPS:

// CIS-5.2.1
// privileged access
func (p *probeState) privilegedAccessRequestIsMarkedForTheKubernetesDeployment(privilegedAccessRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var pa bool
	if privilegedAccessRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := psp.CreatePODSettingSecurityContext(&pa, nil, nil)

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPNoPrivilege, e, err)
}

func (p *probeState) someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster() error {

	return p.runControlTest(psp.PrivilegedAccessIsRestricted, "PrivilegedAccessIsRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatRequiresPrivilegedAccess() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.Chroot, ExpectedExitCode: 1})
}

// CIS-5.2.2
// hostPID
func (p *probeState) hostPIDRequestIsMarkedForTheKubernetesDeployment(hostPIDRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var hostPID bool
	if hostPIDRequested == "True" {
		hostPID = true
	} else {
		hostPID = false
	}

	pd, err := psp.CreatePODSettingAttributes(&hostPID, nil, nil)

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPHostNamespace, e, err)
}

func (p *probeState) someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster() error {

	return p.runControlTest(psp.HostPIDIsRestricted, "HostPIDIsRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostPIDNamespace() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.EnterHostPIDNS, ExpectedExitCode: 1})
}

//CIS-5.2.3
func (p *probeState) hostIPCRequestIsMarkedForTheKubernetesDeployment(hostIPCRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var hostIPC bool
	if hostIPCRequested == "True" {
		hostIPC = true
	} else {
		hostIPC = false
	}

	pd, err := psp.CreatePODSettingAttributes(nil, &hostIPC, nil)

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPHostNamespace, e, err)
}

func (p *probeState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.HostIPCIsRestricted, "HostIPCIsRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostIPCNamespace() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.EnterHostIPCNS, ExpectedExitCode: 1})
}

//CIS-5.2.4
func (p *probeState) hostNetworkRequestIsMarkedForTheKubernetesDeployment(hostNetworkRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var hostNetwork bool
	if hostNetworkRequested == "True" {
		hostNetwork = true
	} else {
		hostNetwork = false
	}

	pd, err := psp.CreatePODSettingAttributes(nil, nil, &hostNetwork)

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPHostNetwork, e, err)
}
func (p *probeState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.HostNetworkIsRestricted, "HostNetworkIsRestricted")
}
func (p *probeState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostNetworkNamespace() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.EnterHostNetworkNS, ExpectedExitCode: 1})
}

//CIS-5.2.5
func (p *probeState) privilegedEscalationIsMarkedForTheKubernetesDeployment(privilegedEscalationRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var pa bool
	if privilegedEscalationRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	pd, err := psp.CreatePODSettingSecurityContext(nil, &pa, nil)

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPNoPrivilegeEscalation, e, err)
}
func (p *probeState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.PrivilegedEscalationIsRestricted, "PrivilegedEscalationIsRestricted")
}

//"but" same as 5.2.1

//CIS-5.2.6
func (p *probeState) theUserRequestedIsForTheKubernetesDeployment(requestedUser string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var runAsUser int64
	if requestedUser == "Root" {
		runAsUser = 0
	} else {
		runAsUser = 1000
	}

	pd, err := psp.CreatePODSettingSecurityContext(nil, nil, &runAsUser)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedUsersGroups, e, err)
}
func (p *probeState) someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.RootUserIsRestricted, "RootUserIsRestricted")
}
func (p *probeState) theKubernetesDeploymentShouldRunWithANonrootUID() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.VerifyNonRootUID, ExpectedExitCode: 1})
}

//CIS-5.2.7
func (p *probeState) nETRAWIsMarkedForTheKubernetesDeployment(netRawRequested string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var c []string
	if netRawRequested == "True" {
		c = make([]string, 1)
		c[0] = "NET_RAW"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedCapabilities, e, err)
}
func (p *probeState) someSystemExistsToPreventAKubernetesDeploymentFromRunningWithNETRAWCapabilityInAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.NETRawIsRestricted, "NETRAWIsRestricted")
}
func (p *probeState) iShouldNotBeAbleToPerformACommandThatRequiresNETRAWCapability() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.NetRawTest, ExpectedExitCode: 1})
}

//CIS-5.2.8
func (p *probeState) additionalCapabilitiesForTheKubernetesDeployment(addCapabilities string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var c []string
	if addCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		c = make([]string, 1)
		c[0] = "NET_ADMIN"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedCapabilities, e, err)
}
func (p *probeState) someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.AllowedCapabilitiesAreRestricted, "AllowedCapabilitiesAreRestricted")
}
func (p *probeState) iShouldNotBeAbleToPerformACommandThatRequiresCapabilitiesOutsideOfTheDefaultSet() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.SpecialCapTest, ExpectedExitCode: 2})
}

//CIS-5.2.9
func (p *probeState) assignedCapabilitiesForTheKubernetesDeployment(assignCapabilities string) error {
	e := audit.AuditLog.GetEventLog(NAME)

	var c []string
	if assignCapabilities == "ARE" {
		//TODO: just add net_admin for now - but is this appropriate?
		//what's the difference with 5.2.8???
		c = make([]string, 1)
		c[0] = "NET_ADMIN"
	}

	pd, err := psp.CreatePODSettingCapabilities(&c)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedCapabilities, e, err)
}
func (p *probeState) someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.AssignedCapabilitiesAreRestricted, "AssignedCapabilitiesAreRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatRequiresAnyCapabilities() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.SpecialCapTest, ExpectedExitCode: 2})
}

//AZ Policy - port range
func (p *probeState) anPortRangeIsRequestedForTheKubernetesDeployment(portRange string) error {
	e := audit.AuditLog.GetEventLog(NAME)

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

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedPortRange, e, err)
}

func (p *probeState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedPortRangeFromBeingDeployedToAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.HostPortsAreRestricted, "HostPortsAreRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatAccessAnUnapprovedPortRange() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.NetCat, ExpectedExitCode: 1})
}

//AZ Policy - volume type
func (p *probeState) anVolumeTypeIsRequestedForTheKubernetesDeployment(volumeType string) error {
	e := audit.AuditLog.GetEventLog(NAME)

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

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPAllowedVolumeTypes, e, err)
}

func (p *probeState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedVolumeTypesFromBeingDeployedToAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.VolumeTypesAreRestricted, "VolumeTypesAreRestricted")
}

func (p *probeState) iShouldNotBeAbleToPerformACommandThatAccessesAnUnapprovedVolumeType() error {

	//TODO: Not sure what the test is here - if any
	return nil
}

//AZ Policy - seccomp profile
func (p *probeState) anSeccompProfileIsRequestedForTheKubernetesDeployment(seccompProfile string) error {
	e := audit.AuditLog.GetEventLog(NAME)

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

	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.PSPSeccompProfile, e, err)
}
func (p *probeState) someSystemExistsToPreventKubernetesDeploymentsWithoutApprovedSeccompProfilesFromBeingDeployedToAnExistingKubernetesCluster() error {

	return p.runControlTest(psp.SeccompProfilesAreRestricted, "SeccompProfilesAreRestricted")
}
func (p *probeState) iShouldNotBeAbleToPerformASystemCallThatIsBlockedByTheSeccompProfile() error {

	return p.runVerificationTest(kubernetes.PSPVerificationProbe{Cmd: kubernetes.Unshare, ExpectedExitCode: 1})
}

func (p *probeState) setup() {

	//just make sure this is reset
	p.state.PodName = ""
	p.state.CreationError = nil
}

func (p *probeState) tearDown() {

	psp.TeardownPodSecurityTest(&p.state.PodName, NAME)
	p.state.PodName = ""
	p.state.CreationError = nil
}

// TestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
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

// ScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the features directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probeState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.setup()
		probes.LogScenarioStart(s)
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
		probes.LogScenarioEnd(s)
	})

}
