package iam

import (
	"fmt"
	"log"
	"strings"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
)

type ProbeStruct struct{}

var Probe ProbeStruct

// IdentityAccessManagement is the section of the kubernetes package which provides the kubernetes interactions required to support
// identity access management scenarios.
var iam IdentityAccessManagement

// SetIAM allows injection of an IdentityAccessManagement helper.
func SetIAM(i IdentityAccessManagement) {
	iam = i
}

// azureIdentitySetupCheck executes the provided function and returns a formatted error
func (s *scenarioState) azureIdentitySetupCheck(f func(arg1 string, arg2 string) (bool, error), namespace, resourceType, resourceName string) error {

	b, err := f(namespace, resourceName)

	if err != nil {
		err = utils.ReformatError("error raised when checking for %v: %v", resourceType, err)
		log.Print(err)
		return err
	}

	if !b {
		err = utils.ReformatError("%v does not exist (result: %t)", resourceType, b)
		log.Print(err)
		return err
	}

	return nil
}

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload, err := kubernetes.ClusterIsDeployed()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()
	return err //  ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

//AZ-AAD-AI-1.0
func (s *scenarioState) aNamedAzureIdentityBindingExistsInNamedNS(aibName string, namespace string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	err = s.azureIdentitySetupCheck(iam.AzureIdentityBindingExists, namespace, "AzureIdentityBinding", aibName)

	description = "Gets the AzureIdentityBindings, then filters according to namespace. Passes if binding is retrieved for namespace."
	payload = struct {
		AzIdentbindName string
		AzNameSpace     string
	}{aibName, namespace}

	return err
}

func (s *scenarioState) iCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding(namespace, aibName string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	y, err := utils.ReadStaticFile(kubernetes.AssetsDir, "iam-azi-test-aib-curl.yaml")
	if err != nil {
		err = utils.ReformatError("error reading yaml for test: %v", err)
		log.Print(err)
	} else {
		if namespace == "the default" {
			s.useDefaultNS = true
		}
		pd, err := iam.CreateIAMProbePod(y, s.useDefaultNS, aibName, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.UndefinedPodCreationErrorReason, err)
	}

	description = fmt.Sprintf("Creating simple pod in %s namespace assigned with the azure identity binding %s", namespace, aibName)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

//AZ-AAD-AI-1.0, AZ-AAD-AI-1.1
func (s *scenarioState) thePodIsDeployedSuccessfully() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//check for pod name
	//note: the pod may still have a creation error if it didn't start up properly, but will have a name if the deployment succeeded
	//i.e.:
	// podName != "" -> successful deploy, potentially non-nil creation error
	// podName == "" -> unsuccessful deploy, non-nil creation error
	if s.podState.PodName == "" {
		err = utils.ReformatError("pod was not deployed successfully - creation error: %v", s.podState.CreationError)
	}

	description = "The Pod Deploying successfully"
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) anAttemptToObtainAnAccessTokenFromThatPodShouldFail() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//reuse the parameterised / scenario outline func
	err = s.anAttemptToObtainAnAccessTokenFromThatPodShould("Fail")

	description = "Attempting to get access token from Pod"
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) anAttemptToObtainAnAccessTokenFromThatPodShould(expectedresult string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	if s.podState.CreationError != nil {
		err = utils.ReformatError("failed to create pod", s.podState.CreationError)
		log.Print(err)
	} else {
		//curl for the auth token ... need to supply appropriate ns
		res, err := iam.GetAccessToken(s.podState.PodName, s.useDefaultNS)

		if err != nil {
			//this is an error from trying to execute the command as opposed to
			//the command itself returning an error
			err = utils.ReformatError("error raised trying to execute auth token command - %v", err)
			log.Print(err)
		} else {
			if expectedresult == "Fail" {
				if res != nil && len(*res) > 0 {
					//we got a token .. error
					err = utils.ReformatError("token was successfully acquired on pod %v (result: %v)", s.podState.PodName, *res)
				}
			} else if expectedresult == "Succeed" {
				if res == nil {
					//we didn't get a token .. error
					err = utils.ReformatError("failed to acquire token on pod %v", s.podState.PodName)
				}
			} else {
				err = utils.ReformatError("unrecognised expected result: %v", expectedresult)
				log.Print(err)
			}
		}
	}

	description = "Attempting to get access toekn from Pod"
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

//AZ-AAD-AI-1.1
func (s *scenarioState) aNamedAzureIdentityExistsInNamedNS(namespace string, aiName string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	err = s.azureIdentitySetupCheck(iam.AzureIdentityExists, namespace, "AzureIdentity", aiName)

	description = fmt.Sprintf("Gets the AzureIdentityBindings from binding %s, then filters according to namespace %s. Passes if binding is retrieved for namespace.", aiName, namespace)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) iCreateAnAzureIdentityBindingCalledInANondefaultNamespace(aibName, aiName string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	err = iam.CreateAIB(false, aibName, aiName) // create an AIB in a non-default NS if it deosn't already exist
	if err != nil {
		err = utils.ReformatError("error returned from CreateAIB: %v", err)
		log.Print(err)
	}
	description = fmt.Sprintf("Gets the AzureIdentityBindings from binding %s, then filters according to namespace %s. Passes if binding is retrieved for namespace.", aiName, aiName)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}
	return err
}

func (s *scenarioState) iDeployAPodAssignedWithTheAzureIdentityBindingIntoTheSameNamespaceAsTheAzureIdentityBinding(aibName string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	y, err := utils.ReadStaticFile(kubernetes.AssetsDir, "iam-azi-test-aib-curl.yaml")
	if err != nil {
		err = utils.ReformatError("error reading yaml for test: %v", err)
		log.Print(err)
	} else {
		pd, err := iam.CreateIAMProbePod(y, false, aibName, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.UndefinedPodCreationErrorReason, err)
	}
	description = fmt.Sprintf("Deploy Pod assigned with the AzureIdentityBindings from binding %s,", aibName)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

//AZ-AAD-AI-1.2
func (s *scenarioState) theClusterHasManagedIdentityComponentsDeployed() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	//look for the mic pods in the default ns
	pl, err := kubernetes.GetKubeInstance().GetPods("")

	if err != nil {
		err = utils.ReformatError("error raised when trying to retrieve pods %v", err)
	} else {
		//a "pass" is the prescence of a "mic*" pod(s)
		//break on the first ...
		for _, pd := range pl.Items {
			if strings.HasPrefix(pd.Name, "mic-") {
				//grab the pod name as we'll execute the cmd against this:
				s.podState.PodName = pd.Name
				err = nil
			}
		}
		if err != nil {
			err = utils.ReformatError("no MIC pods found - test fail")
		}
	}

	description = "Cluster has managed identity componenet deployed"
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) iExecuteTheCommandAgainstTheMICPod(arg1 string) error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	c := CatAzJSON
	res, err := iam.ExecuteVerificationCmd(s.podState.PodName, c, true)

	if err != nil {
		//this is an error from trying to execute the command as opposed to
		//the command itself returning an error
		err = utils.ReformatError("error raised trying to execute verification command (%v) - %v", c, err)
		log.Print(err)
	} else if res == nil {
		err = utils.ReformatError("<nil> result received when trying to execute verification command (%v)", c)
		log.Print(err)
	} else if res.Err != nil && res.Internal {
		//we have an error which was raised before reaching the cluster (i.e. it's "internal")
		//this indicates that the command was not successfully executed
		err = utils.ReformatError("%s: %v - (%v)", utils.CallerName(0), c, res.Err)
		log.Print(err)
	}
	if err != nil {
		// store the result code
		s.podState.CommandExitCode = res.Code
	}

	description = fmt.Sprintf("Executing command %s,", arg1)
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (s *scenarioState) kubernetesShouldPreventMeFromRunningTheCommand() error {
	// Standard auditing logic to ensures panics are also audited
	description, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(description, payload, err)
	}()

	if s.podState.CommandExitCode == 0 {
		//bad! don't want the command to succeed
		err = utils.ReformatError("verification command was not blocked")
	}

	description = "Examines scenario state to ensure that verification command was blocked."
	payload = struct {
		PodState kubernetes.PodState
		PodName  string
	}{s.podState, s.podState.PodName}

	return err
}

func (p ProbeStruct) Name() string {
	return "iam"
}

func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		//check dependancies ...
		if iam == nil {
			// not been given one so set default
			iam = NewDefaultIAM()
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
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	//general/all
	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterIsDeployed)

	//AZ-AAD-AI-1.0
	ctx.Step(`^an AzureIdentityBinding called "([^"]*)" exists in the namespace called "([^"]*)"$`, ps.aNamedAzureIdentityBindingExistsInNamedNS)

	ctx.Step(`^I create a simple pod in "([^"]*)" namespace assigned with the "([^"]*)" AzureIdentityBinding$`, ps.iCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding)

	//AZ-AAD-AI-1.0, AZ-AAD-AI-1.1
	ctx.Step(`^the pod is deployed successfully$`, ps.thePodIsDeployedSuccessfully)

	//AZ-AAD-AI-1.0
	ctx.Step(`^an attempt to obtain an access token from that pod should "([^"]*)"$`, ps.anAttemptToObtainAnAccessTokenFromThatPodShould)
	//AZ-AAD-AI-1.1 (same as above but just single shot scenario)
	ctx.Step(`^an attempt to obtain an access token from that pod should fail$`, ps.anAttemptToObtainAnAccessTokenFromThatPodShouldFail)

	//AZ-AAD-AI-1.1
	ctx.Step(`^the namespace called "([^"]*)" has an AzureIdentity called "([^"]*)"$`, ps.aNamedAzureIdentityExistsInNamedNS)
	ctx.Step(`^I create an AzureIdentityBinding called "([^"]*)" in the Probr namespace bound to the "([^"]*)" AzureIdentity$`, ps.iCreateAnAzureIdentityBindingCalledInANondefaultNamespace)
	ctx.Step(`^I deploy a pod assigned with the "([^"]*)" AzureIdentityBinding into the Probr namespace$`, ps.iDeployAPodAssignedWithTheAzureIdentityBindingIntoTheSameNamespaceAsTheAzureIdentityBinding)

	//AZ-AAD-AI-1.2
	ctx.Step(`^I execute the command "([^"]*)" against the MIC pod$`, ps.iExecuteTheCommandAgainstTheMICPod)
	ctx.Step(`^Kubernetes should prevent me from running the command$`, ps.kubernetesShouldPreventMeFromRunningTheCommand)
	ctx.Step(`^the cluster has managed identity components deployed$`, ps.theClusterHasManagedIdentityComponentsDeployed)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		iam.DeleteIAMProbePod(ps.podState.PodName, ps.useDefaultNS, p.Name())
		coreengine.LogScenarioEnd(s)
	})
}
