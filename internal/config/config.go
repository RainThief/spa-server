package config

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var Config Configuration

// Configuration is the configuration loaded from config.yaml
type Configuration struct {
	LogLevel            string        `yaml:"logLevel"`
	RequirementPath     string        `yaml:"requirementPath"`
	HealthCheckPeriod   time.Duration `yaml:"healthCheckPeriod"`
	CertFile            string        `yaml:"certFile"`
	KeyFile             string        `yaml:"keyFile"`
	Port                string        `yaml:"port"`
	AllowDirectoryIndex bool          `yaml:"allowDirectoryIndex"`
}

// ReadConfig reads the config from the file provided and parses it as Yaml
// returning a Config object if parsed successfully.
func ReadConfig(filePath string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	return &Config, yaml.Unmarshal(data, &Config)
}
