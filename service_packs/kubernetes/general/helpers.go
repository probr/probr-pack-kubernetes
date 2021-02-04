package general

import (
	"github.com/citihub/probr/audit"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"
)

type scenarioState struct {
	name          string
	audit         *audit.ScenarioAudit
	probe         *audit.Probe
	podState      kubernetes.PodState
	wildcardRoles interface{}
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = audit.State.GetProbeLog(probeName)
	s.audit = audit.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}
