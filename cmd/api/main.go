package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pastorenue/kinance/internal/auth"
	"github.com/pastorenue/kinance/internal/app"
	"github.com/pastorenue/kinance/internal/budget"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/expense"
	"github.com/pastorenue/kinance/internal/income"
	"github.com/pastorenue/kinance/internal/receipt"
	"github.com/pastorenue/kinance/internal/repository"
	"github.com/pastorenue/kinance/internal/transaction"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
	"github.com/pastorenue/kinance/pkg/database"
	"github.com/pastorenue/kinance/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize services
	userService := user.NewService(db, logger)
	authService := auth.NewService(db, cfg.JWT, logger)
	budgetService := budget.NewService(db, logger)
	receiptService := receipt.NewService(db, cfg.AI, logger)
	expenseService := expense.NewService(db, logger)
	categoryService := category.NewService(db, logger)
	transactionService := transaction.NewService(db, logger)
	incomeService := income.NewService(db, logger)

	// Initialize OAuth and token repository
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	tokenRepo := repository.NewTokenRepository(redisAddr, cfg.Redis.Password, cfg.Redis.DB)
	oauthSvc := auth.NewOAuthService(cfg, tokenRepo)
	authHandler := auth.NewHandler(authService, userService)
	oauthHandler := auth.NewOAuthHandler(oauthSvc, userService)
	googleHandler := auth.NewGoogleHandler(cfg, tokenRepo, userService, authService)

	// Initialize router
	router := app.SetupRouter(
		cfg,
		authService,
		userService,
		budgetService,
		transactionService,
		receiptService,
		expenseService,
		categoryService,
		incomeService,
		oauthHandler,
		googleHandler,
		authHandler,
		tokenRepo,
	)

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting server on port %d", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start:", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Info("Server exited")
}
