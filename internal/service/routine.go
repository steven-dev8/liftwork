package service

import (
	"context"
	"errors"
	"fmt"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
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
