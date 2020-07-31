package azure

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

type azPolicy struct {
	scope       *string
	policyType  *string
	displayName *string
	uuid        *string
}

// AZSecurityPolicyProvider ...
type AZSecurityPolicyProvider struct {
	policiesByType map[string]*azPolicy
}

const (
	azPSPLinuxRestricted              = "AZPSPLinuxRestricted"
	azPSPContainerImage               = "AZPSPContainerImage"
	azPSPContainerPrivilegeEscalation = "AZPSPContainerPrivilegeEscalation"
	azPSPHostPIDHostIPCNS             = "AZPSPHostPIDHostIPCNS"
)

var azPolicyUUIDToProbrPolicy = make(map[string]string)

func init() {
	//TODO: should really map these to types, but not sure of the how important this is, ie. is it sufficent to say "some polices"?
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policySetDefinitions/42b8ef37-b724-4e24-bbc8-7a7708edfe00"] = azPSPLinuxRestricted
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/febd0533-8e55-448f-b837-bd0e06f16469"] = azPSPContainerImage
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/1c6e92c9-99f0-4e55-9cf2-0c234dc48f99"] = azPSPContainerPrivilegeEscalation
	azPolicyUUIDToProbrPolicy["/providers/Microsoft.Authorization/policyDefinitions/47a1ee2f-2a2a-4576-bf2a-e0e36709c2b8"] = azPSPHostPIDHostIPCNS
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
	pc, err := p.getPolicies()

	if err != nil {
		return nil, err
	}

	//interograte the policy types ..
	//in Azure, "Privilege access" can be via a policy set or specific policy
	//the "AZPSPLinuxRestricted" set has this restriction.  Try this first
	//and then "AZPSPContainerPrivilegeEscalation" if that's not set
	_, b := (*pc)[azPSPLinuxRestricted]
	if b {
		return &b, nil
	}

	//try "AZPSPContainerPrivilegeEscalation"
	_, b = (*pc)[azPSPContainerPrivilegeEscalation]

	return &b, nil
}

// HasHostPIDRestriction ...
func (p *AZSecurityPolicyProvider) HasHostPIDRestriction() (*bool, error) {
	pc, err := p.getPolicies()

	if err != nil {
		return nil, err
	}

	//interograte the policy types ..
	//in Azure, "HostPID Restriction" can be via a policy set or specific policy
	//the "AZPSPLinuxRestricted" set has this restriction.  Try this first
	//and then "AZPSPHostPIDHostIPCNS" if that's not set
	_, b := (*pc)[azPSPLinuxRestricted]
	if b {
		return &b, nil
	}

	//try "AZPSPHostPIDHostIPCNS"
	_, b = (*pc)[azPSPHostPIDHostIPCNS]

	return &b, nil
}

func (p *AZSecurityPolicyProvider) getPolicies() (*map[string]*azPolicy, error) {

	if len(p.policiesByType) > 0 {
		//already got 'em
		return &p.policiesByType, nil
	}

	s := os.Getenv("AZURE_SUBSCRIPTION_ID")
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

		log.Printf("[NOTICE] azPolicy %v %v %v", *azp.uuid, *azp.displayName, *azp.scope)

		//look up the "type" based on the uuid of the policy
		t, exists := azPolicyUUIDToProbrPolicy[*azp.uuid]
		if !exists {
			//set t to a default type
			t = "AZPSPUnkown"
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
