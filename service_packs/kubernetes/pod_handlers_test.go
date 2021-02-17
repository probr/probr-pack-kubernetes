package kubernetes

import (
	"reflect"
	"testing"

	"github.com/citihub/probr/config"
	apiv1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func TestGetDropCapabilitiesFromConfig(t *testing.T) {

	// Default value
	dropCapabilitiesFromDefaultValue := []apiv1.Capability{"NET_RAW"} //This test will fail if default value changes for config.Vars.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities. Adjust as needed.

	// Config value
	configValue := []string{"CAP_SETUID"}
	dropCapabilitiesFromConfigValue := []apiv1.Capability{"CAP_SETUID"}

	tests := []struct {
		testName       string
		configValue    []string
		expectedResult []apiv1.Capability
	}{
		{
			testName:       "ShouldReturn_UsingDefaultValue",
			configValue:    nil, //Use default
			expectedResult: dropCapabilitiesFromDefaultValue,
		},
		{
			testName:       "ShouldReturn_UsingConfigValue",
			configValue:    configValue,
			expectedResult: dropCapabilitiesFromConfigValue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.configValue == nil {
				// Create default config
				err := config.Init("")
				if err != nil {
					t.Errorf("[ERROR] error returned from config.Init: %v", err)
				}
			} else {
				// Set config vars
				// Only ContainerDropCapabilities for now. Extend if more config vars are used.
				config.Vars.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities = tt.configValue
			}

			if got := GetContainerDropCapabilitiesFromConfig(); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("defaultContainerSecurityContext() = %v, Expected result: %v", got, tt.expectedResult)
			}
		})
	}
}

func TestGetCapabilitiesFromList(t *testing.T) {

	// Single value
	singleItemList := []string{"CAP_SETUID"}
	capabilitiesSingleValue := []apiv1.Capability{"CAP_SETUID"}

	// Multiple values
	multiItemsList := []string{"NET_RAW", "CAP_SETUID", "CAP_SYS_ADMIN"}
	capabilitiesMultiValues := []apiv1.Capability{"NET_RAW", "CAP_SETUID", "CAP_SYS_ADMIN"}

	type args struct {
		capList []string
	}
	tests := []struct {
		testName       string
		testArgs       args
		expectedResult []apiv1.Capability
	}{
		{
			testName:       "ShouldReturn_ApiCapabilities_SingleValue",
			testArgs:       args{singleItemList},
			expectedResult: capabilitiesSingleValue,
		},
		{
			testName:       "ShouldReturn_ApiCapabilities_MultipleValues",
			testArgs:       args{multiItemsList},
			expectedResult: capabilitiesMultiValues,
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if got := GetCapabilitiesFromList(tt.testArgs.capList); !reflect.DeepEqual(got, tt.expectedResult) {
				t.Errorf("GetCapabilitiesFromList() = %v, Expected: %v", got, tt.expectedResult)
			}
		})
	}
}
