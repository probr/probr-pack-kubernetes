module github.com/citihub/probr-pack-kubernetes

go 1.14

require (
	github.com/citihub/probr-sdk v0.0.27
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/markbates/pkger v0.17.1
	k8s.io/api v0.19.6
)

// For Development Only
// replace github.com/citihub/probr-sdk => ../probr-sdk
