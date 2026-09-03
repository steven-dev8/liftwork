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
}

func NewWorkoutRepository(dbtx database.Transactor) *WorkoutRepository {
	return &WorkoutRepository{
		querier: db.New(dbtx),
	}
}

func (w *WorkoutRepository) Create(
	ctx context.Context,
	userID int64,
	workout domain.WorkoutSession,
) (domain.WorkoutSession, error) {
	workoutS, err := w.querier.CreateWorkoutSession(
		ctx,
		db.CreateWorkoutSessionParams{
			UserID:    userID,
			RoutineID: workout.RoutineID,
			Notes:     workout.Notes,
		},
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.WorkoutSession{}, repository.ErrRoutineNotFound
		}

		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "one_open_workout_per_user" {
			return domain.WorkoutSession{},
				repository.ErrWorkoutAlreadyOpen
		}

		return domain.WorkoutSession{}, fmt.Errorf(
			"create workout session: %w",
			err,
		)
	}

	workout.ID = workoutS.ID
	workout.CreatedAt = workoutS.CreatedAt.Time

	return workout, nil
}
