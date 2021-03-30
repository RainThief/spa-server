package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/RainThief/spa-server/pkg/logging"
)

var logger = logging.NewLogger(os.Stdout, os.Stderr, Config.LogLevel)

// ExpiresParsed type is processed config of what duration to add to each file type
type ExpiresParsed map[string]time.Duration

type expires map[string]string

func generateExpiresConfig() {
	parsedConfig := func(data expires) (config ExpiresParsed) {
		config = make(ExpiresParsed)
		for i, v := range data {
			duration, err := stringToDuration(strings.ToLower(v))
			if err != nil {
				logger.Error("Cannot set expires duration from %s, %s", data[i], err)
				continue
			}
			config[i] = duration
		}
		return
	}

	// create expires durations for global expires
	Config.ExpiresParsed = parsedConfig(Config.Expires)

	// create expires durations for individual sites
	for i, site := range Config.SitesAvailable {
		Config.SitesAvailable[i].ExpiresParsed = parsedConfig(site.Expires)
	}
}

func stringToDuration(interval string) (expiry time.Duration, err error) {
	re := regexp.MustCompile(`^([0-9]+) ([a-z]+?)s?$`)
	s := re.FindStringSubmatch(strings.ToLower(interval))

	if len(s) != 3 {
		err = fmt.Errorf("invalid expiry interval provided '%s'", interval)
		return
	}

	now := time.Now()

	increment, err := strconv.Atoi(s[1])
	if err != nil {
		return
	}

	switch s[2] {
	case "second":
		expiry = time.Duration(increment) * time.Second
	case "minute":
		expiry = time.Duration(increment) * time.Minute
	case "hour":
		expiry = time.Duration(increment) * time.Hour
	case "day":
		expiry = time.Until(now.AddDate(0, 0, increment))
	case "month":
		expiry = time.Until(now.AddDate(0, increment, 0))
	case "year":
		expiry = time.Until(now.AddDate(increment, 0, 0))
	}
	return
}
