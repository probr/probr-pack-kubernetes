package encryption_at_rest

import (
	"context"
	"log"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/summary"
	"github.com/cucumber/godog"
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

// EncryptionAtRestAzure azure implementation of the encryption in flight for Object Storage feature
type EncryptionAtRestAzure struct {
	ctx                       context.Context
	tags                      map[string]*string
	httpOption                bool
	httpsOption               bool
	policyAssignmentMgmtGroup string
}

var state EncryptionAtRestAzure

func (state *EncryptionAtRestAzure) securityControlsThatRestrictDataFromBeingUnencryptedAtRest() error {
	// It is available
	log.Printf("[DEBUG] Azure Storage account is encrypted by default and cannot be turned off. No test to run. Checking Azure Policy. (Unless customise this test to check for specific key usage.")
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) weProvisionAnObjectStorageBucket() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) encryptionAtRestIs(encryptionOption string) error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) creationWillWithAnErrorMatching(result string) error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) createContainerWithoutEncryption() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) detectiveDetectsNonCompliant() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) containerIsRemediated() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) setup() {
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) teardown() {
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) policyOrRuleAvailable() error {
	// It is available
	log.Printf("[DEBUG] Azure Storage account is encrypted by default and cannot be turned off. No test to run. Checking Azure Policy. (Unless customise this test to check for specific key usage.")
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) checkPolicyOrRuleAssignment() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) policyOrRuleAssigned() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) prepareToCreateContainer() error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) createContainerWithEncryptionOption(encryptionOption string) error {
	return nil
}

// PENDING IMPLEMENTATION
func (state *EncryptionAtRestAzure) createResult(result string) error {
	return nil
}

// Return this probe's name
func (p ProbeStruct) Name() string {
	return "encryption_at_rest"
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

	ctx.Step(`^security controls that restrict data from being unencrypted at rest$`, state.securityControlsThatRestrictDataFromBeingUnencryptedAtRest)
	ctx.Step(`^we provision an Object Storage bucket$`, state.weProvisionAnObjectStorageBucket)
	ctx.Step(`^encryption at rest is "([^"]*)"$`, state.encryptionAtRestIs)
	ctx.Step(`^creation will "([^"]*)" with an error matching "([^"]*)"$`, state.creationWillWithAnErrorMatching)

	ctx.Step(`^there is a detective capability for creation of Object Storage without encryption at rest$`, state.policyOrRuleAvailable)
	ctx.Step(`^the capability for detecting the creation of Object Storage without encryption at rest is active$`, state.checkPolicyOrRuleAssignment)
	ctx.Step(`^the detective measure is enabled$`, state.policyOrRuleAssigned)
	ctx.Step(`^Object Storage is created with without encryption at rest$`, state.createContainerWithoutEncryption)
	ctx.Step(`^the detective capability detects the creation of Object Storage without encryption at rest$`, state.detectiveDetectsNonCompliant)
	ctx.Step(`^the detective capability enforces encryption at rest on the Object Storage Bucket$`, state.containerIsRemediated)

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
