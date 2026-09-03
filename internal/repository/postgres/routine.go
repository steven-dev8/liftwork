package postgres

import (
	"context"
	"errors"
	"fmt"
	"liftwork/internal/database"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
	"liftwork/internal/repository"

	"github.com/jackc/pgx/v5"
)

type RoutineRepository struct {
	querier *db.Queries
	dbtx    database.Transactor
}

func NewRoutineRepository(dbtx database.Transactor) *RoutineRepository {
	return &RoutineRepository{
		querier: db.New(dbtx),
		dbtx:    dbtx,
	}
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
				ID:            exercise.ID,
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

func (r *RoutineRepository) Update(
	ctx context.Context,
	routineInfo repository.UpdateRoutineParams,
) (domain.Routine, error) {
	routine, err := r.querier.UpdateRoutine(ctx, db.UpdateRoutineParams{
		Code:        routineInfo.Code,
		Name:        routineInfo.Name,
		Description: routineInfo.Description,
		ID:          routineInfo.ID,
		UserID:      routineInfo.UserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Routine{}, repository.ErrRoutineNotFound
		}

		return domain.Routine{}, fmt.Errorf("update routine: %w", err)
	}

	return domain.Routine{
		ID:          routine.ID,
		UserID:      routine.UserID,
		Code:        domain.RoutineCode(routine.Code),
		Name:        routine.Name,
		Description: routine.Description,
		CreatedAt:   routine.CreatedAt.Time,
		UpdatedAt:   routine.UpdatedAt.Time,
	}, nil
}

func (r *RoutineRepository) Delete(
	ctx context.Context,
	userID int64,
	routineID int64,
) error {
	rows, err := r.querier.DeleteRoutine(ctx, db.DeleteRoutineParams{
		ID:     routineID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("delete routine: %w", err)
	}

	if rows == 0 {
		return repository.ErrRoutineNotFound
	}

	return nil
}

func (r *RoutineRepository) AddExerciseRoutine(
	ctx context.Context,
	userID int64,
	routineExercise domain.RoutineExercise,
) error {
	rows, err := r.querier.AddExerciseRoutine(ctx, db.AddExerciseRoutineParams{
		UserID:        userID,
		RoutineID:     routineExercise.RoutineID,
		ExerciseID:    routineExercise.ExerciseID,
		Position:      routineExercise.Position,
		TargetSets:    routineExercise.TargetSets,
		TargetRepsMin: routineExercise.TargetRepsMin,
		TargetRepsMax: routineExercise.TargetRepsMax,
	})
	if err != nil {
		return fmt.Errorf("add exercise to routine: %w", err)
	}

	if rows == 0 {
		return repository.ErrRoutineOrExerciseNotFound
	}

	return nil
}

func (r *RoutineRepository) DeleteExerciseRoutine(
	ctx context.Context,
	userID int64,
	routineID int64,
	exerciseID int64,
) error {
	tx, err := r.dbtx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	qtx := r.querier.WithTx(tx)

	position, err := qtx.DeleteExerciseRoutine(
		ctx,
		db.DeleteExerciseRoutineParams{
			UserID:     userID,
			RoutineID:  routineID,
			ExerciseID: exerciseID,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrRoutineExerciseNotFound
		}

		return fmt.Errorf("delete exercise from routine: %w", err)
	}

	err = qtx.ReorderRoutineExercises(
		ctx,
		db.ReorderRoutineExercisesParams{
			RoutineID:       routineID,
			DeletedPosition: position,
		},
	)
	if err != nil {
		return fmt.Errorf("reorder routine exercises: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
