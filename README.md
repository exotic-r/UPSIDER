# UPSIDER

Super Payment API is a Golang-based REST API for managing invoices and payments in a fictional web service called "Super Payment". It utilizes the Gin framework for routing and Gorm as the ORM for database operations.

## Features

- Create new invoices with calculated amounts.
- Retrieve a list of invoices within a specified date range.
- Authentication (TO BE DONE)

## Getting Started

### Prerequisites

- Golang installed on your machine.
- [Gin](https://github.com/gin-gonic/gin) web framework.
- [Gorm](https://gorm.io/) ORM library.
- SQLite (for in-memory database).

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/super-payment-api.git
   cd super-payment-api
   ```

2. Install dependencies:
   
   ```bash 
   go mod download
    ```

Get Invoices
Endpoint: GET /api/invoices

Query Parameters:

startDate (string, required): Start date for filtering invoices.
endDate (string, required): End date for filtering invoices.
Example:

GET /api/invoices?startDate=2024-02-01&endDate=2024-02-28

Response:


### Run application
    go run main.go models.go auth.go handlers.go

## API Endpoints

### Create Invoice
Endpoint: POST /api/invoices

Example Curl Request:
```
curl --location 'http://localhost:8080/api/v1/invoices' \
--header 'Authorization: Bearer your_token_here' \
--header 'Content-Type: application/json' \
--data-raw '{
    "paymentAmount": 15000,
    "dueDate": "2024-03-15",
    "company": {
      "legalName": "Example Company",
      "representativeName": "John Doe",
      "phoneNumber": "123-456-7890",
      "postalCode": "12345",
      "address": "123 Main St"
    },
    "user": {
      "name": "Alice",
      "email": "alice@example.com",
      "password": "securepassword"
    },
    "client": {
      "legalName": "Client Corp",
      "representativeName": "Jane Smith",
      "phoneNumber": "987-654-3210",
      "postalCode": "54321",
      "address": "456 Park Ave"
    },
    "clientBankAccount": {
      "bankName": "Bank XYZ",
      "branchName": "Main Branch",
      "accountNumber": "987654321",
      "accountName": "Client Corp Account"
    }
}'
```

Response:
```
{
    "ID": 1,
    "CreatedAt": "2024-02-25T22:06:12.074367+08:00",
    "UpdatedAt": "2024-02-25T22:06:12.074367+08:00",
    "DeletedAt": null,
    "issueDate": "2024-02-25T22:06:12.074281+08:00",
    "paymentAmount": 15000,
    "fee": 600,
    "feeRate": 0.04,
    "tax": 60,
    "taxRate": 0.1,
    "totalAmount": 17160,
    "dueDate": "2024-03-15T00:00:00Z",
    "status": "",
    "CompanyID": 0,
    "ClientID": 0
}
```


### Get Invoices
Endpoint: GET /api/invoices

Query Parameters:
- startDate (string, required): Start date for filtering invoices.
- endDate (string, required): End date for filtering invoices.

Example Curl Request:
    ```curl --location 'localhost:8080/api/v1/invoices?startDate=2021-01-04&endDateStr=2023-01-04' --header 'Authorization: test'```

Response:
```
[
  {
    "ID": 1,
    "CreatedAt": "2024-02-24T12:00:00Z",
    "UpdatedAt": "2024-02-24T12:00:00Z",
    "DeletedAt": null,
    "issueDate": "2024-02-24T12:00:00Z",
    "paymentAmount": 1000,
    "fee": 40,
    "feeRate": 0.04,
    "tax": 44,
    "taxRate": 0.1,
    "totalAmount": 1084,
    "dueDate": "2024-03-01T00:00:00Z",
    "status": "未処理",
    "companyID": 1,
    "clientID": 1
  }
]
```

## Testing
    go test


## TODO
### Finish Authentication Implementation:
- Complete the implementation of authentication, ensuring that user identity is properly verified. Consider using a secure authentication method (e.g., JWT, OAuth).
Authorization Middleware: Implement middleware for authorization to control access to different API endpoints based on user roles and permissions.

### Database Management:
- Replace In-Memory SQLite with Persistent Database: Transition from the in-memory SQLite database to a persistent database (e.g., PostgreSQL, MySQL). This ensures data persistence across application restarts.

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvement, please open an issue or create a pull request.