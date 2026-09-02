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

type addExerciseRoutineRequest struct {
	ExerciseID    int64 `json:"exercise_id"`
	Position      int32 `json:"position"`
	TargetSets    int32 `json:"target_sets"`
	TargetRepsMin int32 `json:"target_reps_min"`
	TargetRepsMax int32 `json:"target_reps_max"`
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
	ExerciseID    int64  `json:"exercise_id"`
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
				ExerciseID:    exercise.ID,
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

func (h *RoutineHandler) AddExerciseRoutine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	idParam := chi.URLParam(r, "id")

	routineID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || routineID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid routine id")
		return
	}

	var request addExerciseRoutineRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := h.service.AddExerciseRoutine(r.Context(), service.AddExerciseRoutineInput{
		UserID:        userID,
		ExerciseID:    request.ExerciseID,
		RoutineID:     routineID,
		Position:      request.Position,
		TargetSets:    request.TargetSets,
		TargetRepsMin: request.TargetRepsMin,
		TargetRepsMax: request.TargetRepsMax,
	}); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *RoutineHandler) DeleteExerciseRoutine(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
		return
	}

	routineIDParam := chi.URLParam(r, "routineID")
	exerciseIDParam := chi.URLParam(r, "exerciseID")

	routineID, err := strconv.ParseInt(routineIDParam, 10, 64)
	if err != nil || routineID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid routine id")
		return
	}

	exerciseID, err := strconv.ParseInt(exerciseIDParam, 10, 64)
	if err != nil || exerciseID <= 0 {
		writeError(w, http.StatusBadRequest, "invalid exercise id")
		return
	}

	routines, err := h.service.DeleteExerciseRoutine(
		r.Context(),
		service.DeleteExerciseRoutineInput{
			UserID:     userID,
			RoutineID:  routineID,
			ExerciseID: exerciseID,
		},
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	output := make([]listRoutineResponse, len(routines))

	for i, routine := range routines {
		exercises := make([]routineExerciseResponse, len(routine.Exercises))

		for j, exercise := range routine.Exercises {
			exercises[j] = routineExerciseResponse{
				ExerciseID:    exercise.ID,
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
