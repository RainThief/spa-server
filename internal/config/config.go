package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var Config Configuration

type Site struct {
	StaticPath string `yaml:"path"`
	IndexFile  string `yaml:"index"`
	HostName   string `yaml:"host"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	Redirect   bool   `yaml:"redirectNonTLS"`
}

// Configuration is the configuration loaded from config.yaml
// @todo remove key and cert
type Configuration struct {
	LogLevel        string `yaml:"logLevel"`
	RequirementPath string `yaml:"requirementPath"`
	// @todo remove certs
	CertFile            string `yaml:"certFile"`
	KeyFile             string `yaml:"keyFile"`
	TLSPort             string `yaml:"TLSPort"`
	Port                string `yaml:"port"`
	AllowDirectoryIndex bool   `yaml:"allowDirectoryIndex"`
	SitesAvailable      []Site `yaml:"sitesAvailable"`
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

func IsTLSsite(site Site) bool {
	if strings.TrimSpace(site.CertFile) == "" ||
		strings.TrimSpace(site.KeyFile) == "" {
		return false
	}
	return true
}
