package storage

import (
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	azureaw "github.com/citihub/probr/service_packs/storage/azure/access_whitelisting"
	azureear "github.com/citihub/probr/service_packs/storage/azure/encryption_at_rest"
	azureeif "github.com/citihub/probr/service_packs/storage/azure/encryption_in_flight"
	"github.com/markbates/pkger"
)

// GetProbes returns a list of probe objects
func GetProbes() []coreengine.Probe {
	if config.Vars.ServicePacks.Storage.IsExcluded() {
		return nil
	}
	switch config.Vars.ServicePacks.Storage.Provider {
	case "Azure":
		return []coreengine.Probe{
			azureaw.Probe,
			azureear.Probe,
			azureeif.Probe,
		}
	default:
		return nil
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/storage/azure/access_whitelisting/access_whitelisting.feature")
	pkger.Include("/service_packs/storage/azure/encryption_at_rest/encryption_at_rest.feature")
	pkger.Include("/service_packs/storage/azure/encryption_in_flight/encryption_in_flight.feature")
}
