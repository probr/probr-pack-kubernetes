package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/citihub/probr-pack-kubernetes/internal/config"
	"github.com/citihub/probr-pack-kubernetes/internal/summary"
	"github.com/citihub/probr-pack-kubernetes/pack"

	sdk "github.com/citihub/probr-sdk"
	audit "github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/logging"
	"github.com/citihub/probr-sdk/plugin"
	probeengine "github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/utils"
)

var (
	// ServicePackName is the name for the service pack
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

	versionCmd, runCmd := setFlags()
	handleCommands(versionCmd, runCmd)
}

func setFlags() (versionCmd, runCmd *flag.FlagSet) {
	// > probr version [-v]
	versionCmd = flag.NewFlagSet("version", flag.ExitOnError)
	config.Vars.Verbose = *versionCmd.Bool("v", false, "Display extended version information") // TODO: Harness '-v' in the standard probr execution

	// > probr
	runCmd = flag.NewFlagSet("run", flag.ExitOnError)
	runCmd.StringVar(&config.Vars.VarsFile, "varsfile", "", "path to config file")
	runCmd.StringVar(&config.Vars.WriteDirectory, "writedirectory", "", "output directory")
	runCmd.StringVar(&config.Vars.LogLevel, "loglevel", "", "set log level")
	runCmd.StringVar(&config.Vars.ResultsFormat, "resultsformat", "", "set the bdd results format (default = cucumber)")
	runCmd.StringVar(&config.Vars.Tags, "tags", "", "feature tags to include or exclude")
	runCmd.StringVar(&config.Vars.KubeConfigPath, "kubeconfig", "", "kube config file")
	return
}

func handleCommands(versionCmd, runCmd *flag.FlagSet) {
	subCommand := ""
	if len(os.Args) > 1 {
		subCommand = os.Args[1]
	}
	switch subCommand {
	case "version":
		versionCmd.Parse(os.Args[2:])
		printVersion(os.Stdout)

	case "debug": // Same cli args as run. Use this to bypass plugin and execute directly for debugging
		// Parse cli args
		runCmd.Parse(os.Args[2:]) // Skip first arg as it will be 'debug'
		log.Printf("[ERROR] ---->>> %s", config.Vars.VarsFile)
		ProbrCoreLogic()

	default:
		// Parse cli args
		runCmd.Parse(os.Args[1:])
		log.Printf("[ERROR] ---->>> %s", config.Vars.VarsFile)

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
		defer sdk.GlobalConfig.CleanupTmp()
		// TODO: Additional cleanup may be needed. For instance, any pods created during tests are not being dropped if aborted.
		os.Exit(0)
	}()
}

// ProbrCoreLogic ...
func ProbrCoreLogic() (err error) {
	defer sdk.GlobalConfig.CleanupTmp()
	setupCloseHandler() // Sigterm protection

	summary.State = audit.NewSummaryState("kubernetes")

	config.Vars.Init()
	config.Vars.LogConfigState() // TODO: Update this func to accept a generic object, so that global and local settings can be logged (if needed)

	logWriter := logging.ProbrLoggerOutput()
	log.SetOutput(logWriter) // TODO: This is a temporary patch, since logger output is being overritten while loading config vars

	store := probeengine.NewProbeStore("kubernetes", config.Vars.Tags, &summary.State)
	s, err := store.RunAllProbes(pack.GetProbes())
	if err != nil {
		log.Printf("[ERROR] Error executing tests %v", err)
		return
	}
	log.Printf("[INFO] Overall test completion status: %v", s)
	summary.State.SetProbrStatus()

	_, success := probeengine.GetAllProbeResults(store) // TODO: This is returning success=true despite failing probes.
	summary.State.PrintSummary()
	summary.State.WriteSummary()

	if !success || summary.State.ProbesFailed > 0 { //Adding this until 'success' can be fixed. See above TODO
		return utils.ReformatError("One or more probe scenarios were not successful. View the output logs for more details.")
	}
	return
}

func printVersion(w io.Writer) {

	if config.Vars.Verbose {
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
