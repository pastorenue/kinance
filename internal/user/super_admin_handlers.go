package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pastorenue/kinance/internal/common"
)

type SuperHandler struct {
	service *Service
}

func NewSuperHandler(service *Service) *SuperHandler {
	return &SuperHandler{
		service: service,
	}
}


func (h *SuperHandler) GetAllUsers(c *gin.Context) {
	// Check if the user is a super admin
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "Unauthorized",
		})
		return
	}

	u, ok := user.(*User)
	if !(ok && u.IsSuperAdmin) {
		c.JSON(http.StatusForbidden, common.APIResponse{
			Success: false,
			Error:   "Forbidden: super admin access required",
		})
		return
	}

	users, err := h.service.ListAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    users,
	})
}
