// Package coreengine contains the types and functions responsible for managing tests and test execution.  This is the primary
// entry point to the core of the application and should be utilised by the probr library to create, execute and report
// on tests.
package coreengine

import (
	"bytes"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
	"gitlab.com/citihub/probr/internal/output"
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
	UUID           string          `json:"uuid,omitempty"`
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
	Tests       map[uuid.UUID]*[]*Test
	FailedTests map[TestStatus]*[]*Test
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
		Tests: make(map[uuid.UUID]*[]*Test),
	}
}

// AddTest adds a test, described by the TestDescriptor, to the TestStore.
func (ts *TestStore) AddTest(td TestDescriptor) *uuid.UUID {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	//add the test

	uid := uuid.New()
	u := uid.String()
	s := Pending
	t := Test{
		TestDescriptor: &td,
		Status:         &s,
		UUID:           u,
	}
	a := []*Test{&t}
	ts.Tests[uid] = &a

	output.AuditLog.Audit(u, "status", t.Status.String())
	output.AuditLog.Audit(u, "descriptor", t.TestDescriptor.Name)

	return &uid
}

// GetTest returns the test identified by the given UUID.
func (ts *TestStore) GetTest(uuid *uuid.UUID) (*[]*Test, error) {
	ts.Lock.Lock()
	defer ts.Lock.Unlock()

	//get the test from the store
	t, exists := ts.Tests[*uuid]

	if !exists {
		return nil, errors.New("test with uuid " + (*uuid).String() + " not found")
	}
	return t, nil
}

//GetTest by TestDescriptor ... TODO

// ExecTest executes the test identified by the supplied UUID.
func (ts *TestStore) ExecTest(uuid *uuid.UUID) (int, error) {
	t, err := ts.GetTest(uuid)

	if err != nil {
		return 1, err
	}

	st, err := ts.RunTest((*t)[0])

	//TODO: manage store
	//move to FAILURE / SUCCESS as approriate ...

	return st, err
}

//ExecTest by TestDescriptor, etc ... TODO.  In this case there may be more than one so we should set up for concurrency

// ExecAllTests executes all tests that are present in the TestStore.
func (ts *TestStore) ExecAllTests() (int, error) {
	status := 0
	var err error

	for uuid := range ts.Tests {
		st, err := ts.ExecTest(&uuid)
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
