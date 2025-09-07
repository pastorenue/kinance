package expense

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/pkg/middleware"
)

type CategoryHandler struct {
	service *CategoryService
}

func NewCategoryHandler(service *CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	category, err := h.service.CreateCategory(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success:   true,
		StatusCode: http.StatusCreated,
		Data:      category,
	})
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	categories, err := h.service.GetCategories(c.Request.Context(), userID.(uuid.UUID))
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
		Data:      categories,
		StatusCode: http.StatusOK,
	})
}

func (h *CategoryHandler) GetCategoryByID(c *gin.Context) {
	userID, categoryID, ok := getUserAndCategoryID(c)
	if !ok {
		return
	}

	category, err := h.service.GetCategoryByID(c.Request.Context(), userID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			Error:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    category,
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID, categoryID, ok := getUserAndCategoryID(c)
	if !ok {
		return
	}

	var req UpdateCategoryRequest
	if !bindJSONOrAbort(c, &req) {
		return
	}

	category, err := h.service.UpdateCategory(c.Request.Context(), userID, categoryID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			Error:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      category,
		StatusCode: http.StatusOK,
	})
}

func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID, categoryID, ok := getUserAndCategoryID(c)
	if !ok {
		return
	}

	err := h.service.DeleteCategory(c.Request.Context(), userID, categoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success:   false,
			Error:     err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      "Category deleted successfully",
		StatusCode: http.StatusOK,
	})
}
