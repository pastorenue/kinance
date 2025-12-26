package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/repository"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
)


func RegisterRoutes(
	versionedGroup *gin.RouterGroup,
	tokenRepo *repository.TokenRepository,
	oauthHandler *OAuthHandler,
	googleHandler *GoogleHandler,
	authHandler *Handler,
	userSvc *user.Service,
	cfg *config.Config,
) {
	versionedGroup.POST("/auth/register", authHandler.Register)
	versionedGroup.POST("/auth/login", authHandler.Login)
	versionedGroup.POST("/auth/refresh", authHandler.RefreshToken)

	// OAuth2 Authorization Server style endpoints
	versionedGroup.GET("/oauth/authorize", gin.WrapF(oauthHandler.Authorize))
	versionedGroup.POST("/oauth/token", gin.WrapF(oauthHandler.Token))
	versionedGroup.POST("/oauth/introspect", gin.WrapF(oauthHandler.Introspect))
	versionedGroup.POST("/oauth/revoke", gin.WrapF(oauthHandler.Revoke))

	// Google OAuth2 login
	versionedGroup.GET("/auth/google/login", googleHandler.Login)
	versionedGroup.GET("/auth/google/callback", googleHandler.Callback)
}