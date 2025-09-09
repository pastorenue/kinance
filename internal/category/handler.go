package category

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
	return &Handler{service: service}
}

func (h *Handler) CreateCategory(c *gin.Context) {
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

func (h *Handler) GetCategories(c *gin.Context) {
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

func (h *Handler) GetCategoryByID(c *gin.Context) {
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

func (h *Handler) UpdateCategory(c *gin.Context) {
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

func (h *Handler) DeleteCategory(c *gin.Context) {
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


// getUserAndCategoryID extracts userID and categoryID from context and returns them, or writes an error response and returns false.
func getUserAndCategoryID(c *gin.Context) (uuid.UUID, uuid.UUID, bool) {
	userIDVal, _ := c.Get(middleware.UserIDKey)
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(500, common.APIResponse{
			Success: false,
			Error:   "Invalid user ID",
		})
		return uuid.Nil, uuid.Nil, false
	}
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, common.APIResponse{
			Success: false,
			Error:   "Invalid category ID",
		})
		return uuid.Nil, uuid.Nil, false
	}
	return userID, categoryID, true
}

// bindJSONOrAbort binds JSON and handles errors in a unified way.
func bindJSONOrAbort(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(400, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return false
	}
	return true
}
