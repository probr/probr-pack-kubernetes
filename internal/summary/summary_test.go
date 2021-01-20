package summary

import (
	"path/filepath"
	"testing"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"
)

func TestSummaryState_LogPodName(t *testing.T) {

	var fakeSummaryState SummaryState
	fakeSummaryState.Probes = make(map[string]*Probe)
	fakeSummaryState.Meta = make(map[string]interface{})
	fakeSummaryState.Meta["names of pods created"] = []string{}

	type args struct {
		podName string
	}
	tests := []struct {
		testName string
		s        *SummaryState
		args     args
	}{
		{
			testName: "MetaShouldContainPodName",
			s:        &fakeSummaryState,
			args:     args{podName: "testPod"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.s.LogPodName(tt.args.podName)
			tt.s.LogPodName("Anotherpod")
			loggedPods := fakeSummaryState.Meta["names of pods created"].([]string)
			actualPosition, actualFound := utils.FindString(loggedPods, tt.args.podName)
			if !(actualPosition >= 0 && actualFound == true) {
				t.Errorf("State.Meta doesn't contain pod name: %v", loggedPods)
			}
		})
	}
}

// createMockProbe - creates a mock summarystate and probe object in it and returns summarySteate object.
func createSummaryStateWithMockProbe(probename string) SummaryState {
	var sumstate SummaryState
	sumstate.Probes = make(map[string]*Probe)
	ap := filepath.Join(config.AuditDir(), (probename + ".json")) // Needed in both Probe and ProbeAudit
	sumstate.Probes[probename] = &Probe{
		name:          probename,
		Meta:          make(map[string]interface{}),
		PodsDestroyed: 0,
		audit: &ProbeAudit{
			Name: probename,
			path: ap,
		},
	}
	sumstate.Probes[probename].Meta["audit_path"] = ap
	sumstate.Probes[probename].audit.PodsDestroyed = &sumstate.Probes[probename].PodsDestroyed
	sumstate.Probes[probename].audit.ScenariosAttempted = &sumstate.Probes[probename].ScenariosAttempted
	sumstate.Probes[probename].audit.ScenariosFailed = &sumstate.Probes[probename].ScenariosFailed
	sumstate.Probes[probename].audit.Result = &sumstate.Probes[probename].Result
	sumstate.Probes[probename].countResults()
	return sumstate
}

func TestSummaryState_initProbe(t *testing.T) {

	type args struct {
		fakeName string
	}
	var probeName = "testProbe"
	var mockSummaryState = createSummaryStateWithMockProbe(probeName)

	tests := []struct {
		testName string
		s        *SummaryState
		args     args
	}{
		{
			testName: "TestProbeInitialized",
			s:        &mockSummaryState,
			args:     args{fakeName: "testProbe"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.s.initProbe(tt.args.fakeName)
			tt.s.initProbe("AnotherProbe")
			createdProbes := mockSummaryState.Probes
			v, found := createdProbes["testProbe"]
			if !found {
				t.Errorf("Summary State doesn't contain probe name: %v", createdProbes)
				t.Logf("probe name found: %v", v)
			}

		})
	}
}
