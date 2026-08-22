package domain

import "time"

type RoutineCode string

const (
	RotutineA RoutineCode = "A"
	RotutineB RoutineCode = "B"
	RotutineC RoutineCode = "C"
	RotutineD RoutineCode = "D"
	RotutineE RoutineCode = "E"
)

type Routine struct {
	ID          int64
	Code        RoutineCode
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoutineExercise struct {
	RoutineID     int64
	ExerciseID    int64
	position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}
