package servicepacks

import (
	"github.com/citihub/probr-pack-kubernetes/service_packs/kubernetes"
	"github.com/citihub/probr-sdk/probeengine"
)

func packs() (packs map[string][]probeengine.Probe) {
	packs = make(map[string][]probeengine.Probe)

	packs["kubernetes"] = kubernetes.GetProbes()

	return
}

func makeGodogProbe(pack string, p probeengine.Probe) *probeengine.GodogProbe {
	descriptor := probeengine.ProbeDescriptor{Group: probeengine.Kubernetes, Name: p.Name()}
	return &probeengine.GodogProbe{
		ProbeDescriptor:     &descriptor,
		ProbeInitializer:    p.ProbeInitialize,
		ScenarioInitializer: p.ScenarioInitialize,
		FeaturePath:         p.Path(),
	}
}

// GetAllProbes returns a list of probes that are ready to be run by Godog
func GetAllProbes() []*probeengine.GodogProbe {
	var allProbes []*probeengine.GodogProbe

	for packName, pack := range packs() {
		for _, probe := range pack {
			allProbes = append(allProbes, makeGodogProbe(packName, probe))
		}
	}
	return allProbes
}
