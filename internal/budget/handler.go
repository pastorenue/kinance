package budget

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

func (h *Handler) CreateBudget(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	budget, err := h.service.CreateBudget(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Data:    budget,
	})
}

func (h *Handler) GetBudgets(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var pagination common.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = common.PaginationParams{Page: 1, PageSize: 20}
	}

	result, err := h.service.GetBudgets(c.Request.Context(), userID.(uuid.UUID), &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    result,
	})
}

func (h *Handler) UpdateBudget(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid budget ID",
		})
		return
	}

	var req UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	budget, err := h.service.UpdateBudget(c.Request.Context(), userID.(uuid.UUID), budgetID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    budget,
	})
}

func (h *Handler) DeleteBudget(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	budgetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid budget ID",
		})
		return
	}

	if err := h.service.DeleteBudget(c.Request.Context(), userID.(uuid.UUID), budgetID); err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Message: "Budget deleted successfully",
	})
}
