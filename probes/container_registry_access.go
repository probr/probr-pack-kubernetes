// Package containerregistryaccess provides the implementation required to execute the
// feature based test cases described in the the 'events' directory.
package probes

import (
	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/utils"
)

const (
	CRA_NAME = "container_registry_access"
)

// init() registers the feature tests descibed in this package with the test runner (coreengine.TestRunner) via the call
// to coreengine.AddTestHandler.  This links the test - described by the TestDescriptor - with the handler to invoke.  In
// this case, the general test handler is being used (probes.GodogTestHandler) and the GodogTest data provides the data
// require to execute the test.  Specifically, the data includes the Test Suite and Scenario Initializers from this package
// which will be called from probes.GodogTestHandler.  Note: a blank import at probr library level should be done to
// invoke this function automatically on initial load.
func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.ContainerRegistryAccess, Name: CRA_NAME}

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: craTestSuiteInitialize,
			ScenarioInitializer:  craScenarioInitialize,
		},
	})
}

// ContainerRegistryAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// container registry probes.
var cra kubernetes.ContainerRegistryAccess

// SetContainerRegistryAccess allows injection of ContainerRegistryAccess helper.
func SetContainerRegistryAccess(c kubernetes.ContainerRegistryAccess) {
	cra = c
}

// TEST STEPS:

// CIS-6.1.3
// Minimize cluster access to read-only
func (p *probeState) iAmAuthorisedToPullFromAContainerRegistry() error {
	pd, err := cra.SetupContainerAccessTestPod(utils.StringPtr("docker.io"))

	s := ProcessPodCreationResult(&p.state, pd, kubernetes.PSPContainerAllowedImages, p.event, err)
	p.event.AuditProbeStep(p.name, s)
	return s
}

// PENDING IMPLEMENTATION
func (p *probeState) iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity() error {
	return godog.ErrPending
}

// PENDING IMPLEMENTATION
func (p *probeState) thePushRequestIsRejectedDueToAuthorization() error {
	return godog.ErrPending
}

// CIS-6.1.4
// Ensure only authorised container registries are allowed
func (p *probeState) aUserAttemptsToDeployAContainerFrom(auth string, registry string) error {
	pd, err := cra.SetupContainerAccessTestPod(&registry)

	e := p.event
	s := ProcessPodCreationResult(&p.state, pd, kubernetes.PSPContainerAllowedImages, e, err)
	e.AuditProbeStep(p.name, s)
	return s
}

func (p *probeState) theDeploymentAttemptIs(res string) error {
	s := AssertResult(&p.state, res, "")
	p.event.AuditProbeStep(p.name, s)
	return s
}

// craTestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func craTestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	//check dependancies ...
	if cra == nil {
		// not been given one so set default
		cra = kubernetes.NewDefaultCRA()
	}
}

// craScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func craScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probeState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.BeforeScenario(CRA_NAME, s)
	})

	//common
	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)

	//CIS-6.1.3
	ctx.Step(`^I am authorised to pull from a container registry$`, ps.iAmAuthorisedToPullFromAContainerRegistry)
	ctx.Step(`^I attempt to push to the container registry using the cluster identity$`, ps.iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity)
	ctx.Step(`^the push request is rejected due to authorization$`, ps.thePushRequestIsRejectedDueToAuthorization)

	//CIS-6.1.4
	ctx.Step(`^a user attempts to deploy a container from "([^"]*)" registry "([^"]*)"$`, ps.aUserAttemptsToDeployAContainerFrom)
	ctx.Step(`^the deployment attempt is "([^"]*)"$`, ps.theDeploymentAttemptIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		cra.TeardownContainerAccessTestPod(&ps.state.PodName, CRA_NAME)

		LogScenarioEnd(s)
	})
}
