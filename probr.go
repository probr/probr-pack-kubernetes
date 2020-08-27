package main

import (
	"flag"
	"log"
	"os"

	"github.com/google/uuid"
	v1 "gitlab.com/citihub/probr/api/v1"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"

	"gitlab.com/citihub/probr/internal/config" //needed for logging
	"gitlab.com/citihub/probr/test/features"
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

	SetIOPaths(*k, *o)

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

func addTest(tm *coreengine.TestStore, n string, g coreengine.Group, c coreengine.Category) {

	td := coreengine.TestDescriptor{Group: g, Category: c, Name: n}

	uuid1 := uuid.New().String()
	sat := coreengine.Pending

	test := coreengine.Test{
		UUID:           &uuid1,
		TestDescriptor: &td,
		Status:         &sat,
	}

	//add - don't worry about the rtn uuid
	tm.AddTest(&test)
}

// RunAllTests ...
func RunAllTests() (int, error) {
	// get the test mgr
	tm := coreengine.NewTestManager()

	//add some tests and add them to the TM - we need to tidy this up!
	addTest(tm, "container_registry_access", coreengine.Kubernetes, coreengine.ContainerRegistryAccess)
	addTest(tm, "internet_access", coreengine.Kubernetes, coreengine.InternetAccess)
	addTest(tm, "pod_security_policy", coreengine.Kubernetes, coreengine.PodSecurityPolicies)
	addTest(tm, "account_manager", coreengine.CloudDriver, coreengine.General)

	//exec 'em all (for now!)
	return tm.ExecAllTests()

}

// SetIOPaths ...
func SetIOPaths(i string, o string) {
	kube.SetKubeConfigFile(&i)
	features.SetOutputDirectory(&o)
}
