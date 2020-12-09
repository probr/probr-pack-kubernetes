// Package container_registry_access provides the implementation required to execute the
// feature based test cases described in the the 'events' directory.
package container_registry_access

import (
	"fmt"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
)

type ProbeStruct struct{}

var Probe ProbeStruct

// ContainerRegistryAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// container registry scenarios.
var cra ContainerRegistryAccess

// SetContainerRegistryAccess allows injection of ContainerRegistryAccess helper.
func SetContainerRegistryAccess(c ContainerRegistryAccess) {
	cra = c
}

// TEST STEPS:

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload := kubernetes.ClusterIsDeployed()
	s.audit.AuditScenarioStep(description, payload, nil)
	return nil // ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

// CIS-6.1.3
// Minimize cluster access to read-only
func (s *scenarioState) iAmAuthorisedToPullFromAContainerRegistry() error {
	pod, podAudit, err := cra.SetupContainerAccessProbePod(config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry, s.probe)

	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.PSPContainerAllowedImages, err)

	description := fmt.Sprintf("Creates a new pod using an image from %s. Passes if image successfully pulls and pod is built.", config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry)
	payload := kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

// PENDING IMPLEMENTATION
func (s *scenarioState) iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity() error {
	return godog.ErrPending
}

// PENDING IMPLEMENTATION
func (s *scenarioState) thePushRequestIsRejectedDueToAuthorization() error {
	return godog.ErrPending
}

// CIS-6.1.4
// Ensure deployment from unauthorised container registries is denied
func (s *scenarioState) aUserAttemptsToDeployAuthorisedContainer() error {
	pod, podAudit, err := cra.SetupContainerAccessProbePod(config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry, s.probe)

	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.PSPContainerAllowedImages, err)

	description := fmt.Sprintf("Attempts to deploy a container from %s. Retains pod creation result in scenario state. Passes so long as user is authorized to deploy containers.", config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry)
	payload := kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) theDeploymentAttemptIsAllowed() error {
	err := kubernetes.AssertResult(&s.podState, "allowed", "")

	description := fmt.Sprintf("Asserts pod creation result in scenario state is denied.")
	s.audit.AuditScenarioStep(description, nil, err)

	return err
}

// CIS-6.1.5
// Ensure deployment from authorised container registries is allowed
func (s *scenarioState) aUserAttemptsToDeployUnauthorisedContainer() error {
	pod, podAudit, err := cra.SetupContainerAccessProbePod(config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry, s.probe)

	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.PSPContainerAllowedImages, err)

	description := fmt.Sprintf("Attempts to deploy a container from %s. Retains pod creation result in scenario state. Passes so long as user is authorized to deploy containers.", config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry)
	payload := kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) theDeploymentAttemptIsDenied() error {
	err := kubernetes.AssertResult(&s.podState, "denied", "")

	description := fmt.Sprintf("Asserts pod creation result in scenario state is denied.")
	s.audit.AuditScenarioStep(description, nil, err)

	return err
}

func (p ProbeStruct) Name() string {
	return "container_registry_access"
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	//check dependencies ...
	if cra == nil {
		// not been given one so set default
		cra = NewDefaultCRA()
	}
}

// craScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	//common
	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)

	//CIS-6.1.3
	ctx.Step(`^I am authorised to pull from a container registry$`, ps.iAmAuthorisedToPullFromAContainerRegistry)
	ctx.Step(`^I attempt to push to the container registry using the cluster identity$`, ps.iAttemptToPushToTheContainerRegistryUsingTheClusterIdentity)
	ctx.Step(`^the push request is rejected due to authorization$`, ps.thePushRequestIsRejectedDueToAuthorization)

	//CIS-6.1.4
	ctx.Step(`^a user attempts to deploy a container from an authorised registry$`, ps.aUserAttemptsToDeployAuthorisedContainer)
	ctx.Step(`^the deployment attempt is allowed$`, ps.theDeploymentAttemptIsAllowed)

	//CIS-6.1.5
	ctx.Step(`^a user attempts to deploy a container from an unauthorised registry$`, ps.aUserAttemptsToDeployUnauthorisedContainer)
	ctx.Step(`^the deployment attempt is denied$`, ps.theDeploymentAttemptIsDenied)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		cra.TeardownContainerAccessProbePod(ps.podState.PodName, p.Name())

		coreengine.LogScenarioEnd(s)
	})
}
