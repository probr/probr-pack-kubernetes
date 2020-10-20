package summary

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/citihub/probr/internal/config"
)

type EventAudit struct {
	path            string
	Name            string
	PodsDestroyed   *int
	ProbesAttempted *int
	ProbesSucceeded *int
	ProbesFailed    *int
	Result          *string
	Probes          map[int]*ProbeAudit
}

type ProbeAudit struct {
	Name   string
	Result string // Passed / Failed / Given Not Met
	Tags   []string
	Steps  map[int]*StepAudit
}

type StepAudit struct {
	Name        string
	Description string      // Long-form exlanation of anything happening in the step
	Result      string      // Passed / Failed
	Error       string      // Log the error text
	Payload     interface{} // Handles any values that are sent across the network
}

func (e *EventAudit) Write() {
	if config.Vars.AuditEnabled == "true" && e.probeRan() {
		_, err := os.Stat(e.path)
		if err == nil && config.Vars.OverwriteHistoricalAudits == "false" {
			// Historical audits are preserved by default
			log.Fatalf("[ERROR] AuditEnabled is set to true, but audit file already exists or Probr could not open: '%s'", e.path)
		}

		json, _ := json.MarshalIndent(e, "", "  ")
		data := []byte(json)
		err = ioutil.WriteFile(e.path, data, 0755)

		if err != nil {
			log.Fatalf("[ERROR] AuditEnabled is set to true, but Probr could not write audit to file: '%s'", e.path)
		}
	}
}

// auditProbeStep sets description, payload, and pass/fail based on err parameter
func (p *ProbeAudit) AuditProbeStep(description string, payload interface{}, err error) {
	// Initialize any empty objects
	// Now do the actual probe summary
	stepName := getCallerName(3)
	stepNumber := len(p.Steps) + 1
	p.Steps[stepNumber] = &StepAudit{
		Name:        stepName,
		Description: description,
		Payload:     payload,
	}
	if err == nil {
		p.Steps[stepNumber].Result = "Passed"
	} else {
		p.Steps[stepNumber].Result = "Failed"
		p.Steps[stepNumber].Error = strings.Replace(err.Error(), "[ERROR] ", "", -1)
		if stepNumber == 1 {
			p.Result = "Given Not Met" // First entry is always a 'given'; failures should be ignored
		} else {
			p.Result = "Failed" // First 'given' was met, but a subsequent step failed
		}
	}
}

// getCallerName retrieves the name of the function prior to the location it is called
func getCallerName(up int) string {
	f := make([]uintptr, 1)
	runtime.Callers(up, f)                     // add full caller path to empty object
	step := runtime.FuncForPC(f[0] - 1).Name() // get full caller path in string form
	s := strings.Split(step, ".")              // split full caller path
	return s[len(s)-1]                         // select last element from caller path
}

func (e *EventAudit) probeRan() bool {
	if len(e.Probes) > 0 {
		return true
	}
	return false
}
