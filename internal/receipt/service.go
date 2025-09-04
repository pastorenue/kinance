package receipt

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/pkg/config"
	"gorm.io/gorm"
)

type Service struct {
	db       *gorm.DB
	aiConfig config.AIConfig
	logger   Logger
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

func NewService(db *gorm.DB, aiConfig config.AIConfig, logger Logger) *Service {
	return &Service{
		db:       db,
		aiConfig: aiConfig,
		logger:   logger,
	}
}

func (s *Service) ProcessReceipt(ctx context.Context, userID uuid.UUID, req *UploadReceiptRequest) (*Receipt, error) {
	// Create receipt record
	receipt := &Receipt{
		UserID:           userID,
		OriginalImageURL: "placeholder", // Will be updated after file upload
		ProcessingStatus: StatusPendingProcessing,
	}

	if err := s.db.WithContext(ctx).Create(receipt).Error; err != nil {
		return nil, err
	}

	// Process image with OCR service asynchronously
	go s.processReceiptAsync(receipt.ID, req.Image)

	return receipt, nil
}

func (s *Service) processReceiptAsync(receiptID uuid.UUID, imageData string) {
	// Update status to processing
	s.db.Model(&Receipt{}).Where("id = ?", receiptID).Update("processing_status", StatusProcessing)

	// Decode base64 image
	imageBytes, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		s.logger.Error("Failed to decode image", "receipt_id", receiptID, "error", err)
		s.updateReceiptStatus(receiptID, StatusFailed)
		return
	}

	// Call OCR service
	ocrResult, err := s.callOCRService(imageBytes)
	if err != nil {
		s.logger.Error("OCR processing failed", "receipt_id", receiptID, "error", err)
		s.updateReceiptStatus(receiptID, StatusFailed)
		return
	}

	// Parse OCR result and update receipt
	if err := s.updateReceiptWithOCRResult(receiptID, ocrResult); err != nil {
		s.logger.Error("Failed to update receipt with OCR result", "receipt_id", receiptID, "error", err)
		s.updateReceiptStatus(receiptID, StatusFailed)
		return
	}

	s.updateReceiptStatus(receiptID, StatusProcessed)
	s.logger.Info("Receipt processed successfully", "receipt_id", receiptID)
}

type OCRResult struct {
	Merchant   string                 `json:"merchant"`
	Total      float64                `json:"total"`
	Tax        float64                `json:"tax"`
	Items      []OCRItem              `json:"items"`
	Confidence float64                `json:"confidence"`
	RawData    map[string]interface{} `json:"raw_data"`
}

type OCRItem struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	Barcode    string  `json:"barcode,omitempty"`
}

func (s *Service) callOCRService(imageData []byte) (*OCRResult, error) {
	// Prepare request payload
	payload := map[string]interface{}{
		"image": base64.StdEncoding.EncodeToString(imageData),
		"type":  "receipt",
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Make HTTP request to OCR service
	req, err := http.NewRequest("POST", s.aiConfig.OCRServiceURL+"/process", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.aiConfig.OCRAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OCR service returned status %d", resp.StatusCode)
	}

	var result OCRResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *Service) updateReceiptWithOCRResult(receiptID uuid.UUID, ocrResult *OCRResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update receipt with OCR data
		updates := map[string]interface{}{
			"merchant":   ocrResult.Merchant,
			"total":      ocrResult.Total,
			"tax":        ocrResult.Tax,
			"confidence": ocrResult.Confidence,
			"ocr_data":   ocrResult.RawData,
		}

		if err := tx.Model(&Receipt{}).Where("id = ?", receiptID).Updates(updates).Error; err != nil {
			return err
		}

		// Create receipt items
		for _, item := range ocrResult.Items {
			receiptItem := &ReceiptItem{
				ReceiptID:  receiptID,
				Name:       item.Name,
				Quantity:   item.Quantity,
				UnitPrice:  item.UnitPrice,
				TotalPrice: item.TotalPrice,
				Barcode:    item.Barcode,
			}

			if err := tx.Create(receiptItem).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *Service) updateReceiptStatus(receiptID uuid.UUID, status ProcessingStatus) {
	s.db.Model(&Receipt{}).Where("id = ?", receiptID).Update("processing_status", status)
}

func (s *Service) GetReceipts(ctx context.Context, userID uuid.UUID) ([]Receipt, error) {
	var receipts []Receipt
	if err := s.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Find(&receipts).Error; err != nil {
		return nil, err
	}
	return receipts, nil
}

func (s *Service) GetReceiptByID(ctx context.Context, userID, receiptID uuid.UUID) (*Receipt, error) {
	var receipt Receipt
	if err := s.db.WithContext(ctx).Preload("Items").Where("id = ? AND user_id = ?", receiptID, userID).First(&receipt).Error; err != nil {
		return nil, err
	}
	return &receipt, nil
}
