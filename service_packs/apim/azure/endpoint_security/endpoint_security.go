package azurees

import (
	"context"
	"fmt"
	"strings"

	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/cucumber/godog"
)

type scenarioState struct {
	name        string
	currentStep string
	audit       *audit.ScenarioAudit
	probe       *audit.Probe
	ctx         context.Context
}

// ProbeStruct allows this probe to be added to the ProbeStore
type ProbeStruct struct {
	state scenarioState
}

// Probe allows this probe to be added to the ProbeStore
var Probe ProbeStruct

func (s *scenarioState) beforeScenario(probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}

// Name returns this probe's name
func (p ProbeStruct) Name() string {
	return "endpoint_security"
}

// Path returns the path to this probe's feature file
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
func (s *scenarioState) anAPIIsDeployedToAPIM() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		s.audit.AuditScenarioStep(s.currentStep, stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// PENDING IMPLEMENTATION
func (s *scenarioState) eachEndpointHasMTLSEmabled() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		s.audit.AuditScenarioStep(s.currentStep, stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// PENDING IMPLEMENTATION
func (s *scenarioState) allEndpointsAreRetrievedFromAPIM() error {

	var err error
	var stepTrace strings.Builder
	payload := struct {
	}{}
	defer func() {
		s.audit.AuditScenarioStep(s.currentStep, stepTrace.String(), payload, err)
	}()
	err = fmt.Errorf("Not Implemented")
	stepTrace.WriteString("TODO: Pending implementation;")

	return godog.ErrPending

}

// ScenarioInitialize initialises the scenario
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

	ctx.BeforeStep(func(st *godog.Step) {
		p.state.currentStep = st.Text
	})

	ctx.AfterStep(func(st *godog.Step, err error) {
		p.state.currentStep = ""
	})
}
