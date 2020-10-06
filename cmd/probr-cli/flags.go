package main

import (
	"flag"
	"log"

	"github.com/citihub/probr"
	"github.com/citihub/probr/internal/config"
)

type flagHandlerFunc func(v *string)

type Flag struct {
	Handler flagHandlerFunc
	Value   *string
}

var flags []Flag

func (f *Flag) executeHandler() {
	f.Handler(f.Value)
}

func handleFlags() {

	createFlag("varsFile", "", "path to config file", varsFileHandler)
	createFlag("outputDir", "", "output directory", outputDirHandler) // Must run prior to creating outputType flag
	createFlag("outputType", "INMEM", "output defaults to write in memory, if 'IO' will write to specified output directory", outputTypeHandler)
	createFlag("tags", "", "test tags, e.g. -tags=\"@CIS-1.2.3, @CIS-4.5.6\".", tagsHandler)
	createFlag("kubeConfig", "", "kube config file", kubeConfigHandler)
	flag.Parse()

	for _, f := range flags {
		f.executeHandler()
	}
}

func createFlag(n string, d string, t string, h flagHandlerFunc) {
	f := Flag{
		Handler: h,
		Value:   new(string),
	}
	flag.StringVar(f.Value, n, d, t)
	flags = append(flags, f)
}

// varsFileHandler initializes configuration with varsFile overriding env vars & defaults
func varsFileHandler(v *string) {
	err := config.Init(*v)
	if err != nil {
		log.Fatalf("[ERROR] Could not create config from provided filepath: %v", *v)
	} else if len(*v) > 0 {
		log.Printf("[NOTICE] Config read from file '%s', but may still be overridden by CLI flags.", *v)
	} else {
		log.Printf("[NOTICE] No configuration variables file specified. Using environment variabls and defaults only.")
	}
}

// outputDirHandler
func outputDirHandler(v *string) {
	if len(*v) > 0 {
		log.Printf("[NOTICE] Output Directory has been overridden via command line")
	}
}

// outputTypeHandler validates provided value and sets output accordingly
func outputTypeHandler(v *string) {
	if *v != "" {
		if *v == "IO" {
			probr.SetIOPaths("", config.Vars.OutputDir)
			log.Printf("[NOTICE] Probr results will be written to files in the specified output directory: %v", *v)
		} else if *v == "INMEM" {
			log.Printf("[NOTICE] Output type specified as INMEM: Results will not be handled by the CLI. Refer to the Summary Log for a results summary.")
		} else {
			log.Fatalf("[ERROR] Unknown output type specified: %s. Please use 'IO' or 'INMEM'", *v)
		}
		config.Vars.OutputType = *v
	}
}

func tagsHandler(v *string) {
	if len(*v) > 0 {
		config.Vars.Tests.Tags = *v
		log.Printf("[NOTICE] Tags have been added via command line.")
	}
	if len(config.Vars.Tests.Tags) == 0 {
		log.Printf("[NOTICE] No tags specified. All probes will be run.")
	}
}

func kubeConfigHandler(v *string) {
	if len(*v) > 0 {
		config.Vars.KubeConfigPath = *v
		log.Printf("[NOTICE] Kube Config has been overridden via command line")
	}
	if len(config.Vars.KubeConfigPath) == 0 {
		log.Printf("[NOTICE] No kubeconfig path specified. Falling back to default paths.")
	}
}
