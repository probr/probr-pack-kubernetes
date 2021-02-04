package encryption_in_flight

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/azureutil"
	"github.com/citihub/probr/internal/azureutil/group"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/storage"
)

type scenarioState struct {
	name                      string
	audit                     *summary.ScenarioAudit
	probe                     *summary.Probe
	ctx                       context.Context
	tags                      map[string]*string
	httpOption                bool
	httpsOption               bool
	policyAssignmentMgmtGroup string
	storageAccounts           []string
}

// Allows this probe to be added to the ProbeStore
type ProbeStruct struct {
	state scenarioState
}

// Allows this probe to be added to the ProbeStore
var Probe ProbeStruct

func (state *scenarioState) setup() {

	log.Println("[DEBUG] Setting up \"scenarioState\"")

}

func (state *scenarioState) teardown() {
	for _, account := range state.storageAccounts {
		log.Printf("[DEBUG] need to delete the storageAccount: %s", account)
		err := storage.DeleteAccount(state.ctx, azureutil.ResourceGroup(), account)

		if err != nil {
			log.Printf("[ERROR] error deleting the storageAccount: %v", err)
		}
	}

	log.Println("[DEBUG] Teardown completed")
}

func (state *scenarioState) anAzureResourceGroupExists() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		AzureSubscriptionID string
		AzureResourceGroup  string
	}{
		AzureSubscriptionID: azureutil.SubscriptionID(),
		AzureResourceGroup:  azureutil.ResourceGroup(),
	}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Check if value for Azure resource group is set in config vars;")
	if azureutil.ResourceGroup() == "" {
		log.Printf("[ERROR] Azure resource group config var not set")
		err = errors.New("Azure resource group config var not set")
	}
	if err == nil {
		stepTrace.WriteString("Check the resource group exists in the specified azure subscription;")
		_, err = group.Get(state.ctx, azureutil.ResourceGroup())
		if err != nil {
			log.Printf("[ERROR] Configured Azure resource group %s does not exists", azureutil.ResourceGroup())
		}
	}
	return err
}

func (state *scenarioState) weProvisionAnObjectStorageBucket() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	// Nothing to do here
	return nil
}

func (state *scenarioState) httpAccessIs(arg1 string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf("Http Option: %s;", arg1))
	if arg1 == "enabled" {
		state.httpOption = true
	} else {
		state.httpOption = false
	}
	return nil
}

func (state *scenarioState) httpsAccessIs(arg1 string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf("Https Option: %s;", arg1))
	if arg1 == "enabled" {
		state.httpsOption = true
	} else {
		state.httpsOption = false
	}
	return nil
}

func (state *scenarioState) creationWillWithAnErrorMatching(expectation, errDescription string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		AccountName    string
		NetworkRuleSet azureStorage.NetworkRuleSet
		HTTPOption     bool
		HTTPSOption    bool
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Generating random value for account name;")
	accountName := utils.RandomString(5) + "storageac"
	payload.AccountName = accountName

	networkRuleSet := azureStorage.NetworkRuleSet{
		DefaultAction: azureStorage.DefaultActionDeny,
		IPRules:       &[]azureStorage.IPRule{},
	}
	payload.NetworkRuleSet = networkRuleSet
	payload.HTTPOption = state.httpOption
	payload.HTTPSOption = state.httpsOption

	// Both true take it as http option is try
	if state.httpsOption && state.httpOption {
		stepTrace.WriteString(fmt.Sprintf("Creating Storage Account with HTTPS: %v;", false))
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v;", false)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			azureutil.ResourceGroup(), state.tags, false, &networkRuleSet)
	} else if state.httpsOption {
		stepTrace.WriteString(fmt.Sprintf("Creating Storage Account with HTTPS: %v;", state.httpsOption))
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v", state.httpsOption)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			azureutil.ResourceGroup(), state.tags, state.httpsOption, &networkRuleSet)
	} else if state.httpOption {
		stepTrace.WriteString(fmt.Sprintf("Creating Storage Account with HTTPS: %v;", state.httpsOption))
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v", state.httpsOption)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			azureutil.ResourceGroup(), state.tags, state.httpsOption, &networkRuleSet)
	}
	if err == nil {
		// storage account created so add to state
		stepTrace.WriteString(fmt.Sprintf("Created Storage Account: %s;", accountName))
		log.Printf("[DEBUG] Created Storage Account: %s", accountName)
		state.storageAccounts = append(state.storageAccounts, accountName)
	}

	if expectation == "Fail" {

		if err == nil {
			err = fmt.Errorf("storage account was created, but should not have been: policy is not working or incorrectly configured")
			return err
		}

		detailedError := err.(autorest.DetailedError)
		originalErr := detailedError.Original
		detailed := originalErr.(*azure.ServiceError)

		log.Printf("[DEBUG] Detailed Error: %v", detailed)

		if strings.EqualFold(detailed.Code, "RequestDisallowedByPolicy") {
			stepTrace.WriteString("Request was Disallowed By Policy;")
			log.Printf("[DEBUG] Request was Disallowed By Policy: [Step PASSED]")
			return nil
		}

		err = fmt.Errorf("storage account was not created but not due to policy non-compliance")
		return err

	} else if expectation == "Succeed" {
		if err != nil {
			log.Printf("[ERROR] Unexpected failure in create storage ac [Step FAILED]")
			return err
		}
		return nil
	}

	err = fmt.Errorf("unsupported `result` option '%s' in the Gherkin feature - use either 'Fail' or 'Succeed'", expectation)
	return err
}

func (state *scenarioState) detectObjectStorageUnencryptedTransferAvailable() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return nil
}

func (state *scenarioState) detectObjectStorageUnencryptedTransferEnabled() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return nil
}

func (state *scenarioState) createUnencryptedTransferObjectStorage() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return nil
}

func (state *scenarioState) detectsTheObjectStorage() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return nil
}

func (state *scenarioState) encryptedDataTrafficIsEnforced() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return nil
}

func (s *scenarioState) beforeScenario(probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	s.ctx = context.Background()
	coreengine.LogScenarioStart(gs)
}

// Return this probe's name
func (p ProbeStruct) Name() string {
	return "encryption_in_flight"
}

func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "storage", "azure", p.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
//func (p ProbeStruct) ProbeInitialize(ctx *godog.Suite) {
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {

	ctx.BeforeSuite(p.state.setup)

	ctx.AfterSuite(p.state.teardown)
}

// initialises the scenario
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		p.state.beforeScenario(p.Name(), s)
	})

	ctx.Step(`^a specified azure resource group exists$`, p.state.anAzureResourceGroupExists)
	ctx.Step(`^we provision an Object Storage bucket$`, p.state.weProvisionAnObjectStorageBucket)
	ctx.Step(`^http access is "([^"]*)"$`, p.state.httpAccessIs)
	ctx.Step(`^https access is "([^"]*)"$`, p.state.httpsAccessIs)
	ctx.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, p.state.creationWillWithAnErrorMatching)

	ctx.Step(`^there is a detective capability for creation of Object Storage with unencrypted data transfer enabled$`, p.state.detectObjectStorageUnencryptedTransferAvailable)
	ctx.Step(`^the capability for detecting the creation of Object Storage with unencrypted data transfer enabled is active$`, p.state.detectObjectStorageUnencryptedTransferEnabled)
	ctx.Step(`^Object Storage is created with unencrypted data transfer enabled$`, p.state.createUnencryptedTransferObjectStorage)
	ctx.Step(`^the detective capability detects the creation of Object Storage with unencrypted data transfer enabled$`, p.state.detectsTheObjectStorage)
	ctx.Step(`^the detective capability enforces encrypted data transfer on the Object Storage Bucket$`, p.state.encryptedDataTrafficIsEnforced)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		coreengine.LogScenarioEnd(s)
	})
}
