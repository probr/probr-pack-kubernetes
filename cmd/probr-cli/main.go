package main

import (
	"flag"
	"log"
	"os"

	"github.com/citihub/probr"
	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/summary"
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func main() {
	handleFlags()
	config.LogConfigState()

	//exec 'em all (for now!)
	s, ts, err := probr.RunAllProbes()
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		os.Exit(2) // Error code 1 is reserved for probe test failures, and should not fail in CI
	}
	log.Printf("[NOTICE] Overall test completion status: %v", s)
	summary.State.SetProbrStatus()

	if config.Vars.OutputType == "IO" {
		out, err := probr.GetAllProbeResults(ts)
		if err != nil {
			log.Printf("[ERROR] Experienced error getting test results: %v", s)
			os.Exit(2) // Error code 1 is reserved for probe test failures, and should not fail in CI
		}
		if out == nil || len(out) == 0 {
			log.Printf("[ERROR] Test results not written to file, possibly due to permissions on the specified output directory: %s", config.Vars.CucumberDir)
		}
	}
	summary.State.PrintSummary()
	os.Exit(s)
}
