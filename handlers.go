package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	if err := db.Create(&invoice).Error; err != nil {
		log.Printf("Error creating invoice: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
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
	if err := query.Find(&filteredInvoices).Error; err != nil {
		log.Printf("Error retrieving invoices: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, filteredInvoices)
}

func parseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{}
	}
	return date
}
