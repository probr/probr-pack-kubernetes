package kubernetes

import (
	"log"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/markbates/pkger"
	apiv1 "k8s.io/api/core/v1"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
)

var AssetsDir string

// podState captures useful pod state data for use in a scenario's state.
type PodState struct {
	PodName         string
	CreationError   *PodCreationError
	ExpectedReason  *PodCreationErrorReason
	CommandExitCode int
}

type scenarioState struct {
	name           string
	audit          *summary.ScenarioAudit
	probe          *summary.Probe
	httpStatusCode int
	podName        string
	podState       PodState
	useDefaultNS   bool
	wildcardRoles  interface{}
}

const (
	Namespace = "probr-general-test-ns"
)

type PodPayload struct {
	Pod      *apiv1.Pod
	PodAudit *PodAudit
}

func init() {
	AssetsDir = filepath.Join("service_packs", "kubernetes", "assets")

	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// It is needed, as variables are not supported in pkger
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/kubernetes/assets")
}

//
// Helper Functions

func BeforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.setup()
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}

// Setup resets scenario-specific values
func (s *scenarioState) setup() {
	s.podState.PodName = ""
	s.podState.CreationError = nil
	s.useDefaultNS = false
}

// ProcessPodCreationResult is a convenience function to process the result of a pod creation attempt.
// It records state information on the supplied state structure.
func ProcessPodCreationResult(state *PodState, pd *apiv1.Pod, expected PodCreationErrorReason, err error) error {
	//first check for errors:
	if err != nil {
		//check if we've got a partial pod creation
		//e.g. pod was created but didn't get to "running" state
		//in this case we need to hold onto the name so it can be deleted
		if pd != nil {
			state.PodName = pd.GetObjectMeta().GetName()
		}

		//check for known error type
		//this means the pod has not been created for an expected reason and
		//is a valid result if the test is addressing prevention of insecure pod creation
		if e, ok := err.(*PodCreationError); ok {
			state.CreationError = e
			state.ExpectedReason = &expected
			return nil
		}
		//unexpected error
		//in this case something unexpected has happened, return an error to cucumber
		return utils.ReformatError("error attempting to create POD: %v", err)
	}

	//No errors: pod creation may or may not have been expected.  This will be determined
	//by the specific test case
	if pd == nil {
		// pod not created, which could be valid for some tests
		return nil
	}

	//if we've got this far, a pod was successfully created which could be
	//valid for some tests
	state.PodName = pd.GetObjectMeta().GetName()

	//we're good
	return nil
}

// AssertResult evaluate the state in the context of the expected condition, e.g. if expected is "fail",
// then the expectation is that a creation error will be present.
func AssertResult(s *PodState, res, msg string) error {

	if res == "Fail" || res == "denied" {
		//expect pod creation error to be non-null
		if s.CreationError == nil {
			//it's a fail:
			return utils.ReformatError("pod %v was created - test failed", s.PodName)
		}
		//should also check code:
		_, exists := s.CreationError.ReasonCodes[*s.ExpectedReason]
		if !exists {
			//also a fail:
			return utils.ReformatError("pod was not created but failure reasons (%v) did not contain expected (%v)- test failed",
				s.CreationError.ReasonCodes, s.ExpectedReason)
		}

		//we're good
		return nil
	}

	if res == "Succeed" || res == "allowed" {
		// then expect the pod creation error to be nil
		if s.CreationError != nil {
			//it's a fail:
			return utils.ReformatError("pod was not created - test failed: %v", s.CreationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	err := utils.ReformatError("desired result %v is not recognised", res)
	log.Print(err)
	return err

}

type ClusterPayload struct {
	KubeConfigPath string
	KubeContext    string
}

//general feature steps:
func ClusterIsDeployed() (string, ClusterPayload, error) {
	var err error
	b := GetKubeInstance().ClusterIsDeployed()
	if b == nil || !*b {
		err = utils.ReformatError("Kubernetes cluster is not deployed")
		log.Print(err)
	}
	description := "Passes if Probr successfully connects to the specified cluster."
	payload := ClusterPayload{config.Vars.ServicePacks.Kubernetes.KubeConfigPath, config.Vars.ServicePacks.Kubernetes.KubeContext}
	return description, payload, err
}
