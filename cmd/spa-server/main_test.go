package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&mainTestSuite{})

type mainTestSuite struct {
	suite.Suite
}

// TestDefaultConfig check with no args default config file path is set
func (*mainTestSuite) TestDefaultConfig(c *C) {
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, defaultConfigPath)
}

// TestParseArgs check with config file path is set with args
func (*mainTestSuite) TestParseArgs(c *C) {
	config := "my.conf"
	os.Args[1] = config
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, config)
}
