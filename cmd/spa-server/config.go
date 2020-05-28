package main

import (
	"io/ioutil"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

// Configuration is the configuration loaded from config.yaml
type Configuration struct {
	LogLevel          string        `yaml:"logLevel"`
	RequirementPath   string        `yaml:"requirementPath"`
	HealthCheckPeriod time.Duration `yaml:"healthCheckPeriod"`
	CertFile          string        `yaml:"certFile"`
	KeyFile           string        `yaml:"keyFile"`
	Port              string        `yaml:"port"`
}

// ReadConfig reads the config from the file provided and parses it as Yaml
// returning a Config object if parsed successfully.
func ReadConfig(filePath string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	config := Configuration{}
	return &config, yaml.Unmarshal(data, &config)
}
