package transaction

import "github.com/gin-gonic/gin"

func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	transHandler := NewHandler(svc)
	protected := versionedGroup.Group("/transaction")
	protected.GET("/", transHandler.ListTransactions)
	protected.GET("/:id", transHandler.GetTransaction)
	protected.POST("/", transHandler.CreateTransaction)
	protected.PUT("/:id", transHandler.UpdateTransaction)
	protected.DELETE("/:id", transHandler.DeleteTransaction)
	protected.POST("/expense", transHandler.CreateExpenseTransaction)
	protected.POST("/income", transHandler.CreateIncomeTransaction)
}