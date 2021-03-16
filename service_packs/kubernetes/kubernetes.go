package kubernetes

import (
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	cra "github.com/citihub/probr/service_packs/kubernetes/container_registry_access"
	"github.com/citihub/probr/service_packs/kubernetes/general"
	"github.com/citihub/probr/service_packs/kubernetes/iam"
	"github.com/citihub/probr/service_packs/kubernetes/podsecurity"
	"github.com/markbates/pkger"
)

// GetProbes returns a list of probe objects
func GetProbes() []coreengine.Probe {
	if config.Vars.ServicePacks.Kubernetes.IsExcluded() {
		return nil
	}
	return []coreengine.Probe{
		cra.Probe,
		general.Probe,
		podsecurity.Probe,
		iam.Probe,
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/kubernetes/container_registry_access/container_registry_access.feature")
	pkger.Include("/service_packs/kubernetes/general/general.feature")
	pkger.Include("/service_packs/kubernetes/podsecurity/podsecurity.feature")
	pkger.Include("/service_packs/kubernetes/iam/iam.feature")
}
