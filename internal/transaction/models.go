package transaction

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Transaction struct {
	common.BaseModel
	UserID          uuid.UUID              `json:"user_id" gorm:"not null;index"`
	Amount          float64                `json:"amount" gorm:"not null"`
	Description     string                 `json:"description" gorm:"not null"`
	Category        string                 `json:"category" gorm:"not null;index"`
	TransactionDate time.Time              `json:"transaction_date" gorm:"not null"`
	Type            TransactionType        `json:"type" gorm:"not null"`
	Status          TransactionStatus      `json:"status" gorm:"default:completed"`
	Tags            []Tag                  `json:"tags" gorm:"many2many:transaction_tags;"`
	Metadata        map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	TransferID      *uuid.UUID             `json:"transfer_id" gorm:"index"`
	ExpenseID       *uuid.UUID             `json:"expense_id" gorm:"index"`
	IncomeID        *uuid.UUID             `json:"income_id" gorm:"index"`
}

type TransactionType string

const (
	TypeIncome  TransactionType = "income"
	TypeExpense TransactionType = "expense"
	TypeTransfer TransactionType = "transfer"
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
