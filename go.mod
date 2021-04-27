module github.com/citihub/probr-pack-kubernetes

go 1.14

require (
	github.com/Azure/aad-pod-identity v1.7.5
	github.com/citihub/probr-sdk v0.0.19
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0
	github.com/markbates/pkger v0.17.1
	k8s.io/api v0.19.6
	k8s.io/apimachinery v0.19.6
	k8s.io/client-go v0.19.6
)

//replace github.com/citihub/probr-sdk => ../probr-sdk

//Line above is intended to be used during dev only when editing modules locally.
