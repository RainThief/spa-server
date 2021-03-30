package logging

import (
	"os"

	"github.com/RainThief/spa-server/internal/config"
	"github.com/RainThief/spa-server/pkg/logging"
)

var Logger *logging.Logger

var cfg *config.Configuration = &config.Config

func init() {
	Logger = logging.NewLogger(os.Stdout, os.Stderr, "INFO")
}
