package kubernetes

import (
	"reflect"
	"testing"

	"github.com/citihub/probr/config"
	apiv1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestGetDropCapabilitiesFromConfig(t *testing.T) {

	// Default config
	dropCapabilitiesFromDefaultValue := []apiv1.Capability{"NET_RAW"} //This test will fail if default value changes for config.Vars.ServicePacks.Kubernetes.ContainerDropCapabilities. Adjust as needed.

	// Custom config - single value
	customConfig := []string{"CAP_SETUID"}
	dropCapabilitiesFromCustomConfig := []apiv1.Capability{"CAP_SETUID"}

	// Custom config - list
	customConfigList := []string{"NET_RAW", "CAP_SETUID", "CAP_SYS_ADMIN"}
	dropCapabilitiesFromCustomConfigList := []apiv1.Capability{"NET_RAW", "CAP_SETUID", "CAP_SYS_ADMIN"}

	tests := []struct {
		testName           string
		customConfigValues []string
		expectedResult     []apiv1.Capability
	}{
		{
			testName:           "ShouldReturn_UsingDefaultValue",
			customConfigValues: nil, //Use default
			expectedResult:     dropCapabilitiesFromDefaultValue,
		},
		{
			testName:           "ShouldReturn_UsingCustomConfig_SingleValue",
			customConfigValues: customConfig,
			expectedResult:     dropCapabilitiesFromCustomConfig,
		},
		{
			testName:           "ShouldReturn_UsingCustomConfig_List",
			customConfigValues: customConfigList,
			expectedResult:     dropCapabilitiesFromCustomConfigList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.customConfigValues == nil {
				// Create default config
				err := config.Init("")
				if err != nil {
					t.Errorf("[ERROR] error returned from config.Init: %v", err)
				}
			} else {
				// Set custom config vars
				// Only ContainerDropCapabilities for now. Extend if more config vars are used.
				config.Vars.ServicePacks.Kubernetes.ContainerDropCapabilities = tt.customConfigValues
			}

			if got := GetContainerDropCapabilitiesFromConfig(); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("defaultContainerSecurityContext() = %v, Expected result: %v", got, tt.expectedResult)
			}
		})
	}
}
