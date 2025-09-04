package budget

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger Logger
}

type Logger interface {
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

func NewService(db *gorm.DB, logger Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) CreateBudget(ctx context.Context, userID uuid.UUID, req *CreateBudgetRequest) (*Budget, error) {
	budget := &Budget{
		UserID:         userID,
		Name:           req.Name,
		Description:    req.Description,
		Amount:         req.Amount,
		Category:       req.Category,
		Period:         req.Period,
		AlertThreshold: req.AlertThreshold,
		IsActive:       true,
	}

	if err := s.db.WithContext(ctx).Create(budget).Error; err != nil {
		return nil, err
	}

	s.logger.Info("Budget created successfully", "budget_id", budget.ID, "user_id", userID)
	return budget, nil
}

func (s *Service) GetBudgets(ctx context.Context, userID uuid.UUID, pagination *common.PaginationParams) (*common.PaginatedResponse, error) {
	var budgets []Budget
	var total int64

	query := s.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true)

	// Count total records
	if err := query.Model(&Budget{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).Find(&budgets).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pagination.PageSize
	if int(total)%pagination.PageSize > 0 {
		totalPages++
	}

	return &common.PaginatedResponse{
		Data: budgets,
		Pagination: common.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) UpdateBudget(ctx context.Context, userID, budgetID uuid.UUID, req *UpdateBudgetRequest) (*Budget, error) {
	var budget Budget
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).First(&budget).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("budget not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		budget.Name = req.Name
	}
	if req.Description != "" {
		budget.Description = req.Description
	}
	if req.Amount > 0 {
		budget.Amount = req.Amount
	}
	if req.Category != "" {
		budget.Category = req.Category
	}
	if req.Period != "" {
		budget.Period = req.Period
	}
	if req.AlertThreshold > 0 {
		budget.AlertThreshold = req.AlertThreshold
	}
	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	if err := s.db.WithContext(ctx).Save(&budget).Error; err != nil {
		return nil, err
	}

	return &budget, nil
}

func (s *Service) DeleteBudget(ctx context.Context, userID, budgetID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", budgetID, userID).Delete(&Budget{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("budget not found")
	}

	s.logger.Info("Budget deleted successfully", "budget_id", budgetID, "user_id", userID)
	return nil
}

func (s *Service) UpdateSpentAmount(ctx context.Context, budgetID uuid.UUID, amount float64) error {
	return s.db.WithContext(ctx).Model(&Budget{}).
		Where("id = ?", budgetID).
		Update("spent", gorm.Expr("spent + ?", amount)).Error
}
