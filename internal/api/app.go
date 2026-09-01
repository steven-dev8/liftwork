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
	sessionRepo := postgres.NewSessionRepository(dbtx)
	userRepo := postgres.NewUserRepository(dbtx)

	exerciseRepo := postgres.NewExerciseRepository(dbtx)
	exerciseService := service.NewExerciseService(exerciseRepo)
	exerciseHandler := handler.NewExerciseHandler(exerciseService)

	authService := service.NewAuthService(userRepo, sessionRepo, cfg.JWTSecretKey, cfg.JWTTTL, cfg.RefreshTokenTTL)
	authHandler := handler.NewAuthHandler(authService)

	routineRepo := postgres.NewRoutineRepository(dbtx)
	routineService := service.NewRoutineService(routineRepo)
	routineHandler := handler.NewRoutineHandler(routineService)

	return &App{
		Handlers: httpapi.Handlers{
			Auth:     authHandler,
			Exercise: exerciseHandler,
			Routine:  routineHandler,
		},
	}
}
