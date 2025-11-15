package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"hitalent-test/internal/config"
	"hitalent-test/internal/handler"
	"hitalent-test/internal/repository"
	"hitalent-test/internal/server"
	"hitalent-test/internal/service"
	"hitalent-test/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	appLogger := logger.New(cfg.Logger.Level, cfg.Logger.Format)

	db, err := setupDatabase(cfg.Database, appLogger)
	if err != nil {
		appLogger.Error("Failed to connect to database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	userRepo := repository.NewUserRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	answerRepo := repository.NewAnswerRepository(db)

	tokenService := service.NewTokenService(&cfg.JWT)
	refreshTokenStore := service.NewRefreshTokenStore()

	authService := service.NewAuthService(userRepo, tokenService, refreshTokenStore)
	questionService := service.NewQuestionService(questionRepo)
	answerService := service.NewAnswerService(answerRepo, questionRepo)

	questionHandler := handler.NewQuestionHandler(questionService, appLogger)
	answerHandler := handler.NewAnswerHandler(answerService, appLogger)
	authHandler := handler.NewAuthHandler(authService, appLogger)

	router := server.NewRouter(
		questionHandler,
		answerHandler,
		authHandler,
		tokenService,
		appLogger,
	)

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			refreshTokenStore.CleanupExpired()
		}
	}()

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLogger.Info("Starting server",
			slog.String("host", cfg.Server.Host),
			slog.Int("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}

	appLogger.Info("Server exited")
}

func setupDatabase(cfg config.DatabaseConfig, logger *slog.Logger) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	}

	db, err := gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logger.Info("Database connection established")

	return db, nil
}
