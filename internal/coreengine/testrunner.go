package coreengine

import (	 
	"sync"
	"fmt"
)

//TestRunner ...
type TestRunner interface {
	RunTest(t *Test) error
}

//TestHandlerFunc ...
type TestHandlerFunc func()(int, error)

var (
	handlers = make(map[TestDescriptor]TestHandlerFunc)
	handlersMux sync.RWMutex
)

//TestHandleFunc - adds the TestHandlerFunc to the handler map, keyed on the TestDescriptor
func TestHandleFunc(td TestDescriptor, handler func()(int, error)) {
	handlersMux.Lock()
	defer handlersMux.Unlock()
	
	handlers[td]=handler
}

//RunTest TODO: remove TestStore?
func (ts *TestStore) RunTest(t *Test) (int, error) {

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
	f, exists := getHandler(t)

	if !exists {
		//update status
		*t.Status = Error
		return 4, fmt.Errorf("no test handler available for %v - cannot run test", *t.TestDescriptor)
	}
	
	s, err := f()

	if s == 0 {
		// success
		*t.Status = CompleteSuccess
	} else {
		// fail
		*t.Status = CompleteFail

		//TODO: this could be adjusted based on test strictness ...
	}

	return s, err
}

func getHandler(t *Test) (func()(int, error), bool) {
	handlersMux.Lock()
	defer handlersMux.Unlock()
	f, exists := handlers[*(*t).TestDescriptor]
	
	return f, exists
}