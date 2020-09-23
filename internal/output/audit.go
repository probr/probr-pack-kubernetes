package output

import (
	"encoding/json"
	"log"

	"gitlab.com/citihub/probr/internal/config"
)

type AuditLog struct {
	Events map[string]map[string]string
}

func (o *AuditLog) PrintAudit() {
	if config.Vars.AuditEnabled == "true" {
		audit, _ := json.MarshalIndent(o.Events, "", "  ")
		log.Printf("[NOTICE] %s", audit)
	} else {
		log.Printf("[NOTICE] Audit Log suppressed by configuration variable AuditLogEnabled.")
	}
}

// Audit accepts a test name with a key and value to insert to the logs for that test. Overwrites existing keys.
func (o *AuditLog) Audit(n string, k string, v string) {
	if o.Events == nil {
		o.Events = make(map[string]map[string]string)
	}
	l := o.GetEventLog(n)
	l[k] = v
	o.Events[n] = l
}

// GetEventLog initializes or returns existing log for the provided test name
func (o *AuditLog) GetEventLog(n string) map[string]string {
	if o.Events[n] == nil {
		o.Events[n] = make(map[string]string)
	}
	return o.Events[n]
}
