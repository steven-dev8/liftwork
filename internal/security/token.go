package security

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrJWTSecretRequired = errors.New("JWT secret is required")
	ErrInvalidTokenTTL   = errors.New("token TTL must be greater than zero")
	ErrJWTRequired       = errors.New("JWT is required")
	ErrInvalidJWT        = errors.New("invalid JWT")
	ErrInvalidUserID     = errors.New("invalid user ID")
)

func CreateToken(userID int64, secret string, ttl time.Duration) (string, error) {
	if strings.TrimSpace(secret) == "" {
		return "", ErrJWTSecretRequired
	}
	if ttl <= 0 {
		return "", ErrInvalidTokenTTL
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign JWT: %w", err)
	}

	return signedToken, nil
}

func ValidateToken(tokenString string, secret string) (int64, error) {
	if strings.TrimSpace(secret) == "" {
		return 0, ErrJWTSecretRequired
	}

	if strings.TrimSpace(tokenString) == "" {
		return 0, ErrJWTRequired
	}

	claims := new(jwt.RegisteredClaims)

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return 0, fmt.Errorf("validate JWT: %w", err)
	}

	if !token.Valid {
		return 0, ErrInvalidJWT
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, ErrInvalidJWT
	}

	return userID, nil
}