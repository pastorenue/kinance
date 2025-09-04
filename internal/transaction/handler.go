package transaction

import "github.com/gin-gonic/gin"

// ...existing code...

func (h *Handler) GetTransactions(c *gin.Context)   {}
func (h *Handler) CreateTransaction(c *gin.Context) {}
func (h *Handler) UpdateTransaction(c *gin.Context) {}
func (h *Handler) DeleteTransaction(c *gin.Context) {}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
