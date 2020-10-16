package config

import (
	"os"
	"path/filepath"
	"strings"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setFromEnvOrDefaults(e *ConfigVars) {
	e.setKubeConfigPath(getDefaultKubeConfigPath()) // KUBE_CONFIG
	e.setKubeContext()                              // KUBE_CONTEXT
	e.setOutputType("IO")                           // PROBR_OUTPUT_TYPE
	e.setOutputDir("cucumber_output")               // PROBR_OUTPUT_DIR
	e.setSummaryEnabled("true")                     // PROBR_SUMMARY_ENABLED
	e.setAuditEnabled("true")                       // PROBR_AUDIT_ENABLED
	e.setProbrTags()                                // PROBR_TAGS
	e.setOverwriteHistoricalAudits("false")         // OVERWRITE_AUDITS

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

func getDefaultOutputDir() string {
	return filepath.Join(homeDir(), "cucumber_output")
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
		e.KubeConfigPath = d
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
		e.SummaryEnabled = d
	}
}

func (e *ConfigVars) setAuditEnabled(d string) {
	if e.AuditEnabled == "" {
		e.AuditEnabled = os.Getenv("PROBR_AUDIT_ENABLED")
	}
	if e.AuditEnabled == "" {
		e.AuditEnabled = d
	}
}

func (e *ConfigVars) setProbrTags() {
	if e.Tags == "" {
		e.Tags = os.Getenv("PROBR_TAGS")
	}
}

// setOutputType ...
func (e *ConfigVars) setOutputType(default_value string) {
	if e.OutputType == "" {
		e.OutputType = os.Getenv("PROBR_OUTPUT_TYPE")
	}
	if e.OutputType == "" {
		e.OutputType = default_value
	}
}

// setOutputDir ...
func (e *ConfigVars) setOutputDir(default_value string) {
	if e.OutputDir == "" {
		e.OutputDir = os.Getenv("PROBR_OUTPUT_DIR")
	}
	if e.OutputDir == "" {
		e.OutputDir = default_value
	}
}

// setOutputDir ...
func (e *ConfigVars) setOverwriteHistoricalAudits(default_value string) {
	if e.OverwriteHistoricalAudits == "" {
		e.OverwriteHistoricalAudits = os.Getenv("OVERWRITE_AUDITS")
	}
	if e.OverwriteHistoricalAudits == "" {
		e.OverwriteHistoricalAudits = default_value
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
func (e *ConfigVars) setDefaultNamespaceAI(default_value string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = os.Getenv("DEFAULT_NS_AZURE_IDENTITY")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = default_value
	}
}

// setDefaultNamespaceAIB ...
func (e *ConfigVars) setDefaultNamespaceAIB(default_value string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = os.Getenv("DEFAULT_NS_AZURE_IDENTITY_BINDING")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = default_value
	}

}

// setImageRepository ...
func (e *ConfigVars) setImageRepository(default_value string) {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
	if e.Images.Repository == "" {
		e.Images.Repository = default_value
	}
}

// setCurlImage ...
func (e *ConfigVars) setCurlImage(default_value string) {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
	if e.Images.Curl == "" {
		e.Images.Curl = default_value
	}
}

// setBusyBoxImage ...
func (e *ConfigVars) setBusyBoxImage(default_value string) {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = default_value
	}
}

// setSystemClusterRoles ...
func (e *ConfigVars) setSystemClusterRoles(default_value []string) {
	//in this case we always want to take the defaults
	//then append anything from the env
	e.SystemClusterRoles = default_value

	t := os.Getenv("SYSTEM_CLUSTER_ROLES") //comma separated
	if len(t) > 0 {
		e.SystemClusterRoles = append(e.SystemClusterRoles, strings.Split(t, ",")...)
	}
}
