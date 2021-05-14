package pack

import (
	"github.com/markbates/pkger"
	cra "github.com/probr/probr-pack-kubernetes/internal/container_registry_access"
	"github.com/probr/probr-pack-kubernetes/internal/general"
	"github.com/probr/probr-pack-kubernetes/internal/iam"
	"github.com/probr/probr-pack-kubernetes/internal/podsecurity"
	"github.com/probr/probr-sdk/probeengine"
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
