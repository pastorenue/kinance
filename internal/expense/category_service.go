package expense

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CategoryService struct {
	db     *gorm.DB
	logger common.Logger
}

func NewCategoryService(db *gorm.DB, logger common.Logger) *CategoryService {
	return &CategoryService{db: db, logger: logger}
}

func (s *CategoryService) CreateCategory(ctx context.Context, userID uuid.UUID, req *CreateCategoryRequest) (*Category, error) {
	if req.BudgetLimit != nil && req.BudgetLimit.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("budget limit must be greater than zero")
	}

	category := &Category{
		Name:             req.Name,
		UserID:           userID,
		ParentCategoryID: req.ParentCategoryID,
		ColorCode:        req.ColorCode,
		BudgetLimit:      req.BudgetLimit,
	}

	category.ID = uuid.New()
	if err := s.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *CategoryService) GetCategories(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	var categories []Category
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req *UpdateCategoryRequest) (*Category, error) {
	var category Category
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	if req.Name != nil {
		category.Name = *req.Name
	}
	if req.ParentCategoryID != nil {
		category.ParentCategoryID = req.ParentCategoryID
	}
	if req.ColorCode != nil {
		category.ColorCode = *req.ColorCode
	}
	if req.BudgetLimit != nil {
		if req.BudgetLimit.LessThanOrEqual(decimal.Zero) {
			return nil, errors.New("budget limit must be greater than zero")
		}
		category.BudgetLimit = req.BudgetLimit
	}

	if err := s.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).Delete(&Category{}).Error; err != nil {
		return err
	}
	return nil
}
