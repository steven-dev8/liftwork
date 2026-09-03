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
	RoutineID int64  `json:"routine_id"`
}

type createWorkoutResponse struct {
	ID        int64     `json:"id"`
	RoutineID int64     `json:"routine_id"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
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

	workoutS, err := h.service.Create(r.Context(), service.CreateWorkoutSessionInput{
		UserID:    userID,
		RoutineID: request.RoutineID,
		Notes:     request.Notes,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createWorkoutResponse{
		ID:        workoutS.ID,
		RoutineID: workoutS.RoutineID,
		Notes:     workoutS.Notes,
		CreatedAt: workoutS.CreatedAt,
	})
}
