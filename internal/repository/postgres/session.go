package postgres

import (
	"context"
	"fmt"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type SessionRepository struct {
	querier *db.Queries
}

func NewSessionRepository(dbtx db.DBTX) *SessionRepository {
	return &SessionRepository{querier: db.New(dbtx)}
}

func (s *SessionRepository) Create(ctx context.Context, params repository.CreateSessionParams) (repository.CreateSessionOutput, error) {
	row, err := s.querier.CreateSession(ctx, db.CreateSessionParams{
		UserID:           params.UserID,
		RefreshTokenHash: params.RefreshTokenHash,
		ExpiresAt: pgtype.Timestamptz{
			Time:  params.ExpiresAt,
			Valid: true,
		},
	})

	if err != nil {
		return repository.CreateSessionOutput{}, fmt.Errorf("create session: %w", err)
	}

	return repository.CreateSessionOutput{
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt.Time,
	}, nil
}

func (s *SessionRepository) RevokeSession(ctx context.Context, refreshTokenHash string) (bool, error) {
	row, err := s.querier.RevokeSession(ctx, refreshTokenHash)
	if err != nil {
		return false, fmt.Errorf("revoke session: %w", err)
	}

	return row > 0, nil
}
