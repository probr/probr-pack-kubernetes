package kubernetes_test

import (
	"testing"

	"citihub.com/probr/internal/clouddriver/kubernetes"
)

func TestClusterHasPSP(t *testing.T) {

	//TODO: THIS IS NOT REALLY A UNIT TEST
	//but if we want to run it as an integration test and
	//have it interact with a cluster we really need to have
	//either example clusters or set them up as part of the
	//test.   This will basically be what we're doing in the
	//feature/bdd tests so that's probably a more relevant place
	//for that.   Here, just do some basic stuff ...
	yesNo, err := kubernetes.ClusterHasPSP()

	handleResult(yesNo, err)

}

func TestPrivilegedAccessIsRestricted(t *testing.T) {
	yesNo, err := kubernetes.PrivilegedAccessIsRestricted()

	handleResult(yesNo, err)
}

func TestHostPIDIsRestricted(t *testing.T) {
	yesNo, err := kubernetes.HostPIDIsRestricted()

	handleResult(yesNo, err)
}
