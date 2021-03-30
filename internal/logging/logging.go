package logging

import (
	"os"

	"github.com/RainThief/spa-server/pkg/logging"
)

var Logger *logging.Logger

func init() {
	Logger = logging.NewLogger(os.Stdout, os.Stderr, "INFO")
}
