package repository

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrSessionNotFound       = errors.New("session not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrExerciseAlreadyExists = errors.New("exercise already exists")
	ErrExerciseNotFound      = errors.New("exercise not found")

	ErrRoutineOrExerciseNotFound = errors.New("routine or exercise not found")
)
