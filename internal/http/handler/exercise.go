package handler

import (
	"encoding/json"
	"liftwork/internal/http/middleware"
	"liftwork/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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

type updateExerciseRequest struct {
	Name        *string `json:"name"`
	MuscleGroup *string `json:"muscle_group"`
	Notes       *string `json:"notes"`
}

type ExerciseResponse struct {
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

	writeJSON(w, http.StatusCreated, ExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	})
}

func (h *ExerciseHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	exercises, err := h.service.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response := make([]ExerciseResponse, len(exercises))

	for i, exercise := range exercises {
		response[i] = ExerciseResponse{
			ID:          exercise.ID,
			Name:        exercise.Name,
			MuscleGroup: exercise.MuscleGroup,
			Notes:       exercise.Notes,
			CreatedAt:   exercise.CreatedAt,
			UpdatedAt:   exercise.UpdatedAt,
		}
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *ExerciseHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	exerciseID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || exerciseID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var request updateExerciseRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	exercise, err := h.service.Update(r.Context(), service.UpdateExerciseInput{
		ID:          exerciseID,
		UserID:      userID,
		Name:        request.Name,
		MuscleGroup: request.MuscleGroup,
		Notes:       request.Notes,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		MuscleGroup: exercise.MuscleGroup,
		Notes:       exercise.Notes,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	})
}

func (h *ExerciseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	exerciseID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || exerciseID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	if err := h.service.Delete(r.Context(), service.DeleteExerciseInput{
		ID:     exerciseID,
		UserID: userID,
	}); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
