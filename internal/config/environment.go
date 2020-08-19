package config

import (
	"os"
	"strings"
	"sync"

	"citihub.com/probr/internal/utils"
)

//EnvVariable ...
type EnvVariable int

//EnvVariable ...
const (
	KubeConfig EnvVariable = iota
	AzSubscriptionID
	AzClientID
	AzClientSecret
	AzTenantID
	AzLocationDefault
	ImageRepository
	CurlImage
	BusyBoxImage
)

func (v EnvVariable) String() string {
	return [...]string{"KUBE_CONFIG",
		"AZURE_SUBSCRIPTION_ID",
		"AZURE_CLIENT_ID",
		"AZURE_CLIENT_SECRET",
		"AZURE_TENANT_ID",
		"AZURE_LOCATION_DEFAULT",
		"IMAGE_REPOSITORY",
		"CURL_IMAGE",
		"BUSYBOX_IMAGE"}[v]
}

// Config ...
type Config interface {
	GetKubeConfigPath() *string
	SetKubeConfigPath(*string)

	GetAzureSubscriptionID() *string
	GetAzureClientID() *string
	GetAzureClientSecret() *string
	GetAzureTenantID() *string
	GetAzureLocationDefault() *string

	GetImageRepository() *string
	GetCurlImage() *string
	GetBusyBoxImage() *string
}

// EnvConfig ...
type EnvConfig struct {
	kubeConfigPath *string

	azureSubscriptionID  *string
	azureClientID        *string
	azureClientSecret    *string
	azureTenantID        *string
	azureLocationDefault *string

	imageRepository *string
	curlImage       *string
	busyBoxImage    *string
}

var instance *EnvConfig
var once sync.Once

// GetEnvConfigInstance ...
func GetEnvConfigInstance() *EnvConfig {
	once.Do(func() {
		instance = &EnvConfig{}
	})

	return instance
}

// GetKubeConfigPath ...
func (e *EnvConfig) GetKubeConfigPath() *string {
	if e.kubeConfigPath == nil {
		e.kubeConfigPath = utils.StringPtr(os.Getenv("KUBE_CONFIG"))
	}

	return e.kubeConfigPath
}

// SetKubeConfigPath ...
func (e *EnvConfig) SetKubeConfigPath(p *string) {
	e.kubeConfigPath = p
}

// GetAzureSubscriptionID ...
func (e *EnvConfig) GetAzureSubscriptionID() *string {
	if e.azureSubscriptionID == nil {
		e.azureSubscriptionID = utils.StringPtr(os.Getenv("AZURE_SUBSCRIPTION_ID"))
	}
	return e.azureSubscriptionID
}

// GetAzureClientID ...
func (e *EnvConfig) GetAzureClientID() *string {
	if e.azureClientID == nil {
		e.azureClientID = utils.StringPtr(os.Getenv("AZURE_CLIENT_ID"))
	}
	return e.azureClientID
}

// GetAzureClientSecret ...
func (e *EnvConfig) GetAzureClientSecret() *string {
	if e.azureClientSecret == nil {
		e.azureClientSecret = utils.StringPtr(os.Getenv("AZURE_CLIENT_SECRET"))
	}

	return e.azureClientSecret
}

// GetAzureTenantID ...
func (e *EnvConfig) GetAzureTenantID() *string {
	if e.azureTenantID == nil {
		e.azureTenantID = utils.StringPtr(os.Getenv("AZURE_TENANT_ID"))
	}
	return e.azureTenantID
}

// GetAzureLocationDefault ...
func (e *EnvConfig) GetAzureLocationDefault() *string {
	if e.azureLocationDefault == nil {
		e.azureLocationDefault = utils.StringPtr(os.Getenv("AZURE_LOCATION_DEFAULT"))
	}
	return e.azureLocationDefault
}

// GetImageRepository ...
func (e *EnvConfig) GetImageRepository() *string {
	if e.imageRepository == nil {
		e.imageRepository = utils.StringPtr(os.Getenv("IMAGE_REPOSITORY"))
	}
	return e.imageRepository
}

// GetCurlImage ...
func (e *EnvConfig) GetCurlImage() *string {
	if e.curlImage == nil {
		e.curlImage = utils.StringPtr(os.Getenv("CURL_IMAGE"))
	}
	return e.curlImage
}

// GetBusyBoxImage ...
func (e *EnvConfig) GetBusyBoxImage() *string {
	if e.busyBoxImage == nil {
		e.busyBoxImage = utils.StringPtr(os.Getenv("BUSYBOX_IMAGE"))
	}
	return e.busyBoxImage
}

func (e *EnvConfig) String() string {
	var b strings.Builder
	b.WriteString("Environment:\n")
	b.WriteString(*e.getString(KubeConfig, e.GetKubeConfigPath) + "\n")
	b.WriteString(*e.getString(AzSubscriptionID, e.GetAzureSubscriptionID) + "\n")
	b.WriteString(*e.getString(AzTenantID, e.GetAzureTenantID) + "\n")
	b.WriteString(*e.getString(AzClientID, e.GetAzureClientID) + "\n")
	b.WriteString(*e.getString(AzClientSecret, e.GetAzureClientSecret) + "\n")
	b.WriteString(*e.getString(AzLocationDefault, e.GetAzureLocationDefault) + "\n")
	b.WriteString(*e.getString(ImageRepository, e.GetImageRepository) + "\n")
	b.WriteString(*e.getString(BusyBoxImage, e.GetBusyBoxImage) + "\n")
	b.WriteString(*e.getString(CurlImage, e.GetCurlImage) + "\n")

	return b.String()
}

func (e *EnvConfig) getString(v EnvVariable, f func() *string) *string {
	s := f()

	if s != nil && len(*s) > 0 {
		//env variable is set, return it ...
		return utils.StringPtr(v.String() + ": " + *s)
	}
	//else it's unset ..
	return utils.StringPtr(v.String() + ": unset, default will be used.")
}
