package transaction

import "github.com/gin-gonic/gin"

// ...existing code...
type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTransactions(c *gin.Context)           {}
func (h *Handler) CreateTransaction(c *gin.Context)         {}
func (h *Handler) UpdateTransaction(c *gin.Context)         {}
func (h *Handler) DeleteTransaction(c *gin.Context)         {}
func (h *Handler) CreateExpenseTransaction(c *gin.Context)  {}
func (h *Handler) CreateIncomeTransaction(c *gin.Context)   {}
func (h *Handler) CreateTransferTransaction(c *gin.Context) {}
func (h *Handler) GetTransaction(c *gin.Context)            {}
func (h *Handler) LinkExpenseToTransaction(c *gin.Context)   {}
func (h *Handler) LinkIncomeToTransaction(c *gin.Context)    {}
func (h *Handler) LinkTransferToTransaction(c *gin.Context)  {}