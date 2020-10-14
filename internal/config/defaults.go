package config

import (
	"os"
	"path/filepath"
	"strings"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setEnvOrDefaults(e *ConfigVars) {
	e.setKubeConfigPath(getDefaultKubeConfigPath()) // KUBE_CONFIG
	e.setKubeContext()                              // KUBE_CONTEXT
	e.setOutputType("IO")                           // OUTPUT_TYPE
	e.setOutputDir("testoutput")                    // OUTPUT_DIR
	e.setSummaryEnabled("true")                     // SUMMARY_ENABLED

	e.setImageRepository("docker.io") // IMAGE_REPOSITORY
	e.setCurlImage("curl")            // CURL_IMAGE
	e.setBusyBoxImage("busybox")      // BUSYBOX_IMAGE

	e.setAzureSubscriptionID()                      // AZURE_SUBSCRIPTION_ID
	e.setAzureClientID()                            // AZURE_CLIENT_ID
	e.setAzureClientSecret()                        // AZURE_CLIENT_SECRET
	e.setAzureTenantID()                            // AZURE_TENANT_ID
	e.setAzureLocationDefault()                     // AZURE_LOCATION_DEFAULT
	e.setDefaultNamespaceAI("probr-defaultns-ai")   // DEFAULT_NS_AZURE_IDENTITY
	e.setDefaultNamespaceAIB("probr-defaultns-aib") // DEFAULT_NS_AZURE_IDENTITY_BINDING

	e.setSystemClusterRoles([]string{"system:", "aks", "cluster-admin", "policy-agent"})
}

func getDefaultKubeConfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

// setKubeConfigPath ...
func (e *ConfigVars) setKubeConfigPath(d string) {
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = os.Getenv("KUBE_CONFIG")
	}
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = d // default is specified in caller: config/defaults.go
	}
}

// setKubeContext ...
func (e *ConfigVars) setKubeContext() {
	if e.KubeContext == "" {
		e.KubeContext = os.Getenv("KUBE_CONTEXT")
	}
}

// setSummaryEnabled
func (e *ConfigVars) setSummaryEnabled(d string) {
	if e.SummaryEnabled == "" {
		e.SummaryEnabled = os.Getenv("PROBR_SUMMARY_ENABLED")
	}
	if e.SummaryEnabled == "" {
		e.SummaryEnabled = d // default is specified in caller: config/defaults.go
	}
}

// setProbrTags ...
func (e *ConfigVars) setProbrTags() {
	if e.Tags == "" {
		e.Tags = os.Getenv("PROBR_TAGS")
	}
}

// setOutputType ...
func (e *ConfigVars) setOutputType(s string) {
	if e.OutputType == "" {
		e.OutputType = os.Getenv("OUTPUT_TYPE")
	}
	if e.OutputType == "" {
		e.OutputType = s // default is specified in caller: config/defaults.go
	}
}

// setOutputDir ...
func (e *ConfigVars) setOutputDir(s string) {
	if e.OutputDir == "" {
		e.OutputDir = os.Getenv("OUTPUT_DIR")
	}
	if e.OutputType == "" {
		e.OutputDir = s // default is specified in caller: config/defaults.go
	}
}

// setAzureSubscriptionID ...
func (e *ConfigVars) setAzureSubscriptionID() {
	if e.Azure.SubscriptionID == "" {
		e.Azure.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
}

// setAzureClientID ...
func (e *ConfigVars) setAzureClientID() {
	if e.Azure.ClientID == "" {
		e.Azure.ClientID = os.Getenv("AZURE_CLIENT_ID")
	}
}

// setAzureClientSecret ...
func (e *ConfigVars) setAzureClientSecret() {
	if e.Azure.ClientSecret == "" {
		e.Azure.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}

}

// setAzureTenantID ...
func (e *ConfigVars) setAzureTenantID() {
	if e.Azure.TenantID == "" {
		e.Azure.TenantID = os.Getenv("AZURE_TENANT_ID")
	}
}

// setAzureLocationDefault ...
func (e *ConfigVars) setAzureLocationDefault() {
	if e.Azure.LocationDefault == "" {
		e.Azure.LocationDefault = os.Getenv("AZURE_LOCATION_DEFAULT")
	}
}

// setDefaultNamespaceAI ...
func (e *ConfigVars) setDefaultNamespaceAI(s string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = os.Getenv("DEFAULT_NS_AZURE_IDENTITY")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = s // default is specified in caller: config/defaults.go
	}
}

// setDefaultNamespaceAIB ...
func (e *ConfigVars) setDefaultNamespaceAIB(s string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = os.Getenv("DEFAULT_NS_AZURE_IDENTITY_BINDING")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = s // default is specified in caller: config/defaults.go
	}

}

// setImageRepository ...
func (e *ConfigVars) setImageRepository(s string) {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
	if e.Images.Repository == "" {
		e.Images.Repository = s // default is specified in caller: config/defaults.go
	}
}

// setCurlImage ...
func (e *ConfigVars) setCurlImage(s string) {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
	if e.Images.Curl == "" {
		e.Images.Curl = s // default is specified in caller: config/defaults.go
	}
}

// setBusyBoxImage ...
func (e *ConfigVars) setBusyBoxImage(s string) {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = s // default is specified in caller: config/defaults.go
	}
}

// setSystemClusterRoles ...
func (e *ConfigVars) setSystemClusterRoles(s []string) {
	//in this case we always want to take the defaults
	//then append anything from the env
	e.SystemClusterRoles = s // default is specified in caller: config/defaults.go

	t := os.Getenv("SYSTEM_CLUSTER_ROLES") //comma separated
	if len(t) > 0 {
		e.SystemClusterRoles = append(e.SystemClusterRoles, strings.Split(t, ",")...)
	}
}
