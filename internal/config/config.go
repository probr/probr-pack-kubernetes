package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	_ "gopkg.in/yaml.v2"
)

// Config contains all possible config vars. May be set by .yml, env, or defaults.
type Config struct {
	KubeConfigPath string
	OutputType     string `yaml:"outputType"`
	Images         struct {
		Repository string `yaml:"repository"`
		Curl       string `yaml:"curl"`
		BusyBox    string `yaml:"busyBox"`
	} `yaml:"Images"`
	Azure struct {
		SubscriptionID  string `yaml:"subscriptionID"`
		ClientID        string `yaml:"clientID"`
		ClientSecret    string `yaml:"clientSecret"`
		TenantID        string `yaml:"tenantID"`
		LocationDefault string `yaml:"locationDefault"`
	} `yaml:"azure"`
}

var Vars Config

// Init will override Vars when it is used
func Init(configPath string) error {
	config, err := NewConfig(configPath)
	if err != nil {
		return err
	}
	Vars = config
	getEnvOrDefaults(&Vars) // Set any values not retrieved from file
	return nil
}

// NewConfig can be used multiple times, if the need arises
func NewConfig(c string) (Config, error) {
	// Create config structure
	config := Config{}
	err := ValidateConfigPath(c)
	if err != nil {
		if c != "./config.yml" {
			// If config path isn't at the default value, panic on failure
			return config, err
		} else {
			log.Printf("[NOTICE] Config vars not found at default filepath. Continuing with env vars and/or defaults.")
			return config, nil
		}
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
