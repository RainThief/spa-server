package main

import (
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&mainTestSuite{})

type mainTestSuite struct{}

// TestDefaultConfig check with no args default config file path is set
func (*mainTestSuite) TestDefaultConfig(c *C) {
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, defaultConfigPath)
}

// TestParseArgs check with config file path is set with args
func (*mainTestSuite) TestParseArgs(c *C) {
	config := "../../test/testdata/httpOnlySites.yml"
	os.Args[1] = config
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, config)
}

func (*mainTestSuite) TestValidConfig(c *C) {
	os.Args[1] = "../../test/testdata/httpOnlySites.yml"
	err := start()
	c.Assert(err, IsNil)
}

func (*mainTestSuite) TestInvalidConfig(c *C) {
	os.Args[1] = ""
	err := start()
	c.Assert(err, NotNil)
}
