package expense

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/pkg/middleware"
)

type ExpenseHandler struct {
	service *ExpenseService
}

func NewExpenseHandler(service *ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: service}
}


func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	expense, err := h.service.CreateExpense(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Data:    expense,
	})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var pagination common.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = common.PaginationParams{Page: 1, PageSize: 20}
	}

	result, err := h.service.GetExpenses(c.Request.Context(), userID.(uuid.UUID), &pagination)
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

func (h *ExpenseHandler) GetExpenseByID(c *gin.Context) {
	userID, expenseID, ok := getUserAndExpenseID(c)
	if !ok {
		return
	}

	expense, err := h.service.GetExpenseByID(c.Request.Context(), userID, expenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    expense,
	})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	userID, expenseID, ok := getUserAndExpenseID(c)
	if !ok {
		return
	}

	err := h.service.DeleteExpense(c.Request.Context(), userID, expenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    "Expense deleted successfully",
	})
}

func (h *ExpenseHandler) UpdateExpense(c *gin.Context) {
	userID, expenseID, ok := getUserAndExpenseID(c)
	if !ok {
		return
	}

	var req UpdateExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	expense, err := h.service.UpdateExpense(c.Request.Context(), userID, expenseID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    expense,
	})
}

func (h *ExpenseHandler) GetExpensesByCategoryID(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return
	}

	expenses, err := h.service.GetExpensesByCategoryID(c.Request.Context(), userID.(uuid.UUID), categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    expenses,
	})
}
