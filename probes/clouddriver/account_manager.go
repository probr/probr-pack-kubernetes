package clouddriver_probes

import (
	"log"

	"github.com/cucumber/godog"
)

//This is the main implementation of the BDD/Cucumber feature test
//keep it in a separate file for clarity
//(as it will be utilised from the "test" file and the "test handler")

// PENDING IMPLEMENTATION
func aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** THEN: resource can be deployed ACCOUNT: %v | *** CREDENTIAL: %v \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** AND: CREDENTIAL: %v with access to ACCOUNT: %v ALREADY EXISTS \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func iAddTheAccountDetailsToTheSystem(arg1 string) error {
	log.Printf("[INFO] *** WHEN: ADD ACCOUNT: %v to the system \n", arg1)
	// return godog.ErrPending
	return nil
}

// PENDING IMPLEMENTATION
func iAmConfiguringAAccount(arg1 string) error {
	log.Printf("[INFO] *** GIVEN: Configuring ACCOUNT: %v \n", arg1)
	// return godog.ErrPending
	return nil
}

// PENDING IMPLEMENTATION
func iLinkTheCredentialToTheAccount(arg1, arg2 string) error {
	// return godog.ErrPending
	log.Printf("[INFO] *** CREDENTIAL: %v | *** ACCOUNT: %v \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func aResouceDeploymentWillWithTheMessage(arg1, arg2 string) error {
	log.Printf("[INFO] *** THEN: %v %v", arg1, arg2)
	return nil
}

//ProbeInitialize ...
func amProbeInitialize(ctx *godog.TestSuiteContext) {
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
