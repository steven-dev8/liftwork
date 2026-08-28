package handler

import (
	"encoding/json"
	"liftwork/internal/service"
	"net/http"
	"strings"
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

type createdUserResponse struct {
	ID               int64     `json:"id"`
	Username         string    `json:"username"`
	CreatedAt        time.Time `json:"created_at"`
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

type loginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginUserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type logoutUserRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type loginUserResponse struct {
	UserInfo         loginUserInfo `json:"user"`
	AccessToken      string        `json:"access_token"`
	TokenType        string        `json:"token_type"`
	ExpiresIn        int64         `json:"expires_in"`
	RefreshToken     string        `json:"refresh_token"`
	RefreshExpiresAt time.Time     `json:"refresh_expires_at"`
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

	writeJSON(w, http.StatusCreated, createdUserResponse{
		ID:               user.ID,
		Username:         user.Username,
		CreatedAt:        user.CreatedAt,
		AccessToken:      user.AccessToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(user.AccessTokenTTL.Seconds()),
		RefreshToken:     user.RefreshToken,
		RefreshExpiresAt: user.RefreshExpiresAt,
	})
}

func (h *User) Login(w http.ResponseWriter, r *http.Request) {
	var request loginUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	user, err := h.service.Login(r.Context(), service.LoginUserInput{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, loginUserResponse{
		UserInfo: loginUserInfo{
			ID:       user.UserInfo.ID,
			Username: user.UserInfo.Username,
		},
		AccessToken:      user.AccessToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(user.AccessTokenTTL.Seconds()),
		RefreshToken:     user.RefreshToken,
		RefreshExpiresAt: user.RefreshExpiresAt,
	})
}

func (h *User) Logout(w http.ResponseWriter, r *http.Request) {
	var request logoutUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if strings.TrimSpace(request.RefreshToken) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "refresh_token is required",
		})
		return
	}

	if err := h.service.Logout(r.Context(), request.RefreshToken); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
