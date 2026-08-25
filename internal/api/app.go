package api

import (
	db "liftwork/internal/database/sqlc"
	httpapi "liftwork/internal/http"
	"liftwork/internal/http/handler"
	"liftwork/internal/repository/postgres"
	"liftwork/internal/service"
	"log/slog"
)

type App struct {
	Handlers httpapi.Handlers
}

func New(dbtx db.DBTX, logger *slog.Logger) *App {
	userRepo := postgres.NewUserRepository(dbtx)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUser(userService)

	return &App{
		Handlers: httpapi.Handlers{
			User: userHandler,
		},
	}
}
