package internetaccess

import (
	"log"

	"gitlab.com/citihub/probr/test/features"

	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
)

type probState struct {
	podName        string
	httpStatusCode int
}

func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.InternetAccess, Name: "internet_access"}

	coreengine.TestHandleFunc(td, &coreengine.GoDogTestTuple{
		Handler: features.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
		},
	})
}

//TODO: revise when interface this bit up ...
var na kubernetes.NetworkAccess

// SetNetworkAccess ...
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
		return features.LogAndReturnError("POD is nil")
	}

	//hold on to the pod name
	p.podName = pod.GetObjectMeta().GetName()

	//else we're good ...
	return nil
}

func (p *probState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := na.AccessURL(&p.podName, &url)

	if err != nil {
		features.LogAndReturnError("[ERROR] Error raised when attempting to access URL: %v", err)
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
			return features.LogAndReturnError("got HTTP Status Code %v - failed", p.httpStatusCode)
		}
	}
	//otherwise good
	return nil
}

func (p *probState) setup() {
	//anything?
}

func (p *probState) tearDown() {
	na.TeardownNetworkAccessTestPod(&p.podName)
}

func (p *probState) scenarioTearDown() {
	//reset the httpcode
	p.httpStatusCode = 0
}

var ps = probState{}

//TestSuiteInitialize ...
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

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.setup()
		features.LogScenarioStart(s)
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ps.accessIs)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ps.scenarioTearDown()
		features.LogScenarioEnd(s)
	})
}
