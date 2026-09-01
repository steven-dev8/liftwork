package handler

import (
	"encoding/json"
	"liftwork/internal/http/middleware"
	"liftwork/internal/service"
	"net/http"
	"time"
)

type RoutineHandler struct {
	service *service.RoutineService
}

func NewRoutineHandler(service *service.RoutineService) *RoutineHandler {
	return &RoutineHandler{service: service}
}

type createRoutineRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type routineResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type routineExerciseResponse struct {
	Name          string `json:"name"`
	Position      int32  `json:"position"`
	TargetSets    int32  `json:"target_sets"`
	TargetRepsMin int32  `json:"target_reps_min"`
	TargetRepsMax int32  `json:"target_reps_max"`
}

type listRoutineResponse struct {
	ID          int64                     `json:"id"`
	Name        string                    `json:"name"`
	Code        string                    `json:"code"`
	Description string                    `json:"description"`
	Exercises   []routineExerciseResponse `json:"exercises"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

func (h *RoutineHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	var request createRoutineRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	routine, err := h.service.Create(r.Context(), service.CreateRoutineInput{
		UserID:      userID,
		Name:        request.Name,
		Code:        request.Code,
		Description: request.Description,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, routineResponse{
		ID:          routine.ID,
		Name:        routine.Name,
		Code:        routine.Code,
		Description: routine.Description,
		CreatedAt:   routine.CreatedAt,
		UpdatedAt:   routine.UpdatedAt,
	})
}

func (h *RoutineHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	routines, err := h.service.List(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	output := make([]listRoutineResponse, len(routines))

	for i, routine := range routines {
		exercises := make([]routineExerciseResponse, len(routine.Exercises))

		for j, exercise := range routine.Exercises {
			exercises[j] = routineExerciseResponse{
				Name:          exercise.Name,
				Position:      exercise.Position,
				TargetSets:    exercise.TargetSets,
				TargetRepsMin: exercise.TargetRepsMin,
				TargetRepsMax: exercise.TargetRepsMax,
			}
		}

		output[i] = listRoutineResponse{
			ID:          routine.Routine.ID,
			Name:        routine.Routine.Name,
			Code:        routine.Routine.Code,
			Description: routine.Routine.Description,
			Exercises:   exercises,
			CreatedAt:   routine.Routine.CreatedAt,
			UpdatedAt:   routine.Routine.UpdatedAt,
		}
	}

	writeJSON(w, http.StatusOK, output)
}
