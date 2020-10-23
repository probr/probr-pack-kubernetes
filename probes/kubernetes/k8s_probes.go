package k8s_probes

import (
	"log"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/summary"
	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"
)

// podState captures useful pod state data for use in a scenario's state.
type podState struct {
	PodName         string
	CreationError   *kubernetes.PodCreationError
	ExpectedReason  *kubernetes.PodCreationErrorReason
	CommandExitCode int
}

type scenarioState struct {
	name           string
	audit          *summary.ScenarioAudit
	probe          *summary.Probe
	httpStatusCode int
	podName        string
	podState       podState
	useDefaultNS   bool
	wildcardRoles  interface{}
}

type Probe int

const (
	ContainerRegistryAccess Probe = iota
	General
	PodSecurityPolicy
	InternetAccess
	IAMControl
)

// Probes contains all probes with helper functions allowing all to be added in a loop
var Probes []Probe

func init() {
	Probes = []Probe{
		ContainerRegistryAccess,
		General,
		PodSecurityPolicy,
		InternetAccess,
		IAMControl,
	}
}

func (p Probe) String() string {
	return [...]string{
		"container_registry_access",
		"general",
		"pod_security_policy",
		"internet_access",
		"iam_control",
	}[p]
}

func (p Probe) TestSuiteContext(s *godog.TestSuiteContext) {
	f := [...]func(*godog.TestSuiteContext){
		craTestSuiteInitialize,
		genTestSuiteInitialize,
		pspTestSuiteInitialize,
		iaTestSuiteInitialize,
		iamTestSuiteInitialize,
	}[p]
	f(s)
}

func (p Probe) ScenarioContext(s *godog.ScenarioContext) {
	f := [...]func(*godog.ScenarioContext){
		craScenarioInitialize,
		genScenarioInitialize,
		pspScenarioInitialize,
		iaScenarioInitialize,
		iamScenarioInitialize,
	}[p]
	f(s)
}

func (p Probe) GetGodogTest() *coreengine.GodogTest {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes, Name: p.String()}

	return &coreengine.GodogTest{
		TestDescriptor:       &td,
		TestSuiteInitializer: p.TestSuiteContext,
		ScenarioInitializer:  p.ScenarioContext,
	}
}

//
// Helper Functions

func (s *scenarioState) BeforeScenario(probeName string, gs *godog.Scenario) {
	if coreengine.TagsNotExcluded(gs.Tags) {
		s.setup()
		s.name = gs.Name
		s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
		coreengine.LogScenarioStart(gs)
	}
}

// Setup resets scenario-specific values
func (s *scenarioState) setup() {
	s.podState.PodName = ""
	s.podState.CreationError = nil
	s.useDefaultNS = false
}

// ProcessPodCreationResult is a convenince function to process the result of a pod creation attempt.
// It records state information on the supplied state structure.
func ProcessPodCreationResult(probe *summary.Probe, s *podState, pd *apiv1.Pod, expected kubernetes.PodCreationErrorReason, err error) error {
	//first check for errors:
	if err != nil {
		//check if we've got a partial pod creation
		//e.g. pod was created but didn't get to "running" state
		//in this case we need to hold onto the name so it can be deleted
		if pd != nil {
			s.PodName = pd.GetObjectMeta().GetName()
			probe.CountPodCreated()
			summary.State.LogPodName(s.PodName)
		}

		//check for known error type
		//this means the pod has not been created for an expected reason and
		//is a valid result if the test is addressing prevention of insecure pod creation
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			s.CreationError = e
			s.ExpectedReason = &expected
			return nil
		}
		//unexpected error
		//in this case something unexpected has happened, return an error to cucumber
		return coreengine.LogAndReturnError("error attempting to create POD: %v", err)
	}

	//No errors: pod creation may or may not have been expected.  This will be determined
	//by the specific test case
	if pd == nil {
		// pod not created, which could be valid for some tests
		return nil
	}

	//if we've got this far, a pod was successfully created which could be
	//valid for some tests
	s.PodName = pd.GetObjectMeta().GetName()
	probe.CountPodCreated()
	summary.State.LogPodName(s.PodName)

	//we're good
	return nil
}

// AssertResult evaluate the state in the context of the expected condition, e.g. if expected is "fail",
// then the expecation is that a creation error will be present.
func AssertResult(s *podState, res, msg string) error {

	if res == "Fail" || res == "denied" {
		//expect pod creation error to be non-null
		if s.CreationError == nil {
			//it's a fail:
			return coreengine.LogAndReturnError("pod %v was created - test failed", s.PodName)
		}
		//should also check code:
		_, exists := s.CreationError.ReasonCodes[*s.ExpectedReason]
		if !exists {
			//also a fail:
			return coreengine.LogAndReturnError("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				s.CreationError.ReasonCodes, s.ExpectedReason)
		}

		//we're good
		return nil
	}

	if res == "Succeed" || res == "allowed" {
		// then expect the pod creation error to be nil
		if s.CreationError != nil {
			//it's a fail:
			return coreengine.LogAndReturnError("pod was not created - test failed: %v", s.CreationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	return coreengine.LogAndReturnError("desired result %v is not recognised", res)

}

//general feature steps:
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	b := kubernetes.GetKubeInstance().ClusterIsDeployed()

	if b == nil || !*b {
		log.Fatalf("[ERROR] Kubernetes cluster is not deployed")
	}

	description := "Passes if Probr successfully connects to the specified cluster."
	payload := struct {
		KubeConfigPath string
		KubeContext    string
	}{config.Vars.KubeConfigPath, config.Vars.KubeContext}
	s.audit.AuditScenarioStep(description, payload, nil)

	return nil
}

func podPayload(pod *apiv1.Pod, podAudit *kubernetes.PodAudit) interface{} {
	return struct {
		Pod      *apiv1.Pod
		PodAudit *kubernetes.PodAudit
	}{
		Pod:      pod,
		PodAudit: podAudit,
	}
}
