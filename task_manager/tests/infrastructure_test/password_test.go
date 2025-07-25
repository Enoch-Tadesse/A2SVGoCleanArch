package infrastructure_test

import (
	"testing"

	"github.com/A2SVTask7/Infrastructure"
	"github.com/stretchr/testify/suite"
)

type PasswordTestSuite struct {
	suite.Suite
}

func (suite *PasswordTestSuite) TestHashPasswordAndCompare_Success() {
	password := "my_secure_password"

	hashed, err := infrastructure.HashPassword(password)
	suite.NoError(err)
	suite.NotEmpty(hashed)

	// Should successfully compare correct password
	err = infrastructure.ComparePassword(hashed, password)
	suite.NoError(err)

	// Should fail for incorrect password
	err = infrastructure.ComparePassword(hashed, "wrong_password")
	suite.Error(err)
}

func (suite *PasswordTestSuite) TestHashPassword_ErrorOnEmpty() {
	hashed, err := infrastructure.HashPassword("secret")
	suite.NoError(err)
	suite.NotEmpty(hashed)

	err = infrastructure.ComparePassword(hashed, "secret")
	suite.NoError(err)
}

func TestPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}
