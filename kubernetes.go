package pack

import (
	cra "github.com/citihub/probr-pack-kubernetes/internal/container_registry_access"
	"github.com/citihub/probr-pack-kubernetes/internal/general"
	"github.com/citihub/probr-pack-kubernetes/internal/iam"
	"github.com/citihub/probr-pack-kubernetes/internal/podsecurity"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/markbates/pkger"
)

// GetProbes returns a list of probe objects
func GetProbes() []probeengine.Probe {
	if config.Vars.ServicePacks.Kubernetes.IsExcluded() {
		return nil
	}
	return []probeengine.Probe{
		cra.Probe,
		general.Probe,
		podsecurity.Probe,
		iam.Probe,
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/internal/container_registry_access/container_registry_access.feature")
	pkger.Include("/internal/general/general.feature")
	pkger.Include("/internal/podsecurity/podsecurity.feature")
	pkger.Include("/internal/iam/iam.feature")
}
