package config

import (
	"fmt"
	"strings"
	"testing"
)

//
// Helpers
//

func assertIsNotExcluded(obj Excludable, t *testing.T) {
	if obj.IsExcluded() {
		t.Log("Non-excluded probe has been exlcuded")
		t.Fail()
	}
}

func assertIsExcluded(obj Excludable, t *testing.T) {
	if !obj.IsExcluded() {
		t.Log("Excluded probe has not been exlcuded")
		t.Fail()
	}
}

func newConfigWithScenarioExclusionAndInclusion() (config ConfigVars, excludedTag string) {
	config, _ = NewConfig("")
	excludedTag = "that guy"
	config.ServicePacks.Kubernetes.Probes = append(
		config.ServicePacks.Kubernetes.Probes,
		Probe{
			Name:      "container_registry_access",
			Scenarios: []Scenario{Scenario{Name: "this guy"}},
		},
	)
	config.ServicePacks.Kubernetes.Probes[0].Scenarios = append(
		config.ServicePacks.Kubernetes.Probes[0].Scenarios,
		Scenario{
			Name:     excludedTag,
			Excluded: "yes",
		},
	)
	return
}

// Checks scenarios excluded by newConfigWithScenarioExclusionAndInclusion
func checkPreformattedScenarioExclusions(config ConfigVars, t *testing.T) {
	assertIsNotExcluded(config.ServicePacks.Kubernetes.Probes[0].Scenarios[0], t)
	assertIsExcluded(config.ServicePacks.Kubernetes.Probes[0].Scenarios[1], t)
}

func checkTagsContainExclusion(config ConfigVars, tag string, t *testing.T) {
	// Only log one of these possible failures, they won't overlap
	if len(config.Tags) == 0 {
		t.Log("Tags string was not modified by addExclusion")
		t.Fail()
	} else {
		if !strings.Contains(config.Tags, tag) {
			t.Log(fmt.Sprintf("Tags does not contain '%s'", tag))
			t.Fail()
		} else if !strings.Contains(config.Tags, "~@") {
			t.Log("Tags does not contain exclusion syntax")
			t.Fail()
		}
	}
}

//
// Tests
//

func TestNewConfig(t *testing.T) {
	// Just use a default config, no file-read for now
	config, err := NewConfig("")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	switch interface{}(config).(type) {
	case ConfigVars:
	default:
		t.Log("NewConfig did not create a ConfigVars object")
		t.Fail()
	}
}

func TestK8sIsExcluded(t *testing.T) {
	config, _ := NewConfig("")

	// 0 required vars set
	assertIsExcluded(config.ServicePacks.Kubernetes, t)

	// 1 of 2 required vars set
	config.ServicePacks.Kubernetes.AuthorisedContainerRegistry = "not-empty"
	assertIsExcluded(config.ServicePacks.Kubernetes, t)

	// All required vars set
	config.ServicePacks.Kubernetes.UnauthorisedContainerRegistry = "not-empty"
	assertIsNotExcluded(config.ServicePacks.Kubernetes, t)
}

func TestStorageIsExcluded(t *testing.T) {
	config, _ := NewConfig("")

	// 0 required vars set
	assertIsExcluded(config.ServicePacks.Storage, t)

	// All required vars set
	config.ServicePacks.Storage.Provider = "not-empty"
	assertIsNotExcluded(config.ServicePacks.Storage, t)
}

func TestProbeIsExcluded(t *testing.T) {
	config, _ := NewConfig("")
	config.ServicePacks.Kubernetes.Probes = append(
		config.ServicePacks.Kubernetes.Probes,
		Probe{Name: "container_registry_access"},
	)
	config.ServicePacks.Kubernetes.Probes = append(
		config.ServicePacks.Kubernetes.Probes,
		Probe{
			Name:     "iam",
			Excluded: "yes",
		},
	)
	assertIsNotExcluded(config.ServicePacks.Kubernetes.Probes[0], t)
	assertIsExcluded(config.ServicePacks.Kubernetes.Probes[1], t)
}
func TestScenarioIsExcluded(t *testing.T) {
	config, _ := newConfigWithScenarioExclusionAndInclusion()
	checkPreformattedScenarioExclusions(config, t)
}

func TestHandleProbeExclusions(t *testing.T) {
	config, excludedTag := newConfigWithScenarioExclusionAndInclusion()
	config.handleProbeExclusions("kubernetes", config.ServicePacks.Kubernetes.Probes)
	checkPreformattedScenarioExclusions(config, t) // verify that exclusions weren't modified
	checkTagsContainExclusion(config, excludedTag, t)
}

func TestAddExclusion(t *testing.T) {
	config, _ := NewConfig("")
	value := "exclude-this-tag"
	tag := "~@" + value
	config.addExclusion(value)
	checkTagsContainExclusion(config, tag, t)
}

func TestSetTags(t *testing.T) {
	vars, _ := NewConfig("")
	tagName := "@probes/kubernetes"
	tagID := "k-cra"
	expected := 0
	vars.Tags = tagName
	tags := map[string][]string{tagName: []string{tagID}}
	t.Log(tags)
	vars.SetTags(tags)
	for _, tag := range strings.Split(vars.Tags, ",") {
		t.Log("Tag found:" + tag)
		if tag == tagName || tag == "@"+tagID {
			expected = expected + 1
		}
	}
	if expected != 2 {
		t.Errorf("Tag name and tag ID were not found in Vars.Tags")
	}
}

// Pending... these may be too integration-y for a unit test
func TestInit(t *testing.T)                       {}
func TestValidateConfigPath(t *testing.T)         {}
func TestLogConfigState(t *testing.T)             {}
func TestAuditDir(t *testing.T)                   {}
func TestHandleConfigFileExclusions(t *testing.T) {}
