package config

import (
	"os"
)

// GetKubeConfigPath ...
func (e *Config) GetKubeConfigPath(d string) {
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = os.Getenv("KUBE_CONFIG")
	}
	if e.KubeConfigPath == "" {
		e.KubeConfigPath = d
	}
}

// GetAzureSubscriptionID ...
func (e *Config) GetAzureSubscriptionID() string {
	if e.Azure.SubscriptionID == "" {
		e.Azure.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	}
	return e.Azure.SubscriptionID
}

// GetAzureClientID ...
func (e *Config) GetAzureClientID() string {
	if e.Azure.ClientID == "" {
		e.Azure.ClientID = os.Getenv("AZURE_CLIENT_ID")
	}
	return e.Azure.ClientID
}

// GetAzureClientSecret ...
func (e *Config) GetAzureClientSecret() string {
	if e.Azure.ClientSecret == "" {
		e.Azure.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	}

	return e.Azure.ClientSecret
}

// GetAzureTenantID ...
func (e *Config) GetAzureTenantID() string {
	if e.Azure.TenantID == "" {
		e.Azure.TenantID = os.Getenv("AZURE_TENANT_ID")
	}
	return e.Azure.TenantID
}

// GetAzureLocationDefault ...
func (e *Config) GetAzureLocationDefault() string {
	if e.Azure.LocationDefault == "" {
		e.Azure.LocationDefault = os.Getenv("AZURE_LOCATION_DEFAULT")
	}
	return e.Azure.LocationDefault
}

// GetImageRepository ...
func (e *Config) GetImageRepository() string {
	if e.Images.Repository == "" {
		e.Images.Repository = os.Getenv("IMAGE_REPOSITORY")
	}
	return e.Images.Repository
}

// GetCurlImage ...
func (e *Config) GetCurlImage() string {
	if e.Images.Curl == "" {
		e.Images.Curl = os.Getenv("CURL_IMAGE")
	}
	return e.Images.Curl
}

// GetBusyBoxImage ...
func (e *Config) GetBusyBoxImage() string {
	if e.Images.BusyBox == "" {
		e.Images.BusyBox = os.Getenv("BUSYBOX_IMAGE")
	}
	return e.Images.BusyBox
}

// GetOutputType ...
func (e *Config) GetOutputType() string {
	if e.OutputType == "" {
		e.OutputType = os.Getenv("OUTPUT_TYPE")
	}
	return e.OutputType
}
