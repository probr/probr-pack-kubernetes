package main

import (
	"flag"
	"log"
	"os"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	"citihub.com/probr/internal/coreengine"
	"github.com/google/uuid"

	_ "citihub.com/probr/internal/config" //needed for logging
	"citihub.com/probr/test/features"
	_ "citihub.com/probr/test/features/clouddriver"
	_ "citihub.com/probr/test/features/kubernetes/containerregistryaccess" //needed to run init on TestHandlers
	_ "citihub.com/probr/test/features/kubernetes/internetaccess"          //needed to run init on TestHandlers
	_ "citihub.com/probr/test/features/kubernetes/podsecuritypolicy"       //needed to run init on TestHandlers
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func main() {
	k := flag.String("kube", "", "kube config file")
	o := flag.String("outputDir", "", "output directory")
	flag.Parse()

	SetIOPaths(*k, *o)

	//TODO: this is the cli and what will be called on Docker run ...
	//use args to figure out what needs to be run / output paths / etc
	//and call TestManager to make it happen :-)

	//(possibly want to create a separate "cli" file)

	// get all the below from args ... just hard code for now
	
	//exec 'em all (for now!)
	s, err := RunAllTests()
	if err != nil {
		log.Fatalf("[ERROR] Error executing tests %v", err)
	}

	log.Printf("[NOTICE] Overall test completion status: %v", s)

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