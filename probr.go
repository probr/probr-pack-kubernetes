package probr

import (
	"os"

	"github.com/citihub/probr/config"
	servicepacks "github.com/citihub/probr/service_packs"
	"github.com/citihub/probr/service_packs/coreengine"
)

var tmpDirFunc = config.Vars.TmpDir // TODO: revise this

// RunAllProbes retrieves and executes all probes that have been included
func RunAllProbes() (int, *coreengine.ProbeStore, error) {
	ts := coreengine.NewProbeStore()

	for _, probe := range servicepacks.GetAllProbes() {
		ts.AddProbe(probe)
	}

	s, err := ts.ExecAllProbes() // Executes all added (queued) tests
	return s, ts, err
}

//GetAllProbeResults maps ProbeStore results to strings
func GetAllProbeResults(ps *coreengine.ProbeStore) map[string]string {
	defer CleanupTmp()

	out := make(map[string]string)
	for name := range ps.Probes {
		results, name, err := readProbeResults(ps, name)
		if err != nil {
			out[name] = err.Error()
		} else {
			out[name] = results
		}
	}
	return out
}

func readProbeResults(ps *coreengine.ProbeStore, name string) (probeResults, probeName string, err error) {
	p, err := ps.GetProbe(name)
	if err != nil {
		return
	}
	probeResults = p.Results.String()
	probeName = p.ProbeDescriptor.Name
	return
}

// CleanupTmp is used to dispose of any temp resources used during execution
func CleanupTmp() {
	// Remove tmp folder and its content
	err := os.RemoveAll(tmpDirFunc())
	if err != nil {
		//log.Printf("[ERROR] Error removing tmp folder %v", err)
	}
}
