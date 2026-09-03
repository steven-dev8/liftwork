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

type RoutineService struct {
	repository repository.RoutineRepository
}

func NewRoutineService(routineRepo repository.RoutineRepository) *RoutineService {
	return &RoutineService{repository: routineRepo}
}

type CreateRoutineInput struct {
	UserID      int64
	Name        string
	Code        string
	Description string
}

type UpdateRoutineInput struct {
	ID          int64
	UserID      int64
	Name        *string
	Code        *string
	Description *string
}

type UpdateExerciseRoutineInput struct {
	UserID        int64
	RoutineID     int64
	ExerciseID    int64
	TargetSets    *int32
	TargetRepsMin *int32
	TargetRepsMax *int32
}

type RoutineOutput struct {
	ID          int64
	Name        string
	Code        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ExerciseRoutine struct {
	ID            int64
	Name          string
	Position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}

type AddExerciseRoutineInput struct {
	UserID        int64
	RoutineID     int64
	ExerciseID    int64
	Position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}

type ListRoutineOutput struct {
	Routine   RoutineOutput
	Exercises []ExerciseRoutine
}

type DeleteExerciseRoutineInput struct {
	UserID     int64
	ExerciseID int64
	RoutineID  int64
}

func (r *RoutineService) Create(
	ctx context.Context,
	input CreateRoutineInput,
) (RoutineOutput, error) {
	routine, err := domain.NewRoutine(
		input.UserID,
		domain.RoutineCode(input.Code),
		input.Name,
		input.Description,
	)

	if err != nil {
		return RoutineOutput{}, err
	}

	routineInfo, err := r.repository.Create(ctx, routine)
	if err != nil {
		return RoutineOutput{}, fmt.Errorf("create routine: %w", err)
	}

	return routineToOutput(routineInfo), nil
}

func (r *RoutineService) List(
	ctx context.Context,
	userID int64,
) ([]ListRoutineOutput, error) {
	routines, err := r.repository.List(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list routines: %w", err)
	}

	output := make([]ListRoutineOutput, len(routines))

	for i, routineInfo := range routines {
		exercises := make([]ExerciseRoutine, len(routineInfo.Exercises))

		for j, exercise := range routineInfo.Exercises {
			exercises[j] = ExerciseRoutine{
				ID:            exercise.ID,
				Name:          exercise.Name,
				Position:      exercise.Position,
				TargetSets:    exercise.TargetSets,
				TargetRepsMin: exercise.TargetRepsMin,
				TargetRepsMax: exercise.TargetRepsMax,
			}
		}

		output[i] = ListRoutineOutput{
			Routine:   routineToOutput(routineInfo.Routine),
			Exercises: exercises,
		}
	}

	return output, nil
}

func (r *RoutineService) Update(
	ctx context.Context,
	input UpdateRoutineInput,
) (RoutineOutput, error) {
	if input.Code == nil &&
		input.Name == nil &&
		input.Description == nil {
		return RoutineOutput{}, ErrEmptyRoutineUpdate
	}

	if input.Code != nil {
		code := strings.TrimSpace(*input.Code)

		if code == "" {
			return RoutineOutput{}, domain.ErrInvalidRoutineCode
		}

		routineCode := domain.RoutineCode(code)
		if !routineCode.IsValid() {
			return RoutineOutput{}, domain.ErrInvalidRoutineCode
		}

		input.Code = &code
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)

		if name == "" {
			return RoutineOutput{}, domain.ErrRoutineNameRequired
		}

		input.Name = &name
	}

	if input.Description != nil {
		description := strings.TrimSpace(*input.Description)
		input.Description = &description
	}

	routine, err := r.repository.Update(ctx, repository.UpdateRoutineParams{
		ID:          input.ID,
		UserID:      input.UserID,
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
	})
	if err != nil {
		if errors.Is(err, repository.ErrRoutineNotFound) {
			return RoutineOutput{}, ErrRoutineNotFound
		}

		return RoutineOutput{}, fmt.Errorf("update routine: %w", err)
	}

	return routineToOutput(routine), nil
}

func (r *RoutineService) Delete(
	ctx context.Context,
	userID int64,
	routineID int64,
) error {
	if err := r.repository.Delete(ctx, userID, routineID); err != nil {
		if errors.Is(err, repository.ErrRoutineNotFound) {
			return ErrRoutineNotFound
		}

		return fmt.Errorf("delete routine: %w", err)
	}

	return nil
}

func (r *RoutineService) AddExerciseRoutine(
	ctx context.Context,
	input AddExerciseRoutineInput,
) error {
	routineExercise, err := domain.NewRoutineExercise(
		input.RoutineID,
		input.ExerciseID,
		input.Position,
		input.TargetSets,
		input.TargetRepsMin,
		input.TargetRepsMax,
	)
	if err != nil {
		return err
	}

	if err := r.repository.AddExerciseRoutine(
		ctx,
		input.UserID,
		routineExercise,
	); err != nil {
		if errors.Is(err, repository.ErrRoutineOrExerciseNotFound) {
			return ErrRoutineOrExerciseNotFound
		}

		return fmt.Errorf("add exercise to routine: %w", err)
	}

	return nil
}

func (r *RoutineService) UpdateExerciseRoutine(
	ctx context.Context,
	input UpdateExerciseRoutineInput,
) (ExerciseRoutine, error) {
	if input.TargetSets == nil &&
		input.TargetRepsMin == nil &&
		input.TargetRepsMax == nil {
		return ExerciseRoutine{}, ErrEmptyRoutineExerciseUpdate
	}

	if input.TargetSets != nil && *input.TargetSets <= 0 {
		return ExerciseRoutine{}, domain.ErrInvalidTargetSets
	}

	if input.TargetRepsMin != nil && *input.TargetRepsMin <= 0 {
		return ExerciseRoutine{}, domain.ErrInvalidTargetRepsMin
	}

	if input.TargetRepsMax != nil && *input.TargetRepsMax <= 0 {
		return ExerciseRoutine{}, domain.ErrInvalidTargetRepsRange
	}

	if input.TargetRepsMin != nil &&
		input.TargetRepsMax != nil &&
		*input.TargetRepsMax < *input.TargetRepsMin {
		return ExerciseRoutine{}, domain.ErrInvalidTargetRepsRange
	}

	exerciseRoutine, err := r.repository.UpdateExerciseRoutine(
		ctx,
		repository.UpdateExerciseRoutineParams{
			UserID:        input.UserID,
			RoutineID:     input.RoutineID,
			ExerciseID:    input.ExerciseID,
			TargetSets:    input.TargetSets,
			TargetRepsMin: input.TargetRepsMin,
			TargetRepsMax: input.TargetRepsMax,
		},
	)
	if err != nil {
		if errors.Is(err, repository.ErrRoutineExerciseNotFound) {
			return ExerciseRoutine{}, ErrRoutineExerciseNotFound
		}

		return ExerciseRoutine{}, fmt.Errorf(
			"update routine exercise: %w",
			err,
		)
	}

	return ExerciseRoutine{
		Position:      exerciseRoutine.Position,
		TargetSets:    exerciseRoutine.TargetSets,
		TargetRepsMin: exerciseRoutine.TargetRepsMin,
		TargetRepsMax: exerciseRoutine.TargetRepsMax,
	}, nil
}

func (r *RoutineService) DeleteExerciseRoutine(
	ctx context.Context,
	input DeleteExerciseRoutineInput,
) ([]ListRoutineOutput, error) {
	if err := r.repository.DeleteExerciseRoutine(
		ctx,
		input.UserID,
		input.RoutineID,
		input.ExerciseID,
	); err != nil {
		if errors.Is(err, repository.ErrRoutineExerciseNotFound) {
			return nil, ErrRoutineExerciseNotFound
		}

		return nil, fmt.Errorf("delete exercise from routine: %w", err)
	}

	return r.List(ctx, input.UserID)
}

func routineToOutput(routine domain.Routine) RoutineOutput {
	return RoutineOutput{
		ID:          routine.ID,
		Name:        routine.Name,
		Code:        string(routine.Code),
		Description: routine.Description,
		CreatedAt:   routine.CreatedAt,
		UpdatedAt:   routine.UpdatedAt,
	}
}
