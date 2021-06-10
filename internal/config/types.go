package config

import kc "github.com/probr/probr-sdk/providers/kubernetes/config"

type varOptions struct {
	VarsFile     string
	Verbose      bool
	ServicePacks servicePacks `yaml:"ServicePacks"`
}

type servicePacks struct {
	Kubernetes kubernetes `yaml:"Kubernetes"`
}

type kubernetes struct {
	kc.Kubernetes                     `yaml:",inline"`
	SystemClusterRoles                []string `yaml:"SystemClusterRoles"`
	UnauthorisedContainerImage        string   `yaml:"UnauthorisedContainerImage"`
	ProbeImage                        string   `yaml:"ProbeImage"`
	ContainerRequiredDropCapabilities []string `yaml:"ContainerRequiredDropCapabilities"`
	ContainerAllowedAddCapabilities   []string `yaml:"ContainerAllowedAddCapabilities"`
	ApprovedVolumeTypes               []string `yaml:"ApprovedVolumeTypes"`
	UnapprovedHostPort                string   `yaml:"UnapprovedHostPort"`
	SystemNamespace                   string   `yaml:"SystemNamespace"`
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
