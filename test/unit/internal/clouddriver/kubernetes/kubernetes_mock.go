package kubernetesunit

import (
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type kubeMock struct {
	mock.Mock
}

func (m *kubeMock) ClusterIsDeployed() *bool {
	b := m.Called().Bool(0)
	return &b
}

func (m *kubeMock) SetKubeConfigFile(f *string) {
	m.Called()
}

func (m *kubeMock) GetClient() (*kubernetes.Clientset, error) {
	c := m.Called().Get(0).(*kubernetes.Clientset)
	e := m.Called().Error(1)
	return c, e
}

func (m *kubeMock) GetPods() (*apiv1.PodList, error) {
	pl := m.Called().Get(0).(*apiv1.PodList)
	e := m.Called().Error(1)
	return pl, e
}
func (m *kubeMock) CreatePod(pname *string, ns *string, cname *string, image *string, w bool, sc *apiv1.SecurityContext) (*apiv1.Pod, error) {
	//The below will check the args are as expected, ie. the security context has the correct attributes
	a := m.Called(pname, ns, cname, image, w, sc)
		
	return a.Get(0).(*apiv1.Pod), a.Error(1)
}
func (m *kubeMock) CreatePodFromObject(p *apiv1.Pod, pname *string, ns *string, w bool) (*apiv1.Pod, error) {
	po := m.Called().Get(0).(*apiv1.Pod)
	e := m.Called().Error(1)
	return po, e
}
func (m *kubeMock) GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod {
	p := m.Called().Get(0).(*apiv1.Pod)
	return p
}
func (m *kubeMock) ExecCommand(cmd, ns, pn *string) (string, string, int, error) {
	so := m.Called().String(0)
	se := m.Called().String(1)
	ec := m.Called().Int(2)
	e := m.Called().Error(3)
	return so, se, ec, e
}
func (m *kubeMock) DeletePod(pname *string, ns *string, w bool) error {
	e := m.Called().Error(0)
	return e
}
func (m *kubeMock) DeleteNamespace(ns *string) error {
	e := m.Called().Error(0)
	return e
}
