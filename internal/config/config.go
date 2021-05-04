package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	sdkConfig "github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/utils"
)

// Vars is a stateful object containing the variables required to execute this pack
var Vars parsedVars

// Init will set values with the content retrieved from a filepath, env vars, or defaults
func (ctx *parsedVars) Init() (err error) {
	if ctx.VarsFile != "" {
		ctx.decode()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
	} else {
		log.Printf("[DEBUG] No vars file provided, unexpected behavior may occur")
	}

	ctx.setEnvAndDefaults()

	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))
	return
}

// decode uses an SDK helper to create a YAML file decoder,
// parse the file to an object, then extracts the values from
// ServicePacks.Kubernetes into this context
func (ctx *parsedVars) decode() (err error) {
	var unParsed varsFile
	configDecoder, file, err := sdkConfig.NewConfigDecoder(ctx.VarsFile)
	if err != nil {
		return
	}
	err = configDecoder.Decode(&unParsed)
	file.Close()
	if err != nil {
		return
	}
	JSON, _ := json.MarshalIndent(unParsed.ServicePacks.Kubernetes, "", "  ")
	err = json.Unmarshal(JSON, ctx)
	return err
}

// LogConfigState will write the config file to the write directory
func (ctx *parsedVars) LogConfigState() {
	json, _ := json.MarshalIndent(ctx, "", "  ")
	log.Printf("[INFO] Config State: %s", json)
	// path := filepath.Join("config.json")
	// if ctx.WriteConfig == "true" && utils.WriteAllowed(path) {
	// 	data := []byte(json)
	// 	ioutil.WriteFile(path, data, 0644)
	// 	//log.Printf("[NOTICE] Config State written to file %s", path)
	// }
}

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func (ctx *parsedVars) setEnvAndDefaults() {
	// Notes on SetVar's values:
	// 1. Pointer to local object; will be overwritten by env or default if empty
	// 2. Name of env var to check
	// 3. Default value to set if flags, vars file, and env have not provided a value

	sdkConfig.SetVar(&ctx.Tags, "PROBR_TAGS", "")
	sdkConfig.SetVar(&ctx.WriteDirectory, "PROBR_WRITE_DIRECTORY", "probr_output")
	sdkConfig.SetVar(&ctx.LogLevel, "PROBR_LOG_LEVEL", "DEBUG")
	sdkConfig.SetVar(&ctx.ResultsFormat, "PROBR_RESULTS_FORMAT", "cucumber")

	sdkConfig.SetVar(&ctx.KeepPods, "PROBR_KEEP_PODS", "false")
	sdkConfig.SetVar(&ctx.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	sdkConfig.SetVar(&ctx.KubeContext, "KUBE_CONTEXT", "")
	sdkConfig.SetVar(&ctx.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
	sdkConfig.SetVar(&ctx.AuthorisedContainerRegistry, "PROBR_AUTHORISED_REGISTRY", "")
	sdkConfig.SetVar(&ctx.UnauthorisedContainerRegistry, "PROBR_UNAUTHORISED_REGISTRY", "")
	sdkConfig.SetVar(&ctx.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")
	sdkConfig.SetVar(&ctx.ContainerRequiredDropCapabilities, "PROBR_REQUIRED_DROP_CAPABILITIES", []string{"NET_RAW"})
	sdkConfig.SetVar(&ctx.ContainerAllowedAddCapabilities, "PROBR_ALLOWED_ADD_CAPABILITIES", []string{""})
	sdkConfig.SetVar(&ctx.ApprovedVolumeTypes, "PROBR_APPROVED_VOLUME_TYPES", []string{"configmap", "emptydir", "persistentvolumeclaim"})
	sdkConfig.SetVar(&ctx.UnapprovedHostPort, "PROBR_UNAPPROVED_HOSTPORT", "22")
	sdkConfig.SetVar(&ctx.SystemNamespace, "PROBR_K8S_SYSTEM_NAMESPACE", "kube-system")
	sdkConfig.SetVar(&ctx.DashboardPodNamePrefix, "PROBR_K8S_DASHBOARD_PODNAMEPREFIX", "kubernetes-dashboard")
	sdkConfig.SetVar(&ctx.ProbeNamespace, "PROBR_K8S_PROBE_NAMESPACE", "probr-general-test-ns")
	sdkConfig.SetVar(&ctx.Azure.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-aib")
	sdkConfig.SetVar(&ctx.Azure.IdentityNamespace, "PROBR_K8S_AZURE_IDENTITY_NAMESPACE", "kube-system")

	// TODO: move this logic to SDK
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.TenantID, "AZURE_TENANT_ID", "")
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.SubscriptionID, "AZURE_SUBSCRIPTION_ID", "")
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.ClientID, "AZURE_CLIENT_ID", "")
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.ClientSecret, "AZURE_CLIENT_SECRET", "")
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.ResourceGroup, "AZURE_RESOURCE_GROUP", "")
	sdkConfig.SetVar(&ctx.CloudProviders.Azure.ResourceLocation, "AZURE_RESOURCE_LOCATION", "")
}

func getDefaultKubeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
