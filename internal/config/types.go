package config

// ConfigVars contains all possible config vars
type ConfigVars struct {
	// NOTE: Env and Defaults are ONLY available if corresponding logic is added to defaults.go and getters.go
	ServicePacks              ServicePacks   `yaml:"ServicePacks"`
	CloudProviders            CloudProviders `yaml:"CloudProviders"`
	OutputType                string         `yaml:"OutputType"`
	CucumberDir               string         `yaml:"CucumberDir"`
	AuditDir                  string         `yaml:"AuditDir"`
	AuditEnabled              string         `yaml:"AuditEnabled"`
	LogLevel                  string         `yaml:"LogLevel"`
	OverwriteHistoricalAudits string         `yaml:"OverwriteHistoricalAudits"`
	TagExclusions             []string       `yaml:"TagExclusions"`
	Tags                      string         // set by flags
	VarsFile                  string         // set by flags only
	NoSummary                 bool           // set by flags only
	Silent                    bool           // set by flags only
}

type ServicePacks struct {
	Kubernetes Kubernetes `yaml:"Kubernetes"`
}

type ServicePack struct {
	Excluded string
	Probes   []Probe
}

type Kubernetes struct {
	ServicePack
	Excluded                      string   `yaml:"Excluded"`
	Probes                        []Probe  `yaml:"Probes"`
	KubeConfigPath                string   `yaml:"KubeConfig"`
	KubeContext                   string   `yaml:"KubeContext"`
	SystemClusterRoles            []string `yaml:"SystemClusterRoles"`
	AuthorisedContainerRegistry   string   `yaml:"AuthorisedContainerRegistry"`
	UnauthorisedContainerRegistry string   `yaml:"UnauthorisedContainerRegistry"`
	ProbeImage                    string   `yaml:"ProbeImage"`
}

type Probe struct {
	Name      string     `yaml:"Name"`
	Excluded  string     `yaml:"Excluded"`
	Scenarios []Scenario `yaml:"Scenarios"`
}

type Scenario struct {
	Name     string `yaml:"Name"`
	Excluded string `yaml:"Excluded"`
}

type CloudProviders struct {
	Azure Azure `yaml:"Azure"`
}

type Azure struct {
	Excluded        string `yaml:"Excluded"`
	SubscriptionID  string `yaml:"SubscriptionID"`
	ClientID        string `yaml:"ClientID"`
	ClientSecret    string `yaml:"ClientSecret"`
	TenantID        string `yaml:"TenantID"`
	LocationDefault string `yaml:"LocationDefault"`
	Identity        struct {
		DefaultNamespaceAI  string `yaml:"DefaultNamespaceAI"`
		DefaultNamespaceAIB string `yaml:"DefaultNamespaceAIB"`
	}
}

// For testing
type Excludable interface {
	isExcluded() bool
}
