package main

import (
	"os"
	"testing"

	"bou.ke/monkey"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&mainTestSuite{})

type mainTestSuite struct{}

func (*mainTestSuite) TestMain(c *C) {
	os.Args[1] = "no file"
	fakeExit := func(int) {
		panic("os.Exit called")
	}
	patch := monkey.Patch(os.Exit, fakeExit)

	defer patch.Unpatch()

	c.Check(main, Panics, "os.Exit called")

	os.Args[1] = "../../configs/config.example.yaml"

	c.Check(main, Not(Panics), "os.Exit called")
}

// TestDefaultConfig check with no args default config file path is set
func (*mainTestSuite) TestDefaultConfig(c *C) {
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, defaultConfigPath)
}

// TestParseArgs check with config file path is set with args
func (*mainTestSuite) TestParseArgs(c *C) {
	config := "../../configs/config.example.yaml"
	os.Args[1] = config
	configFilePath := parseArgs()
	c.Assert(configFilePath, Equals, config)
}

func (*mainTestSuite) TestValidConfig(c *C) {
	os.Args[1] = "../../configs/config.example.yaml"
	err := start()
	c.Assert(err, IsNil)
}

func (*mainTestSuite) TestInvalidConfig(c *C) {
	os.Args[1] = ""
	err := start()
	c.Assert(err, NotNil)
}
