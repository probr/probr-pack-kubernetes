package internetaccess

import (
	"fmt"
	"log"

	"citihub.com/probr/test/features"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"citihub.com/probr/internal/coreengine"
	"github.com/cucumber/godog"
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

func (p *probState) aKubernetesClusterIsDeployed() error {
	c, err := kubernetes.GetClient()
	if err != nil {
		return err
	}

	if c == nil {
		return fmt.Errorf("client is nil")
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

	pod, err := kubernetes.SetupNetworkAccessTestPod()

	if err != nil {
		return err
	}

	if pod == nil {
		return fmt.Errorf("POD is nil")
	}

	//hold on to the pod name
	p.podName = pod.GetObjectMeta().GetName()

	//else we're good ...
	return nil
}

func (p *probState) aProcessInsideThePodEstablishesADirectHTTPSConnectionTo(url string) error {
	code, err := kubernetes.AccessURL(&p.podName, &url)

	if err != nil {
		log.Printf("[ERROR] Error raised when attempting to access URL: %v", err)
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
			return fmt.Errorf("got HTTP Status Code %v - failed", p.httpStatusCode)
		}
	}
	//otherwise good
	return nil
}

func (p *probState) setup() {
	//anything?
}

func (p *probState) tearDown() {
	kubernetes.TeardownNetworkAccessTestPod(&p.podName)
}

func (p *probState) scenarioTearDown() {
	//reset the httpcode
	p.httpStatusCode=0
}

var ps = probState{}
//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now	

	ctx.AfterSuite(func() {
		ps.tearDown()
	})
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {
		ps.setup()
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ps.accessIs)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		ps.scenarioTearDown()
	})
}
