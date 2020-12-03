package pod_security_policy

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/kubernetes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
)

func TestPSPProbeCommand(t *testing.T) {
	assert.Equal(t, "chroot .", Chroot.String())
	assert.Equal(t, "nsenter -t 1 -p ps", EnterHostPIDNS.String())
	assert.Equal(t, "nsenter -t 1 -i ps", EnterHostIPCNS.String())
	assert.Equal(t, "nsenter -t 1 -n ps", EnterHostNetworkNS.String())
	assert.Equal(t, "id -u > 0 ", VerifyNonRootUID.String())
	assert.Equal(t, "ping google.com", NetRawProbe.String())
	assert.Equal(t, "ip link add dummy0 type dummy", SpecialCapProbe.String())

}

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

	psp := NewPSP(nil, &[]SecurityPolicyProvider{tmp})

	b, _ := reflectiveCall(psp, testMethod)
	assert.True(t, *b)
	tmp.AssertNumberOfCalls(t, mockMethod, 1)
	tmp.AssertExpectations(t)

	//create and set a 'false' provider
	fmp := &securityProviderMock{}
	fmp.On(mockMethod).Return(false, nil)

	//need a different PSP with the false provider
	psp = NewPSP(nil, &[]SecurityPolicyProvider{fmp})

	b, _ = reflectiveCall(psp, testMethod)
	assert.False(t, *b)
	fmp.AssertNumberOfCalls(t, mockMethod, 1)
	fmp.AssertExpectations(t)

	//add both true & false providers:
	//true first ...
	psp = NewPSP(nil, &[]SecurityPolicyProvider{tmp, fmp})

	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)                        //expect this to be true ...
	tmp.AssertNumberOfCalls(t, mockMethod, 2) //expect another call to true (so two in total)
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 1) //but false should only be at 1 as PSP should return on first +ve
	fmp.AssertExpectations(t)

	//switch order of true/false (so false first)
	psp = NewPSP(nil, &[]SecurityPolicyProvider{fmp, tmp})
	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)                        //expect this to be true ...
	tmp.AssertNumberOfCalls(t, mockMethod, 3) //expect another call to true (three in total now)
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //false should be up to 2 by now ...
	fmp.AssertExpectations(t)

	//add nil provider in the first slot ...
	psp = NewPSP(nil, &[]SecurityPolicyProvider{nil, tmp})
	b, _ = reflectiveCall(psp, testMethod)
	assert.True(t, *b)                        //should still get an overall true result...
	tmp.AssertNumberOfCalls(t, mockMethod, 4) //true up to four
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //no call to false
	fmp.AssertExpectations(t)

	//try with just an "error" provider:
	emp := &securityProviderMock{}
	emp.On(mockMethod).Return(false, fmt.Errorf("SPM ERROR"))

	psp = NewPSP(nil, &[]SecurityPolicyProvider{emp})
	b, err := reflectiveCall(psp, testMethod)
	assert.False(t, *b)                       //overall should be false...
	assert.NotNil(t, err)                     // and we should have an error in this case
	tmp.AssertNumberOfCalls(t, mockMethod, 4) //no call to true (so same as above)
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //no call to false (so same as above)
	fmp.AssertExpectations(t)
	emp.AssertNumberOfCalls(t, mockMethod, 1) //one call to "error provider"
	emp.AssertExpectations(t)

	//now, add a true provider to the second slot
	psp = NewPSP(nil, &[]SecurityPolicyProvider{emp, tmp})
	b, err = reflectiveCall(psp, testMethod)
	assert.True(t, *b)                        //in this case, should still get an overall true result...
	assert.Nil(t, err)                        // and error should be nil (as we've had at least one success)
	tmp.AssertNumberOfCalls(t, mockMethod, 5) //true up to five
	tmp.AssertExpectations(t)
	fmp.AssertNumberOfCalls(t, mockMethod, 2) //no call to false (so as above)
	fmp.AssertExpectations(t)
	emp.AssertNumberOfCalls(t, mockMethod, 2) //another call to "error provider" so up to two
	emp.AssertExpectations(t)

}

func reflectiveCall(p *PSP, tm string) (*bool, error) {

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

	var ret *bool
	var err error
	if bb, ok := b.Interface().(bool); ok {
		ret = &bb
	}
	if er, ok := e.(error); ok {
		err = er
	}

	return ret, err
}

//TODO: this feel rough - we need access to some 'kube' functions, but want to short circuit any external calls
//should move these more general functions out to a utility/helper
var k = kubernetes.GetKubeInstance()

func TestCreatePODSettingSecurityContext(t *testing.T) {
	//need a mock kube
	mk := &kubernetes.KubeMock{}
	psp := NewPSP(mk, nil)

	//set up the mock
	sc := apiv1.SecurityContext{
		Privileged:               utils.BoolPtr(true),
		AllowPrivilegeEscalation: utils.BoolPtr(true),
		RunAsUser:                utils.Int64Ptr(2000),
	}

	mk.On("CreatePod", mock.Anything, mock.Anything, mock.Anything, mock.Anything, true, &sc, mock.Anything, mock.Anything).
		Return(k.GetPodObject("n", "ns", "c", "i", &sc), nil).Once()

	//privileged and privileged access true, runasuser 2000
	p, err := psp.CreatePODSettingSecurityContext(utils.BoolPtr(true), utils.BoolPtr(true), utils.Int64Ptr(2000), nil)

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
	mk.On("CreatePod", mock.Anything, mock.Anything, mock.Anything, mock.Anything, true, &sc, mock.Anything, mock.Anything).
		Return(k.GetPodObject("n", "ns", "c", "i", &sc), nil).Once()

	p, err = psp.CreatePODSettingSecurityContext(utils.BoolPtr(false), utils.BoolPtr(true), nil, nil)

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

func TestCreatePODSettingAttributes(t *testing.T) {
	//need a mock kube
	mk := &kubernetes.KubeMock{}
	psp := NewPSP(mk, nil)

	//set up the mock
	po := k.GetPodObject("n", "ns", "c", "i", nil)
	mk.On("GetPodObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(po, nil)
	mk.On("CreatePodFromObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(po, nil)

	//hostPID, hostIPC & hostNetwork all true:
	p, err := psp.CreatePODSettingAttributes(utils.BoolPtr(true), utils.BoolPtr(true), utils.BoolPtr(true), nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//check hostPID, hostIPC & hostNetwork values:
	assert.Equal(t, true, p.Spec.HostPID)
	assert.Equal(t, true, p.Spec.HostIPC)
	assert.Equal(t, true, p.Spec.HostNetwork)
	mk.AssertNumberOfCalls(t, "GetPodObject", 1)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 1)
	mk.AssertExpectations(t)

	//hostPID, hostIPC & hostNetwork all false:
	p, err = psp.CreatePODSettingAttributes(utils.BoolPtr(false), utils.BoolPtr(false), utils.BoolPtr(false), nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//check hostPID, hostIPC & hostNetwork values:
	assert.Equal(t, false, p.Spec.HostPID)
	assert.Equal(t, false, p.Spec.HostIPC)
	assert.Equal(t, false, p.Spec.HostNetwork)
	mk.AssertNumberOfCalls(t, "GetPodObject", 2)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 2)
	mk.AssertExpectations(t)

	//hostPID, hostIPC & hostNetwork mixed:
	p, err = psp.CreatePODSettingAttributes(utils.BoolPtr(false), utils.BoolPtr(true), nil, nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//check hostPID, hostIPC & hostNetwork values:
	assert.Equal(t, false, p.Spec.HostPID)
	assert.Equal(t, true, p.Spec.HostIPC)
	assert.Equal(t, false, p.Spec.HostNetwork)
	mk.AssertNumberOfCalls(t, "GetPodObject", 3)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 3)
	mk.AssertExpectations(t)

	//hostPID, hostIPC & hostNetwork all nil (should default to false):
	p, err = psp.CreatePODSettingAttributes(nil, nil, nil, nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//check hostPID, hostIPC & hostNetwork values:
	assert.Equal(t, false, p.Spec.HostPID)
	assert.Equal(t, false, p.Spec.HostIPC)
	assert.Equal(t, false, p.Spec.HostNetwork)
	mk.AssertNumberOfCalls(t, "GetPodObject", 4)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 4)
	mk.AssertExpectations(t)

}

func TestCreatePODSettingCapabilities(t *testing.T) {
	//need a mock kube
	mk := &kubernetes.KubeMock{}
	psp := NewPSP(mk, nil)

	//set up the mock
	po := k.GetPodObject("n", "ns", "c", "i", nil)
	mk.On("GetPodObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(po, nil)
	mk.On("CreatePodFromObject", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(po, nil)

	//no capabilities:
	p, err := psp.CreatePODSettingCapabilities(nil, nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//only expect one container
	assert.Equal(t, 1, len(p.Spec.Containers))
	//don't expect any capabilities, and container sec context should be non-nil though
	assert.NotNil(t, p.Spec.Containers[0].SecurityContext)
	assert.Nil(t, p.Spec.Containers[0].SecurityContext.Capabilities)
	mk.AssertNumberOfCalls(t, "GetPodObject", 1)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 1)
	mk.AssertExpectations(t)

	//some capabilities:
	c := []string{"NET_RAW"}
	p, err = psp.CreatePODSettingCapabilities(&c, nil)

	//don't expect an error
	assert.Nil(t, err)
	//do expect pod
	assert.NotNil(t, p)
	//only expect one container
	assert.Equal(t, 1, len(p.Spec.Containers))
	//don't expect any capabilities, and container sec context should be non-nil though
	assert.NotNil(t, p.Spec.Containers[0].SecurityContext)
	assert.NotNil(t, p.Spec.Containers[0].SecurityContext.Capabilities)
	//expect one capability
	assert.Equal(t, 1, len(p.Spec.Containers[0].SecurityContext.Capabilities.Add))
	//should be "NET_RAW"
	assert.Equal(t, "NET_RAW", string(p.Spec.Containers[0].SecurityContext.Capabilities.Add[0]))

	mk.AssertNumberOfCalls(t, "GetPodObject", 2)
	mk.AssertNumberOfCalls(t, "CreatePodFromObject", 2)
	mk.AssertExpectations(t)

}

type securityProviderMock struct {
	mock.Mock
}

func (m *securityProviderMock) HasSecurityPolicies() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasPrivilegedAccessRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasHostPIDRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasHostIPCRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasHostNetworkRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasRootUserRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasNETRAWRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasAllowedCapabilitiesRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasAssignedCapabilitiesRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasHostPortRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasVolumeTypeRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) HasSeccompProfileRestriction() (*bool, error) {
	return m.returnArgs(m.Called())
}
func (m *securityProviderMock) returnArgs(a mock.Arguments) (*bool, error) {
	b := a.Bool(0)
	e := a.Error(1)
	return &b, e
}
