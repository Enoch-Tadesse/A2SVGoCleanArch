package infrastructure

import (
	domain "github.com/A2SVTask7/Domain"
	"golang.org/x/crypto/bcrypt"
)

type PasswordService struct{}

func NewPasswordService() domain.IPasswordService {
	return &PasswordService{}
}

// HashPassword hashes a plain text password using bcrypt.
// Returns the hashed password as a string.
func (ps PasswordService) HashPassword(passowrd string) (string, error) {
	// Generate hashed password with default cost
	hashedByte, err := bcrypt.GenerateFromPassword([]byte(passowrd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedByte), nil
}

// ComparePassword compares a hashed password with a plain text password.
// Returns nil if they match, or an error if they do not.
func (ps PasswordService) ComparePassword(hashed string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}
