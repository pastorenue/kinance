package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/user"
)

type Handler struct {
	authService *Service
	userService *user.Service
}

func NewHandler(authService *Service, userService *user.Service) *Handler {
	return &Handler{
		authService: authService,
		userService: userService,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	token, err := h.authService.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	result := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    result,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	user, err := h.userService.GetUserByID(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	_, err = h.authService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "Invalid or expired refresh token",
		})
		return
	}

	accessToken, err := h.authService.generateAccessToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to generate access token",
		})
		return
	}

	newRefreshToken, err := h.authService.generateRefreshToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: map[string]interface{}{
			"token":         accessToken,
			"refresh_token": newRefreshToken,
			"user":          user,
		},
	})
}
