package logging

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var log *Logger

var STDBuf bytes.Buffer
var ERRBuf bytes.Buffer

type LoggingTestSuite struct {
	suite.Suite
}

func (suite *LoggingTestSuite) SetupTest() {
	log = NewLogger(os.Stdout, os.Stderr, "INFO")
	STDBuf = bytes.Buffer{}
	ERRBuf = bytes.Buffer{}
	log.SetSTDOutput(&STDBuf)
	log.SetSTDError(&ERRBuf)
}

func (suite *LoggingTestSuite) TestLogError() {
	log.Error("test log")
	assert.Equal(suite.T(), "ERROR: test log\n", ERRBuf.String())
}

func (suite *LoggingTestSuite) TestLogFatal() {
	log.Fatal("test log")
	assert.Equal(suite.T(), "FATAL: test log\n", ERRBuf.String())
}

func (suite *LoggingTestSuite) TestLogDebug() {
	log.Debug("test log")
	assert.Equal(suite.T(), "DEBUG: test log\n", STDBuf.String())
}

func (suite *LoggingTestSuite) TestLogInfo() {
	log.Info("test log")
	assert.Equal(suite.T(), "INFO: test log\n", STDBuf.String())
}

func (suite *LoggingTestSuite) TestLogAndRaiseError() {
	err := log.LogAndRaiseError("test log")
	assert.Equal(suite.T(), "ERROR: test log\n", ERRBuf.String())
	assert.NotNil(suite.T(), err)
}

func (suite *LoggingTestSuite) TestFatalLogOnly() {
	log = NewLogger(&STDBuf, &ERRBuf, "FATAL")
	log.Info("%s", "info event")
	assert.Equal(suite.T(), "", STDBuf.String(), "")
	log.Error("%s", "error event")
	assert.Equal(suite.T(), "", ERRBuf.String(), "")
	log.Fatal("%s", "critical event")
	assert.Equal(suite.T(), "FATAL: critical event\n", ERRBuf.String())
}

func TestLoggingTestSuite(t *testing.T) {
	suite.Run(t, new(LoggingTestSuite))
}
