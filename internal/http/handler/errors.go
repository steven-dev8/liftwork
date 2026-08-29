package handler

import (
	"errors"
	"net/http"

	"liftwork/internal/domain"
	"liftwork/internal/service"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, service.ErrInvalidCredentials.Error())
	case errors.Is(err, service.ErrInvalidRefreshToken):
		writeError(w, http.StatusUnauthorized, service.ErrInvalidRefreshToken.Error())
	case errors.Is(err, service.ErrUsernameAlreadyExists):
		writeError(w, http.StatusConflict, service.ErrUsernameAlreadyExists.Error())
	case errors.Is(err, service.ErrEmailAlreadyExists):
		writeError(w, http.StatusConflict, service.ErrEmailAlreadyExists.Error())
	default:
		message, ok := domainValidationMessage(err)
		if ok {
			writeError(w, http.StatusUnprocessableEntity, message)
			return
		}

		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func domainValidationMessage(err error) (string, bool) {
	validationErrors := []error{
		domain.ErrInvalidEmail,
		domain.ErrInvalidUsernameLength,
		domain.ErrInvalidUsernameChars,
		domain.ErrNumericUsername,
		domain.ErrPasswordTooShort,
		domain.ErrPasswordContainsSpace,
		domain.ErrPasswordRequiresDigit,
		domain.ErrPasswordRequiresUpper,
	}

	for _, validationErr := range validationErrors {
		if errors.Is(err, validationErr) {
			return validationErr.Error(), true
		}
	}

	return "", false
}
