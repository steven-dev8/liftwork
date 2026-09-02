package service

import "errors"

var (
	ErrInvalidCredentials    = errors.New("invalid username or password")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")

	ErrExerciseAlreadyExists      = errors.New("exercise already exists")
	ErrInvalidExerciseName        = errors.New("exercise name cannot be empty")
	ErrInvalidExerciseMuscleGroup = errors.New("muscle group name cannot be empty")
	ErrEmptyExerciseUpdate        = errors.New("no exercise fields to update")
	ErrExerciseNotFound           = errors.New("exercise not found")

	ErrRoutineOrExerciseNotFound = errors.New("routine or exercise not found")
	ErrRoutineExerciseNotFound   = errors.New("routine exercise not found")
)
