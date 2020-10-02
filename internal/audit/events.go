package audit

import (
	"runtime"
	"strings"
)

type Probe struct {
	Steps map[string]string
}

type Event struct {
	Meta          map[string]string
	PodsCreated   int
	PodsDestroyed int
	ProbesFailed  int
	Probes        map[string]*Probe
}

// CountPodCreated increments PodsCreated for event
func (e *Event) CountPodCreated() {
	e.PodsCreated = e.PodsCreated + 1
}

// CountPodDestroyed increments PodsDestroyed for event
func (e *Event) CountPodDestroyed() {
	e.PodsDestroyed = e.PodsDestroyed + 1
}

// AuditProbe
func (e *Event) AuditProbe(name string, err error) {
	// get most recent caller function name as probe step
	f := make([]uintptr, 1)
	runtime.Callers(2, f)                      // add full caller path to empty object
	step := runtime.FuncForPC(f[0] - 1).Name() // get full caller path in string form
	s := strings.Split(step, ".")              // split full caller path
	step = s[len(s)-1]                         // select last element from caller path

	// Initialize any empty objects
	if e.Probes == nil {
		e.Probes = make(map[string]*Probe)
	}
	if e.Probes[name] == nil {
		e.Probes[name] = new(Probe)
		e.Probes[name].Steps = make(map[string]string)
	}

	// Now do the actual probe audit
	if err == nil {
		e.Probes[name].Steps[step] = "Passed"
	} else {
		e.Probes[name].Steps[step] = "Failed"
	}
}
