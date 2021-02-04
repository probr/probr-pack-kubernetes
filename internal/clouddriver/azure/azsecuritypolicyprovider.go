package azure

import (
	"context"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	azureutil "github.com/citihub/probr/service_packs/storage/azure"
)

type azPolicy struct {
	scope       *string
	policyType  *string
	displayName *string
	uuid        *string
}

//AZSecurityPolicyProvider queries the policies applied on the supplied subscription and resource group.  This may be deprecated in favour
//of azk8sconstrainttemplate which queries the kubernetes cluster directly.
//TODO: decide if this should be kept.
type AZSecurityPolicyProvider struct {
	policiesByType map[string]*azPolicy
}

const (
	azPSPLinuxRestricted              = "AZPSPLinuxRestricted"
	azPSPContainerImage               = "AZPSPContainerImage"
	azPSPContainerPrivilegeEscalation = "AZPSPContainerPrivilegeEscalation"
	azPSPHostPIDHostIPCNS             = "AZPSPHostPIDHostIPCNS"
	azPSPContainerPrivileged          = "AZPSPContainerPrivileged"
	azPSPApprovedUsersAndGroups       = "AZPSPApprovedUsersAndGroups"
	azPSPAllowedCapabilitiesOnly      = "AZPSPAllowedCapabilitiesOnly"
	azPSPApprovedPortRangeOnly        = "AZPSPApprovedPortRangeOnly"
	azPSPApprovedVolumeTypeOnly       = "AZPSPApprovedVolumeTypeOnly"
	azPSPApprovedSeccompProfile       = "AZPSPApprovedSeccompProfile"
)

var azPolicyUUIDToProbrPolicy = make(map[string]string)

func init() {
	//TODO: should really map these to types, but not sure of the how important this is, ie. is it sufficient to say "some polices"?
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policySetDefinitions/42b8ef37-b724-4e24-bbc8-7a7708edfe00"] = azPSPLinuxRestricted
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/febd0533-8e55-448f-b837-bd0e06f16469"] = azPSPContainerImage
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/1c6e92c9-99f0-4e55-9cf2-0c234dc48f99"] = azPSPContainerPrivilegeEscalation
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/47a1ee2f-2a2a-4576-bf2a-e0e36709c2b8"] = azPSPHostPIDHostIPCNS
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/95edb821-ddaf-4404-9732-666045e056b4"] = azPSPContainerPrivileged
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/f06ddb64-5fa3-4b77-b166-acb36f7f6042"] = azPSPApprovedUsersAndGroups
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/c26596ff-4d70-4e6a-9a30-c2506bd2f80c"] = azPSPAllowedCapabilitiesOnly
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/82985f06-dc18-4a48-bc1c-b9f4f0098cfe"] = azPSPApprovedPortRangeOnly
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/16697877-1118-4fb1-9b65-9898ec2509ec"] = azPSPApprovedVolumeTypeOnly
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/975ce327-682c-4f2e-aa46-b9598289b86c"] = azPSPApprovedSeccompProfile
}

//NewAzPolicyProvider ...
func NewAzPolicyProvider() *AZSecurityPolicyProvider {
	return &AZSecurityPolicyProvider{
		policiesByType: make(map[string]*azPolicy),
	}
}

// HasSecurityPolicies ...
func (p *AZSecurityPolicyProvider) HasSecurityPolicies() (*bool, error) {
	pc, err := p.getPolicies()

	if err != nil {
		return nil, err
	}

	b := len(*pc) > 0

	return &b, nil
}

// HasPrivilegedAccessRestriction ...
func (p *AZSecurityPolicyProvider) HasPrivilegedAccessRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPContainerPrivileged})
}

// For the following, the hostPID, hostIPC and hostNetwork restrictions are wrapped together
// in Azure policies.  This is either in the general 'Linux Restricted' policy set or in the
// HostPID/HostIPC/HostNetwork policy:

// HasHostPIDRestriction ...
func (p *AZSecurityPolicyProvider) HasHostPIDRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPHostPIDHostIPCNS})
}

// HasHostIPCRestriction ...
func (p *AZSecurityPolicyProvider) HasHostIPCRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPHostPIDHostIPCNS})
}

// HasHostNetworkRestriction ...
func (p *AZSecurityPolicyProvider) HasHostNetworkRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPHostPIDHostIPCNS})
}

// HasAllowPrivilegeEscalationRestriction ...
func (p *AZSecurityPolicyProvider) HasAllowPrivilegeEscalationRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPContainerPrivilegeEscalation})
}

// HasRootUserRestriction ...
func (p *AZSecurityPolicyProvider) HasRootUserRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPApprovedUsersAndGroups})
}

// HasNETRAWRestriction ...
func (p *AZSecurityPolicyProvider) HasNETRAWRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted})
}

// HasAllowedCapabilitiesRestriction ...
func (p *AZSecurityPolicyProvider) HasAllowedCapabilitiesRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPAllowedCapabilitiesOnly})
}

// HasAssignedCapabilitiesRestriction ...
func (p *AZSecurityPolicyProvider) HasAssignedCapabilitiesRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted})
}

// HasHostPortRestriction ...
func (p *AZSecurityPolicyProvider) HasHostPortRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPApprovedPortRangeOnly})
}

// HasVolumeTypeRestriction ...
func (p *AZSecurityPolicyProvider) HasVolumeTypeRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPApprovedPortRangeOnly})
}

// HasSeccompProfileRestriction ...
func (p *AZSecurityPolicyProvider) HasSeccompProfileRestriction() (*bool, error) {
	return p.checkForRestrictions(&[]string{azPSPLinuxRestricted, azPSPApprovedPortRangeOnly})
}

func (p *AZSecurityPolicyProvider) checkForRestrictions(res *[]string) (*bool, error) {
	pc, err := p.getPolicies()

	if err != nil {
		return nil, err
	}

	//check for the policies we've been given
	for _, r := range *res {
		_, b := (*pc)[r]
		if b {
			return &b, nil
		}
	}

	//haven't found any if we get to here ...
	b := false
	return &b, nil
}

func (p *AZSecurityPolicyProvider) getPolicies() (*map[string]*azPolicy, error) {

	if len(p.policiesByType) > 0 {
		//already got 'em
		return &p.policiesByType, nil
	}

	s := azureutil.SubscriptionID()
	log.Printf("[INFO] Using Azure Sub: %v", s)

	scope := "/subscriptions/" + s
	log.Printf("[DEBUG] Getting Policy Assignment with subscriptionID: %v", scope)

	ac := assignmentClient(s)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	l, err := ac.List(ctx, "")

	if err != nil {
		log.Printf("[ERROR] Error getting Azure Policies: %v", err)
		return nil, err
	}

	alr := l.Response()
	// convert the AZ representation into something sensible ..
	for _, r := range *alr.Value {
		azp := azPolicy{}

		//TODO: also map the type here
		azp.uuid = r.AssignmentProperties.PolicyDefinitionID //name in azure policy is the uuid
		azp.displayName = r.AssignmentProperties.DisplayName //display name == 'name'
		azp.scope = r.AssignmentProperties.Scope

		log.Printf("[DEBUG] azPolicy %v %v %v", *azp.uuid, *azp.displayName, *azp.scope)

		//look up the "type" based on the uuid of the policy
		t, exists := azPolicyUUIDToProbrPolicy[*azp.uuid]
		if !exists {
			//set t to a default type
			t = "AZPSPUnknown"
		}

		p.policiesByType[t] = &azp
	}

	return &p.policiesByType, nil
}

func assignmentClient(sub string) policy.AssignmentsClient {

	c := policy.NewAssignmentsClient(sub)
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Printf("[ERROR] Unable to authorise Azure Policy Assignment client: %v", err)
	}
	return c
}
