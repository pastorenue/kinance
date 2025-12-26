package app

import (
	"time"

	"github.com/gin-gonic/gin"
	redoc "github.com/mvrilo/go-redoc"
	"github.com/pastorenue/kinance/internal/auth"
	"github.com/pastorenue/kinance/internal/budget"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/expense"
	"github.com/pastorenue/kinance/internal/income"
	"github.com/pastorenue/kinance/internal/receipt"
	"github.com/pastorenue/kinance/internal/repository"
	"github.com/pastorenue/kinance/internal/transaction"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
	"github.com/pastorenue/kinance/pkg/middleware"
)

func SetupRouter(
	cfg *config.Config,
	authSvc *auth.Service,
	userSvc *user.Service,
	budgetSvc *budget.Service,
	tranxSvc *transaction.Service,
	receiptSvc *receipt.Service,
	expenseSvc *expense.Service,
	categorySvc *category.Service,
	incomeSvc *income.Service,
	oauthHandler *auth.OAuthHandler,
	googleHandler *auth.GoogleHandler,
	authHandler *auth.Handler,
	tokenRepo *repository.TokenRepository,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.Trace())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(cfg.MiddlewareConf.RateLimit))

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
		auth.RegisterRoutes(
			v1,
			tokenRepo,
			oauthHandler,
			googleHandler,
			authHandler,
			userSvc,
			cfg,
		)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthRequired(func(token string) (any, error) {
			return authSvc.ValidateToken(token)
		}))
		{
			user.RegisterRoutes(protected, userSvc)
			budget.RegisterRoutes(protected, budgetSvc)
			transaction.RegisterRoutes(protected, tranxSvc)
			receipt.RegisterRoutes(protected, receiptSvc)
			expense.RegisterRoutes(protected, expenseSvc)
			category.RegisterRoutes(protected, categorySvc)
			income.RegisterRoutes(protected, incomeSvc)
		}
	}

	return router
}