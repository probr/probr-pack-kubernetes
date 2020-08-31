package config

// getEnvOrDefaults will set value from os.Getenv and default to the specified value
func getEnvOrDefaults(e *Config) {
	e.GetKubeConfigPath(".kube/config")
	e.GetOutputType()

	e.GetImageRepository()
	e.GetCurlImage()
	e.GetBusyBoxImage()

	e.GetAzureSubscriptionID()
	e.GetAzureClientID()
	e.GetAzureClientSecret()
	e.GetAzureTenantID()
	e.GetAzureLocationDefault()
}
