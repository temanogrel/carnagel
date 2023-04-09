package integrity

import (
	"testing"

	"git.misc.vee.bz/carnagel/minion/pkg"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	ecosystem_mocks "git.misc.vee.bz/carnagel/go-ecosystem/mocks"
)

type filesystemIntegritySuite struct {
	suite.Suite

	// dependencies
	consul *ecosystem_mocks.ConsulClient

	fileSystemIntegrity minion.FilesystemIntegrity
}

func (suite *filesystemIntegritySuite) SetupTest() {

	suite.consul = &ecosystem_mocks.ConsulClient{}

	app := &minion.Application{
		Consul: suite.consul,
	}

	suite.fileSystemIntegrity = &fileSystemIntegrity{
		app: app,
		log: logrus.New(),
	}
}

func (suite *filesystemIntegritySuite) TestIsRunning() {
	assert.Equal(suite.T(), false, suite.fileSystemIntegrity.IsRunning())
}

func (suite *filesystemIntegritySuite) TestIsScanning() {
	assert.Equal(suite.T(), false, suite.fileSystemIntegrity.IsScanning())
}

func TestNewFileSystemIntegrity(t *testing.T) {
	suite.Run(t, new(filesystemIntegritySuite))
}
