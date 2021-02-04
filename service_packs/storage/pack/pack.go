package storage_pack

import (
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
	azure_access_whitelisting "github.com/citihub/probr/service_packs/storage/azure/access_whitelisting"
	azure_encryption_at_rest "github.com/citihub/probr/service_packs/storage/azure/encryption_at_rest"
	azure_encryption_in_flight "github.com/citihub/probr/service_packs/storage/azure/encryption_in_flight"
	"github.com/markbates/pkger"
)

func GetProbes() []coreengine.Probe {
	config.Vars.SetTags(tags)
	if config.Vars.ServicePacks.Storage.IsExcluded() {
		return nil
	}
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

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/storage/azure/access_whitelisting/access_whitelisting.feature")
	pkger.Include("/service_packs/storage/azure/encryption_at_rest/encryption_at_rest.feature")
	pkger.Include("/service_packs/storage/azure/encryption_in_flight/encryption_in_flight.feature")
}
