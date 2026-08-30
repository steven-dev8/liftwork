package httpapi

import (
	"liftwork/internal/config"
	"liftwork/internal/http/handler"
	"liftwork/internal/http/middleware"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handlers struct {
	Auth     *handler.AuthHandler
	Exercise *handler.ExerciseHandler
}

func NewRouter(
	logger *slog.Logger,
	cfg *config.Config,
	handlers Handlers,
) http.Handler {

	router := chi.NewRouter()
	router.Use(middleware.Logging(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.CORSAllowedOrigins,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	router.Get("/ping", handler.Ping)

	// Auth
	router.Post("/auth/register", handlers.Auth.Register)
	router.Post("/auth/login", handlers.Auth.Login)
	router.Post("/auth/logout", handlers.Auth.Logout)
	router.Post("/auth/refresh", handlers.Auth.Refresh)

	router.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecretKey))

		r.Post("/exercise", handlers.Exercise.Create)
	})

	return router
}
