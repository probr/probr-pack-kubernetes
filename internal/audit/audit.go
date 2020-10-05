package audit

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/citihub/probr/internal/config"
)

type AuditLogStruct struct {
	Status        string
	EventsPassed  int
	EventsFailed  int
	EventsSkipped int
	PodNames      []string
	Events        map[string]*Event
}

var AuditLog AuditLogStruct

// PrintAudit will print the current Events object state, formatted to JSON, if AuditEnabled is not "false"
func (a *AuditLogStruct) PrintAudit() {
	if config.Vars.AuditEnabled == "false" {
		log.Printf("[NOTICE] Audit Log suppressed by configuration AuditEnabled=false.")
	} else {
		audit, _ := json.MarshalIndent(a, "", "  ")
		fmt.Printf("%s", audit) // Audit output should not be handled by log levels
	}
}

// SetProbrStatus evaluates the current AuditLogStruct state to set the Status
func (a *AuditLogStruct) SetProbrStatus() {
	if a.EventsPassed > 0 && a.EventsFailed == 0 {
		a.Status = "Complete - All Events Completed Successfully"
	} else {
		a.Status = fmt.Sprintf("Complete - %v of %v Events Failed", a.EventsFailed, (len(a.Events) - a.EventsSkipped))
	}
}

// AuditMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (a *AuditLogStruct) AuditMeta(name string, key string, value string) {
	e := a.GetEventLog(name)
	e.Meta[key] = value
	a.Events[name] = e
}

// EventComplete takes an event name and status then updates the audit & event meta information
func (a *AuditLogStruct) EventComplete(name string) {
	e := a.GetEventLog(name)
	e.CountFailures()
	if len(e.Probes) < 1 {
		e.Meta["status"] = "Skipped"
		a.EventsSkipped = a.EventsSkipped + 1
	} else if e.ProbesFailed < 1 {
		e.Meta["status"] = "Success"
		a.EventsPassed = a.EventsPassed + 1
	} else {
		e.Meta["status"] = "Failed"
		a.EventsFailed = a.EventsFailed + 1
	}
}

// GetEventLog initializes or returns existing log event for the provided test name
func (a *AuditLogStruct) GetEventLog(n string) *Event {
	a.logInit(n)
	return a.Events[n]
}

func (a *AuditLogStruct) AuditPodName(n string) {
	a.PodNames = append(a.PodNames, n)
}

// GetEventLog initializes log event if it doesn't already exist
func (a *AuditLogStruct) logInit(n string) {
	if a.Events == nil {
		a.Events = make(map[string]*Event)
		a.Status = "Running"
	}
	if a.Events[n] == nil {
		a.Events[n] = &Event{
			Meta:          make(map[string]string),
			PodsCreated:   0,
			PodsDestroyed: 0,
		}
	}
}
