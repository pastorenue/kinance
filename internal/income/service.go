package income

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/user"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger common.Logger
}

func NewService(db *gorm.DB, logger common.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) CreateIncome(ctx context.Context, userID uuid.UUID, req *CreateIncomeRequest) (*Income, error) {
	var createdIncome *Income

	// Check if Source exists using swift code
	var source Source
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("swift_code = ?", req.SwiftCode).First(&source).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Create new source if not found
				source = Source{
					Name:      req.SwiftCode, // Should be updated to a proper name
					SwiftCode: req.SwiftCode,
				}
				source.ID = uuid.New()
				if err := tx.Create(&source).Error; err != nil {
					s.logger.Error("Failed to create source", "error", err)
					return err
				}
			} else {
				s.logger.Error("Failed to retrieve source", "error", err)
				return err
			}
		}

		income := &Income{
			Amount:   req.Amount,
			SourceID: source.ID,
			UserID:   userID,
			Source:   &source, // Assign the source directly
		}

		income.ID = uuid.New()
		if err := tx.Create(income).Error; err != nil {
			tx.Rollback()
			s.logger.Error("Failed to create income", "error", err)
			return err
		}

		s.logger.Info("Income created successfully", "income_id", income.ID)
		createdIncome = income
		return nil
	}); err != nil {
		return nil, err
	}
	return createdIncome, nil
}

func (s *Service) GetIncomes(
	ctx context.Context,
	userID uuid.UUID,
	pagination *common.PaginationParams,
) (*common.PaginatedResponse, error) {
	var incomes []Income
	if err := s.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("Source"). // Preload source when fetching incomes
		Limit(pagination.PageSize).
		Offset((pagination.Page - 1) * pagination.PageSize).
		Find(&incomes).Error; err != nil {
		return nil, err
	}

	// Preloading is now done in the initial query
	// No need for separate preloading

	totalPages := int(len(incomes)) / int(pagination.PageSize)
	if int(len(incomes))%int(pagination.PageSize) > 0 {
		totalPages++
	}

	return &common.PaginatedResponse{
		Data: incomes,
		Pagination: common.Pagination{
			Page:       pagination.Page,
			PageSize:   pagination.PageSize,
			Total:      int64(len(incomes)),
			TotalPages: totalPages,
		},
	}, nil
}

func (s *Service) GetIncomeByID(ctx context.Context, userID uuid.UUID, incomeID uuid.UUID) (*Income, error) {
	var income Income
	if err := s.db.WithContext(ctx).
		Preload("Source").
		Where("id = ? AND user_id = ?", incomeID, userID).
		First(&income).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("income not found")
		}
		return nil, err
	}
	return &income, nil
}

func (s *Service) UpdateIncome(ctx context.Context, userID uuid.UUID, incomeID uuid.UUID, req *UpdateIncomeRequest) (*Income, error) {
	var income Income
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", incomeID, userID).First(&income).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("income not found")
		}
		return nil, err
	}

	if req.Status != nil {
		income.Status = *req.Status
	}
	if req.Metadata != nil {
		income.Metadata = *req.Metadata
	}
	if req.CategoryID != nil {
		income.CategoryID = req.CategoryID
	}
	if req.Note != nil {
		income.Note = *req.Note
	}

	if err := s.db.WithContext(ctx).Save(&income).Error; err != nil {
		return nil, err
	}
	return &income, nil
}

func (s *Service) DeleteIncome(ctx context.Context, userID uuid.UUID, incomeID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", incomeID, userID).Delete(&Income{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) GetIncomeByCategory(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]Income, error) {
	var incomes []Income
	if err := s.db.WithContext(ctx).Where("user_id = ? AND category_id = ?", userID, categoryID).Find(&incomes).Error; err != nil {
		return nil, err
	}
	return incomes, nil
}

// Sources
func (s *Service) CreateSource(ctx context.Context, userID uuid.UUID, req *CreateSourceRequest) (*Source, error) {
	if !isSuperAdmin(s.db, userID) {
		s.logger.Error("Unauthorized access attempt", "user_id", userID)
		return nil, errors.New("forbidden")
	}

	source := &Source{
		Name:        req.Name,
		Description: req.Description,
		SwiftCode: req.SwiftCode,
		LogoURL:    req.LogoURL,
	}

	source.ID = uuid.New()
	if err := s.db.WithContext(ctx).Create(source).Error; err != nil {
		return nil, err
	}
	return source, nil
}

func (s *Service) GetSources(ctx context.Context, userID uuid.UUID) ([]Source, error) {
	if !isSuperAdmin(s.db, userID) {
		s.logger.Error("Unauthorized access attempt", "user_id", userID)
		return nil, errors.New("forbidden")
	}

	var sources []Source
	if err := s.db.WithContext(ctx).Find(&sources).Error; err != nil {
		return nil, err
	}
	return sources, nil
}

func (s *Service) GetSourceBySwiftCode(ctx context.Context, userID uuid.UUID, swiftCode string) (*Source, error) {
	if !isSuperAdmin(s.db, userID) {
		s.logger.Error("Unauthorized access attempt", "user_id", userID)
		return nil, errors.New("forbidden")
	}

	var source Source
	if err := s.db.WithContext(ctx).Where("swift_code = ?", swiftCode).First(&source).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source record not found")
		}
		return nil, err
	}
	return &source, nil
}

func (s *Service) UpdateSource(ctx context.Context, userID uuid.UUID, swiftCode string, req *UpdateSourceRequest) (*Source, error) {
	if !isSuperAdmin(s.db, userID) {
		s.logger.Error("Unauthorized access attempt", "user_id", userID)
		return nil, errors.New("forbidden")
	}

	var source Source
	if err := s.db.
		WithContext(ctx).
		Where("swift_code = ?", swiftCode).
		First(&source).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("source not found")
		}
		return nil, err
	}

	if req.Name != nil {
		source.Name = *req.Name
	}
	if req.Description != nil {
		source.Description = *req.Description
	}
	if req.LogoURL != nil {
		source.LogoURL = *req.LogoURL
	}
	if req.IsValidated != nil {
		source.IsValidated = *req.IsValidated
	}

	if err := s.db.WithContext(ctx).Save(&source).Error; err != nil {
		return nil, err
	}
	return &source, nil
}

func isSuperAdmin(db *gorm.DB, userID uuid.UUID) bool {
	var user user.User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return false
	}
	// return user.IsSuperAdmin
	return true
}