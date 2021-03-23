package clouddriverprobes

import (
	"github.com/cucumber/godog"
)

//This is the main implementation of the BDD/Cucumber feature test
//keep it in a separate file for clarity
//(as it will be utilised from the "test" file and the "test handler")

// PENDING IMPLEMENTATION
func aResourceCanBeDeployedIntoTheAccountUsingTheLinkedCredential(arg1, arg2 string) error {
	// return godog.ErrPending
	////log.Printf("[DEBUG] *** THEN: resource can be deployed ACCOUNT: %v | *** CREDENTIAL: %v \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem(arg1, arg2 string) error {
	// return godog.ErrPending
	////log.Printf("[DEBUG] *** AND: CREDENTIAL: %v with access to ACCOUNT: %v ALREADY EXISTS \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func iAddTheAccountDetailsToTheSystem(arg1 string) error {
	////log.Printf("[DEBUG] *** WHEN: ADD ACCOUNT: %v to the system \n", arg1)
	// return godog.ErrPending
	return nil
}

// PENDING IMPLEMENTATION
func iAmConfiguringAAccount(arg1 string) error {
	////log.Printf("[DEBUG] *** GIVEN: Configuring ACCOUNT: %v \n", arg1)
	// return godog.ErrPending
	return nil
}

// PENDING IMPLEMENTATION
func iLinkTheCredentialToTheAccount(arg1, arg2 string) error {
	// return godog.ErrPending
	////log.Printf("[DEBUG] *** CREDENTIAL: %v | *** ACCOUNT: %v \n", arg1, arg2)
	return nil
}

// PENDING IMPLEMENTATION
func aResourceDeploymentWillWithTheMessage(arg1, arg2 string) error {
	////log.Printf("[DEBUG] *** THEN: %v %v", arg1, arg2)
	return nil
}

//ProbeInitialize ...
func amProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {}) //nothing for now
}

//amScenarioInitialize ...
func amScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {}) //nothing for now

	ctx.Step(`^a resource can be deployed into the "([^"]*)" Account using the linked "([^"]*)" Credential$`, aResourceCanBeDeployedIntoTheAccountUsingTheLinkedCredential)
	ctx.Step(`^"([^"]*)" Credential with access to the "([^"]*)" Account is already configured in the system$`, credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem)
	ctx.Step(`^I add the "([^"]*)" Account details to the system$`, iAddTheAccountDetailsToTheSystem)
	ctx.Step(`^I am configuring a "([^"]*)" Account$`, iAmConfiguringAAccount)
	ctx.Step(`^I link the "([^"]*)" Credential to the "([^"]*)" Account$`, iLinkTheCredentialToTheAccount)
	ctx.Step(`^a resource deployment will "([^"]*)" with the message "([^"]*)"$`, aResourceDeploymentWillWithTheMessage)
}
