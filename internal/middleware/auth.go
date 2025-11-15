package middleware

import (
	"context"
	"errors"
	"hitalent-test/internal/handler"
	"hitalent-test/internal/service"
	"log/slog"
	"net/http"
	"strings"
)

func Auth(tokenService *service.TokenService, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Context().Value("request_id").(string)

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				err := errors.New("missing authorization header")
				handler.HandleError(w, logger, err, requestID)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				err := errors.New("invalid authorization header format")
				handler.HandleError(w, logger, err, requestID)
				return
			}

			token := parts[1]
			claims, err := tokenService.VerifyToken(token)
			if err != nil {
				logger.Warn("invalid token",
					slog.String("request_id", requestID),
					slog.String("error", err.Error()),
				)
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
			ctx = context.WithValue(ctx, "user_email", claims.Email)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
