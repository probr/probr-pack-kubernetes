package service_packs

import (
	"path/filepath"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/utils"
	kubernetes_pack "github.com/citihub/probr/service_packs/kubernetes/pack"
	storage_pack "github.com/citihub/probr/service_packs/storage/pack"
)

func packs() (packs map[string][]coreengine.Probe) {
	packs = make(map[string][]coreengine.Probe)

	// Kubernetes pack requires the following vars:
	//   AuthorisedContainerRegistry, UnauthorisedContainerRegistry
	packs["kubernetes"] = kubernetes_pack.GetProbes()

	// Storage pack requires the following vars:
	//   Provider
	packs["storage"] = storage_pack.GetProbes()
	return
}

func makeGodogProbe(pack string, p coreengine.Probe) *coreengine.GodogProbe {
	box := utils.BoxStaticFile(pack+p.Name(), "service_packs", pack, p.Name()) // Establish static files for binary build
	descriptor := coreengine.ProbeDescriptor{Group: coreengine.Kubernetes, Name: p.Name()}
	path := filepath.Join(box.ResolutionDir, p.Name()+".feature")
	return &coreengine.GodogProbe{
		ProbeDescriptor:     &descriptor,
		ProbeInitializer:    p.ProbeInitialize,
		ScenarioInitializer: p.ScenarioInitialize,
		FeaturePath:         path,
	}
}

func GetAllProbes() []*coreengine.GodogProbe {
	var allProbes []*coreengine.GodogProbe

	for packName, pack := range packs() {
		for _, probe := range pack {
			allProbes = append(allProbes, makeGodogProbe(packName, probe))
		}
	}
	return allProbes
}
