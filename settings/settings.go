package settings

import (
	"log"
	"os"

	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/utils"
)

// PackSettings ...
type PackSettings struct {
	Global     *config.VarOptions
	MySettings Kubernetes
}

// Kubernetes config options
type Kubernetes struct {
	KeepPods                          string   `yaml:"KeepPods"` // TODO: Change type to bool, this would allow us to remove logic from kubernetes.GetKeepPodsFromConfig()
	KubeConfigPath                    string   `yaml:"KubeConfig"`
	KubeContext                       string   `yaml:"KubeContext"`
	SystemClusterRoles                []string `yaml:"SystemClusterRoles"`
	AuthorisedContainerRegistry       string   `yaml:"AuthorisedContainerRegistry"`
	UnauthorisedContainerRegistry     string   `yaml:"UnauthorisedContainerRegistry"`
	ProbeImage                        string   `yaml:"ProbeImage"`
	ContainerRequiredDropCapabilities []string `yaml:"ContainerRequiredDropCapabilities"`
	ContainerAllowedAddCapabilities   []string `yaml:"ContainerAllowedAddCapabilities"`
	ApprovedVolumeTypes               []string `yaml:"ApprovedVolumeTypes"`
	UnapprovedHostPort                string   `yaml:"UnapprovedHostPort"`
	SystemNamespace                   string   `yaml:"SystemNamespace"`
	ProbeNamespace                    string   `yaml:"ProbeNamespace"`
	DashboardPodNamePrefix            string   `yaml:"DashboardPodNamePrefix"`
	Azure                             K8sAzure `yaml:"Azure"`
}

// K8sAzure contains Azure-specific options for the Kubernetes service pack
type K8sAzure struct {
	DefaultNamespaceAIB string
	IdentityNamespace   string
}

// NewSettings ...
func NewSettings() PackSettings {

	return PackSettings{
		Global:     &config.Vars,
		MySettings: Kubernetes{},
	}
}

// Load populate all config settings
func (s *PackSettings) Load() {
	// TODO: Make this a singleton
	log.Printf("[INFO] Loading settings")

	// Priority shall be:
	// 1. Default Values
	// 2. Env Variables
	// 3. Config file
	// 4. CLI args

	// Global settings
	var globalSettingsErr error
	s.Global, globalSettingsErr = loadGlobalSettings(*RunCliFlags.VarsFile, RunCliFlags)
	if globalSettingsErr != nil {
		log.Fatalf("[ERROR] Unexpected error occured: %v", globalSettingsErr)
	}

	// My settings
	var mySettingsErr error
	s.MySettings, mySettingsErr = loadMySettings(*RunCliFlags.VarsFile, RunCliFlags)
	if globalSettingsErr != nil {
		log.Fatalf("[ERROR] Unexpected error occured: %v", mySettingsErr)
	}
}

func loadGlobalSettings(configFile string, cliArgs RunFlags) (globalSettings *config.VarOptions, err error) {

	// Initialize config
	configErr := config.Init(configFile)
	if configErr != nil {
		err = utils.ReformatError("Error occurred while loading global config settings: %v", configErr)
		return
	}

	// Apply CLI Args if any
	if *RunCliFlags.VarsFile != "" {
		config.Vars.VarsFile = *RunCliFlags.VarsFile
	}
	if *RunCliFlags.WriteDirectory != "" {
		config.Vars.WriteDirectory = *RunCliFlags.WriteDirectory
	}
	if *RunCliFlags.LogLevel != "" {

		levels := []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"}
		_, found := utils.FindString(levels, *RunCliFlags.LogLevel)
		if !found {
			err = utils.ReformatError("Unexpected value provided for loglevel: '%s' Expected values: %v", *RunCliFlags.LogLevel, levels)
			return
		}

		config.Vars.LogLevel = *RunCliFlags.LogLevel
		config.SetLogFilter(config.Vars.LogLevel, os.Stderr) // TODO: Remove this line once SDK has been updated. Once this is removed, the pacth in main can be removed as well since log output will not be overriden.
		// TODO: logging.SetLoglevel(level)
	}
	if *RunCliFlags.ResultsFormat != "" {
		options := []string{"cucumber", "events", "junit", "pretty", "progress"}
		_, found := utils.FindString(options, *RunCliFlags.ResultsFormat)
		if !found {
			err = utils.ReformatError("Unexpected value provided for resultsformat: '%s' Expected values: %v", *RunCliFlags.ResultsFormat, options)
			return
		}

		config.Vars.ResultsFormat = *RunCliFlags.ResultsFormat
		config.SetLogFilter(config.Vars.ResultsFormat, os.Stderr) // TODO: Taken from cliflags.ResultsformatHandler. Clarify this, looks like a copy/paste bug.
	}
	if *RunCliFlags.Tags != "" {
		config.Vars.Tags = *RunCliFlags.Tags
	}
	if *RunCliFlags.KubeConfig != "" {
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath = *RunCliFlags.KubeConfig // TODO: Extract this to PackSettings
	}

	return &config.Vars, err
}

func loadMySettings(configFile string, cliArgs RunFlags) (mySettings Kubernetes, err error) {

	// TODO: Transform this into loading config file and populating my settings only
	// TODO: Create generic logic to parse config file and allow passing a Settings oject (interface? generics?)
	// Initialize config
	// configErr := config.Init(configFile)
	// if configErr != nil {
	// 	err = utils.ReformatError("Error occurred while loading global config settings: %v", configErr)
	// 	return
	// }

	// TODO: This is temporary. Remove once above logic to read from file has been added.
	mySettings = Kubernetes{
		KeepPods:                          config.Vars.ServicePacks.Kubernetes.KeepPods,
		KubeConfigPath:                    config.Vars.ServicePacks.Kubernetes.KubeConfigPath,
		KubeContext:                       config.Vars.ServicePacks.Kubernetes.KubeContext,
		SystemClusterRoles:                config.Vars.ServicePacks.Kubernetes.SystemClusterRoles,
		AuthorisedContainerRegistry:       config.Vars.ServicePacks.Kubernetes.AuthorisedContainerRegistry,
		UnauthorisedContainerRegistry:     config.Vars.ServicePacks.Kubernetes.UnauthorisedContainerRegistry,
		ProbeImage:                        config.Vars.ServicePacks.Kubernetes.ProbeImage,
		ContainerRequiredDropCapabilities: config.Vars.ServicePacks.Kubernetes.ContainerRequiredDropCapabilities,
		ContainerAllowedAddCapabilities:   config.Vars.ServicePacks.Kubernetes.ContainerAllowedAddCapabilities,
		ApprovedVolumeTypes:               config.Vars.ServicePacks.Kubernetes.ApprovedVolumeTypes,
		UnapprovedHostPort:                config.Vars.ServicePacks.Kubernetes.UnapprovedHostPort,
		SystemNamespace:                   config.Vars.ServicePacks.Kubernetes.SystemNamespace,
		ProbeNamespace:                    config.Vars.ServicePacks.Kubernetes.ProbeNamespace,
		DashboardPodNamePrefix:            config.Vars.ServicePacks.Kubernetes.DashboardPodNamePrefix,
		Azure: K8sAzure{
			DefaultNamespaceAIB: config.Vars.ServicePacks.Kubernetes.Azure.DefaultNamespaceAIB,
			IdentityNamespace:   config.Vars.ServicePacks.Kubernetes.Azure.IdentityNamespace,
		},
	}

	// Apply CLI Args if any
	if *RunCliFlags.KubeConfig != "" {
		mySettings.KubeConfigPath = *RunCliFlags.KubeConfig
	}

	return mySettings, err
}
