package postgres

import (
	"context"
	"fmt"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
)

type RoutineRepository struct {
	querier *db.Queries
}

func NewRoutineRepository(dbtx db.DBTX) *RoutineRepository {
	return &RoutineRepository{querier: db.New(dbtx)}
}

func (r *RoutineRepository) Create(
	ctx context.Context,
	routine domain.Routine,
) (domain.Routine, error) {
	row, err := r.querier.CreateRoutine(ctx, db.CreateRoutineParams{
		UserID:      routine.UserID,
		Code:        string(routine.Code),
		Name:        routine.Name,
		Description: routine.Description,
	})
	if err != nil {
		return domain.Routine{}, err
	}

	routine.ID = row.ID
	routine.CreatedAt = row.CreatedAt.Time
	routine.UpdatedAt = row.UpdatedAt.Time

	return routine, nil
}

func (r *RoutineRepository) List(
	ctx context.Context,
	userID int64,
) ([]repository.RoutineWithExercises, error) {
	routines, err := r.querier.ListRoutine(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list routines: %w", err)
	}

	listRoutines := make([]repository.RoutineWithExercises, len(routines))

	for i, routine := range routines {
		exercises, err := r.querier.GetExerciseRoutine(ctx, routine.ID)
		if err != nil {
			return nil, fmt.Errorf(
				"list exercises for routine %d: %w",
				routine.ID,
				err,
			)
		}

		routineExercises := make(
			[]repository.RoutineExerciseInfo,
			len(exercises),
		)

		for j, exercise := range exercises {
			routineExercises[j] = repository.RoutineExerciseInfo{
				Name:          exercise.Name,
				Position:      exercise.Position,
				TargetSets:    exercise.TargetSets,
				TargetRepsMin: exercise.TargetRepsMin,
				TargetRepsMax: exercise.TargetRepsMax,
			}
		}

		listRoutines[i] = repository.RoutineWithExercises{
			Routine: domain.Routine{
				ID:          routine.ID,
				UserID:      userID,
				Code:        domain.RoutineCode(routine.Code),
				Name:        routine.Name,
				Description: routine.Description,
				CreatedAt:   routine.CreatedAt.Time,
				UpdatedAt:   routine.UpdatedAt.Time,
			},
			Exercises: routineExercises,
		}
	}

	return listRoutines, nil
}
