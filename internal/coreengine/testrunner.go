package coreengine

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/cucumber/godog"
)

// TestRunner describes the interface that should be implemented to support the execution of tests.
type TestRunner interface {
	RunTest(t *Test) error
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
}

// GoDogTestTuple holds the tuple of data required when excuting the test case, namely the function to call as
// denoted by Handler and the data to pass to the function, as denoted by Data.
type GoDogTestTuple struct {
	Handler TestHandlerFunc
	Data    *GodogTest
}

var (
	handlers    = make(map[TestDescriptor]*GoDogTestTuple)
	handlersMux sync.RWMutex
)

// TestHandleFunc adds the TestHandlerFunc to the handler map, keyed on the TestDescriptor, and is effectively 
// a register of the test cases.  This is the mechanism which links the test case handler to the TestRunner, 
// therefore it is essential that the test case register itself with the TestRunner by calling this function 
// supplying a description of the test and the GoDogTestTuple.  See pod_security_feature.init() for an example.
func TestHandleFunc(td TestDescriptor, gd *GoDogTestTuple) {
	handlersMux.Lock()
	defer handlersMux.Unlock()

	handlers[td] = gd
}

// RunTest runs the test case described by the supplied Test.  It looks in it's test register (the handlers global
// variable) for an entry with the same TestDescriptor as the supplied test.  If found, it uses the 
// function and data held in the GoDogTestTuple to execute the test: it calls the handler function with the 
// GodogTest data structure.
func (ts *TestStore) RunTest(t *Test) (int, error) {
	ts.AuditLog.Audit(t.UUID, "status", "running")
	
	if t == nil {
		return 2, fmt.Errorf("test is nil - cannot run test")
	}

	if t.TestDescriptor == nil {
		//update status
		*t.Status = Error
		return 3, fmt.Errorf("test descriptor is nil - cannot run test")
	}

	// get the handler (based on the test supplied)
	g, exists := getHandler(t)

	if !exists {
		//update status
		*t.Status = Error
		return 4, fmt.Errorf("no test handler available for %v - cannot run test", *t.TestDescriptor)
	}

	s, o, err := g.Handler(g.Data)

	if s == 0 {
		// success
		*t.Status = CompleteSuccess
	} else {
		// fail
		*t.Status = CompleteFail

		//TODO: this could be adjusted based on test strictness ...
	}

	t.Results = o // If in-mem output provided, store as Results
	return s, err
}

func getHandler(t *Test) (*GoDogTestTuple, bool) {
	handlersMux.Lock()
	defer handlersMux.Unlock()
	g, exists := handlers[*(*t).TestDescriptor]

	return g, exists
}
