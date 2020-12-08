// Provides the implementation required to execute the feature based test cases described in the
// the 'events' directory.  The 'assets' directory holds any assets required for the test cases.   Assets are 'embedded'
// via the 'go-bindata.exe' tool which is invoked via the 'go generate' tool.  It is important, therefore, that the
//'go:generate' comment is present in order to include this package in the scope of the 'go generate' tool.  This can be
// invoked directly on the command line of via the Makefile (e.g. make clean-build).
package iam

//go:generate go-bindata.exe -pkg $GOPACKAGE -o assets/iam/assets.go assets/iam/yaml probe_specifications/iamcontrol

import (
	"log"
	"strings"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/utils"
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
func (s *scenarioState) azureIdentitySetupCheck(f func(bool) (bool, error), useDefaultNS bool, k string) error {

	b, err := f(useDefaultNS)

	if err != nil {
		err = utils.ReformatError("error raised when checking for %v: %v", k, err)
		log.Print(err)
		return err
	}

	if !b {
		err = utils.ReformatError("%v does not exist (result: %t)", k, b)
		log.Print(err)
		return err
	}

	return nil
}

// General
func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	description, payload := kubernetes.ClusterIsDeployed()
	s.audit.AuditScenarioStep(description, payload, nil)
	return nil // ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

//AZ-AAD-AI-1.0
func (s *scenarioState) theDefaultNamespaceHasAnAzureIdentityBinding() error {
	err := s.azureIdentitySetupCheck(iam.AzureIdentityBindingExists, true, "AzureIdentityBinding")

	description := "Gets the AzureIdentityBindings, then filters according to namespace (default, if none supplied). Passes if binding is retrieved for namespace."
	s.audit.AuditScenarioStep(description, nil, err)

	return err

}

func (s *scenarioState) iCreateASimplePodInNamespaceAssignedWithThatAzureIdentityBinding(namespace string) error {
	description := ""
	var payload interface{}

	y, err := utils.ReadStaticFile(kubernetes.AssetsDir, "iam-azi-test-aib-curl.yaml")
	if err != nil {
		err = utils.ReformatError("error reading yaml for test: %v", err)
		log.Print(err)
	} else {
		if namespace == "the default" {
			s.useDefaultNS = true
		}
		pd, err := iam.CreateIAMProbePod(y, s.useDefaultNS, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.UndefinedPodCreationErrorReason, err)
	}

	s.audit.AuditScenarioStep(description, payload, err)

	return err

}

//AZ-AAD-AI-1.0, AZ-AAD-AI-1.1
func (s *scenarioState) thePodIsDeployedSuccessfully() error {
	//check for pod name
	//note: the pod may still have a creation error if it didn't start up properly, but will have a name if the deployment succeeded
	//i.e.:
	// podName != "" -> successful deploy, potentially non-nil creation error
	// podName == "" -> unsuccessful deploy, non-nil creation error
	var err error
	if s.podState.PodName == "" {
		err = utils.ReformatError("pod was not deployed successfully - creation error: %v", s.podState.CreationError)
	}

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) anAttemptToObtainAnAccessTokenFromThatPodShouldFail() error {
	//reuse the parameterised / scenario outline func
	err := s.anAttemptToObtainAnAccessTokenFromThatPodShould("Fail")

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) anAttemptToObtainAnAccessTokenFromThatPodShould(expectedresult string) error {
	var err error
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

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

//AZ-AAD-AI-1.1
func (s *scenarioState) theDefaultNamespaceHasAnAzureIdentity() error {
	err := s.azureIdentitySetupCheck(iam.AzureIdentityExists, true, "AzureIdentity")

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err

}

func (s *scenarioState) iCreateAnAzureIdentityBindingCalledInANondefaultNamespace(arg1 string) error {

	err := iam.CreateAIB()
	log.Printf("[DEBUG] error returned from CreateAIB: %v", err)

	return err
}

func (s *scenarioState) iDeployAPodAssignedWithTheAzureIdentityBindingIntoTheSameNamespaceAsTheAzureIdentityBinding(arg1, arg2 string) error {
	description := ""
	var payload interface{}

	y, err := utils.ReadStaticFile(kubernetes.AssetsDir, "iam-azi-test-aib-curl.yaml")
	if err != nil {
		err = utils.ReformatError("error reading yaml for test: %v", err)
		log.Print(err)
	} else {
		pd, err := iam.CreateIAMProbePod(y, false, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.UndefinedPodCreationErrorReason, err)
	}

	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

//AZ-AAD-AI-1.2
func (s *scenarioState) theClusterHasManagedIdentityComponentsDeployed() error {
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

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) iExecuteTheCommandAgainstTheMICPod(arg1 string) error {

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
		err = utils.ReformatError("error raised trying to execute verification command (%v)", c)
		log.Print(err)
	}
	if err != nil {
		// store the result code
		s.podState.CommandExitCode = res.Code
	}

	description := ""
	var payload interface{}
	s.audit.AuditScenarioStep(description, payload, err)

	return err
}

func (s *scenarioState) kubernetesShouldPreventMeFromRunningTheCommand() error {
	var err error
	if s.podState.CommandExitCode == 0 {
		//bad! don't want the command to succeed
		err = utils.ReformatError("verification command was not blocked")
	}

	description := "Examines scenario state to ensure that verification command was blocked."
	s.audit.AuditScenarioStep(description, nil, err)

	return err
}

func (p ProbeStruct) Name() string {
	return "iam"
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
		iam.DeleteIAMProbePod(ps.podState.PodName, ps.useDefaultNS, p.Name())
		coreengine.LogScenarioEnd(s)
	})
}
