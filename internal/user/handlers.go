package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/pkg/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	user, err := h.service.GetUserByID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Error:   "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    user,
	})
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    user,
	})
}

func (h *Handler) GetFamilyMembers(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	members, err := h.service.GetFamilyMembers(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    members,
	})
}
