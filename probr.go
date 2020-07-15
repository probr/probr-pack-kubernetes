package main

import (
	"log"
	"os"

	"citihub.com/probr/internal/coreengine"
	"github.com/google/uuid"

	_ "citihub.com/probr/test/features/clouddriver"
	_ "citihub.com/probr/test/features/kubernetes/podsecuritypolicy" //needed to run init on TestHandlers
)

func main() {
	//TODO: this is the cli and what will be called on Docker run ...
	//use args to figure out what needs to be run / output paths / etc
	//and call TestManager to make it happen :-)

	//(possibly want to create a separate "cli" file)

	// get all the below from args ... just hard code for now 
	
	// get the test mgr
	tm := coreengine.NewTestManager()

	//add some tests and add them to the TM - we need to tidy this up!
	addTest(tm, "account_manager", coreengine.General)
	addTest(tm, "pod_security_policy", coreengine.PodSecurityPolicies)

	//exec 'em all (for now!)
	s, err := tm.ExecAllTests()
	if err != nil {
		log.Fatalf("Error executing tests %v", err)
	}

	os.Exit(s)

}

func addTest(tm *coreengine.TestStore, testname string, category coreengine.Category) {

	cat := category	
	name := testname
	td := coreengine.TestDescriptor {Category: cat, Name: name}

	uuid1 := uuid.New().String()	
	sat := coreengine.Pending	

	test := coreengine.Test {
		UUID: &uuid1,
		TestDescriptor: &td,
		Status: &sat,
	}

	//add - don't worry about the rtn uuid
	tm.AddTest(&test)
}
