package storage_pack

import (
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/service_packs/storage/access_whitelisting"
	"github.com/citihub/probr/service_packs/storage/encryption_at_rest"
	"github.com/citihub/probr/service_packs/storage/encryption_in_flight"
)

func GetProbes() []coreengine.Probe {
	if config.Vars.ServicePacks.Storage.IsExcluded() {
		return nil
	}
	return []coreengine.Probe{
		access_whitelisting.Probe,
		encryption_at_rest.Probe,
		encryption_in_flight.Probe,
	}
}
