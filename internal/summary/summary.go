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
	ProbesPassed          int
	ProbesFailed          int
	ProbesSkipped         int
	names_of_pods_created []string
	ProbeTags             []config.Probe // config.Probe contains user-specified tagging options
	Probes                map[string]*Probe
}

var State SummaryState

func init() {
	State.Probes = make(map[string]*Probe)
	State.Meta = make(map[string]interface{})
	State.Meta["names_of_pods_created"] = new([]string)
}

// PrintSummary will print the current Probes object state, formatted to JSON, if SummaryEnabled is not "false"
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
	if s.ProbesPassed > 0 && s.ProbesFailed == 0 {
		s.Status = "Complete - All Probes Completed Successfully"
	} else {
		s.Status = fmt.Sprintf("Complete - %v of %v Probes Failed", s.ProbesFailed, (len(s.Probes) - s.ProbesSkipped))
	}
	if config.Vars.Probes != nil {
		s.Meta["probe_tags_from_config"] = config.Vars.Probes
	}
}

// LogProbeMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (s *SummaryState) LogProbeMeta(name string, key string, value interface{}) {
	e := s.GetProbeLog(name)
	e.Meta[key] = value
	s.Probes[name] = e
	s.Probes[name].name = name // Probe must be able to access its own name, but it is not publicly printed
}

// ProbeComplete takes an probe name and status then updates the summary & probe meta information
func (s *SummaryState) ProbeComplete(name string) {
	e := s.GetProbeLog(name)
	s.completeProbe(e)
	e.audit.Write()
}

// GetProbeLog initializes or returns existing log probe for the provided test name
func (s *SummaryState) GetProbeLog(n string) *Probe {
	s.initProbe(n)
	return s.Probes[n]
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

func (s *SummaryState) initProbe(n string) {
	if s.Probes[n] == nil {
		ap := filepath.Join(config.Vars.AuditDir, (n + ".json")) // Needed in both Probe and ProbeAudit
		s.Probes[n] = &Probe{
			name:          n,
			Meta:          make(map[string]interface{}),
			PodsDestroyed: 0,
			audit: &ProbeAudit{
				Name: n,
				path: ap,
			},
		}
		s.Probes[n].Meta["audit_path"] = ap // Meta is open for extension, any similar data can be stored there as needed

		// The probe auditor should have pointers to the summary information
		s.Probes[n].audit.PodsDestroyed = &s.Probes[n].PodsDestroyed
		s.Probes[n].audit.ScenariosAttempted = &s.Probes[n].ScenariosAttempted
		s.Probes[n].audit.ScenariosSucceeded = &s.Probes[n].ScenariosSucceeded
		s.Probes[n].audit.ScenariosFailed = &s.Probes[n].ScenariosFailed
		s.Probes[n].audit.Result = &s.Probes[n].Result
	}
}

func (s *SummaryState) completeProbe(e *Probe) {
	e.countResults()
	if e.Result == "Excluded" {
		e.Meta["audit_path"] = ""
		s.ProbesSkipped = s.ProbesSkipped + 1
	} else if len(e.audit.Scenarios) < 1 {
		e.Result = "No Scenarios Executed"
		e.Meta["audit_path"] = ""
		s.ProbesSkipped = s.ProbesSkipped + 1
	} else if e.ScenariosFailed < 1 {
		e.Result = "Success"
		s.ProbesPassed = s.ProbesPassed + 1
	} else {
		e.Result = "Failed"
		s.ProbesFailed = s.ProbesFailed + 1
	}
}
