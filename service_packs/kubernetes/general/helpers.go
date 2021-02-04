package general

import (
	"github.com/citihub/probr/internal/summary"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/cucumber/godog"
)

type scenarioState struct {
	name          string
	audit         *summary.ScenarioAudit
	probe         *summary.Probe
	podState      kubernetes.PodState
	wildcardRoles interface{}
}

func beforeScenario(s *scenarioState, probeName string, gs *godog.Scenario) {
	s.name = gs.Name
	s.probe = summary.State.GetProbeLog(probeName)
	s.audit = summary.State.GetProbeLog(probeName).InitializeAuditor(gs.Name, gs.Tags)
	coreengine.LogScenarioStart(gs)
}
