package coreengine

import (
	"testing"
)

const (
	probeName = "good_probe"
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

func TestAddProbe(t *testing.T) {
	ps := NewProbeStore()
	ps.AddProbe(createProbeObj(probeName))

	// Verify correct conditions succeed
	if ps.Probes[probeName] == nil {
		t.Logf("Probe not added to probe store")
		t.Fail()
	} else if ps.Probes[probeName].ProbeDescriptor.Name != probeName {
		t.Logf("Probe name not set properly in test store")
		t.Fail()
	}
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
