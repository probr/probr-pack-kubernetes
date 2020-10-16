package summary

import (
	"github.com/cucumber/messages-go/v10"
)

type Event struct {
	name            string
	audit           *EventAudit
	Meta            map[string]interface{}
	PodsDestroyed   int
	ProbesFailed    int
	ProbesSucceeded int
	Status          string
}

// CountPodCreated increments pods_created for event
func (e *Event) CountPodCreated() {
	e.Meta["pods_created"] = e.Meta["pods_created"].(int) + 1
}

// CountPodDestroyed increments pods_destroyed for event
func (e *Event) CountPodDestroyed() {
	e.Meta["pods_destroyed"] = e.Meta["pods_destroyed"].(int) + 1
}

// countResults stores the current total number of failures as e.ProbesFailed. Run at event end
func (e *Event) countResults() {
	for _, v := range e.audit.Probes {
		if v.Result == "Failed" {
			e.ProbesFailed = e.ProbesFailed + 1
		} else {
			e.ProbesSucceeded = e.ProbesSucceeded + 1
		}
	}
}

func (e *Event) AuditProbeStep(name string, err error) {
	e.audit.logProbeStep(name, err)
}

func (e *Event) AuditProbeMeta(name string, tags []*messages.Pickle_PickleTag) {
	if e.audit.Probes == nil {
		e.audit.Probes = make(map[string]*ProbeAudit)
	}
	var t []string
	for _, v := range tags {
		t = append(t, v.Name)
	}
	e.audit.Probes[name] = &ProbeAudit{
		Steps: make(map[string]*StepAudit),
		Tags:  t,
	}
}
