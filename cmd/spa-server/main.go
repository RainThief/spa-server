package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/martinfleming/spa-server/internal/config"
	"gitlab.com/martinfleming/spa-server/internal/logging"
	"gitlab.com/martinfleming/spa-server/internal/server"
)

const (
	defaultConfigPath = "/etc/spa-server/config.yaml"
)

func main() {
	_, err := config.ReadConfig(parseArgs())
	if err != nil {
		logging.Error("Failed to read config file: %s", err)
		return
	}
	server := server.NewServer()
	defer server.Stop()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		server.Start()
	}()
	<-sigint
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
	logging.Debug("Using user-supplied configuration file %s", args[0])
	return args[0]
}
