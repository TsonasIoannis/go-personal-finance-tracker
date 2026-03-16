# Personal Finance Tracker API

This repository is a portfolio-style Go backend project that illustrates how to build a small API layer with:

- Gin for HTTP routing
- GORM for persistence
- PostgreSQL as the primary database
- layered controller/service/repository structure
- token-based authentication
- Docker for local development
- unit tests across handlers, services, database, and repositories

## What The Project Demonstrates

The current implementation focuses on a realistic backend slice rather than a huge feature set.

It demonstrates:

- user registration and login
- signed auth token generation
- protected routes with middleware
- budget and transaction APIs tied to the authenticated user
- service-layer validation
- validated environment-based configuration
- explicit HTTP server timeouts and graceful shutdown
- readiness and health endpoints
- Docker-based and native local startup
- a passing `go test ./...` suite

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- Docker Compose

## Current Endpoints

### Public

| Method | Endpoint    | Description                            |
| ------ | ----------- | -------------------------------------- |
| POST   | `/register` | Register a user and return a token     |
| POST   | `/login`    | Authenticate a user and return a token |
| GET    | `/health`   | Liveness probe                         |
| GET    | `/ready`    | Readiness probe backed by the database |

### Protected

These endpoints require `Authorization: Bearer <token>`.

| Method | Endpoint            | Description                                         |
| ------ | ------------------- | --------------------------------------------------- |
| GET    | `/transactions`     | List the authenticated user's transactions          |
| POST   | `/transactions`     | Create a transaction for the authenticated user     |
| DELETE | `/transactions/:id` | Delete one of the authenticated user's transactions |
| GET    | `/budgets`          | List the authenticated user's budgets               |
| POST   | `/budgets`          | Create a budget for the authenticated user          |
| DELETE | `/budgets/:id`      | Delete one of the authenticated user's budgets      |

## Project Structure

```text
cmd/
  main.go                  application entrypoint

internal/
  auth/                    token generation and parsing
  controllers/             HTTP handlers and request/response binding
  database/                database connection and migrations
  handlers/                health and readiness handlers
  middleware/              route middleware
  models/                  GORM models
  repositories/            repository interfaces
  repositories/gorm/       GORM-backed repository implementations
  routes/                  route registration
  services/                service interfaces
  services/default/        default service implementations
```

## Configuration

The app uses these environment variables:

| Variable       | Required | Description                     |
| -------------- | -------- | ------------------------------- |
| `DATABASE_URL` | Yes      | PostgreSQL connection string    |
| `JWT_SECRET`   | Yes      | Secret used to sign auth tokens |
| `PORT`         | No       | HTTP port, defaults to `8080`   |

Optional runtime tuning:

| Variable                   | Required | Description                           | Default |
| -------------------------- | -------- | ------------------------------------- | ------- |
| `HTTP_READ_TIMEOUT`        | No       | Maximum time to read the full request | `5s`    |
| `HTTP_READ_HEADER_TIMEOUT` | No       | Maximum time to read request headers  | `2s`    |
| `HTTP_WRITE_TIMEOUT`       | No       | Maximum time to write the response    | `10s`   |
| `HTTP_IDLE_TIMEOUT`        | No       | Maximum keep-alive idle time          | `60s`   |
| `HTTP_SHUTDOWN_TIMEOUT`    | No       | Grace period for graceful shutdown    | `10s`   |
| `AUTH_TOKEN_TTL`           | No       | Signed token lifetime                 | `24h`   |

An example file is included at [.env.example](c:/Users/Trelobarbouni/Documents/GitHub/go-personal-finance-tracker/.env.example).

## Run With Docker

Start the API and PostgreSQL:

```sh
docker-compose up --build
```

The containerized app uses the same environment variables as the native process, so you can move between `go run` and Compose without code changes.

Then verify the service:

```sh
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

## Run Locally Without Docker For The App

You can still use Docker for PostgreSQL and run the Go process directly.

1. Start the database:

```sh
docker-compose up -d db
```

2. Set environment variables:

```powershell
$env:DATABASE_URL="postgres://user:password@localhost:5432/personal_finance_db?sslmode=disable"
$env:JWT_SECRET="dev-secret"
$env:PORT="8080"
```

3. Run the API:

```sh
go run .\cmd\main.go
```

The application now fails fast when required configuration is missing and shuts down gracefully on `Ctrl+C` or `SIGTERM`.

## Quick API Walkthrough

Register:

```sh
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alan","email":"alan@example.com","password":"secure123"}'
```

Login:

```sh
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alan@example.com","password":"secure123"}'
```

Create a budget:

```sh
curl -X POST http://localhost:8080/budgets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"limit":500,"start_date":"2026-03-01T00:00:00Z","end_date":"2026-03-31T23:59:59Z"}'
```

Create a transaction:

```sh
curl -X POST http://localhost:8080/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"type":"expense","amount":42.5,"category_id":1,"date":"2026-03-15T12:00:00Z","note":"Groceries"}'
```

List data:

```sh
curl http://localhost:8080/budgets -H "Authorization: Bearer <token>"
curl http://localhost:8080/transactions -H "Authorization: Bearer <token>"
```

## Testing

Run the full suite:

```sh
go test ./...
```

## Notes

- This project is intentionally scoped as an illustration repository rather than a production-complete finance platform.
- The current API returns model-shaped JSON for budgets and transactions. That is a reasonable next refinement for a follow-up branch.
- Budget enforcement currently validates the transaction against the matching budget limit, but not yet against accumulated spending over a period.

## License

This project is licensed under the MIT License.
