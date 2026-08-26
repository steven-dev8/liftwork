package service

import (
	"context"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"liftwork/internal/security"
	"time"
)

type UserService struct {
	user          repository.UserRepository
	jwtSecret     string
	acessTokenTTL time.Duration
}

func NewUserService(
	user repository.UserRepository,
	jwtSecret string,
	acessTokenTTL time.Duration,
) *UserService {
	return &UserService{
		user:          user,
		jwtSecret:     jwtSecret,
		acessTokenTTL: acessTokenTTL,
	}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUserOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	TokenJWT  string
	Duration  time.Duration
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

	tokenJWT, err := security.CreateToken(user.ID, s.jwtSecret, s.acessTokenTTL)
	if err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		TokenJWT:  tokenJWT,
		Duration:  s.acessTokenTTL,
	}, nil
}
