package cliflags

import (
	"flag"
	"os"

	"github.com/citihub/probr/config"
	"github.com/citihub/probr/utils"
)

type flagHandlerFunc func(v interface{})

// Flag holds the user-provided value for the flag, and the function to be run within executeHandler
type Flag struct {
	// exported to avoid conflict with the default `flag` provided in go
	Handler flagHandlerFunc
	Value   interface{}
}

var flags []Flag

func (f Flag) executeHandler() {
	f.Handler(f.Value)
}

// HandleFlags executes the logic for any flags that are provided via `./probr (--<FLAG>)`
func HandleFlags() {

	stringFlag("varsfile", "path to config file", varsFileHandler)
	stringFlag("loglevel", "set log level", loglevelHandler)
	stringFlag("kubeconfig", "kube config file", kubeConfigHandler)
	stringFlag("writedirectory", "output directory", writeDirHandler)
	stringFlag("tags", "feature tags to include or exclude", tagsHandler)
	stringFlag("resultsformat", "set the bdd results format (default = cucumber)", resultsformatHandler)
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
		//log.Fatalf("[ERROR] error returned from config.Init: %v", err)
	} else if len(*v.(*string)) > 0 {
		config.Vars.VarsFile = *v.(*string)
		//log.Printf("[INFO] Config read from file '%v', but may still be overridden by CLI flags.", v.(*string))
	} else {
		//log.Printf("[NOTICE] No configuration variables file specified. Using environment variabls and defaults only.")
	}
}

// writeDirHandler
func writeDirHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		//log.Printf("[NOTICE] Output Directory has been overridden via command line")
		config.Vars.WriteDirectory = *v.(*string)
	}
}

// loglevelHandler validates provided value and sets output accordingly
func loglevelHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		levels := []string{"DEBUG", "INFO", "NOTICE", "WARN", "ERROR"}
		_, found := utils.FindString(levels, *v.(*string))
		if !found {
			//log.Fatalf("[ERROR] Unknown loglevel specified: '%s'. Must be one of %v", *v.(*string), levels)
		} else {
			config.Vars.LogLevel = *v.(*string)
			config.SetLogFilter(config.Vars.LogLevel, os.Stderr)
		}
	}
}

func resultsformatHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		options := []string{"cucumber", "events", "junit", "pretty", "progress"}
		_, found := utils.FindString(options, *v.(*string))
		if !found {
			//log.Fatalf("[ERROR] Unknown resultsformat specified: '%s'. Must be one of %v", *v.(*string), options)
		} else {
			config.Vars.ResultsFormat = *v.(*string)
			config.SetLogFilter(config.Vars.ResultsFormat, os.Stderr)
		}
	}
}

func tagsHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.Tags = *v.(*string)
		//log.Printf("[INFO] tags have been added via command line.")
	}
}

func kubeConfigHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath = *v.(*string)
		//log.Printf("[NOTICE] Kubeconfig path has been overridden via command line")
	}
	if len(config.Vars.ServicePacks.Kubernetes.KubeConfigPath) == 0 {
		//log.Printf("[NOTICE] No kubeconfig path specified. Falling back to default paths.")
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
