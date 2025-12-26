package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	catHandler := NewHandler(svc)
	protected := versionedGroup.Group("/categories")
	protected.POST("/", catHandler.CreateCategory)
	protected.GET("/", catHandler.GetCategories)
	protected.GET("/:id", catHandler.GetCategoryByID)
	protected.PUT("/:id", catHandler.UpdateCategory)
	protected.DELETE("/:id", catHandler.DeleteCategory)
}