package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"stockapp/internal/db"
	"stockapp/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite" // Pure go sqlite
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Use in-memory SQLite for testing
	d, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Override the ID field for tests because SQLite doesn't support 'gen_random_uuid()'
	// We can use AutoMigrate but we need to be careful with the UUID default
	d.AutoMigrate(&models.Stock{}, &models.StockPrice{})
	return d
}

func TestGetStocks(t *testing.T) {
	// Setup
	testDB := setupTestDB()
	db.DB = testDB // Overwrite global DB

	// Seed data
	testDB.Create(&models.Stock{
		Symbol:       "TEST",
		Name:         "Test Stock",
		CurrentPrice: 100.0,
		UpdatedAt:    time.Now(),
	})
	testDB.Create(&models.Stock{
		Symbol:       "ABC",
		Name:         "ABC Corp",
		CurrentPrice: 50.0,
		UpdatedAt:    time.Now(),
	})

	// Router
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/stocks", GetStocks)

	// Test 1: List all
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stocks", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var stocks []models.Stock
	err := json.Unmarshal(w.Body.Bytes(), &stocks)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(stocks))

	// Test 2: Search
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/stocks?search=ABC", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &stocks)
	if assert.Equal(t, 1, len(stocks)) {
		assert.Equal(t, "ABC Corp", stocks[0].Name)
	}
}

func TestGetStockDetail(t *testing.T) {
	testDB := setupTestDB()
	db.DB = testDB

	stock := models.Stock{
		ID:           "uuid-1",
		Symbol:       "XYZ",
		Name:         "XYZ Inc",
		CurrentPrice: 200.0,
	}
	testDB.Create(&stock)

	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/stocks/:symbol", GetStockDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stocks/XYZ", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	stockData := response["stock"].(map[string]interface{})
	assert.Equal(t, "XYZ", stockData["symbol"])
}
