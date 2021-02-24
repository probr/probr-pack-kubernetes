package servicepacks

import (
	apim_pack "github.com/citihub/probr/service_packs/apim/pack"
	"github.com/citihub/probr/service_packs/coreengine"
	kubernetes_pack "github.com/citihub/probr/service_packs/kubernetes/pack"
	storage_pack "github.com/citihub/probr/service_packs/storage/pack"
)

func packs() (packs map[string][]coreengine.Probe) {
	packs = make(map[string][]coreengine.Probe)

	packs["kubernetes"] = kubernetes_pack.GetProbes()
	packs["storage"] = storage_pack.GetProbes()
	packs["apim"] = apim_pack.GetProbes()

	return
}

func makeGodogProbe(pack string, p coreengine.Probe) *coreengine.GodogProbe {
	descriptor := coreengine.ProbeDescriptor{Group: coreengine.Kubernetes, Name: p.Name()}
	return &coreengine.GodogProbe{
		ProbeDescriptor:     &descriptor,
		ProbeInitializer:    p.ProbeInitialize,
		ScenarioInitializer: p.ScenarioInitialize,
		FeaturePath:         p.Path(),
	}
}

// GetAllProbes returns a list of probes that are ready to be run by Godog
func GetAllProbes() []*coreengine.GodogProbe {
	var allProbes []*coreengine.GodogProbe

	for packName, pack := range packs() {
		for _, probe := range pack {
			allProbes = append(allProbes, makeGodogProbe(packName, probe))
		}
	}
	return allProbes
}
