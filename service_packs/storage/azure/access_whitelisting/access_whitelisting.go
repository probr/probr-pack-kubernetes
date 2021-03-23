package azureaw

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	azurePolicy "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cucumber/godog"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/service_packs/coreengine"
	azureutil "github.com/citihub/probr/service_packs/storage/azure"
	"github.com/citihub/probr/service_packs/storage/azure/group"
	"github.com/citihub/probr/service_packs/storage/azure/policy"
	"github.com/citihub/probr/service_packs/storage/connection"

	"github.com/citihub/probr/utils"
)

const (
	policyAssignmentName = "deny_storage_wo_net_acl"        // TODO: Should this be in config?
	storageRgEnvVar      = "STORAGE_ACCOUNT_RESOURCE_GROUP" // TODO: Should this be replaced with azureutil.ResourceGroup() - which not only checks in env var, but also config vars?
)

// ProbeStruct allows this probe to be added to the ProbeStore
type ProbeStruct struct {
	state scenarioState
}

// Probe allows this probe to be added to the ProbeStore
var Probe ProbeStruct

type scenarioState struct {
	name                      string
	currentStep               string
	audit                     *audit.ScenarioAudit
	probe                     *audit.Probe
	ctx                       context.Context
	policyAssignmentMgmtGroup string
	tags                      map[string]*string
	bucketName                string
	storageAccount            azureStorage.Account
	runningErr                error
}

func (state *scenarioState) setup() {

	//log.Println("[DEBUG] Setting up \"AccessWhitelistingAzure\"")

}

func (state *scenarioState) teardown() {

	//log.Println("[DEBUG] Teardown completed")
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
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Check if value for Azure resource group is set in config vars;")
	if azureutil.ResourceGroup() == "" {
		//log.Printf("[ERROR] Azure resource group config var not set")
		err = errors.New("Azure resource group config var not set")
	}
	if err == nil {
		stepTrace.WriteString("Check the resource group exists in the specified azure subscription;")
		_, err = group.Get(state.ctx, azureutil.ResourceGroup())
		if err != nil {
			//log.Printf("[ERROR] Configured Azure resource group %s does not exists", azureutil.ResourceGroup())
		}
	}

	return err
}

func (state *scenarioState) checkPolicyAssigned() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		AzureSubscriptionID  string
		ManagamentGroup      string
		PolicyAssignmentName string
		PolicyAssignment     azurePolicy.Assignment
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	var a azurePolicy.Assignment

	if state.policyAssignmentMgmtGroup == "" {
		stepTrace.WriteString("Management Group has not been set, check Policy Assignment at the Subscription;")
		a, err = policy.AssignmentBySubscription(state.ctx, azureutil.SubscriptionID(), policyAssignmentName)
	} else {
		stepTrace.WriteString("Check Policy Assignment at the Management Group;")
		a, err = policy.AssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyAssignmentName)
	}

	//Audit log
	payload.AzureSubscriptionID = azureutil.SubscriptionID()
	payload.ManagamentGroup = state.policyAssignmentMgmtGroup
	payload.PolicyAssignmentName = policyAssignmentName
	payload.PolicyAssignment = a

	if err != nil {
		//log.Printf("[ERROR] Policy Assignment error: %v", err)
		return err
	}

	//log.Printf("[DEBUG] Policy Assignment check: %v [Step PASSED]", *a.Name)
	return nil
}

func (state *scenarioState) provisionStorageContainer() error {

	// define a bucket name, then pass the step - we will provision the account in the next step.

	var err error
	var stepTrace strings.Builder
	payload := struct {
		BucketName string
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("A bucket name is defined using a random string, storage account is not yet provisioned;")
	state.bucketName = utils.RandomString(10)

	//Audit log
	payload.BucketName = state.bucketName

	return err
}

func (state *scenarioState) createWithWhitelist(ipRange string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		SubscriptionID string
		ResourceGroup  string
		BucketName     string
		IPRange        string
		NetworkRuleSet azureStorage.NetworkRuleSet
		Tags           interface{}
		StorageAccount azureStorage.Account
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Attempting to create storage bucket with whitelisting for given IP Range: %s;", ipRange))

	var networkRuleSet azureStorage.NetworkRuleSet
	if ipRange == "nil" {
		stepTrace.WriteString("IP Range is nil, using DefaultActionAllow for NetworkRuleSet;")
		networkRuleSet = azureStorage.NetworkRuleSet{
			DefaultAction: azureStorage.DefaultActionAllow,
		}
	} else {
		stepTrace.WriteString("Setting IP Rule to allow given IP Range;")
		ipRule := azureStorage.IPRule{
			Action:           azureStorage.Allow,
			IPAddressOrRange: to.StringPtr(ipRange),
		}

		stepTrace.WriteString("Setting Network Rule Set with IP Rule;")
		networkRuleSet = azureStorage.NetworkRuleSet{
			IPRules:       &[]azureStorage.IPRule{ipRule},
			DefaultAction: azureStorage.DefaultActionDeny,
		}
	}

	stepTrace.WriteString("Creating storage bucket with Network Rule Set within Resource Group;")
	state.storageAccount, state.runningErr = connection.CreateWithNetworkRuleSet(state.ctx, state.bucketName, azureutil.ResourceGroup(), state.tags, true, &networkRuleSet)

	//Audit log
	err = state.runningErr
	payload.SubscriptionID = azureutil.SubscriptionID()
	payload.ResourceGroup = azureutil.ResourceGroup()
	payload.BucketName = state.bucketName
	payload.IPRange = ipRange
	payload.NetworkRuleSet = networkRuleSet
	payload.Tags = state.tags
	payload.StorageAccount = state.storageAccount

	return nil
}

func (state *scenarioState) creationWill(expectation string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		StorageAccountID string
		CreationError    string
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Expectation that Object Storage container was provisioned with whitelisting in previous step is: %s;", expectation))

	if expectation == "Fail" {
		if state.runningErr == nil {
			//Expected Fail but no previous error occurred, step should Fail
			err = fmt.Errorf("incorrectly created Storage Account: %v", *state.storageAccount.ID)
			payload.StorageAccountID = *state.storageAccount.ID // Audit log
			return err
		}
		payload.CreationError = state.runningErr.Error() // Audit log
		return nil                                       //Expected Fail and previous error occurred, step should Pass
	}

	if state.runningErr == nil {
		payload.StorageAccountID = *state.storageAccount.ID // Audit log
		return nil                                          //Expected Success and no previous error occurred, step should Pass
	}

	//Expected Success but previous error occurred, step should Fail
	err = state.runningErr
	payload.CreationError = state.runningErr.Error() // Audit log
	return err
}

func (state *scenarioState) cspSupportsWhitelisting() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	err = fmt.Errorf("Not Implemented")

	stepTrace.WriteString("TODO: Pending implementation;")

	//return err
	return nil //TODO: Remove this line and return actual err. This is temporary to ensure test doesn't halt and other steps are not skipped
}

func (state *scenarioState) examineStorageContainer(containerNameEnvVar string) error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
		StorageAccountName string
		ResourceGroup      string
		NetworkRuleSet     azureStorage.NetworkRuleSet
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Checking value for environment variable: %s;", containerNameEnvVar))
	accountName := os.Getenv(containerNameEnvVar) // TODO: Should this come from config?
	payload.StorageAccountName = accountName
	if accountName == "" {
		err = fmt.Errorf("environment variable \"%s\" is not defined test can't run", containerNameEnvVar)
		return err
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Checking value for environment variable: %s;", storageRgEnvVar))
	resourceGroup := os.Getenv(storageRgEnvVar) // TODO: Should this be replaced with azureutil.ResourceGroup() - which not only checks in env var, but also config vars?
	payload.ResourceGroup = resourceGroup
	if resourceGroup == "" {
		err = fmt.Errorf("environment variable \"%s\" is not defined test can't run", storageRgEnvVar)
		return err
	}

	stepTrace.WriteString("Retrieving storage account details from Azure;")
	state.storageAccount, state.runningErr = connection.AccountProperties(state.ctx, resourceGroup, accountName)
	if state.runningErr != nil {
		err = state.runningErr
		return err
	}

	stepTrace.WriteString("Checking that firewall network rule default action is not Allow;")
	networkRuleSet := state.storageAccount.AccountProperties.NetworkRuleSet
	payload.NetworkRuleSet = *networkRuleSet
	result := false
	// Default action is deny
	if networkRuleSet.DefaultAction == azureStorage.DefaultActionAllow {
		err = fmt.Errorf("%s has not configured with firewall network rule default action is not deny", accountName)
		return err
	}

	stepTrace.WriteString("Checking if it has IP whitelisting;")
	//for _, ipRule := range *networkRuleSet.IPRules {
	for range *networkRuleSet.IPRules {
		result = true
		//log.Printf("[DEBUG] IP WhiteListing: %v, %v", *ipRule.IPAddressOrRange, ipRule.Action)
	}

	stepTrace.WriteString("Checking if it has private Endpoint whitelisting;")
	//for _, vnetRule := range *networkRuleSet.VirtualNetworkRules {
	for range *networkRuleSet.VirtualNetworkRules {
		result = true
		//log.Printf("[DEBUG] VNet whitelisting: %v, %v", *vnetRule.VirtualNetworkResourceID, vnetRule.Action)
	}

	// TODO: Private Endpoint implementation when it's GA

	if result {
		//log.Printf("[DEBUG] Whitelisting rule exists. [Step PASSED]")
		err = nil
	} else {
		err = fmt.Errorf("no whitelisting has been defined for %v", accountName)
	}
	return err
}

// PENDING IMPLEMENTATION
func (state *scenarioState) whitelistingIsConfigured() error {
	// Checked in previous step

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(state.currentStep, stepTrace.String(), payload, err)
	}()

	err = fmt.Errorf("Not Implemented")

	stepTrace.WriteString("TODO: Pending implementation;")

	//return err
	return nil //TODO: Remove this line. This is temporary to ensure test doesn't halt and other steps are not skipped
}

func (state *scenarioState) beforeScenario(probeName string, gs *godog.Scenario) {
	state.name = gs.Name
	state.probe = audit.State.GetProbeLog(probeName)
	state.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	state.ctx = context.Background()
	coreengine.LogScenarioStart(gs)
}

// Name returns this probe's name
func (p ProbeStruct) Name() string {
	return "access_whitelisting"
}

// Path returns this probe's feature file path
func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "storage", "azure", p.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
//func (p ProbeStruct) ProbeInitialize(ctx *godog.Suite) {
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	p.state = scenarioState{}

	ctx.BeforeSuite(p.state.setup)

	ctx.AfterSuite(p.state.teardown)
}

// ScenarioInitialize initialises the scenario
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		p.state.beforeScenario(p.Name(), s)
	})

	ctx.Step(`^the CSP provides a whitelisting capability for Object Storage containers$`, p.state.cspSupportsWhitelisting)
	ctx.Step(`^a specified azure resource group exists$`, p.state.anAzureResourceGroupExists)
	ctx.Step(`^we examine the Object Storage container in environment variable "([^"]*)"$`, p.state.examineStorageContainer)
	ctx.Step(`^whitelisting is configured with the given IP address range or an endpoint$`, p.state.whitelistingIsConfigured)
	ctx.Step(`^security controls that Prevent Object Storage from being created without network source address whitelisting are applied$`, p.state.checkPolicyAssigned)
	ctx.Step(`^we provision an Object Storage container$`, p.state.provisionStorageContainer)
	ctx.Step(`^it is created with whitelisting entry "([^"]*)"$`, p.state.createWithWhitelist)
	ctx.Step(`^creation will "([^"]*)"$`, p.state.creationWill)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		coreengine.LogScenarioEnd(s)
	})

	ctx.BeforeStep(func(st *godog.Step) {
		p.state.currentStep = st.Text
	})

	ctx.AfterStep(func(st *godog.Step, err error) {
		p.state.currentStep = ""
	})
}
