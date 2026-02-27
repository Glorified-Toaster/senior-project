package helpers

import (
	"fmt"
	"time"

	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/domain"

	"github.com/golang-jwt/jwt"
)

type JWTAuth struct {
	cfg *config.Config
}

type Claims struct {
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	UserName string `json:"user_name"`
	IsActive bool   `json:"is_active"`
	UserID   string `json:"user_id"`

	jwt.StandardClaims
}

func NewJWT(cfg *config.Config) *JWTAuth {
	return &JWTAuth{
		cfg: cfg,
	}
}

func (j *JWTAuth) GetTokenFromConfig() (string, error) {
	secret := j.cfg.JWTAuth.Secret
	if secret == "" {
		return "", fmt.Errorf("failed to get secret from config file, please add it to the yaml config")
	}
	return secret, nil
}

func (j *JWTAuth) GenerateToken(user domain.User) (string, error) {

	if j == nil {
		return "", fmt.Errorf("jwt auth is not initialized")
	}

	if j.cfg == nil {
		return "", fmt.Errorf("jwt auth config is not initialized")
	}

	tokenExpiry := time.Now().Add(time.Hour * 24)

	claims := Claims{
		FullName: user.FullName,
		Role:     string(user.Role),
		UserName: user.Username,
		IsActive: user.IsActive,
		UserID:   user.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiry.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "uot-exam",
		},
	}

	secret, err := j.GetTokenFromConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	if j == nil {
		return nil, fmt.Errorf("jwt auth is not initialized")
	}

	if j.cfg == nil {
		return nil, fmt.Errorf("jwt auth config is not initialized")
	}

	if tokenString == "" {
		return nil, fmt.Errorf("token string is empty")
	}

	secret, err := j.GetTokenFromConfig()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token")
	}
}
