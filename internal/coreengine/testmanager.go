package coreengine

import (
	"bytes"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"gitlab.com/citihub/probr/internal/output"
)

//TestStatus ..
type TestStatus int

//TestStatus enum ...
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

//Group ... TODO: not sure if this is the correct name for this?
type Group int

//Group enum
const (
	Kubernetes Group = iota
	CloudDriver
	CoreEngine
)

func (g Group) String() string {
	return [...]string{"kubernetes", "clouddriver", "coreengine"}[g]
}

//Category ...
type Category int

//Category enum
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

//Test - structure to hold test data
type Test struct {
	UUID           string          `json:"uuid,omitempty"`
	TestDescriptor *TestDescriptor `json:"test_descriptor,omitempty"`

	Status *TestStatus `json:"status,omitempty"`

	Results *bytes.Buffer
}

//TestDescriptor - struct to hold description of test, name, category, strictness?? etc.
type TestDescriptor struct {
	Group    Group    `json:"group,omitempty"`
	Category Category `json:"category,omitempty"`
	Name     string   `json:"name,omitempty"`
}

//TestStore - holds the current test suite.
//TODO: still not sure about the structure.
//Below is:
// [uuid] -> pointer to array of pointers to tests
// this implies that we'd be setting a test run up with multiple "sub" test runs,
// with each run being identified by a uuid which is mapped to the array of tests
// I think this could be too complicated, and it's just a simple uuid -> test,
// i.e. the "test run" is a map of multiple entries, each uuid simply pointing to one
// test, i.e.
// [uuid] -> pointer to test
// (done ... simples :-) )
type TestStore struct {
	Tests       map[uuid.UUID]*[]*Test
	FailedTests map[TestStatus]*[]*Test
	Lock        sync.RWMutex
	AuditLog    *output.AuditLog
}

//GetAvailableTests - return the universe of available tests
func GetAvailableTests() *[]TestDescriptor {

	//not sure if this is needed
	//TODO: get this from the TestRunner handler store - basically it's the collection of
	//tests that have registered a handler ..

	// return &p
	return nil
}

//NewTestManager - create a new test manager, backed by TestStore
func NewTestManager() *TestStore {
	return &TestStore{
		Tests:    make(map[uuid.UUID]*[]*Test),
		AuditLog: new(output.AuditLog),
	}
}

//AddTest ...
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

	ts.AuditLog.Audit(u, "status", t.Status.String())
	ts.AuditLog.Audit(u, "descriptor", t.TestDescriptor.Name)

	return &uid
}

//GetTest by UUID ...
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

//ExecTest by UUID in this case, but could be name, category, etc.  Probably need an ExecTests as well ...
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

//ExecAllTests ...
func (ts *TestStore) ExecAllTests() (int, error) {
	status := 0
	var err error

	for uuid := range ts.Tests {
		st, err := ts.ExecTest(&uuid)
		ts.AuditLog.Audit(uuid.String(), "status", "Exited: "+strconv.Itoa(st))
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
