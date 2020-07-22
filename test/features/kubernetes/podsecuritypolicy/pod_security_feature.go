package podsecuritypolicy

import (
	"fmt"
	"github.com/cucumber/godog"		
	"citihub.com/probr/internal/clouddriver/kubernetes"
)

func aDeploymentIsCreated() error {
	return godog.ErrPending
}

func accessIsRequested(arg1 string) error {
	return godog.ErrPending
}

func controlExistsToPreventPrivilegedAccess() error {
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

func creationWillWithAMessage(arg1, arg2 string) error {
	return godog.ErrPending
}

//TestSuiteInitialize ...
func TestSuiteInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func(){}) //nothing for now
}

//ScenarioInitialize ...
func ScenarioInitialize(ctx *godog.ScenarioContext) {
	ctx.BeforeScenario(func(*godog.Scenario) {}) //nothing for now

	ctx.Step(`^a deployment is created$`, aDeploymentIsCreated)
	ctx.Step(`^"([^"]*)" access is requested$`, accessIsRequested)
	ctx.Step(`^control exists to prevent privileged access$`, controlExistsToPreventPrivilegedAccess)
	ctx.Step(`^creation will "([^"]*)" with a message "([^"]*)"$`, creationWillWithAMessage)
}