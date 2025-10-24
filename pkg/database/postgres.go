package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

"github.com/pastorenue/kinance/internal/budget"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/document"
	"github.com/pastorenue/kinance/internal/expense"
	"github.com/pastorenue/kinance/internal/income"
	// "github.com/pastorenue/kinance/internal/receipt"
	"github.com/pastorenue/kinance/internal/transaction"
	"github.com/pastorenue/kinance/internal/user"
	"github.com/pastorenue/kinance/pkg/config"
)

func NewPostgres(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Create custom enum types first
	err = createEnumTypes(db)
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&user.User{},
		// &user.Family{},
		&category.Category{},
		&expense.RecurringExpense{},
		&expense.Expense{},
		&budget.Budget{},
		&transaction.Transaction{},
		&income.Income{},
		&document.Statement{},
		&document.Report{},
		// &transaction.Tag{},
		// &receipt.Receipt{},
		// &receipt.ReceiptItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createEnumTypes(db *gorm.DB) error {
	// Create payment_method enum type
	paymentMethodSQL := `
		DO $$ BEGIN
			CREATE TYPE payment_method AS ENUM ('cash', 'card', 'bank_transfer');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`

	if err := db.Exec(paymentMethodSQL).Error; err != nil {
		return fmt.Errorf("failed to create payment_method enum: %w", err)
	}

	// Create recurring_frequency enum type
	recurringFrequencySQL := `
		DO $$ BEGIN
			CREATE TYPE recurring_frequency AS ENUM ('daily', 'weekly', 'monthly', 'yearly');
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`

	if err := db.Exec(recurringFrequencySQL).Error; err != nil {
		return fmt.Errorf("failed to create recurring_frequency enum: %w", err)
	}

	return nil
}

func DropTransactionTable(db *gorm.DB) error {
	if err := db.Migrator().DropTable(&transaction.Transaction{}); err != nil {
		return err
	}
	return nil
}

func DropAllTables(db *gorm.DB) error {
	if err := db.Migrator().DropTable(
		&user.User{},
		// &user.Family{},
		&category.Category{},
		&expense.RecurringExpense{},
		&expense.Expense{},
		&budget.Budget{},
		&transaction.Transaction{},
		// &transaction.Tag{},
		// &receipt.Receipt{},
		// &receipt.ReceiptItem{},
		&income.Income{},
	); err != nil {
		return err
	}
	return nil
}