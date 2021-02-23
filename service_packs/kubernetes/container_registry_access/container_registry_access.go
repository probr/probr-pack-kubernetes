// Package container_registry_access provides the implementation required to execute the
// feature based test cases described in the the 'events' directory.
package container_registry_access

import (
	"fmt"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/citihub/probr/utils"
)

type probeStruct struct{}

// Probe meets the service pack interface for adding the logic from this file
var Probe probeStruct

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
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	description, payload, err := kubernetes.ClusterIsDeployed()
	stepTrace.WriteString(description)
	return err //  ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

// CIS-6.1.3
// Minimize cluster access to read-only
func (s *scenarioState) iAmAuthorisedToPullFromAContainerRegistry() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	// TODO: We are assuming too much here- if the image successfully pulls but fails to build, this will still fail
	pod, podAudit, err := cra.SetupContainerAccessProbePod(config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.PSPContainerAllowedImages, err)

	stepTrace.WriteString("Attempted to create a new pod using an image pulled from authorized registry; ")
	payload = struct {
		AuthorizedRegistry string
		PodAudit           *kubernetes.PodAudit
	}{
		AuthorizedRegistry: config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry,
		PodAudit:           podAudit,
	}
	return err
}

func (s *scenarioState) theDeploymentAttemptIsAllowed() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	// TODO: Extending the comment in iAmAuthorisedToPullFromAContainerRegistry...
	//       This step doesn't validate the attempt being allowed, it validates the success of the deployment
	err = kubernetes.AssertResult(&s.podState, "allowed", "")
	stepTrace.WriteString("Asserts pod creation result in scenario state is successful; ")
	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

// CIS-6.1.5
// Ensure deployment from authorised container registries is allowed
func (s *scenarioState) aUserAttemptsToDeployUnauthorisedContainer() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	pod, podAudit, err := cra.SetupContainerAccessProbePod(config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry, s.probe)

	err = kubernetes.ProcessPodCreationResult(&s.podState, pod, kubernetes.PSPContainerAllowedImages, err)

	stepTrace.WriteString(fmt.Sprintf(
		"Attempts to deploy a container from %s. Retains pod creation result in scenario state; ",
		config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry))
	payload = struct {
		UnauthorizedRegistry string
		PodAudit             *kubernetes.PodAudit
	}{
		UnauthorizedRegistry: config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry,
		PodAudit:             podAudit,
	}
	return err
}

func (s *scenarioState) theDeploymentAttemptIsDenied() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	err = kubernetes.AssertResult(&s.podState, "denied", "")
	stepTrace.WriteString("Asserts pod creation result in scenario state is denied; ")
	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

// Name presents the name of this probe for external reference
func (p probeStruct) Name() string {
	return "container_registry_access"
}

// Path presents the path of these feature files for external reference
func (p probeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	//check dependencies ...
	if cra == nil {
		// not been given one so set default
		cra = NewDefaultCRA()
	}
}

// ScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	//common
	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)

	//CIS-6.1.4
	ctx.Step(`^a user attempts to deploy a container from an authorised registry$`, ps.iAmAuthorisedToPullFromAContainerRegistry) // TODO: This step should be modified in the feature file, or a unique function should be written for it
	ctx.Step(`^the deployment attempt is allowed$`, ps.theDeploymentAttemptIsAllowed)

	//CIS-6.1.5
	ctx.Step(`^a user attempts to deploy a container from an unauthorised registry$`, ps.aUserAttemptsToDeployUnauthorisedContainer)
	ctx.Step(`^the deployment attempt is denied$`, ps.theDeploymentAttemptIsDenied)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		cra.TeardownContainerAccessProbePod(ps.podState.PodName, p.Name())

		coreengine.LogScenarioEnd(s)
	})
}
