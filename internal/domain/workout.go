package domain

import "time"

// WorkoutSession is one execution of a training routine
type WorkoutSession struct {
	ID         int64
	RoutineID  int64
	StartedAt  time.Time
	FinishedAt time.Time
	Notes      string
	CreatedAt  time.Time
}

// PerformedSet records load, repetition and repetition in reserve (RIR)
type PerformedSet struct {
	ID               int64
	WorkoutSessionID int64
	ExerciseID       int64
	SetNumber        int32
	WeightKG         float64
	Repetition       int32
	RIR              int32
	PerformedAt      time.Time
}
