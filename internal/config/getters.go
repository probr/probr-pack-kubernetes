package config

import (
	"os"
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

// GetImageRepository ...
func (e *ConfigVars) GetImageRepository() {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
}

// GetCurlImage ...
func (e *ConfigVars) GetCurlImage() {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
}

// GetBusyBoxImage ...
func (e *ConfigVars) GetBusyBoxImage() {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
}
