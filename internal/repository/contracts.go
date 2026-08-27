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

type CreateSessionOutput = CreateSessionParams

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
}
