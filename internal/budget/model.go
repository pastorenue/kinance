package budget

import (
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Budget struct {
	common.BaseModel
	UserID         uuid.UUID  `json:"user_id" gorm:"not null;index"`
	FamilyID       *uuid.UUID `json:"family_id" gorm:"index"`
	Name           string     `json:"name" gorm:"not null"`
	Description    string     `json:"description"`
	Amount         float64    `json:"amount" gorm:"not null"`
	Spent          float64    `json:"spent" gorm:"default:0"`
	Category       string     `json:"category" gorm:"not null"`
	Period         Period     `json:"period" gorm:"default:monthly"`
	IsActive       bool       `json:"is_active" gorm:"default:true"`
	AlertThreshold float64    `json:"alert_threshold" gorm:"default:80"` // Alert when 80% spent
}

type Period string

const (
	PeriodWeekly  Period = "weekly"
	PeriodMonthly Period = "monthly"
	PeriodYearly  Period = "yearly"
)

type CreateBudgetRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Category       string  `json:"category" binding:"required"`
	Period         Period  `json:"period" binding:"required"`
	AlertThreshold float64 `json:"alert_threshold" binding:"min=0,max=100"`
}

type UpdateBudgetRequest struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Amount         float64 `json:"amount" binding:"gt=0"`
	Category       string  `json:"category"`
	Period         Period  `json:"period"`
	AlertThreshold float64 `json:"alert_threshold" binding:"min=0,max=100"`
	IsActive       *bool   `json:"is_active"`
}
