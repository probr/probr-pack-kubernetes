package k8s_probes

import (
	"log"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"
)

var ia_ps scenarioState

func init() {
	ia_ps = scenarioState{}
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes, Name: ia_name}

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: coreengine.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: iaTestSuiteInitialize,
			ScenarioInitializer:  iaScenarioInitialize,
		},
	})
}

// NetworkAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// network access scenarios.
var na kubernetes.NetworkAccess

// SetNetworkAccess allows injection of a specific NetworkAccess helper.
func SetNetworkAccess(n kubernetes.NetworkAccess) {
	na = n
}

func (s *scenarioState) aPodIsDeployedInTheCluster() error {
	var err error
	var podAudit *kubernetes.PodAudit
	var pod *apiv1.Pod
	if s.podName != "" {
		//only one pod is needed for all scenarios in this probe
		log.Printf("[DEBUG] Pod %v has already been created - reusing the pod", s.podName)
	} else {
		pd, pa, e := na.SetupNetworkAccessTestPod()
		podAudit = pa
		pod = pd
		if e != nil {
			err = e
		} else if pod == nil {
			err = coreengine.LogAndReturnError("Failed to setup network access test pod")
		} else {
			s.podName = pod.GetObjectMeta().GetName()
		}
	}

	description := ""
	payload := podPayload(pod, podAudit)
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := na.AccessURL(&s.podName, &url)

	if err != nil {
		err = coreengine.LogAndReturnError("[ERROR] Error raised when attempting to access URL: %v", err)
	}

	//hold on to the code
	s.httpStatusCode = code

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) accessIs(accessResult string) error {
	var err error
	if accessResult == "blocked" {
		//then the result should be anything other than 200
		if s.httpStatusCode == 200 {
			//it's a fail:
			err = coreengine.LogAndReturnError("got HTTP Status Code %v - failed", s.httpStatusCode)
		}
	}

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

// iaTestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func iaTestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {
		na.TeardownNetworkAccessTestPod(&ia_ps.podName, ia_name)
	})

	//check dependancies ...
	if na == nil {
		// not been given one so set default
		na = kubernetes.NewDefaultNA()
	}
}

// iaScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func iaScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		ia_ps.BeforeScenario(ia_name, s)
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ia_ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ia_ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ia_ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ia_ps.accessIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ia_ps.httpStatusCode = 0
		coreengine.LogScenarioEnd(s)
	})
}
