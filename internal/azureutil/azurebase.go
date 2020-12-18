package azureutil

import (
	"log"
	"os"

	"github.com/citihub/probr/internal/config"
)

const (
	PolicyAssignmentManagementGroup string = "AZURE_POLICY_ASSIGNMENT_MANAGEMENT_GROUP"
)

var prefix string
var rgName string

//ResourceGroup is a singleton that generates or returns a test Resource Group name in the form 'test[a-z]{6}resourceGP'.
func ResourceGroup() string {
	if rgName == "" {
		rgName = randomPrefix() + "resourceGP"
	}
	return rgName
}

//Location returns the location in which the tests should be executed, driven by environment variable AZURE_LOCATION.
func Location() string {
	return getFromEnvVar("AZURE_LOCATION")
}

//Location returns the Subscription in which the tests should be executed, driven by environment variable AZURE_SUBSCRIPTION_ID.
func SubscriptionID() string {
	return getFromEnvVar("AZURE_SUBSCRIPTION_ID")
}

//
func ManagementGroup() string {
	return config.Vars.CloudProviders.Azure.ManagementGroup
}

func randomPrefix() string {
	if prefix == "" {
		prefix = "test" + RandString(6) + ""
	}
	return "test" + RandString(6) + ""
}

func getFromEnvVar(varName string) string {
	v, b := os.LookupEnv(varName)
	if !b {
		log.Printf("[ERROR] Environment variable \"%v\" is not defined", varName)
	}
	return v
}
