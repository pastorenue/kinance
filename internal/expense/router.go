package expense

import "github.com/gin-gonic/gin"

func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	expHandler := NewHandler(svc)
	protected := versionedGroup.Group("/expenses")
	// Expense routes
	protected.POST("/", expHandler.CreateExpense)
	protected.GET("/", expHandler.GetExpenses)
	protected.GET("/:id", expHandler.GetExpenseByID)
	protected.PUT("/:id", expHandler.UpdateExpense)
	protected.DELETE("/:id", expHandler.DeleteExpense)
	protected.GET("/category/:id", expHandler.GetExpensesByCategoryID)

	// Recurring Expense routes
	protected.POST("/recurring", expHandler.CreateRecurringExpense)
	protected.GET("/recurring", expHandler.GetRecurringExpenses)
	protected.GET("/recurring/:id", expHandler.GetRecurringExpenseByID)
	protected.PUT("/recurring/:id", expHandler.UpdateRecurringExpense)
	protected.DELETE("/recurring/:id", expHandler.DeleteRecurringExpense)
	protected.GET("/recurring/:id/history", expHandler.GetRecurringExpenseHistory)
	protected.GET("/recurring/upcoming", expHandler.GetUpcomingRecurringExpenses)
}
