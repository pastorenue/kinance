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

func NewHandler(authService *Service) *Handler {
	return &Handler{
		authService: authService,
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

	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
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
	// Implementation for refresh token logic
	c.JSON(http.StatusNotImplemented, common.APIResponse{
		Success: false,
		Message: "Refresh token not implemented yet",
	})
}
