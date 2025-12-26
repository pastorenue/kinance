package income

import "github.com/gin-gonic/gin"

func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	incomeHandler := NewHandler(svc)
	versionedGroup.GET("/all-sources", incomeHandler.GetSources)

	protected := versionedGroup.Group("/incomes")
	protected.POST("/", incomeHandler.CreateIncome)
	protected.PUT("/source/:swift_code", incomeHandler.UpdateSource)
	protected.GET("/source/:swift_code", incomeHandler.GetSourceBySwiftCode)
	protected.POST("/source", incomeHandler.CreateSource)
	protected.GET("/", incomeHandler.GetIncomes)
	protected.GET("/:id", incomeHandler.GetIncomeByID)
	protected.PUT("/:id", incomeHandler.UpdateIncome)
	protected.DELETE("/:id", incomeHandler.DeleteIncome)
}