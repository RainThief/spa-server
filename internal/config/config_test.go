package config

import (
	"testing"

	//revive:disable:dot-imports
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&configTestSuite{})

type configTestSuite struct{}

func (s *configTestSuite) TestConfigStructHasExpectedValues(c *C) {
	cfg, err := ReadConfig("../../configs/config.example.yaml")
	c.Assert(err, IsNil)
	c.Check(cfg.Port, Equals, "80")
	c.Check(cfg.TLSPort, Equals, "443")
	c.Check(cfg.Expires["png"], Equals, "2 days")
	c.Check(cfg.SitesAvailable[0].CertFile, Equals, "/etc/spa-server/certs/spa-server.pem")
	c.Check(cfg.SitesAvailable[0].Expires["html"], Equals, "2 months")
	c.Check(len(cfg.SitesAvailable), Equals, 3)
	c.Check(cfg.DisableHealthCheck, Equals, false)
}

func (s *configTestSuite) TestFailsWithNoConfigFileAvailable(c *C) {
	cfg, err := ReadConfig("invalidpath/not_a_file")
	c.Check(err, NotNil)
	c.Check(cfg, IsNil)

	cfg, err = ReadConfig("../../configs/config.invalid.yaml")
	c.Check(err, NotNil)
	c.Check(cfg, IsNil)
}

func (s *configTestSuite) TestDetectTLSsite(c *C) {
	c.Check(IsTLSsite(Site{}), Equals, false)

	c.Check(IsTLSsite(Site{
		CertFile: "CertFile",
	}), Equals, false)

	c.Check(IsTLSsite(Site{
		CertFile: "CertFile",
		KeyFile:  "KeyFile",
	}), Equals, true)
}
