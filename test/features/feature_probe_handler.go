package features

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/citihub/probr/internal/config"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"gitlab.com/citihub/probr/internal/coreengine"
)

//this is the "TEST HANDLER" impl  and will get called when probr is invoked from the CLI or API
//all we do here is set the godog args based on what has been supplied (e.g. output path)
//and call to the "feature" implementation (i.e the same impl when godog / go test is invoked)

//GodogTestHandler ...
func GodogTestHandler(gd *coreengine.GodogTest) (int, error) {
	if config.GetOutputType() == "INMEM" {
		return InMemGodogTestHandler(gd)
	}
	return ToFileGodogTestHandler(gd)
}

// GetOutputType

func ToFileGodogTestHandler(gd *coreengine.GodogTest) (int, error) {
	o, err := GetOutputPath(&gd.TestDescriptor.Name)
	if err != nil {
		return -1, err
	}

func InMemGodogTestHandler(gd *coreengine.GodogTest) (int, error) {
	var t []byte
	o := bytes.NewBuffer(t)
	status, err := runTestSuite(o, gd)
	return status, err
}

func InMemGodogTestHandler(gd *coreengine.GodogTest) (int, error) {
	var t []byte
	o := bytes.NewBuffer(t)
	status, err := runTestSuite(o, gd)
	return status, err
}

func runTestSuite(o io.Writer, gd *coreengine.GodogTest) (int, error) {
	f, err := getFeaturesPath(gd)
	if err != nil {
		return -2, err
	}

	opts := godog.Options{
		Format: "cucumber",
		Output: colors.Colored(o),
		Paths:  []string{f},
	}

	status := godog.TestSuite{
		Name:                 gd.TestDescriptor.Name,
		TestSuiteInitializer: gd.TestSuiteInitializer,
		ScenarioInitializer:  gd.ScenarioInitializer,
		Options:              &opts,
	}.Run()

<<<<<<< Updated upstream
<<<<<<< Updated upstream
	return status, nil
}

=======
>>>>>>> Stashed changes
=======
>>>>>>> Stashed changes
func getFeaturesPath(gd *coreengine.GodogTest) (string, error) {
	r, err := GetRootDir()
	if err != nil {
		return "", fmt.Errorf("unable to determine root directory - not able to perform tests")
	}

	if gd.FeaturePath != nil {
		//if we've been given a feature path, add to root and return:
		return filepath.Join(r, *gd.FeaturePath), nil
	}

	//otherwise derive it from the group and category data:
	var g = gd.TestDescriptor.Group.String()
	var c = gd.TestDescriptor.Category.String()

	return filepath.Join(r, "test", "features",
		strings.ReplaceAll(strings.ToLower(g), " ", ""),
		strings.ReplaceAll(strings.ToLower(c), " ", ""), "features"), nil

}
