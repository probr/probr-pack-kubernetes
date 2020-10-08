package probes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/citihub/probr/internal/clouddriver/kubernetes"
	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/summary"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	apiv1 "k8s.io/api/core/v1"
)

const rootDirName = "probr"

var outputDir *string

// GetRootDir gets the root directory of the probr executable.
func GetRootDir() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	log.Printf("[DEBUG] GetRootDir pwd is: %v", pwd)

	b := strings.Contains(pwd, rootDirName)
	if !b {
		return "", fmt.Errorf("could not find '%v' root directory in %v", rootDirName, pwd)
	}

	s := strings.SplitAfter(pwd, rootDirName)
	log.Printf("[DEBUG] path(s) after splitting: %v\n", s)

	if len(s) < 1 {
		//expect at least one result
		return "", fmt.Errorf("could not split out '%v' from directory in %v", rootDirName, pwd)
	}

	if !strings.HasSuffix(s[0], rootDirName) {
		//the first path should end with "probr"
		return "", fmt.Errorf("first path after split (%v) does not end with '%v'", s[0], rootDirName)
	}

	return s[0], nil
}

// SetOutputDirectory allows specification of the output directory for the test output json files.
func SetOutputDirectory(d *string) {
	outputDir = d
}

// GetOutputPath gets the output path for the test based on the output directory
// plus the test name supplied
func GetOutputPath(t *string) (*os.File, error) {
	outPath, err := getOutputDirectory()
	if err != nil {
		return nil, err
	}
	os.Mkdir(*outPath, os.ModeDir)

	//filename is test name (supplied) + .json
	fn := *t + ".json"
	return os.Create(filepath.Join(*outPath, fn))
}

func getOutputDirectory() (*string, error) {
	if outputDir == nil || len(*outputDir) < 1 {
		log.Printf("[INFO] output directory not set - attempting to default")
		//default it:
		r, err := GetRootDir()
		if err != nil {
			return nil, fmt.Errorf("output directory not set - attempt to default resulted in error: %v", err)
		}

		f := filepath.Join(r, "testoutput")
		outputDir = &f
	}

	log.Printf("[INFO] output directory is: %v", *outputDir)

	return outputDir, nil
}

// LogAndReturnError logs the given string and raise an error with the same string.  This is useful in Godog steps
// where an error is displayed in the test report but not logged.
func LogAndReturnError(e string, v ...interface{}) error {
	var b strings.Builder
	b.WriteString("[ERROR] ")
	b.WriteString(e)

	s := fmt.Sprintf(b.String(), v...)
	log.Print(s)

	return fmt.Errorf(s)
}

// LogScenarioStart logs the name and tags associtated with the supplied scenario.
func LogScenarioStart(s *godog.Scenario) {
	log.Print(scenarioString(true, s))
}

// LogScenarioEnd logs the name and tags associtated with the supplied scenario.
func LogScenarioEnd(s *godog.Scenario) {
	log.Print(scenarioString(false, s))
}

func scenarioString(st bool, s *godog.Scenario) string {
	var b strings.Builder
	if st {
		b.WriteString("[INFO] >>> Scenario Start: ")
	} else {
		b.WriteString("[INFO] <<< Scenario End: ")
	}

	b.WriteString(s.Name)
	b.WriteString(". (Tags: ")

	for _, t := range s.Tags {
		b.WriteString(t.GetName())
		b.WriteString(" ")
	}
	b.WriteString(").")
	return b.String()
}

func notExcluded(tags []*messages.Pickle_PickleTag) bool {
	for _, exclusion := range config.Vars.TagExclusions {
		for _, tag := range tags {
			if tag.Name == "@"+exclusion {
				return false
			}
		}
	}
	return true
}

func BeforeScenario(name string, ps *probeState, s *godog.Scenario) {
	if notExcluded(s.Tags) {
		ps.setup()
		ps.name = s.Name
		ps.event = summary.State.GetEventLog(name)
		LogScenarioStart(s)
	}
}

type probeState struct {
	name             string
	event            *summary.Event
	httpStatusCode   int
	podName          string
	state            State
	useDefaultNS     bool
	hasWildcardRoles bool
}

// Setup resets scenario-specific values
func (p *probeState) setup() {
	p.state.PodName = ""
	p.state.CreationError = nil
	p.useDefaultNS = false
}

// State captures useful state data for use in tests.
type State struct {
	PodName         string
	CreationError   *kubernetes.PodCreationError
	ExpectedReason  *kubernetes.PodCreationErrorReason
	CommandExitCode int
}

// ProcessPodCreationResult is a convenince function to process the result of a pod creation attempt.
// It records state information on the supplied state structure.
func ProcessPodCreationResult(s *State, pd *apiv1.Pod, expected kubernetes.PodCreationErrorReason, e *summary.Event, err error) error {

	//first check for errors:
	if err != nil {
		//check if we've got a partial pod creation
		//e.g. pod was created but didn't get to "running" state
		//in this case we need to hold onto the name so it can be deleted
		if pd != nil {
			s.PodName = pd.GetObjectMeta().GetName()
			e.CountPodCreated()
			summary.State.LogPodName(s.PodName)
		}

		//check for known error type
		//this means the pod has not been created for an expected reason and
		//is a valid result if the test is addressing prevention of insecure pod creation
		if e, ok := err.(*kubernetes.PodCreationError); ok {
			s.CreationError = e
			s.ExpectedReason = &expected
			return nil
		}
		//unexpected error
		//in this case something unexpected has happened, return an error to cucumber
		return LogAndReturnError("error attempting to create POD: %v", err)
	}

	//No errors: pod creation may or may not have been expected.  This will be determined
	//by the specific test case
	if pd == nil {
		// pod not created, which could be valid for some tests
		return nil
	}

	//if we've got this far, a pod was successfully created which could be
	//valid for some tests
	s.PodName = pd.GetObjectMeta().GetName()
	e.CountPodCreated()
	summary.State.LogPodName(s.PodName)

	//we're good
	return nil
}

// AssertResult evaluate the state in the context of the expected condition, e.g. if expected is "fail",
// then the expecation is that a creation error will be present.
func AssertResult(s *State, res, msg string) error {

	if res == "Fail" || res == "denied" {
		//expect pod creation error to be non-null
		if s.CreationError == nil {
			//it's a fail:
			return LogAndReturnError("pod %v was created - test failed", s.PodName)
		}
		//should also check code:
		_, exists := s.CreationError.ReasonCodes[*s.ExpectedReason]
		if !exists {
			//also a fail:
			return LogAndReturnError("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				s.CreationError.ReasonCodes, s.ExpectedReason)
		}

		//we're good
		return nil
	}

	if res == "Succeed" || res == "allowed" {
		// then expect the pod creation error to be nil
		if s.CreationError != nil {
			//it's a fail:
			return LogAndReturnError("pod was not created - test failed: %v", s.CreationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	return LogAndReturnError("desired result %v is not recognised", res)

}

//general feature steps:
func (p *probeState) aKubernetesClusterIsDeployed() error {
	b := kubernetes.GetKubeInstance().ClusterIsDeployed()

	if b == nil || !*b {
		log.Fatalf("[ERROR] Kubernetes cluster is not deployed")
	}
	p.event.LogProbe(p.name, nil) // If not fatal, success
	return nil
}
