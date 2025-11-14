package handler

import (
	"encoding/json"
	"hitalent-test/internal/domain"
	"hitalent-test/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

type QuestionHandler struct {
	service *service.QuestionService
	logger  *slog.Logger
}

func NewQuestionHandler(service *service.QuestionService, logger *slog.Logger) *QuestionHandler {
	return &QuestionHandler{
		service: service,
		logger:  logger,
	}
}

func (h *QuestionHandler) Create(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	var req domain.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	question, err := h.service.Create(&req)
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusCreated, question)
}

func (h *QuestionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		HandleError(w, h.logger, domain.ErrInvalidInput, requestID)
		return
	}

	question, err := h.service.GetByID(uint(id))
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusOK, question)
}

func (h *QuestionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	requestID := r.Context().Value("request_id").(string)

	questions, err := h.service.GetAll()
	if err != nil {
		HandleError(w, h.logger, err, requestID)
		return
	}

	respondJSON(w, http.StatusOK, questions)
}

func (h *QuestionHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
