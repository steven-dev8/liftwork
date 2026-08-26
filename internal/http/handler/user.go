package handler

import (
	"encoding/json"
	"liftwork/internal/service"
	"net/http"
	"time"
)

type User struct {
	service *service.UserService
}

func NewUser(userService *service.UserService) *User {
	return &User{service: userService}
}

type createUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreatedUserResponse struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	AccessToken string    `json:"access_token"`
	TokenType  string    `json:"token_type"`
	ExpiresIn  int64     `json:"expires_in"`
}

func (h *User) Create(w http.ResponseWriter, r *http.Request) {
	var request createUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	user, err := h.service.Create(r.Context(), service.CreateUserInput{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, CreatedUserResponse{
		ID:         user.ID,
		Username:   user.Username,
		CreatedAt:  user.CreatedAt,
		AccessToken: user.TokenJWT,
		TokenType:  "Bearer",
		ExpiresIn:  int64(user.Duration.Seconds()),
	})
}
