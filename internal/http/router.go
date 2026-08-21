package httpapi

import (
	"liftwork/internal/http/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewRouter(allowedOrigins []string) http.Handler {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		MaxAge:         300,
	}))

	router.Get("/ping", handler.Ping)

	return router
}
