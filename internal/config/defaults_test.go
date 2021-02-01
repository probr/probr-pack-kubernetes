package config

import (
	"os"
	"testing"
)

func Test_setFromEnvOrDefaults(t *testing.T) {

	// Note:
	// This test is only verifying WriteDirectory.
	// It could be extended to test other config vars, or all. If so, it should be refactored to avoid code duplication. Keeping it simple for now (YAGNI).

	envVarCurrentValuePROBR_WRITE_DIRECTORY := os.Getenv("PROBR_WRITE_DIRECTORY") // Used to restore env to original state after test
	defer func() {
		os.Setenv("PROBR_WRITE_DIRECTORY", envVarCurrentValuePROBR_WRITE_DIRECTORY)
	}()
	defaultValuePROBR_WRITE_DIRECTORY := "probr_output"
	envVarValuePROBR_WRITE_DIRECTORY := "ValueFromEnvVar_WriteDirectory"

	type args struct {
		e *ConfigVars
	}
	tests := []struct {
		testName                     string
		testArgs                     args
		setEnvVar                    bool
		expectedResultWriteDirectory string
	}{
		{
			testName:                     "setFromEnvOrDefaults_GivenEnvVar_ShouldSetConfigVarToEnvVarValue",
			testArgs:                     args{e: &ConfigVars{}},
			setEnvVar:                    true,
			expectedResultWriteDirectory: envVarValuePROBR_WRITE_DIRECTORY,
		},
		{
			testName:                     "setFromEnvOrDefaults_WithoutEnvVar_ShouldSetConfigVarToDefaultValue",
			testArgs:                     args{e: &ConfigVars{}},
			setEnvVar:                    false,
			expectedResultWriteDirectory: defaultValuePROBR_WRITE_DIRECTORY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			//Based on test case, modify env vars
			if tt.setEnvVar {
				os.Setenv("PROBR_WRITE_DIRECTORY", envVarValuePROBR_WRITE_DIRECTORY)
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
