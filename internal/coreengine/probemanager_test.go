package coreengine

import (
	"testing"

	"github.com/citihub/probr/internal/config"
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

func TestTagIsExcluded(t *testing.T) {
	config.Vars.TagExclusions = []string{"tag_name"}
	if tagIsExcluded("not_tag_name") {
		t.Logf("Non-excluded tag was excluded")
		t.Fail()
	}
	if !tagIsExcluded("tag_name") {
		t.Logf("Excluded tag was not excluded")
		t.Fail()
	}
}

func TestIsExcluded(t *testing.T) {
	config.Vars.TagExclusions = []string{"excluded_probe"}
	pd := ProbeDescriptor{Group: Kubernetes, Name: "good_probe"}
	pd_excluded := ProbeDescriptor{Group: Kubernetes, Name: "excluded_probe"}

	if pd.isExcluded() {
		t.Logf("Non-excluded probe was excluded")
		t.Fail()
	}
	if !pd_excluded.isExcluded() {
		t.Logf("Excluded probe was not excluded")
		t.Fail()
	}
}

func TestAddProbe(t *testing.T) {
	probe_name := "test probe"
	excluded_probe_name := "different test probe"
	config.Vars.TagExclusions = []string{excluded_probe_name}
	ps := NewProbeStore()
	ps.AddProbe(createProbeObj(probe_name))
	ps.AddProbe(createProbeObj(excluded_probe_name))

	// Verify correct conditions succeed
	if ps.Probes[probe_name] == nil {
		t.Logf("Probe not added to probe store")
		t.Fail()
	} else if ps.Probes[probe_name].ProbeDescriptor.Name != probe_name {
		t.Logf("Probe name not set properly in test store")
		t.Fail()
	}

	// Verify probe1 and probe2 are different
	if ps.Probes[probe_name] == ps.Probes[excluded_probe_name] {
		t.Logf("Probes that should not match are equal to each other")
		t.Fail()
	}

	// Verify status is properly set
	if *ps.Probes[excluded_probe_name].Status != Excluded {
		t.Logf("Excluded probe was not excluded from probe store")
		t.Fail()
	}
	if *ps.Probes[probe_name].Status == Excluded {
		t.Logf("Excluded probe was not excluded from probe store")
		t.Fail()
	}
	// Note: this is not currently testing whether the summary or audit
	// are properly set for this because we may change how that is handled
	// without effecting probr functionality
}

func TestGetProbe(t *testing.T) {
	probe_name := "test probe"
	ps := NewProbeStore()
	probe := createProbeObj(probe_name)
	ps.AddProbe(probe)

	retrieved_probe, err := ps.GetProbe(probe_name)
	if err != nil {
		t.Logf(err.Error())
		t.Fail()
	}
	if retrieved_probe != probe {
		t.Logf("Retrieved probe does not match added probe")
		t.Fail()
	}
}

// Integration methods:
// TestExecProbe
// TestExecAllProbes
