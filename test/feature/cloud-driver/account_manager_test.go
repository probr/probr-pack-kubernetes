package driver

import (
	"fmt"
	"os"
	"flag"
	"testing"	
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opt = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

	// godog v0.10.0 (latest)
	status := godog.TestSuite{
		Name: "account_manager",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options: &opt,
	}.Run()

	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}



func aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential(arg1, arg2 string) error {
	// return godog.ErrPending
	fmt.Printf("*** THEN: resource can be deployed ACCOUNT: %v | *** CREDENTIAL: %v \n", arg1, arg2)
	return nil
}

func credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem(arg1, arg2 string) error {
	// return godog.ErrPending
	fmt.Printf("*** AND: CREDENTIAL: %v with access to ACCOUNT: %v ALREADY EXISTS \n", arg1, arg2)
	return nil
}

func iAddTheAccountDetailsToTheSystem(arg1 string) error {
	fmt.Printf("*** WHEN: ADD ACCOUNT: %v to the system \n", arg1)
	// return godog.ErrPending
	return nil
}

func iAmConfiguringAAccount(arg1 string) error {
	fmt.Printf("*** GIVEN: Configuring ACCOUNT: %v \n", arg1)
	// return godog.ErrPending
	return nil
}

func iLinkTheCredentialToTheAccount(arg1, arg2 string) error {
	// return godog.ErrPending
	fmt.Printf("*** CREDENTIAL: %v | *** ACCOUNT: %v \n", arg1, arg2)
	return nil
}

func aResouceDeploymentWillWithTheMessage(arg1, arg2 string) error {
	fmt.Printf("*** THEN: %v %v", arg1, arg2)
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a resouce can be deployed into the "([^"]*)" Account using the linked "([^"]*)" Credential$`, aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential)
	s.Step(`^"([^"]*)" Credential with access to the "([^"]*)" Account is already configured in the system$`, credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem)
	s.Step(`^I add the "([^"]*)" Account details to the system$`, iAddTheAccountDetailsToTheSystem)
	s.Step(`^I am configuring a "([^"]*)" Account$`, iAmConfiguringAAccount)
	s.Step(`^I link the "([^"]*)" Credential to the "([^"]*)" Account$`, iLinkTheCredentialToTheAccount)
	s.Step(`^a resouce deployment will "([^"]*)" with the message "([^"]*)"$`, aResouceDeploymentWillWithTheMessage)
}

// godog v0.10.0 (latest)
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func(){}) //nothing for now
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {}) //nothing for now

	ctx.Step(`^a resouce can be deployed into the "([^"]*)" Account using the linked "([^"]*)" Credential$`, aResouceCanBeDeployedIntoTheAccountUsingTheLinkedCredential)
	ctx.Step(`^"([^"]*)" Credential with access to the "([^"]*)" Account is already configured in the system$`, credentialWithAccessToTheAccountIsAlreadyConfiguredInTheSystem)
	ctx.Step(`^I add the "([^"]*)" Account details to the system$`, iAddTheAccountDetailsToTheSystem)
	ctx.Step(`^I am configuring a "([^"]*)" Account$`, iAmConfiguringAAccount)
	ctx.Step(`^I link the "([^"]*)" Credential to the "([^"]*)" Account$`, iLinkTheCredentialToTheAccount)
	ctx.Step(`^a resouce deployment will "([^"]*)" with the message "([^"]*)"$`, aResouceDeploymentWillWithTheMessage)
}