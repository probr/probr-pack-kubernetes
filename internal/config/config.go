package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	sdkConfig "github.com/probr/probr-sdk/config"
	"github.com/probr/probr-sdk/config/setter"
	"github.com/probr/probr-sdk/utils"
)

// Vars is a stateful object containing the variables required to execute this pack
var Vars varOptions

// Init will set values with the content retrieved from a filepath, env vars, or defaults
func (ctx *varOptions) Init() (err error) {
	if ctx.VarsFile != "" {
		ctx.decode()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
	} else {
		log.Printf("[DEBUG] No vars file provided, unexpected behavior may occur")
	}
	sdkConfig.GlobalConfig.VarsFile = ctx.VarsFile
	sdkConfig.GlobalConfig.Init()
	ctx.ServicePacks.Kubernetes.setEnvAndDefaults()

	log.Printf("[DEBUG] Config initialized by %s", utils.CallerName(1))
	return
}

// decode uses an SDK helper to create a YAML file decoder,
// parse the file to an object, then extracts the values from
// ServicePacks.Kubernetes into this context
func (ctx *varOptions) decode() (err error) {
	configDecoder, file, err := sdkConfig.NewConfigDecoder(ctx.VarsFile)
	if err != nil {
		return
	}
	err = configDecoder.Decode(&ctx)
	file.Close()
	return err
}

// LogConfigState will write the config file to the write directory
func (ctx *varOptions) LogConfigState() {
	json, _ := json.MarshalIndent(ctx, "", "  ")
	log.Printf("[INFO] Config State: %s", json)
	// path := filepath.Join("config.json")
	// if ctx.WriteConfig == "true" && utils.WriteAllowed(path) {
	// 	data := []byte(json)
	// 	ioutil.WriteFile(path, data, 0644)
	// 	//log.Printf("[NOTICE] Config State written to file %s", path)
	// }
}

func (ctx *varOptions) Tags() string {
	return sdkConfig.ParseTags(ctx.ServicePacks.Kubernetes.TagInclusions, ctx.ServicePacks.Kubernetes.TagExclusions)
}

// setEnvOrDefaults will set value from os.Getenv and default to the specified value
func (ctx *kubernetes) setEnvAndDefaults() {
	// Notes on SetVar's values:
	// 1. Pointer to local object; will be overwritten by env or default if empty
	// 2. Name of env var to check
	// 3. Default value to set if flags, vars file, and env have not provided a value

	setter.SetVar(&ctx.KeepPods, "PROBR_KEEP_PODS", "false")
	setter.SetVar(&ctx.KubeConfigPath, "KUBE_CONFIG", getDefaultKubeConfigPath())
	setter.SetVar(&ctx.KubeContext, "KUBE_CONTEXT", "")
	setter.SetVar(&ctx.SystemClusterRoles, "", []string{"system:", "aks", "cluster-admin", "policy-agent"})
	setter.SetVar(&ctx.AuthorisedContainerImage, "PROBR_AUTHORISED_IMAGE", "")
	setter.SetVar(&ctx.UnauthorisedContainerImage, "PROBR_UNAUTHORISED_IMAGE", "")
	setter.SetVar(&ctx.ProbeImage, "PROBR_PROBE_IMAGE", "citihub/probr-probe")
	setter.SetVar(&ctx.ContainerRequiredDropCapabilities, "PROBR_REQUIRED_DROP_CAPABILITIES", []string{"NET_RAW"})
	setter.SetVar(&ctx.ContainerAllowedAddCapabilities, "PROBR_ALLOWED_ADD_CAPABILITIES", []string{""})
	setter.SetVar(&ctx.ApprovedVolumeTypes, "PROBR_APPROVED_VOLUME_TYPES", []string{"configmap", "emptydir", "persistentvolumeclaim"})
	setter.SetVar(&ctx.UnapprovedHostPort, "PROBR_UNAPPROVED_HOSTPORT", "22")
	setter.SetVar(&ctx.SystemNamespace, "PROBR_K8S_SYSTEM_NAMESPACE", "kube-system")
	setter.SetVar(&ctx.DashboardPodNamePrefix, "PROBR_K8S_DASHBOARD_PODNAMEPREFIX", "kubernetes-dashboard")
	setter.SetVar(&ctx.ProbeNamespace, "PROBR_K8S_PROBE_NAMESPACE", "probr-general-test-ns")
	setter.SetVar(&ctx.Azure.DefaultNamespaceAIB, "DEFAULT_NS_AZURE_IDENTITY_BINDING", "probr-aib")
	setter.SetVar(&ctx.Azure.IdentityNamespace, "PROBR_K8S_AZURE_IDENTITY_NAMESPACE", "kube-system")
}

func getDefaultKubeConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}
