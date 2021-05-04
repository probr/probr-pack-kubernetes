module github.com/citihub/probr-pack-kubernetes

go 1.14

require (
	github.com/Azure/aad-pod-identity v1.7.5
	github.com/Azure/azure-sdk-for-go v53.3.0+incompatible
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/citihub/probr-sdk v0.0.24
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/markbates/pkger v0.17.1
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
	k8s.io/client-go v0.19.6
)

// For Development Only
// replace github.com/citihub/probr-sdk => ../probr-sdk
