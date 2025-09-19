package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/shopspring/decimal"
)

type TransactionType string

const (
	TypeIncome   TransactionType = "income"
	TypeExpense  TransactionType = "expense"
	TypeTransfer TransactionType = "transfer"
)

type Merchant struct {
	common.BaseModel
	Name    string    `json:"name" gorm:"not null;uniqueIndex"`
	Website string    `json:"website"`
	LogoURL string    `json:"logo_url"`
	UserID  uuid.UUID `json:"user_id" gorm:"not null;index"`
}

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusCanceled  TransactionStatus = "canceled"
)

type Transaction struct {
	common.BaseModel
	UserID               uuid.UUID              `json:"user_id" gorm:"not null;index"`
	Amount               decimal.Decimal        `json:"amount" gorm:"not null"`
	Description          string                 `json:"description"`
	CategoryID           uuid.UUID              `json:"category_id" gorm:"index"`
	Category             *category.Category     `json:"category" gorm:"foreignKey:CategoryID"`
	TransactionDate      time.Time              `json:"transaction_date" gorm:"not null"`
	Status               TransactionStatus      `json:"status" gorm:"default:completed"`
	Tags                 []Tag                  `json:"tags" gorm:"many2many:transaction_tags;"`
	Type                 TransactionType        `json:"type" gorm:"not null"`
	ProcessingObjectID   *uuid.UUID             `json:"processing_object_id" gorm:"uniqueIndex"`
	Currency             *common.Currency       `json:"currency" gorm:"not null;default:EUR"`
	ExcludeFromAnalytics bool                   `json:"exclude_from_analytics" gorm:"default:false"`
	MerchantID           *uuid.UUID             `json:"merchant" gorm:"index"`
	Merchant             *Merchant              `json:"merchant_details" gorm:"foreignKey:MerchantID"`
	ReceiptID            *uuid.UUID             `json:"receipt" gorm:"index"`
	Metadata             map[string]interface{} `json:"metadata" gorm:"type:json;serializer:json"`
	PaymentMethod        common.PaymentMethod   `json:"payment_method" gorm:"type:payment_method"`
}

type Tag struct {
	common.BaseModel
	Name   string    `json:"name" gorm:"uniqueIndex;not null"`
	Color  string    `json:"color" gorm:"default:#007bff"`
	UserID uuid.UUID `json:"user_id" gorm:"not null"`
}

type CreateTransactionRequest struct {
	Amount          decimal.Decimal        `json:"amount" binding:"required"`
	Description     string                 `json:"description" binding:"required"`
	CategoryID      uuid.UUID              `json:"category_id" binding:"required"`
	Merchant        string                 `json:"merchant"`
	TransactionDate time.Time              `json:"transaction_date" binding:"required"`
	Type            TransactionType        `json:"type" binding:"required"`
	Tags            []string               `json:"tags"`
	Currency        common.Currency        `json:"currency" binding:"required,oneof=USD EUR GBP JPY CHF NGN"`
	PaymentMethod   common.PaymentMethod   `json:"payment_method" binding:"required,oneof=cash card bank_transfer"`
	Metadata        map[string]interface{} `json:"metadata"`
}

type CreateIncomeTransactionRequest struct {
	CreateTransactionRequest
	SwiftCode string `json:"swift_code" binding:"required,len=8|len=11"`
}

type TransactionResponse struct {
	Transaction Transaction `json:"transaction"`
	Entity      interface{} `json:"entity"`
}
