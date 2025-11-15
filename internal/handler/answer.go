package handler

import (
	"encoding/json"
	"hitalent-test/internal/domain"
	"hitalent-test/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type AnswerHandler struct {
	service *service.AnswerService
	logger  *slog.Logger
}

func NewAnswerHandler(service *service.AnswerService, logger *slog.Logger) *AnswerHandler {
	return &AnswerHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AnswerHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)
	userID := r.Context().Value("user_id").(string)

	idStr := r.PathValue("id")
	questionID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	var req domain.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	req.UserID = userID

	answer, err := h.service.Create(uint(questionID), &req)
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusCreated, answer)
}

func (h *AnswerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	answer, err := h.service.GetByID(uint(id))
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusOK, answer)
}

func (h *AnswerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
