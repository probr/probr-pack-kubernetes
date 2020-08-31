package main

import (
	"flag"
	"log"
	"os"

	"gitlab.com/citihub/probr"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"

	"gitlab.com/citihub/probr/internal/config" //needed for logging
	// _ "gitlab.com/citihub/probr/test/features/clouddriver"
	// _ "gitlab.com/citihub/probr/test/features/kubernetes/containerregistryaccess" //needed to run init on TestHandlers
	// _ "gitlab.com/citihub/probr/test/features/kubernetes/internetaccess"          //needed to run init on TestHandlers
	// _ "gitlab.com/citihub/probr/test/features/kubernetes/podsecuritypolicy"       //needed to run init on TestHandlers
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func main() {
	var v string
	var ot string
	flag.StringVar(&v, "varsFile", "", "path to config file")
	flag.StringVar(&ot, "outputType", "INMEM", "output defaults to write in memory, if 'IO' will write to specified output directory")
	i := flag.String("kubeConfig", "", "kube config file")
	o := flag.String("outputDir", "", "output directory")
	flag.Parse()

	if ot == "IO" {
		probr.SetIOPaths(*i, *o)
	}

	if i != nil && len(*i) > 0 {
		log.Printf("[NOTICE] Kube Config has been overridden via command line to: " + *i)
	}
	if o != nil && len(*o) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden via command line to: " + *o)
	}

	// Will make config.Vars.XYZ available for the rest of the runtime
	err := config.Init(v)
	if err != nil {
		log.Fatalf("[ERROR] Could not create config from provided filepath: %v", err)
	}
	log.Printf("[NOTICE] Probr running with environment: ")
	log.Printf("[NOTICE] %v", config.Vars)

	//exec 'em all (for now!)
	s, ts, err := probr.RunAllTests()

	if err != nil {
		log.Fatalf("[ERROR] Error executing tests %v", err)
	}
	log.Printf("[NOTICE] Overall test completion status: %v", s)

	out, err := probr.GetAllTestResults(ts)
	if err != nil {
		log.Fatalf("[ERROR] Experienced error getting test results: %v", s)
	}
	for k, _ := range out {
		log.Printf("Test results in memory: %v", k)
	}
	os.Exit(s)
}
