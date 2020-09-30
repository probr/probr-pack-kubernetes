package coreengine_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/citihub/probr/internal/coreengine"
)

func TestTestStatus(t *testing.T) {
	ts := coreengine.Pending
	assert.Equal(t, ts.String(), "Pending")

	ts = coreengine.Running
	assert.Equal(t, ts.String(), "Running")

	ts = coreengine.CompleteSuccess
	assert.Equal(t, ts.String(), "CompleteSuccess")

	ts = coreengine.CompleteFail
	assert.Equal(t, ts.String(), "CompleteFail")

	ts = coreengine.Error
	assert.Equal(t, ts.String(), "Error")
}

func TestGetAvailableTests(t *testing.T) {
	alltests := coreengine.GetAvailableTests()

	//not implemented yet, so expect alltests to be nil
	assert.Nil(t, alltests)
}

func TestAddGetTest(t *testing.T) {
	// create a test and add it to the TestStore

	//test descriptor ... (general)
	grp := coreengine.CloudDriver
	cat := coreengine.General
	name := "account_manager"
	td := coreengine.TestDescriptor{Group: grp, Category: cat, Name: name}

	sat1 := coreengine.Pending

	test1 := coreengine.Test{
		TestDescriptor: &td,
		Status:         &sat1,
	}

	assert.NotNil(t, test1)

	// get the test mgr
	tm := coreengine.NewTestManager()

	assert.NotNil(t, tm)

	tsuuid := tm.AddTest(td)

	// now try and get it back ...
	rtntest, err := tm.GetTest(tsuuid)

	assert.Nil(t, err)
	assert.NotNil(t, rtntest)
}

func addTest(tm *coreengine.TestStore, testname string, grp coreengine.Group, cat coreengine.Category) {
	td := coreengine.TestDescriptor{Group: grp, Category: cat, Name: testname}
	tm.AddTest(td)
}
