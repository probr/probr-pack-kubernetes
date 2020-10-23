package coreengine_test

import (
	"testing"

	"github.com/citihub/probr/internal/coreengine"
	"github.com/stretchr/testify/assert"
)

// Test manager integration tests, so actually calling out to Kube ...
// TODO: complete

func TestExecTest(t *testing.T) {
	// create a test and add it to the TestStore

	//test descriptor ... (general)
	grp := coreengine.CloudDriver
	name := "account_manager"
	td := coreengine.TestDescriptor{Group: grp, Name: name}

	tm := coreengine.NewTestManager() // get the test mgr

	assert.NotNil(t, tm)

	tsuuid := tm.AddTest(td)

	s, err := tm.ExecTest(tsuuid)
	if err != nil {
		t.Fatalf("Error executing test: %v", err)
	}

	assert.True(t, s == 0)
}

func TestExecAllTests(t *testing.T) {

	tm := coreengine.NewTestManager()

	//add some tests and add them to the TM
	addTest(tm, "account_manager", coreengine.CloudDriver)
	addTest(tm, "pod_security_policy", coreengine.Kubernetes)
	addTest(tm, "internet_access", coreengine.Kubernetes)

	tm.ExecAllTests()
}

func addTest(tm *coreengine.TestStore, testname string, grp coreengine.Group) {
	td := coreengine.TestDescriptor{Group: grp, Name: testname}
	tm.AddTest(td)
}
