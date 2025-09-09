package expense

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/shopspring/decimal"
)

type PaymentMethod string

const (
	Cash   PaymentMethod = "cash"
	Card   PaymentMethod = "card"
	BankTransfer PaymentMethod = "bank_transfer"
)

type RecurringFrequency string

const (
	Daily   RecurringFrequency = "daily"
	Weekly  RecurringFrequency = "weekly"
	Monthly RecurringFrequency = "monthly"
	Yearly  RecurringFrequency = "yearly"
)
type RecurringExpense struct {
	common.BaseModel
	Amount        decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`
	Description   string          `json:"description"`
	CategoryID    uuid.UUID      `gorm:"not null" json:"category_id"`
	Category      *category.Category        `json:"category"`
	UserID        uuid.UUID      `gorm:"not null" json:"user_id"`
	Frequency     RecurringFrequency `gorm:"type:recurring_frequency;not null" json:"frequency"`
	PaymentMethod PaymentMethod   `gorm:"type:payment_method" json:"payment_method"`
	IsActive	  bool            `gorm:"default:true" json:"is_active"`
	Expenses      []Expense       `gorm:"foreignKey:RecurringExpenseID" json:"expenses"`
	StartDate     time.Time       `gorm:"not null" json:"start_date"`
	EndDate       *time.Time      `json:"end_date,omitempty"`
	NextDueDate   time.Time       `gorm:"not null" json:"next_due_date_parsed"`
	LastProcessed time.Time       `json:"last_processed,omitempty"`
}

type Expense struct {
	common.BaseModel
	Amount      decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount" binding:"required"`
	Description string            `json:"description" binding:"required"`
	CategoryID  uuid.UUID      `gorm:"not null" json:"category_id" binding:"required"`
	Category    *category.Category        `gorm:"foreignKey:CategoryID" json:"category"`
	UserID      uuid.UUID      `gorm:"not null" json:"user_id" binding:"required"`
	RecurringExpenseID *uuid.UUID        `gorm:"index" json:"recurring_expense_id,omitempty"`
	RecurringExpense   *RecurringExpense `gorm:"foreignKey:RecurringExpenseID" json:"recurring_expense,omitempty"`
	ReceiptURL  string            `json:"receipt_url,omitempty"`
	PaymentMethod PaymentMethod   `gorm:"type:payment_method" json:"payment_method"`
}

// Request payloads

type CreateExpenseRequest struct {
	Amount            decimal.Decimal `json:"amount" binding:"required"`
	Description       string          `json:"description"`
	CategoryID        uuid.UUID      `json:"category_id" binding:"required"`
	PaymentMethod     PaymentMethod   `json:"payment_method" binding:"required,oneof=cash card bank_transfer"`
	RecurringExpenseID *uuid.UUID      `json:"recurring_expense_id,omitempty"`
	ReceiptURL       string          `json:"receipt_url,omitempty"`
}


type UpdateExpenseRequest struct {
	Amount        *decimal.Decimal `json:"amount" binding:"omitempty"`
	Description   *string          `json:"description"`
	CategoryID    *uuid.UUID      `json:"category_id"`
	PaymentMethod *PaymentMethod   `json:"payment_method" binding:"omitempty,oneof=cash card bank_transfer"`
	ReceiptURL   *string          `json:"receipt_url,omitempty"`
}

type CreateRecurringExpenseRequest struct {
	Amount        decimal.Decimal    `json:"amount" binding:"required"`
	Description   string             `json:"description"`
	CategoryID    uuid.UUID          `json:"category_id" binding:"required"`
	Frequency     RecurringFrequency `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	PaymentMethod PaymentMethod      `json:"payment_method" binding:"required,oneof=cash card bank_transfer"`
	StartDate     time.Time          `json:"start_date" binding:"required"`
	EndDate	   *time.Time         `json:"end_date,omitempty"`
}

type UpdateRecurringExpenseRequest struct {
	Amount        *decimal.Decimal    `json:"amount" binding:"omitempty"`
	Description   *string             `json:"description"`
	CategoryID    *uuid.UUID          `json:"category_id"`
	Frequency     *RecurringFrequency `json:"frequency" binding:"omitempty,oneof=daily weekly monthly yearly"`
	PaymentMethod *PaymentMethod      `json:"payment_method" binding:"omitempty,oneof=cash card bank_transfer"`
	StartDate     *time.Time          `json:"start_date"`
	EndDate	   *time.Time         `json:"end_date,omitempty"`
	IsActive	  *bool               `json:"is_active"`
}

type RecurringExpenseResponse struct {
	RecurringExpense RecurringExpense `json:"recurring_expense"`
	DaysUntilDue    int              `json:"days_until_due"`
}
