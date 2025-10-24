package document

import (
	"time"

	"github.com/pastorenue/kinance/internal/common"
)

type DocumentStatus string

const (
	StatusProcessing DocumentStatus = "processing"
	StatusCompleted  DocumentStatus = "completed"
	StatusFailed     DocumentStatus = "failed"
)

type ReportType string

const (
	ReportTypeExpenseSummary ReportType = "expense_summary"
	ReportTypeIncomeReport   ReportType = "income_report"
	ReportTypeBudgetReport   ReportType = "budget_report"
	ReportTypeTaxReport      ReportType = "tax_report"
	ReportTypeInvestmentReport ReportType = "investment_report"
	ReportTypeCustomReport   ReportType = "custom_report"
)

type Statement struct {
	common.BaseModel
	UserID     string    `gorm:"index" json:"user_id"`
	Description string    `json:"description"`
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	FileURL    string    `json:"file_url"`
	Status     DocumentStatus `json:"status"`
}

type Report struct {
	common.BaseModel
	UserID     string    `gorm:"index" json:"user_id"`
	Title      string    `json:"title"`
	GeneratedAt time.Time `json:"generated_at"`
	FileURL    string    `json:"file_url"`
	Status     DocumentStatus `json:"status"`
	ReportType ReportType `json:"report_type"`
}

// Request types
type GenerateReportRequest struct {
	Title      string     `json:"title" binding:"required"`
	ReportType ReportType `json:"report_type" binding:"required"`
	From       time.Time  `json:"from" binding:"required"`
	To         time.Time  `json:"to" binding:"required"`
}

type GenerateStatementRequest struct {
	Description string    `json:"description" binding:"required"`
	From        time.Time `json:"from" binding:"required"`
	To          time.Time `json:"to" binding:"required"`
}