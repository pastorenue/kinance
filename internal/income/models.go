package income

import (
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/shopspring/decimal"
)

type IncomeStatus string

const (
	IncomeStatusPending   IncomeStatus = "pending"
	IncomeStatusCompleted IncomeStatus = "completed"
	IncomeStatusFailed    IncomeStatus = "failed"
)

type Source struct {
	common.BaseModel
	Name        string `json:"name" gorm:"type:varchar(255);not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	LogoURL     string `json:"logo_url" gorm:"type:varchar(255)"`
	SwiftCode   string `json:"swift_code" gorm:"type:varchar(11)"`
	IsValidated bool   `json:"is_validated" gorm:"default:false"`
}

type Income struct {
	common.BaseModel
	Amount     decimal.Decimal    `json:"amount" gorm:"type:decimal(20,8);not null"`
	SourceID   uuid.UUID          `json:"source" gorm:"index;not null"`
	Source     *Source            `json:"source_details" gorm:"foreignKey:SourceID"`
	UserID     uuid.UUID          `json:"user_id" gorm:"type:varchar(36);not null"`
	Status     IncomeStatus       `json:"status" gorm:"default:'pending'"`
	Note       string             `json:"note" gorm:"type:text"`
	Metadata   string             `json:"metadata" gorm:"type:text"`
	CategoryID *uuid.UUID         `json:"category_id" gorm:"index"`
	Category   *category.Category `json:"category" gorm:"foreignKey:CategoryID"`
}

type CreateIncomeRequest struct {
	Amount     decimal.Decimal `json:"amount" binding:"required"`
	SwiftCode  string          `json:"swift_code" binding:"omitempty,len=8|len=11"`
	Note       *string         `json:"note"`
	CategoryID *uuid.UUID      `json:"category_id"`
}

type UpdateIncomeRequest struct {
	Status     *IncomeStatus `json:"status"`
	Note       *string       `json:"note"`
	Metadata   *string       `json:"metadata"`
	CategoryID *uuid.UUID    `json:"category_id"`
}

type CreateSourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	SwiftCode   string `json:"swift_code" binding:"omitempty,len=8|len=11"`
}

type UpdateSourceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	LogoURL     *string `json:"logo_url"`
	IsValidated *bool   `json:"is_validated"`
}
