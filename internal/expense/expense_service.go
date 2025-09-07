package expense

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ExpenseService struct {
	db     *gorm.DB
	logger common.Logger
}

func NewExpenseService(db *gorm.DB, logger common.Logger) *ExpenseService {
	return &ExpenseService{db: db, logger: logger}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, userID uuid.UUID, req *CreateExpenseRequest) (*Expense, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than zero")
	}

	expense := &Expense{
		Amount:             req.Amount,
		Description:        req.Description,
		CategoryID:         req.CategoryID,
		UserID:             userID,
		PaymentMethod:      req.PaymentMethod,
		ReceiptURL:         req.ReceiptURL,
		RecurringExpenseID: req.RecurringExpenseID,
	}

	expense.ID = uuid.New()

	if err := s.db.WithContext(ctx).Create(expense).Error; err != nil {
		return nil, err
	}

	// Load the category relationship
	if err := s.db.WithContext(ctx).Preload("Category").First(expense, expense.ID).Error; err != nil {
		return nil, err
	}

	return expense, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID, req *UpdateExpenseRequest) (*Expense, error) {
	var expense Expense
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}

	if req.Amount != nil {
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			return nil, errors.New("amount must be greater than zero")
		}
		expense.Amount = *req.Amount
	}
	if req.Description != nil {
		expense.Description = *req.Description
	}
	if req.CategoryID != nil {
		expense.CategoryID = *req.CategoryID
	}
	if req.PaymentMethod != nil {
		expense.PaymentMethod = *req.PaymentMethod
	}
	if req.ReceiptURL != nil {
		expense.ReceiptURL = *req.ReceiptURL
	}

	if err := s.db.WithContext(ctx).Save(&expense).Error; err != nil {
		return nil, err
	}

	// Load the category relationship
	if err := s.db.WithContext(ctx).Preload("Category").First(&expense, expense.ID).Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).Delete(&Expense{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetExpenses(ctx context.Context, userID uuid.UUID, pagination *common.PaginationParams) ([]Expense, error) {
	var expenses []Expense

	offset := (pagination.Page - 1) * pagination.PageSize
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Category").Offset(offset).Limit(pagination.PageSize).Find(&expenses).Error; err != nil {
		return nil, err
	}

	return expenses, nil
}

func (s *ExpenseService) GetExpenseByID(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID) (*Expense, error) {
	var expense Expense
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).Preload("Category").First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}
	return &expense, nil
}

func (s *ExpenseService) GetExpensesByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.WithContext(ctx).Where("user_id = ? AND category_id = ?", userID, categoryID).Preload("Category").Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (s *ExpenseService) GetTotalExpensesByCategory(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]decimal.Decimal, error) {
	type Result struct {
		CategoryID uuid.UUID
		Total      decimal.Decimal
	}

	var results []Result
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("category_id, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("category_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalMap := make(map[uuid.UUID]decimal.Decimal)
	for _, r := range results {
		totalMap[r.CategoryID] = r.Total
	}

	return totalMap, nil
}

func (s *ExpenseService) GetTotalExpenses(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (s *ExpenseService) CreateRecurringExpense(ctx context.Context, userID uuid.UUID, req *CreateRecurringExpenseRequest) (*RecurringExpense, error) {
	if req.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("amount must be greater than zero")
	}

	recurringExpense := &RecurringExpense{
		Amount:        req.Amount,
		Description:   req.Description,
		CategoryID:    req.CategoryID,
		UserID:        userID,
		Frequency:     req.Frequency,
		PaymentMethod: req.PaymentMethod,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		NextDueDate:   req.StartDate, // Initialize with start date
	}

	recurringExpense.ID = uuid.New()

	if err := s.db.WithContext(ctx).Create(recurringExpense).Error; err != nil {
		return nil, err
	}
	return recurringExpense, nil
}

func (s *ExpenseService) GetRecurringExpenses(ctx context.Context, userID uuid.UUID) ([]RecurringExpense, error) {
	var recurringExpenses []RecurringExpense
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&recurringExpenses).Error; err != nil {
		return nil, err
	}
	return recurringExpenses, nil
}

func (s *ExpenseService) DeleteRecurringExpense(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", recurringExpenseID, userID).Delete(&RecurringExpense{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *ExpenseService) GetTotalExpensesByMonth(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
	type Result struct {
		Month string
		Total decimal.Decimal
	}

	var results []Result
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("TO_CHAR(created_at, 'YYYY-MM') as month, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("month").
		Order("month").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalMap := make(map[string]decimal.Decimal)
	for _, r := range results {
		totalMap[r.Month] = r.Total
	}

	return totalMap, nil
}

func (s *ExpenseService) GetTotalExpensesByDay(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
	type Result struct {
		Day   string
		Total decimal.Decimal
	}

	var results []Result
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as day, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("day").
		Order("day").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalMap := make(map[string]decimal.Decimal)
	for _, r := range results {
		totalMap[r.Day] = r.Total
	}

	return totalMap, nil
}

func (s *ExpenseService) GetTotalExpensesByPaymentMethod(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
	type Result struct {
		PaymentMethod string
		Total         decimal.Decimal
	}

	var results []Result
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("payment_method, SUM(amount) as total").
		Where("user_id = ?", userID).
		Group("payment_method").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	totalMap := make(map[string]decimal.Decimal)
	for _, r := range results {
		totalMap[r.PaymentMethod] = r.Total
	}

	return totalMap, nil
}

func (s *ExpenseService) UpdateRecurringExpense(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID, req *UpdateRecurringExpenseRequest) (*RecurringExpense, error) {
	var recurringExpense RecurringExpense
	if err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", recurringExpenseID, userID).
		First(&recurringExpense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurring expense not found")
		}
		return nil, err
	}

	if req.Amount != nil {
		if req.Amount.LessThanOrEqual(decimal.Zero) {
			return nil, errors.New("amount must be greater than zero")
		}
		recurringExpense.Amount = *req.Amount
	}
	if req.Description != nil {
		recurringExpense.Description = *req.Description
	}
	if req.CategoryID != nil {
		recurringExpense.CategoryID = *req.CategoryID
	}
	if req.Frequency != nil {
		recurringExpense.Frequency = *req.Frequency
	}
	if req.PaymentMethod != nil {
		recurringExpense.PaymentMethod = *req.PaymentMethod
	}
	if req.StartDate != nil {
		recurringExpense.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		recurringExpense.EndDate = req.EndDate
	}

	if err := s.db.WithContext(ctx).Save(&recurringExpense).Error; err != nil {
		return nil, err
	}
	return &recurringExpense, nil
}

func (s *ExpenseService) GetRecurringExpenseByID(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) (*RecurringExpense, error) {
	var recurringExpense RecurringExpense
	if err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", recurringExpenseID, userID).
		First(&recurringExpense).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurring expense not found")
		}
		return nil, err
	}
	return &recurringExpense, nil
}

func (s *ExpenseService) GetRecurringExpenseHistory(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND recurring_expense_id = ?", userID, recurringExpenseID).
		Preload("Category").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}
