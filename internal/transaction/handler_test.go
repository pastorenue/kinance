package transaction

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pastorenue/kinance/internal/category"
	"github.com/pastorenue/kinance/internal/common"
	"github.com/pastorenue/kinance/internal/income"
	"github.com/pastorenue/kinance/pkg/middleware"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRouter(svc *Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHandler(svc)

	r.Use(func(c *gin.Context) {
		var userID uuid.UUID
		if uid := c.Request.Header.Get("X-User-ID"); uid != "" {
			parsed, err := uuid.Parse(uid)
			if err == nil {
				userID = parsed
			}
		}
		if userID == uuid.Nil {
			userID = uuid.New()
		}
		c.Set(middleware.UserIDKey, userID)
		c.Next()
	})

	// Set up routes
	r.GET("/transactions", h.ListTransactions)
	r.POST("/transactions/expense", h.CreateExpenseTransaction)
	r.POST("/transactions/income", h.CreateIncomeTransaction)
	r.GET("/transactions/:id", h.GetTransaction)
	return r
}

func TestHandlerListTransactions_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Create a test transaction
	userID := uuid.New()
	testCategory := &category.Category{
		Name:      "Test Category",
		UserID:    userID,
		ColorCode: "#000000",
	}
	testCategory.ID = uuid.New()
	err := db.Create(testCategory).Error
	require.NoError(t, err)

	testTransaction := &Transaction{
		UserID:          userID,
		Amount:          decimal.NewFromInt(100),
		Description:     "Test transaction",
		CategoryID:      testCategory.ID,
		TransactionDate: time.Now(),
		Type:            TypeExpense,
		Currency:        ptrCurrency(common.EUR),
		PaymentMethod:   common.Card,
	}
	testTransaction.ID = uuid.New()
	err = db.Create(testTransaction).Error
	require.NoError(t, err)

	// Make the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions", nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response []Transaction
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, testTransaction.ID, response[0].ID)
}

func TestHandlerCreateExpenseTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Create a test category
	userID := uuid.New()
	testCategory := &category.Category{
		Name:      "Test Category",
		UserID:    userID,
		ColorCode: "#000000",
	}
	testCategory.ID = uuid.New()
	err := db.Create(testCategory).Error
	require.NoError(t, err)

	// Create request payload
	req := CreateTransactionRequest{
		Type:            TypeExpense,
		Amount:          decimal.NewFromInt(100),
		Description:     "Test expense",
		CategoryID:      testCategory.ID,
		Currency:        common.EUR,
		TransactionDate: time.Now(),
		PaymentMethod:   common.Card,
	}
	payload, err := json.Marshal(req)
	require.NoError(t, err)

	// Make the request
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions/expense", bytes.NewBuffer(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, req.Amount, response.Transaction.Amount)
	assert.Equal(t, req.Description, response.Transaction.Description)
}

func TestHandlerCreateIncomeTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Create test source
	testSource := &income.Source{
		Name:      "Test Bank",
		SwiftCode: "BARCGB22",
	}
	testSource.ID = uuid.New()
	err := db.Create(testSource).Error
	require.NoError(t, err)

	// Create test category
	userID := uuid.New()
	testCategory := &category.Category{
		Name:      "Test Category",
		UserID:    userID,
		ColorCode: "#000000",
	}
	testCategory.ID = uuid.New()
	err = db.Create(testCategory).Error
	require.NoError(t, err)

	// Create request payload
	req := CreateIncomeTransactionRequest{
		CreateTransactionRequest: CreateTransactionRequest{
			Type:            TypeIncome,
			Amount:          decimal.NewFromInt(1000),
			Description:     "Test income",
			CategoryID:      testCategory.ID,
			Currency:        common.EUR,
			TransactionDate: time.Now(),
			PaymentMethod:   common.BankTransfer,
		},
		SwiftCode: "BARCGB22",
	}
	payload, err := json.Marshal(req)
	require.NoError(t, err)

	// Make the request
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions/income", bytes.NewBuffer(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response TransactionResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, req.Amount, response.Transaction.Amount)
	assert.Equal(t, req.Description, response.Transaction.Description)
}

func TestHandlerGetTransaction_Success(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Create a test transaction
	userID := uuid.New()
	testCategory := &category.Category{
		Name:      "Test Category",
		UserID:    userID,
		ColorCode: "#000000",
	}
	testCategory.ID = uuid.New()
	err := db.Create(testCategory).Error
	require.NoError(t, err)

	testTransaction := &Transaction{
		UserID:          userID,
		Amount:          decimal.NewFromInt(100),
		Description:     "Test transaction",
		CategoryID:      testCategory.ID,
		TransactionDate: time.Now(),
		Type:            TypeExpense,
		Currency:        ptrCurrency(common.EUR),
		PaymentMethod:   common.Card,
	}
	testTransaction.ID = uuid.New()
	err = db.Create(testTransaction).Error
	require.NoError(t, err)

	// Make the request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactions/%s", testTransaction.ID), nil)
	req.Header.Set("X-User-ID", userID.String())
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response Transaction
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, testTransaction.ID, response.ID)
	assert.Equal(t, testTransaction.Amount, response.Amount)
}

func TestHandlerGetTransaction_NotFound(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Make request with non-existent ID
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/transactions/%s", uuid.New()), nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "not found")
}

func TestHandlerCreateExpenseTransaction_ValidationError(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	mockLog := new(mockLogger)
	svc := NewService(db, mockLog)
	router := setupTestRouter(svc)

	// Create an invalid request payload (missing required fields)
	req := CreateTransactionRequest{
		Type:        TypeExpense,
		Description: "Test expense",
	}
	payload, err := json.Marshal(req)
	require.NoError(t, err)

	// Make the request
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/transactions/expense", bytes.NewBuffer(payload))
	httpReq.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Contains(t, response["error"], "validation")
}
