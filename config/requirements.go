package config

var Requirements = map[string][]string{
	"Storage":    []string{"Provider"},
	"APIM":       []string{"Provider"},
	"Kubernetes": []string{"AuthorisedContainerRegistry", "UnauthorisedContainerRegistry"},
}
