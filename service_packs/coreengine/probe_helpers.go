package coreengine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/utils"
)

// Probe is an interface used by probes that are to be exported from any service pack
type Probe interface {
	ProbeInitialize(*godog.TestSuiteContext)
	ScenarioInitialize(*godog.ScenarioContext)
	Name() string
	Path() string
}

const rootDirName = "probr"

var outputDir *string

// These variables points to the functions. they are used in oder to be able to mock oiginal behavior during testing.
var cucumberDirFunc = config.Vars.CucumberDir // see TestGetOutputPath
var getTmpFeatureFileFunc = getTmpFeatureFile // See TestGeatFeaturePath
var tmpDirFunc = config.Vars.TmpDir           // See Test_getTmpFeatureFile

// getRootDir gets the root directory of the probr executable.
func getRootDir() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	//log.Printf("[DEBUG] getRootDir pwd is: %v", pwd)

	b := strings.Contains(pwd, rootDirName)
	if !b {
		return "", fmt.Errorf("could not find '%v' root directory in %v", rootDirName, pwd)
	}

	s := strings.SplitAfter(pwd, rootDirName)
	//log.Printf("[DEBUG] path(s) after splitting: %v\n", s)

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

	////filename is test name (supplied) + .json
	fn := t + ".json"
	return os.Create(filepath.Join(cucumberDirFunc(), fn))
}

// GetFeaturePath parses a list of strings into a standardized file path
func GetFeaturePath(path ...string) string {
	featureName := path[len(path)-1] + ".feature"
	dirPath := ""
	for _, folder := range path {
		dirPath = filepath.Join(dirPath, folder)
	}
	//return filepath.Join(dirPath, featureName)
	featurePath := filepath.Join(dirPath, featureName) // This is the original path to feature file in source code

	// Unpacking/copying feature file to tmp location
	tmpFeaturePath, err := getTmpFeatureFileFunc(featurePath)
	if err != nil {
		//log.Printf("Error unpacking feature file '%v' - Error: %v", featurePath, err)
		return ""
	}
	return tmpFeaturePath
}

// getTmpFeatureFile checks if feature file exists in -tmp- folder.
// If so returns the file path, otherwise unpacks the original file using pkger and copies it to -tmp- location before returning file path.
func getTmpFeatureFile(featurePath string) (string, error) {

	tmpFeaturePath := filepath.Join(tmpDirFunc(), featurePath)

	// If file already exists return it
	_, e := os.Stat(tmpFeaturePath)
	if e == nil {
		return tmpFeaturePath, nil
	}

	// If file doesn't exist, extract it from pkger inmemory buffer
	if os.IsNotExist(e) {

		err := unpackFileAndSave(featurePath, tmpFeaturePath)
		if err != nil {
			return "", fmt.Errorf("Error unpacking file: '%v' - Error: %v", featurePath, err)
		}

		return tmpFeaturePath, err
	}

	return "", fmt.Errorf("Error getting os stat for tmp file: '%v' - Error: %v", tmpFeaturePath, e)
}

func unpackFileAndSave(origFilePath string, newFilePath string) error {

	// TODO: This function could be extracted to a separate object i.e: Bundler interface?

	fileBytes, readFileErr := utils.ReadStaticFile(origFilePath) // Read bytes using pkger memory bundle
	if readFileErr != nil {
		return fmt.Errorf("Error reading file content: '%v' - Error: %v", origFilePath, readFileErr)
	}

	createFilePathErr := os.MkdirAll(filepath.Dir(newFilePath), 0755) // Create directory and sub directories for file
	if createFilePathErr != nil {
		return fmt.Errorf("Error creating path for file: '%v' - Error: %v", newFilePath, createFilePathErr)
	}

	writeFileErr := ioutil.WriteFile(newFilePath, fileBytes, 0755) // Save file to new location
	if writeFileErr != nil {
		return fmt.Errorf("Error saving file: '%v' - Error: %v", newFilePath, writeFileErr)
	}

	return nil // File created
}

// LogScenarioStart logs the name and tags associated with the supplied scenario.
func LogScenarioStart(s *godog.Scenario) {
	//log.Print(scenarioString(true, s))
}

// LogScenarioEnd logs the name and tags associated with the supplied scenario.
func LogScenarioEnd(s *godog.Scenario) {
	//log.Print(scenarioString(false, s))
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
