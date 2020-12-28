package kubernetes_pack

import (
	"log"

	"github.com/citihub/probr/internal/config"
	"github.com/citihub/probr/internal/coreengine"
	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/kubernetes/container_registry_access"
	"github.com/citihub/probr/service_packs/kubernetes/general"
	"github.com/citihub/probr/service_packs/kubernetes/iam"
	"github.com/citihub/probr/service_packs/kubernetes/internet_access"
	"github.com/citihub/probr/service_packs/kubernetes/pod_security_policy"
)

func GetProbes() []coreengine.Probe {
	if config.Vars.ServicePacks.Kubernetes.IsExcluded() {
		file, line := utils.CallerFileLine()
		log.Printf("[WARN] %s:%v: Ignoring Kubernetes service pack due to required vars not being present.", file, line)
		return nil
	}
	return []coreengine.Probe{
		container_registry_access.Probe,
		general.Probe,
		pod_security_policy.Probe,
		internet_access.Probe,
		iam.Probe,
	}
}
