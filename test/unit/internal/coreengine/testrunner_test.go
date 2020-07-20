package coreengine_test

import (
	"os"
	"fmt"
	"log"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"

	"citihub.com/probr/internal/coreengine"
	_ "citihub.com/probr/test/features/clouddriver" //needed to run init on TestHandlers
	_ "citihub.com/probr/test/features/kubernetes/podsecuritypolicy" //needed to run init on TestHandlers
	_ "citihub.com/probr/test/features/kubernetes/internetaccess" //needed to run init on TestHandlers
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

func TestMain(m *testing.M) {
	flag.Parse()

	if ! *integrationTest {
		//skip
		log.Print("testrunner_test: Integration Test Flag not set. SKIPPING TEST.")
		return
	}

	result := m.Run()

	os.Exit(result)
}

func TestTestRunner (t *testing.T) {
	
	tr := coreengine.TestStore{}

	//test descriptor ... (general)
	cat := coreengine.General	
	name := "account_manager"
	td := coreengine.TestDescriptor {Category: cat, Name: name}

	//specific terms for *this* test
	uuid1 := uuid.New().String()
	sat1 := coreengine.Pending

	//construct the test to run
	test1 := coreengine.Test {
		UUID: &uuid1,
		TestDescriptor: &td,
		Status: &sat1,
	}

	assert.NotNil(t, test1)

	st, err := tr.RunTest(&test1)

	//not expecting this one to fail:
	assert.Nil(t,err)
	assert.Equal(t, 0, st)
	assert.Equal(t, coreengine.CompleteSuccess, *test1.Status)

	//update the name and watch it fail:	
	td.Name = "not_a_test"
	st, err = tr.RunTest(&test1)

	//expecting an error in this case:
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("no test handler available for %v - cannot run test",td), err.Error() )
	assert.Equal(t, 4, st)
	assert.Equal(t, coreengine.Error, *test1.Status )

	//test for nil descriptor
	test1.TestDescriptor = nil
	
	st, err = tr.RunTest(&test1)

	//expecting an error in this case:
	assert.NotNil(t, err)
	assert.Equal(t, "test descriptor is nil - cannot run test", err.Error() )
	assert.Equal(t, 3, st)
	assert.Equal(t, coreengine.Error, *test1.Status)

	//try another one ..
	//test descriptor ... (general)
	cat2 := coreengine.PodSecurityPolicies	
	name2 := "pod_security_policy"
	td2 := coreengine.TestDescriptor {Category: cat2, Name: name2}

	//specific terms for *this* test
	uuid2 := uuid.New().String()
	sat2 := coreengine.Pending

	//construct the test to run
	test2 := coreengine.Test {
		UUID: &uuid2,
		TestDescriptor: &td2,
		Status: &sat2,
	}

	assert.NotNil(t, test2)

	st, err = tr.RunTest(&test2)

	//not expecting this one to fail:
	assert.Nil(t,err)
	assert.Equal(t, 0, st)
	assert.Equal(t, coreengine.CompleteSuccess, *test2.Status)

}

