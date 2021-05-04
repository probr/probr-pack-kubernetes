// Package cra provides the implementation required to execute the BDD tests described in container_registry_access.feature file
package cra

import (
	"fmt"
	"log"

	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
	"github.com/citihub/probr-pack-kubernetes/internal/connection"
	"github.com/citihub/probr-pack-kubernetes/internal/constructors"
	"github.com/citihub/probr-pack-kubernetes/internal/errors"
	"github.com/citihub/probr-pack-kubernetes/internal/summary"
	audit "github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/utils"
)

type probeStruct struct{}

// scenarioState holds the steps and state for any scenario in this probe
type scenarioState struct {
	name        string
	currentStep string
	namespace   string
	audit       *audit.Scenario
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
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that a cluster can be reached using the specified kube config and context; ")

	payload = struct {
		KubeConfigPath string
		KubeContext    string
	}{
		config.Vars.KubeConfigPath,
		config.Vars.KubeContext,
	}

	err = connection.State.ClusterIsDeployed() // Must be assigned to 'err' be audited
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
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
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

	stepTrace.WriteString("Build a pod spec with default values; ")
	podObject := constructors.PodSpec(Probe.Name(), scenario.namespace, config.Vars.AuthorisedContainerImage)

	stepTrace.WriteString(fmt.Sprintf("Set container image registry to '%s' value in pod spec; ", registryAccess))
	podObject.Spec.Containers[0].Image = imageFromConfig(isRegistryAuthorized)

	stepTrace.WriteString("Create pod from spec; ")
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
			stepTrace.WriteString("Check that pod creation failed due to expected reason (403 Forbidden); ")
			if !errors.IsStatusCode(403, creationErr) {
				err = utils.ReformatError("Unexpected error during Pod creation : %v", creationErr)
			}
		}
	}

	payload = struct {
		ExpectedResult string
		RegistryAccess string
		RequestedPod   *apiv1.Pod
		CreatedPod     *apiv1.Pod
		CreationError  error
	}{
		ExpectedResult: expectedResult,
		RegistryAccess: registryAccess,
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
	return probeengine.GetFeaturePath("internal", probe.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (probe probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
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
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.pods = make([]string, 0)
	s.namespace = config.Vars.ProbeNamespace
	probeengine.LogScenarioStart(gs)
}

func afterScenario(scenario scenarioState, probe probeStruct, gs *godog.Scenario, err error) {
	if config.Vars.KeepPods == "false" {
		for _, podName := range scenario.pods {
			err = connection.State.DeletePodIfExists(podName, scenario.namespace, probe.Name())
			if err != nil {
				log.Printf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err)
			}
		}
	}
	probeengine.LogScenarioEnd(gs)
}

func imageFromConfig(authorized bool) string {
	if authorized {
		return config.Vars.AuthorisedContainerImage
	}
	return config.Vars.UnauthorisedContainerImage
}

func (scenario *scenarioState) createPodfromObject(podObject *apiv1.Pod) (createdPodObject *apiv1.Pod, err error) {
	createdPodObject, err = connection.State.CreatePodFromObject(podObject, Probe.Name())
	if err == nil {
		scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	}
	return
}
