package api

import (
	"liftwork/internal/config"
	db "liftwork/internal/database/sqlc"
	httpapi "liftwork/internal/http"
	"liftwork/internal/http/handler"
	"liftwork/internal/repository/postgres"
	"liftwork/internal/service"
)

type App struct {
	Handlers httpapi.Handlers
}

func New(dbtx db.DBTX, cfg *config.Config) *App {
	userRepo := postgres.NewUserRepository(dbtx)
	userService := service.NewUserService(userRepo, cfg.JWTSecretKey, cfg.JWTTTL)
	userHandler := handler.NewUser(userService)

	return &App{
		Handlers: httpapi.Handlers{
			User: userHandler,
		},
	}
}
