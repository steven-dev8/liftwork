package domain

import "time"

// Exercise describes a movement independently from any routine
type Exercise struct {
	ID          int64
	Name        string
	MuscleGroup string
	Notes       string
	CreatedAt   time.Time
	UpdateAt    time.Time
}

// LastExercisePerformance summarizes the most recent set performed for an exercise.
type LastExercisePerformance struct {
	ID               int64
	WorkoutSessionID int64
	ExerciseID       int64
	SetNumber        int32
	WeightKG         float64
	Repetitions      int32
	RIR              int32
	PerformedAt      time.Time
	Found            bool
}

