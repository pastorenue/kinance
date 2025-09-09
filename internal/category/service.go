package category

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
	return &Service{db: db, logger: logger}
}

func (s *Service) CreateCategory(ctx context.Context, userID uuid.UUID, req *CreateCategoryRequest) (*Category, error) {

	category := &Category{
		Name:             req.Name,
		UserID:           userID,
		ParentCategoryID: req.ParentCategoryID,
		ColorCode:        req.ColorCode,
	}

	category.ID = uuid.New()
	if err := s.db.WithContext(ctx).Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) GetCategories(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	var categories []Category
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *Service) GetCategoryByID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) (*Category, error) {
	var category Category
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (s *Service) UpdateCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID, req *UpdateCategoryRequest) (*Category, error) {
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

	if err := s.db.WithContext(ctx).Save(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *Service) DeleteCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", categoryID, userID).Delete(&Category{}).Error; err != nil {
		return err
	}
	return nil
}
