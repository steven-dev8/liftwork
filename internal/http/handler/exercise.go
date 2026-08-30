package handler

import (
	"encoding/json"
	"liftwork/internal/http/middleware"
	"liftwork/internal/service"
	"net/http"
	"time"
)

type ExerciseHandler struct {
	service *service.ExerciseService
}

func NewExerciseHandler(service *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{service: service}
}

type createExerciseRequest struct {
	Name        string `json:"name"`
	MuscleGroup string `json:"muscle_group"`
	Notes       string `json:"notes"`
}

type createExerciseResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	MuscleGroup string    `json:"muscle_group"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h *ExerciseHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var request createExerciseRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	exercise, err := h.service.Create(r.Context(), service.CreateExerciseInput{
		UserID:      userID,
		Name:        request.Name,
		MuscleGroup: request.MuscleGroup,
		Notes:       request.Notes,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, createExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	})
}
