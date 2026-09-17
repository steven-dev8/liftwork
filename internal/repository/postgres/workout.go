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
	"github.com/jackc/pgx/v5/pgconn"
)

type WorkoutRepository struct {
	querier *db.Queries
	dbtx    database.Transactor
}

func NewWorkoutRepository(dbtx database.Transactor) *WorkoutRepository {
	return &WorkoutRepository{
		querier: db.New(dbtx),
		dbtx:    dbtx,
	}
}

func (w *WorkoutRepository) CreateWithoutRoutine(
	ctx context.Context,
	userID int64,
	workout domain.WorkoutSession,
) (domain.WorkoutSession, error) {
	workoutS, err := w.querier.CreateWorkoutSessionWithoutRoutine(
		ctx,
		db.CreateWorkoutSessionWithoutRoutineParams{
			UserID: userID,
			Notes:  workout.Notes,
		},
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "one_open_workout_per_user" {
			return domain.WorkoutSession{},
				repository.ErrWorkoutAlreadyOpen
		}

		return domain.WorkoutSession{}, fmt.Errorf(
			"create workout session without routine: %w",
			err,
		)
	}

	return workoutSessionFromRow(workoutS), nil
}

func (w *WorkoutRepository) CreateFromRoutine(
	ctx context.Context,
	userID int64,
	workout domain.WorkoutSession,
) (repository.WorkoutWithExercises, error) {
	if workout.RoutineID == nil {
		return repository.WorkoutWithExercises{},
			repository.ErrRoutineNotFound
	}

	tx, err := w.dbtx.Begin(ctx)
	if err != nil {
		return repository.WorkoutWithExercises{},
			fmt.Errorf("begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	qtx := w.querier.WithTx(tx)

	workoutS, err := qtx.CreateWorkoutSessionWithRoutine(
		ctx,
		db.CreateWorkoutSessionWithRoutineParams{
			UserID:    userID,
			RoutineID: *workout.RoutineID,
			Notes:     workout.Notes,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.WorkoutWithExercises{},
				repository.ErrRoutineNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "one_open_workout_per_user" {
			return repository.WorkoutWithExercises{},
				repository.ErrWorkoutAlreadyOpen
		}

		return repository.WorkoutWithExercises{},
			fmt.Errorf("create workout session from routine: %w", err)
	}

	workoutExercises, err := qtx.CreateWorkoutExercisesFromRoutine(
		ctx,
		db.CreateWorkoutExercisesFromRoutineParams{
			WorkoutSessionID: workoutS.ID,
			RoutineID:        *workout.RoutineID,
		},
	)
	if err != nil {
		return repository.WorkoutWithExercises{},
			fmt.Errorf("create workout exercises from routine: %w", err)
	}

	exercises := make(
		[]repository.WorkoutExerciseInfo,
		len(workoutExercises),
	)

	for i, exercise := range workoutExercises {
		exercises[i] = repository.WorkoutExerciseInfo{
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

	if err := tx.Commit(ctx); err != nil {
		return repository.WorkoutWithExercises{},
			fmt.Errorf("commit transaction: %w", err)
	}

	return repository.WorkoutWithExercises{
		Workout:   workoutSessionFromRow(workoutS),
		Exercises: exercises,
	}, nil
}

func workoutSessionFromRow(row db.WorkoutSession) domain.WorkoutSession {
	workout := domain.WorkoutSession{
		ID:        row.ID,
		RoutineID: row.RoutineID,
		Notes:     row.Notes,
		CreatedAt: row.CreatedAt.Time,
	}

	if row.StartedAt.Valid {
		startedAt := row.StartedAt.Time
		workout.StartedAt = &startedAt
	}

	if row.FinishedAt.Valid {
		finishedAt := row.FinishedAt.Time
		workout.FinishedAt = &finishedAt
	}

	return workout
}
