package password

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// ValidatePasswordCriteria validates that a password meets the minimum security criteria:
// at least 8 characters, one uppercase letter, and one number.
func ValidatePasswordCriteria(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var hasNum, hasUpper bool
	for _, char := range password {
		if unicode.IsNumber(char) {
			hasNum = true
		}
		if unicode.IsUpper(char) {
			hasUpper = true
		}
	}

	if !hasNum || !hasUpper {
		return fmt.Errorf("password must contain at least one number and one uppercase letter")
	}

	return nil
}

// Hash returns the bcrypt hash of the given plaintext password.
func Hash(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashed), nil
}

// Compare checks a plaintext password against a bcrypt hash.
func Compare(password, hash string) error {
	if password == "" || hash == "" {
		return fmt.Errorf("password and hash cannot be empty")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
