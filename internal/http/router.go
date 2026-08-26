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
	User *handler.User
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
	router.Post("/users/create_user", handlers.User.Create)

	return router
}
