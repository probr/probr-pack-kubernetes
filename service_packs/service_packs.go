package service_packs

import (
	"path/filepath"

	"github.com/cucumber/godog"
	packr "github.com/gobuffalo/packr/v2"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes/container_registry_access"
	"github.com/citihub/probr/service_packs/kubernetes/general"
	"github.com/citihub/probr/service_packs/kubernetes/iam"
	"github.com/citihub/probr/service_packs/kubernetes/internet_access"
	"github.com/citihub/probr/service_packs/kubernetes/pod_security_policy"
)

type probe interface {
	ProbeInitialize(*godog.TestSuiteContext)
	ScenarioInitialize(*godog.ScenarioContext)
	Name() string
}

var packrBox *packr.Box
var packs map[string][]probe

func init() {
	packrBox = packr.New("box", "") // Allows static filepaths within this directory to be referenced even within the binary
	packs = make(map[string][]probe)
	packs["kubernetes"] = []probe{
		container_registry_access.Probe,
		general.Probe,
		pod_security_policy.Probe,
		internet_access.Probe,
		iam.Probe,
	}
}

func makeGodogProbe(pack string, p probe) *coreengine.GodogProbe {
	pd := coreengine.ProbeDescriptor{Group: coreengine.Kubernetes, Name: p.Name()}
	return &coreengine.GodogProbe{
		ProbeDescriptor:     &pd,
		ProbeInitializer:    p.ProbeInitialize,
		ScenarioInitializer: p.ScenarioInitialize,
		FeaturePath:         filepath.Join(packrBox.ResolutionDir, pack, p.Name(), p.Name()+".feature"),
	}
}

func GetAllProbes() []*coreengine.GodogProbe {
	var allProbes []*coreengine.GodogProbe

	for packName, pack := range packs {
		for _, probe := range pack {
			allProbes = append(allProbes, makeGodogProbe(packName, probe))
		}
	}
	return allProbes
}
