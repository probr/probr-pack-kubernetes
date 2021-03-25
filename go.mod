module github.com/citihub/probr-pack-kubernetes

go 1.14

require (
	github.com/Azure/aad-pod-identity v1.7.0
	github.com/Azure/azure-sdk-for-go v44.2.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.12.0
	github.com/Azure/go-autorest/autorest v0.11.0
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.0
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/briandowns/spinner v1.12.0
	github.com/citihub/probr-sdk v0.0.7
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/hashicorp/go-hclog v0.15.0
	github.com/hashicorp/logutils v1.0.0
	github.com/markbates/pkger v0.17.1
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
)

//replace github.com/citihub/probr-sdk => ../probr-sdk

//Line above is intended to be used during dev only when editing modules locally.
