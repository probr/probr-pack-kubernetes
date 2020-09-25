package main

import (
	"flag"
	"log"
	"os"

	"gitlab.com/citihub/probr"
	"gitlab.com/citihub/probr/internal/audit"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/config"
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func main() {
	handleFlags()

	log.Printf("[NOTICE] Probr running with environment: ")
	log.Printf("[NOTICE] %+v", config.Vars)

	//exec 'em all (for now!)
	s, ts, err := probr.RunAllTests()
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		os.Exit(2) // Error code 1 is reserved for probe test failures, and should not fail in CI
	}
	log.Printf("[NOTICE] Overall test completion status: %v", s)

	out, err := probr.GetAllTestResults(ts)
	if err != nil {
		log.Printf("[ERROR] Experienced error getting test results: %v", s)
		os.Exit(2) // Error code 1 is reserved for probe test failures, and should not fail in CI
	}
	if config.Vars.OutputType == "IO" && (out == nil || len(out) == 0) {
		log.Printf("[ERROR] Test results not written to file, possibly due to permissions on the specified output directory: %s", config.Vars.OutputDir)
	}
	audit.AuditLog.PrintAudit()
	os.Exit(s)
}
