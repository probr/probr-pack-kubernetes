package config

import (
	"os"
	"testing"
)

func Test_setFromEnvOrDefaults(t *testing.T) {

	// Note:
	// This test is only verifying WriteDirectory.
	// It could be extended to test other config vars, or all. If so, it should be refactored to avoid code duplication. Keeping it simple for now (YAGNI).

	envVarCurrentValuePROBRWRITEDIRECTORY := os.Getenv("PROBR_WRITE_DIRECTORY") // Used to restore env to original state after test
	defer func() {
		os.Setenv("PROBR_WRITE_DIRECTORY", envVarCurrentValuePROBRWRITEDIRECTORY)
	}()
	defaultValuePROBRWRITEDIRECTORY := "probr_output"
	envVarValuePROBRWRITEDIRECTORY := "ValueFromEnvVar_WriteDirectory"

	type args struct {
		e *VarOptions
	}
	tests := []struct {
		testName                     string
		testArgs                     args
		setEnvVar                    bool
		expectedResultWriteDirectory string
	}{
		{
			testName:                     "setFromEnvOrDefaults_GivenEnvVar_ShouldSetConfigVarToEnvVarValue",
			testArgs:                     args{e: &VarOptions{}},
			setEnvVar:                    true,
			expectedResultWriteDirectory: envVarValuePROBRWRITEDIRECTORY,
		},
		{
			testName:                     "setFromEnvOrDefaults_WithoutEnvVar_ShouldSetConfigVarToDefaultValue",
			testArgs:                     args{e: &VarOptions{}},
			setEnvVar:                    false,
			expectedResultWriteDirectory: defaultValuePROBRWRITEDIRECTORY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			//Based on test case, modify env vars
			if tt.setEnvVar {
				os.Setenv("PROBR_WRITE_DIRECTORY", envVarValuePROBRWRITEDIRECTORY)
			} else {
				os.Setenv("PROBR_WRITE_DIRECTORY", "")
			}

			setFromEnvOrDefaults(tt.testArgs.e) //This function will modify config object

			//Check WriteDirectory
			if tt.testArgs.e.WriteDirectory != tt.expectedResultWriteDirectory {
				t.Errorf("setFromEnvOrDefaults(); PROBR_WRITE_DIRECTORY = %v, Expected: %v", tt.testArgs.e.WriteDirectory, tt.expectedResultWriteDirectory)
				return
			}

		})
	}
}
