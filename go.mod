module github.com/citihub/probr-pack-kubernetes

go 1.14

require (
	github.com/Azure/aad-pod-identity v1.7.0
	github.com/Azure/go-autorest/autorest v0.11.0 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.2 // indirect
	github.com/citihub/probr-sdk v0.0.15
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/markbates/pkger v0.17.1
	golang.org/x/sys v0.0.0-20200828194041-157a740278f4 // indirect
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
)

// replace github.com/citihub/probr-sdk => ../probr-sdk

//Line above is intended to be used during dev only when editing modules locally.
