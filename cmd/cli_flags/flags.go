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

	stringFlag("varsFile", "path to config file", varsFileHandler)
	stringFlag("outputDir", "output directory", outputDirHandler) // Must run prior to creating outputType flag
	stringFlag("outputType", "output defaults to write in memory, if 'IO' will write to specified output directory", outputTypeHandler)
	stringFlag("tags", "test tags, e.g. -tags=\"@CIS-1.2.3, @CIS-4.5.6\".", tagsHandler)
	stringFlag("kubeConfig", "kube config file", kubeConfigHandler)
	boolFlag("silent", "Disable visual runtime indicator, useful for CI tasks", silentHandler)
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

// varsFileHandler initializes configuration with varsFile overriding env vars & defaults
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

// outputDirHandler
func outputDirHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden via command line")
		config.Vars.CucumberDir = *v.(*string)
	}
}

// outputTypeHandler validates provided value and sets output accordingly
func outputTypeHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		if *v.(*string) == "IO" {
			log.Printf("[NOTICE] Probr results will be written to files in the specified output directory: %v", v.(*string))
		} else if *v.(*string) == "INMEM" {
			log.Printf("[NOTICE] Output type specified as INMEM: Results will not be handled by the CLI. Refer to the Summary Log for a results summary.")
		} else {
			log.Fatalf("[ERROR] Unknown output type specified: %v. Please use 'IO' or 'INMEM'", v.(*string))
		}
		config.Vars.OutputType = *v.(*string)
	}
}

func tagsHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.Tags = *v.(*string)
		log.Printf("[NOTICE] Tags have been added via command line.")
	}
	if len(config.Vars.GetTags()) == 0 {
		log.Printf("[NOTICE] No tags specified.")
	}
}

func kubeConfigHandler(v interface{}) {
	if len(*v.(*string)) > 0 {
		config.Vars.KubeConfigPath = *v.(*string)
		log.Printf("[NOTICE] Kubeconfig path has been overridden via command line")
	}
	if len(config.Vars.KubeConfigPath) == 0 {
		log.Printf("[NOTICE] No kubeconfig path specified. Falling back to default paths.")
	}
}

func silentHandler(v interface{}) {
	config.Vars.Silent = isFlagPassed("silent")
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
