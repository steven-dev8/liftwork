package repository

import (
	"context"
	"liftwork/internal/domain"
	"time"
)

type CreateSessionParams struct {
	UserID           int64
	RefreshTokenHash string
	ExpiresAt        time.Time
}

type RotateRefreshTokenParams struct {
	SessionID           int64
	OldRefreshTokenHash string
	NewRefreshTokenHash string
}

type ExerciseUpdateParams struct {
	ID          int64
	UserID      int64
	Name        *string
	MuscleGroup *string
	Notes       *string
}

type CreateSessionOutput = CreateSessionParams

type RoutineExerciseInfo struct {
	ID            int64
	Name          string
	Position      int32
	TargetSets    int32
	TargetRepsMin int32
	TargetRepsMax int32
}

type RoutineWithExercises struct {
	Routine   domain.Routine
	Exercises []RoutineExerciseInfo
}

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, params CreateSessionParams) (CreateSessionOutput, error)
	RevokeSession(ctx context.Context, refreshTokenHash string) (bool, error)
	FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (domain.Session, error)
	RotateRefreshToken(ctx context.Context, params RotateRefreshTokenParams) (bool, error)
}

type ExerciseRepository interface {
	Create(ctx context.Context, userID int64, exercise domain.Exercise) (domain.Exercise, error)
	List(ctx context.Context, userID int64) ([]domain.Exercise, error)
	Update(ctx context.Context, exerciseInfo ExerciseUpdateParams) (domain.Exercise, error)
	Delete(ctx context.Context, id int64, userID int64) error
}

type RoutineRepository interface {
	Create(ctx context.Context, routine domain.Routine) (domain.Routine, error)
	List(ctx context.Context, userID int64) ([]RoutineWithExercises, error)
	AddExerciseRoutine(ctx context.Context, userID int64, routineExercise domain.RoutineExercise) error
	DeleteExerciseRoutine(ctx context.Context, userID int64, routineID int64, exerciseID int64) error
}
