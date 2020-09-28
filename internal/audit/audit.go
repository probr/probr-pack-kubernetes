package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"gitlab.com/citihub/probr/internal/config"
)

type Event struct {
	Meta          map[string]string
	PodsCreated   int
	PodsDestroyed int
}

type AuditLogStruct struct {
	Status       string
	ProbesPassed int
	ProbesFailed int
	Events       map[string]*Event
}

var AuditLog AuditLogStruct

// PrintAudit will print the current Events object state, formatted to JSON, if AuditEnabled is not "false"
func (o *AuditLogStruct) PrintAudit() {
	if config.Vars.AuditEnabled == "false" {
		log.Printf("[NOTICE] Audit Log suppressed by configuration AuditEnabled=false.")
	} else {
		audit, _ := json.MarshalIndent(o, "", "  ")
		log.Printf("[NOTICE] %s", audit)
	}
}

// SetProbrStatus evaluates the current AuditLogStruct state to set ProbesPassed, ProbesFailed, and Status
func (o *AuditLogStruct) SetProbrStatus() {
	for _, v := range o.Events {
		if strings.Contains(v.Meta["status"], "Passed") {
			o.ProbesPassed = o.ProbesPassed + 1
		} else if strings.Contains(v.Meta["status"], "Failed") {
			o.ProbesFailed = o.ProbesFailed + 1
		}
	}
	if o.ProbesPassed > 0 && o.ProbesFailed == 0 {
		o.Status = "Completed - All Tests Passed"
	} else {
		o.Status = fmt.Sprintf("Completed - %v of %v Probes Failed", o.ProbesFailed, len(o.Events))
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
		o.Status = "Running"
	}
	if o.Events[n] == nil {
		o.Events[n] = &Event{
			Meta:          make(map[string]string),
			PodsCreated:   0,
			PodsDestroyed: 0,
		}
	}
}
