package server

import (
	"hitalent-test/internal/handler"
	"hitalent-test/internal/middleware"
	"hitalent-test/internal/service"
	"log/slog"
	"net/http"
)

func NewRouter(
	questionService *service.QuestionService,
	answerService *service.AnswerService,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	questionHandler := handler.NewQuestionHandler(questionService, logger)
	answerHandler := handler.NewAnswerHandler(answerService, logger)

	mux.HandleFunc("GET /questions/", questionHandler.GetAll)
	mux.HandleFunc("POST /questions/", questionHandler.Create)
	mux.HandleFunc("GET /questions/{id}", questionHandler.GetByID)
	mux.HandleFunc("DELETE /questions/{id}", questionHandler.Delete)

	mux.HandleFunc("POST /questions/{id}/answers/", answerHandler.Create)
	mux.HandleFunc("GET /answers/{id}", answerHandler.GetByID)
	mux.HandleFunc("DELETE /answers/{id}", answerHandler.Delete)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	var h http.Handler = mux
	h = middleware.Logger(logger)(h)

	return h
}
