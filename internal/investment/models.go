package investment

import (
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
)

type Investment struct {
	common.BaseModel
	UserID   uuid.UUID `json:"user_id" gorm:"not null;index"`
	FamilyID *uuid.UUID `json:"family_id" gorm:"index"`
	Name     string    `json:"name" gorm:"not null"`
	Type     string    `json:"type" gorm:"not null"`
	Amount   float64   `json:"amount" gorm:"not null"`
	Date     time.Time `json:"date" gorm:"not null"`
	Status   string    `json:"status" gorm:"not null"`
}

type CreateInvestmentRequest struct {
	Name   string    `json:"name" binding:"required"`
	Type   string    `json:"type" binding:"required"`
	Amount float64   `json:"amount" binding:"required,gt=0"`
	Date   time.Time `json:"date" binding:"required"`
	Status string    `json:"status" binding:"required"`
}

type UpdateInvestmentRequest struct {
	Name   string    `json:"name"`
	Type   string    `json:"type"`
	Amount float64   `json:"amount" binding:"gt=0"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
}
