// Package coreengine contains the types and functions responsible for managing tests and test execution.  This is the primary
// entry point to the core of the application and should be utilised by the probr library to create, execute and report
// on tests.
package coreengine

import (
	"errors"
	"log"
	"sync"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/summary"
)

// TestStatus type describes the status of the test, e.g. Pending, Running, CompleteSuccess, CompleteFail and Error
type TestStatus int

//TestStatus enumeration for the TestStatus type.
const (
	Pending TestStatus = iota
	Running
	CompleteSuccess
	CompleteFail
	Error
	Excluded
)

func (s TestStatus) String() string {
	return [...]string{"Pending", "Running", "CompleteSuccess", "CompleteFail", "Error", "Excluded"}[s]
}

// Group type describes the group to which the test belongs, e.g. kubernetes, clouddriver, coreengine, etc.
type Group int

// Group type enumeration
const (
	Kubernetes Group = iota
	CloudDriver
	CoreEngine
)

func (g Group) String() string {
	return [...]string{"kubernetes", "clouddriver", "coreengine"}[g]
}

// TestDescriptor describes the specific test case and includes name and group.
type TestDescriptor struct {
	Group Group  `json:"group,omitempty"`
	Name  string `json:"name,omitempty"`
}

// TestStore maintains a collection of tests to be run and their status.  FailedTests is an explicit
// collection of failed tests.
type TestStore struct {
	Tests       map[string]*GodogTest
	FailedTests map[TestStatus]*GodogTest
	Lock        sync.RWMutex
}

// GetAvailableTests return the collection of available tests.
func GetAvailableTests() *[]TestDescriptor {
	//TODO: to implement
	//get this from the TestRunner handler store - basically it's the collection of
	//tests that have registered a handler ..

	// return &p
	return nil
}

// NewTestManager creates a new test manager, backed by TestStore
func NewTestManager() *TestStore {
	return &TestStore{
		Tests: make(map[string]*GodogTest),
	}
}

// AddTest provided GodogTest to the TestStore.
func (ts *TestStore) AddTest(test *GodogTest) string {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	var status TestStatus
	if test.TestDescriptor.isExcluded() {
		status = Excluded
	} else {
		status = Pending
	}

	//add the test
	test.Status = &status
	ts.Tests[test.TestDescriptor.Name] = test

	summary.State.GetProbeLog(test.TestDescriptor.Name).Result = test.Status.String()
	summary.State.LogProbeMeta(test.TestDescriptor.Name, "group", test.TestDescriptor.Group.String())

	return test.TestDescriptor.Name
}

// GetTest returns the test identified by the given name.
func (ts *TestStore) GetTest(name string) (*GodogTest, error) {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	//get the test from the store
	t, exists := ts.Tests[name]

	if !exists {
		return nil, errors.New("test with name '" + name + "' not found")
	}
	return t, nil
}

// ExecTest executes the test identified by the specified name.
func (ts *TestStore) ExecTest(name string) (int, error) {
	t, err := ts.GetTest(name)
	if err != nil {
		return 1, err // Failure
	}
	if t.Status.String() != Excluded.String() {
		return ts.RunTest(t) // Return test results
	}
	return 0, nil // Succeed if test is excluded
}

// ExecAllTests executes all tests that are present in the TestStore.
func (ts *TestStore) ExecAllTests() (int, error) {
	status := 0
	var err error

	for name := range ts.Tests {
		st, err := ts.ExecTest(name)
		summary.State.ProbeComplete(name)
		if err != nil {
			//log but continue with remaining tests
			log.Printf("[ERROR] error executing test: %v", err)
		}
		if st > status {
			status = st
		}
	}
	return status, err
}

func (td *TestDescriptor) isExcluded() bool {
	v := []string{td.Name, td.Group.String()} // iterable name & group strings
	for _, r := range v {
		if tagIsExcluded(r) {
			return true
		}
	}
	return false
}

func tagIsExcluded(tag string) bool {
	for _, exclusion := range config.Vars.TagExclusions {
		if tag == exclusion {
			return true
		}
	}
	return false
}
