package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func calculateInvoiceAmount(paymentAmount float64) float64 {
	fee := paymentAmount * 0.04
	tax := (paymentAmount + fee) * 0.1
	return paymentAmount + fee + tax
}

func createInvoiceHandler(c *gin.Context) {
	var requestBody struct {
		PaymentAmount     float64           `json:"paymentAmount"`
		DueDate           string            `json:"dueDate"`
		Company           Company           `json:"company"`
		User              User              `json:"user"`
		Client            Client            `json:"client"`
		ClientBankAccount ClientBankAccount `json:"clientBankAccount"`
		Status            StatusEnum        `json:"status"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	db.Create(&requestBody.Company)

	requestBody.User.CompanyID = int(requestBody.Company.ID)
	db.Create(&requestBody.User)

	requestBody.Client.CompanyID = int(requestBody.Company.ID)
	db.Create(&requestBody.Client)

	requestBody.ClientBankAccount.ClientID = int(requestBody.Client.ID)
	db.Create(&requestBody.ClientBankAccount)

	invoice := Invoice{
		IssueDate:     time.Now(),
		PaymentAmount: requestBody.PaymentAmount,
		Fee:           requestBody.PaymentAmount * 0.04,
		FeeRate:       0.04,
		Tax:           requestBody.PaymentAmount * 0.04 * 0.1,
		TaxRate:       0.1,
		TotalAmount:   calculateInvoiceAmount(requestBody.PaymentAmount),
		DueDate:       parseDate(requestBody.DueDate),
		Status:        requestBody.Status,
		CompanyID:     int(requestBody.Company.ID),
		ClientID:      int(requestBody.Client.ID),
	}

	db.Create(&invoice)
	c.JSON(http.StatusOK, invoice)
}

func getInvoicesHandler(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	query := db.Model(&Invoice{}).Where("due_date > ?", time.Now())

	startDateStr := c.Query("startDate")
	endDateStr := c.Query("endDate")

	if startDateStr != "" {
		query = query.Where("due_date > ?", parseDate(startDateStr))
	}
	if endDateStr != "" {
		query = query.Where("due_date < ?", parseDate(endDateStr))
	}

	var filteredInvoices []Invoice
	query.Find(&filteredInvoices)
	c.JSON(http.StatusOK, filteredInvoices)
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return date
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Invoice{})

	r := gin.Default()
	r.Use(authenticate) // TODO: implmenet authentication logic

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/api/v1/invoices", createInvoiceHandler)
	r.GET("/api/v1/invoices", getInvoicesHandler)

	fmt.Println("Server is running on :8080...")
	r.Run("localhost:8080")
}
