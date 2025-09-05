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

	"github.com/gin-gonic/gin"
	redoc "github.com/mvrilo/go-redoc"
	"github.com/pastorenue/kinance/internal/auth"
	"github.com/pastorenue/kinance/internal/budget"
	"github.com/pastorenue/kinance/internal/receipt"
	"github.com/pastorenue/kinance/internal/transaction"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
	"github.com/pastorenue/kinance/pkg/database"
	"github.com/pastorenue/kinance/pkg/logger"
	"github.com/pastorenue/kinance/pkg/middleware"
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
	transactionService := transaction.NewService(db)
	receiptService := receipt.NewService(db, cfg.AI, logger)

	// Initialize router
	router := setupRouter(cfg, authService, userService, budgetService, transactionService, receiptService)

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

func setupRouter(cfg *config.Config, authSvc *auth.Service, userSvc *user.Service,
	budgetSvc *budget.Service, transSvc *transaction.Service,
	receiptSvc *receipt.Service) *gin.Engine {
	router := gin.New()

	// Global middleware
	// router.Use(middleware.Logger())
	// router.Use(middleware.Recovery())
	// router.Use(middleware.CORS())
	// router.Use(middleware.RateLimit())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// Serve OpenAPI spec file
	router.StaticFile("/openapi.yaml", "api/docs/openapi.yaml")

	// Serve Redoc UI
	redocHandler := redoc.Redoc{
		Title:    "Kinance API Docs",
		SpecFile: "api/docs/openapi.yaml",
		SpecPath: "/openapi.yaml",
	}
	router.GET("/docs", gin.WrapH(redocHandler.Handler()))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		authHandler := auth.NewHandler(authSvc, userSvc)
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(func(token string) (interface{}, error) {
			return authSvc.ValidateToken(token)
		}))
		{
			// User routes
			userHandler := user.NewHandler(userSvc)
			protected.GET("/users/profile", userHandler.GetProfile)
			protected.PUT("/users/profile", userHandler.UpdateProfile)
			protected.GET("/users/family", userHandler.GetFamilyMembers)
			protected.POST("/users/address", userHandler.CreateAddress)
			protected.PUT("/users/address", userHandler.UpdateAddress)

			// Budget routes
			budgetHandler := budget.NewHandler(budgetSvc)
			protected.GET("/budgets", budgetHandler.GetBudgets)
			protected.POST("/budgets", budgetHandler.CreateBudget)
			protected.PUT("/budgets/:id", budgetHandler.UpdateBudget)
			protected.DELETE("/budgets/:id", budgetHandler.DeleteBudget)

			// Transaction routes
			transHandler := transaction.NewHandler(transSvc)
			protected.GET("/transactions", transHandler.GetTransactions)
			protected.POST("/transactions", transHandler.CreateTransaction)
			protected.PUT("/transactions/:id", transHandler.UpdateTransaction)
			protected.DELETE("/transactions/:id", transHandler.DeleteTransaction)

			// Receipt routes
			receiptHandler := receipt.NewHandler(receiptSvc)
			protected.POST("/receipts/upload", receiptHandler.UploadReceipt)
			protected.GET("/receipts", receiptHandler.GetReceipts)
			protected.GET("/receipts/:id", receiptHandler.GetReceipt)
		}
	}

	return router
}
