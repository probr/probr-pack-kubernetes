package containerregistryaccess

import (
	"fmt"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"citihub.com/probr/internal/coreengine"
	"citihub.com/probr/test/features"
	"github.com/cucumber/godog"
)

type probState struct {
	podName       string
	creationError *kubernetes.PodCreationError
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
		//check for expected error
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			p.creationError = e
			return nil
		}
		//unexpected error
		return fmt.Errorf("error attempting to create POD: %v", err)
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
		//expect pod creation error to be non-null (i.e. creation was prevented)
		if p.creationError == nil {
			//it's a fail:
			return fmt.Errorf("pod %v was created - test failed", p.podName)
		}
		//should also check code:
		_, exists := p.creationError.ReasonCodes[kubernetes.PSPContainerAllowedImages]
		if !exists {		
			//also a fail:
			return fmt.Errorf("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				p.creationError.ReasonCodes, kubernetes.PSPContainerAllowedImages)
		}

		//we're good
		return nil
	}

	if res == "allowed" {
		// then expect the pod name to be present
		if p.podName == "" {
			//it's a fail:
			return fmt.Errorf("pod was not created - test failed: %v", p.creationError)
		}
	}

	//else we're good ...
	return nil
}

func (p *probState) setup() {
	//just make sure this is reset
	p.podName = ""
	p.creationError = nil
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
