package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// Error constants
var (
	ErrInvalidPassword = errors.New("invalid password")
)

// PasswordManager handles password hashing and verification
type PasswordManager struct{}

// NewPasswordManager creates a new PasswordManager
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword hashes a password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	// Generate a salt with a cost of 12
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if a password matches a hash
func (pm *PasswordManager) CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}