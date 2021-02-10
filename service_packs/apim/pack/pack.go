package apim_pack

import (
	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/apim/azure/endpoint_security"
	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/markbates/pkger"
)

func GetProbes() []coreengine.Probe {
	//	config.Vars.SetTags(tags)
	if config.Vars.ServicePacks.APIM.IsExcluded() {
		return nil
	}
	switch config.Vars.ServicePacks.APIM.Provider {
	case "Azure":
		return []coreengine.Probe{
			endpoint_security.Probe,
		}
	default:
		return nil
	}
}

func init() {
	// This line will ensure that all static files are bundled into pked.go file when using pkger cli tool
	// See: https://github.com/markbates/pkger
	pkger.Include("/service_packs/apim/azure/endpoint_security/endpoint_security.feature")

}
