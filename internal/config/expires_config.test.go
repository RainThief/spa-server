package config

import (
	"time"

	. "gopkg.in/check.v1"
)

var _ = Suite(&expiresConfigTestSuite{})

type expiresConfigTestSuite struct{}

var cfg *Configuration

func (s *expiresConfigTestSuite) SetUpTest(c *C) {
	cfg, _ = ReadConfig("../../configs/config.example.yaml")
}

func (s *expiresConfigTestSuite) TestConfigStructHasExpectedExpiresValues(c *C) {
	expires := cfg.ExpiresParsed["png"]
	c.Check(expires > (48*time.Hour)-time.Second, Equals, true)
	c.Check(expires < 48*time.Hour, Equals, true)
}

func (s *expiresConfigTestSuite) TestConfigSiteStructHasExpectedValues(c *C) {
	expires := cfg.SitesAvailable[0].ExpiresParsed["png"]
	c.Check(expires > (24*time.Hour)-time.Second, Equals, true)
	c.Check(expires < 24*time.Hour, Equals, true)
}

func (s *expiresConfigTestSuite) TestStringsCanBeConvertedToDuration(c *C) {
	second, _ := stringToDuration("1 second")
	c.Check(second, Equals, 1*time.Second)
	minute, _ := stringToDuration("1 minute")
	c.Check(minute, Equals, 1*time.Minute)
	hour, _ := stringToDuration("1 hour")
	c.Check(hour, Equals, 1*time.Hour)
	day, _ := stringToDuration("1 day")
	c.Check(day.Round(time.Minute), Equals, 24*time.Hour)
	month, _ := stringToDuration("1 month")
	c.Check(
		time.Until(time.Now().AddDate(0, 1, 0)).Round(time.Minute),
		Equals,
		month.Round(time.Minute),
	)
	year, _ := stringToDuration("1 year")
	c.Check(year.Round(time.Minute), Equals, (365 * (24 * time.Hour)))
}

func (s *expiresConfigTestSuite) TestPluralStringsCanBeConvertedToDuration(c *C) {
	second, _ := stringToDuration("1 seconds")
	c.Check(second, Equals, 1*time.Second)
	minute, _ := stringToDuration("2 minutes")
	c.Check(minute, Equals, 2*time.Minute)
}

func (s *expiresConfigTestSuite) TestCapitalStringsCanBeConvertedToDuration(c *C) {
	second, _ := stringToDuration("1 SECOND")
	c.Check(second, Equals, 1*time.Second)
	minute, _ := stringToDuration("2 SECONDS")
	c.Check(minute, Equals, 2*time.Second)
}
