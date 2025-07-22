package infrastructure

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain text password using bcrypt.
// Returns the hashed password as a string.
func HashPassword(passowrd string) (string, error) {
	// Generate hashed password with default cost
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedByte), nil
}

// ComparePassword compares a hashed password with a plain text password.
// Returns nil if they match, or an error if they do not.
func ComparePassword(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
