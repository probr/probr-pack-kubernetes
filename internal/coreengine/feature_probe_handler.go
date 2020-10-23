package coreengine

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"github.com/citihub/probr/internal/config"
)

// GodogTestHandler is a general implmentation of TestHandlerFunc.  Based on the
// output type, the test will either be executed using an in-memory or file output.  In
// both cases, the handler uses the data supplied in GodogTest to call the underlying
// GoDog test suite.
func GodogTestHandler(gd *GodogTest) (int, *bytes.Buffer, error) {
	if config.Vars.OutputType == "INMEM" {
		return inMemGodogTestHandler(gd)
	}
	return toFileGodogTestHandler(gd)
}

func toFileGodogTestHandler(gd *GodogTest) (int, *bytes.Buffer, error) {
	o, err := getOutputPath(gd.TestDescriptor.Name)
	if err != nil {
		return -1, nil, err
	}
	status, err := runTestSuite(o, gd)

	//TODO - review!
	//FUDGE! If the tests are skipped due to tags, then an empty file may
	//be left lingering.  This will have a non-zero size as we've actually
	//had to create the file prior to the test run (see line 31).  If it's
	//less than 4 bytes, it's fairly certain that this will indeed be empty
	//and can be removed.
	i, err := o.Stat()
	s := i.Size()

	if s < 4 {
		err = os.Remove(o.Name())
		if err != nil {
			log.Printf("[WARN] error removing empty test result file: %v", err)
		}
	}
	return status, nil, err
}

func inMemGodogTestHandler(gd *GodogTest) (int, *bytes.Buffer, error) {
	var t []byte
	o := bytes.NewBuffer(t)
	status, err := runTestSuite(o, gd)
	return status, o, err
}

func runTestSuite(o io.Writer, gd *GodogTest) (int, error) {
	f, err := getProbesPath(gd)
	if err != nil {
		return -2, err
	}

	tags := config.Vars.GetTags()
	opts := godog.Options{
		Format: "cucumber",
		Output: colors.Colored(o),
		Paths:  []string{f},
		Tags:   tags,
	}

	status := godog.TestSuite{
		Name:                 gd.TestDescriptor.Name,
		TestSuiteInitializer: gd.TestSuiteInitializer,
		ScenarioInitializer:  gd.ScenarioInitializer,
		Options:              &opts,
	}.Run()

	return status, nil
}

func getProbesPath(gd *GodogTest) (string, error) {
	r, err := getRootDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine root directory - not able to perform tests")
	}

	if gd.FeaturePath != nil {
		//if we've been given a feature path, add to root and return:
		return filepath.Join(r, *gd.FeaturePath), nil
	}

	//otherwise derive it from the group and name data:
	g := gd.TestDescriptor.Group.String()
	group := strings.ReplaceAll(strings.ToLower(g), " ", "")
	name := strings.ReplaceAll(gd.TestDescriptor.Name, "_", "")

	return filepath.Join(r, "probes", group, "probe_definitions", name), nil
}
