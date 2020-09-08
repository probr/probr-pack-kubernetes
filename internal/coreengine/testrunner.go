package coreengine

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/cucumber/godog"
)

//TestRunner ...
type TestRunner interface {
	RunTest(t *Test) error
}

//TestHandlerFunc ...
type TestHandlerFunc func(t *GodogTest) (int, *bytes.Buffer, error)

//GodogTest ...
type GodogTest struct {
	TestDescriptor       *TestDescriptor
	TestSuiteInitializer func(*godog.TestSuiteContext)
	ScenarioInitializer  func(*godog.ScenarioContext)
	FeaturePath          *string
}

//GoDogTestTuple ...
type GoDogTestTuple struct {
	Handler TestHandlerFunc
	Data    *GodogTest
}

var (
	handlers    = make(map[TestDescriptor]*GoDogTestTuple)
	handlersMux sync.RWMutex
)

//TestHandleFunc - adds the TestHandlerFunc to the handler map, keyed on the TestDescriptor
func TestHandleFunc(td TestDescriptor, gd *GoDogTestTuple) {
	handlersMux.Lock()
	defer handlersMux.Unlock()

	handlers[td] = gd
}

//RunTest TODO: remove TestStore?
func (ts *TestStore) RunTest(t *Test) (int, error) {
	ts.AuditLog.Audit(t.UUID, "status", "running")

	//TODO: error codes!
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
