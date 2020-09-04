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

// GetAzureSubscriptionID ...
func (e *ConfigVars) GetAzureSubscriptionID() string {
	if e.Azure.SubscriptionID == "" {
		e.Azure.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
	return e.Azure.SubscriptionID
}

// GetAzureClientID ...
func (e *ConfigVars) GetAzureClientID() string {
	if e.Azure.ClientID == "" {
		e.Azure.ClientID = os.Getenv("AZURE_CLIENT_ID")
	}
	return e.Azure.ClientID
}

// GetAzureClientSecret ...
func (e *ConfigVars) GetAzureClientSecret() string {
	if e.Azure.ClientSecret == "" {
		e.Azure.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}

	return e.Azure.ClientSecret
}

// GetAzureTenantID ...
func (e *ConfigVars) GetAzureTenantID() string {
	if e.Azure.TenantID == "" {
		e.Azure.TenantID = os.Getenv("AZURE_TENANT_ID")
	}
	return e.Azure.TenantID
}

// GetAzureLocationDefault ...
func (e *ConfigVars) GetAzureLocationDefault() string {
	if e.Azure.LocationDefault == "" {
		e.Azure.LocationDefault = os.Getenv("AZURE_LOCATION_DEFAULT")
	}
	return e.Azure.LocationDefault
}

// GetImageRepository ...
func (e *ConfigVars) GetImageRepository() string {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
	return e.Images.Repository
}

// GetCurlImage ...
func (e *ConfigVars) GetCurlImage() string {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
	return e.Images.Curl
}

// GetBusyBoxImage ...
func (e *ConfigVars) GetBusyBoxImage() string {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
	return e.Images.BusyBox
}

// GetOutputType ...
func (e *ConfigVars) GetOutputType(s string) string {
	if e.OutputType == "" {
		e.OutputType = os.Getenv("OUTPUT_TYPE")
	}
	if e.OutputType == "" {
		e.OutputType = s
	}
	return e.OutputType
}
