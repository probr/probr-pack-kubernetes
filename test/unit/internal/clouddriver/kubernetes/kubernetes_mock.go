package kubernetesunit

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	apiv1 "k8s.io/api/core/v1"
	k8s "k8s.io/client-go/kubernetes"
	rbacv1 "k8s.io/api/rbac/v1"
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

func (m *kubeMock) GetClient() (*k8s.Clientset, error) {
	c := m.Called().Get(0).(*k8s.Clientset)
	e := m.Called().Error(1)
	return c, e
}

func (m *kubeMock) GetPods(ns string) (*apiv1.PodList, error) {
	a := m.Called()
	pl := a.Get(0).(*apiv1.PodList)
	e := a.Error(1)
	return pl, e
}
func (m *kubeMock) CreatePod(pname *string, ns *string, cname *string, image *string, w bool, sc *apiv1.SecurityContext) (*apiv1.Pod, error) {
	//The below will check the args are as expected, ie. the security context has the correct attributes
	a := m.Called(pname, ns, cname, image, w, sc)

	return a.Get(0).(*apiv1.Pod), a.Error(1)
}
func (m *kubeMock) CreatePodFromObject(p *apiv1.Pod, pname *string, ns *string, w bool) (*apiv1.Pod, error) {
	//The below will check the args are as expected, ie. the Pod has the correct attributes
	a := m.Called(p, pname, ns, w)

	//This time, return the pod we've been given, so ignore what's been supplied on the mock call:
	return p, a.Error(1)
}
func (m *kubeMock) CreatePodFromYaml(y []byte, pname *string, ns *string, image *string, w bool) (*apiv1.Pod, error) {
	po := m.Called().Get(0).(*apiv1.Pod)
	e := m.Called().Error(1)
	return po, e
}
func (m *kubeMock) GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod {
	p := m.Called().Get(0).(*apiv1.Pod)
	return p
}
func (m *kubeMock) ExecCommand(cmd, ns, pn *string) *kubernetes.CmdExecutionResult {
	a := m.Called()
	return a.Get(0).(*kubernetes.CmdExecutionResult)
}
func (m *kubeMock) DeletePod(pname *string, ns *string, w bool) error {
	e := m.Called().Error(0)
	return e
}
func (m *kubeMock) DeleteNamespace(ns *string) error {
	e := m.Called().Error(0)
	return e
}
func (m *kubeMock) CreateConfigMap(n *string, ns *string) (*apiv1.ConfigMap, error) {
	cm := m.Called().Get(0).(*apiv1.ConfigMap)
	e := m.Called().Error(1)
	return cm, e
}
func (m *kubeMock) DeleteConfigMap(n *string, ns *string) error {
	e := m.Called().Error(0)
	return e
}
func (m *kubeMock) GetConstraintTemplates(prefix *string) (*map[string]interface{}, error) {
	a := m.Called()
	return a.Get(0).(*map[string]interface{}), a.Error(1)
}
func (m *kubeMock) GetClusterRolesByResource(r string) (*[]rbacv1.ClusterRole, error) {
	a := m.Called()
	return a.Get(0).(*[]rbacv1.ClusterRole), a.Error(1)
}
func (m *kubeMock) GetClusterRoles() (*rbacv1.ClusterRoleList, error){
	a := m.Called()
	return a.Get(0).(*rbacv1.ClusterRoleList), a.Error(1)
}
