package infrastructure_test

import (
	"os"
	"testing"
	"time"

	"github.com/A2SVTask7/Infrastructure"
	"github.com/stretchr/testify/suite"
)

type JWTServiceSuite struct {
	suite.Suite
	service infrastructure.JWTService
	secret  string
}

func (s *JWTServiceSuite) SetupSuite() {
	s.secret = "testsecret123"
	os.Setenv("JWT_SECRET", s.secret)
	s.service = infrastructure.NewJWTService(s.secret)
}

func (s *JWTServiceSuite) TearDownSuite() {
	os.Unsetenv("JWT_SECRET")
}

// Test Generate creates a valid token string
func (s *JWTServiceSuite) TestGenerate() {
	claims := map[string]any{
		"username": "user1",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"sub":      "123",
	}
	tokenStr, err := s.service.Generate(claims)
	s.Require().NoError(err)
	s.NotEmpty(tokenStr)
}

// Test Validate parses and validates a good token
func (s *JWTServiceSuite) TestValidate_ValidToken() {
	claims := map[string]any{
		"username": "user1",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"sub":      "123",
	}
	tokenStr, err := s.service.Generate(claims)
	s.Require().NoError(err)

	parsedClaims, err := s.service.Validate(tokenStr)
	s.Require().NoError(err)
	s.Equal("user1", parsedClaims.Username)
	s.Equal("123", parsedClaims.Subject)
}

// Test Validate returns error on expired token
func (s *JWTServiceSuite) TestValidate_ExpiredToken() {
	claims := map[string]any{
		"username": "user1",
		"exp":      time.Now().Add(-time.Hour).Unix(), // expired
		"sub":      "123",
	}
	tokenStr, err := s.service.Generate(claims)
	s.Require().NoError(err)

	_, err = s.service.Validate(tokenStr)
	s.ErrorContains(err, "expired")
}

// Test Validate returns error on invalid token string
func (s *JWTServiceSuite) TestValidate_InvalidTokenString() {
	_, err := s.service.Validate("this-is-not-a-token")
	s.Error(err)
}

func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceSuite))
}
