package postgres

import (
	"context"
	"errors"
	"fmt"

	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
	"liftwork/internal/repository"

	"github.com/jackc/pgx/v5"
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

func (s *SessionRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (domain.Session, error) {
	row, err := s.querier.GetSessionByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Session{}, repository.ErrSessionNotFound
		}

		return domain.Session{}, fmt.Errorf("get session: %w", err)
	}

	session := domain.Session{
		ID:        row.ID,
		UserID:    row.UserID,
		ExpiresAt: row.ExpiresAt.Time,
	}

	if row.RevokedAt.Valid {
		session.RevokedAt = &row.RevokedAt.Time
	}

	return session, nil
}

func (s *SessionRepository) RotateRefreshToken(
	ctx context.Context,
	params repository.RotateRefreshTokenParams,
) (bool, error) {
	rowsAffected, err := s.querier.RotateRefreshToken(
		ctx,
		db.RotateRefreshTokenParams{
			SessionID:           params.SessionID,
			NewRefreshTokenHash: params.NewRefreshTokenHash,
			OldRefreshTokenHash: params.OldRefreshTokenHash,
		},
	)
	if err != nil {
		return false, fmt.Errorf("rotate refresh token: %w", err)
	}

	return rowsAffected > 0, nil
}
