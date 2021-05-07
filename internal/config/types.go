package config

type varOptions struct {
	VarsFile     string
	Verbose      bool
	ServicePacks servicePacks `yaml:"ServicePacks"`
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
	TagInclusions                     []string `yaml:"TagInclusions"`
	TagExclusions                     []string `yaml:"TagExclusions"`
}

// K8sAzure contains Azure-specific options for the Kubernetes service pack
type k8sAzure struct {
	DefaultNamespaceAIB string
	IdentityNamespace   string
}
