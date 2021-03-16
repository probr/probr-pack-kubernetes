package storage

import (
	"testing"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/service_packs/coreengine"
)

func TestGetProbes(t *testing.T) {
	pack := make([]coreengine.Probe, 0)
	pack = GetProbes()
	if len(pack) > 0 {
		t.Logf("Unexpected value returned from GetProbes")
		t.Fail()
	}

	config.Vars.ServicePacks.Storage.Provider = "Azure"
	pack = GetProbes()
	if len(pack) == 0 {
		t.Logf("Expected value not returned from GetProbes")
		t.Fail()
	}
}
