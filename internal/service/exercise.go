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
	repository repository.ExerciseRepository
}

func NewExerciseService(exerciseRepo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repository: exerciseRepo}
}

type CreateExerciseInput struct {
	UserID      int64
	Name        string
	MuscleGroup string
	Notes       string
}

type ExerciseOutput struct {
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
) (ExerciseOutput, error) {
	exercise := domain.Exercise{
		Name:        input.Name,
		MuscleGroup: input.MuscleGroup,
		Notes:       input.Notes,
	}

	if err := exercise.VerifyFields(); err != nil {
		return ExerciseOutput{}, err
	}

	exercise, err := s.repository.Create(ctx, input.UserID, exercise)
	if err != nil {
		if errors.Is(err, repository.ErrExerciseAlreadyExists) {
			return ExerciseOutput{}, ErrExerciseAlreadyExists
		}

		return ExerciseOutput{}, fmt.Errorf("create exercise: %w", err)
	}

	return ExerciseOutput{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdateAt,
	}, nil
}

func (s *ExerciseService) List(ctx context.Context, userID int64) ([]ExerciseOutput, error) {
	exercises, err := s.repository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list exercises: %w", err)
	}

	listExercises := make([]ExerciseOutput, len(exercises))

	for index, exercise := range exercises {
		listExercises[index] = ExerciseOutput{
			ID:          exercise.ID,
			Name:        exercise.Name,
			MuscleGroup: exercise.MuscleGroup,
			Notes:       exercise.Notes,
			CreatedAt:   exercise.CreatedAt,
			UpdatedAt:   exercise.UpdateAt,
		}
	}

	return listExercises, nil
}
