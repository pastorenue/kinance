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

	var pagination common.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = common.PaginationParams{Page: 1, PageSize: 20}
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

func (h *ExpenseHandler) CreateRecurringExpense(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateRecurringExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusBadRequest,
			Error:     err.Error(),
		})
		return
	}

	result, err := h.service.CreateRecurringExpense(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusInternalServerError,
			Error:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusCreated,
		Data:      result,
	})
}

func (h *ExpenseHandler) GetRecurringExpenses(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var pagination common.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = common.PaginationParams{Page: 1, PageSize: 20}
	}

result, err := h.service.GetRecurringExpenses(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusInternalServerError,
			Error:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusOK,
		Data:      result,
	})
}

func (h *ExpenseHandler) GetRecurringExpenseByID(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	recurringExpenseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusBadRequest,
			Error:     "Invalid recurring expense ID",
		})
		return
	}

	result, err := h.service.GetRecurringExpenseByID(c.Request.Context(), userID.(uuid.UUID), recurringExpenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusInternalServerError,
			Error:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusOK,
		Data:	  result,
	})
}

func (h *ExpenseHandler) UpdateRecurringExpense(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	recurringExpenseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusBadRequest,
			Error:     "Invalid recurring expense ID",
		})
		return
	}

	var req UpdateRecurringExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusBadRequest,
			Error:     err.Error(),
		})
		return
	}

	result, err := h.service.UpdateRecurringExpense(c.Request.Context(), userID.(uuid.UUID), recurringExpenseID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusInternalServerError,
			Error:     err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusOK,
		Data:      result,
	})
}

func (h *ExpenseHandler) DeleteRecurringExpense(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	recurringExpenseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusBadRequest,
			Error:     "Invalid recurring expense ID",
		})
		return
	}

	err = h.service.DeleteRecurringExpense(c.Request.Context(), userID.(uuid.UUID), recurringExpenseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			StatusCode: http.StatusInternalServerError,
			Error:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusOK,
		Data:      "Recurring expense deleted successfully",
	})
}