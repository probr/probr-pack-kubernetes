package pod_security_policy

import (
	"log"

	"github.com/citihub/probr/internal/utils"
	"github.com/citihub/probr/service_packs/kubernetes"
)

//AzK8sConstraintTemplate captures the Azure specific constraint templates that are the result of applying
//an Azure Policy which can be used to support PodSecurityPolicy like behaviour.
//Implements securitypolicyprovider and is the preferred way of determining constraints on
//an AKS cluster.
type AzK8sConstraintTemplate struct {
	k kubernetes.Kubernetes

	//set of constraints applied to cluster
	//(can view via kubectl get constrainttemplate)
	azConstraints *map[string]interface{}
}

const (
	azK8sPrefix                         = "k8sazure"
	azK8sAllowedCapabilities            = "k8sazureallowedcapabilities"
	azK8sAllowedSeccomp                 = "k8sazureallowedseccomp"
	azK8sAllowedUsersGroups             = "k8sazureallowedusersgroups"
	azK8sBlockHostNamespace             = "k8sazureblockhostnamespace"
	azK8sContainerNoPrivilege           = "k8sazurecontainernoprivilege"
	azK8sContainerNoPrivilegeEscalation = "k8sazurecontainernoprivilegeescalation"
	azK8sHostNetworkingPorts            = "k8sazurehostnetworkingports"
	azK8sAllowedVolumeTypes             = "k8sazurevolumetypes"
)

//NewAzK8sConstraintTemplate constructs a new AzK8sConstraintTemplate using the supplied kubernetes instance.
func NewAzK8sConstraintTemplate(k kubernetes.Kubernetes) *AzK8sConstraintTemplate {
	a := &AzK8sConstraintTemplate{}
	a.k = k

	return a
}

//NewDefaultAzK8sConstraintTemplate constructs a new AzK8sConstraintTemplate using the default kubernetes instance.
func NewDefaultAzK8sConstraintTemplate() *AzK8sConstraintTemplate {
	a := &AzK8sConstraintTemplate{}
	a.k = kubernetes.GetKubeInstance()

	return a
}

//HasSecurityPolicies provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasSecurityPolicies() (*bool, error) {
	c, err := az.getConstraints()

	if err != nil {
		return nil, err
	}
	if c != nil && len(*c) > 0 {
		return utils.BoolPtr(true), nil
	}
	return utils.BoolPtr(false), nil
}

//HasPrivilegedAccessRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasPrivilegedAccessRestriction() (*bool, error) {
	return az.hasConstraint(azK8sContainerNoPrivilege)
}

//HasHostPIDRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasHostPIDRestriction() (*bool, error) {
	//"blockhostnamespace" covers host PID, IPC & network
	return az.hasConstraint(azK8sBlockHostNamespace)
}

//HasHostIPCRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasHostIPCRestriction() (*bool, error) {
	//"blockhostnamespace" covers host PID, IPC & network
	return az.hasConstraint(azK8sBlockHostNamespace)
}

//HasHostNetworkRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasHostNetworkRestriction() (*bool, error) {
	//"blockhostnamespace" covers host PID, IPC & network
	return az.hasConstraint(azK8sBlockHostNamespace)
}

//HasAllowPrivilegeEscalationRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	return az.hasConstraint(azK8sContainerNoPrivilegeEscalation)
}

//HasRootUserRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasRootUserRestriction() (*bool, error) {
	return az.hasConstraint(azK8sAllowedUsersGroups)
}

//HasNETRAWRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasNETRAWRestriction() (*bool, error) {
	//TODO: at time of writing, not clear that any AZ policy/constraint convers NET_RAW
	return utils.BoolPtr(false), nil
}

//HasAllowedCapabilitiesRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasAllowedCapabilitiesRestriction() (*bool, error) {
	//Az AllowedCapabilities covers both allowed & assigned capabilities
	return az.hasConstraint(azK8sAllowedCapabilities)
}

//HasAssignedCapabilitiesRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasAssignedCapabilitiesRestriction() (*bool, error) {
	//Az AllowedCapabilities covers both allowed & assigned capabilities
	return az.hasConstraint(azK8sAllowedCapabilities)
}

//HasHostPortRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasHostPortRestriction() (*bool, error) {
	return az.hasConstraint(azK8sHostNetworkingPorts)
}

//HasVolumeTypeRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasVolumeTypeRestriction() (*bool, error) {
	return az.hasConstraint(azK8sAllowedVolumeTypes)
}

//HasSeccompProfileRestriction provides the AzK8sConstraintTemplate implementation of SecurityPolicyProvider.
func (az *AzK8sConstraintTemplate) HasSeccompProfileRestriction() (*bool, error) {
	return az.hasConstraint(azK8sAllowedSeccomp)
}

func (az *AzK8sConstraintTemplate) hasConstraint(cst string) (*bool, error) {
	c, err := az.getConstraints()

	if err != nil {
		return nil, err
	}
	if c == nil {
		return utils.BoolPtr(false), nil
	}

	_, b := (*c)[cst]

	log.Printf("[DEBUG] Azure Constraint template %q. Result %t.", cst, b)
	return &b, nil
}

func (az *AzK8sConstraintTemplate) getConstraints() (*map[string]interface{}, error) {
	if az.azConstraints == nil {
		//get "k8sazure" constraints:
		c, err := az.k.GetConstraintTemplates(azK8sPrefix)

		if err != nil {
			return nil, err
		}

		//otherwise set it
		az.azConstraints = c

		log.Printf("[DEBUG] Azure Constraints (%d): %v ", len(*az.azConstraints), az.azConstraints)
	}

	return az.azConstraints, nil
}
