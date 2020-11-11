package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setFromEnvOrDefaults(e *ConfigVars) {

	e.set(&e.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	e.set(&e.KubeContext, "KUBE_CONTEXT", "")
	e.set(&e.Tags, "PROBR_TAGS", "")
	e.set(&e.AuditEnabled, "PROBR_AUDIT_ENABLED", "true")
	e.set(&e.SummaryEnabled, "PROBR_SUMMARY_ENABLED", "true")
	e.set(&e.OutputType, "PROBR_OUTPUT_TYPE", "IO")
	e.set(&e.CucumberDir, "PROBR_CUCUMBER_DIR", "cucumber_output")
	e.set(&e.AuditDir, "PROBR_AUDIT_DIR", "audit_output")
	e.set(&e.LogLevel, "PROBR_LOG_LEVEL", "ERROR")
	e.set(&e.OverwriteHistoricalAudits, "OVERWRITE_AUDITS", "true")
	e.set(&e.ContainerRegistry, "PROBR_CONTAINER_REGISTRY", "docker.io")
	e.set(&e.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")
	e.set(&e.Azure.SubscriptionID, "AZURE_SUBSCRIPTION_ID", "")
	e.set(&e.Azure.ClientID, "AZURE_CLIENT_ID", "")
	e.set(&e.Azure.ClientSecret, "AZURE_CLIENT_SECRET", "")
	e.set(&e.Azure.TenantID, "AZURE_TENANT_ID", "")
	e.set(&e.Azure.LocationDefault, "AZURE_LOCATION_DEFAULT", "")
	e.set(&e.Azure.Identity.DefaultNamespaceAI, "DEFAULT_NS_AZURE_IDENTITY", "probr-defaultns-ai")
	e.set(&e.Azure.Identity.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-defaultns-aib")

	e.set(&e.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
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

// set fetches the env var or sets the default value as needed for the specified field from ConfigVars
func (e *ConfigVars) set(field interface{}, var_name string, default_value interface{}) {
	switch v := field.(type) {
	default:
		log.Fatalf("unexpected type for %v, %T", var_name, v)
	case *string:
		if *field.(*string) == "" {
			*field.(*string) = os.Getenv(var_name)
		}
		if *field.(*string) == "" {
			*field.(*string) = default_value.(string)
		}
	case *[]string:
		if len(*field.(*[]string)) == 0 {
			t := os.Getenv(var_name) // if []string, env var should be comma separated values
			if len(t) > 0 {
				*field.(*[]string) = append(*field.(*[]string), strings.Split(t, ",")...)
			}
		}
		if len(*field.(*[]string)) == 0 {
			*field.(*[]string) = default_value.([]string)
		}
	}

}
