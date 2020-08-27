package main

import (
	"flag"
	"log"
	"os"

	v1 "gitlab.com/citihub/probr/api/v1"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"

	"gitlab.com/citihub/probr/internal/config" //needed for logging
	_ "gitlab.com/citihub/probr/test/features/clouddriver"
	_ "gitlab.com/citihub/probr/test/features/kubernetes/containerregistryaccess" //needed to run init on TestHandlers
	_ "gitlab.com/citihub/probr/test/features/kubernetes/internetaccess"          //needed to run init on TestHandlers
	_ "gitlab.com/citihub/probr/test/features/kubernetes/podsecuritypolicy"       //needed to run init on TestHandlers
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func main() {
	//TODO: this (to line 45) will all move when we merge the change to move to a library
	//just dumping in here for now ...
	k := flag.String("kube", "", "kube config file")
	o := flag.String("outputDir", "", "output directory")
	flag.Parse()

	v1.SetIOPaths(*k, *o)

	log.Printf("[NOTICE] Probr running with environment: ")
	log.Printf("[NOTICE] %v", config.GetEnvConfigInstance())

	if k != nil && len(*k) > 0 {
		log.Printf("[NOTICE] Kube Config has been overridden on command line to: " + *k)
	}
	if o != nil && len(*o) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden on command line to: " + *o)
	}

	//TODO: this is the cli and what will be called on Docker run ...
	//use args to figure out what needs to be run / output paths / etc
	//and call TestManager to make it happen :-)

	//(possibly want to create a separate "cli" file)

	// get all the below from args ... just hard code for now

	//exec 'em all (for now!)
	s, ts, err := v1.RunAllTests()

	if err != nil {
		log.Fatalf("[ERROR] Error executing tests %v", err)
	}
	log.Printf("[NOTICE] Overall test completion status: %v", s)

	out, err := v1.GetAllTestResults(ts)
	if err != nil {
		log.Fatalf("[ERROR] Experienced error getting test results: %v", s)
	}
	for k, _ := range out {
		log.Printf("Test results in memory: %v", k)
	}
	os.Exit(s)
}
