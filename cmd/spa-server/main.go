package main

import (
	"flag"

	"gitlab.com/martinfleming/spa-server/internal/logging"
)

const (
	defaultConfigPath = "/etc/spa-server/config.yaml"
)

// Configuration needs to be accessed at global level
// @todo config to own internal package
var config *Configuration

func main() {
	cfg, err := ReadConfig(parseArgs())
	if err != nil {
		logging.Error("Failed to read config file: %s", err)
		return
	}
	// fmt.Println(cfg)
	config = cfg
	// fmt.Println(cfg)
	// fmt.Println(Configuration)

	server := NewServer(config.Port)
	server.Start()
	defer server.Stop()
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
