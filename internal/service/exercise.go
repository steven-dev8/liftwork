package service

import (
	"context"
	"errors"
	"fmt"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"time"
)

type ExerciseService struct {
	exerciseRepo repository.ExerciseRepository
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{exerciseRepo: exerciseRepo}
}

type CreateExerciseInput struct {
	UserID      int64
	Name        string
	MuscleGroup string
	Notes       string
}

type CreateExerciseOutput struct {
	ID          int64
	Name        string
	MuscleGroup string
	Notes       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *ExerciseService) Create(
	ctx context.Context,
	input CreateExerciseInput,
) (CreateExerciseOutput, error) {
	exercise := domain.Exercise{
		Name:        input.Name,
		MuscleGroup: input.MuscleGroup,
		Notes:       input.Notes,
	}

	if err := exercise.VerifyFields(); err != nil {
		return CreateExerciseOutput{}, err
	}

	exercise, err := s.exerciseRepo.Create(ctx, input.UserID, exercise)
	if err != nil {
		if errors.Is(err, repository.ErrExerciseAlreadyExists) {
			return CreateExerciseOutput{}, ErrExerciseAlreadyExists
		}

		return CreateExerciseOutput{}, fmt.Errorf("create exercise: %w", err)
	}

	return CreateExerciseOutput{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdateAt,
	}, nil
}
