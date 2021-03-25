package policy

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	azureutil "github.com/citihub/probr/service_packs/storage/azure"
)

// AssignmentBySubscription gets a Policy Assignment by Policy Assignment name, scoped to a Subscription.
func AssignmentBySubscription(ctx context.Context, subscriptionID, name string) (policy.Assignment, error) {
	scope := "/subscriptions/" + subscriptionID
	log.Printf("[DEBUG] Getting Policy Assignment with subscriptionID: %v", scope)
	return assignmentClient().Get(ctx, scope, name)
}

// AssignmentByManagementGroup gets a Policy Assignment by Policy Assignment name, scoped to a Managed Group.
func AssignmentByManagementGroup(ctx context.Context, managementGroup, name string) (policy.Assignment, error) {
	scope := "/providers/Microsoft.Management/managementGroups/" + managementGroup
	log.Printf("[DEBUG] Getting Policy Assignment with scope: %v", scope)
	return assignmentClient().Get(ctx, scope, name)
}

func assignmentClient() policy.AssignmentsClient {
	c := policy.NewAssignmentsClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Policy Assignment client: %v", err)
	}
	return c
}
