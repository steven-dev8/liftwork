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

type UpdateExerciseInput struct {
	ID          int64
	UserID      int64
	Name        *string
	MuscleGroup *string
	Notes       *string
}

type DeleteExerciseInput struct {
	ID     int64
	UserID int64
}

func (s *ExerciseService) Create(
	ctx context.Context,
	input CreateExerciseInput,
) (ExerciseOutput, error) {
	exercise := domain.Exercise{
		Name:        strings.TrimSpace(input.Name),
		MuscleGroup: strings.TrimSpace(input.MuscleGroup),
		Notes:       strings.TrimSpace(input.Notes),
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
		UpdatedAt:   exercise.UpdatedAt,
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
			UpdatedAt:   exercise.UpdatedAt,
		}
	}

	return listExercises, nil
}

func (s *ExerciseService) Update(
	ctx context.Context,
	input UpdateExerciseInput,
) (ExerciseOutput, error) {
	if input.Name == nil &&
		input.MuscleGroup == nil &&
		input.Notes == nil {
		return ExerciseOutput{}, ErrEmptyExerciseUpdate
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)

		if name == "" {
			return ExerciseOutput{}, ErrInvalidExerciseName
		}

		input.Name = &name
	}

	if input.MuscleGroup != nil {
		muscleGroup := strings.TrimSpace(*input.MuscleGroup)

		if muscleGroup == "" {
			return ExerciseOutput{}, ErrInvalidExerciseMuscleGroup
		}

		input.MuscleGroup = &muscleGroup
	}

	if input.Notes != nil {
		notes := strings.TrimSpace(*input.Notes)
		input.Notes = &notes
	}

	exercise, err := s.repository.Update(ctx, repository.ExerciseUpdateParams{
		ID:          input.ID,
		UserID:      input.UserID,
		Name:        input.Name,
		MuscleGroup: input.MuscleGroup,
		Notes:       input.Notes,
	})
	if err != nil {
		if errors.Is(err, repository.ErrExerciseNotFound) {
			return ExerciseOutput{}, ErrExerciseNotFound
		}

		return ExerciseOutput{}, fmt.Errorf("update exercise: %w", err)
	}

	return ExerciseOutput{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	}, nil
}

func (s *ExerciseService) Delete(ctx context.Context, input DeleteExerciseInput) error {
	if err := s.repository.Delete(ctx, input.ID, input.UserID); err != nil {
		if errors.Is(err, repository.ErrExerciseNotFound) {
			return ErrExerciseNotFound
		}
		return fmt.Errorf("delete exercise: %w", err)
	}

	return nil
}
