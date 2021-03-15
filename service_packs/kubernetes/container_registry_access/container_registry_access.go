// Package cra provides the implementation required to execute the BDD tests described in container_registry_access.feature file
package cra

import (
	"fmt"
	"log"

	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes/connection"
	"github.com/citihub/probr/service_packs/kubernetes/constructors"
	"github.com/citihub/probr/service_packs/kubernetes/errors"
	"github.com/citihub/probr/utils"
)

type probeStruct struct{}

// Will provide functionality to interact with K8s cluster
var conn connection.Connection

// scenarioState holds the steps and state for any scenario in this probe
type scenarioState struct {
	name        string
	currentStep string
	namespace   string
	audit       *audit.ScenarioAudit
	probe       *audit.Probe
	pods        []string
}

// Probe meets the service pack interface for adding the logic from this file
var Probe probeStruct
var scenario scenarioState

func (scenario *scenarioState) aKubernetesClusterIsDeployed() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()
	stepTrace.WriteString(fmt.Sprintf("Validate that a cluster can be reached using the specified kube config and context; "))

	payload = struct {
		KubeConfigPath string
		KubeContext    string
	}{
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		config.Vars.ServicePacks.Kubernetes.KubeContext,
	}

	err = conn.ClusterIsDeployed() // Must be assigned to 'err' be audited
	return err
}

func (scenario *scenarioState) podCreationXWithContainerImageFromYRegistry(expectedResult, registryAccess string) error {
	// Supported values for 'expectedResult':
	//	'succeeds'
	//	'is denied'

	// Supported values for 'registryAccess':
	//	'authorized'
	//	'unauthorized'

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	var shouldCreatePod bool
	// Validate input values
	switch expectedResult {
	case "succeeds":
		shouldCreatePod = true
	case "is denied":
		shouldCreatePod = false
	default:
		err = utils.ReformatError("Unexpected value provided for expectedResult: '%s' Expected values: ['succeeds', 'is denied']", expectedResult)
		return err
	}

	var isRegistryAuthorized bool
	// Validate input values
	switch registryAccess {
	case "authorized":
		isRegistryAuthorized = true
	case "unauthorized":
		isRegistryAuthorized = false
	default:
		err = utils.ReformatError("Unexpected value provided for registryAccess: '%s' Expected values: ['authorized', 'unauthorized']", registryAccess)
		return err
	}

	stepTrace.WriteString(fmt.Sprintf("Get appropriate container image from an '%s' registry; ", registryAccess))
	imageRegistry := getImageFromConfig(isRegistryAuthorized)

	stepTrace.WriteString(fmt.Sprintf("Build a pod spec with default values; "))
	podObject := constructors.PodSpec(Probe.Name(), scenario.namespace)

	stepTrace.WriteString(fmt.Sprintf("Set container image registry to appropriate value in pod spec; "))
	podObject.Spec.Containers[0].Image = imageRegistry

	stepTrace.WriteString(fmt.Sprintf("Create pod from spec; "))
	createdPodObject, creationErr := scenario.createPodfromObject(podObject) // Pod name is saved to scenario state if successful

	stepTrace.WriteString(fmt.Sprintf("Validate pod creation %s; ", expectedResult))
	switch shouldCreatePod {
	case true:
		if creationErr != nil {
			err = utils.ReformatError("Pod creation did not succeed: %v", creationErr)
		}
	case false:
		if creationErr == nil {
			err = utils.ReformatError("Pod creation succeeded, but should have been denied")
		} else {
			stepTrace.WriteString(fmt.Sprintf("Check that pod creation failed due to expected reason (403 Forbidden); "))
			if !errors.IsStatusCode(403, creationErr) {
				err = utils.ReformatError("Unexpected error during Pod creation : %v", creationErr)
			}
		}
	}

	payload = struct {
		ExpectedResult string
		RegistryAccess string
		ImageRegistry  string
		RequestedPod   *apiv1.Pod
		CreatedPod     *apiv1.Pod
		CreationError  error
	}{
		ExpectedResult: expectedResult,
		RegistryAccess: registryAccess,
		ImageRegistry:  imageRegistry,
		RequestedPod:   podObject,
		CreatedPod:     createdPodObject,
		CreationError:  creationErr,
	}

	return err
}

// Name presents the name of this probe for external reference
func (probe probeStruct) Name() string {
	return "container_registry_access"
}

// Path presents the path of these feature files for external reference
func (probe probeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", probe.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (probe probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		conn = connection.Get()
	})

	ctx.AfterSuite(func() {
	})
}

// ScenarioInitialize provides initialization logic before each scenario is executed
func (probe probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&scenario, probe.Name(), s)
	})

	// Background
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, scenario.aKubernetesClusterIsDeployed)

	// Steps
	ctx.Step(`^pod creation "([^"]*)" with container image from "([^"]*)" registry$`, scenario.podCreationXWithContainerImageFromYRegistry)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		afterScenario(scenario, probe, s, err)
	})

	ctx.BeforeStep(func(st *godog.Step) {
		scenario.currentStep = st.Text
	})

	ctx.AfterStep(func(st *godog.Step, err error) {
		scenario.currentStep = ""
	})
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.pods = make([]string, 0)
	s.namespace = config.Vars.ServicePacks.Kubernetes.ProbeNamespace
	coreengine.LogScenarioStart(gs)
}

func afterScenario(scenario scenarioState, probe probeStruct, gs *godog.Scenario, err error) {
	if config.Vars.ServicePacks.Kubernetes.KeepPods == "false" {
		for _, podName := range scenario.pods {
			err = conn.DeletePodIfExists(podName, scenario.namespace, probe.Name())
			if err != nil {
				log.Printf(fmt.Sprintf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err))
			}
		}
	}
	coreengine.LogScenarioEnd(gs)
}

func getContainerRegistryFromConfig(accessLevel bool) string {
	if accessLevel {
		return config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry
	}
	return config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry
}

func getImageFromConfig(accessLevel bool) string {
	registry := getContainerRegistryFromConfig(accessLevel)
	//full image is the repository + the configured image
	return registry + "/" + config.Vars.ServicePacks.Kubernetes.ProbeImage
}

func (scenario *scenarioState) createPodfromObject(podObject *apiv1.Pod) (createdPodObject *apiv1.Pod, err error) {
	createdPodObject, err = conn.CreatePodFromObject(podObject, Probe.Name())
	if err == nil {
		scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	}
	return
}
