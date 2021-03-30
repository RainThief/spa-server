package main

import (
	"flag"
	"fmt"
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

var logger = logging.Logger

func main() {
	if err := start(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func start() error {
	if _, err := config.ReadConfig(parseArgs()); err != nil {
		return fmt.Errorf("Failed to read config file: %s", err)
	}
	httpServer := server.NewServer()
	defer httpServer.Stop()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		httpServer.Start(sigint)
	}()
	<-sigint
	return nil
}

// parseArgs gets the path to config file if supplied from commandline
// if not supplied, returns the default value
func parseArgs() string {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		logger.Debug("No user-supplied configuration file, using default")
		return defaultConfigPath
	}
	logger.Debug("Using user-supplied configuration file %s", args[0])
	return args[0]
}
