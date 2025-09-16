package transaction

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/expense"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	logger common.Logger
}

func NewService(db *gorm.DB, logger common.Logger) *Service {
	return &Service{db: db, logger: logger}
}

func (s *Service) CreateExpenseTransaction(
	ctx context.Context,
	userID uuid.UUID,
	req *CreateTransactionRequest,
) (*ExpenseTransactionResponse, error) {
	if req.CategoryID == uuid.Nil {
		return nil, fmt.Errorf("category ID cannot be empty")
	}

	transaction := &Transaction{
		UserID:          userID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		TransactionDate: req.TransactionDate,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		Metadata:        req.Metadata,
		PaymentMethod:   req.PaymentMethod,
	}
	transaction.ID = uuid.New()

	err := s.db.WithContext(ctx).Transaction((func(tx *gorm.DB) error {
		if err := tx.Create(transaction).Error; err != nil {
			s.logger.Error("Failed to create transaction", "error", err)
			return err
		} else {
			s.logger.Info("Transaction created successfully", "transaction_id", transaction.ID)
		}

		expense := &expense.Expense{
			UserID:        userID,
			Amount:        req.Amount,
			Currency:      req.Currency,
			TransactionID: &transaction.ID,
			Description:   req.Description,
			CategoryID:    req.CategoryID,
			PaymentMethod: req.PaymentMethod,
		}
		expense.ID = uuid.New()
		if err := tx.Create(expense).Error; err != nil {
			tx.Rollback()
			s.logger.Error("Failed to create expense", "error", err)
			return err
		} else {
			s.logger.Info("Expense created successfully", "expense_id", expense.ID)
		}

		transaction.ProcessingObjectID = &expense.ID
		if err := tx.Save(transaction).Error; err != nil {
			tx.Rollback()
			s.logger.Error("Failed to link expense to transaction", "error", err)
			return err
		}
		return nil
	}))

	if err != nil {
		return nil, err
	}

	// Preload Category before returning the transaction
	if err := s.db.Preload("Category").First(transaction, transaction.ID).Error; err != nil {
		s.logger.Error("Failed to preload category", "error", err)
		return nil, err
	}

	// Use the processing object ID to fetch the linked expense
	var linkedExpense expense.Expense
	if transaction.ProcessingObjectID != nil {
		if err := s.db.First(&linkedExpense, "id = ?", *transaction.ProcessingObjectID).Error; err != nil {
			s.logger.Error("Failed to fetch linked expense", "error", err)
			return nil, err
		}
	}

	// Create a response struct that includes both transaction and linked expense if needed
	response := &ExpenseTransactionResponse{
		StatusCode: http.StatusOK,
		Message:    "Success",
		Transaction: *transaction,
		Expense:     linkedExpense,
	}

	return response, nil
}

func (s *Service) GetTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	var transactions []Transaction
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&transactions).Error; err != nil {
		s.logger.Error("Failed to fetch transactions", "error", err)
		return nil, err
	}
	return transactions, nil
}

func (s *Service) GetTransactionByID(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*Transaction, error) {
	var transaction Transaction
	if err := s.db.WithContext(ctx).Where("user_id = ? AND id = ?", userID, transactionID).First(&transaction).Error; err != nil {
		s.logger.Error("Failed to fetch transaction", "error", err)
		return nil, err
	}
	return &transaction, nil
}

func (s *Service) CreateIncomeTransaction(ctx context.Context, userID uuid.UUID, req *CreateTransactionRequest) (*Transaction, error) {
	if req.CategoryID == uuid.Nil {
		return nil, fmt.Errorf("category ID cannot be empty")
	}

	transaction := &Transaction{
		UserID:          userID,
		Type:            req.Type,
		Amount:          req.Amount,
		Currency:        req.Currency,
		TransactionDate: req.TransactionDate,
		CategoryID:      req.CategoryID,
		Description:     req.Description,
		Metadata:        req.Metadata,
		PaymentMethod:   req.PaymentMethod,
	}
	transaction.ID = uuid.New()

	if err := s.db.WithContext(ctx).Create(transaction).Error; err != nil {
		s.logger.Error("Failed to create income transaction", "error", err)
		return nil, err
	}
	return transaction, nil
}

func (s *Service) CreateTransferTransaction(ctx context.Context) {}
func (s *Service) LinkExpenseToTransaction(ctx context.Context)  {}
func (s *Service) LinkIncomeToTransaction(ctx context.Context)   {}
func (s *Service) LinkTransferToTransaction(ctx context.Context) {}
