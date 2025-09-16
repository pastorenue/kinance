package income

import (
	"net/http"

	"strings"

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

func (h *Handler) CreateIncome(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	var req CreateIncomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid request",
			StatusCode: http.StatusBadRequest,
		})
		return
	}
	income, err := h.service.CreateIncome(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to create income",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	c.JSON(http.StatusCreated, common.APIResponse{
		Success:   true,
		Data:     income,
		StatusCode: http.StatusCreated,
	})
}

func (h *Handler) GetIncomes(c *gin.Context)   {
	userID, _ := c.Get(middleware.UserIDKey)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}
	var pagination common.PaginationParams
	if err := c.ShouldBindQuery(&pagination); err != nil {
		pagination = common.PaginationParams{Page: 1, PageSize: 20}
	}

	result, err := h.service.GetIncomes(c.Request.Context(), userID.(uuid.UUID), &pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to fetch incomes",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      result,
		StatusCode: http.StatusOK,
	})
}

func (h *Handler) GetIncomeByID(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	if userID == nil {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
			StatusCode: http.StatusUnauthorized,
		})
		return
	}

	incomeIDParam := c.Param("id")
	incomeID, err := uuid.Parse(incomeIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid income ID",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	income, err := h.service.GetIncomeByID(c.Request.Context(), userID.(uuid.UUID), incomeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to fetch income",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      income,
		StatusCode: http.StatusOK,
	})
}

func (h *Handler) UpdateIncome(c *gin.Context) {}
func (h *Handler) DeleteIncome(c *gin.Context) {}

// Sources
func (h *Handler) CreateSource(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	var req CreateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid request",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	source, err := h.service.CreateSource(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		h.service.logger.Error("Failed to create source", "error", err)
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to create source",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success:   true,
		Data:      source,
		StatusCode: http.StatusCreated,
	})
}

func (h *Handler) GetSources(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	result, err := h.service.GetSources(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		h.service.logger.Error("Failed to fetch sources", "error", err)
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to fetch sources",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      result,
		StatusCode: http.StatusOK,
	})
}

func (h *Handler) GetSourceBySwiftCode(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	swiftCode := strings.ToUpper(c.Param("swift_code"))
	// Validate swiftCode format
	if !isValidSwiftCode(swiftCode) {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid SWIFT code format",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	source, err := h.service.GetSourceBySwiftCode(c.Request.Context(), userID.(uuid.UUID), swiftCode)
	if err != nil {
		h.service.logger.Error("Failed to fetch source", "error", err)
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to fetch source",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	if source == nil {
		c.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Error:   "Source not found",
			StatusCode: http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      source,
		StatusCode: http.StatusOK,
	})
}

func isValidSwiftCode(swiftCode string) bool {
	if swiftCode != "" {
		if len(swiftCode) < 8 || len(swiftCode) > 11 {
			return false
		}
	} else {
		return false
	}
	return true
}

func (h *Handler) UpdateSource(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	// Get swift_code from URL param
	swiftCode := strings.ToUpper(c.Param("swift_code"))
	var req UpdateSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.service.logger.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid request",
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	source, err := h.service.UpdateSource(c.Request.Context(), userID.(uuid.UUID), swiftCode, &req)
	if err != nil {
		h.service.logger.Error("Failed to update source", "error", err)
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   "Failed to update source",
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success:   true,
		Data:      source,
		StatusCode: http.StatusOK,
	})
}