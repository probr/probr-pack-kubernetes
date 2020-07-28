package containerregistryaccess

import (
	"fmt"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"citihub.com/probr/internal/coreengine"
	"citihub.com/probr/test/features"
	"github.com/cucumber/godog"
)

type probState struct {
	podName        string
	httpStatusCode int
}

func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.ContainerRegistryAccess, Name: "container_registry_access"}

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

func (p *probState) aUserAttemptsToDeployAContainerFrom(registry string) error {

	pd, err := kubernetes.SetupContainerAccessTestPod(&registry)

	if err != nil {
		return err
	}

	if pd == nil {
		// this is valid if the registry should be denied
		return nil
	}

	//hold on to pod name
	p.podName = pd.GetObjectMeta().GetName()

	//we're good ...
	return nil
}

func (p *probState) theDeploymentAttemptIs(res string) error {
	if res == "denied" {
		//expect pod name to be empty in this case (i.e. wasn't created)
		if p.podName != "" {
			//it's a fail:
			return fmt.Errorf("pod %v was created - test failed", p.podName)
		}
	}

	if res == "allowed" {
		// then expect the pod name to have a value
		if p.podName == "" {
			//it's a fail:
			return fmt.Errorf("pod was not created - test failed")
		}
	}

	//else we're good ...
	return nil
}

func (p *probState) setup() {
	//just make sure this is reset
	p.podName = ""
	p.httpStatusCode = 0
}

func (p *probState) tearDown() {
	kubernetes.TeardownContainerAccessTestPod(&p.podName)
}

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probState{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		ps.setup()
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a user attempts to deploy a container from "([^"]*)"$`, ps.aUserAttemptsToDeployAContainerFrom)
	ctx.Step(`^the deployment attempt is "([^"]*)"$`, ps.theDeploymentAttemptIs)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		ps.tearDown()
	})
}
