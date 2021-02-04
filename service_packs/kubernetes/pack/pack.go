package kubernetes_pack

import (
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes/container_registry_access"
	"github.com/citihub/probr/service_packs/kubernetes/general"
	"github.com/citihub/probr/service_packs/kubernetes/iam"
	"github.com/citihub/probr/service_packs/kubernetes/internet_access"
	"github.com/citihub/probr/service_packs/kubernetes/pod_security_policy"
	"github.com/markbates/pkger"
)

func GetProbes() []coreengine.Probe {
	config.Vars.SetTags(tags)
	if config.Vars.ServicePacks.Kubernetes.IsExcluded() {
		return nil
	}
	return []coreengine.Probe{
		container_registry_access.Probe,
		general.Probe,
		pod_security_policy.Probe,
		internet_access.Probe,
		iam.Probe,
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/kubernetes/container_registry_access/container_registry_access.feature")
	pkger.Include("/service_packs/kubernetes/general/general.feature")
	pkger.Include("/service_packs/kubernetes/pod_security_policy/pod_security_policy.feature")
	pkger.Include("/service_packs/kubernetes/internet_access/internet_access.feature")
	pkger.Include("/service_packs/kubernetes/iam/iam.feature")
}
