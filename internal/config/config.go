package config

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config global app configuration
var Config Configuration

// type ExpiresList map[string]string

// Site holds configuration for individual spa sites
type Site struct {
	StaticPath    string  `yaml:"path"`
	IndexFile     string  `yaml:"index"`
	HostName      string  `yaml:"host"`
	CertFile      string  `yaml:"certFile"`
	KeyFile       string  `yaml:"keyFile"`
	Redirect      bool    `yaml:"redirectNonTLS"`
	Compress      bool    `yaml:"compress"`
	CompressLevel int     `yaml:"compressionLevel"`
	Expires       expires `yaml:"expires"`
	ExpiresParsed ExpiresParsed
}

// Configuration is the configuration loaded from config.yaml
type Configuration struct {
	LogLevel            string  `yaml:"logLevel"`
	RequirementPath     string  `yaml:"requirementPath"`
	TLSPort             string  `yaml:"TLSPort"`
	Port                string  `yaml:"port"`
	AllowDirectoryIndex bool    `yaml:"allowDirectoryIndex"`
	SitesAvailable      []Site  `yaml:"sitesAvailable"`
	DisableHealthCheck  bool    `yaml:"disableHealthCheck"`
	HealthCheckPort     int     `yaml:"healthCheckPort"`
	Compress            bool    `yaml:"compress"`
	CompressLevel       int     `yaml:"compressionLevel"`
	Expires             expires `yaml:"expires"`
	ExpiresParsed       ExpiresParsed
}

// ReadConfig reads the config from the file provided and parses it as Yaml
// returning a Config object if parsed successfully.
func ReadConfig(filePath string) (*Configuration, error) {
	data, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &Config); err != nil {
		return &Config, err
	}
	generateExpiresConfig()

	return &Config, nil
}

// IsTLSsite inspects site config to see if can be served under ssl
func IsTLSsite(site Site) bool {
	// @todo check cert file exists in location, check is cert file or key file
	if strings.TrimSpace(site.CertFile) == "" ||
		strings.TrimSpace(site.KeyFile) == "" {
		return false
	}
	return true
}
