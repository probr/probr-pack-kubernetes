package audit

import (
	"encoding/json"
	"log"

	"gitlab.com/citihub/probr/internal/config"
)

type probeState struct {
	PodName         string
	CreationError   *interface{}
	ExpectedReason  *interface{}
	CommandExitCode int
}

type Event struct {
	Meta          map[string]string
	PodsCreated   int
	PodsDestroyed int
	Tests         map[string]probeState
}

type AuditLogStruct struct {
	Events map[string]*Event
}

var AuditLog AuditLogStruct

func (o *AuditLogStruct) PrintAudit() {
	if config.Vars.AuditEnabled == "false" {
		log.Printf("[NOTICE] Audit Log suppressed by configuration AuditEnabled=false.")
	} else {
		audit, _ := json.MarshalIndent(o.Events, "", "  ")
		log.Printf("[NOTICE] %s", audit)
	}
}

// AuditMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (o *AuditLogStruct) AuditMeta(name string, key string, value string) {
	e := o.GetEventLog(name)
	e.Meta[key] = value
	o.Events[name] = e
}

// GetEventLog initializes or returns existing log event for the provided test name
func (o *AuditLogStruct) GetEventLog(n string) *Event {
	o.logInit(n)
	return o.Events[n]
}

// Writes or updates given probe state
func (o *Event) AuditProbeState(n string, p probeState) {
	o.Tests[n] = p
}

func (o *Event) LogPodCreated() {
	o.PodsCreated = o.PodsCreated + 1
}

func (o *Event) LogPodDestroyed() {
	o.PodsDestroyed = o.PodsDestroyed + 1
}

// GetEventLog initializes log event if it doesn't already exist
func (o *AuditLogStruct) logInit(n string) {
	if o.Events == nil {
		o.Events = make(map[string]*Event)
	}
	if o.Events[n] == nil {
		o.Events[n] = &Event{
			Meta:          make(map[string]string),
			PodsCreated:   0,
			PodsDestroyed: 0,
		}
	}
}
