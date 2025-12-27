package receipt

import (
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Receipt struct {
	common.BaseModel
	UserID            uuid.UUID              `json:"user_id" gorm:"not null;index"`
	TransactionID     *uuid.UUID             `json:"transaction_id"`
	OriginalImageURL  string                 `json:"original_image_url" gorm:"not null"`
	ProcessedImageURL string                 `json:"processed_image_url"`
	Merchant          string                 `json:"merchant"`
	Total             float64                `json:"total"`
	Tax               float64                `json:"tax"`
	ProcessingStatus  ProcessingStatus       `json:"processing_status" gorm:"default:pending"`
	OCRData           map[string]interface{} `json:"ocr_data" gorm:"type:jsonb"`
	Items             []ReceiptItem          `json:"items"`
	Confidence        float64                `json:"confidence"` // OCR confidence score
}

type ReceiptItem struct {
	common.BaseModel
	ReceiptID  uuid.UUID `json:"receipt_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	Quantity   int       `json:"quantity" gorm:"default:1"`
	UnitPrice  float64   `json:"unit_price"`
	TotalPrice float64   `json:"total_price"`
	Category   string    `json:"category"`
	Barcode    string    `json:"barcode"`
}

type ProcessingStatus string

const (
	StatusPendingProcessing ProcessingStatus = "pending"
	StatusProcessing        ProcessingStatus = "processing"
	StatusProcessed         ProcessingStatus = "processed"
	StatusFailed            ProcessingStatus = "failed"
)

type UploadReceiptRequest struct {
	Image       string `json:"image" binding:"required"` // Base64 encoded image
	Description string `json:"description"`
}
