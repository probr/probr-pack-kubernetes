package coreengine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"

	"github.com/citihub/probr/internal/config"
)

func TestGetRootDir(t *testing.T) {
	// Make sure it doesn't catch one of the several fail conditions
	_, err := getRootDir()
	if err != nil {
		t.Fail()
	}
}

func TestGetOutputPath(t *testing.T) {
	var file *os.File
	d := "test_output_dir"
	f := "test_file"
	desired_file := filepath.Join(d, f) + ".json"
	defer func() {
		// Cleanup test assets
		file.Close()
		err := os.RemoveAll(d)
		if err != nil {
			t.Logf("%s", err)
		}

		// Swallow any panics and print a verbose error message
		if err := recover(); err != nil {
			t.Logf("Panicked when trying to create directory or file: '%s'", desired_file)
			t.Fail()
		}
	}()
	config.Vars.CucumberDir = d
	file, _ = getOutputPath(f)
	if desired_file != file.Name() {
		t.Logf("Desired filepath '%s' does not match '%s'", desired_file, file.Name())
		t.Fail()
	}
}

func TestScenarioString(t *testing.T) {
	gs := &godog.Scenario{Name: "test scenario"}

	// Start scenario
	s := scenarioString(true, gs)
	s_contains_string := strings.Contains(s, "Start")
	if !s_contains_string {
		t.Logf("Test string does not contain 'Start'")
		t.Fail()
	}

	// End scenario
	s = scenarioString(false, gs)
	s_contains_string = strings.Contains(s, "End")
	if !s_contains_string {
		t.Logf("Test string does not contain 'End'")
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
