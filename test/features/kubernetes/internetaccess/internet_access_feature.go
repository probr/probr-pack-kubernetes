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
	code, err := kubernetes.AccessURL(&url)

	if err != nil {
		log.Printf("Error raised when attempting to access URL: %v", err)
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
	kubernetes.TeardownNetworkAccessTestPod()
}

// func FeatureContext(s *godog.Suite) {
// 	s.Step(`^a Kubernetes cluster is deployed$`, aKubernetesClusterIsDeployed)
// 	s.Step(`^a pod is deployed in the cluster$`, aPodIsDeployedInTheCluster)
// 	s.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, aProcessInsideThePodEstablishesADirectHttpsConnectionTo)
// 	s.Step(`^access is "([^"]*)"$`, accessIs)
// }

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	ps := probState{}
	ctx.AfterSuite(ps.tearDown)
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probState{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		ps.setup()
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a pod is deployed in the cluster$`, ps.aPodIsDeployedInTheCluster)
	ctx.Step(`^a process inside the pod establishes a direct http\(s\) connection to "([^"]*)"$`, ps.aProcessInsideThePodEstablishesADirectHTTPSConnectionTo)
	ctx.Step(`^access is "([^"]*)"$`, ps.accessIs)

}
