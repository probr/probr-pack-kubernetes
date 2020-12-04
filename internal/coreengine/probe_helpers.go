package coreengine

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/internal/config"
)

const rootDirName = "probr"

var outputDir *string

// getRootDir gets the root directory of the probr executable.
func getRootDir() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	log.Printf("[DEBUG] getRootDir pwd is: %v", pwd)

	b := strings.Contains(pwd, rootDirName)
	if !b {
		return "", fmt.Errorf("could not find '%v' root directory in %v", rootDirName, pwd)
	}

	s := strings.SplitAfter(pwd, rootDirName)
	log.Printf("[DEBUG] path(s) after splitting: %v\n", s)

	if len(s) < 1 {
		//expect at least one result
		return "", fmt.Errorf("could not split out '%v' from directory in %v", rootDirName, pwd)
	}

	if !strings.HasSuffix(s[0], rootDirName) {
		//the first path should end with "probr"
		return "", fmt.Errorf("first path after split (%v) does not end with '%v'", s[0], rootDirName)
	}

	return s[0], nil
}

// getOutputPath gets the output path for the test based on the output directory
// plus the test name supplied
func getOutputPath(t string) (*os.File, error) {

	_ = os.Mkdir(config.Vars.CucumberDir, 0755)

	//filename is test name (supplied) + .json
	fn := t + ".json"
	return os.Create(filepath.Join(config.Vars.CucumberDir, fn))
}

// LogScenarioStart logs the name and tags associated with the supplied scenario.
func LogScenarioStart(s *godog.Scenario) {
	log.Print(scenarioString(true, s))
}

// LogScenarioEnd logs the name and tags associated with the supplied scenario.
func LogScenarioEnd(s *godog.Scenario) {
	log.Print(scenarioString(false, s))
}

func scenarioString(st bool, s *godog.Scenario) string {
	var b strings.Builder
	if st {
		b.WriteString("[INFO] >>> Scenario Start: ")
	} else {
		b.WriteString("[INFO] <<< Scenario End: ")
	}

	b.WriteString(s.Name)
	b.WriteString(". (Tags: ")

	for _, t := range s.Tags {
		b.WriteString(t.GetName())
		b.WriteString(" ")
	}
	b.WriteString(").")
	return b.String()
}
