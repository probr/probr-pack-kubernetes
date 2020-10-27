package coreengine

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/citihub/probr/internal/config"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

func TestGetRootDir(t *testing.T) {
	// Make sure it doesn't catch one of the several fail conditions
	_, err := getRootDir()
	if err != nil {
		t.Fail()
	}
}

func TestGetProbesPath(t *testing.T) {
	var failed bool
	r, _ := getRootDir()
	desired_path := filepath.Join(r, "probes", "clouddriver", "probe_definitions", "accountmanager")

	// Test with feature path provided
	p := filepath.Join("probes", "clouddriver", "probe_definitions", "accountmanager")
	test := &GodogProbe{FeaturePath: &p}
	path, err := getProbesPath(test)
	if err != nil || desired_path != path {
		t.Logf("Custom feature path not handled properly")
		failed = true
	}

	// Test building path from properties
	test = &GodogProbe{ProbeDescriptor: &ProbeDescriptor{Group: CloudDriver, Name: "account_manager"}}
	path, err = getProbesPath(test)
	if err != nil || desired_path != path {
		t.Logf("Failed to build probe path from GodogProbe properties")
		failed = true
	}

	// Allow both failures to log before ending, if applicable
	if failed {
		t.Fail()
	}
}

func TestLogAndReturnError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf) // Intercept expected Stderr output
	defer func() {
		log.SetOutput(os.Stderr) // Return to normal Stderr handling after function
	}()

	long_string := "Verify that this somewhat long string remains unchanged in the output after being handled"
	err := LogAndReturnError(long_string)
	err_contains_string := strings.Contains(err.Error(), long_string)
	if !err_contains_string {
		t.Logf("Test string was not properly included in retured error")
		t.Fail()
	}
}

func TestScenarioString(t *testing.T) {
	var failed bool
	gs := &godog.Scenario{Name: "test scenario"}

	// Start scenario
	s := scenarioString(true, gs)
	s_contains_string := strings.Contains(s, "Start")
	if !s_contains_string {
		t.Logf("Test string does not contain 'Start'")
		failed = true
	}

	// End scenario
	s = scenarioString(false, gs)
	s_contains_string = strings.Contains(s, "End")
	if !s_contains_string {
		t.Logf("Test string does not contain 'End'")
		failed = true
	}

	// Allow both failures to log before ending, if applicable
	if failed {
		t.Fail()
	}
}

func TestTagsNotExcluded(t *testing.T) {
	tags := []*messages.Pickle_PickleTag{
		&messages.Pickle_PickleTag{Name: "@test-tag", AstNodeId: "123"},
	}

	if !TagsNotExcluded(tags) {
		t.Logf("Non-excluded tag is being reported as excluded")
		t.Fail()
	}

	config.Vars.TagExclusions = []string{"test-tag"}
	if TagsNotExcluded(tags) {
		t.Logf("Excluded tag is being reported as not excluded")
		t.Fail()
	}
}
