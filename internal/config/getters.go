package config

import (
	"os"
	"strings"
)

// GetKubeConfigPath ...
func (e *ConfigVars) GetKubeConfigPath(d string) {
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = os.Getenv("KUBE_CONFIG")
	}
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = d // default is specified in caller: config/defaults.go
	}
}

// GetKubeContext ...
func (e *ConfigVars) GetKubeContext() {
	if e.KubeContext == "" {
		e.KubeContext = os.Getenv("KUBE_CONTEXT")
	}
}

// GetProbrTags ...
func (e *ConfigVars) GetProbrTags() {
	if e.Tests.Tags == "" {
		e.Tests.Tags = os.Getenv("PROBR_TAGS")
	}
}

// GetOutputType ...
func (e *ConfigVars) GetOutputType(s string) {
	if e.OutputType == "" {
		e.OutputType = os.Getenv("OUTPUT_TYPE")
	}
	if e.OutputType == "" {
		e.OutputType = s // default is specified in caller: config/defaults.go
	}
}

// GetAzureSubscriptionID ...
func (e *ConfigVars) GetAzureSubscriptionID() {
	if e.Azure.SubscriptionID == "" {
		e.Azure.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
}

// GetAzureClientID ...
func (e *ConfigVars) GetAzureClientID() {
	if e.Azure.ClientID == "" {
		e.Azure.ClientID = os.Getenv("AZURE_CLIENT_ID")
	}
}

// GetAzureClientSecret ...
func (e *ConfigVars) GetAzureClientSecret() {
	if e.Azure.ClientSecret == "" {
		e.Azure.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}

}

// GetAzureTenantID ...
func (e *ConfigVars) GetAzureTenantID() {
	if e.Azure.TenantID == "" {
		e.Azure.TenantID = os.Getenv("AZURE_TENANT_ID")
	}
}

// GetAzureLocationDefault ...
func (e *ConfigVars) GetAzureLocationDefault() {
	if e.Azure.LocationDefault == "" {
		e.Azure.LocationDefault = os.Getenv("AZURE_LOCATION_DEFAULT")
	}
}

// GetDefaultNamespaceAI ...
func (e *ConfigVars) GetDefaultNamespaceAI(s string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = os.Getenv("DEFAULT_NS_AZURE_IDENTITY")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAI == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAI = s // default is specified in caller: config/defaults.go
	}
}

// GetDefaultNamespaceAIB ...
func (e *ConfigVars) GetDefaultNamespaceAIB(s string) {
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = os.Getenv("DEFAULT_NS_AZURE_IDENTITY_BINDING")
	}
	if e.Azure.AzureIdentity.DefaultNamespaceAIB == "" {
		e.Azure.AzureIdentity.DefaultNamespaceAIB = s // default is specified in caller: config/defaults.go
	}

}

// GetImageRepository ...
func (e *ConfigVars) GetImageRepository(s string) {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
	if e.Images.Repository == "" {
		e.Images.Repository = s // default is specified in caller: config/defaults.go
	}
}

// GetCurlImage ...
func (e *ConfigVars) GetCurlImage(s string) {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
	if e.Images.Curl == "" {
		e.Images.Curl = s // default is specified in caller: config/defaults.go
	}
}

// GetBusyBoxImage ...
func (e *ConfigVars) GetBusyBoxImage(s string) {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = s // default is specified in caller: config/defaults.go
	}
}

// GetSystemClusterRoles ...
func (e *ConfigVars) GetSystemClusterRoles(s []string) {
	//in this case we always want to take the defaults
	//then append anything from the env
	e.SystemClusterRoles = s // default is specified in caller: config/defaults.go

	t := os.Getenv("SYSTEM_CLUSTER_ROLES") //comma separated
	if len(t) > 0 {
		e.SystemClusterRoles = append(e.SystemClusterRoles, strings.Split(t, ",")...)
	}
}
