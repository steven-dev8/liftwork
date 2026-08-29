package repository

import (
	"context"
	"liftwork/internal/domain"
	"time"
)

type CreateSessionParams struct {
	UserID           int64
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type RotateRefreshTokenParams struct {
	SessionID           int64
	OldRefreshTokenHash string
	NewRefreshTokenHash string
}

type CreateSessionOutput = CreateSessionParams

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
	RevokeSession(ctx context.Context, refreshTokenHash string) (bool, error)
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (domain.Session, error)
	RotateRefreshToken(ctx context.Context, params RotateRefreshTokenParams) (bool, error)
}
