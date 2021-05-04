package pack

import (
	cra "github.com/citihub/probr-pack-kubernetes/internal/container_registry_access"
	"github.com/citihub/probr-pack-kubernetes/internal/general"
	"github.com/citihub/probr-pack-kubernetes/internal/iam"
	"github.com/citihub/probr-pack-kubernetes/internal/podsecurity"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/markbates/pkger"
)

// GetProbes returns a list of probe objects
func GetProbes() []probeengine.Probe {
	return []probeengine.Probe{
		cra.Probe,
		general.Probe,
		podsecurity.Probe,
		iam.Probe,
	}
}

func init() {
	// pkger.Include is a no-op that directs the pkger tool to include the desired file or folder.
	pkger.Include("/internal/container_registry_access/container_registry_access.feature")
	pkger.Include("/internal/general/general.feature")
	pkger.Include("/internal/podsecurity/podsecurity.feature")
	pkger.Include("/internal/iam/iam.feature")
}
