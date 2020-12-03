package coreengine

import (
	"testing"

	"github.com/citihub/probr/internal/config"
)

const (
	probeName         = "good_probe"
	excludedProbeName = "excluded_probe"
)

func createProbeObj(name string) *GodogProbe {
	return &GodogProbe{
		ProbeDescriptor: &ProbeDescriptor{
			Name:  name,
			Group: Kubernetes,
		},
	}
}

func TestNewProbeStore(t *testing.T) {
	ts := NewProbeStore()
	if ts == nil {
		t.Logf("Probe store was not initialized")
		t.Fail()
	} else if ts.Probes == nil {
		t.Logf("Probe store was not ready to add probes")
		t.Fail()
	}
}

func TestProbeIsExcluded(t *testing.T) {
	config.Vars.ProbeExclusions = []config.ProbeExclusion{config.ProbeExclusion{
		Name:          excludedProbeName,
		Excluded:      true,
		Justification: "testing",
	}}
	if probeIsExcluded(probeName) {
		t.Logf("Non-excluded probe was excluded")
		t.Fail()
	}
	if !probeIsExcluded(excludedProbeName) {
		t.Logf("Excluded probe was not excluded:\n%v", config.Vars.ProbeExclusions)
		t.Fail()
	}
}

func TestIsExcluded(t *testing.T) {
	config.Vars.ProbeExclusions = []config.ProbeExclusion{config.ProbeExclusion{
		Name:          excludedProbeName,
		Excluded:      true,
		Justification: "testing",
	}}
	pd := ProbeDescriptor{Group: Kubernetes, Name: probeName}
	pdExcluded := ProbeDescriptor{Group: Kubernetes, Name: excludedProbeName}

	if pd.isExcluded() {
		t.Logf("Non-excluded probe was excluded")
		t.Fail()
	}
	if !pdExcluded.isExcluded() {
		t.Logf("Excluded probe was not excluded")
		t.Fail()
	}
}

func TestAddProbe(t *testing.T) {
	config.Vars.ProbeExclusions = []config.ProbeExclusion{config.ProbeExclusion{
		Name:          excludedProbeName,
		Excluded:      true,
		Justification: "testing",
	}}
	ps := NewProbeStore()
	ps.AddProbe(createProbeObj(probeName))
	ps.AddProbe(createProbeObj(excludedProbeName))

	// Verify correct conditions succeed
	if ps.Probes[probeName] == nil {
		t.Logf("Probe not added to probe store")
		t.Fail()
	} else if ps.Probes[probeName].ProbeDescriptor.Name != probeName {
		t.Logf("Probe name not set properly in test store")
		t.Fail()
	}

	// Verify probe1 and probe2 are different
	if ps.Probes[probeName] == ps.Probes[excludedProbeName] {
		t.Logf("Probes that should not match are equal to each other")
		t.Fail()
	}

	// Verify status is properly set
	if *ps.Probes[excludedProbeName].Status != Excluded {
		t.Logf("Excluded probe was not excluded from probe store")
		t.Fail()
	}
	if *ps.Probes[probeName].Status == Excluded {
		t.Logf("Excluded probe was not excluded from probe store")
		t.Fail()
	}
	// Note: this is not currently testing whether the summary or audit
	// are properly set for this because we may change how that is handled
	// without effecting probr functionality
}

func TestGetProbe(t *testing.T) {
	ps := NewProbeStore()
	probe := createProbeObj(probeName)
	ps.AddProbe(probe)

	retrievedProbe, err := ps.GetProbe(probeName)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
	if retrievedProbe != probe {
		t.Logf("Retrieved probe does not match added probe")
		t.Fail()
	}
}

// Integration methods:
// TestExecProbe
// TestExecAllProbes
