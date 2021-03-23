package coreengine

import (
	"bytes"
	"io"
	"os"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"github.com/citihub/probr/config"
)

// GodogProbeHandler is a general implementation of ProbeHandlerFunc.  Based on the
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

	//FUDGE! If the tests are skipped due to tags, then an empty file may
	//be left lingering.  This will have a non-zero size as we've actually
	//had to create the file prior to the test run (see line 31).  If it's
	//less than 4 bytes, it's fairly certain that this will indeed be empty
	//and can be removed.
	i, err := o.Stat()
	s := i.Size()
	o.Close()
	if s < 4 {
		err = os.Remove(o.Name())
		if err != nil {
			//log.Printf("[WARN] unable to remove empty test result file: %v", err)
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
	tags := config.Vars.GetTags()
	opts := godog.Options{
		Format: config.Vars.ResultsFormat,
		Output: colors.Colored(o),
		Paths:  []string{gd.FeaturePath},
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
