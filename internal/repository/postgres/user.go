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

type UserRepository struct {
	querier *db.Queries
}

func NewUserRepository(dbtx db.DBTX) *UserRepository {
	return &UserRepository{querier: db.New(dbtx)}
}

func (u *UserRepository) Create(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
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
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_username_key":
				return domain.User{},
					repository.ErrUsernameAlreadyExists

			case "users_email_key":
				return domain.User{},
					repository.ErrEmailAlreadyExists
			}
		}

		return domain.User{}, fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	user.ID = row.ID
	user.Username = row.Username
	user.CreatedAt = row.CreatedAt.Time

	return user, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := u.querier.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, repository.ErrUserNotFound
		}

		return domain.User{}, fmt.Errorf(
			"find user by username: %w",
			err,
		)
	}

	return domain.User{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}, nil
}
