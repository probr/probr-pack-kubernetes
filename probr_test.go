package probr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/coreengine"
)

func TestGetAllProbeResults(t *testing.T) {
	testTmpDir := filepath.Join("testdata", utils.RandomString(10))

	// Faking original behavior of config.Vars.tmpDir()
	tmpDirFunc = func() string {
		return testTmpDir
	}
	defer func() {
		tmpDirFunc = config.Vars.TmpDir //Restoring to original function after test

		// Delete test data after tests
		os.RemoveAll(testTmpDir)
	}()

	type args struct {
		ps *coreengine.ProbeStore
	}
	tests := []struct {
		testName       string
		testArgs       args
		expectedResult map[string]string
		expectedErr    bool
	}{
		{
			testName:       "ShouldCleanupTmpDir",
			testArgs:       args{coreengine.NewProbeStore()},
			expectedResult: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			//Create temp dir with test data
			if err := createTestData(testTmpDir); err != nil {
				t.Fatalf("Test file couldn't be created. Error: %v", err)
			}

			if got := GetAllProbeResults(tt.testArgs.ps); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("GetAllProbeResults() = %v, Expected: %v", got, tt.expectedResult)
			}
			//Check if tmp folder was removed
			if _, err := os.Stat(testTmpDir); !os.IsNotExist(err) {
				t.Errorf("Temp folder was not removed: %v - Error: %v", testTmpDir, err)
			}
		})
	}
}

func createTestData(testDir string) error {
	createTestFolderErr := os.MkdirAll(testDir, 0755) // Creates if not already existing
	if createTestFolderErr != nil {
		return createTestFolderErr
	}
	testFileContent := []byte("hello\ngo\n")
	return ioutil.WriteFile(filepath.Join(testDir, "testfile"), testFileContent, 0644)
}
