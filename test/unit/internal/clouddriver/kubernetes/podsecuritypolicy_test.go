package kubernetesunit

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/config"
	"gitlab.com/citihub/probr/internal/utils"
	apiv1 "k8s.io/api/core/v1"
)

func TestPSPTestCommand(t *testing.T) {
	assert.Equal(t, "chroot .", kubernetes.Chroot.String())
	assert.Equal(t, "nsenter -t 1 -p ps", kubernetes.EnterHostPIDNS.String())
	assert.Equal(t, "nsenter -t 1 -i ps", kubernetes.EnterHostIPCNS.String())
	assert.Equal(t, "nsenter -t 1 -n ps", kubernetes.EnterHostNetworkNS.String())
	assert.Equal(t, "id -u > 0 ", kubernetes.VerifyNonRootUID.String())
	assert.Equal(t, "ping google.com", kubernetes.NetRawTest.String())
	assert.Equal(t, "ip link add dummy0 type dummy", kubernetes.SpecialCapTest.String())

}

//TODO: need to add tests in for all "Has..." functions:
//TODO: more tests for all interface methods

func TestClusterHasPSP(t *testing.T) {
	runTest(t, "HasSecurityPolicies", "ClusterHasPSP")
}

func TestPrivilegedAccessIsRestricted(t *testing.T) {
	runTest(t, "HasPrivilegedAccessRestriction", "PrivilegedAccessIsRestricted")	
}

func TestHostPIDIsRestricted(t *testing.T) {	
	runTest(t, "HasHostPIDRestriction", "HostPIDIsRestricted")
}

func TestHostIPCIsRestricted(t *testing.T) {	
	runTest(t, "HasHostIPCRestriction", "HostIPCIsRestricted")
}

func TestHostNetworkIsRestricted(t *testing.T) {	
	runTest(t, "HasHostNetworkRestriction", "HostNetworkIsRestricted")
}

func TestPrivilegedEscalationIsRestricted(t *testing.T) {	
	runTest(t, "HasAllowPrivilegeEscalationRestriction", "PrivilegedEscalationIsRestricted")
}

func TestRootUserIsRestricted(t *testing.T) {	
	runTest(t, "HasRootUserRestriction", "RootUserIsRestricted")
}

func TestNETRawIsRestricted(t *testing.T) {	
	runTest(t, "HasNETRAWRestriction", "NETRawIsRestricted")
}

func TestAllowedCapabilitiesAreRestricted(t *testing.T) {	
	runTest(t, "HasAllowedCapabilitiesRestriction", "AllowedCapabilitiesAreRestricted")
}

func TestAssignedCapabilitiesAreRestricted(t *testing.T) {	
	runTest(t, "HasAssignedCapabilitiesRestriction", "AssignedCapabilitiesAreRestricted")
}

func TestHostPortsAreRestricted(t *testing.T) {	
	runTest(t, "HasHostPortRestriction", "HostPortsAreRestricted")
}

func TestVolumeTypesAreRestricted(t *testing.T) {	
	runTest(t, "HasVolumeTypeRestriction", "VolumeTypesAreRestricted")
}

func TestSeccompProfilesAreRestricted(t *testing.T) {	
	runTest(t, "HasSeccompProfileRestriction", "SeccompProfilesAreRestricted")
}

func runTest(t *testing.T, mockMethod string, testMethod string) {
	fmt.Printf("==== Running PSP test for %v with mock %v \n", testMethod, mockMethod)

	//create and set a 'true' provider
	tmp := &securityProviderMock{}
	tmp.On(mockMethod).Return(true, nil)

	psp := kubernetes.NewPSP(nil, &[]kubernetes.SecurityPolicyProvider{tmp}, config.GetEnvConfigInstance())

	b, _ := reflectiveCall(psp, testMethod)
	assert.True(t, *b)
	tmp.AssertNumberOfCalls(t, mockMethod, 1)
	tmp.AssertExpectations(t)

	//create and set a 'false' provider
	fmp := &securityProviderMock{}
	fmp.On(mockMethod).Return(false, nil)

	//need a different PSP with the false provider
	psp = kubernetes.NewPSP(nil, &[]kubernetes.SecurityPolicyProvider{fmp}, config.GetEnvConfigInstance())

	b, _ = reflectiveCall(psp, testMethod)
	assert.False(t, *b)
	fmp.AssertNumberOfCalls(t, mockMethod, 1)
	fmp.AssertExpectations(t)

	//add both true & false providers:
	//true first ...
	psp = kubernetes.NewPSP(nil, &[]kubernetes.SecurityPolicyProvider{tmp, fmp}, config.GetEnvConfigInstance())

	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)               //expect this to be true ...
	tmp.AssertNumberOfCalls(t, mockMethod, 2) //expect another call to true (so two in total)
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 1) //but false should only be at 1 as PSP should return on first +ve

	//switch order of true/false (so false first)
	psp = kubernetes.NewPSP(nil, &[]kubernetes.SecurityPolicyProvider{fmp, tmp}, config.GetEnvConfigInstance())
	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)               //expect this to be true ...
	tmp.AssertNumberOfCalls(t, mockMethod, 3) //expect another call to true (three in total now)
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //false should be up to 2 by now ...

	//add nil provider in the first slot ...
	psp = kubernetes.NewPSP(nil, &[]kubernetes.SecurityPolicyProvider{nil, tmp}, config.GetEnvConfigInstance())
	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)               //should still get an overall true result...
	tmp.AssertNumberOfCalls(t, mockMethod, 4) //true up to four
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //no call to false

}

func reflectiveCall(p *kubernetes.PSP, tm string) (*bool, error) {

	fmt.Printf("Reflectively calling %v on %T\n", tm, p)

	res := reflect.ValueOf(p).MethodByName(tm).Call([]reflect.Value{})	

	fmt.Printf("Reflection Result: %v\n", res)

	//expect two params in the result:
	if len(res) != 2 {
		panic("unexpected number of values in function return") 
	}

	b := res[0].Elem()
	e := res[1].Interface()

	fmt.Printf("Extracted Results: %v, %v\n", b, e)

	if bb, ok := b.Interface().(bool); ok {		
		return &bb, nil
	}
	if err, ok := e.(error); ok {
		return nil, err
	}

	//unexpected:
	return nil, nil
}

//TODO: this feel rough - we need access to some 'kube' functions, but want to short circuit any external calls
//should move these more general functions out to a utility/helper
var k = kubernetes.GetKubeInstance()

func TestCreatePODSettingSecurityContext(t *testing.T) {
	//need a mock kube
	mk := &kubeMock{}
	psp := kubernetes.NewPSP(mk, nil, config.GetEnvConfigInstance())

	//set up the mock
	sc := apiv1.SecurityContext{
		Privileged:               utils.BoolPtr(true),
		AllowPrivilegeEscalation: utils.BoolPtr(true),
		RunAsUser:                utils.Int64Ptr(2000),
	}

	mk.On("CreatePod", mock.Anything, mock.Anything, mock.Anything, mock.Anything, true, &sc).
		Return(k.GetPodObject("n", "ns", "c", "i", &sc), nil).Once()

	//privileged and privileged access true, runasuser 2000
	p, err := psp.CreatePODSettingSecurityContext(utils.BoolPtr(true), utils.BoolPtr(true), utils.Int64Ptr(2000))

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//with a security context (container):
	assert.NotNil(t, p.Spec.Containers[0].SecurityContext) //should only be one
	//and with privileged, allowprivileged = true and runasuser = 2000
	assert.Equal(t, utils.BoolPtr(true), p.Spec.Containers[0].SecurityContext.Privileged)
	assert.Equal(t, utils.BoolPtr(true), p.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, utils.Int64Ptr(2000), p.Spec.Containers[0].SecurityContext.RunAsUser)
	mk.AssertNumberOfCalls(t, "CreatePod", 1)
	mk.AssertExpectations(t)

	//privileged false, privileged access true, runasuser nil
	sc = apiv1.SecurityContext{
		Privileged:               utils.BoolPtr(false),
		AllowPrivilegeEscalation: utils.BoolPtr(true),
		RunAsUser:                utils.Int64Ptr(1000),
	}
	mk.On("CreatePod", mock.Anything, mock.Anything, mock.Anything, mock.Anything, true, &sc).
		Return(k.GetPodObject("n", "ns", "c", "i", &sc), nil).Once()

	p, err = psp.CreatePODSettingSecurityContext(utils.BoolPtr(false), utils.BoolPtr(true), nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//with a security context (container):
	assert.NotNil(t, p.Spec.Containers[0].SecurityContext) //should only be one
	//and with privileged, allowprivileged = true and runasuser = 1000 (default)
	assert.Equal(t, utils.BoolPtr(false), p.Spec.Containers[0].SecurityContext.Privileged)
	assert.Equal(t, utils.BoolPtr(true), p.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, utils.Int64Ptr(1000), p.Spec.Containers[0].SecurityContext.RunAsUser)
	mk.AssertNumberOfCalls(t, "CreatePod", 2) //2 calls now.  TODO: need to figure out how to reset this!
	mk.AssertExpectations(t)

}

type securityProviderMock struct {
	mock.Mock
}

func (m *securityProviderMock) HasSecurityPolicies() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasPrivilegedAccessRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasHostPIDRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasHostIPCRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasHostNetworkRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasRootUserRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasNETRAWRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasAllowedCapabilitiesRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasAssignedCapabilitiesRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasHostPortRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasVolumeTypeRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
func (m *securityProviderMock) HasSeccompProfileRestriction() (*bool, error) {
	b := m.Called().Bool(0)
	return &b, nil
}
