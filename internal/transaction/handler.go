package transaction

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/pkg/middleware"
)

// ...existing code...
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListTransactions(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)

	transactions, err := h.service.GetTransactions(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, transactions)
}

func (h *Handler) CreateTransaction(c *gin.Context)         {}
func (h *Handler) UpdateTransaction(c *gin.Context)         {}
func (h *Handler) DeleteTransaction(c *gin.Context)         {}
func (h *Handler) CreateExpenseTransaction(c *gin.Context)  {
	userID, _ := c.Get(middleware.UserIDKey)
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	trnxResponse, err := h.service.CreateExpenseTransaction(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, trnxResponse)
}

func (h *Handler) CreateIncomeTransaction(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	var req CreateIncomeTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	trnxResponse, err := h.service.CreateIncomeTransaction(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(201, trnxResponse)
}

func (h *Handler) CreateTransferTransaction(c *gin.Context) {}
func (h *Handler) GetTransaction(c *gin.Context) {
	userID, _ := c.Get(middleware.UserIDKey)
	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid transaction ID"})
		return
	}
	
	transaction, err := h.service.GetTransactionByID(c.Request.Context(), userID.(uuid.UUID), transactionID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, transaction)
}
func (h *Handler) LinkExpenseToTransaction(c *gin.Context)   {}
func (h *Handler) LinkIncomeToTransaction(c *gin.Context)    {}
func (h *Handler) LinkTransferToTransaction(c *gin.Context)  {}