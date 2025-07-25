package infrastructure_test

import (
	"os"
	"testing"
	"time"

	"github.com/A2SVTask7/Infrastructure"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
	envVars map[string]string
}

func (s *ConfigSuite) SetupSuite() {
	// Define all env vars and their values once
	s.envVars = map[string]string{
		"MONGO_URI":       "mongodb://customhost:12345",
		"COLLECTION_TASK": "mytasks",
		"COLLECTION_USER": "myusers",
		"JWT_SECRET":      "myjwtsecret",
		"DBName":          "mydatabase",
		"Port":            "9090",
		"APP_TIMEOUT":     "10s",
	}

	// Set all env vars
	for k, v := range s.envVars {
		os.Setenv(k, v)
	}
}

func (s *ConfigSuite) TearDownSuite() {
	// Unset all env vars
	for k := range s.envVars {
		os.Unsetenv(k)
	}
}

func (s *ConfigSuite) TestLoadConfig() {
	infrastructure.LoadConfig()
	cfg := infrastructure.AppConfig

	s.Equal(s.envVars["MONGO_URI"], cfg.MongoURI)
	s.Equal(s.envVars["COLLECTION_TASK"], cfg.CollectionTask)
	s.Equal(s.envVars["COLLECTION_USER"], cfg.CollectionUser)
	s.Equal(s.envVars["JWT_SECRET"], cfg.JWTSecret)
	s.Equal(s.envVars["DBName"], cfg.DBName)
	s.Equal(s.envVars["Port"], cfg.Port)
	s.Equal(10*time.Second, cfg.Timeout)
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
