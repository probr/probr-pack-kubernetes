package podsecurity

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cucumber/godog"

	"github.com/probr/probr-pack-kubernetes/internal/config"
	"github.com/probr/probr-pack-kubernetes/internal/connection"
	"github.com/probr/probr-pack-kubernetes/internal/summary"
	audit "github.com/probr/probr-sdk/audit"
	"github.com/probr/probr-sdk/probeengine"
	"github.com/probr/probr-sdk/providers/kubernetes/constructors"
	"github.com/probr/probr-sdk/providers/kubernetes/errors"
	"github.com/probr/probr-sdk/utils"

	apiv1 "k8s.io/api/core/v1"
)

type probeStruct struct {
}

// scenarioState holds the steps and state for any scenario in this probe
type scenarioState struct {
	name        string
	currentStep string
	namespace   string
	probeAudit  *audit.Probe
	audit       *audit.Scenario
	pods        []string
}

// Probe meets the service pack interface for adding the logic from this file
var Probe probeStruct
var scenario scenarioState

func (scenario *scenarioState) createPodfromObject(podObject *apiv1.Pod) (createdPodObject *apiv1.Pod, err error) {
	createdPodObject, err = connection.State.CreatePodFromObject(podObject, Probe.Name())
	if createdPodObject != nil && createdPodObject.ObjectMeta.Name != "" {
		scenario.pods = append(scenario.pods, createdPodObject.ObjectMeta.Name)
	}
	return
}

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
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		config.Vars.ServicePacks.Kubernetes.KubeContext,
	}

	err = connection.State.ClusterIsDeployed() // Must be assigned to 'err' be audited
	return err
}

func (scenario *scenarioState) toDo(todo string) error {
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("This step was included to inform developers that a scenario is incomplete; ")
	payload = struct {
		TODO string
	}{TODO: todo}
	return godog.ErrPending
}

// Attempt to deploy a pod from a default pod spec, with specified modification
func (scenario *scenarioState) podCreationResultsWithXSetToYInThePodSpec(result, key, value string) (err error) {
	// Supported key/values:
	// | Key                        | Value                                                     |
	// | 'allowPrivilegeEscalation' | 'true', 'false', 'not have a value provided'              |
	// | 'hostPID'                  | 'true', 'false', 'not have a value provided'              |
	// | 'hostIPC'                  | 'true', 'false', 'not have a value provided'              |
	// | 'hostNetwork'              | 'true', 'false', 'not have a value provided'              |
	// | 'user'                     | Any whole number (such as '0' or '1000')                  |
	// | 'annotations'              | 'include seccomp profile', 'not include seccomp profile'  |

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	podShouldCreate, err := shouldPodCreate(result)
	if err != nil {
		return
	}

	stepTrace.WriteString("Build a pod spec with default values; ")
	pod := constructors.PodSpec(Probe.Name(), config.Vars.ServicePacks.Kubernetes.ProbeNamespace, config.Vars.ServicePacks.Kubernetes.AuthorisedContainerImage)

	// Any key that expects a non-bool value should have it's own case here to handle the pod modification

	switch key {
	case "user":
		err = userPodSpecModifier(pod, value)
	case "annotations":
		err = annotationsPodSpecModifier(pod, value)
	case "capabilities":
		err = capabilitiesPodSpecModifier(pod, value)
	default:
		if value == "true" || value == "false" {
			err = boolPodSpecModifier(pod, key, value)
		} else if value != "not have a value provided" {
			err = utils.ReformatError("Expected 'true', 'false', or 'not have a value provided', but found '%s'", value)
		}
	}
	if err != nil {
		return
	}

	stepTrace.WriteString("Create pod from spec; ")
	createdPod, creationErr := scenario.createPodfromObject(pod)

	stepTrace.WriteString(fmt.Sprintf("Validate pod creation %s; ", result))
	switch podShouldCreate {
	case true:
		if creationErr != nil {
			err = utils.ReformatError("Pod creation did not succeed: %v", creationErr)
		}
	case false:
		if creationErr == nil {
			err = utils.ReformatError("Pod creation succeeded, but should have failed")
		} else {
			if !errors.IsStatusCode(403, creationErr) && !strings.Contains(creationErr.Error(), "ErrImagePull") {
				err = utils.ReformatError("Unexpected error during Pod creation : %v", creationErr)
			}
		}
	}

	payload = struct {
		RequestedPod  *apiv1.Pod
		CreatedPod    *apiv1.Pod
		CreationError error
	}{
		RequestedPod:  pod,
		CreatedPod:    createdPod,
		CreationError: creationErr,
	}
	return
}

func (scenario *scenarioState) theExecutionOfAXCommandInsideThePodIsY(cmdType, result string) error {
	// Supported cmdType:
	//     'non-privileged'
	//     'privileged'
	//     'root'
	//     'ping'
	//
	// Supported results:
	//     'successful'
	//     'prevented'

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	// Guard clause
	if len(scenario.pods) == 0 {
		err = utils.ReformatError("Pod failed to create in the previous step")
		return err
	}
	var expectedExitCodes []int
	var expectedFailureCodes []int

	var cmd string
	switch cmdType {
	case "non-privileged":
		cmd = "ls"
	case "privileged":
		cmd = "mount /fake /fake"
		expectedFailureCodes = []int{1, 32}
	case "root":
		cmd = "touch /dev/probr"
		expectedFailureCodes = []int{1}
	case "ping":
		cmd = "ping google.com"
		expectedFailureCodes = []int{1, 2}
	default:
		err = utils.ReformatError("Unexpected value provided for command type: %s", cmdType) // No payload is necessary if an invalid value was provided
		return err
	}

	switch result {
	case "successful":
		expectedExitCodes = []int{0}
	case "prevented":
		expectedExitCodes = expectedFailureCodes
	default:
		err = utils.ReformatError("Unexpected value provided for expected command result: %s", result) // No payload is necessary if an invalid value was provided
		return err

	}
	stepTrace.WriteString("Attempt to run a command in the pod that was created by the previous step; ")
	exitCode, stdout, stderr, err := connection.State.ExecCommand(cmd, scenario.namespace, scenario.pods[0])

	payload = struct {
		Command           string
		StdOut            string
		StdErr            string
		ExecErr           error
		ExitCode          int
		ExpectedExitCodes []int
	}{
		Command:           cmd,
		StdOut:            stdout,
		StdErr:            stderr,
		ExecErr:           err,
		ExitCode:          exitCode,
		ExpectedExitCodes: expectedExitCodes,
	}

	// Validate that no internal error occurred during execution of curl command
	if stderr != "" && exitCode == 0 {
		err = utils.ReformatError("Unknown error raised when attempting to execute '%s' inside container. Please review audit output for more information.", cmd)
		return err
	}

	var exitKnown bool
	for _, expectedCode := range expectedExitCodes {
		if exitCode == expectedCode {
			exitKnown = true
			err = nil
		}
	}
	if !exitKnown {
		err = utils.ReformatError("Unexpected exit code: %d. Please review audit output for more information.", exitCode)
	}
	return err
}

func (scenario *scenarioState) aXInspectionShouldOnlyShowTheContainerProcesses(inspectionType string) (err error) {
	// Supported inspection types:
	//     'process'
	//     'namespace'

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	var command string
	switch inspectionType {
	case "process":
		command = "ps"
	case "namespace":
		command = "lsns -n"
	default:
		err = utils.ReformatError("Unsupported value provided for inspection type")
		return
	}
	entrypoint := strings.Join(constructors.DefaultEntrypoint(), " ")
	exitCode, stdout, _, err := connection.State.ExecCommand(command, scenario.namespace, scenario.pods[0])

	if err != nil {
		// TODO: Validate that this fails as expected
		switch command {
		case "ps":
			stepTrace.WriteString("Validate that the container's entrypoint is PID 1 in the process tree; ")
			// NOTE: This particular expectation depends on using DefaultPodSecurityContext during the previous step
			//       Also, this explicitly assumes that we're using the alpine distro 'ps' (output is different for ubuntu, for example)
			expected := fmt.Sprintf("1 1000      0:00 %s", entrypoint)
			if !strings.Contains(stdout, expected) {
				err = utils.ReformatError("An entrypoint different from the container's was found for PID 1, suggesting hostPID was used")
			}
		case "lsns -n":
			stepTrace.WriteString("Validate that no namespace has an entrypoint different from the container's entrypoint; ")
			stdoutLines := strings.Split(stdout, "\n")
			for _, entry := range stdoutLines {
				if entry != "" && !strings.Contains(entry, entrypoint) {
					err = utils.ReformatError("A namespace is visible that uses a different entrypoint from the container, suggesting that hostIPC was used")
				}
			}
		}
	}

	payload = struct {
		Command    string
		ExitCode   int
		Stdout     string
		Entrypoint string
	}{
		Command:    command,
		ExitCode:   exitCode,
		Stdout:     stdout,
		Entrypoint: entrypoint,
	}
	return
}

func (scenario *scenarioState) thePodIPAndHostIPHaveDifferentValues() (err error) {
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = utils.ReformatError("[ERROR] Unexpected behavior occured: %s", panicErr)
		}
		scenario.audit.AuditScenarioStep(scenario.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Retrieve IP values from created pod; ")
	podIP, hostIP, err := connection.State.GetPodIPs(config.Vars.ServicePacks.Kubernetes.ProbeNamespace, scenario.pods[0])

	stepTrace.WriteString("Validate that PodIP and HostIP have different values; ")
	if err != nil && podIP == hostIP {
		err = utils.ReformatError("Pod IP and Host IP are identical, but should not be")
	}

	payload = struct {
		PodName string
		PodIP   string
		HostIP  string
	}{
		PodName: scenario.pods[0],
		PodIP:   podIP,
		HostIP:  hostIP,
	}
	return
}

// Name presents the name of this probe for external reference
func (probe probeStruct) Name() string {
	return "podsecurity"
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

// ScenarioInitialize initializes the specific test steps
func (probe probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&scenario, probe.Name(), s)
	})

	// Background
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, scenario.aKubernetesClusterIsDeployed)

	// Use for steps that have yet to be written
	ctx.Step(`^TODO: "([^"]*)"$`, scenario.toDo)

	// Parameterized Scenarios
	ctx.Step(`^pod creation "([^"]*)" with "([^"]*)" set to "([^"]*)" in the pod spec$`, scenario.podCreationResultsWithXSetToYInThePodSpec)
	ctx.Step(`^pod creation "([^"]*)" with "([^"]*)" set to "([^"]*)" in the pod spec$`, scenario.podCreationResultsWithXSetToYInThePodSpec)
	ctx.Step(`^the execution of a "([^"]*)" command inside the pod is "([^"]*)"$`, scenario.theExecutionOfAXCommandInsideThePodIsY)
	ctx.Step(`^a "([^"]*)" inspection should only show the container processes$`, scenario.aXInspectionShouldOnlyShowTheContainerProcesses)
	ctx.Step(`^the PodIP and HostIP have different values$`, scenario.thePodIPAndHostIPHaveDifferentValues)

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
	s.probeAudit = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.pods = make([]string, 0)
	s.namespace = config.Vars.ServicePacks.Kubernetes.ProbeNamespace
	probeengine.LogScenarioStart(gs)
}

func afterScenario(scenario scenarioState, probe probeStruct, gs *godog.Scenario, err error) { // TODO: err is overwritten before first use
	if config.Vars.ServicePacks.Kubernetes.KeepPods == "false" {
		for _, podName := range scenario.pods {
			err = connection.State.DeletePodIfExists(podName, scenario.namespace, probe.Name())
			if err != nil {
				log.Printf("[ERROR] Could not retrieve pod from namespace '%s' for deletion: %s", scenario.namespace, err)
			}
		}
	}
	probeengine.LogScenarioEnd(gs)
}

func boolPodSpecModifier(pod *apiv1.Pod, key, value string) (err error) {
	// Supported keys:
	//     'allowPrivilegeEscalation'
	//     'hostPID'
	//     'hostIPC'
	//     'hostNetwork'
	// Supported values:
	//     'true'
	//     'false'

	boolValue, _ := strconv.ParseBool(value)
	switch key {
	case "allowPrivilegeEscalation":
		pod.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation = &boolValue
	case "hostPID":
		pod.Spec.HostPID = boolValue
	case "hostIPC":
		pod.Spec.HostIPC = boolValue
	case "hostNetwork":
		pod.Spec.HostNetwork = boolValue
	default:
		err = utils.ReformatError("Unsupported key provided: %s", key) // No payload is necessary if an invalid key was provided
	}
	return
}

func annotationsPodSpecModifier(pod *apiv1.Pod, value string) (err error) {
	switch value {
	case "include seccomp profile":
		return // default
	case "not include seccomp profile":
		pod.ObjectMeta.Annotations = nil
	default:
		err = utils.ReformatError("Expected 'include seccomp profile' or 'not include seccomp profile', but found '%s'", value) // No payload is necessary if an invalid value was provided
	}
	return
}

func capabilitiesPodSpecModifier(pod *apiv1.Pod, value string) (err error) {
	switch value {
	case "drop NET_RAW":
		// default probe pod does this already
	case "add NET_RAW":
		pod.Spec.Containers[0].SecurityContext.Capabilities.Drop = []apiv1.Capability{} // clear default cap drop
		pod.Spec.Containers[0].SecurityContext.Capabilities.Add = append(pod.Spec.Containers[0].SecurityContext.Capabilities.Add, "NET_RAW")
	case "not have a value provided":
		pod.Spec.Containers[0].SecurityContext.Capabilities.Drop = []apiv1.Capability{}
	default:
		err = utils.ReformatError("Expected 'include NET_RAW' or 'not include NET_RAW', but found '%s'", value) // No payload is necessary if an invalid value was provided
	}
	return
}

func userPodSpecModifier(pod *apiv1.Pod, value string) (err error) {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		err = utils.ReformatError("Expected value to be a whole number, but found '%s' (%s)", value, err) // No payload is necessary if an invalid value was provided
		return err
	}
	pod.Spec.SecurityContext.RunAsUser = &intValue
	return
}

func shouldPodCreate(result string) (shouldCreate bool, err error) {
	switch result {
	case "succeeds":
		shouldCreate = true
	case "fails":
		shouldCreate = false
	default:
		err = utils.ReformatError("Unexpected value provided for expected pod creation result: %s", result) // No payload is necessary if an invalid value was provided
	}
	return
}
