package handler

import (
	"encoding/json"
	"errors"
	"hitalent-test/internal/domain"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func HandleError(w http.ResponseWriter, logger *slog.Logger, err error, requestID string) {
	var statusCode int
	var message string

	switch {
	case errors.Is(err, domain.ErrQuestionNotFound),
		errors.Is(err, domain.ErrAnswerNotFound):
		statusCode = http.StatusNotFound
		message = err.Error()
	case errors.Is(err, domain.ErrInvalidInput):
		statusCode = http.StatusBadRequest
		message = err.Error()
	default:
		statusCode = http.StatusInternalServerError
		message = "internal server error"
		logger.Error("internal error",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
		)
	}

	respondJSON(w, statusCode, ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	})
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
