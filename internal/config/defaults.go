package config

import (
	"os"
	"path/filepath"
)

// getEnvOrDefaults will set value from os.Getenv and default to the specified value
func getEnvOrDefaults(e *ConfigVars) {
	e.GetKubeConfigPath(getDefaultKubeConfigPath()) // KUBE_CONFIG
	e.GetKubeContext()                              // KUBE_CONTEXT
	e.GetOutputType("IO")                           // OUTPUT_TYPE
	e.GetOutputDir("./testoutput")                  // OUTPUT_DIR
	e.GetAuditEnabled("true")                       // AUDIT_ENABLED
	e.GetProbrTags()                                // PROBR_TAGS

	e.GetImageRepository("docker.io") // IMAGE_REPOSITORY
	e.GetCurlImage("curl")            // CURL_IMAGE
	e.GetBusyBoxImage("busybox")      // BUSYBOX_IMAGE

	e.GetAzureSubscriptionID()                      // AZURE_SUBSCRIPTION_ID
	e.GetAzureClientID()                            // AZURE_CLIENT_ID
	e.GetAzureClientSecret()                        // AZURE_CLIENT_SECRET
	e.GetAzureTenantID()                            // AZURE_TENANT_ID
	e.GetAzureLocationDefault()                     // AZURE_LOCATION_DEFAULT
	e.GetDefaultNamespaceAI("probr-defaultns-ai")   // DEFAULT_NS_AZURE_IDENTITY
	e.GetDefaultNamespaceAIB("probr-defaultns-aib") // DEFAULT_NS_AZURE_IDENTITY_BINDING

	e.GetSystemClusterRoles([]string{"system:", "aks", "cluster-admin", "policy-agent"})
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
