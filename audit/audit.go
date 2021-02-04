package audit

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"
)

type ProbeAudit struct {
	path               string
	Name               string
	PodsDestroyed      *int
	ScenariosAttempted *int
	ScenariosSucceeded *int
	ScenariosFailed    *int
	Result             *string
	Scenarios          map[int]*ScenarioAudit
}

type ScenarioAudit struct {
	Name   string
	Result string // Passed / Failed / Given Not Met
	Tags   []string
	Steps  map[int]*StepAudit
}

type StepAudit struct {
	Name        string
	Description string      // Long-form explanation of anything happening in the step
	Result      string      // Passed / Failed
	Error       string      // Log the error text
	Payload     interface{} // Handles any values that are sent across the network
}

func (e *ProbeAudit) Write() {
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

// auditScenarioStep sets description, payload, and pass/fail based on err parameter.
// This function should be deferred to catch panic behavior, otherwise the audit will not be logged on panic
func (p *ScenarioAudit) AuditScenarioStep(description string, payload interface{}, err error) {
	stepName := utils.CallerName(2) // returns name if deferred and not panicking
	switch stepName {
	case "call":
		stepName = utils.CallerName(1) // returns name if this function was not deferred in the caller
	case "gopanic":
		stepName = utils.CallerName(3) // returns name if caller panicked and this function was deferred
	}
	p.audit(stepName, description, payload, err)
}

func (p *ScenarioAudit) audit(stepName string, description string, payload interface{}, err error) {
	stepNumber := len(p.Steps) + 1
	p.Steps[stepNumber] = &StepAudit{
		Name:        stepName,
		Description: description,
		Payload:     payload,
	}
	if err == nil {
		p.Steps[stepNumber].Result = "Passed"
		p.Result = "Passed"
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
func (e *ProbeAudit) probeRan() bool {
	if len(e.Scenarios) > 0 {
		return true
	}
	return false
}
