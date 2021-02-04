package audit

import (
	"github.com/cucumber/messages-go/v10"
)

type Probe struct {
	name               string
	audit              *ProbeAudit
	Meta               map[string]interface{}
	PodsCreated        int
	PodsDestroyed      int
	ScenariosAttempted int
	ScenariosSucceeded int
	ScenariosFailed    int
	Result             string
}

// CountPodCreated increments pods_created for probe
func (e *Probe) CountPodCreated(podName string) {
	State.LogPodName(podName)
	e.PodsCreated = e.PodsCreated + 1
}

// CountPodDestroyed increments pods_destroyed for probe
func (e *Probe) CountPodDestroyed() {
	e.PodsDestroyed = e.PodsDestroyed + 1
}

// countResults stores the current total number of failures as e.ScenariosFailed. Run at probe end
func (e *Probe) countResults() {
	e.ScenariosAttempted = len(e.audit.Scenarios)
	for _, v := range e.audit.Scenarios {
		if v.Result == "Failed" {
			e.ScenariosFailed = e.ScenariosFailed + 1
		} else if v.Result == "Passed" {
			e.ScenariosSucceeded = e.ScenariosSucceeded + 1
		}
	}
}

func (e *Probe) InitializeAuditor(name string, tags []*messages.Pickle_PickleTag) *ScenarioAudit {
	if e.audit.Scenarios == nil {
		e.audit.Scenarios = make(map[int]*ScenarioAudit)
	}
	scenarioCounter := len(e.audit.Scenarios) + 1
	var t []string
	for _, tag := range tags {
		t = append(t, tag.Name)
	}
	e.audit.Scenarios[scenarioCounter] = &ScenarioAudit{
		Name:  name,
		Steps: make(map[int]*StepAudit),
		Tags:  t,
	}
	return e.audit.Scenarios[scenarioCounter]
}
