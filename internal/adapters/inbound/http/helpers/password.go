package helpers

import (
	"uot-exam/internal/utils/password"
)

// ValidatePasswordCriteria validates the password against minimum security criteria.
// Delegates to the shared internal/utils/password package.
func ValidatePasswordCriteria(p string) error {
	return password.ValidatePasswordCriteria(p)
}

// HashPassword returns the bcrypt hash of the given plaintext password.
// Delegates to the shared internal/utils/password package.
func HashPassword(p string) (string, error) {
	return password.Hash(p)
}

// CheckWithHashedPassword verifies a plaintext password against a bcrypt hash.
// Delegates to the shared internal/utils/password package.
func CheckWithHashedPassword(plain, hash string) error {
	return password.Compare(plain, hash)
}
