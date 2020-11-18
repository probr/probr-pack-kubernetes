package internet_access

import (
	"log"

	"github.com/cucumber/godog"
	apiv1 "k8s.io/api/core/v1"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/kubernetes"
)

type ProbeStruct struct{}

var Probe ProbeStruct

var ia_ps scenarioState

func init() {
	ia_ps = scenarioState{}
}

// NetworkAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// network access scenarios.
var na NetworkAccess

// SetNetworkAccess allows injection of a specific NetworkAccess helper.
func SetNetworkAccess(n NetworkAccess) {
	na = n
}

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload := kubernetes.ClusterIsDeployed()
	s.audit.AuditScenarioStep(description, payload, nil)
	return nil // ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

func (s *scenarioState) aPodIsDeployedInTheCluster() error {
	var err error
	var podAudit *kubernetes.PodAudit
	var pod *apiv1.Pod
	if s.podName != "" {
		//only one pod is needed for all scenarios in this probe
		log.Printf("[DEBUG] Pod %v has already been created - reusing the pod", s.podName)
	} else {
		pd, pa, e := na.SetupNetworkAccessProbePod()
		podAudit = pa
		pod = pd
		if e != nil {
			err = e
		} else if pod == nil {
			err = utils.ReformatError("Failed to setup network access test pod")
			log.Print(err)
		} else {
			s.podName = pod.GetObjectMeta().GetName()
		}
	}

	description := ""
	payload := kubernetes.PodPayload{Pod: pod, PodAudit: podAudit}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := na.AccessURL(&s.podName, &url)

	if err != nil {
		err = utils.ReformatError("[ERROR] Error raised when attempting to access URL: %v", err)
		log.Print(err)
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
			err = utils.ReformatError("got HTTP Status Code %v - failed", s.httpStatusCode)
		}
	}

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (p ProbeStruct) Name() string {
	return "internet_access"
}

// iaProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {
		na.TeardownNetworkAccessProbePod(&ia_ps.podName, p.Name())
	})

	//check dependancies ...
	if na == nil {
		// not been given one so set default
		na = NewDefaultNA()
	}
}

// iaScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ia_ps, p.Name(), s)
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
