// Package coreengine contains the types and functions responsible for managing tests and test execution.  This is the primary
// entry point to the core of the application and should be utilised by the probr library to create, execute and report
// on tests.
package coreengine

import (
	"bytes"
	"errors"
	"log"
	"sync"

	"gitlab.com/citihub/probr/internal/audit"
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
)

func (s TestStatus) String() string {
	return [...]string{"Pending", "Running", "CompleteSuccess", "CompleteFail", "Error"}[s]
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

// Category type describes the category to with the test belongs, e.g. PodSecurity, Network, etc.
type Category int

// Category type enumeration.
const (
	RBAC Category = iota
	PodSecurityPolicies
	NetworkPolicies
	SecretsMgmt
	General
	ContainerRegistryAccess
	ImageScanning
	IAM
	KeyMgmt
	Authentication
	Storage
	InternetAccess
)

func (c Category) String() string {
	return [...]string{"RBAC", "Pod Security Policy", "Network Policies", "Secrets Mgmt", "General", "Container Registry Access", "Image Scanning", "IAM",
		"Key Mgmt", "Authentication", "Storage", "InternetAccess"}[c]
}

// Test encapsulates the data required to support test execution.
type Test struct {
	TestDescriptor *TestDescriptor `json:"test_descriptor,omitempty"`

	Status *TestStatus `json:"status,omitempty"`

	Results *bytes.Buffer
}

// TestDescriptor describes the specific test case and includes name, category and group.
type TestDescriptor struct {
	Group    Group    `json:"group,omitempty"`
	Category Category `json:"category,omitempty"`
	Name     string   `json:"name,omitempty"`
}

// TestStore maintains a collection of tests to be run and their status.  FailedTests is an explicit
// collection of failed tests.
type TestStore struct {
	Tests       map[string]*Test
	FailedTests map[TestStatus]*Test
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
		Tests: make(map[string]*Test),
	}
}

// AddTest adds a test, described by the TestDescriptor, to the TestStore.
func (ts *TestStore) AddTest(td TestDescriptor) string {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	//add the test

	s := Pending
	t := Test{
		TestDescriptor: &td,
		Status:         &s,
	}
	ts.Tests[td.Name] = &t

	audit.AuditLog.AuditMeta(td.Name, "status", t.Status.String())
	audit.AuditLog.AuditMeta(td.Name, "group", td.Group.String())
	audit.AuditLog.AuditMeta(td.Name, "category", td.Category.String())

	return td.Name
}

// GetTest returns the test identified by the given name.
func (ts *TestStore) GetTest(name string) (*Test, error) {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	//get the test from the store
	t, exists := ts.Tests[name]

	if !exists {
		return nil, errors.New("test with name '" + name + "' not found")
	}
	return t, nil
}

//GetTest by TestDescriptor ... TODO

// ExecTest executes the test identified by the specified name.
func (ts *TestStore) ExecTest(name string) (int, error) {
	t, err := ts.GetTest(name)

	if err != nil {
		return 1, err
	}

	st, err := ts.RunTest(t)

	//TODO: manage store
	//move to FAILURE / SUCCESS as approriate ...

	return st, err
}

//ExecTest by TestDescriptor, etc ... TODO.  In this case there may be more than one so we should set up for concurrency

// ExecAllTests executes all tests that are present in the TestStore.
func (ts *TestStore) ExecAllTests() (int, error) {
	status := 0
	var err error

	for name := range ts.Tests {
		st, err := ts.ExecTest(name)
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
