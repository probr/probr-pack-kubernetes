package cli_flags

import (
	"fmt"
	"os"

	"github.com/citihub/probr/internal/config"
)

func HandleRequestForRequiredVars() {
	if os.Args[1] == "show-requirements" {
		for req := range config.Requirements {
			if len(os.Args) > 2 {
				// Show specified
				if os.Args[2] == req {
					respond(req, config.Requirements[req]...)
					os.Exit(0)
				}
			} else {
				// Show all
				respond(req, config.Requirements[req]...)
			}
		}
		os.Exit(0) // Never run probr if 'show-requirements' is called
	}
}

func respond(pack string, vars ...string) {
	fmt.Printf("Required variables for %s:\n", pack)
	for _, v := range vars {
		fmt.Printf("    %s\n", v)
	}
}
