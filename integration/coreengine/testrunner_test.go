package coreengine_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
)

// Test runner integration tests, so actually calling out to Kube ...
// TODO: complete

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

func TestMain(m *testing.M) {
	flag.Parse()

	//TODO: for now, skip if integration flag isn't set
	//need to figure out how to set the kube config in the CI pipeline
	//before this can be run in the pipeline
	if !*integrationTest {
		//skip
		log.Print("[NOTICE] testrunner_test: Integration Test Flag not set. SKIPPING TEST.")
		return
	}

	result := m.Run()

	os.Exit(result)
}

func TestTestRunner(t *testing.T) {

	tr := coreengine.TestStore{}

	//test descriptor ... (general)
	grp := coreengine.CloudDriver
	name := "account_manager"
	td := coreengine.TestDescriptor{Group: grp, Name: name}

	//specific terms for *this* test
	sat1 := coreengine.Pending

	//construct the test to run
	test1 := coreengine.GodogTest{
		TestDescriptor: &td,
		Status:         &sat1,
	}

	assert.NotNil(t, test1)

	st, err := tr.RunTest(&test1)

	//not expecting this one to fail:
	assert.Nil(t, err)
	assert.Equal(t, 0, st)
	assert.Equal(t, coreengine.CompleteSuccess, *test1.Status)

	//update the name and watch it fail:
	td.Name = "not_a_test"
	st, err = tr.RunTest(&test1)

	//expecting an error in this case:
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("no test handler available for %v - cannot run test", td), err.Error())
	assert.Equal(t, 4, st)
	assert.Equal(t, coreengine.Error, *test1.Status)

	//test for nil descriptor
	test1.TestDescriptor = nil

	st, err = tr.RunTest(&test1)

	//expecting an error in this case:
	assert.NotNil(t, err)
	assert.Equal(t, "test descriptor is nil - cannot run test", err.Error())
	assert.Equal(t, 3, st)
	assert.Equal(t, coreengine.Error, *test1.Status)

	//try another one ..
	//test descriptor ... (general)
	name2 := "internet_access"
	td2 := coreengine.TestDescriptor{Name: name2}

	//specific terms for *this* test
	sat2 := coreengine.Pending

	//construct the test to run
	test2 := coreengine.GodogTest{
		TestDescriptor: &td2,
		Status:         &sat2,
	}

	assert.NotNil(t, test2)

	st, err = tr.RunTest(&test2)

	//now testing against an evironment which should have the correct
	//network access rules, hence this test should succeed
	assert.Nil(t, err)
	assert.Equal(t, 0, st)
	assert.Equal(t, coreengine.CompleteSuccess, *test2.Status)

}

func TestTestRunnerInMem(t *testing.T) {
	config.Vars.OutputType = "INMEM"

	tr := coreengine.TestStore{}

	//test descriptor ... (general)
	grp := coreengine.CloudDriver
	name := "account_manager"
	td := coreengine.TestDescriptor{Group: grp, Name: name}

	//specific terms for *this* test
	sat1 := coreengine.Pending

	//construct the test to run
	test1 := coreengine.GodogTest{
		TestDescriptor: &td,
		Status:         &sat1,
	}

	assert.NotNil(t, test1)

	st, err := tr.RunTest(&test1)

	//not expecting this one to fail:
	assert.Nil(t, err)
	assert.Equal(t, 0, st)
	assert.Equal(t, coreengine.CompleteSuccess, *test1.Status)
}
