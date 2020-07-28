package features

import (
	"fmt"
	"strings"
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"citihub.com/probr/internal/coreengine"
)

//this is the "TEST HANDLER" impl  and will get called when probr is invoked from the CLI or API
//all we do here is set the godog args based on what has been supplied (e.g. output path)
//and call to the "feature" implementation (i.e the same impl when godog / go test is invoked)

//GodogTestHandler ...
func GodogTestHandler(gd *coreengine.GodogTest) (int, error) {
	r, err := GetRootDir()

	if err != nil {
		return -1, fmt.Errorf("unable to determine root directory - not able to perform tests")
	}

	var g = gd.TestDescriptor.Group.String()
	var c = gd.TestDescriptor.Category.String()
	featPath := filepath.Join(r, "test", "features", 
		strings.TrimSpace(strings.ToLower(g)), strings.TrimSpace(strings.ToLower(c)), "features")

	f, err := GetOutputPath(&c)
	if err != nil {
		return -2, err
	}

	opts := godog.Options{
		Format: "cucumber",
		Output: colors.Colored(f),
		Paths:  []string{featPath},
	}

	status := godog.TestSuite{
		Name:                 gd.TestDescriptor.Name,
		TestSuiteInitializer: gd.TestSuiteInitializer,
		ScenarioInitializer:  gd.ScenarioInitializer,
		Options:              &opts,
	}.Run()

	return status, nil
}
