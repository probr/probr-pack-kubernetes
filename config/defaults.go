package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func setFromEnvOrDefaults(e *VarOptions) {

	e.set(&e.Tags, "PROBR_TAGS", "")
	e.set(&e.AuditEnabled, "PROBR_AUDIT_ENABLED", "true")
	e.set(&e.OutputType, "PROBR_OUTPUT_TYPE", "IO")
	e.set(&e.WriteDirectory, "PROBR_WRITE_DIRECTORY", "probr_output")
	e.set(&e.LogLevel, "PROBR_LOG_LEVEL", "ERROR")
	e.set(&e.OverwriteHistoricalAudits, "OVERWRITE_AUDITS", "true")
	e.set(&e.WriteConfig, "PROBR_LOG_CONFIG", "true")
	e.set(&e.ResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")

	e.set(&e.ServicePacks.Kubernetes.KeepPods, "PROBR_KEEP_PODS", "false")
	e.set(&e.ServicePacks.Kubernetes.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	e.set(&e.ServicePacks.Kubernetes.KubeContext, "KUBE_CONTEXT", "")
	e.set(&e.ServicePacks.Kubernetes.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
	e.set(&e.ServicePacks.Kubernetes.AuthorisedContainerRegistry, "PROBR_AUTHORISED_REGISTRY", "")
	e.set(&e.ServicePacks.Kubernetes.UnauthorisedContainerRegistry, "PROBR_UNAUTHORISED_REGISTRY", "")
	e.set(&e.ServicePacks.Kubernetes.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")
	e.set(&e.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities, "PROBR_REQUIRED_DROP_CAPABILITIES", []string{"NET_RAW"})
	e.set(&e.ServicePacks.Kubernetes.ContainerAllowedAddCapabilities, "PROBR_ALLOWED_ADD_CAPABILITIES", []string{""})
	e.set(&e.ServicePacks.Kubernetes.ApprovedVolumeTypes, "PROBR_APPROVED_VOLUME_TYPES", []string{"configmap", "emptydir", "persistentvolumeclaim"})
	e.set(&e.ServicePacks.Kubernetes.UnapprovedHostPort, "PROBR_UNAPPROVED_HOSTPORT", "22")
	e.set(&e.ServicePacks.Kubernetes.SystemNamespace, "PROBR_K8S_SYSTEM_NAMESPACE", "kube-system")
	e.set(&e.ServicePacks.Kubernetes.DashboardPodNamePrefix, "PROBR_K8S_DASHBOARD_PODNAMEPREFIX", "kubernetes-dashboard")
	e.set(&e.ServicePacks.Kubernetes.ProbeNamespace, "PROBR_K8S_PROBE_NAMESPACE", "probr-general-test-ns")
	e.set(&e.ServicePacks.Kubernetes.Azure.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-aib")
	e.set(&e.ServicePacks.Kubernetes.Azure.IdentityNamespace, "PROBR_K8S_AZURE_IDENTITY_NAMESPACE", "kube-system")

	e.set(&e.CloudProviders.Azure.TenantID, "AZURE_TENANT_ID", "")
	e.set(&e.CloudProviders.Azure.SubscriptionID, "AZURE_SUBSCRIPTION_ID", "")
	e.set(&e.CloudProviders.Azure.ClientID, "AZURE_CLIENT_ID", "")
	e.set(&e.CloudProviders.Azure.ClientSecret, "AZURE_CLIENT_SECRET", "")
	e.set(&e.CloudProviders.Azure.ResourceGroup, "AZURE_RESOURCE_GROUP", "")
	e.set(&e.CloudProviders.Azure.ResourceLocation, "AZURE_RESOURCE_LOCATION", "")
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

// set fetches the env var or sets the default value as needed for the specified field from VarOptions
func (e *VarOptions) set(field interface{}, varName string, defaultValue interface{}) {
	switch v := field.(type) {
	default:
		//log.Fatalf("unexpected type for %v, %T", varName, v)
		panic(fmt.Sprintf("unexpected type for %v, %T", varName, v))
	case *string:
		if *field.(*string) == "" {
			*field.(*string) = os.Getenv(varName)
		}
		if *field.(*string) == "" {
			*field.(*string) = defaultValue.(string)
		}
	case *[]string:
		if len(*field.(*[]string)) == 0 {
			t := os.Getenv(varName) // if []string, env var should be comma separated values
			if len(t) > 0 {
				*field.(*[]string) = append(*field.(*[]string), strings.Split(t, ",")...)
			}
		}
		if len(*field.(*[]string)) == 0 {
			*field.(*[]string) = defaultValue.([]string)
		}
	}

}
