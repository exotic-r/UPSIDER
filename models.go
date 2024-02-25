package main

import (
	"time"

	"gorm.io/gorm"
)

type StatusEnum string

const (
	Unprocessed StatusEnum = "未処理"
	InProgress  StatusEnum = "処理中"
	Paid        StatusEnum = "支払い済み"
	Error       StatusEnum = "エラー"
)

type Company struct {
	gorm.Model
	LegalName          string `json:"legalName"`
	RepresentativeName string `json:"representativeName"`
	PhoneNumber        string `json:"phoneNumber"`
	PostalCode         string `json:"postalCode"`
	Address            string `json:"address"`
}

type User struct {
	gorm.Model
	Name      string `json:"name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
	CompanyID int    `json:"companyID" gorm:"unique"`
}

type Client struct {
	gorm.Model
	LegalName          string `json:"legalName"`
	RepresentativeName string `json:"representativeName"`
	PhoneNumber        string `json:"phoneNumber"`
	PostalCode         string `json:"postalCode"`
	Address            string `json:"address"`
	CompanyID          int    `json:"companyID" gorm:"unique"`
}

type ClientBankAccount struct {
	gorm.Model
	ClientBankAccountID int    `gorm:"unique"`
	BankName            string `json:"bankName"`
	BranchName          string `json:"branchName"`
	AccountNumber       string `json:"accountNumber" gorm:"unique"`
	AccountName         string `json:"accountName"`
	ClientID            int    `json:"clientID" gorm:"unique"`
}

type Invoice struct {
	gorm.Model
	IssueDate     time.Time  `json:"issueDate"`
	PaymentAmount float64    `json:"paymentAmount"`
	Fee           float64    `json:"fee"`
	FeeRate       float64    `json:"feeRate"`
	Tax           float64    `json:"tax"`
	TaxRate       float64    `json:"taxRate"`
	TotalAmount   float64    `json:"totalAmount"`
	DueDate       time.Time  `json:"dueDate"`
	Status        StatusEnum `json:"status"`
	CompanyID     int
	ClientID      int
}
