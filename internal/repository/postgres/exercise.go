package postgres

import (
	"context"
	"errors"
	"fmt"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
	"liftwork/internal/repository"

	"github.com/jackc/pgx/v5/pgconn"
)

type ExerciseRepository struct {
	querier *db.Queries
}

func NewExerciseRepository(dbtx db.DBTX) *ExerciseRepository {
	return &ExerciseRepository{querier: db.New(dbtx)}
}

func (e *ExerciseRepository) Create(
	ctx context.Context,
	userID int64,
	exercise domain.Exercise,
) (domain.Exercise, error) {
	row, err := e.querier.CreateExercise(ctx, db.CreateExerciseParams{
		UserID:      &userID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
	})

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "exercises_user_id_name_key" {

			return domain.Exercise{}, repository.ErrExerciseAlreadyExists
		}

		return domain.Exercise{}, fmt.Errorf("create exercise: %w", err)
	}

	exercise.ID = row.ID
	exercise.CreatedAt = row.CreatedAt.Time
	exercise.UpdateAt = row.UpdatedAt.Time

	return exercise, nil
}
