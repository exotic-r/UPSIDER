package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateInvoiceHandler(t *testing.T) {
	payload := map[string]interface{}{
		"paymentAmount": 1000.00,
		"dueDate":       "2024-03-01",
		"status":        "未処理",
		"company": map[string]interface{}{
			"legalName":          "Example Company",
			"representativeName": "John Doe",
			"phoneNumber":        "123-456-7890",
			"postalCode":         "12345",
			"address":            "123 Main St",
		},
		"user": map[string]interface{}{
			"userID":   "user123",
			"name":     "Alice",
			"email":    "alice@example.com",
			"password": "securepassword",
		},
		"client": map[string]interface{}{
			"legalName":          "Client Corp",
			"representativeName": "Jane Smith",
			"phoneNumber":        "987-654-3210",
			"postalCode":         "54321",
			"address":            "456 Second St",
		},
		"clientBankAccount": map[string]interface{}{
			"bankName":      "Example Bank",
			"branchName":    "Main Branch",
			"accountNumber": "12345678",
			"accountName":   "Client Account",
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", "/api/invoices", bytes.NewBuffer(payloadBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := setupTestDB()

	gin.SetMode(gin.TestMode)
	context, _ := gin.CreateTestContext(rr)
	context.Request = req
	context.Set("db", db)

	createInvoiceHandler(context)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Contains(t, response, "ID")
	assert.Contains(t, response, "issueDate")
	assert.Contains(t, response, "paymentAmount")
	assert.Contains(t, response, "dueDate")
	assert.Contains(t, response, "status")

	db.Migrator().DropTable(&Invoice{}, &ClientBankAccount{}, &Client{}, &User{}, &Company{})
}

func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the test database")
	}

	db.AutoMigrate(&Company{}, &User{}, &Client{}, &ClientBankAccount{}, &Invoice{})

	return db
}
