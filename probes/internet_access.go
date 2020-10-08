package probes

import (
	"log"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/cucumber/godog"
)

const IA_NAME = "internet_access"

var IA_PROBE probeState

func init() {
	IA_PROBE = probeState{}
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.InternetAccess, Name: IA_NAME}

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: iaTestSuiteInitialize,
			ScenarioInitializer:  iaScenarioInitialize,
		},
	})
}

// NetworkAccess is the section of the kubernetes package which provides the kubernetes interactions required to support
// network access probes.
var na kubernetes.NetworkAccess

// SetNetworkAccess allows injection of a specific NetworkAccess helper.
func SetNetworkAccess(n kubernetes.NetworkAccess) {
	na = n
}

func (p *probeState) aPodIsDeployedInTheCluster() error {
	var err error
	if p.podName != "" {
		//only one pod is needed for all probes in this event
		log.Printf("[DEBUG] Pod %v has already been created - reusing the pod", p.podName)
	} else {
		pod, e := na.SetupNetworkAccessTestPod()
		if e != nil {
			err = e
		} else if pod == nil {
			err = LogAndReturnError("Failed to setup network access test pod")
		} else {
			p.podName = pod.GetObjectMeta().GetName()
		}
	}
	p.event.LogProbe(p.name, err)
	return err
}

func (p *probeState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := na.AccessURL(&p.podName, &url)

	if err != nil {
		err = LogAndReturnError("[ERROR] Error raised when attempting to access URL: %v", err)
	}

	//hold on to the code
	p.httpStatusCode = code
	p.event.LogProbe(p.name, err)
	return err
}

func (p *probeState) accessIs(accessResult string) error {
	var err error
	if accessResult == "blocked" {
		//then the result should be anything other than 200
		if p.httpStatusCode == 200 {
			//it's a fail:
			err = LogAndReturnError("got HTTP Status Code %v - failed", p.httpStatusCode)
		}
	}
	p.event.LogProbe(p.name, err)
	return err
}

// iaTestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func iaTestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {
		na.TeardownNetworkAccessTestPod(&IA_PROBE.podName, IA_NAME)
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
		BeforeScenario(IA_NAME, &IA_PROBE, s)
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, IA_PROBE.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, IA_PROBE.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, IA_PROBE.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, IA_PROBE.accessIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		IA_PROBE.httpStatusCode = 0
		LogScenarioEnd(s)
	})
}
