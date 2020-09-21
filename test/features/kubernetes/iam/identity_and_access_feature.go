// Package iam provides the implementation required to execute the feature based test cases described in the
// the 'features' directory.  The 'assets' directory holds any assets required for the test cases.   Assets are 'embedded'
// via the 'go-bindata.exe' tool which is invoked via the 'go generate' tool.  It is important, therefore, that the 
//'go:generate' comment is present in order to include this package in the scope of the 'go generate' tool.  This can be 
// invoked directly on the command line of via the Makefile (e.g. make clean-build).
package iam

//go:generate go-bindata.exe -pkg $GOPACKAGE -o assets/assets.go assets/yaml

import (
	"strings"

	"github.com/cucumber/godog"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"
	"gitlab.com/citihub/probr/test/features"
	"gitlab.com/citihub/probr/test/features/kubernetes/probe"

	iamassets "gitlab.com/citihub/probr/test/features/kubernetes/iam/assets"
)

type probeState struct {
	state        probe.State
	useDefaultNS bool
}

// IdentityAccessManagement is the section of the kubernetes package which provides the kubernetes interactions required to support
// identity access management probes.
var iam kubernetes.IdentityAccessManagement

// SetIAM allows injection of an IdentityAccessManagement helper.
func SetIAM(i kubernetes.IdentityAccessManagement) {
	iam = i
}

// init() registers the feature tests descibed in this package with the test runner (coreengine.TestRunner) via the call
// to coreengine.TestHandleFunc.  This links the test - described by the TestDescriptor - with the handler to invoke.  In
// this case, the general test handler is being used (features.GodogTestHandler) and the GodogTest data provides the data
// require to execute the test.  Specifically, the data includes the Test Suite and Scenario Initializers from this package
// which will be called from features.GodogTestHandler.  Note: a blank import at probr library level should be done to
// invoke this function automatically on initial load.
func init() {
	td := coreengine.TestDescriptor{Group: coreengine.Kubernetes,
		Category: coreengine.IAM, Name: "iam_control"}

	coreengine.TestHandleFunc(td, &coreengine.GoDogTestTuple{
		Handler: features.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: TestSuiteInitialize,
			ScenarioInitializer:  ScenarioInitialize,
		},
	})
}

//general/misc helpers:
func (p *probeState) runAISetupCheck(f func(bool) (bool, error), useDefaultNS bool, k string) error {

	b, err := f(useDefaultNS)

	if err != nil {
		return features.LogAndReturnError("error raised when checking for %v: %v", k, err)
	}

	if !b {
		return features.LogAndReturnError("%v does not exist (result: %t)", k, b)
	}

	return nil
}

//general feature steps:
func (p *probeState) aKubernetesClusterExistsWhichWeCanDeployInto() error {

	b := kubernetes.GetKubeInstance().ClusterIsDeployed()

	if b == nil || !*b {
		return features.LogAndReturnError("kubernetes cluster is NOT deployed")
	}

	//else we're good ...
	return nil
}

//AZ-AAD-AI-1.0
func (p *probeState) theDefaultNamespaceHasAnAzureIdentityBinding() error {
	return p.runAISetupCheck(iam.AzureIdentityBindingExists, true, "AzureIdentityBinding")
}
func (p *probeState) iCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding(namespace string) error {

	y, err := iamassets.Asset("assets/yaml/iam-azi-test-aib-curl.yaml")
	if err != nil {
		return features.LogAndReturnError("error reading yaml for test: %v", err)
	}

	if namespace == "the default" {
		p.useDefaultNS = true
	}

	pd, err := iam.CreateIAMTestPod(y, p.useDefaultNS)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.UndefinedPodCreationErrorReason, err)
}

//AZ-AAD-AI-1.0, AZ-AAD-AI-1.1
func (p *probeState) thePodIsDeployedSuccessfully() error {
	//check for pod name
	//note: the pod may still have a creation error if it didn't start up properly, but will have a name if the deployment succeeded
	//i.e.:
	// podName != "" -> successful deploy, potentially non-nil creation error
	// podName == "" -> unsuccessful deploy, non-nil creation error
	if p.state.PodName == "" {
		return features.LogAndReturnError("pod was not deployed successfully - creation error: %v", p.state.CreationError)
	}

	return nil
}
func (p *probeState) anAttemptToObtainAnAccessTokenFromThatPodShouldFail() error {
	//reuse the parameterised / scenario outline func
	return p.anAttemptToObtainAnAccessTokenFromThatPodShould("Fail")
}
func (p *probeState) anAttemptToObtainAnAccessTokenFromThatPodShould(expectedresult string) error {

	if p.state.CreationError != nil {
		return features.LogAndReturnError("failed to create pod", p.state.CreationError)
	}

	//curl for the auth token ... need to supply appropiate ns	
	res, err := iam.GetAccessToken(p.state.PodName, p.useDefaultNS)

	if err != nil {
		//this is an error from trying to execute the command as opposed to
		//the command itself returning an error
		return features.LogAndReturnError("error raised trying to execute auth token command - %v", err)
	}

	if expectedresult == "Fail" {
		if res != nil && len(*res) > 0 {
			//we got a token .. error
			return features.LogAndReturnError("token was successfully acquired on pod %v (result: %v)", p.state.PodName, *res)
		}
	} else if expectedresult == "Succeed" {
		if res == nil {
			//we didn't get a token .. error
			return features.LogAndReturnError("failed to acquire token on pod %v", p.state.PodName)
		}
	} else {
		return features.LogAndReturnError("unrecognised expected result: %v", expectedresult)
	}

	//else we're good
	return nil
}

//AZ-AAD-AI-1.1
func (p *probeState) theDefaultNamespaceHasAnAzureIdentity() error {
	return p.runAISetupCheck(iam.AzureIdentityExists, true, "AzureIdentity")
}
func (p *probeState) iCreateAnAzureIdentityBindingCalledInANondefaultNamespace(arg1 string) error {
	return p.runAISetupCheck(iam.AzureIdentityBindingExists, false, "AzureIdentityBinding")
}
func (p *probeState) iDeployAPodAssignedWithTheAzureIdentityBindingIntoTheSameNamespaceAsTheAzureIdentityBinding(arg1, arg2 string) error {
	y, err := iamassets.Asset("assets/yaml/iam-azi-test-aib-curl.yaml")
	if err != nil {
		return features.LogAndReturnError("error reading yaml for test: %v", err)
	}

	pd, err := iam.CreateIAMTestPod(y, false)
	return probe.ProcessPodCreationResult(&p.state, pd, kubernetes.UndefinedPodCreationErrorReason, err)
}

//AZ-AAD-AI-1.2
func (p *probeState) theClusterHasManagedIdentityComponentsDeployed() error {
	//look for the mic pods in the default ns
	pl, err := kubernetes.GetKubeInstance().GetPods("")

	if err != nil {
		return features.LogAndReturnError("error raised when trying to retrieve pods %v", err)
	}

	//a "pass" is the prescence of a "mic*" pod(s)
	//break on the first ...
	for _, pd := range pl.Items {
		if strings.HasPrefix(pd.Name, "mic-") {
			//grab the pod name as we'll execute the cmd against this:
			p.state.PodName = pd.Name
			return nil
		}
	}

	//fail if we get to here
	return features.LogAndReturnError("no MIC pods found - test fail")
}
func (p *probeState) iExecuteTheCommandAgainstTheMICPod(arg1 string) error {

	c := kubernetes.CatAzJSON
	res, err := iam.ExecuteVerificationCmd(p.state.PodName, c, true)

	if err != nil {
		//this is an error from trying to execute the command as opposed to
		//the command itself returning an error
		return features.LogAndReturnError("error raised trying to execute verification command (%v) - %v", c, err)
	}
	if res == nil {
		return features.LogAndReturnError("<nil> result received when trying to execute verification command (%v)", c)
	}
	if res.Err != nil && res.Internal {
		//we have an error which was raised before reaching the cluster (i.e. it's "internal")
		//this indicates that the command was not successfully executed
		return features.LogAndReturnError("error raised trying to execute verification command (%v)", c)
	}

	//otherwise, store the result code and return
	p.state.CommandExitCode = res.Code

	return nil
}
func (p *probeState) kubernetesShouldPreventMeFromRunningTheCommand() error {

	if p.state.CommandExitCode == 0 {
		//bad! don't want the command to succeed
		return features.LogAndReturnError("verification command was not blocked - test fail")
	}

	return nil
}

//setup, initialisation, etc.
func (p *probeState) setup() {

	//just make sure this is reset
	p.state.PodName = ""
	p.state.CreationError = nil
	p.useDefaultNS = false
}

func (p *probeState) tearDown() {
	iam.DeleteIAMTestPod(p.state.PodName, p.useDefaultNS)
}

// TestSuiteInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		//check dependancies ...
		if iam == nil {
			// not been given one so set default
			iam = kubernetes.NewDefaultIAM()
		}
		//setup AzureIdentity stuff ..??  Or should this be a pre-test setup
		// psp.CreateConfigMap()
	})

	ctx.AfterSuite(func() {
		//tear down AzureIdentity stuff?
		// psp.DeleteConfigMap()
	})
}

// ScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test 
// which reflects the tests described in the features directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example. 
func ScenarioInitialize(ctx *godog.ScenarioContext) {

	ps := probeState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		ps.setup()
		features.LogScenarioStart(s)
	})

	//general/all
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterExistsWhichWeCanDeployInto)

	//AZ-AAD-AI-1.0
	ctx.Step(`^the default namespace has an AzureIdentityBinding$`, ps.theDefaultNamespaceHasAnAzureIdentityBinding)
	ctx.Step(`^I create a simple pod in "([^"]*)" namespace assigned with that AzureIdentityBinding$`, ps.iCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding)

	//AZ-AAD-AI-1.0, AZ-AAD-AI-1.1
	ctx.Step(`^the pod is deployed successfully$`, ps.thePodIsDeployedSuccessfully)

	//AZ-AAD-AI-1.0
	ctx.Step(`^an attempt to obtain an access token from that pod should "([^"]*)"$`, ps.anAttemptToObtainAnAccessTokenFromThatPodShould)
	//AZ-AAD-AI-1.1 (same as above but just single shot scenario)
	ctx.Step(`^an attempt to obtain an access token from that pod should fail$`, ps.anAttemptToObtainAnAccessTokenFromThatPodShouldFail)

	//AZ-AAD-AI-1.1
	ctx.Step(`^the default namespace has an AzureIdentity$`, ps.theDefaultNamespaceHasAnAzureIdentity)
	ctx.Step(`^I create an AzureIdentityBinding called "([^"]*)" in a non-default namespace$`, ps.iCreateAnAzureIdentityBindingCalledInANondefaultNamespace)
	ctx.Step(`^I deploy a pod assigned with the "([^"]*)" AzureIdentityBinding into the same namespace as the "([^"]*)" AzureIdentityBinding$`, ps.iDeployAPodAssignedWithTheAzureIdentityBindingIntoTheSameNamespaceAsTheAzureIdentityBinding)

	//AZ-AAD-AI-1.2
	ctx.Step(`^I execute the command "([^"]*)" against the MIC pod$`, ps.iExecuteTheCommandAgainstTheMICPod)
	ctx.Step(`^Kubernetes should prevent me from running the command$`, ps.kubernetesShouldPreventMeFromRunningTheCommand)
	ctx.Step(`^the cluster has managed identity components deployed$`, ps.theClusterHasManagedIdentityComponentsDeployed)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		ps.tearDown()
		features.LogScenarioEnd(s)
	})
}
