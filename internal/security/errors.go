package security

import "errors"

var (
	// JWT
	ErrJWTSecretRequired = errors.New("jwt secret is required")
	ErrInvalidTokenTTL   = errors.New("token TTL must be greater than zero")
	ErrJWTRequired       = errors.New("jwt is required")
	ErrInvalidJWT        = errors.New("invalid jwt")
	ErrInvalidUserID     = errors.New("invalid user ID")

	// Password hashing
	ErrInvalidPasswordHashFormat = errors.New("invalid password hash format")
	ErrInvalidArgonVersion       = errors.New("invalid argon2 version")
	ErrInvalidArgonParams        = errors.New("invalid argon2 parameters")
	ErrInvalidHashSalt           = errors.New("invalid password hash salt")
	ErrInvalidPasswordHash       = errors.New("invalid password hash")
)
