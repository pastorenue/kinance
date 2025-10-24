package document

import (
	"net/http"
	"time"

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

// GenerateReport generates a new report
func (h *Handler) GenerateReport(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req GenerateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate date range
	if req.From.After(req.To) {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "From date cannot be after To date",
		})
		return
	}

	report, err := h.service.GenerateReport(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Data:    report,
	})
}

// GenerateStatement generates a new statement
func (h *Handler) GenerateStatement(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	var req GenerateStatementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Validate date range
	if req.From.After(req.To) {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "From date cannot be after To date",
		})
		return
	}

	statement, err := h.service.GenerateStatement(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, common.APIResponse{
		Success: true,
		Data:    statement,
	})
}

// GetReport retrieves a specific report by ID
func (h *Handler) GetReport(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	reportIDStr := c.Param("id")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid report ID",
		})
		return
	}

	report, err := h.service.GetReportByID(c.Request.Context(), userID.(uuid.UUID), reportID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    report,
	})
}

// GetStatement retrieves a specific statement by ID
func (h *Handler) GetStatement(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	statementIDStr := c.Param("id")
	statementID, err := uuid.Parse(statementIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.APIResponse{
			Success: false,
			Error:   "Invalid statement ID",
		})
		return
	}

	statement, err := h.service.GetStatementByID(c.Request.Context(), userID.(uuid.UUID), statementID)
	if err != nil {
		c.JSON(http.StatusNotFound, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    statement,
	})
}

// ListReports lists all reports for the authenticated user
func (h *Handler) ListReports(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	// Parse optional query parameters for date filtering
	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if parsedFrom, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &parsedFrom
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsedTo, err := time.Parse("2006-01-02", toStr); err == nil {
			to = &parsedTo
		}
	}

	reports, err := h.service.ListReports(c.Request.Context(), userID.(uuid.UUID), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    reports,
	})
}

// ListStatements lists all statements for the authenticated user
func (h *Handler) ListStatements(c *gin.Context) {
	userID, exists := c.Get(middleware.UserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, common.APIResponse{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	// Parse optional query parameters for date filtering
	var from, to *time.Time
	if fromStr := c.Query("from"); fromStr != "" {
		if parsedFrom, err := time.Parse("2006-01-02", fromStr); err == nil {
			from = &parsedFrom
		}
	}
	if toStr := c.Query("to"); toStr != "" {
		if parsedTo, err := time.Parse("2006-01-02", toStr); err == nil {
			to = &parsedTo
		}
	}

	statements, err := h.service.ListStatements(c.Request.Context(), userID.(uuid.UUID), from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, common.APIResponse{
		Success: true,
		Data:    statements,
	})
}