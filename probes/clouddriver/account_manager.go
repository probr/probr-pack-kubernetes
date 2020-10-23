package clouddriver_probes

import (
	"log"
	"path/filepath"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/cucumber/godog"
)

func init() {
	td := coreengine.TestDescriptor{Group: coreengine.CloudDriver, Name: "account_manager"}

	fp := filepath.Join("probes", "clouddriver", "probe_definitions")

	coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
		Handler: coreengine.GodogTestHandler,
		Data: &coreengine.GodogTest{
			TestDescriptor:       &td,
			TestSuiteInitializer: amTestSuiteInitialize,
			ScenarioInitializer:  amScenarioInitialize,
			FeaturePath:          &fp,
		},
	})
}

//This is the main implementation of the BDD/Cucumber feature test
//keep it in a separate file for clarity
//(as it will be utilised from the "test" file and the "test handler")

func aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** THEN: resource can be deployed ACCOUNT: %v | *** CREDENTIAL: %v \n", arg1, arg2)
	return nil
}

func credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** AND: CREDENTIAL: %v with access to ACCOUNT: %v ALREADY EXISTS \n", arg1, arg2)
	return nil
}

func iAddTheAccountDetailsToTheSystem(arg1 string) error {
	log.Printf("[INFO] *** WHEN: ADD ACCOUNT: %v to the system \n", arg1)
	// return godog.ErrPending
	return nil
}

func iAmConfiguringAAccount(arg1 string) error {
	log.Printf("[INFO] *** GIVEN: Configuring ACCOUNT: %v \n", arg1)
	// return godog.ErrPending
	return nil
}

func iLinkTheCredentialToTheAccount(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** CREDENTIAL: %v | *** ACCOUNT: %v \n", arg1, arg2)
	return nil
}

func aResouceDeploymentWillWithTheMessage(arg1, arg2 string) error {
	log.Printf("[INFO] *** THEN: %v %v", arg1, arg2)
	return nil
}

//TestSuiteInitialize ...
func amTestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now
}

//amScenarioInitialize ...
func amScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {}) //nothing for now

	ctx.Step(`^a resouce can be deployed into the "([^"]*)" Account using the linked "([^"]*)" Credential$`, aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential)
	ctx.Step(`^"([^"]*)" Credential with access to the "([^"]*)" Account is already configured in the system$`, credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem)
	ctx.Step(`^I add the "([^"]*)" Account details to the system$`, iAddTheAccountDetailsToTheSystem)
	ctx.Step(`^I am configuring a "([^"]*)" Account$`, iAmConfiguringAAccount)
	ctx.Step(`^I link the "([^"]*)" Credential to the "([^"]*)" Account$`, iLinkTheCredentialToTheAccount)
	ctx.Step(`^a resouce deployment will "([^"]*)" with the message "([^"]*)"$`, aResouceDeploymentWillWithTheMessage)
}
