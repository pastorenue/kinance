package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Transaction struct {
	common.BaseModel
	UserID          uuid.UUID              `json:"user_id" gorm:"not null;index"`
	BudgetID        *uuid.UUID             `json:"budget_id" gorm:"index"`
	ReceiptID       *uuid.UUID             `json:"receipt_id" gorm:"index"`
	Amount          float64                `json:"amount" gorm:"not null"`
	Description     string                 `json:"description" gorm:"not null"`
	Category        string                 `json:"category" gorm:"not null;index"`
	Merchant        string                 `json:"merchant"`
	TransactionDate time.Time              `json:"transaction_date" gorm:"not null"`
	Type            TransactionType        `json:"type" gorm:"not null"`
	Status          TransactionStatus      `json:"status" gorm:"default:completed"`
	Tags            []Tag                  `json:"tags" gorm:"many2many:transaction_tags;"`
	Metadata        map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
}

type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusCanceled  TransactionStatus = "canceled"
)

type Tag struct {
	common.BaseModel
	Name   string    `json:"name" gorm:"uniqueIndex;not null"`
	Color  string    `json:"color" gorm:"default:#007bff"`
	UserID uuid.UUID `json:"user_id" gorm:"not null"`
}

type CreateTransactionRequest struct {
	Amount          float64         `json:"amount" binding:"required"`
	Description     string          `json:"description" binding:"required"`
	Category        string          `json:"category" binding:"required"`
	Merchant        string          `json:"merchant"`
	TransactionDate string          `json:"transaction_date" binding:"required"`
	Type            TransactionType `json:"type" binding:"required"`
	Tags            []string        `json:"tags"`
}
