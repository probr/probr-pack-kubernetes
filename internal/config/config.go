package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/citihub/probr/internal/utils"
	"gopkg.in/yaml.v2"
)

// Vars is a singleton instance of ConfigVars
var Vars ConfigVars
var Spinner *spinner.Spinner

// GetTags returns Tags, prioritising command line parameter over vars file
func (ctx *ConfigVars) GetTags() string {
	if ctx.Tags == "" {
		ctx.handleTagExclusions() // only process tag exclusions from vars file if not supplied via the command line
	}
	return ctx.Tags
}

// Handle tag exclusions provided via the config vars file
func (ctx *ConfigVars) handleTagExclusions() {
	for _, tag := range ctx.TagExclusions {
		if ctx.Tags == "" {
			ctx.Tags = "~@" + tag
		} else {
			ctx.Tags = fmt.Sprintf("%s && ~@%s", ctx.Tags, tag)
		}
	}
}

// Init will override config.Vars with the content retrieved from a filepath
func Init(configPath string) error {
	log.Printf("[NOTICE] Initialized by %s", utils.CallerName(1))
	config, err := NewConfig(configPath)

	if err != nil {
		log.Printf("[ERROR] %v", err)
		return err
	}
	Vars = config
	setFromEnvOrDefaults(&Vars) // Set any values not retrieved from file

	SetLogFilter(Vars.LogLevel, os.Stderr) // Set the minimum log level obtained from Vars
	Vars.handleConfigFileExclusions()

	return nil
}

// NewConfig overrides the current config.Vars values
func NewConfig(c string) (ConfigVars, error) {
	// Create config structure
	config := ConfigVars{}
	if c == "" {
		return config, nil // No file path provided, return empty config
	}
	err := ValidateConfigPath(c)
	if err != nil {
		return config, err
	}
	// Open config file
	file, err := os.Open(c)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}

// ValidateConfigPath simply ensures the file exists
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

func LogConfigState() {
	s, _ := json.MarshalIndent(Vars, "", "  ")
	log.Printf("[NOTICE] Config State: %s", s)
}

func AuditDir() string {
	_ = os.Mkdir(Vars.AuditDir, 0755) // Creates if not already existing
	return Vars.AuditDir
}

func (ctx *ConfigVars) handleConfigFileExclusions() {
	if ctx.ServicePacks.Kubernetes.isExcluded() {
		ctx.addExclusion("probes/kubernetes")
	} else {
		ctx.handleProbeExclusions("kubernetes", ctx.ServicePacks.Kubernetes.Probes)
	}
	if ctx.ServicePacks.Storage.isExcluded() {
		ctx.addExclusion("probes/storage")
	} else {
		ctx.handleProbeExclusions("storage", ctx.ServicePacks.Storage.Probes)
	}
}

func (ctx *ConfigVars) handleProbeExclusions(packName string, probes []Probe) {
	for _, probe := range probes {
		if probe.isExcluded() {
			ctx.addExclusion(fmt.Sprintf("probes/%s/%s", packName, probe.Name))
		} else {
			for _, scenario := range probe.Scenarios {
				if scenario.isExcluded() {
					ctx.addExclusion(fmt.Sprintf("probes/%s/%s/%s", packName, probe.Name, scenario.Name))
				}
			}
		}
	}
}

func (ctx *ConfigVars) addExclusion(tag string) {
	if len(ctx.Tags) > 0 {
		ctx.Tags = ctx.Tags + " && "
	}
	ctx.Tags = fmt.Sprintf("%s~@%s", ctx.Tags, tag)
}

// Log and return exclusion configuration
func (k Kubernetes) isExcluded() bool {
	if k.Excluded != "" {
		log.Printf("[NOTICE] Excluding Kubernetes service pack. Justification: %s", k.Excluded)
		return true
	}
	log.Printf("[NOTICE] Kubernetes service pack included.")
	return false
}

// Log and return exclusion configuration
func (k Storage) isExcluded() bool {
	if k.Excluded != "" {
		log.Printf("[NOTICE] Excluding Storage service pack. Justification: %s", k.Excluded)
		return true
	}
	log.Printf("[NOTICE] Storage service pack included.")
	return false
}

// Log and return exclusion configuration
func (p Probe) isExcluded() bool {
	if p.Excluded != "" {
		log.Printf("[NOTICE] Excluding %s probe. Justification: %s", strings.Replace(p.Name, "_", " ", -1), p.Excluded)
		return true
	}
	return false
}

// Log and return exclusion configuration
func (s Scenario) isExcluded() bool {
	if s.Excluded != "" {
		log.Printf("[NOTICE] Excluding scenario '%s'. Justification: %s", s.Name, s.Excluded)
		return true
	}
	return false
}
