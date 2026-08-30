package service

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrExerciseAlreadyExists = errors.New("exercise already exists")
)
