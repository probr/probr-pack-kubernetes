package service_packs

import (
	"path/filepath"

	"github.com/cucumber/godog"
	"github.com/gobuffalo/packr/v2"

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

var packs map[string][]probe

func init() {
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
	box := packr.New(pack+p.Name(), filepath.Join(pack, p.Name())) // Establish static files for binary build
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

	for packName, pack := range packs {
		for _, probe := range pack {
			allProbes = append(allProbes, makeGodogProbe(packName, probe))
		}
	}
	return allProbes
}
