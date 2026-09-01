package service

import (
	"context"
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
