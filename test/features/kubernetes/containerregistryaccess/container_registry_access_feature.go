package containerregistryaccess

import (
	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
	"gitlab.com/citihub/probr/test/features"
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

//TODO: revise when interface this bit up ...
var cra kubernetes.ContainerRegistryAccess

// SetContainerRegistryAccess ...
func SetContainerRegistryAccess(c kubernetes.ContainerRegistryAccess) {
	cra = c
}

func (p *probState) aKubernetesClusterIsDeployed() error {
	b := cra.ClusterIsDeployed()

	if b == nil || !*b {
		return features.LogAndReturnError("kubernetes cluster is NOT deployed")
	}

	//else we're good ...
	return nil
}

func (p *probState) aUserAttemptsToDeployAContainerFrom(auth string, registry string) error {

	pd, err := cra.SetupContainerAccessTestPod(&registry)

	if err != nil {
		//check for partial creation, if we've got a pod, hold onto it's name so we can delete
		//(this could happen with imagepullerr etc)
		if pd != nil {
			p.podName = pd.GetObjectMeta().GetName()
		}

		//check for expected error
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			p.creationError = e
			return nil
		}
		//unexpected error
		return features.LogAndReturnError("error attempting to create POD: %v", err)
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
			return features.LogAndReturnError("pod %v was created - test failed", p.podName)
		}
		//should also check code:
		_, exists := p.creationError.ReasonCodes[kubernetes.PSPContainerAllowedImages]
		if !exists {
			//also a fail:
			return features.LogAndReturnError("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				p.creationError.ReasonCodes, kubernetes.PSPContainerAllowedImages)
		}

		//we're good
		return nil
	}

	if res == "allowed" {
		// then expect the pod name to be present
		if p.podName == "" {
			//it's a fail:
			return features.LogAndReturnError("pod was not created - test failed: %v", p.creationError)
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
	cra.TeardownContainerAccessTestPod(&p.podName)
}

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now

	//check dependancies ...
	if cra == nil {
		// not been given one so set default
		cra = kubernetes.NewDefaultCRA()
	}
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := probState{}

	ctx.BeforeScenario(func(*godog.Scenario) {
		ps.setup()
	})

	ctx.Step(`^a Kubernetes cluster is deployed$`, ps.aKubernetesClusterIsDeployed)
	ctx.Step(`^a user attempts to deploy a container from "([^"]*)" registry "([^"]*)"$`, ps.aUserAttemptsToDeployAContainerFrom)
	ctx.Step(`^the deployment attempt is "([^"]*)"$`, ps.theDeploymentAttemptIs)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		ps.tearDown()
	})
}
