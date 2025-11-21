package helpers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Role       string `json:"role"`
	Department string `json:"department,omitempty"`
	StudentID  string `json:"student_id"`
	Email      string `json:"email"`
	IsActive   bool   `json:"is_active"`
	UserID     string `json:"user_id"`

	jwt.StandardClaims
}

var (
	jwtKey []byte
	// Define common errors
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
	ErrMissingKey   = errors.New("JWT key not set")
)

// SetJWTKey sets the JWT secret key
func SetJWTKey(key string) error {
	if key == "" {
		return errors.New("JWT key cannot be empty")
	}
	jwtKey = []byte(key)
	log.Printf("JWT Key set successfully (length: %d bytes)", len(jwtKey))
	return nil
}

// GetJWTKey safely retrieves the JWT key
func GetJWTKey() ([]byte, error) {
	if len(jwtKey) == 0 {
		return nil, ErrMissingKey
	}
	return jwtKey, nil
}

// ValidateToken validates and parses a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token string is empty")
	}

	secretKey, err := GetJWTKey()
	if err != nil {
		return nil, err
	}

	// Parse the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		// Provide more specific error messages
		var validationErr *jwt.ValidationError
		if errors.As(err, &validationErr) {
			if validationErr.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrExpiredToken
			}
			if validationErr.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("invalid token signature")
			}
		}
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// Check if token is valid and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation for required fields
		if claims.Email == "" {
			return nil, errors.New("token missing required claims")
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// GenerateToken generates a new JWT token
func GenerateToken(email, userID, role string, additionalClaims map[string]any) (string, error) {
	secretKey, err := GetJWTKey()
	if err != nil {
		return "", err
	}

	// Token expiration time
	tokenExpiry := time.Now().Add(24 * time.Hour)

	// Create base claims
	claims := &Claims{
		Email:  email,
		Role:   role,
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "e-exam",
			Subject:   userID,
		},
	}

	// Add additional claims if provided
	if additionalClaims != nil {
		if firstName, ok := additionalClaims["first_name"].(string); ok {
			claims.FirstName = firstName
		}
		if lastName, ok := additionalClaims["last_name"].(string); ok {
			claims.LastName = lastName
		}
		if department, ok := additionalClaims["department"].(string); ok {
			claims.Department = department
		}
		if studentID, ok := additionalClaims["student_id"].(string); ok {
			claims.StudentID = studentID
		}
		if isActive, ok := additionalClaims["is_active"].(bool); ok {
			claims.IsActive = isActive
		}
	}

	// Create and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// GenerateTokens generates both access and refresh tokens (if needed)
func GenerateTokens(email, userID, role string, additionalClaims map[string]any) (accessToken string, err error) {
	return GenerateToken(email, userID, role, additionalClaims)
}

// RefreshToken generates a new token with extended expiration
func RefreshToken(oldToken string) (string, error) {
	claims, err := ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	// Generate new token with same claims but new expiration
	additionalClaims := map[string]any{
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
		"department": claims.Department,
		"student_id": claims.StudentID,
		"is_active":  claims.IsActive,
	}

	return GenerateToken(claims.Email, claims.UserID, claims.Role, additionalClaims)
}

func GenerateRandomKey() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal("Failed to generate key", err)
	}

	return base64.URLEncoding.EncodeToString(bytes)
}
