package probe

import (
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/test/features"
	apiv1 "k8s.io/api/core/v1"
)

// State captures useful state data for use in tests.
type State struct {
	PodName        string
	CreationError  *kubernetes.PodCreationError
	ExpectedReason *kubernetes.PodCreationErrorReason
}

// ProcessPodCreationResult is a convenince function to process the result of a pod creation attempt. 
// It records state information on the supplied state structure.
func ProcessPodCreationResult(s *State, pd *apiv1.Pod, expected kubernetes.PodCreationErrorReason, err error) error {

	//first check for errors:
	if err != nil {
		//check if we've got a partial pod creation
		//e.g. pod was created but did't get to "running" state
		//in this case we need to hold onto the name so it can be deleted
		if pd != nil {
			s.PodName = pd.GetObjectMeta().GetName()
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
		return features.LogAndReturnError("error attempting to create POD: %v", err)
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
			return features.LogAndReturnError("pod %v was created - test failed", s.PodName)
		}
		//should also check code:
		_, exists := s.CreationError.ReasonCodes[*s.ExpectedReason]
		if !exists {
			//also a fail:
			return features.LogAndReturnError("pod not was created but failure reasons (%v) did not contain expected (%v)- test failed",
				s.CreationError.ReasonCodes, s.ExpectedReason)
		}

		//we're good
		return nil
	}

	if res == "Succeed" || res == "allowed" {
		// then expect the pod creation error to be nil
		if s.CreationError != nil {
			//it's a fail:
			return features.LogAndReturnError("pod was not created - test failed: %v", s.CreationError)
		}

		//else we're good ...
		return nil
	}

	// we've been given a result that we don't know about ...
	return features.LogAndReturnError("desired result %v is not recognised", res)

}
