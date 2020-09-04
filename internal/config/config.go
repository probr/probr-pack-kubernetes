package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// ConfigVars contains all possible config vars. May be set by .yml, env, or defaults.
type ConfigVars struct {
	KubeConfigPath string `yaml:"kubeConfig"`
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
	Tests struct {
		Tags string `yaml:"tags"`
	} `yaml:"Tests"`
}

//Vars ...
var Vars ConfigVars

func init() {
	//create a defaulted config
	Init("")
}

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
func NewConfig(c string) (ConfigVars, error) {
	// Create config structure
	config := ConfigVars{}
	if c == "" {
		return config, nil
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
