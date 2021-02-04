package cli_flags

import (
	"os"
	"testing"

	"github.com/citihub/probr/config"
)

func TestHandleFlags(t *testing.T) {

	// Note:
	// This test is only verifying WriteDirectory cli flag.
	// It could be extended to test other flags, or all. If so, it should be refactored to avoid code duplication. Keeping it simple for now (YAGNI).

	tests := []struct {
		testName                  string
		addCliFlag                string
		expectedResultInConfigVar string
	}{
		{
			testName:                  "HandleFlag_WithCliFlag_ShouldAddCliFlagValueToGlobalConfig",
			addCliFlag:                "-writedirectory=newdirectoryfromcliflag",
			expectedResultInConfigVar: "newdirectoryfromcliflag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			// Simulating cli arguments
			os.Args = append(os.Args, "-writedirectory=newdirectoryfromcliflag")

			// This function is expected to modify global configVar object, setting values based on tags
			// It cannot be called more than once, since global flag object is used and it would raise a "flag redefined" error.
			// An alternative for a potential refactoring is to use FlagSet instead. See: https://stackoverflow.com/questions/24504024/defining-independent-flagsets-in-golang
			HandleFlags()

			//Check WriteDirectory was set in global ConfigVars
			if config.Vars.WriteDirectory != tt.expectedResultInConfigVar {
				t.Errorf("HandleFlags(); config.Vars.WiteDirectory = %v, Expected: %v", config.Vars.WriteDirectory, tt.expectedResultInConfigVar)
				return
			}
		})
	}
}
