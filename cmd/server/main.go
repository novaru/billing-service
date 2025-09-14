package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/novaru/billing-service/db/generated"
	"github.com/novaru/billing-service/internal/app/handler"
	"github.com/novaru/billing-service/internal/app/repository"
	"github.com/novaru/billing-service/internal/app/service"
	"github.com/novaru/billing-service/internal/config"
	"github.com/novaru/billing-service/internal/database"
	"github.com/novaru/billing-service/internal/router"
	"github.com/novaru/billing-service/pkg/logger"
)

func main() {
	cfg := config.Load()

	// Initialize logger
	if err := logger.Initialize(cfg.Env); err != nil {
		panic(err)
	}

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	q := generated.New(db.Pool)

	// Initialize repositories
	userRepo := repository.NewUserRepository(q)

	// Initialize services
	userService := service.NewUserService(cfg, userRepo)

	// Initialize handlers
	handlers := handler.New(
		userService,
	)

	// Setup router
	r := router.New(cfg, handlers).Setup()

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprint(":", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Server starting on " + cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited gracefully")
}
