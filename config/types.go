package config

// VarOptions contains all top-level config vars
type VarOptions struct {
	// NOTE: Env and Defaults are ONLY available if corresponding logic is added to defaults.go
	ServicePacks              ServicePacks   `yaml:"ServicePacks"`
	CloudProviders            CloudProviders `yaml:"CloudProviders"`
	OutputType                string         `yaml:"OutputType"`
	WriteDirectory            string         `yaml:"WriteDirectory"`
	AuditEnabled              string         `yaml:"AuditEnabled"`
	LogLevel                  string         `yaml:"LogLevel"`
	OverwriteHistoricalAudits string         `yaml:"OverwriteHistoricalAudits"`
	TagExclusions             []string       `yaml:"TagExclusions"`
	WriteConfig               string         `yaml:"WriteConfig"`
	Tags                      string         // set by flags
	VarsFile                  string         // set by flags only
	NoSummary                 bool           // set by flags only
	Silent                    bool           // set by flags only
	Meta                      Meta           // set by CLI options only
	ResultsFormat             string         // set by flags only
}

// Meta config options
type Meta struct {
	RunOnly string // set by CLI 'run' option
}

// ServicePacks config options
type ServicePacks struct {
	Kubernetes Kubernetes `yaml:"Kubernetes"`
	Storage    Storage    `yaml:"Storage"`
	APIM       APIM       `yaml:"APIM"`
}

// Kubernetes config options
type Kubernetes struct {
	exclusionLogged                   bool
	KeepPods                          string   `yaml:"KeepPods"` // TODO: Change type to bool, this would allow us to remove logic from kubernetes.GetKeepPodsFromConfig()
	Probes                            []Probe  `yaml:"Probes"`
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

// Storage service pack config options
type Storage struct {
	exclusionLogged bool
	Provider        string  `yaml:"Provider"` // Placeholder!
	Probes          []Probe `yaml:"Probes"`
}

// APIM service pack config options
type APIM struct {
	exclusionLogged bool
	Provider        string  `yaml:"Provider"` // Placeholder!
	Probes          []Probe `yaml:"Probes"`
}

// Probe config options
type Probe struct {
	Name      string     `yaml:"Name"`
	Excluded  string     `yaml:"Excluded"`
	Scenarios []Scenario `yaml:"Scenarios"`
}

// Scenario config options
type Scenario struct {
	Name     string `yaml:"Name"`
	Excluded string `yaml:"Excluded"`
}

// CloudProviders config options
type CloudProviders struct {
	Azure Azure `yaml:"Azure"`
}

// Azure config options that may be required by any service pack
type Azure struct {
	Excluded         string `yaml:"Excluded"`
	TenantID         string `yaml:"TenantID"`
	SubscriptionID   string `yaml:"SubscriptionID"`
	ClientID         string `yaml:"ClientID"`
	ClientSecret     string `yaml:"ClientSecret"`
	ResourceGroup    string `yaml:"ResourceGroup"`
	ResourceLocation string `yaml:"ResourceLocation"`
	ManagementGroup  string `yaml:"ManagementGroup"`
}

// Excludable is used for testing purposes only
type Excludable interface {
	IsExcluded() bool
}
