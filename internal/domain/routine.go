package domain

import (
	"strings"
	"time"
)

type RoutineCode string

const (
	RoutineA RoutineCode = "A"
	RoutineB RoutineCode = "B"
	RoutineC RoutineCode = "C"
	RoutineD RoutineCode = "D"
	RoutineE RoutineCode = "E"
)

type Routine struct {
	ID          int64
	UserID      int64
	Code        RoutineCode
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoutineExercise struct {
	RoutineID     int64
	ExerciseID    int64
	Position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}

func NewRoutine(
	userID int64,
	code RoutineCode,
	name string,
	description string,
) (Routine, error) {
	code = RoutineCode(strings.TrimSpace(string(code)))
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if !code.IsValid() {
		return Routine{}, ErrInvalidRoutineCode
	}

	if name == "" {
		return Routine{}, ErrRoutineNameRequired
	}

	return Routine{
		UserID:      userID,
		Code:        code,
		Name:        name,
		Description: description,
	}, nil
}

func NewRoutineExercise(
	routineID int64,
	exerciseID int64,
	position int32,
	targetSets int32,
	targetRepsMin int32,
	targetRepsMax int32,
) (RoutineExercise, error) {
	if position <= 0 {
		return RoutineExercise{}, ErrInvalidExercisePosition
	}

	if targetSets <= 0 {
		return RoutineExercise{}, ErrInvalidTargetSets
	}

	if targetRepsMin <= 0 {
		return RoutineExercise{}, ErrInvalidTargetRepsMin
	}

	if targetRepsMax < targetRepsMin {
		return RoutineExercise{}, ErrInvalidTargetRepsRange
	}

	return RoutineExercise{
		RoutineID:     routineID,
		ExerciseID:    exerciseID,
		Position:      position,
		TargetSets:    targetSets,
		TargetRepsMin: targetRepsMin,
		TargetRepsMax: targetRepsMax,
	}, nil
}

func (c RoutineCode) IsValid() bool {
	switch c {
	case RoutineA, RoutineB, RoutineC, RoutineD, RoutineE:
		return true
	default:
		return false
	}
}
