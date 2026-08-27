package postgres

import (
	"context"
	"fmt"
	db "liftwork/internal/database/sqlc"
	"liftwork/internal/domain"
)

type UserRepository struct {
	querier *db.Queries
}

func NewUserRepository(dbtx db.DBTX) *UserRepository {
	return &UserRepository{querier: db.New(dbtx)}
}

func (u *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	var email *string

	if user.Email != "" {
		email = &user.Email
	}

	row, err := u.querier.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	})

	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	user.ID = row.ID
	user.Username = row.Username
	user.CreatedAt = row.CreatedAt.Time

	return user, nil
}


func (u *UserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := u.querier.GetUser(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID: row.ID,
		Username: row.Username,
		PasswordHash: row.PasswordHash,
	}, nil
}