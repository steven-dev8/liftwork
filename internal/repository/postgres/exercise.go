package postgres

import (
	"context"
	"errors"
	"fmt"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
	"liftwork/internal/repository"

	"github.com/jackc/pgx/v5"
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
	exercise.UpdatedAt = row.UpdatedAt.Time

	return exercise, nil
}

func (e *ExerciseRepository) List(
	ctx context.Context,
	userID int64,
) ([]domain.Exercise, error) {
	rows, err := e.querier.GetExercises(ctx, &userID)

	if err != nil {
		return nil, err
	}

	exercises := make([]domain.Exercise, len(rows))

	for index, exercise := range rows {
		exercises[index] = domain.Exercise{
			ID:          exercise.ID,
			Name:        exercise.Name,
			MuscleGroup: exercise.MuscleGroup,
			Notes:       exercise.Notes,
			CreatedAt:   exercise.CreatedAt.Time,
			UpdatedAt:   exercise.UpdatedAt.Time,
		}
	}

	return exercises, nil
}

func (e *ExerciseRepository) Update(
	ctx context.Context,
	exerciseInfo repository.ExerciseUpdateParams,
) (domain.Exercise, error) {
	exercise, err := e.querier.UpdateExerciseById(ctx, db.UpdateExerciseByIdParams{
		ID:          exerciseInfo.ID,
		UserID:      &exerciseInfo.UserID,
		Name:        exerciseInfo.Name,
		MuscleGroup: exerciseInfo.MuscleGroup,
		Notes:       exerciseInfo.Notes,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Exercise{}, repository.ErrExerciseNotFound
		}

		return domain.Exercise{}, fmt.Errorf("update exercise: %w", err)
	}

	return domain.Exercise{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt.Time,
		UpdatedAt:   exercise.UpdatedAt.Time,
	}, nil
}

func (e *ExerciseRepository) Delete(
	ctx context.Context,
	id int64,
	userID int64,
) error {
	rows, err := e.querier.DeleteExerciseByID(ctx, db.DeleteExerciseByIDParams{
		ID:     id,
		UserID: &userID,
	})

	if err != nil {
		return fmt.Errorf("delete exercise: %w", err)
	}

	if rows == 0 {
		return repository.ErrExerciseNotFound
	}

	return nil
}
