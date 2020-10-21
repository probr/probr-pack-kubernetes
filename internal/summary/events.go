package summary

import (
	"github.com/cucumber/messages-go/v10"
)

type Event struct {
	name            string
	audit           *EventAudit
	Meta            map[string]interface{}
	PodsCreated     int
	PodsDestroyed   int
	ProbesAttempted int
	ProbesSucceeded int
	ProbesFailed    int
	Result          string
}

// CountPodCreated increments pods_created for event
func (e *Event) CountPodCreated() {
	e.PodsCreated = e.PodsCreated + 1
}

// CountPodDestroyed increments pods_destroyed for event
func (e *Event) CountPodDestroyed() {
	e.PodsDestroyed = e.PodsDestroyed + 1
}

// countResults stores the current total number of failures as e.ProbesFailed. Run at event end
func (e *Event) countResults() {
	e.ProbesAttempted = len(e.audit.Probes)
	for _, v := range e.audit.Probes {
		if v.Result == "Failed" {
			e.ProbesFailed = e.ProbesFailed + 1
		} else if v.Result == "Passed" {
			e.ProbesSucceeded = e.ProbesSucceeded + 1
		}
	}
}

func (e *Event) InitializeAuditor(name string, tags []*messages.Pickle_PickleTag) *ProbeAudit {
	if e.audit.Probes == nil {
		e.audit.Probes = make(map[int]*ProbeAudit)
	}
	probeCounter := len(e.audit.Probes) + 1
	var t []string
	for _, tag := range tags {
		t = append(t, tag.Name)
	}
	e.audit.Probes[probeCounter] = &ProbeAudit{
		Name:  name,
		Steps: make(map[int]*StepAudit),
		Tags:  t,
	}
	return e.audit.Probes[probeCounter]
}
