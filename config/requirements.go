package config

var Requirements = map[string][]string{
	"Storage":    []string{"Provider"},
	"Kubernetes": []string{"AuthorisedContainerRegistry", "UnauthorisedContainerRegistry"},
}
