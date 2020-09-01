package config

// getEnvOrDefaults will set value from os.Getenv and default to the specified value
func getEnvOrDefaults(e *ConfigVars) {
	e.GetKubeConfigPath(".kube/config") // KUBE_CONFIG
	e.GetOutputType()                   // OUTPUT_TYPE

	e.GetImageRepository() // IMAGE_REPOSITORY
	e.GetCurlImage()       // CURL_IMAGE
	e.GetBusyBoxImage()    // BUSYBOX_IMAGE

	e.GetAzureSubscriptionID()  // AZURE_SUBSCRIPTION_ID
	e.GetAzureClientID()        // AZURE_CLIENT_ID
	e.GetAzureClientSecret()    // AZURE_CLIENT_SECRET
	e.GetAzureTenantID()        // AZURE_TENANT_ID
	e.GetAzureLocationDefault() // AZURE_LOCATION_DEFAULT
}
