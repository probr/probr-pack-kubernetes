package kubernetes_test

import (
	"testing"

	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
)

//TODO: this will be revised when the unit/integration tests are refactored to be properly mocked
var psp = kubernetes.NewDefaultPSP()

func TestClusterHasPSP(t *testing.T) {

	//TODO: THIS IS NOT REALLY A UNIT TEST
	//but if we want to run it as an integration test and
	//have it interact with a cluster we really need to have
	//either example clusters or set them up as part of the
	//test.   This will basically be what we're doing in the
	//feature/bdd tests so that's probably a more relevant place
	//for that.   Here, just do some basic stuff ...
	yesNo, err := psp.ClusterHasPSP()

	handleResult(yesNo, err)

}

func TestPrivilegedAccessIsRestricted(t *testing.T) {
	yesNo, err := psp.PrivilegedAccessIsRestricted()

	handleResult(yesNo, err)
}

func TestHostPIDIsRestricted(t *testing.T) {
	yesNo, err := psp.HostPIDIsRestricted()

	handleResult(yesNo, err)
}

func TestCreatePODSettingPrivilegedAccess(t *testing.T) {
	tr := true
	p, err := psp.CreatePODSettingSecurityContext(&tr, &tr, nil)

	//pod creation should fail so p should be nil
	res := p == nil
	handleResult(&res, err)

}

func TestCreatePODSettingCapabilities(t *testing.T) {
	var c = make([]string, 1)
	c[0] = "NET_ADMIN"
	
	p, err := psp.CreatePODSettingCapabilities(&c)

	//pod creation should fail so p should be nil
	res := p == nil
	handleResult(&res, err)

}

func TestPrivilegedEscalationPrevented(t *testing.T) {
	res, err := psp.ExecPSPTestCmd(nil, kubernetes.Chroot)

	//this should fail against a secured cluster
	//non-zero result required
	b := res > 0
	handleResult(&b, err)
}