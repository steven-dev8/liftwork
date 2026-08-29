package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"liftwork/internal/domain"
	"liftwork/internal/repository"
	"liftwork/internal/security"
)

type AuthService struct {
	user            repository.UserRepository
	session         repository.SessionRepository
	jwtSecret       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthService(
	user repository.UserRepository,
	session repository.SessionRepository,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		user:            user,
		session:         session,
		jwtSecret:       jwtSecret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

type RegisterUserInput struct {
	Username string
	Email    string
	Password string
}

type RegisterUserOutput struct {
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

type RefreshAccessTokenOutput struct {
	AccessToken      string
	AccessTokenTTL   time.Duration
	RefreshToken     string
	RefreshExpiresAt time.Time
}

func (s *AuthService) Register(ctx context.Context, input RegisterUserInput) (RegisterUserOutput, error) {
	if err := domain.ValidatePassword(input.Password); err != nil {
		return RegisterUserOutput{}, err
	}

	user, err := domain.NewUser(
		input.Username,
		input.Email,
	)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	passwordHash, err := security.HashPassword(input.Password)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	user.PasswordHash = passwordHash

	user, err = s.user.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUsernameAlreadyExists):
			return RegisterUserOutput{}, ErrUsernameAlreadyExists
		case errors.Is(err, repository.ErrEmailAlreadyExists):
			return RegisterUserOutput{}, ErrEmailAlreadyExists
		default:
			return RegisterUserOutput{}, fmt.Errorf("create user: %w", err)
		}
	}

	refreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return RegisterUserOutput{}, err
	}

	refreshTokenHash := security.HashRefreshToken(refreshToken)
	session, err := s.session.Create(ctx, repository.CreateSessionParams{
		UserID:           user.ID,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(s.refreshTokenTTL),
	})

	if err != nil {
		return RegisterUserOutput{}, err
	}

	accessToken, err := security.CreateToken(user.ID, s.jwtSecret, s.accessTokenTTL)
	if err != nil {
		return RegisterUserOutput{}, err
	}

	return RegisterUserOutput{
		ID:               user.ID,
		Username:         user.Username,
		CreatedAt:        user.CreatedAt,
		AccessToken:      accessToken,
		AccessTokenTTL:   s.accessTokenTTL,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginUserInput) (LoginUserOutput, error) {
	username := strings.ToLower(strings.TrimSpace(input.Username))
	user, err := s.user.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return LoginUserOutput{}, ErrInvalidCredentials
		}

		return LoginUserOutput{}, fmt.Errorf("find user by username: %w", err)
	}

	passwordEqual, err := security.ComparePassword(input.Password, user.PasswordHash)
	if err != nil {
		return LoginUserOutput{}, err
	}

	if !passwordEqual {
		return LoginUserOutput{}, ErrInvalidCredentials
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

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	refreshTokenHash := security.HashRefreshToken(refreshToken)

	_, err := s.session.RevokeSession(ctx, refreshTokenHash)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAccessToken(
	ctx context.Context,
	refreshToken string,
) (RefreshAccessTokenOutput, error) {
	refreshTokenHash := security.HashRefreshToken(refreshToken)

	session, err := s.session.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, repository.ErrSessionNotFound) {
			return RefreshAccessTokenOutput{}, ErrInvalidRefreshToken
		}

		return RefreshAccessTokenOutput{}, fmt.Errorf("find refresh session: %w", err)
	}

	if !session.IsActive(time.Now().UTC()) {
		return RefreshAccessTokenOutput{}, ErrInvalidRefreshToken
	}

	accessToken, err := security.CreateToken(
		session.UserID,
		s.jwtSecret,
		s.accessTokenTTL,
	)
	if err != nil {
		return RefreshAccessTokenOutput{}, err
	}

	newRefreshToken, err := security.GenerateRefreshToken()
	if err != nil {
		return RefreshAccessTokenOutput{}, err
	}

	newRefreshTokenHash := security.HashRefreshToken(newRefreshToken)

	rotated, err := s.session.RotateRefreshToken(ctx, repository.RotateRefreshTokenParams{
		SessionID:           session.ID,
		OldRefreshTokenHash: refreshTokenHash,
		NewRefreshTokenHash: newRefreshTokenHash,
	})
	if err != nil {
		return RefreshAccessTokenOutput{}, err
	}

	if !rotated {
		return RefreshAccessTokenOutput{}, ErrInvalidRefreshToken
	}

	return RefreshAccessTokenOutput{
		AccessToken:      accessToken,
		AccessTokenTTL:   s.accessTokenTTL,
		RefreshToken:     newRefreshToken,
		RefreshExpiresAt: session.ExpiresAt,
	}, nil
}
