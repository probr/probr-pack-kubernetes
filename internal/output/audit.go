package output

type AuditLog struct {
	Events map[string]map[string]string
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
