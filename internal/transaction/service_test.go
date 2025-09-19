package transaction

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/expense"
	"github.com/pastorenue/kinance/internal/income"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// mockLogger implements common.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string, keysAndValues ...interface{})  {}
func (m *mockLogger) Error(msg string, keysAndValues ...interface{}) {}
func (m *mockLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Helper function to create Currency pointer
func ptrCurrency(c common.Currency) *common.Currency {
	return &c
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Configure SQLite for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",
			SingularTable: false,
		},
	})
	assert.NoError(t, err)

	// Create test tables
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			name TEXT NOT NULL,
			user_id TEXT NOT NULL,
			parent_category_id TEXT,
			color_code CHAR(7),
			UNIQUE(name, user_id)
		);
		CREATE INDEX IF NOT EXISTS idx_categories_user_id ON categories(user_id);
		CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_category_id);

		CREATE TABLE IF NOT EXISTS transactions (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			user_id TEXT NOT NULL,
			amount DECIMAL NOT NULL,
			description TEXT,
			category_id TEXT,
			transaction_date DATETIME NOT NULL,
			status TEXT DEFAULT 'completed',
			type TEXT NOT NULL,
			processing_object_id TEXT UNIQUE,
			currency TEXT NOT NULL DEFAULT 'EUR',
			exclude_from_analytics BOOLEAN DEFAULT 0,
			merchant_id TEXT,
			receipt_id TEXT,
			metadata JSON DEFAULT '{}',
			payment_method TEXT,
			FOREIGN KEY(category_id) REFERENCES categories(id),
			FOREIGN KEY(merchant_id) REFERENCES merchants(id)
		);
		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_category_id ON transactions(category_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_merchant_id ON transactions(merchant_id);

		CREATE TABLE IF NOT EXISTS expenses (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			amount DECIMAL NOT NULL,
			description TEXT,
			category_id TEXT,
			user_id TEXT NOT NULL,
			currency TEXT NOT NULL DEFAULT 'EUR',
			receipt_url TEXT,
			payment_method TEXT,
			transaction_id TEXT,
			recurring_expense_id TEXT,
			FOREIGN KEY(category_id) REFERENCES categories(id),
			FOREIGN KEY(recurring_expense_id) REFERENCES recurring_expenses(id)
		);
		CREATE INDEX IF NOT EXISTS idx_expenses_user_id ON expenses(user_id);
		CREATE INDEX IF NOT EXISTS idx_expenses_category_id ON expenses(category_id);
		CREATE INDEX IF NOT EXISTS idx_expenses_transaction_id ON expenses(transaction_id);

		CREATE TABLE IF NOT EXISTS incomes (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			amount DECIMAL NOT NULL,
			currency TEXT NOT NULL DEFAULT 'EUR',
			source_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			status TEXT DEFAULT 'pending',
			note TEXT,
			metadata TEXT DEFAULT '{}',
			category_id TEXT,
			FOREIGN KEY(category_id) REFERENCES categories(id),
			FOREIGN KEY(source_id) REFERENCES sources(id)
		);
		CREATE INDEX IF NOT EXISTS idx_incomes_user_id ON incomes(user_id);
		CREATE INDEX IF NOT EXISTS idx_incomes_category_id ON incomes(category_id);
		CREATE INDEX IF NOT EXISTS idx_incomes_source_id ON incomes(source_id);

		CREATE TABLE IF NOT EXISTS sources (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			logo_url TEXT,
			swift_code TEXT,
			is_validated BOOLEAN DEFAULT 0
		);
		CREATE INDEX IF NOT EXISTS idx_sources_swift_code ON sources(swift_code);

		CREATE TABLE IF NOT EXISTS merchants (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			name TEXT NOT NULL,
			website TEXT,
			logo_url TEXT,
			user_id TEXT NOT NULL,
			UNIQUE(name)
		);
		CREATE INDEX IF NOT EXISTS idx_merchants_user_id ON merchants(user_id);

		CREATE TABLE IF NOT EXISTS tags (
			id TEXT PRIMARY KEY,
			created_at DATETIME,
			updated_at DATETIME,
			name TEXT NOT NULL,
			color TEXT DEFAULT '#007bff',
			user_id TEXT NOT NULL,
			UNIQUE(name)
		);
		CREATE INDEX IF NOT EXISTS idx_tags_user_id ON tags(user_id);

		CREATE TABLE IF NOT EXISTS transaction_tags (
			transaction_id TEXT NOT NULL,
			tag_id TEXT NOT NULL,
			PRIMARY KEY(transaction_id, tag_id),
			FOREIGN KEY(transaction_id) REFERENCES transactions(id),
			FOREIGN KEY(tag_id) REFERENCES tags(id)
		);
	`).Error
	assert.NoError(t, err)

	// Create a test category to satisfy foreign key constraint
	testCategory := &category.Category{
		BaseModel: common.BaseModel{ID: uuid.New()},
		Name:      "Test Category",
		UserID:    uuid.New(),
		ColorCode: "#000000",
	}
	err = db.Create(testCategory).Error
	assert.NoError(t, err)

	return db
}

func TestCreateExpenseTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()

	// Get test category ID
	var testCategory category.Category
	err := db.First(&testCategory).Error
	assert.NoError(t, err)

	// Test data
	testTime := time.Now()
	req := &CreateTransactionRequest{
		Type:            TypeExpense,
		Amount:          decimal.NewFromInt(100),
		Description:     "Test expense",
		CategoryID:      testCategory.ID,
		Currency:        common.EUR,
		TransactionDate: testTime,
		PaymentMethod:   common.Card,
		Metadata: map[string]interface{}{
			"location": "San Francisco",
			"store_id": 123,
		},
	}

	// Execute
	response, err := svc.CreateExpenseTransaction(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Amount, response.Transaction.Amount)
	assert.Equal(t, req.Description, response.Transaction.Description)
	assert.Equal(t, req.CategoryID, response.Transaction.CategoryID)
	assert.Equal(t, userID, response.Transaction.UserID)

	// Verify the expense was created
	expense, ok := response.Entity.(*expense.Expense)
	assert.True(t, ok)
	assert.Equal(t, req.Amount, expense.Amount)
	assert.Equal(t, req.Description, expense.Description)
	assert.Equal(t, req.CategoryID, expense.CategoryID)
	assert.Equal(t, req.Metadata["location"], "San Francisco")
	assert.Equal(t, req.Metadata["store_id"], 123)
	assert.Equal(t, userID, expense.UserID)
}

func TestCreateExpenseTransaction_InvalidCategory(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()

	// Test data with nil category
	req := &CreateTransactionRequest{
		Type:            "expense",
		Amount:          decimal.NewFromInt(100),
		Description:     "Test expense",
		CategoryID:      uuid.Nil, // Invalid category
		Currency:        "EUR",
		TransactionDate: time.Now(),
		PaymentMethod:   "card",
	}

	// Execute
	response, err := svc.CreateExpenseTransaction(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "category ID cannot be empty")
}

func TestCreateExpenseTransaction_InvalidType(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	// Test data with wrong type
	req := &CreateTransactionRequest{
		Type:            "income", // Wrong type
		Amount:          decimal.NewFromInt(100),
		Description:     "Test expense",
		CategoryID:      categoryID,
		Currency:        "EUR",
		TransactionDate: time.Now(),
		PaymentMethod:   "card",
	}

	// Execute
	response, err := svc.CreateExpenseTransaction(ctx, userID, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "transaction type must be 'expense'")
}

func TestGetTransactions_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	// Create test transactions
	transactions := []Transaction{
		{
			BaseModel:       common.BaseModel{ID: uuid.New()},
			UserID:          userID,
			Type:            TypeExpense,
			Amount:          decimal.NewFromInt(100),
			TransactionDate: time.Now(),
			CategoryID:      categoryID,
			Currency:        ptrCurrency(common.EUR),
		},
		{
			BaseModel:       common.BaseModel{ID: uuid.New()},
			UserID:          userID,
			Type:            TypeIncome,
			Amount:          decimal.NewFromInt(200),
			TransactionDate: time.Now(),
			CategoryID:      categoryID,
			Currency:        ptrCurrency(common.EUR),
		},
	}

	for _, tx := range transactions {
		assert.NoError(t, db.Create(&tx).Error)
	}

	// Execute
	result, err := svc.GetTransactions(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, len(transactions))
	assert.Equal(t, transactions[0].Amount, result[0].Amount)
	assert.Equal(t, transactions[1].Amount, result[1].Amount)
}

func TestGetTransactionByID_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()
	transactionID := uuid.New()

	// Create test transaction
	transaction := Transaction{
		BaseModel:       common.BaseModel{ID: transactionID},
		UserID:          userID,
		Type:            TypeExpense,
		Amount:          decimal.NewFromInt(100),
		TransactionDate: time.Now(),
		CategoryID:      categoryID,
		Currency:        ptrCurrency(common.EUR),
	}
	assert.NoError(t, db.Create(&transaction).Error)

	// Execute
	result, err := svc.GetTransactionByID(ctx, userID, transactionID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, transaction.ID, result.ID)
	assert.Equal(t, transaction.Amount, result.Amount)
}

func TestGetTransactionByID_NotFound(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()

	// Execute
	result, err := svc.GetTransactionByID(ctx, userID, transactionID)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, err == gorm.ErrRecordNotFound)
}

func TestCreateIncomeTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	// Create test source
	testSource := &income.Source{
		Name:      "Test Bank",
		SwiftCode: "BARCGB22",
	}
	testSource.ID = uuid.New()
	err := db.Create(testSource).Error
	assert.NoError(t, err)

	// Create test category
	testCategory := &category.Category{
		Name:      "Test Category",
		UserID:    userID,
		ColorCode: "#000000",
	}
	testCategory.ID = categoryID
	err = db.Create(testCategory).Error
	assert.NoError(t, err)

	// Test data
	req := &CreateIncomeTransactionRequest{
		CreateTransactionRequest: CreateTransactionRequest{
			Type:            "income",
			Amount:          decimal.NewFromInt(1000),
			Description:     "Test income",
			CategoryID:      categoryID,
			Currency:        common.EUR,
			TransactionDate: time.Now(),
			PaymentMethod:   common.BankTransfer,
		},
		SwiftCode: "BARCGB22",
	}

	// Execute
	response, err := svc.CreateIncomeTransaction(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Amount, response.Transaction.Amount)
	assert.Equal(t, req.Description, response.Transaction.Description)
	assert.Equal(t, req.CategoryID, response.Transaction.CategoryID)
	assert.Equal(t, userID, response.Transaction.UserID)

	// Verify income instance was created
	income, ok := response.Entity.(income.Income)
	assert.True(t, ok)
	assert.Equal(t, req.Amount, income.Amount)
	assert.Equal(t, req.CategoryID, income.CategoryID)
	assert.Equal(t, userID, income.UserID)
}

func TestCreateTransferTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	categoryID := uuid.New()

	// Test data
	req := &CreateTransactionRequest{
		Type:            "transfer",
		Amount:          decimal.NewFromInt(500),
		Description:     "Test transfer",
		CategoryID:      categoryID,
		Currency:        common.EUR,
		TransactionDate: time.Now(),
		PaymentMethod:   common.BankTransfer,
	}

	// Execute
	transaction, err := svc.CreateTransferTransaction(ctx, userID, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, req.Amount, transaction.Amount)
	assert.Equal(t, req.Description, transaction.Description)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, req.Type, transaction.Type)
}

func TestLinkToTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	ctx := context.Background()
	userID := uuid.New()
	transactionID := uuid.New()
	processingObjectID := uuid.New()

	// Create test transaction
	transaction := Transaction{
		BaseModel:       common.BaseModel{ID: transactionID},
		UserID:          userID,
		Type:            TypeExpense,
		Amount:          decimal.NewFromInt(100),
		TransactionDate: time.Now(),
		Currency:        ptrCurrency(common.EUR),
	}
	assert.NoError(t, db.Create(&transaction).Error)

	// Execute
	err := svc.linkToTransaction(ctx, processingObjectID, transactionID)

	// Assert
	assert.NoError(t, err)

	// Verify transaction was updated
	var updatedTransaction Transaction
	err = db.First(&updatedTransaction, transactionID).Error
	assert.NoError(t, err)
	assert.Equal(t, processingObjectID, *updatedTransaction.ProcessingObjectID)
}
