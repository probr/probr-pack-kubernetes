package config

type parsedVars struct {
	// Must be set via CLI flag
	VarsFile string

	// Can be set by CLI flags or vars file
	WriteDirectory string
	LogLevel       string
	ResultsFormat  string
	Tags           string
	KubeConfigPath string
	Verbose        bool

	// Set by vars file; defined in SDK
	CloudProviders cloudProviders

	// Set by vars file; must match struct 'kubernetes'
	ApprovedVolumeTypes               []string
	AuthorisedContainerImage          string
	Azure                             k8sAzure
	ContainerRequiredDropCapabilities []string
	ContainerAllowedAddCapabilities   []string
	DashboardPodNamePrefix            string
	KeepPods                          string
	KubeContext                       string
	ProbeImage                        string
	ProbeNamespace                    string
	SystemClusterRoles                []string
	SystemNamespace                   string
	UnapprovedHostPort                string
	UnauthorisedContainerImage        string
}

// VarsFile contains all top-level config vars
type varsFile struct {
	CloudProviders cloudProviders `yaml:"CloudProviders"`
	ServicePacks   servicePacks   `yaml:"ServicePacks"`
	TagExclusions  []string       `yaml:"TagExclusions"`
	TagInclusions  []string       `yaml:"TagInclusions"`
}

type servicePacks struct {
	Kubernetes kubernetes `yaml:"Kubernetes"`
}

type kubernetes struct {
	KeepPods                          string   `yaml:"KeepPods"` // TODO: Change type to bool, this would allow us to remove logic from kubernetes.GetKeepPodsFromConfig()
	KubeConfigPath                    string   `yaml:"KubeConfig"`
	KubeContext                       string   `yaml:"KubeContext"`
	SystemClusterRoles                []string `yaml:"SystemClusterRoles"`
	AuthorisedContainerImage          string   `yaml:"AuthorisedContainerImage"`
	UnauthorisedContainerImage        string   `yaml:"UnauthorisedContainerImage"`
	ProbeImage                        string   `yaml:"ProbeImage"`
	ContainerRequiredDropCapabilities []string `yaml:"ContainerRequiredDropCapabilities"`
	ContainerAllowedAddCapabilities   []string `yaml:"ContainerAllowedAddCapabilities"`
	ApprovedVolumeTypes               []string `yaml:"ApprovedVolumeTypes"`
	UnapprovedHostPort                string   `yaml:"UnapprovedHostPort"`
	SystemNamespace                   string   `yaml:"SystemNamespace"`
	ProbeNamespace                    string   `yaml:"ProbeNamespace"`
	DashboardPodNamePrefix            string   `yaml:"DashboardPodNamePrefix"`
	Azure                             k8sAzure `yaml:"Azure"`
}

type azure struct {
	Excluded         string `yaml:"Excluded"`
	TenantID         string `yaml:"TenantID"`
	SubscriptionID   string `yaml:"SubscriptionID"`
	ClientID         string `yaml:"ClientID"`
	ClientSecret     string `yaml:"ClientSecret"`
	ResourceGroup    string `yaml:"ResourceGroup"`
	ResourceLocation string `yaml:"ResourceLocation"`
	ManagementGroup  string `yaml:"ManagementGroup"`
}

// TODO: move this to SDK
type cloudProviders struct {
	Azure azure `yaml:"Azure"`
}

// K8sAzure contains Azure-specific options for the Kubernetes service pack
type k8sAzure struct {
	DefaultNamespaceAIB string
	IdentityNamespace   string
}
