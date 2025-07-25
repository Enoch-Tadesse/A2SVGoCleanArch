package mocks

import (
	infrastructure "github.com/A2SVTask7/Infrastructure"
	"github.com/stretchr/testify/mock"
)

type MockJWTService struct {
	mock.Mock
}

// Validate(tokenString string) (*customClaims, error) // Validates a token and returns claims
// Generate(claims map[string]any) (string, error)     // Generates a signed JWT string from claims

func (m *MockJWTService) Validate(token string) (*infrastructure.CustomClaims, error) {
	args := m.Called(token)

	// try to assert to a pointer
	if claims, ok := args.Get(0).(*infrastructure.CustomClaims); ok {
		return claims, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockJWTService) Generate(claims map[string]any) (string, error) {
	args := m.Called(claims)
	return args.String(0), args.Error(0)
}
