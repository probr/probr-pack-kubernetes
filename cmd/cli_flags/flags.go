package cli_flags

import (
	"flag"
	"log"

	"github.com/citihub/probr/internal/config"
)

type flagHandlerFunc func(v interface{})

type Flag struct {
	Handler flagHandlerFunc
	Value   interface{}
}

var flags []Flag

func (f Flag) executeHandler() {
	f.Handler(f.Value)
}

func HandleFlags() {

	stringFlag("varsfile", "path to config file", varsFileHandler)
	stringFlag("kubeconfig", "kube config file", kubeConfigHandler)
	stringFlag("cucumberdir", "cucumber output directory", cucumberDirHandler)
	stringFlag("loglevel", "set log level", loglevelHandler)
	stringFlag("tags", "feature tags to include or exclude", tagsHandler)
	boolFlag("silent", "disable visual runtime indicator, useful for CI tasks", silentHandler)
	boolFlag("nosummary", "switch off summary output", nosummaryHandler)
	flag.Parse()

	for _, f := range flags {
		f.executeHandler()
	}
}

func stringFlag(name string, usage string, handler flagHandlerFunc) {
	f := Flag{
		Handler: handler,
		Value:   new(string),
	}
	v := f.Value.(*string)
	flag.StringVar(v, name, "", usage)
	flags = append(flags, f)
}

func boolFlag(name string, usage string, handler flagHandlerFunc) {
	f := Flag{
		Handler: handler,
		Value:   new(bool),
	}
	v := f.Value.(*bool)
	flag.BoolVar(v, name, false, usage)
	flags = append(flags, f)
}

// Note:
// Even though it's a bit ugly, using things like `*v.(*string)` comes from accepting bool, string, and other flag types

// varsFileHandler initializes configuration with VarsFile overriding env vars & defaults
func varsFileHandler(v interface{}) {
	err := config.Init(*v.(*string))
	if err != nil {
		log.Fatalf("[ERROR] Could not create config from provided filepath: %v", v.(*string))
	} else if len(*v.(*string)) > 0 {
		log.Printf("[NOTICE] Config read from file '%v', but may still be overridden by CLI flags.", v.(*string))
	} else {
		log.Printf("[NOTICE] No configuration variables file specified. Using environment variabls and defaults only.")
	}
}

// cucumberDirHandler
func cucumberDirHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden via command line")
		config.Vars.CucumberDir = *v.(*string)
	}
}

// loglevelHandler validates provided value and sets output accordingly
func loglevelHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		if (*v.(*string) != "DEBUG") && (*v.(*string) != "INFO") && (*v.(*string) != "NOTICE") && (*v.(*string) != "WARN") && (*v.(*string) != "ERROR") {
			log.Fatalf("[ERROR] Unknown loglevel specified: %v. Must be one of 'DEBUG', 'INFO', 'NOTICE', 'WARN', 'ERROR'", v.(*string))
		}
		config.Vars.LogLevel = *v.(*string)
	}
}

func tagsHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.Tags = *v.(*string)
		log.Printf("[NOTICE] tags have been added via command line.")
	}
	if len(config.Vars.GetTags()) == 0 {
		log.Printf("[NOTICE] No tags specified.")
	}
}

func kubeConfigHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath = *v.(*string)
		log.Printf("[NOTICE] Kubeconfig path has been overridden via command line")
	}
	if len(config.Vars.ServicePacks.Kubernetes.KubeConfigPath) == 0 {
		log.Printf("[NOTICE] No kubeconfig path specified. Falling back to default paths.")
	}
}

func silentHandler(v interface{}) {
	config.Vars.Silent = isFlagPassed("silent")
}

func nosummaryHandler(v interface{}) {
	config.Vars.NoSummary = isFlagPassed("nosummary")
}

func isFlagPassed(flagName string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			found = true
		}
	})
	return found
}
