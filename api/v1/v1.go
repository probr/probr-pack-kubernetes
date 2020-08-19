package v1

import (
	"gitlab.com/citihub/probr/internal/clouddriver/kubernetes"
	"gitlab.com/citihub/probr/internal/coreengine"

	"github.com/google/uuid"

	_ "gitlab.com/citihub/probr/internal/config" //needed for logging
	"gitlab.com/citihub/probr/test/features"
	_ "gitlab.com/citihub/probr/test/features/clouddriver"
	_ "gitlab.com/citihub/probr/test/features/kubernetes/containerregistryaccess" //needed to run init on TestHandlers
	_ "gitlab.com/citihub/probr/test/features/kubernetes/internetaccess"          //needed to run init on TestHandlers
	_ "gitlab.com/citihub/probr/test/features/kubernetes/podsecuritypolicy"       //needed to run init on TestHandlers
)

//TODO: revise when interface this bit up ...
var kube = kubernetes.GetKubeInstance()

func addTest(tm *coreengine.TestStore, n string, g coreengine.Group, c coreengine.Category) {

	td := coreengine.TestDescriptor{Group: g, Category: c, Name: n}

	uuid1 := uuid.New().String()
	sat := coreengine.Pending

	test := coreengine.Test{
		UUID:           &uuid1,
		TestDescriptor: &td,
		Status:         &sat,
	}

	//add - don't worry about the rtn uuid
	tm.AddTest(&test)
}

// RunAllTests ...
func RunAllTests() (int, error) {
	// MUST run after SetIOPaths
	// get the test mgr
	tm := coreengine.NewTestManager()

	//add some tests and add them to the TM - we need to tidy this up!
	addTest(tm, "container_registry_access", coreengine.Kubernetes, coreengine.ContainerRegistryAccess)
	addTest(tm, "internet_access", coreengine.Kubernetes, coreengine.InternetAccess)
	addTest(tm, "pod_security_policy", coreengine.Kubernetes, coreengine.PodSecurityPolicies)
	addTest(tm, "account_manager", coreengine.CloudDriver, coreengine.General)

	//exec 'em all (for now!)
	return tm.ExecAllTests()

}

// SetIOPaths ...
func SetIOPaths(i string, o string) {
	kube.SetKubeConfigFile(&i)
	features.SetOutputDirectory(&o)
}
