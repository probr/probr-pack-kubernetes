package audit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/utils"
)

type summaryState struct {
	Meta          map[string]interface{}
	Status        string
	ProbesPassed  int
	ProbesFailed  int
	ProbesSkipped int
	Probes        map[string]*Probe
}

// State holds the values for all probe and scenario audits throughout the runtime
var State summaryState

func init() {
	State.Probes = make(map[string]*Probe)
	State.Meta = make(map[string]interface{})
	State.Meta["names of pods created"] = []string{}
}

// PrintSummary will print the current Probes object state, formatted to JSON, if NoSummary is not "true"
func (s *summaryState) PrintSummary() {
	if config.Vars.NoSummary == true {
		//log.Printf("[NOTICE] Summary Log suppressed by configuration NoSummary=true.")
	} else {
		//summary, _ := json.MarshalIndent(s, "", "  ")
		json.MarshalIndent(s, "", "  ")
		//log.Printf("Finished\n%s", summary) // Summary output should not be handled by log levels
	}
}

// WriteSummary will write the summary to the audit directory
func (s *summaryState) WriteSummary() {
	if config.Vars.AuditEnabled == "true" {
		path := filepath.Join(config.Vars.GetWriteDirectory(), "summary.json")
		if utils.WriteAllowed(path, config.Vars.Overwrite()) {
			json, _ := json.MarshalIndent(s, "", "  ")
			data := []byte(json)
			ioutil.WriteFile(path, data, 0755)
		}
	}
}

// SetProbrStatus evaluates the current summaryState state to set the Status
func (s *summaryState) SetProbrStatus() {
	if s.ProbesPassed > 0 && s.ProbesFailed == 0 {
		s.Status = "Complete - All Probes Completed Successfully"
	} else {
		s.Status = fmt.Sprintf("Complete - %v of %v Probes Failed", s.ProbesFailed, (len(s.Probes) - s.ProbesSkipped))
	}
}

// LogProbeMeta accepts a test name with a key and value to insert to the meta logs for that test. Overwrites key if already present.
func (s *summaryState) LogProbeMeta(name string, key string, value interface{}) {
	e := s.GetProbeLog(name)
	e.Meta[key] = value
	s.Probes[name] = e
	s.Probes[name].name = name // probe must be able to access its own name, but it is not publicly printed
}

// ProbeComplete takes an probe name and status then updates the summary & probe meta information
func (s *summaryState) ProbeComplete(name string) {
	e := s.GetProbeLog(name)
	s.completeProbe(e)
	e.audit.Write()
}

// GetProbeLog initializes or returns existing log probe for the provided test name
func (s *summaryState) GetProbeLog(n string) *Probe {
	s.initProbe(n)
	return s.Probes[n]
}

// LogPodName adds pod names to a list for user's debugging purposes
func (s *summaryState) LogPodName(n string) {
	podNames := s.Meta["names of pods created"].([]string)
	podNames = append(podNames, n)

	s.Meta["names of pods created"] = podNames
}

func (s *summaryState) initProbe(n string) {
	if s.Probes[n] == nil {
		ap := filepath.Join(config.Vars.AuditDir(), (n + ".json")) // Needed in both probe and ProbeAudit
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

func (s *summaryState) completeProbe(e *Probe) {
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
	if len(s.Meta["names of pods created"].([]string)) == 0 {
		s.Meta["pod creation error"] = "An error has occurred while creating pods. This may be due to a configuration error in the cluster, or in the config that was passed in to Probr"
	}
}
