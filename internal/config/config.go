package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHTTPAddr        = ":8080"
	defaultReadHeader      = 5 * time.Second
	defaultConnectTimeout  = 5 * time.Second
	defaultShutdown        = 10 * time.Second
	defaultJWTTTL          = 15 * time.Minute
	defaultRefreshTokenTTL = 30 * 24 * time.Hour
)

// Config contains all runtime settings needed by the API
type Config struct {
	HTTPAddr              string
	HTTPReadHeaderTimeout time.Duration
	JWTSecretKey          string
	JWTTTL                time.Duration
	RefreshTokenTTL       time.Duration
	DatabaseURL           string
	CORSAllowedOrigins    []string
	DBMaxConns            int32
	DBMinConns            int32
	DBConnectTimeout      time.Duration
	ShutdownTimeout       time.Duration
}

// Load all configuration from environment variables (.env)
func Load() (Config, error) {
	maxConns, err := int32FromEnv("DB_MAX_CONNS", 10)
	if err != nil {
		return Config{}, err
	}

	minConns, err := int32FromEnv("DB_MIN_CONNS", 1)
	if err != nil {
		return Config{}, err
	}

	connectTimeout, err := durationFromEnv("DB_CONNECT_TIMEOUT", defaultConnectTimeout)
	if err != nil {
		return Config{}, err
	}

	shutdownTimeout, err := durationFromEnv("SHUTDOWN_TIMEOUT", defaultShutdown)
	if err != nil {
		return Config{}, err
	}

	readHeaderTimeout, err := durationFromEnv("HTTP_READ_HEADER_TIMEOUT", defaultReadHeader)
	if err != nil {
		return Config{}, err
	}

	jwtTTL, err := durationFromEnv("JWT_ACCESS_TOKEN_TTL", defaultJWTTTL)
	if err != nil {
		return Config{}, err
	}

	refreshTokenTTL, err := durationFromEnv("REFRESH_TOKEN_TTL", defaultRefreshTokenTTL)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		HTTPAddr:              valueOrDefault("HTTP_ADDR", defaultHTTPAddr),
		HTTPReadHeaderTimeout: readHeaderTimeout,
		JWTSecretKey:          strings.TrimSpace(os.Getenv("JWT_SECRET")),
		JWTTTL:                jwtTTL,
		RefreshTokenTTL:       refreshTokenTTL,
		DatabaseURL:           strings.TrimSpace(os.Getenv("DATABASE_URL")),
		CORSAllowedOrigins:    commaSeparatedEnv("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"}),
		DBMaxConns:            maxConns,
		DBMinConns:            minConns,
		DBConnectTimeout:      connectTimeout,
		ShutdownTimeout:       shutdownTimeout,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}
	if cfg.JWTSecretKey == "" {
		return Config{}, errors.New("JWT_SECRET is required")
	}
	if cfg.JWTTTL <= 0 {
		return Config{}, errors.New("JWT_ACCESS_TOKEN_TTL must be greater than zero")
	}
	if cfg.RefreshTokenTTL <= 0 {
		return Config{}, errors.New("REFRESH_TOKEN_TTL must be greater than zero")
	}
	if cfg.DBMaxConns < 1 {
		return Config{}, errors.New("DB_MAX_CONNS must be at least 1")
	}
	if cfg.DBMinConns < 0 || cfg.DBMinConns > cfg.DBMaxConns {
		return Config{}, errors.New("DB_MIN_CONNS must be between 0 and DB_MAX_CONNS")
	}

	return cfg, nil
}

// Assistant functions //

func valueOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func commaSeparatedEnv(key string, fallback []string) []string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parts := strings.Split(value, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

func int32FromEnv(key string, fallback int32) (int32, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return int32(parsed), nil
}

func durationFromEnv(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", key, err)
	}
	return parsed, nil
}
