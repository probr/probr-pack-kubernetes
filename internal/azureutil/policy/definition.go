package policy

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-01-01/policy"
	"github.com/Azure/go-autorest/autorest/azure/auth"

	"github.com/citihub/probr/internal/azureutil"
)

// DefinitionByName get a Policy Definition by name.
func DefinitionByName(ctx context.Context, name string) (policy.Definition, error) {
	return definitionClient().Get(ctx, name)
}

func definitionClient() policy.DefinitionsClient {
	c := policy.NewDefinitionsClient(azureutil.SubscriptionID())
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		c.Authorizer = a
	} else {
		log.Fatalf("Unable to authorise Policy Definition client: %v", err)
	}
	return c
}
