package domain

// JWTService defines the interface for generating and validating JWT tokens
type JWTService interface {
	Validate(tokenString string) (map[string]any, error) // Validates a token and returns claims
	Generate(claims map[string]any) (string, error)      // Generates a signed JWT string from claims
}

type IPasswordService interface {
	HashPassword(string) (string, error)  // Hashes a password
	ComparePassword(string, string) error // validates if two passwords are the same
}
