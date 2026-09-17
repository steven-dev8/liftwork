package handler

import (
	"encoding/json"
	"liftwork/internal/http/middleware"
	"liftwork/internal/service"
	"net/http"
	"time"
)

type WorkoutHandler struct {
	service *service.WorkoutService
}

func NewWorkoutHandler(service *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{service: service}
}

type createWorkoutRequest struct {
	Notes     string `json:"notes"`
	RoutineID *int64 `json:"routine_id"`
}

type workoutExerciseResponse struct {
	ID            int64  `json:"id"`
	ExerciseID    int64  `json:"exercise_id"`
	Name          string `json:"name"`
	MuscleGroup   string `json:"muscle_group"`
	Notes         string `json:"notes"`
	Position      int32  `json:"position"`
	TargetSets    int32  `json:"target_sets"`
	TargetRepsMin int32  `json:"target_reps_min"`
	TargetRepsMax int32  `json:"target_reps_max"`
}

type createWorkoutResponse struct {
	ID         int64                     `json:"id"`
	RoutineID  *int64                    `json:"routine_id"`
	StartedAt  *time.Time                `json:"started_at"`
	FinishedAt *time.Time                `json:"finished_at"`
	Notes      string                    `json:"notes"`
	CreatedAt  time.Time                 `json:"created_at"`
	Exercises  []workoutExerciseResponse `json:"exercises"`
}

func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var request createWorkoutRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	workoutS, err := h.service.Create(
		r.Context(),
		service.CreateWorkoutSessionInput{
			UserID:    userID,
			RoutineID: request.RoutineID,
			Notes:     request.Notes,
		},
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	exercises := make(
		[]workoutExerciseResponse,
		len(workoutS.Exercises),
	)

	for i, exercise := range workoutS.Exercises {
		exercises[i] = workoutExerciseResponse{
			ID:            exercise.ID,
			ExerciseID:    exercise.ExerciseID,
			Name:          exercise.Name,
			MuscleGroup:   exercise.MuscleGroup,
			Notes:         exercise.Notes,
			Position:      exercise.Position,
			TargetSets:    exercise.TargetSets,
			TargetRepsMin: exercise.TargetRepsMin,
			TargetRepsMax: exercise.TargetRepsMax,
		}
	}

	writeJSON(w, http.StatusCreated, createWorkoutResponse{
		ID:         workoutS.ID,
		RoutineID:  workoutS.RoutineID,
		StartedAt:  workoutS.StartedAt,
		FinishedAt: workoutS.FinishedAt,
		Notes:      workoutS.Notes,
		CreatedAt:  workoutS.CreatedAt,
		Exercises:  exercises,
	})
}
