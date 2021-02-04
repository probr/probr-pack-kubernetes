package coreengine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"
	"github.com/cucumber/godog"
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

	// Using -testdata- folder to ensure no test resources are included in build
	// Once we migrate to go v1.15, we should use t.TempDir() to ensure built-in test directory is automatically removed by cleanup when test and subtests complete. See: https://golang.org/pkg/testing/#pkg-subdirectories
	d := filepath.Join("testdata", "test_output_dir")

	f := "test_file"
	desiredFile := filepath.Join(d, f) + ".json"
	defer func() {

		cucumberDirFunc = config.Vars.CucumberDir //Restoring to original function after test

		// Cleanup test assets
		file.Close()
		err := os.RemoveAll(d)
		if err != nil {
			t.Logf("%s", err)
		}

		// Swallow any panics and print a verbose error message
		if err := recover(); err != nil {
			t.Logf("Panicked when trying to create directory or file: '%s'", desiredFile)
			t.Fail()
		}
	}()
	// Faking result for config.CucumberDir(). This is used inside getOutputPath.
	cucumberDirFunc = func() string {
		_ = os.MkdirAll(d, 0755) // Creates if not already existing
		return d
	}
	file, _ = getOutputPath(f)
	if desiredFile != file.Name() {
		t.Logf("Desired filepath '%s' does not match '%s'", desiredFile, file.Name())
		t.Fail()
	}
}

func TestScenarioString(t *testing.T) {
	gs := &godog.Scenario{Name: "test scenario"}

	// Start scenario
	s := scenarioString(true, gs)
	sContainsString := strings.Contains(s, "Start")
	if !sContainsString {
		t.Logf("Test string does not contain 'Start'")
		t.Fail()
	}

	// End scenario
	s = scenarioString(false, gs)
	sContainsString = strings.Contains(s, "End")
	if !sContainsString {
		t.Logf("Test string does not contain 'End'")
		t.Fail()
	}
}

func TestGetFeaturePath(t *testing.T) {
	// Faking result for getTmpFeatureFileFunc() to avoid creating -tmp- folder and feature file.
	getTmpFeatureFileFunc = func(featurePath string) (string, error) {
		tmpFeaturePath := filepath.Join("tmp", featurePath)
		return tmpFeaturePath, nil
	}
	defer func() {
		getTmpFeatureFileFunc = getTmpFeatureFile //Restoring to original function after test
	}()

	type args struct {
		path []string
	}
	tests := []struct {
		testName       string
		testArgs       args
		expectedResult string
	}{
		{
			testName:       "GetFeaturePath_WithTwoSubfoldersAndFeatureName_ShouldReturnFeatureFilePath",
			testArgs:       args{path: []string{"service_packs", "kubernetes", "container_registry_access"}},
			expectedResult: filepath.Join("tmp", "service_packs", "kubernetes", "container_registry_access", "container_registry_access.feature"), // Using filepath.join() instead of literal string in order to run test in Windows (\\) and Linux (/)
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := GetFeaturePath(tt.testArgs.path...); got != tt.expectedResult {
				t.Errorf("GetFeaturePath() = %v, Expected: %v", got, tt.expectedResult)
			}
		})
	}
}

func Test_getTmpFeatureFile(t *testing.T) {
	testTmpDir := filepath.Join("testdata", utils.RandomString(10))

	// Faking original behavior
	tmpDirFunc = func() string {
		return testTmpDir
	}
	defer func() {
		tmpDirFunc = config.Vars.TmpDir //Restoring to original function after test

		// Delete test data after tests
		os.RemoveAll(testTmpDir)
	}()

	type args struct {
		featurePath string
	}
	tests := []struct {
		testName       string
		testArgs       args
		expectedResult string
		expectedErr    bool
	}{
		{
			testName:       "ShouldCreateTmpFolderWithFeatureFile",
			testArgs:       args{featurePath: filepath.Join("service_packs", "kubernetes", "container_registry_access", "container_registry_access.feature")},
			expectedResult: filepath.Join(testTmpDir, "service_packs", "kubernetes", "container_registry_access", "container_registry_access.feature"),
			expectedErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := getTmpFeatureFile(tt.testArgs.featurePath)
			if (err != nil) != tt.expectedErr {
				t.Errorf("getTmpFeatureFile() error = %v, expected error: %v", err, tt.expectedErr)
				return
			}
			if got != tt.expectedResult {
				t.Errorf("getTmpFeatureFile() = %v, expected %v", got, tt.expectedResult)
			}
			// Check if file was saved to tmp location
			_, e := os.Stat(tt.expectedResult)
			if e != nil {
				t.Errorf("File not found in tmp location: %v - Error: %v", tt.expectedResult, e)
			}
		})
	}
}

func Test_unpackFileAndSave(t *testing.T) {
	testTmpDir := filepath.Join("testdata", utils.RandomString(10))
	defer func() {
		// Delete test data after tests
		os.RemoveAll(testTmpDir)
	}()

	type args struct {
		origFilePath string
		newFilePath  string
	}
	tests := []struct {
		testName    string
		testArgs    args
		expectedErr bool
	}{
		{
			testName: "ShouldCreateFileInNewLocation",
			testArgs: args{
				origFilePath: filepath.Join("service_packs", "kubernetes", "container_registry_access", "container_registry_access.feature"),
				newFilePath:  filepath.Join(testTmpDir, "service_packs", "kubernetes", "container_registry_access", "container_registry_access.feature"),
			},
			expectedErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if err := unpackFileAndSave(tt.testArgs.origFilePath, tt.testArgs.newFilePath); (err != nil) != tt.expectedErr {
				t.Errorf("unpackFileAndSave() error = %v, expected error: %v", err, tt.expectedErr)
			}
			// Check if file was saved to tmp location
			_, e := os.Stat(tt.testArgs.newFilePath)
			if e != nil {
				t.Errorf("File not found in tmp location: %v - Error: %v", tt.testArgs.newFilePath, e)
			}
		})
	}
}
