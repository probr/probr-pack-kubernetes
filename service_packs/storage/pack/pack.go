package storage_pack

import (
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	azure_access_whitelisting "github.com/citihub/probr/service_packs/storage/azure/access_whitelisting"
	azure_encryption_at_rest "github.com/citihub/probr/service_packs/storage/azure/encryption_at_rest"
	azure_encryption_in_flight "github.com/citihub/probr/service_packs/storage/azure/encryption_in_flight"
)

func GetProbes() []coreengine.Probe {
	switch config.Vars.ServicePacks.Storage.Provider {
	case "Azure":
		return []coreengine.Probe{
			azure_access_whitelisting.Probe,
			azure_encryption_at_rest.Probe,
			azure_encryption_in_flight.Probe,
		}
	default:
		return nil
	}
}
