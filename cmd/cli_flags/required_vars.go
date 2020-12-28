package cli_flags

import (
	"fmt"
	"os"
)

func HandleRequestForRequiredVars() {
	if os.Args[1] == "show-required" {
		switch os.Args[2] {
		case "kubernetes":
			respond("Kubernetes", "AuthorisedContainerRegistry", "UnauthorisedContainerRegistry")
		case "storage":
			respond("Storage", "Provider")
		default:
			fmt.Printf("Unknown service pack specified, cannot get required variables")
		}
		os.Exit(0) // Don't continue if this option is called
	}
}

func respond(pack string, vars ...string) {
	fmt.Printf("Required variables for %s service pack:\n", pack)
	for _, v := range vars {
		fmt.Printf("    %s\n", v)
	}
}
