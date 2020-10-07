package summary

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/citihub/probr/internal/config"
)

type SummaryStateStruct struct {
	Status        string
	EventsPassed  int
	EventsFailed  int
	EventsSkipped int
	PodNames      []string
	Events        map[string]*Event
	EventTags     []config.Event
}

var State SummaryStateStruct

// PrintSummary will print the current Events object state, formatted to JSON, if SummaryEnabled is not "false"
func (a *SummaryStateStruct) PrintSummary() {
	if config.Vars.SummaryEnabled == "false" {
		log.Printf("[NOTICE] Summary Log suppressed by configuration SummaryEnabled=false.")
	} else {
		summary, _ := json.MarshalIndent(a, "", "  ")
		fmt.Printf("%s", summary) // Summary output should not be handled by log levels
	}
}

// SetProbrStatus evaluates the current SummaryStateStruct state to set the Status
func (a *SummaryStateStruct) SetProbrStatus() {
	if a.EventsPassed > 0 && a.EventsFailed == 0 {
		a.Status = "Complete - All Events Completed Successfully"
	} else {
		a.Status = fmt.Sprintf("Complete - %v of %v Events Failed", a.EventsFailed, (len(a.Events) - a.EventsSkipped))
	}
	if config.Vars.Events != nil {
		a.EventTags = config.Vars.Events
	}
}

// LogEventMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (a *SummaryStateStruct) LogEventMeta(name string, key string, value string) {
	e := a.GetEventLog(name)
	e.Meta[key] = value
	a.Events[name] = e
}

// EventComplete takes an event name and status then updates the summary & event meta information
func (a *SummaryStateStruct) EventComplete(name string) {
	e := a.GetEventLog(name)
	e.CountFailures()
	if e.Meta["status"] == "Excluded" {
		a.EventsSkipped = a.EventsSkipped + 1
	} else if len(e.Probes) < 1 {
		e.Meta["status"] = "No Probes Executed"
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
func (a *SummaryStateStruct) GetEventLog(n string) *Event {
	a.logInit(n)
	return a.Events[n]
}

func (a *SummaryStateStruct) LogPodName(n string) {
	a.PodNames = append(a.PodNames, n)
}

// GetEventLog initializes log event if it doesn't already exist
func (a *SummaryStateStruct) logInit(n string) {
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
