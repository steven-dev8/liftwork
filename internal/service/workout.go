package service

import (
	"context"
	"errors"
	"fmt"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"strings"
	"time"
)

type WorkoutService struct {
	repository repository.WorkoutRepository
}

func NewWorkoutService(WorkoutRepo repository.WorkoutRepository) *WorkoutService {
	return &WorkoutService{repository: WorkoutRepo}
}

type CreateWorkoutSessionInput struct {
	UserID    int64
	RoutineID int64
	Notes     string
}

type WorkoutSessionOutput struct {
	ID         int64
	RoutineID  int64
	StartedAt  *time.Time
	FinishedAt *time.Time
	Notes      string
	CreatedAt  time.Time
}

func (w *WorkoutService) Create(
	ctx context.Context,
	input CreateWorkoutSessionInput,
) (WorkoutSessionOutput, error) {
	if input.RoutineID <= 0 {
		return WorkoutSessionOutput{}, ErrInvalidRoutineID
	}

	notes := strings.TrimSpace(input.Notes)

	workout := domain.WorkoutSession{
		RoutineID: input.RoutineID,
		Notes:     notes,
	}

	createdWorkout, err := w.repository.Create(
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
			"create workout session: %w",
			err,
		)
	}

	return workoutSessionToOutput(createdWorkout), nil
}

func workoutSessionToOutput(
	workout domain.WorkoutSession,
) WorkoutSessionOutput {
	return WorkoutSessionOutput{
		ID:         workout.ID,
		RoutineID:  workout.RoutineID,
		StartedAt:  workout.StartedAt,
		FinishedAt: workout.FinishedAt,
		Notes:      workout.Notes,
		CreatedAt:  workout.CreatedAt,
	}
}
