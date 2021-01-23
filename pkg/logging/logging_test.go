package logging

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var logger = Logger{os.Stdout, os.Stderr}

var STDBuf bytes.Buffer
var ERRBuf bytes.Buffer

type ExampleTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
}

func (suite *ExampleTestSuite) SetupTest() {
	STDBuf = bytes.Buffer{}
	ERRBuf = bytes.Buffer{}
	logger.SetSTDOutput(&STDBuf)
	logger.SetSTDError(&ERRBuf)
}

func (suite *ExampleTestSuite) TestLogError() {
	logger.Error("test log")
	assert.Equal(suite.T(), "ERROR: test log\n", ERRBuf.String())

}

func (suite *ExampleTestSuite) TestLogFatal() {
	logger.Fatal("test log")
	assert.Equal(suite.T(), "FATAL: test log\n", ERRBuf.String())
}

func (suite *ExampleTestSuite) TestLogDebug() {
	logger.Debug("test log")
	assert.Equal(suite.T(), "DEBUG: test log\n", STDBuf.String())
}

func (suite *ExampleTestSuite) TestLogInfo() {
	logger.Info("test log")
	assert.Equal(suite.T(), "INFO: test log\n", STDBuf.String())
}

func (suite *ExampleTestSuite) TestLogAndRaiseError() {
	err := logger.LogAndRaiseError("test log")
	assert.Equal(suite.T(), "ERROR: test log\n", ERRBuf.String())
	assert.NotNil(suite.T(), err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}
