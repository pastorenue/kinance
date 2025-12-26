package budget

import "github.com/gin-gonic/gin"


func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	budgetHandler := NewHandler(svc)
	protected := versionedGroup.Group("/budgets")
	protected.GET("/", budgetHandler.GetBudgets)
	protected.POST("/", budgetHandler.CreateBudget)
	protected.GET("/:id", budgetHandler.GetBudget)
	protected.PUT("/:id", budgetHandler.UpdateBudget)
	protected.DELETE("/:id", budgetHandler.DeleteBudget)
}