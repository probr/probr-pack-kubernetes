module github.com/probr/probr-pack-kubernetes

go 1.14

require (
	github.com/cucumber/godog v0.11.0
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/markbates/pkger v0.17.1
	github.com/probr/probr-sdk v0.1.5
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211015210444-4f30a5c0130f // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	golang.org/x/text v0.3.7 // indirect
	k8s.io/api v0.19.6
)

// For Development Only
// replace github.com/probr/probr-sdk => ../probr-sdk
