package document

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"gorm.io/gorm"
)

// StorageProvider defines the interface for document storage (local or cloud)
type StorageProvider interface {
	UploadFile(ctx context.Context, userID uuid.UUID, fileName string, content []byte) (string, error)
	GetFile(ctx context.Context, userID uuid.UUID, fileName string) ([]byte, error)
	GetFileURL(ctx context.Context, userID uuid.UUID, fileName string) (string, error)
}

// LocalStorageProvider implements StorageProvider for local file system storage
type LocalStorageProvider struct {
	BaseDir string
}

func NewLocalStorageProvider(baseDir string) *LocalStorageProvider {
	return &LocalStorageProvider{BaseDir: baseDir}
}

func (l *LocalStorageProvider) UploadFile(ctx context.Context, userID uuid.UUID, fileName string, content []byte) (string, error) {
	userDir := filepath.Join(l.BaseDir, userID.String())
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user directory: %w", err)
	}

	filePath := filepath.Join(userDir, fileName)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return "", fmt.Errorf("failed to write file to local storage: %w", err)
	}
	return filePath, nil // Return local path as URL for now
}

func (l *LocalStorageProvider) GetFile(ctx context.Context, userID uuid.UUID, fileName string) ([]byte, error) {
	filePath := filepath.Join(l.BaseDir, userID.String(), fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file from local storage: %w", err)
	}
	return content, nil
}

func (l *LocalStorageProvider) GetFileURL(ctx context.Context, userID uuid.UUID, fileName string) (string, error) {
	return filepath.Join(l.BaseDir, userID.String(), fileName), nil
}

type Service struct {
	db      *gorm.DB
	logger  common.Logger
	storage StorageProvider
}

func NewService(db *gorm.DB, logger common.Logger, storage StorageProvider) *Service {
	return &Service{
		db:      db,
		logger:  logger,
		storage: storage,
	}
}

// GenerateReport generates a report document (placeholder for now)
func (s *Service) GenerateReport(ctx context.Context, userID uuid.UUID, req *GenerateReportRequest) (*Report, error) {
	report := &Report{
		UserID:     userID.String(),
		Title:      req.Title,
		GeneratedAt: time.Now(),
		Status:     StatusProcessing,
		ReportType: req.ReportType,
	}
	report.ID = uuid.New()

	if err := s.db.WithContext(ctx).Create(report).Error; err != nil {
		return nil, fmt.Errorf("failed to create report record: %w", err)
	}

	// Simulate PDF generation and storage
	go func() {
		// In a real scenario, this would involve fetching data, generating PDF, and uploading
		// For now, create a dummy PDF content
		dummyContent := []byte(fmt.Sprintf("Report: %s\nType: %s\nFrom: %s\nTo: %s\nGenerated for User: %s",
			req.Title, req.ReportType, req.From.Format("2006-01-02"), req.To.Format("2006-01-02"), userID.String()))

		fileName := fmt.Sprintf("report_%s_%s.pdf", report.ID.String(), time.Now().Format("20060102150405"))
		fileURL, err := s.storage.UploadFile(ctx, userID, fileName, dummyContent)
		if err != nil {
			s.logger.Error("Failed to upload report file", "error", err, "report_id", report.ID)
			report.Status = StatusFailed
		} else {
			report.FileURL = fileURL
			report.Status = StatusCompleted
		}

		if err := s.db.WithContext(ctx).Save(report).Error; err != nil {
			s.logger.Error("Failed to update report status/file URL", "error", err, "report_id", report.ID)
		}
	}()

	return report, nil
}

// GenerateStatement generates a statement document (placeholder for now)
func (s *Service) GenerateStatement(ctx context.Context, userID uuid.UUID, req *GenerateStatementRequest) (*Statement, error) {
	statement := &Statement{
		UserID:     userID.String(),
		Description: req.Description,
		From:       req.From,
		To:         req.To,
		Status:     StatusProcessing,
	}
	statement.ID = uuid.New()

	if err := s.db.WithContext(ctx).Create(statement).Error; err != nil {
		return nil, fmt.Errorf("failed to create statement record: %w", err)
	}

	// Simulate PDF generation and storage
	go func() {
		// In a real scenario, this would involve fetching data, generating PDF, and uploading
		// For now, create a dummy PDF content
		dummyContent := []byte(fmt.Sprintf("Statement: %s\nFrom: %s\nTo: %s\nGenerated for User: %s",
			req.Description, req.From.Format("2006-01-02"), req.To.Format("2006-01-02"), userID.String()))

		fileName := fmt.Sprintf("statement_%s_%s.pdf", statement.ID.String(), time.Now().Format("20060102150405"))
		fileURL, err := s.storage.UploadFile(ctx, userID, fileName, dummyContent)
		if err != nil {
			s.logger.Error("Failed to upload statement file", "error", err, "statement_id", statement.ID)
			statement.Status = StatusFailed
		} else {
			statement.FileURL = fileURL
			statement.Status = StatusCompleted
		}

		if err := s.db.WithContext(ctx).Save(statement).Error; err != nil {
			s.logger.Error("Failed to update statement status/file URL", "error", err, "statement_id", statement.ID)
		}
	}()

	return statement, nil
}

// GetReportByID retrieves a specific report by its ID
func (s *Service) GetReportByID(ctx context.Context, userID uuid.UUID, reportID uuid.UUID) (*Report, error) {
	var report Report
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", reportID, userID.String()).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("report not found")
		}
		return nil, fmt.Errorf("failed to retrieve report: %w", err)
	}
	return &report, nil
}

// GetStatementByID retrieves a specific statement by its ID
func (s *Service) GetStatementByID(ctx context.Context, userID uuid.UUID, statementID uuid.UUID) (*Statement, error) {
	var statement Statement
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", statementID, userID.String()).First(&statement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("statement not found")
		}
		return nil, fmt.Errorf("failed to retrieve statement: %w", err)
	}
	return &statement, nil
}

// ListReports lists all reports for a user, with optional date range filtering
func (s *Service) ListReports(ctx context.Context, userID uuid.UUID, from, to *time.Time) ([]Report, error) {
	var reports []Report
	query := s.db.WithContext(ctx).Where("user_id = ?", userID.String())

	if from != nil {
		query = query.Where("generated_at >= ?", *from)
	}
	if to != nil {
		query = query.Where("generated_at <= ?", *to)
	}

	if err := query.Order("generated_at DESC").Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("failed to list reports: %w", err)
	}
	return reports, nil
}

// ListStatements lists all statements for a user, with optional date range filtering
func (s *Service) ListStatements(ctx context.Context, userID uuid.UUID, from, to *time.Time) ([]Statement, error) {
	var statements []Statement
	query := s.db.WithContext(ctx).Where("user_id = ?", userID.String())

	if from != nil {
		query = query.Where("from >= ?", *from)
	}
	if to != nil {
		query = query.Where("to <= ?", *to)
	}

	if err := query.Order("created_at DESC").Find(&statements).Error; err != nil {
		return nil, fmt.Errorf("failed to list statements: %w", err)
	}
	return statements, nil
}
