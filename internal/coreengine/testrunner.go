package coreengine

import (
	"bytes"
	"fmt"

	"github.com/citihub/probr/internal/summary"
	"github.com/cucumber/godog"
)

// TestRunner describes the interface that should be implemented to support the execution of tests.
type TestRunner interface {
	RunTest(t *GodogTest) error
}

// TestHandlerFunc describes a callback that should be implemented by test cases in order for TestRunner
// to be able to execute the test case.
type TestHandlerFunc func(t *GodogTest) (int, *bytes.Buffer, error)

// GodogTest encapsulates the specific data that GoDog feature based tests require in order to run.   This
// structure will be passed to the test handler callback.
type GodogTest struct {
	TestDescriptor       *TestDescriptor
	TestSuiteInitializer func(*godog.TestSuiteContext)
	ScenarioInitializer  func(*godog.ScenarioContext)
	FeaturePath          *string
	Status               *TestStatus `json:"status,omitempty"`
	Results              *bytes.Buffer
}

// RunTest runs the test case described by the supplied Test.  It looks in it's test register (the handlers global
// variable) for an entry with the same TestDescriptor as the supplied test.  If found, it uses the provided GodogTest
func (ts *TestStore) RunTest(test *GodogTest) (int, error) {

	if test == nil {
		summary.State.GetProbeLog(test.TestDescriptor.Name).Result = "Internal Error - Test not found"
		return 2, fmt.Errorf("test is nil - cannot run test")
	}

	if test.TestDescriptor == nil {
		//update status
		*test.Status = Error
		summary.State.GetProbeLog(test.TestDescriptor.Name).Result = "Internal Error - Test descriptor not found"
		return 3, fmt.Errorf("test descriptor is nil - cannot run test")
	}

	s, o, err := GodogTestHandler(test)

	if s == 0 {
		// success
		*test.Status = CompleteSuccess
	} else {
		// fail
		*test.Status = CompleteFail
	}

	test.Results = o // If in-mem output provided, store as Results
	return s, err
}
