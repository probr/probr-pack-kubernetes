package kubernetes

import (
	"github.com/citihub/probr/audit"
	"github.com/stretchr/testify/mock"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8s "k8s.io/client-go/kubernetes"
)

type KubeMock struct {
	mock.Mock
}

func (m *KubeMock) ClusterIsDeployed() *bool {
	b := m.Called().Bool(0)
	return &b
}

func (m *KubeMock) SetKubeConfigFile(f *string) {
	m.Called()
}

func (m *KubeMock) GetClient() (*k8s.Clientset, error) {
	c := m.Called().Get(0).(*k8s.Clientset)
	e := m.Called().Error(1)
	return c, e
}

func (m *KubeMock) GetPods(ns string) (*apiv1.PodList, error) {
	a := m.Called()
	pl := a.Get(0).(*apiv1.PodList)
	e := a.Error(1)
	return pl, e
}
func (m *KubeMock) CreatePod(pname string, ns string, cname string, image string, w bool, sc *apiv1.SecurityContext, probe *audit.Probe) (*apiv1.Pod, *PodAudit, error) {
	//The below will check the args are as expected, ie. the security context has the correct attributes
	a := m.Called(pname, ns, cname, image, w, sc)

	return a.Get(0).(*apiv1.Pod), &PodAudit{}, a.Error(1)
}
func (m *KubeMock) CreatePodFromObject(p *apiv1.Pod, pname string, ns string, w bool, probe *audit.Probe) (*apiv1.Pod, error) {
	//The below will check the args are as expected, ie. the Pod has the correct attributes
	a := m.Called(p, pname, ns, w)

	//This time, return the pod we've been given, so ignore what's been supplied on the mock call:
	return p, a.Error(1)
}
func (m *KubeMock) CreatePodFromYaml(y []byte, pname string, ns string, image string, identityBinding string, w bool, probe *audit.Probe) (*apiv1.Pod, error) {
	po := m.Called().Get(0).(*apiv1.Pod)
	e := m.Called().Error(1)
	return po, e
}
func (m *KubeMock) GetPodObject(pname string, ns string, cname string, image string, sc *apiv1.SecurityContext) *apiv1.Pod {
	p := m.Called().Get(0).(*apiv1.Pod)
	return p
}
func (m *KubeMock) ExecCommand(cmd string, ns string, pn *string) *CmdExecutionResult {
	a := m.Called()
	return a.Get(0).(*CmdExecutionResult)
}
func (m *KubeMock) DeletePod(pname string, ns string, e string) error {
	x := m.Called().Error(0)
	return x
}
func (m *KubeMock) DeleteNamespace(ns *string) error {
	e := m.Called().Error(0)
	return e
}
func (m *KubeMock) CreateConfigMap(n *string, ns string) (*apiv1.ConfigMap, error) {
	cm := m.Called().Get(0).(*apiv1.ConfigMap)
	e := m.Called().Error(1)
	return cm, e
}
func (m *KubeMock) DeleteConfigMap(n string) error {
	e := m.Called().Error(0)
	return e
}
func (m *KubeMock) GetConstraintTemplates(prefix string) (*map[string]interface{}, error) {
	a := m.Called()
	return a.Get(0).(*map[string]interface{}), a.Error(1)
}
func (m *KubeMock) GetRawResourcesByGrp(g string) (*K8SJSON, error) {
	a := m.Called()
	return a.Get(0).(*K8SJSON), a.Error(1)
}

func (m *KubeMock) GetClusterRolesByResource(r string) (*[]rbacv1.ClusterRole, error) {
	a := m.Called()
	return a.Get(0).(*[]rbacv1.ClusterRole), a.Error(1)
}
func (m *KubeMock) GetClusterRoles() (*rbacv1.ClusterRoleList, error) {
	a := m.Called()
	return a.Get(0).(*rbacv1.ClusterRoleList), a.Error(1)
}
