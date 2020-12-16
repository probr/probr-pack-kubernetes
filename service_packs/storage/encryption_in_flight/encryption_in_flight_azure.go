package encryption_in_flight

import (
	"context"
	"fmt"
	"log"
	"strings"

	azurePolicy "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/azureutil"
	"github.com/citihub/probr/internal/azureutil/policy"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/service_packs/storage"
)

const (
	policyName = "deny_http_storage"
)

// Allows this probe to be added to the ProbeStore
type ProbeStruct struct{}

// Allows this probe to be added to the ProbeStore
var Probe ProbeStruct

type scenarioState struct {
	name  string
	audit *summary.ScenarioAudit
	probe *summary.Probe
}

// EncryptionInFlightAzure azure implementation of the encryption in flight for Object Storage feature
type EncryptionInFlightAzure struct {
	ctx                       context.Context
	tags                      map[string]*string
	httpOption                bool
	httpsOption               bool
	policyAssignmentMgmtGroup string
	resourceGroupName         string
}

var state EncryptionInFlightAzure

func (state *EncryptionInFlightAzure) setup() {

	log.Println("[DEBUG] Setting up \"EncryptionInFlightAzure\"")
	state.ctx = context.Background()

}

func (state *EncryptionInFlightAzure) teardown() {
	log.Println("[DEBUG] Teardown completed")
}

func (state *EncryptionInFlightAzure) securityControlsThatRestrictDataFromBeingUnencryptedInFlight() error {
	var policyAssignment azurePolicy.Assignment
	var aerr error
	// Search assignment from Management Group instead of subscription
	if state.policyAssignmentMgmtGroup != "" {
		policyAssignment, aerr = policy.AssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyName)
	} else {
		policyAssignment, aerr = policy.AssignmentBySubscription(state.ctx, azureutil.SubscriptionID(), policyName)
	}

	if aerr != nil {
		log.Printf("[ERROR] Get policy assignment error: %v", aerr)
		return aerr
	}

	log.Printf("[DEBUG] Policy assignment check: %v [Step PASSED]", *policyAssignment.Name)
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionInFlightAzure) anAzureResourceGroupExists() error {

	// check the resource group has been configured
	if config.Vars.CloudProviders.Azure.ResourceGroup == "" {
		log.Printf("[ERROR] Azure resource group config var not set")
	} else {
		log.Printf("[NOTICE] Azure resource group config var is %s", config.Vars.CloudProviders.Azure.ResourceGroup)
	}

	state.resourceGroupName = config.Vars.CloudProviders.Azure.ResourceGroup
	// Check the resource group exists in the specified azure subscription

	return nil
}

func (state *EncryptionInFlightAzure) weProvisionAnObjectStorageBucket() error {
	// Nothing to do here
	return nil
}

func (state *EncryptionInFlightAzure) httpAccessIs(arg1 string) error {
	if arg1 == "enabled" {
		state.httpOption = true
	} else {
		state.httpOption = false
	}
	return nil
}

func (state *EncryptionInFlightAzure) httpsAccessIs(arg1 string) error {
	if arg1 == "enabled" {
		state.httpsOption = true
	} else {
		state.httpsOption = false
	}
	return nil
}

func (state *EncryptionInFlightAzure) creationWillWithAnErrorMatching(expectation, errDescription string) error {
	accountName := azureutil.RandString(5) + "storageac"

	var err error

	networkRuleSet := azureStorage.NetworkRuleSet{
		DefaultAction: azureStorage.DefaultActionDeny,
		IPRules:       &[]azureStorage.IPRule{},
	}

	// Both true take it as http option is try
	if state.httpsOption && state.httpOption {
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v", false)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			state.resourceGroupName, state.tags, false, &networkRuleSet)
	} else if state.httpsOption {
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v", state.httpsOption)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			state.resourceGroupName, state.tags, state.httpsOption, &networkRuleSet)
	} else if state.httpOption {
		log.Printf("[DEBUG] Creating Storage Account with HTTPS: %v", state.httpsOption)
		_, err = storage.CreateWithNetworkRuleSet(state.ctx, accountName,
			state.resourceGroupName, state.tags, state.httpsOption, &networkRuleSet)
	}

	if expectation == "Fail" {

		if err == nil {
			return fmt.Errorf("storage account was created, but should not have been: policy is not working or incorrectly configured")
		}

		detailedError := err.(autorest.DetailedError)
		originalErr := detailedError.Original
		detailed := originalErr.(*azure.ServiceError)

		log.Printf("[DEBUG] Detailed Error: %v", detailed)

		if strings.EqualFold(detailed.Code, "RequestDisallowedByPolicy") {
			// Now check if it is the right policy
			if strings.Contains(detailed.Message, policyName) {
				log.Printf("[DEBUG] Request was Disallowed By Policy: %v [Step PASSED]", policyName)
				return nil
			}
			return fmt.Errorf("storage account was not created but blocked not by the right policy: %v", detailed.Message)
		}

		return fmt.Errorf("storage account was not created")
	} else if expectation == "Succeed" {
		if err != nil {
			log.Printf("[ERROR] Unexpected failure in create storage ac [Step FAILED]")
			return err
		}
		return nil
	}

	return fmt.Errorf("unsupported `result` option '%s' in the Gherkin feature - use either 'Fail' or 'Succeed'", expectation)
}

func (state *EncryptionInFlightAzure) detectObjectStorageUnencryptedTransferAvailable() error {
	return nil
}

func (state *EncryptionInFlightAzure) detectObjectStorageUnencryptedTransferEnabled() error {
	return nil
}

func (state *EncryptionInFlightAzure) createUnencryptedTransferObjectStorage() error {
	return nil
}

func (state *EncryptionInFlightAzure) detectsTheObjectStorage() error {
	return nil
}

func (state *EncryptionInFlightAzure) encryptedDataTrafficIsEnforced() error {
	return nil
}

// Return this probe's name
func (p ProbeStruct) Name() string {
	return "encryption_in_flight"
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
//func (p ProbeStruct) ProbeInitialize(ctx *godog.Suite) {
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {

	ctx.BeforeSuite(state.setup)

	ctx.AfterSuite(state.teardown)
}

// initialises the scenario
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	ctx.Step(`^a specified azure resource group exists$`, state.anAzureResourceGroupExists)
	ctx.Step(`^we provision an Object Storage bucket$`, state.weProvisionAnObjectStorageBucket)
	ctx.Step(`^http access is "([^"]*)"$`, state.httpAccessIs)
	ctx.Step(`^https access is "([^"]*)"$`, state.httpsAccessIs)
	ctx.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, state.creationWillWithAnErrorMatching)

	ctx.Step(`^there is a detective capability for creation of Object Storage with unencrypted data transfer enabled$`, state.detectObjectStorageUnencryptedTransferAvailable)
	ctx.Step(`^the capability for detecting the creation of Object Storage with unencrypted data transfer enabled is active$`, state.detectObjectStorageUnencryptedTransferEnabled)
	ctx.Step(`^Object Storage is created with unencrypted data transfer enabled$`, state.createUnencryptedTransferObjectStorage)
	ctx.Step(`^the detective capability detects the creation of Object Storage with unencrypted data transfer enabled$`, state.detectsTheObjectStorage)
	ctx.Step(`^the detective capability enforces encrypted data transfer on the Object Storage Bucket$`, state.encryptedDataTrafficIsEnforced)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		coreengine.LogScenarioEnd(s)
	})
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}
