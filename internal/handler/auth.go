package handler

import (
	"encoding/json"
	"hitalent-test/internal/domain"
	"hitalent-test/internal/service"
	"log/slog"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	user, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	authResp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusOK, authResp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	var req domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	accessToken, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
