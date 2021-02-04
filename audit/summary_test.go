package audit

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/citihub/probr/config"
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

// createMockProbe - creates a mock SummaryState and probe object in it and returns SummaryState object.
func createSummaryStateWithMockProbe(probename string) SummaryState {
	var sumstate SummaryState
	sumstate.Probes = make(map[string]*Probe)
	sumstate.Meta = make(map[string]interface{})
	sumstate.Meta["names of pods created"] = []string{}
	ap := filepath.Join(config.Vars.AuditDir(), (probename + ".json")) // Needed in both Probe and ProbeAudit
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

func TestSummaryState_SetProbrStatus(t *testing.T) {
	var probeName = "testProbe"
	var mockSummaryState = createSummaryStateWithMockProbe(probeName)
	type args struct {
		probeName    string
		probesPassed int
		probesFailed int
	}
	tests := []struct {
		testName       string
		s              *SummaryState
		expectedResult string
		args           args
	}{
		{
			testName:       "SetProbrStatus",
			s:              &mockSummaryState,
			expectedResult: "Complete - All Probes Completed Successfully",
			args:           args{probeName: "testProbe", probesPassed: 1, probesFailed: 0},
		}, {
			testName:       "SetProbrStatus_WithFailedProbes",
			s:              &mockSummaryState,
			expectedResult: fmt.Sprintf("Complete - %v of %v Probes Failed", 1, 2),
			args:           args{probeName: "testProbe", probesPassed: 0, probesFailed: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.s.initProbe(tt.args.probeName)
			tt.s.initProbe("anotherPod")
			tt.s.ProbesPassed = tt.args.probesPassed
			tt.s.ProbesFailed = tt.args.probesFailed
			tt.s.ProbesSkipped = 0
			tt.s.SetProbrStatus()
			if string(tt.s.Status) != string(tt.expectedResult) {
				t.Errorf("\nCall: SetProbrStatus()\nExpected: %q", string(tt.expectedResult))
			}
		})
	}
}

func TestSummaryState_LogProbeMeta(t *testing.T) {

	var mockSummaryState SummaryState
	mockSummaryState.Probes = make(map[string]*Probe)
	mockSummaryState.Meta = make(map[string]interface{})
	mockSummaryState.Meta["names of pods created"] = []string{}

	type args struct {
		name  string
		key   string
		value interface{}
	}
	tests := []struct {
		testName       string
		s              *SummaryState
		expectedResult string
		args           args
	}{
		{
			testName:       "LogProbeMeta",
			s:              &mockSummaryState,
			expectedResult: "valueTest",
			args:           args{name: "testProbe", key: "testKey", value: "valueTest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.s.LogProbeMeta(tt.args.name, tt.args.key, tt.args.value)
			if string(tt.s.Probes[tt.args.name].Meta[tt.args.key].(string)) != string(tt.expectedResult) {
				t.Errorf("\nCall: LogProbeMeta()\nExpected: %q", string(tt.expectedResult))
			}
		})
	}
}

func TestSummaryState_GetProbeLog(t *testing.T) {
	var probeName = "testProbe"
	var mockSummaryState SummaryState
	mockSummaryState.Probes = make(map[string]*Probe)
	mockSummaryState.Meta = make(map[string]interface{})
	mockSummaryState.Meta["names of pods created"] = []string{}
	var expectedProbeWithSummaryState = createSummaryStateWithMockProbe(probeName)

	type args struct {
		name string
	}
	tests := []struct {
		testName string
		s        *SummaryState
		want     *Probe
		args     args
	}{
		{
			testName: "GetProbeLog",
			s:        &mockSummaryState,
			want:     expectedProbeWithSummaryState.Probes[probeName],
			args:     args{name: "testProbe"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := tt.s.GetProbeLog(tt.args.name); strings.Compare(got.name, tt.want.name) > 0 {
				t.Errorf("SummaryState.GetProbeLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummaryState_completeProbe(t *testing.T) {

	var probeName = "testProbe"
	var mockSummaryState = createSummaryStateWithMockProbe(probeName)

	mockSummaryState.Probes[probeName].audit.Scenarios = make(map[int]*ScenarioAudit)
	scenarioCounter := len(mockSummaryState.Probes[probeName].audit.Scenarios) + 1
	mockSummaryState.Probes[probeName].Result = "Excluded"
	mockSummaryState.Probes[probeName].audit.Scenarios[scenarioCounter] = &ScenarioAudit{
		Name:  "scena1",
		Steps: make(map[int]*StepAudit),
		Tags:  []string{"scenario"},
	}
	mockSummaryState.Probes[probeName].ScenariosFailed = 0
	mockSummaryState.Probes[probeName].audit.ScenariosFailed = &mockSummaryState.Probes[probeName].ScenariosFailed

	type args struct {
		e *Probe
	}
	tests := []struct {
		testName string
		s        *SummaryState
		args
		expectedResult string
	}{
		{
			testName:       "completeProbe",
			s:              &mockSummaryState,
			args:           args{e: mockSummaryState.Probes[probeName]},
			expectedResult: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			tt.s.completeProbe(tt.args.e)

			if strings.Compare(tt.args.e.Result, tt.expectedResult) > 0 {
				t.Errorf("SummaryState.completeProbe() = %v, want %v", tt.args.e.Result, tt.expectedResult)
			}

		})
	}
}
