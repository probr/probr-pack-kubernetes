package podsecuritypolicy

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
		Category: coreengine.PodSecurityPolicies, Name: "pod_security_policy"}

	coreengine.TestHandleFunc(td, &coreengine.GoDogTestTuple{
		Handler: features.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
		},
	})
}

func (p *probState) creationWillWithAMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

func (p *probState) aKubernetesClusterExistsWhichWeCanDeployInto() error {
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

func (p *probState) aKubernetesDeploymentIsAppliedToTheActiveKubernetesCluster() error {
	//TODO: not sure this step is adding value ... return "pass" for now ...
	return nil
}

func (p *probState) privilegedAccessRequestIsMarkedForTheKubernetesDeployment(privilegedAccessRequested string) error {	
	
	var pa kubernetes.PrivilegedAccess
	if privilegedAccessRequested == "True" {
		pa = kubernetes.WithPrivilegedAccess
	} else {
		pa = kubernetes.WithoutPrivilegedAccess
	}

	pd, err := kubernetes.CreatePODSettingPrivilegedAccess(pa)
	if err != nil {
		return fmt.Errorf("error attempting to create POD: %v", err)
	}

	if pd == nil {
		// valid if the request was for privileged (i.e. pod creation should fail)
		return nil
	}

	//hold on to the pod name
	p.podName = pd.GetObjectMeta().GetName()

	//we're good
	return nil
}

func (p *probState) someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster() error {
	yesNo, err := kubernetes.PrivilegedAccessIsRestricted()

	if err != nil {
		return fmt.Errorf("error determining Pod Security Policy %v", err)
	}
	if yesNo == nil {
		return fmt.Errorf("result of PrivilegedAccessIsRestricted is nil despite no error being raised from the call")
	}

	if !*yesNo {
		return fmt.Errorf("Privileged Access is NOT restricted (result: %t)", *yesNo)
	}

	return nil
}

func (p *probState) theOperationWillWithAnError(res, msg string) error {
	if res == "Fail" {
		//expect pod name to be empty in this case (i.e. wasn't created)
		if p.podName != "" {
			//it's a fail:
			return fmt.Errorf("pod %v was created - test failed", p.podName)
		}
	}

	if res == "Succeed" {
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
	kubernetes.TeardownPodSecurityTestPod(&p.podName)
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

	ctx.Step(`^A kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterExistsWhichWeCanDeployInto)
	ctx.Step(`^a Kubernetes deployment is applied to the active Kubernetes cluster$`, ps.aKubernetesDeploymentIsAppliedToTheActiveKubernetesCluster)
	ctx.Step(`^privileged access request is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedAccessRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some control exists to prevent privileged access for kubernetes deployments to an active kubernetes cluster$`, ps.someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster)
	ctx.Step(`^the operation will "([^"]*)" with an error "([^"]*)"$`, ps.theOperationWillWithAnError)

	ctx.AfterScenario(func(sc *godog.Scenario, err error) {
		ps.tearDown()
	})
}
