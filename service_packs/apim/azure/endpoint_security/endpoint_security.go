package endpoint_security

import (
	"context"
	"fmt"
	"strings"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/cucumber/godog"
)

type scenarioState struct {
	name  string
	audit *audit.ScenarioAudit
	probe *audit.Probe
	ctx   context.Context
}

// Allows this probe to be added to the ProbeStore
type ProbeStruct struct {
	state scenarioState
}

// Allows this probe to be added to the ProbeStore
var Probe ProbeStruct

func (s *scenarioState) beforeScenario(probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}

// Return this probe's name
func (p ProbeStruct) Name() string {
	return "endpoint_security"
}

func (p ProbeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "apim", "azure", p.Name())
}

// ProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
//func (p ProbeStruct) ProbeInitialize(ctx *godog.Suite) {
func (p ProbeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	p.state = scenarioState{}

	//	ctx.BeforeSuite(p.state.setup)

	//	ctx.AfterSuite(p.state.teardown)
}

// PENDING IMPLEMENTATION
func (state *scenarioState) anAPIIsDeployedToAPIM() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// PENDING IMPLEMENTATION
func (state *scenarioState) eachEndpointHasMTLSEmabled() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// PENDING IMPLEMENTATION
func (state *scenarioState) allEndpointsAreRetrievedFromAPIM() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		state.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// initialises the scenario
func (p ProbeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {

	ctx.BeforeScenario(func(s *godog.Scenario) {
		p.state.beforeScenario(p.Name(), s)
	})

	ctx.Step(`^an API that is deployed to APIM$`, p.state.anAPIIsDeployedToAPIM)
	ctx.Step(`^all endpoints are retrieved from APIM$`, p.state.allEndpointsAreRetrievedFromAPIM)
	ctx.Step(`^each endpoint has mTLS enabled$`, p.state.eachEndpointHasMTLSEmabled)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		coreengine.LogScenarioEnd(s)
	})
}
