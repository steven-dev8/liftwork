package postgres

import (
	"context"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
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
