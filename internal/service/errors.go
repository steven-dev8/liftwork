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

	ErrRoutineNotFound            = errors.New("routine not found")
	ErrInvalidRoutineID           = errors.New("routine ID invalid")
	ErrEmptyRoutineUpdate         = errors.New("no routine fields to update")
	ErrRoutineOrExerciseNotFound  = errors.New("routine or exercise not found")
	ErrRoutineExerciseNotFound    = errors.New("routine exercise not found")
	ErrEmptyRoutineExerciseUpdate = errors.New("no routine exercise fields to update")

	ErrWorkoutAlreadyOpen = errors.New("workout already open")
)
