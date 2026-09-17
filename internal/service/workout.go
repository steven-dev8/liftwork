package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"liftwork/internal/domain"
	"liftwork/internal/repository"
)

type WorkoutService struct {
	repository repository.WorkoutRepository
}

func NewWorkoutService(
	workoutRepo repository.WorkoutRepository,
) *WorkoutService {
	return &WorkoutService{
		repository: workoutRepo,
	}
}

type CreateWorkoutSessionInput struct {
	UserID    int64
	RoutineID *int64
	Notes     string
}

type WorkoutExerciseRoutine struct {
	ID            int64
	ExerciseID    int64
	Name          string
	MuscleGroup   string
	Notes         string
	Position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}

type WorkoutSessionOutput struct {
	ID         int64
	RoutineID  *int64
	StartedAt  *time.Time
	FinishedAt *time.Time
	Notes      string
	CreatedAt  time.Time
	Exercises  []WorkoutExerciseRoutine
}

func (w *WorkoutService) Create(
	ctx context.Context,
	input CreateWorkoutSessionInput,
) (WorkoutSessionOutput, error) {
	if input.RoutineID == nil {
		return w.createFreeWorkout(ctx, input)
	}
	if *input.RoutineID <= 0 {
		return WorkoutSessionOutput{}, ErrInvalidRoutineID
	}

	return w.createWorkoutFromRoutine(ctx, input)
}

func (w *WorkoutService) createFreeWorkout(
	ctx context.Context,
	input CreateWorkoutSessionInput,
) (WorkoutSessionOutput, error) {
	notes := strings.TrimSpace(input.Notes)

	workout := domain.WorkoutSession{
		RoutineID: nil,
		Notes:     notes,
	}

	createdWorkout, err := w.repository.CreateWithoutRoutine(
		ctx,
		input.UserID,
		workout,
	)
	if err != nil {
		if errors.Is(err, repository.ErrWorkoutAlreadyOpen) {
			return WorkoutSessionOutput{}, ErrWorkoutAlreadyOpen
		}

		return WorkoutSessionOutput{}, fmt.Errorf(
			"create free workout session: %w",
			err,
		)
	}

	return WorkoutSessionOutput{
		ID:         createdWorkout.ID,
		RoutineID:  createdWorkout.RoutineID,
		StartedAt:  createdWorkout.StartedAt,
		FinishedAt: createdWorkout.FinishedAt,
		Notes:      createdWorkout.Notes,
		CreatedAt:  createdWorkout.CreatedAt,
		Exercises:  []WorkoutExerciseRoutine{},
	}, nil
}

func (w *WorkoutService) createWorkoutFromRoutine(
	ctx context.Context,
	input CreateWorkoutSessionInput,
) (WorkoutSessionOutput, error) {
	notes := strings.TrimSpace(input.Notes)

	workout := domain.WorkoutSession{
		RoutineID: input.RoutineID,
		Notes:     notes,
	}

	createdWorkout, err := w.repository.CreateFromRoutine(
		ctx,
		input.UserID,
		workout,
	)
	if err != nil {
		if errors.Is(err, repository.ErrRoutineNotFound) {
			return WorkoutSessionOutput{}, ErrRoutineNotFound
		}

		if errors.Is(err, repository.ErrWorkoutAlreadyOpen) {
			return WorkoutSessionOutput{}, ErrWorkoutAlreadyOpen
		}

		return WorkoutSessionOutput{}, fmt.Errorf(
			"create workout session from routine: %w",
			err,
		)
	}

	exercises := make(
		[]WorkoutExerciseRoutine,
		len(createdWorkout.Exercises),
	)

	for i, exercise := range createdWorkout.Exercises {
		exercises[i] = WorkoutExerciseRoutine{
			ID:            exercise.ID,
			ExerciseID:    exercise.ExerciseID,
			Name:          exercise.Name,
			MuscleGroup:   exercise.MuscleGroup,
			Notes:         exercise.Notes,
			Position:      exercise.Position,
			TargetSets:    exercise.TargetSets,
			TargetRepsMin: exercise.TargetRepsMin,
			TargetRepsMax: exercise.TargetRepsMax,
		}
	}

	return WorkoutSessionOutput{
		ID:         createdWorkout.Workout.ID,
		RoutineID:  createdWorkout.Workout.RoutineID,
		StartedAt:  createdWorkout.Workout.StartedAt,
		FinishedAt: createdWorkout.Workout.FinishedAt,
		Notes:      createdWorkout.Workout.Notes,
		CreatedAt:  createdWorkout.Workout.CreatedAt,
		Exercises:  exercises,
	}, nil
}
