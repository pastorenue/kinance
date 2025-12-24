package expense

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/shopspring/decimal"
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

func (s *Service) CreateExpense(ctx context.Context, userID uuid.UUID, req *CreateExpenseRequest) (*Expense, error) {
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

func (s *Service) UpdateExpense(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID, req *UpdateExpenseRequest) (*Expense, error) {
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

func (s *Service) DeleteExpense(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).Delete(&Expense{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) GetExpenses(ctx context.Context, userID uuid.UUID, pagination *common.PaginationParams) ([]Expense, error) {
	var expenses []Expense

	offset := (pagination.Page - 1) * pagination.PageSize
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Preload("Category").Offset(offset).Limit(pagination.PageSize).Find(&expenses).Error; err != nil {
		return nil, err
	}

	return expenses, nil
}

func (s *Service) GetExpenseByID(ctx context.Context, userID uuid.UUID, expenseID uuid.UUID) (*Expense, error) {
	var expense Expense
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", expenseID, userID).Preload("Category").First(&expense).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("expense not found")
		}
		return nil, err
	}
	return &expense, nil
}

func (s *Service) GetExpensesByCategoryID(ctx context.Context, userID uuid.UUID, categoryID uuid.UUID) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.WithContext(ctx).Where("user_id = ? AND category_id = ?", userID, categoryID).Preload("Category").Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

func (s *Service) GetTotalExpensesByCategory(ctx context.Context, userID uuid.UUID) (map[uuid.UUID]decimal.Decimal, error) {
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

func (s *Service) GetTotalExpenses(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := s.db.WithContext(ctx).Model(&Expense{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ?", userID).
		Scan(&total).Error; err != nil {
		return decimal.Zero, err
	}
	return total, nil
}

func (s *Service) CreateRecurringExpense(ctx context.Context, userID uuid.UUID, req *CreateRecurringExpenseRequest) (*RecurringExpense, error) {
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
	recurringExpense.LastProcessed = time.Now()

	// Use a transaction to ensure atomicity
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(recurringExpense).Error; err != nil {
			return err
		}

		// Create an initial expense if the start date is today or in the past
		if !recurringExpense.StartDate.After(time.Now()) {
			expense := &Expense{
				Amount:             recurringExpense.Amount,
				Description:        recurringExpense.Description,
				CategoryID:         recurringExpense.CategoryID,
				UserID:             recurringExpense.UserID,
				PaymentMethod:      recurringExpense.PaymentMethod,
				RecurringExpenseID: &recurringExpense.ID,
			}
			expense.ID = uuid.New()

			if err := tx.Create(expense).Error; err != nil {
				s.logger.Error(
					"Failed to create initial expense for recurring expense",
					"error",
					err,
					"recurring_expense_id",
					recurringExpense.ID,
				)
				return err
			} else {
				s.logger.Info(
					"Created initial expense for recurring expense",
					"recurring_expense_id",
					recurringExpense.ID,
					"expense_id",
					expense.ID,
				)
			}

			// Calculate the next due date after creating the initial expense
			recurringExpense.CalculateNextDueDate()
			if err := tx.Save(recurringExpense).Error; err != nil {
				s.logger.Error(
					"Failed to update next due date for recurring expense",
					"error",
					err,
					"recurring_expense_id",
					recurringExpense.ID,
				)
				return err
			}
		}
		// Load the category and expenses relationship
		if err := tx.WithContext(ctx).Preload("Category").Preload("Expenses").First(recurringExpense, recurringExpense.ID).Error; err != nil {
			s.logger.Error(
				"Failed to preload category for recurring expense",
				"error",
				err,
				"recurring_expense_id",
				recurringExpense.ID,
			)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return recurringExpense, nil
}

func (s *Service) GetRecurringExpenses(ctx context.Context, userID uuid.UUID) ([]RecurringExpenseResponse, error) {
	var recurringExpenses []RecurringExpense
	if err := s.db.WithContext(ctx).Where("user_id = ?", userID).Find(&recurringExpenses).Error; err != nil {
		return nil, err
	}

	var responses []RecurringExpenseResponse
	for _, rec := range recurringExpenses {
		responses = append(responses, RecurringExpenseResponse{
			RecurringExpense: rec,
			DaysUntilDue:     rec.DaysUntilNextDue(time.Now()),
		})
	}
	return responses, nil
}

func (s *Service) DeleteRecurringExpense(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) error {
	if err := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", recurringExpenseID, userID).Delete(&RecurringExpense{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTotalExpensesByMonth(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
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

func (s *Service) GetTotalExpensesByDay(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
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

func (s *Service) GetTotalExpensesByPaymentMethod(ctx context.Context, userID uuid.UUID) (map[string]decimal.Decimal, error) {
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

func (s *Service) UpdateRecurringExpense(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID, req *UpdateRecurringExpenseRequest) (*RecurringExpense, error) {
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

func (s *Service) GetRecurringExpenseByID(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) (*RecurringExpenseResponse, error) {
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
	return &RecurringExpenseResponse{
		RecurringExpense: recurringExpense,
		DaysUntilDue:     recurringExpense.DaysUntilNextDue(time.Now()),
	}, nil
}

/* Return the history of expenses associated with a specific recurring expense for a given user.
 *
 * @param ctx - The context for the request
 * @param userID - The ID of the user
 * @param recurringExpenseID - The ID of the recurring expense
 * @returns A slice of Expense and an error if the database query fails.
 */
func (s *Service) GetRecurringExpenseHistory(ctx context.Context, userID uuid.UUID, recurringExpenseID uuid.UUID) ([]Expense, error) {
	var expenses []Expense
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND recurring_expense_id = ?", userID, recurringExpenseID).
		Preload("Category").
		Find(&expenses).Error; err != nil {
		return nil, err
	}
	return expenses, nil
}

/* Get upcoming recurring expenses due in the next offset days
*
* @param ctx - The context for the request
* @param userID - The ID of the user
* @param dueInterval - The number of days to check for upcoming expenses
* @returns A list of upcoming recurring expenses
 */
func (s *Service) GetUpcomingRecurringExpenses(ctx context.Context, userID uuid.UUID, dueInterval int) ([]RecurringExpenseResponse, error) {
	var recurringExpenses []RecurringExpense
	// Calculate the target date in Go and pass it as a parameter
	targetDate := time.Now().AddDate(0, 0, dueInterval)

	// Return only recurring expenses with next_due_date just targetDate away only
	if err := s.db.WithContext(ctx).
		Where("user_id = ? AND next_due_date = ?", userID, targetDate).
		Find(&recurringExpenses).Error; err != nil {
		return nil, err
	}

	// Preload the Category relationship for each recurring expense
	for i := range recurringExpenses {
		if err := s.db.WithContext(ctx).Model(&recurringExpenses[i]).Association("Category").Find(&recurringExpenses[i].Category); err != nil {
			return nil, err
		}
	}

	// Return the list of upcoming recurring expenses
	result := make([]RecurringExpenseResponse, len(recurringExpenses))
	for i, re := range recurringExpenses {
		result[i] = RecurringExpenseResponse{
			RecurringExpense: re,
			DaysUntilDue:     re.DaysUntilNextDue(time.Now()),
		}
	}
	return result, nil
}

/* Background Job to process all recurring expenses that are due as of the current time.
 *
 * For each due recurring expense, it creates a corresponding expense record, updates the next due date,
 * and logs the operation. If an error occurs while creating an expense or updating a recurring expense,
 * it logs the error and continues processing the remaining items.
 * Returns an error if the initial query for due recurring expenses fails.
 * @param ctx - The context for the request
 * @returns An error if processing fails
 */
func (s *Service) ProcessRecurringExpenses(ctx context.Context) error {
	var dueExpenses []RecurringExpense
	if err := s.db.WithContext(ctx).
		Where("next_due_date <= ?", time.Now()).
		Find(&dueExpenses).Error; err != nil {
		return err
	}

	for _, re := range dueExpenses {
		// Create an expense for each due recurring expense
		expense := &Expense{
			Amount:             re.Amount,
			Description:        re.Description,
			CategoryID:         re.CategoryID,
			UserID:             re.UserID,
			PaymentMethod:      re.PaymentMethod,
			RecurringExpenseID: &re.ID,
		}
		expense.ID = uuid.New()

		if err := s.db.WithContext(ctx).Create(expense).Error; err != nil {
			s.logger.Error(
				"Failed to create expense for recurring expense",
				"error",
				err,
				"recurring_expense_id",
				re.ID,
			)
			continue
		}
		re.CalculateNextDueDate()
		re.LastProcessed = time.Now()

		// Update the recurring expense with the new next due date
		if err := s.db.WithContext(ctx).Save(&re).Error; err != nil {
			s.logger.Error(
				"Failed to update recurring expense",
				"error",
				err,
				"recurring_expense_id",
				re.ID,
			)
		}
		s.logger.Info(
			"Processed recurring expense",
			"recurring_expense_id",
			re.ID,
			"next_due_date",
			re.NextDueDate,
		)
	}

	return nil
}
