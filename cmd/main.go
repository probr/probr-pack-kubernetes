package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	pack "github.com/citihub/probr-pack-kubernetes"
	"github.com/citihub/probr-sdk/audit"
	cliflags "github.com/citihub/probr-sdk/cli_flags"
	"github.com/citihub/probr-sdk/config"
	"github.com/citihub/probr-sdk/logging"
	"github.com/citihub/probr-sdk/plugin"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/utils"
)

var (
	// ServicePackName is the display name for the service pack
	ServicePackName = "Kubernetes" // TODO: Return binary name instead?

	// Version is the main version number that is being run at the moment
	Version = "0.0.2"

	// Prerelease is a marker for the version. If this is "" (empty string)
	// then it means that it is a final release. Otherwise, this is a pre-release
	// such as "dev" (in development), "beta", "rc", etc.
	// This should only be modified thru ldflags in make file. See 'make release'
	Prerelease = "dev"

	// GitCommitHash references the commit id at build time
	// This should only be modified thru ldflags in make file. See 'make release'
	GitCommitHash = ""

	// BuiltAt is the build date
	// This should only be modified thru ldflags in make file. See 'make release'
	BuiltAt = ""
)

// ServicePack ...
type ServicePack struct {
}

// RunProbes ...
func (sp *ServicePack) RunProbes() error {
	log.Printf("[DEBUG] message from ServicePack_Probr.RunProbes")
	log.Printf("[DEBUG] args... %v", os.Args)

	return ProbrCoreLogic()
}

func main() {

	// > probr version [-v]
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	verboseVersionFlag := versionCmd.Bool("v", false, "Display extended version information")

	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	case "version":
		versionCmd.Parse(os.Args[2:])
		printVersion(os.Stdout, *verboseVersionFlag)

	default:
		if len(os.Args) > 1 && os.Args[1] == "debug" { // TODO: Check this logic
			ProbrCoreLogic()
			return
		}

		// Serve plugin
		spProbr := &ServicePack{}
		serveOpts := &plugin.ServeOpts{
			Pack: spProbr,
		}

		plugin.Serve(serveOpts)
	}
}

// setupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
// Ref: https://golangcode.com/handle-ctrl-c-exit-in-terminal/
func setupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Execution aborted - %v", "SIGTERM")
		probeengine.CleanupTmp()
		// TODO: Additional cleanup may be needed. For instance, any pods created during tests are not being dropped if aborted.
		os.Exit(0)
	}()
}

// ProbrCoreLogic ...
func ProbrCoreLogic() (err error) {
	log.Printf("[INFO] message from ProbCoreLogic: %s", "Start")
	defer probeengine.CleanupTmp()
	setupCloseHandler()

	// Setup for handling SIGTERM (Ctrl+C)
	//setupCloseHandler()

	err = config.Init("") // Create default config
	if err != nil {
		log.Printf("[ERROR] error returned from config.Init: %v", err)
		return
	}
	if len(os.Args[1:]) > 0 {
		log.Printf("[DEBUG] Checking for CLI options or flags")
		cliflags.HandleRequestForRequiredVars()
		log.Printf("[DEBUG] Handle pack option")
		cliflags.HandlePackOption()
		parseFlags()
	}

	config.Vars.LogConfigState()

	logWriter := logging.ProbrLoggerOutput()
	log.SetOutput(logWriter) // TODO: This is a temporary patch, since logger output is being overritten while loading config vars

	s, ts, err := probeengine.RunAllProbes("kubernetes", pack.GetProbes())
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		return
	}
	log.Printf("[INFO] Overall test completion status: %v", s)
	audit.State.SetProbrStatus()

	_, success := probeengine.GetAllProbeResults(ts) // TODO: Use the results provided here
	audit.State.PrintSummary()
	audit.State.WriteSummary()

	log.Printf("[INFO] message from ProbCoreLogic: %s", "End")

	if !success {
		return utils.ReformatError("One or more probe scenarios were not successful. View the output logs for more details.")
	}
	return
}

func parseFlags() {
	var flags cliflags.Flags

	flags.NewStringFlag("varsfile", "path to config file", cliflags.VarsFileHandler)
	flags.NewStringFlag("writedirectory", "output directory", cliflags.WriteDirHandler)
	flags.NewStringFlag("loglevel", "set log level", cliflags.LoglevelHandler)
	flags.NewStringFlag("resultsformat", "set the bdd results format (default = cucumber)", cliflags.ResultsformatHandler)
	flags.NewStringFlag("tags", "feature tags to include or exclude", cliflags.TagsHandler)

	flags.NewStringFlag("kubeconfig", "kube config file", kubeConfigHandler)

	flags.ExecuteHandlers()

}

func kubeConfigHandler(v *string) {
	value := *v
	if len(value) > 0 {
		config.Vars.ServicePacks.Kubernetes.KubeConfigPath = value
		log.Printf("[NOTICE] Kubeconfig path has been overridden via command line")
	}
	if len(config.Vars.ServicePacks.Kubernetes.KubeConfigPath) == 0 {
		log.Printf("[NOTICE] No kubeconfig path specified. Falling back to default paths.")
	}
}

func printVersion(w io.Writer, verbose bool) {

	if verbose {
		fmt.Fprintf(w, "Service Pack : %s", ServicePackName)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Version      : %s", getVersion())
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Commit       : %s", GitCommitHash)
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Built at     : %s", BuiltAt)
	} else {
		fmt.Fprintf(w, "Version      : %s", getVersion())
	}

}

func getVersion() string {
	if Prerelease != "" {
		return fmt.Sprintf("%s-%s", Version, Prerelease)
	}
	return Version
}
