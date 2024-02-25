package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
