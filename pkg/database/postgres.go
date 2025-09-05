package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	// "github.com/pastorenue/kinance/internal/budget"
	// "github.com/pastorenue/kinance/internal/receipt"
	// "github.com/pastorenue/kinance/internal/transaction"
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

	// Auto-migrate models
	err = db.AutoMigrate(
		&user.User{},
		// &user.Family{},
		// &budget.Budget{},
		// &transaction.Transaction{},
		// &transaction.Tag{},
		// &receipt.Receipt{},
		// &receipt.ReceiptItem{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
