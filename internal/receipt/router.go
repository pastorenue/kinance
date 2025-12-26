package receipt

import "github.com/gin-gonic/gin"

func RegisterRoutes(versionedGroup *gin.RouterGroup, svc *Service) {
	receiptHandler := NewHandler(svc)
	protected := versionedGroup.Group("/receipts")
	protected.POST("/upload", receiptHandler.UploadReceipt)
	protected.GET("", receiptHandler.GetReceipts)
	protected.GET("/:id", receiptHandler.GetReceipt)
}
