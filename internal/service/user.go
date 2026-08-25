package service

import (
	"context"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"liftwork/internal/security"
	"time"
)

type UserService struct {
	user repository.UserRepository
}

func NewUserService(user repository.UserRepository) *UserService {
	return &UserService{user: user}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUserOutput struct {
	Username  string
	CreatedAt time.Time
}

func (s *UserService) Create(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	if err := domain.ValidatePassword(input.Password); err != nil {
		return CreateUserOutput{}, err
	}

	user, err := domain.NewUser(
		input.Username,
		input.Email,
	)
	if err != nil {
		return CreateUserOutput{}, err
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return CreateUserOutput{}, err
	}

	user.PasswordHash = passwordHash

	user, err = s.user.Create(ctx, user)
	if err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}, nil
}
