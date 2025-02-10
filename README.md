# Personal Finance Tracker

## Description
A simple personal finance tracker API built with Golang and PostgreSQL. This project allows users to track their income, expenses, and budgets in a structured manner. It provides authentication, basic financial reports, and a RESTful API for managing transactions.

## Tech Stack
- **Language**: Golang
- **Framework**: Gin (for API handling)
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (golang-jwt)
- **Deployment**: Docker

## Features
- User authentication (JWT-based)
- CRUD operations for transactions (income & expenses)
- Monthly budget tracking
- Secure API with middleware
- Docker support for easy deployment

## Project Structure
```
personal-finance-tracker/
│── cmd/                     # Entry point for the app
│   ├── main.go               # Initializes and runs the app
│
│── config/                   # Configuration files
│   ├── config.go             # Loads environment variables
│
│── models/                   # Database models
│   ├── user.go               # User model
│   ├── transaction.go        # Transaction model
│
│── controllers/              # Handles API requests
│   ├── user_controller.go    # User-related endpoints
│   ├── transaction_controller.go # Transaction-related endpoints
│
│── routes/                   # API routes
│   ├── routes.go             # Initializes all routes
│
│── services/                 # Business logic
│   ├── user_service.go       # Handles user-related logic
│   ├── transaction_service.go # Handles transaction logic
│
│── repository/               # Database interaction
│   ├── user_repository.go    # User-related database queries
│   ├── transaction_repository.go # Transaction-related queries
│
│── middleware/               # Authentication and validation
│   ├── auth_middleware.go    # JWT authentication middleware
│
│── db/                       # Database connection
│   ├── database.go           # Initializes database connection
│
│── utils/                    # Helper functions
│   ├── jwt.go                # JWT token functions
│
│── .env                      # Environment variables (DB credentials, JWT secret)
│── docker-compose.yml        # Docker setup for PostgreSQL
│── go.mod                    # Dependencies
│── go.sum                    # Checksums
│── LICENSE                   # Project license
│── README.md                 # Project documentation
```

## Installation
### Prerequisites
- Golang installed
- Docker & Docker Compose
- PostgreSQL database (or use Docker setup)

### Setup
```sh
git clone https://github.com/TsonasIoannis/go-personal-finance-tracker.git
cd go-personal-finance-tracker
cp .env.example .env  # Configure environment variables
```

### Running Locally
```sh
docker-compose up -d  # Starts PostgreSQL

# Run the application
go run cmd/main.go
```

## API Endpoints
| Method | Endpoint                 | Description                  |
|--------|--------------------------|------------------------------|
| POST   | `/register`               | Register a new user          |
| POST   | `/login`                  | Authenticate user & get JWT  |
| GET    | `/transactions`           | Get all transactions         |
| POST   | `/transactions`           | Add a new transaction        |
| DELETE | `/transactions/:id`       | Delete a transaction         |

## License
This project is licensed under the MIT License.

