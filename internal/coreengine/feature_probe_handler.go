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

// GodogProbeHandler is a general implmentation of ProbeHandlerFunc.  Based on the
// output type, the test will either be executed using an in-memory or file output.  In
// both cases, the handler uses the data supplied in GodogProbe to call the underlying
// GoDog test suite.
func GodogProbeHandler(probe *GodogProbe) (int, *bytes.Buffer, error) {
	if config.Vars.OutputType == "INMEM" {
		return inMemGodogProbeHandler(probe)
	}
	return toFileGodogProbeHandler(probe)
}

func toFileGodogProbeHandler(gd *GodogProbe) (int, *bytes.Buffer, error) {
	o, err := getOutputPath(gd.ProbeDescriptor.Name)
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

func inMemGodogProbeHandler(gd *GodogProbe) (int, *bytes.Buffer, error) {
	var t []byte
	o := bytes.NewBuffer(t)
	status, err := runTestSuite(o, gd)
	return status, o, err
}

func runTestSuite(o io.Writer, gd *GodogProbe) (int, error) {
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
		Name:                 gd.ProbeDescriptor.Name,
		TestSuiteInitializer: gd.ProbeInitializer,
		ScenarioInitializer:  gd.ScenarioInitializer,
		Options:              &opts,
	}.Run()

	return status, nil
}

func getProbesPath(gd *GodogProbe) (string, error) {
	r, err := getRootDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine root directory - not able to perform tests")
	}

	if gd.FeaturePath != nil {
		//if we've been given a feature path, add to root and return:
		return filepath.Join(r, *gd.FeaturePath), nil
	}

	//otherwise derive it from the group and name data:
	g := gd.ProbeDescriptor.Group.String()
	group := strings.ReplaceAll(strings.ToLower(g), " ", "")
	name := strings.ReplaceAll(gd.ProbeDescriptor.Name, "_", "")
	return filepath.Join(r, "probes", group, "probe_definitions", name), nil
}
