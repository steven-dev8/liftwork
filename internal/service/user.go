package service

import (
	"context"
	"errors"
	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"liftwork/internal/security"
	"strings"
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
	ID               int64
	Username         string
	CreatedAt        time.Time
	AccessToken      string
	AccessTokenTTL   time.Duration
	RefreshToken     string
	RefreshExpiresAt time.Time
}

type LoginUserInput struct {
	Username string
	Password string
}

type UserInfoOutput struct {
	ID       int64
	Username string
}

type LoginUserOutput struct {
	UserInfo         UserInfoOutput
	AccessToken      string
	AccessTokenTTL   time.Duration
	RefreshToken     string
	RefreshExpiresAt time.Time
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

	accessToken, err := security.CreateToken(user.ID, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{
		ID:               user.ID,
		Username:         user.Username,
		CreatedAt:        user.CreatedAt,
		AccessToken:      accessToken,
		AccessTokenTTL:   s.accessTokenTTL,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *UserService) Login(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	username := strings.ToLower(strings.TrimSpace(input.Username))
	user, err := s.user.FindByUsername(ctx, username)
	if err != nil {
		return LoginUserOutput{}, err
	}

	passwordEqual, err := security.ComparePassword(input.Password, user.PasswordHash)
	if err != nil {
		return LoginUserOutput{}, err
	}

	if !passwordEqual {
		return LoginUserOutput{}, errors.New("invalid username or password")
	}

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return LoginUserOutput{}, err
	}

	refreshTokenHash := security.HashRefreshToken(refreshToken)
	session, err := s.session.Create(ctx, repository.CreateSessionParams{
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(s.refreshTokenTTL),
	})

	if err != nil {
		return LoginUserOutput{}, err
	}

	accessToken, err := security.CreateToken(user.ID, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return LoginUserOutput{}, err
	}

	return LoginUserOutput{
		UserInfo: UserInfoOutput{
			ID:       user.ID,
			Username: user.Username,
		},
		AccessToken:      accessToken,
		AccessTokenTTL:   s.accessTokenTTL,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *UserService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := security.HashRefreshToken(refreshToken)

	_, err := s.session.RevokeSession(ctx, refreshTokenHash)
	if err != nil {
		return err
	}

	return nil
}
