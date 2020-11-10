package k8s_probes_test

import (
	"os"
	"path/filepath"
	"testing"

	k8s_probes "github.com/citihub/probr/probes/kubernetes"
)

func TestSpecifications(t *testing.T) {
	if k8s_probes.Specifications == nil {
		t.Logf("Packr box for Specifications was not successfully created.")
		t.Fail()
	}

	var files []string
	dir := "probe_specifications"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		t.Log("Probe specifications were not found in the expected location.")
		t.Fail()
	}

	if len(files) != len(k8s_probes.Specifications.List()) {
		t.Logf("Files within probe specifications directory do not match the files in the packr box:\n %v\n\n%v", files, k8s_probes.Specifications.List())
		t.Fail()
	}
}

func TestProbes(t *testing.T) {
	if len(k8s_probes.Probes) != len(k8s_probes.Specifications.List()) {
		t.Log("The number of probes does not match the number of probe specifications")
		t.Fail()
	}
}

func TestGetGodogProbe(t *testing.T) {
	p := k8s_probes.Probes[0].GetGodogProbe()
	if p.ProbeDescriptor == nil {
		t.Log("Probe Descriptor was not properly set")
		t.Fail()
	}
	if p.ProbeInitializer == nil {
		t.Log("Probe Initializer was not properly set")
		t.Fail()
	}
	if p.ScenarioInitializer == nil {
		t.Log("Probe's Scenario Initializer was not properly set")
		t.Fail()
	}
	if p.FeaturePath == "" {
		t.Log("Probe feature path was not properly set")
		t.Fail()
	}
}
