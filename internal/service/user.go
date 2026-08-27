package service

import (
	"context"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"liftwork/internal/security"
	"time"
)

type UserService struct {
	user            repository.UserRepository
	session         repository.SessionRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewUserService(
	user repository.UserRepository,
	session repository.SessionRepository,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *UserService {
	return &UserService{
		user:            user,
		session:         session,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type CreateUserInput struct {
	Username string
	Email    string
	Password string
}

type CreateUserOutput struct {
	ID              int64
	Username        string
	CreatedAt       time.Time
	AccessToken     string
	DurationAcess   time.Duration
	RefreshToken    string
	DurationRefresh time.Time
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

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return CreateUserOutput{}, err
	}

	refreshTokenHash := security.HashRefreshToken(refreshToken)
	session, err := s.session.Create(ctx, repository.CreateSessionParams{
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(s.refreshTokenTTL),
	})

	if err != nil {
		return CreateUserOutput{}, err
	}

	AccessToken, err := security.CreateToken(user.ID, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{
		ID:              user.ID,
		Username:        user.Username,
		CreatedAt:       user.CreatedAt,
		AccessToken:     AccessToken,
		DurationAcess:   s.accessTokenTTL,
		RefreshToken:    refreshToken,
		DurationRefresh: session.ExpiresAt,
	}, nil
}
