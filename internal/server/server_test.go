package server

import (
	"testing"

	"github.com/RainThief/spa-server/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

// func (suite *ServerTestSuite) SetupTest() {

// }

func (suite *ServerTestSuite) TestSetsUpServer() {
	_, _ = config.ReadConfig("../../configs/config.example.yaml")
	server := NewServer()
	// assert.Equal(suite.T(), server.sites[0].HostName, "{subdomain:(?:www.)?}localhost")
	assert.Len(suite.T(), server.sites, 3)
	assert.Len(suite.T(), server.tlsSites, 1)
}
