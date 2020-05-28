package main

import (
	"flag"

	"gitlab.com/martinfleming/spa-server/logging"
)

const (
	defaultConfigPath = "/etc/spa-server/config.yaml"
)

// Configuration needs to be accessed at package level
var Configuration *Config

func main() {
	cfg, err := ReadConfig(parseArgs())
	if err != nil {
		logging.Error("Failed to read config file: %s", err)
		return
	}
	Configuration = cfg

	server := NewServer(Configuration.Port)
	server.Start()
	defer server.Stop()

	// Register healthcheck function
	if err != nil {
		logging.Fatal("Error parsing health check duration period from config: %s", err)
	}
}

// parseArgs gets the path to config file if supplied from commandline
// if not supplied, returns the default value
func parseArgs() string {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		logging.Debug("No user-supplied configuration file, using default")
		return defaultConfigPath
	}
	return args[0]
}
