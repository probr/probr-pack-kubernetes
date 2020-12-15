package access_whitelisting

import (
	"context"
	"fmt"
	"log"
	"os"

	azurePolicy "github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	azureStorage "github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/azureutil"
	"github.com/citihub/probr/internal/azureutil/group"
	"github.com/citihub/probr/internal/azureutil/policy"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/service_packs/storage"
)

const (
	policyAssignmentName = "deny_storage_wo_net_acl"
	storageRgEnvVar      = "STORAGE_ACCOUNT_RESOURCE_GROUP"
)

// Allows this probe to be added to the ProbeStore
type ProbeStruct struct{}

// Allows this probe to be added to the ProbeStore
var Probe ProbeStruct

type scenarioState struct {
	name  string
	audit *summary.ScenarioAudit
	probe *summary.Probe
	//	httpStatusCode int
	//	podName        string
	//	podState       kubernetes.PodState
	//	useDefaultNS   bool
	//	wildcardRoles  interface{}
}

type accessWhitelistingAzure struct {
	ctx                       context.Context
	policyAssignmentMgmtGroup string
	tags                      map[string]*string
	bucketName                string
	storageAccount            azureStorage.Account
	runningErr                error
}

var state accessWhitelistingAzure

func (state *accessWhitelistingAzure) setup() {

	log.Println("[DEBUG] Setting up 'accessWhitelistingAzure'")
	state.ctx = context.Background()

	state.policyAssignmentMgmtGroup = azureutil.ManagementGroup()
	if state.policyAssignmentMgmtGroup == "" {
		log.Printf("[NOTICE] Mgmt Group config variable is not defined. Policy assignment check against subscription")
	}

	state.tags = map[string]*string{
		"project": to.StringPtr("CICD"),
		"env":     to.StringPtr("test"),
		"tier":    to.StringPtr("internal"),
	}

	_, err := group.CreateWithTags(state.ctx, azureutil.ResourceGroup(), state.tags)
	if err != nil {
		log.Fatalf("failed to create group: %v\n", err.Error())
	}

	log.Printf("[DEBUG] Created Resource Group: %v", azureutil.ResourceGroup())
}

func (state *accessWhitelistingAzure) teardown() {
	err := group.Cleanup(state.ctx)
	if err != nil {
		log.Fatalf("Failed to teardown: %v\n", err.Error())
	}
	log.Println("[DEBUG] Teardown completed")
}

func (state *accessWhitelistingAzure) checkPolicyAssigned() error {

	var a azurePolicy.Assignment
	var err error

	// If a Management Group has not been set, check Policy Assignment at the Subscription
	if state.policyAssignmentMgmtGroup == "" {
		a, err = policy.AssignmentBySubscription(state.ctx, azureutil.SubscriptionID(), policyAssignmentName)
	} else {
		a, err = policy.AssignmentByManagementGroup(state.ctx, state.policyAssignmentMgmtGroup, policyAssignmentName)
	}

	if err != nil {
		log.Printf("[ERROR] Policy Assignment error: %v", err)
		return err
	}

	log.Printf("[DEBUG] Policy Assignment check: %v [Step PASSED]", *a.Name)
	return nil
}

func (state *accessWhitelistingAzure) provisionStorageContainer() error {
	// define a bucket name, then pass the step - we will provision the account in the next step.
	state.bucketName = azureutil.RandString(10)
	return nil
}

func (state *accessWhitelistingAzure) createWithWhitelist(ipRange string) error {
	var networkRuleSet azureStorage.NetworkRuleSet
	if ipRange == "nil" {
		networkRuleSet = azureStorage.NetworkRuleSet{
			DefaultAction: azureStorage.DefaultActionAllow,
		}
	} else {
		ipRule := azureStorage.IPRule{
			Action:           azureStorage.Allow,
			IPAddressOrRange: to.StringPtr(ipRange),
		}

		networkRuleSet = azureStorage.NetworkRuleSet{
			IPRules:       &[]azureStorage.IPRule{ipRule},
			DefaultAction: azureStorage.DefaultActionDeny,
		}
	}

	state.storageAccount, state.runningErr = storage.CreateWithNetworkRuleSet(state.ctx, state.bucketName, azureutil.ResourceGroup(), state.tags, true, &networkRuleSet)
	return nil
}

func (state *accessWhitelistingAzure) creationWill(expectation string) error {
	if expectation == "Fail" {
		if state.runningErr == nil {
			return fmt.Errorf("incorrectly created Storage Account: %v", *state.storageAccount.ID)
		}
		return nil
	}

	if state.runningErr == nil {
		return nil
	}

	return state.runningErr
}

func (state *accessWhitelistingAzure) cspSupportsWhitelisting() error {
	return nil
}

func (state *accessWhitelistingAzure) examineStorageContainer(containerNameEnvVar string) error {
	accountName := os.Getenv(containerNameEnvVar)
	if accountName == "" {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", containerNameEnvVar)
	}

	resourceGroup := os.Getenv(storageRgEnvVar)
	if resourceGroup == "" {
		return fmt.Errorf("environment variable \"%s\" is not defined test can't run", storageRgEnvVar)
	}

	state.storageAccount, state.runningErr = storage.AccountProperties(state.ctx, resourceGroup, accountName)

	if state.runningErr != nil {
		return state.runningErr
	}

	networkRuleSet := state.storageAccount.AccountProperties.NetworkRuleSet
	result := false
	// Default action is deny
	if networkRuleSet.DefaultAction == azureStorage.DefaultActionAllow {
		return fmt.Errorf("%s has not configured with firewall network rule default action is not deny", accountName)
	}

	// Check if it has IP whitelisting
	for _, ipRule := range *networkRuleSet.IPRules {
		result = true
		log.Printf("[DEBUG] IP WhiteListing: %v, %v", *ipRule.IPAddressOrRange, ipRule.Action)
	}

	// Check if it has private Endpoint whitelisting
	for _, vnetRule := range *networkRuleSet.VirtualNetworkRules {
		result = true
		log.Printf("[DEBUG] VNet whitelisting: %v, %v", *vnetRule.VirtualNetworkResourceID, vnetRule.Action)
	}

	// TODO: Private Endpoint implementation when it's GA

	if result {
		log.Printf("[DEBUG] Whitelisting rule exists. [Step PASSED]")
		return nil
	}
	return fmt.Errorf("no whitelisting has been defined for %v", accountName)
}

func (state *accessWhitelistingAzure) whitelistingIsConfigured() error {
	// Checked in previous step
	return nil
}

// Return this probe's name
func (p ProbeStruct) Name() string {
	return "access_whitelisting"
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

	ctx.Step(`^the CSP provides a whitelisting capability for Object Storage containers$`, state.cspSupportsWhitelisting)
	ctx.Step(`^we examine the Object Storage container in environment variable "([^"]*)"$`, state.examineStorageContainer)
	ctx.Step(`^whitelisting is configured with the given IP address range or an endpoint$`, state.whitelistingIsConfigured)
	ctx.Step(`^security controls that Prevent Object Storage from being created without network source address whitelisting are applied$`, state.checkPolicyAssigned)
	ctx.Step(`^we provision an Object Storage container$`, state.provisionStorageContainer)
	ctx.Step(`^it is created with whitelisting entry "([^"]*)"$`, state.createWithWhitelist)
	ctx.Step(`^creation will "([^"]*)"$`, state.creationWill)

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
