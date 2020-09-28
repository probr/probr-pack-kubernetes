package general

import (	
	"os"
	"flag"
	"log"
	"testing"	
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

// This is the "test" implementation of the associated "_feature" file and allows execution
// of the godog test from the command line using regular "go test" syntax. This is useful
// for exercising specific test cases directly.

var (
	opts = godog.Options{Output: colors.Colored(os.Stdout)}
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)


func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

func TestMain(m *testing.M) {	
	flag.Parse()
	opts.Paths = flag.Args()	

	//TODO: for now, skip if integration flag isn't set
	//need to figure out how to set the kube config in the CI pipeline 
	//before this can be run in the pipeline
	if ! *integrationTest {
		//skip
		log.Print("[NOTICE] general_test: Integration Test Flag not set. SKIPPING TEST.")		
		return
	}

	// godog testing (v0.10.0 (latest))
	status := godog.TestSuite{
		Name: "general",		
		TestSuiteInitializer: TestSuiteInitialize,
		ScenarioInitializer:  ScenarioInitialize,
		Options: &opts,
	}.Run()

	// TestMain may have been invoked as part of a "go test" call
	// so we need to ensure the standard/non-godog tests are also run
	// go testing
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}