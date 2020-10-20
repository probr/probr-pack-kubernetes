package summary

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"reflect"

	"github.com/citihub/probr/internal/config"
)

type SummaryState struct {
	Meta                  map[string]interface{}
	Status                string
	EventsPassed          int
	EventsFailed          int
	EventsSkipped         int
	names_of_pods_created []string
	EventTags             []config.Event // config.Event contains user-specified tagging options
	Events                map[string]*Event
}

var State SummaryState

func init() {
	State.Events = make(map[string]*Event)
	State.Meta = make(map[string]interface{})
	State.Meta["names_of_pods_created"] = new([]string)
}

// PrintSummary will print the current Events object state, formatted to JSON, if SummaryEnabled is not "false"
func (s *SummaryState) PrintSummary() {
	if config.Vars.SummaryEnabled == "false" {
		log.Printf("[NOTICE] Summary Log suppressed by configuration SummaryEnabled=false.")
	} else {
		summary, _ := json.MarshalIndent(s, "", "  ")
		fmt.Printf("%s", summary) // Summary output should not be handled by log levels
	}
}

// SetProbrStatus evaluates the current SummaryState state to set the Status
func (s *SummaryState) SetProbrStatus() {
	if s.EventsPassed > 0 && s.EventsFailed == 0 {
		s.Status = "Complete - All Events Completed Successfully"
	} else {
		s.Status = fmt.Sprintf("Complete - %v of %v Events Failed", s.EventsFailed, (len(s.Events) - s.EventsSkipped))
	}
	if config.Vars.Events != nil {
		s.Meta["event_tags_from_config"] = config.Vars.Events
	}
}

// LogEventMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (s *SummaryState) LogEventMeta(name string, key string, value interface{}) {
	e := s.GetEventLog(name)
	e.Meta[key] = value
	s.Events[name] = e
	s.Events[name].name = name // Event must be able to access its own name, but it is not publicly printed
}

// EventComplete takes an event name and status then updates the summary & event meta information
func (s *SummaryState) EventComplete(name string) {
	e := s.GetEventLog(name)
	s.completeEvent(e)
	e.audit.Write()
}

// GetEventLog initializes or returns existing log event for the provided test name
func (s *SummaryState) GetEventLog(n string) *Event {
	s.initEvent(n)
	return s.Events[n]
}

// LogPodName adds pod names to a list for user's debugging purposes
func (s *SummaryState) LogPodName(n string) {
	// A bit of effort is needed to keep this list in the generic "Meta"
	pn := reflect.ValueOf(s.Meta["names_of_pods_created"])
	var items []interface{}
	var result []string
	for i := 0; i < pn.Len(); i++ {
		items = append(items, pn.Index(i).Interface())
	}
	for _, v := range items {
		item := reflect.ValueOf(v)
		var record []string
		for i := 0; i < item.NumField(); i++ {
			itm := item.Field(i).Interface()
			record = append(record, fmt.Sprintf("%v", itm))
		}
		result = append(result, fmt.Sprintf("%v", record))
	}
	s.Meta["names_of_pods_created"] = result
}

func (s *SummaryState) initEvent(n string) {
	if s.Events[n] == nil {
		ap := filepath.Join(config.Vars.AuditDir, (n + ".json")) // Needed in both Event and EventAudit
		s.Events[n] = &Event{
			name:          n,
			Meta:          make(map[string]interface{}),
			PodsDestroyed: 0,
			audit: &EventAudit{
				Name: n,
				path: ap,
			},
		}
		s.Events[n].Meta["audit_path"] = ap // Meta is open for extension, any similar data can be stored there as needed

		// The event auditor should have pointers to the summary information
		s.Events[n].audit.PodsDestroyed = &s.Events[n].PodsDestroyed
		s.Events[n].audit.ProbesAttempted = &s.Events[n].ProbesAttempted
		s.Events[n].audit.ProbesSucceeded = &s.Events[n].ProbesSucceeded
		s.Events[n].audit.ProbesFailed = &s.Events[n].ProbesFailed
		s.Events[n].audit.Result = &s.Events[n].Result
	}
}

func (s *SummaryState) completeEvent(e *Event) {
	e.countResults()
	if e.Result == "Excluded" {
		e.Meta["audit_path"] = ""
		s.EventsSkipped = s.EventsSkipped + 1
	} else if len(e.audit.Probes) < 1 {
		e.Result = "No Probes Executed"
		e.Meta["audit_path"] = ""
		s.EventsSkipped = s.EventsSkipped + 1
	} else if e.ProbesFailed < 1 {
		e.Result = "Success"
		s.EventsPassed = s.EventsPassed + 1
	} else {
		e.Result = "Failed"
		s.EventsFailed = s.EventsFailed + 1
	}
}
