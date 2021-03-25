package audit

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/utils"
)

// ProbeAudit is used to hold all information related to probe execution
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

// ScenarioAudit is used by scenario states to audit progress through each step
type ScenarioAudit struct {
	Name   string
	Result string // Passed / Failed / Given Not Met
	Tags   []string
	Steps  map[int]*stepAudit
}

type stepAudit struct {
	Function    string
	Name        string
	Description string      // Long-form explanation of anything happening in the step
	Result      string      // Passed / Failed
	Error       string      // Log the error text
	Payload     interface{} // Handles any values that are sent across the network
}

func (e *ProbeAudit) Write() {
	if config.Vars.AuditEnabled == "true" && e.probeRan() {
		if utils.WriteAllowed(e.path, config.Vars.Overwrite()) {
			json, _ := json.MarshalIndent(e, "", "  ")
			data := []byte(json)
			ioutil.WriteFile(e.path, data, 0755)
		}
	}
}

// AuditScenarioStep sets description, payload, and pass/fail based on err parameter.
// This function should be deferred to catch panic behavior, otherwise the audit will not be logged on panic
func (p *ScenarioAudit) AuditScenarioStep(stepName, description string, payload interface{}, err error) {
	stepFunctionName := utils.CallerName(2) // returns name if deferred and not panicking
	switch stepFunctionName {
	case "call":
		stepFunctionName = utils.CallerName(1) // returns name if this function was not deferred in the caller
	case "gopanic":
		stepFunctionName = utils.CallerName(3) // returns name if caller panicked and this function was deferred
	}

	p.audit(stepFunctionName, stepName, description, payload, err)
}

func (p *ScenarioAudit) audit(functionName string, stepName string, description string, payload interface{}, err error) {
	// TODO: This function should replace audit. Added here to avoid breaking existing probes.
	stepNumber := len(p.Steps) + 1
	p.Steps[stepNumber] = &stepAudit{
		Function:    functionName,
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
