package helpers

import (
	"fmt"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// ValidatePassword : validate the password if it's met the criteria
func ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	num := false
	upperCase := false

	for _, char := range password {
		if unicode.IsNumber(char) {
			num = true
		}
		if unicode.IsUpper(char) {
			upperCase = true
		}
	}

	if !num || !upperCase {
		return fmt.Errorf("password must contain at least one number and one uppercase letter")
	}

	return nil
}

// HashPassword : return the hashed password and error
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckWithHashedPassword(password, hashedPassword string) error {
	if password == "" || hashedPassword == "" {
		return fmt.Errorf("password and hash cannot be empty")
	}

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	return nil
}
