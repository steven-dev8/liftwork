package domain

import "errors"

var (
	ErrInvalidEmail          = errors.New("invalid email")
	ErrInvalidUsernameLength = errors.New("username must have between 3 and 15 characters")
	ErrInvalidUsernameChars  = errors.New("username can only contain lowercase letters, numbers and underscores")
	ErrNumericUsername       = errors.New("username cannot contain only numbers")
	ErrPasswordTooShort      = errors.New("password must have at least 10 characters")
	ErrPasswordContainsSpace = errors.New("password must not contain spaces")
	ErrPasswordRequiresDigit = errors.New("password must contain at least one number")
	ErrPasswordRequiresUpper = errors.New("password must contain at least one uppercase letter")

	ErrExerciseNameRequired       = errors.New("exercise name cannot be empty.")
	ErrExerciseMucleGroupRequired = errors.New("muscle group name cannot be empty.")
)
