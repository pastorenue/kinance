package user

import "github.com/gin-gonic/gin"


func RegisterRoutes(versionedGroup *gin.RouterGroup, userSvc *Service) {
	userHandler := NewHandler(userSvc)
	protected := versionedGroup.Group("/users")
	protected.GET("/profile", userHandler.GetProfile)
	protected.PUT("/profile", userHandler.UpdateProfile)
	protected.GET("/family", userHandler.GetFamilyMembers)
	protected.POST("/address", userHandler.CreateAddress)
	protected.PUT("/address", userHandler.UpdateAddress)
}