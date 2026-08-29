package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"liftwork/internal/service"
)

type Auth struct {
	service *service.AuthService
}

func NewAuth(authService *service.AuthService) *Auth {
	return &Auth{service: authService}
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registeredUserResponse struct {
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

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	RefreshToken     string    `json:"refresh_token"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
}

func (h *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var request registerUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(request.Username) == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	if strings.TrimSpace(request.Password) == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	user, err := h.service.Register(r.Context(), service.RegisterUserInput{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	})

	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, registeredUserResponse{
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

func (h *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var request loginUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(request.Username) == "" {
		writeError(w, http.StatusBadRequest, "username is required")
		return
	}

	if strings.TrimSpace(request.Password) == "" {
		writeError(w, http.StatusBadRequest, "password is required")
		return
	}

	user, err := h.service.Login(r.Context(), service.LoginUserInput{
		Username: request.Username,
		Password: request.Password,
	})

	if err != nil {
		writeServiceError(w, err)
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

func (h *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	var request logoutUserRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(request.RefreshToken) == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	if err := h.service.Logout(r.Context(), request.RefreshToken); err != nil {
		writeServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Auth) Refresh(w http.ResponseWriter, r *http.Request) {
	var request refreshTokenRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if strings.TrimSpace(request.RefreshToken) == "" {
		writeError(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	output, err := h.service.RefreshAccessToken(
		r.Context(),
		request.RefreshToken,
	)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, refreshTokenResponse{
		AccessToken:      output.AccessToken,
		TokenType:        "Bearer",
		ExpiresIn:        int64(output.AccessTokenTTL.Seconds()),
		RefreshToken:     output.RefreshToken,
		RefreshExpiresAt: output.RefreshExpiresAt,
	})
}
