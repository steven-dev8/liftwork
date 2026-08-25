package httpapi

import (
	"liftwork/internal/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Handlers struct {
	User *handler.User
}

func NewRouter(allowedOrigins []string, handlers Handlers) http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	router.Get("/ping", handler.Ping)
	router.Post("/users/create_user", handlers.User.Create)

	return router
}
