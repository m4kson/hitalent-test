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
	authService *service.AuthService,
	tokenService *service.TokenService,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	questionHandler := handler.NewQuestionHandler(questionService, logger)
	answerHandler := handler.NewAnswerHandler(answerService, logger)
	authHandler := handler.NewAuthHandler(authService, logger)

	authMiddleware := middleware.Auth(tokenService, logger)

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)

	mux.HandleFunc("GET /questions/", questionHandler.GetAll)
	mux.HandleFunc("POST /questions/", questionHandler.Create)
	mux.HandleFunc("GET /questions/{id}", questionHandler.GetByID)
	mux.HandleFunc("DELETE /questions/{id}", questionHandler.Delete)

	mux.HandleFunc("POST /questions/{id}/answers/",
		authMiddleware(http.HandlerFunc(answerHandler.Create)).ServeHTTP)

	mux.HandleFunc("GET /answers/{id}", answerHandler.GetByID)

	mux.HandleFunc("DELETE /answers/{id}",
		authMiddleware(http.HandlerFunc(answerHandler.Delete)).ServeHTTP)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	var h http.Handler = mux
	h = middleware.Logger(logger)(h)

	return h
}
