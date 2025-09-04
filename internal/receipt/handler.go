package receipt

import "github.com/gin-gonic/gin"

// ...existing code...

func (h *Handler) UploadReceipt(c *gin.Context) {}
func (h *Handler) GetReceipts(c *gin.Context)   {}
func (h *Handler) GetReceipt(c *gin.Context)    {}

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}
