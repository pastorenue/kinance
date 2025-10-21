package investment

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger common.Logger
}

func NewService(db *gorm.DB, logger common.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) CreateInvestment(ctx context.Context, userID uuid.UUID, req *CreateInvestmentRequest) (*Investment, error) {
	investment := &Investment{
		UserID: userID,
		Name:   req.Name,
		Type:   req.Type,
		Amount: req.Amount,
		Date:   req.Date,
		Status: req.Status,
	}

	if err := s.db.WithContext(ctx).Create(investment).Error; err != nil {
		return nil, err
	}

	s.logger.Info("Investment created successfully", "investment_id", investment.ID, "user_id", userID)
	return investment, nil
}

func (s *Service) GetInvestments(ctx context.Context, userID uuid.UUID, pagination *common.PaginationParams) (*common.PaginatedResponse, error) {
	var investments []Investment
	var total int64

	query := s.db.WithContext(ctx).Where("user_id = ?", userID)

	// Count total records
	if err := query.Model(&Investment{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	if err := query.Offset(offset).Limit(pagination.PageSize).Find(&investments).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / pagination.PageSize
	if int(total)%pagination.PageSize > 0 {
		totalPages++
	}

	return &common.PaginatedResponse{
		Data: investments,
		Pagination: common.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetInvestment(ctx context.Context, userID, investmentID uuid.UUID) (*Investment, error) {
	var investment Investment
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("investment not found")
		}
		return nil, err
	}
	return &investment, nil
}

func (s *Service) UpdateInvestment(ctx context.Context, userID, investmentID uuid.UUID, req *UpdateInvestmentRequest) (*Investment, error) {
	var investment Investment
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).First(&investment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("investment not found")
		}
		return nil, err
	}

	// Update fields
	if req.Name != "" {
		investment.Name = req.Name
	}
	if req.Type != "" {
		investment.Type = req.Type
	}
	if req.Amount > 0 {
		investment.Amount = req.Amount
	}
	if !req.Date.IsZero() {
		investment.Date = req.Date
	}
	if req.Status != "" {
		investment.Status = req.Status
	}

	if err := s.db.WithContext(ctx).Save(&investment).Error; err != nil {
		return nil, err
	}

	return &investment, nil
}

func (s *Service) DeleteInvestment(ctx context.Context, userID, investmentID uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", investmentID, userID).Delete(&Investment{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("investment not found")
	}

	s.logger.Info("Investment deleted successfully", "investment_id", investmentID, "user_id", userID)
	return nil
}
