package infrastructure

import (
	"fmt"
	"os"
	"time"

	domain "github.com/A2SVTask7/Domain"
	"github.com/golang-jwt/jwt/v5"
)

// customClaims represents the custom payload we embed inside JWT tokens
type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// jwtService is a concrete implementation of JWTService
type jwtService struct {
	secret []byte // Secret key used for signing and verifying JWTs
}

// NewJWTService creates a new instance of jwtService using the provided secret
func NewJWTService(secret string) domain.JWTService {
	return &jwtService{
		secret: []byte(secret),
	}
}

// Generate creates a JWT token string from a map of claims
func (js *jwtService) Generate(claims map[string]any) (string, error) {
	mapClaims := jwt.MapClaims{}
	for k, v := range claims {
		mapClaims[k] = v
	}
	// Create a new token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)

	// Sign and return the token string
	tokenString, err := token.SignedString(js.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// Validate parses and verifies a JWT token string, returning its claims if valid
func (js *jwtService) Validate(tokenString string) (map[string]any, error) {
	// Load the JWT secret from environment variables
	jwt_secret := []byte(os.Getenv("JWT_SECRET"))

	// Parse the token with expected signing method and custom claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		return jwt_secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	// Assert the claims and check token validity
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims: %w", err)
	}

	// Check for expiration
	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}
	return map[string]any{
		"sub":      claims.Subject,
		"username": claims.Username,
		"exp":      claims.ExpiresAt.Time.Unix(),
	}, err
}
