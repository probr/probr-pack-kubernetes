package cliflags

import (
	"fmt"
	"os"
	"strings"

	"github.com/citihub/probr/config"
)

// HandleRequestForRequiredVars will execute the logic for `./probr show-requirements (<PACK>)`
func HandleRequestForRequiredVars() {
	////log.Printf("[DEBUG] Checking for CLI options or flags")
	if os.Args[1] == "show-requirements" {
		////log.Printf("[INFO] CLI option 'show-requirements' was found")
		for req := range config.Requirements {
			if len(os.Args) > 2 {
				// Show specified
				if strings.ToLower(req) == strings.ToLower(os.Args[2]) {
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

// HandlePackOption will execute the logic necessary for `./probr run <PACK>`
func HandlePackOption() {
	if os.Args[1] == "run" {
		////log.Printf("[DEBUG] CLI option 'run' was found. Args: %s", os.Args)
		for _, pack := range config.GetPacks() {
			if strings.ToLower(pack) == strings.ToLower(os.Args[2]) {
				////log.Printf("[INFO] CLI Option specified to run only %s service pack", pack)
				config.Vars.Meta.RunOnly = pack
			}
		}
		if config.Vars.Meta.RunOnly == "" {
			// If run was specified without a valid pack name, exit
			////log.Printf("[ERROR] Expected a service pack name.\n\nUsage: ./probr run <PACK-NAME>\n\n")
			os.Exit(2)
		}
		// Remove the "run" and "PACK-NAME" arguments to prevent interference with flag handling
		copy(os.Args[1:], os.Args[3:])
		os.Args = os.Args[:len(os.Args)-2]
		////log.Printf("[DEBUG] Args after 'run %s': %s", config.Vars.Meta.RunOnly, os.Args)
	}
}
