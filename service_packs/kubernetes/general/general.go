// Package general provides the implementation required to execute the BDD tests described in general.feature file
package general

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	apiv1 "k8s.io/api/core/v1"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes/connection"
	"github.com/citihub/probr/service_packs/kubernetes/constructors"
	"github.com/citihub/probr/utils"
)

type probeStruct struct{}

var conn connection.Connection

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

func (scenario *scenarioState) theKubernetesWebUIIsDisabled() error {

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	kubeSystemNamespace := config.Vars.ServicePacks.Kubernetes.SystemNamespace
	dashboardPodNamePrefix := config.Vars.ServicePacks.Kubernetes.DashboardPodNamePrefix
	stepTrace.WriteString(fmt.Sprintf("Attempt to find a pod in the '%s' namespace with the prefix '%s'; ", kubeSystemNamespace, dashboardPodNamePrefix))

	stepTrace.WriteString(fmt.Sprintf("Get all pods from '%s' namespace; ", kubeSystemNamespace))
	podList, getError := conn.GetPodsByNamespace(kubeSystemNamespace) // Also validates if provided namespace is valid
	if getError != nil {
		err = utils.ReformatError("An error occurred while retrieving pods from '%s' namespace. Error: %s", kubeSystemNamespace, getError)
		return err
	}

	stepTrace.WriteString(fmt.Sprintf("Confirm a pod with '%s' prefix doesn't exist; ", dashboardPodNamePrefix))
	for _, pod := range podList.Items {
		if strings.HasPrefix(pod.Name, dashboardPodNamePrefix) {
			err = utils.ReformatError("Dashboard UI Pod was found: '%s'", pod.Name)
			break
		}
	}

	payload = struct {
		KubeSystemNamespace    string
		DashboardPodNamePrefix string
	}{
		KubeSystemNamespace:    kubeSystemNamespace,
		DashboardPodNamePrefix: dashboardPodNamePrefix,
	}

	return err
}

func (scenario *scenarioState) aPodIsDeployedInTheCluster() error {

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf("Build a pod spec with default values; "))
	podObject := constructors.PodSpec(Probe.Name(), config.Vars.ServicePacks.Kubernetes.ProbeNamespace)

	stepTrace.WriteString(fmt.Sprintf("Create pod from spec; "))
	createdPodObject, creationErr := scenario.createPodfromObject(podObject)

	if creationErr != nil {
		err = utils.ReformatError("Pod creation did not succeed: %v", creationErr)
	}

	payload = struct {
		RequestedPod  *apiv1.Pod
		CreatedPod    *apiv1.Pod
		CreationError error
	}{
		RequestedPod:  podObject,
		CreatedPod:    createdPodObject,
		CreationError: creationErr,
	}

	return err
}

func (scenario *scenarioState) theResultOfAProcessInsideThePodEstablishingADirectHTTPConnectionToXIsY(urlAddress, result string) error {
	// Supported values for urlAddress:
	//	A valid absolute path URL with http(s) prefix
	//
	// Supported values for result:
	//	'blocked'

	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	// Guard clause - Validate url
	if _, urlErr := url.ParseRequestURI(urlAddress); urlErr != nil {
		err = utils.ReformatError("Invalid url provided.")
		return err
	}

	// Guard clause - Ensure pod was created in previous step
	if len(scenario.pods) == 0 {
		err = utils.ReformatError("Pod failed to create in the previous step")
		return err
	}

	// Validate input value
	var expectedCurlExitCode int
	switch result {
	case "blocked":
		expectedCurlExitCode = 6
		// Expecting curl exit code 6 (Couldn't resolve host) //TODO Confirm exit code is 7 in GCP
		// Ref: https://everything.curl.dev/usingcurl/returns
	default:
		err = utils.ReformatError("Unexpected value provided for expected command result: %s", result)
		return err
	}

	// Create a curl command to access the supplied url and only show http response in stdout.
	cmd := "curl -s -o /dev/null -I -L -w %{http_code} " + urlAddress

	stepTrace.WriteString(fmt.Sprintf("Attempt to run command in the pod: '%s'; ", cmd))
	exitCode, stdOut, cmdErr := conn.ExecCommand(cmd, scenario.namespace, scenario.pods[0])

	// Validate that no internal error occurred during execution of curl command
	if cmdErr != nil && exitCode == -1 {
		err = utils.ReformatError("Error raised when attempting to execute curl command inside container: %v", cmdErr)
		return err
	}

	// TODO: Confirm this implementation:
	// Expected Exit Code from Curl command is:	6

	stepTrace.WriteString("Check expected exit code was raised from curl command; ")
	if exitCode != expectedCurlExitCode {
		err = utils.ReformatError("Unexpected exit code: %d Error: %v", exitCode, cmdErr)
		return err
	}

	payload = struct {
		PodName              string
		Namespace            string
		Command              string
		ExpectedCurlExitCode int
		CurlExitCode         int
		StdOut               string
	}{
		PodName:              scenario.pods[0],
		Namespace:            scenario.namespace,
		Command:              cmd,
		ExpectedCurlExitCode: expectedCurlExitCode,
		CurlExitCode:         exitCode,
		StdOut:               stdOut,
	}

	return err
}

// Name presents the name of this probe for external reference
func (probe probeStruct) Name() string {
	return "general"
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

	ctx.AfterSuite(func() {})

}

// ScenarioInitialize provides initialization logic before each scenario is executed
func (probe probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&scenario, probe.Name(), s)
	})

	// Background
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, scenario.aKubernetesClusterIsDeployed)

	// Steps
	ctx.Step(`^the Kubernetes Web UI is disabled$`, scenario.theKubernetesWebUIIsDisabled)
	ctx.Step(`^a pod is deployed in the cluster$`, scenario.aPodIsDeployedInTheCluster)
	ctx.Step(`^the result of a process inside the pod establishing a direct http\(s\) connection to "([^"]*)" is "([^"]*)"$`, scenario.theResultOfAProcessInsideThePodEstablishingADirectHTTPConnectionToXIsY)

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

func (scenario *scenarioState) createPodfromObject(podObject *apiv1.Pod) (createdPodObject *apiv1.Pod, err error) {
	createdPodObject, err = conn.CreatePodFromObject(podObject, Probe.Name())
	if err == nil {
		scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	}
	return
}
