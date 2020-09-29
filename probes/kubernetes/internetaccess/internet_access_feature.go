// Package internetaccess provides the implementation required to execute the feature based test cases described in the
// the 'features' directory.
package internetaccess

import (
	"log"

	"gitlab.com/citihub/probr/probes"

	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
)

type probState struct {
	podName        string
	httpStatusCode int
}

const NAME = "internet_access"

func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.InternetAccess, Name: NAME}

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: probes.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
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

func (p *probState) aKubernetesClusterIsDeployed() error {
	b := na.ClusterIsDeployed()

	if b == nil || !*b {
		log.Fatalf("[ERROR] Kubernetes cluster is not deployed")
	}

	//else we're good ...
	return nil
}

func (p *probState) aPodIsDeployedInTheCluster() error {
	//only one pod is needed for all scenarios
	//if we have a pod name, then it's already created so
	//this step can be skipped and the pod will be reused
	if p.podName != "" {
		log.Printf("[INFO] Pod %v has already been created - reusing the pod", p.podName)
		return nil
	}

	pod, err := na.SetupNetworkAccessTestPod()

	if err != nil {
		return err
	}

	if pod == nil {
		return probes.LogAndReturnError("POD is nil")
	}

	//hold on to the pod name
	p.podName = pod.GetObjectMeta().GetName()

	//else we're good ...
	return nil
}

func (p *probState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := na.AccessURL(&p.podName, &url)

	if err != nil {
		probes.LogAndReturnError("[ERROR] Error raised when attempting to access URL: %v", err)
		return err
	}

	//hold on to the code
	p.httpStatusCode = code

	return nil
}

func (p *probState) accessIs(accessResult string) error {
	if accessResult == "blocked" {
		//then the result should be anything other than 200
		if p.httpStatusCode == 200 {
			//it's a fail:
			return probes.LogAndReturnError("got HTTP Status Code %v - failed", p.httpStatusCode)
		}
	}
	//otherwise good
	return nil
}

func (p *probState) setup() {
	//anything?
}

func (p *probState) tearDown() {
	na.TeardownNetworkAccessTestPod(&p.podName, NAME)
}

func (p *probState) scenarioTearDown() {
	//reset the httpcode
	p.httpStatusCode = 0
}

var ps = probState{}

// TestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ctx.AfterSuite(func() {
		ps.tearDown()
	})

	//check dependancies ...
	if na == nil {
		// not been given one so set default
		na = kubernetes.NewDefaultNA()
	}
}

// ScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the features directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.setup()
		probes.LogScenarioStart(s)
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ps.accessIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ps.scenarioTearDown()
		probes.LogScenarioEnd(s)
	})
}
