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
- versioned startup migrations with a separate persistence bootstrap
- readiness and health endpoints
- request IDs, structured logs, Prometheus metrics, and OpenTelemetry tracing hooks
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
| POST   | `/api/v1/register` | Register a user and return a token |
| POST   | `/api/v1/login`    | Authenticate a user and return a token |
| GET    | `/health`   | Liveness probe                         |
| GET    | `/ready`    | Readiness probe backed by the database |
| GET    | `/metrics`  | Prometheus metrics endpoint            |

### Protected

These endpoints require `Authorization: Bearer <token>`.

| Method | Endpoint            | Description                                         |
| ------ | ------------------- | --------------------------------------------------- |
| GET    | `/api/v1/transactions`     | List the authenticated user's transactions with `page` and `page_size` |
| POST   | `/api/v1/transactions`     | Create a transaction for the authenticated user |
| DELETE | `/api/v1/transactions/:id` | Delete one of the authenticated user's transactions |
| GET    | `/api/v1/budgets`          | List the authenticated user's budgets with `page` and `page_size` |
| POST   | `/api/v1/budgets`          | Create a budget for the authenticated user    |
| DELETE | `/api/v1/budgets/:id`      | Delete one of the authenticated user's budgets |

Legacy unversioned endpoints remain available for compatibility during the transition to `/api/v1`.

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

Optional tracing:

| Variable                      | Required | Description                                                    | Default                       |
| ----------------------------- | -------- | -------------------------------------------------------------- | ----------------------------- |
| `OTEL_SERVICE_NAME`           | No       | Service name attached to exported traces                       | `go-personal-finance-tracker` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | No       | OTLP/HTTP endpoint for trace export. If empty, tracing is noop | unset                         |
| `OTEL_EXPORTER_OTLP_INSECURE` | No       | Use insecure OTLP transport for non-TLS collectors             | `false`                       |
| `OTEL_TRACES_SAMPLER_ARG`     | No       | Trace sampling ratio between `0` and `1`                       | `1.0`                         |

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
curl http://localhost:8080/metrics
```

### Optional Observability Profile

You can also run an OTLP collector plus Jaeger to simulate a more production-like tracing setup locally without making it part of the default stack.

PowerShell:

```powershell
$env:OTEL_SERVICE_NAME="go-personal-finance-tracker"
$env:OTEL_EXPORTER_OTLP_ENDPOINT="http://otel-collector:4318"
$env:OTEL_EXPORTER_OTLP_INSECURE="true"
$env:OTEL_TRACES_SAMPLER_ARG="0.25"
docker compose --profile observability up --build
```

Then open Jaeger at `http://localhost:16686` and search for traces from `go-personal-finance-tracker`.

If you want to go back to the normal lightweight local stack, unset `OTEL_EXPORTER_OTLP_ENDPOINT` and run plain `docker compose up --build`.

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
$env:OTEL_EXPORTER_OTLP_ENDPOINT=""
```

3. Run the API:

```sh
go run .\cmd\main.go
```

The application now fails fast when required configuration is missing and shuts down gracefully on `Ctrl+C` or `SIGTERM`.

## Observability

The API includes a small observability baseline:

- Every request gets an `X-Request-ID` header. If the client sends one, the API reuses it.
- Request logs are structured JSON and include request ID, route, status, latency, and trace IDs when tracing is enabled.
- `/metrics` exposes Prometheus-format HTTP metrics.
- OpenTelemetry tracing is instrumented in the request pipeline and exports spans when `OTEL_EXPORTER_OTLP_ENDPOINT` is configured.

Example local tracing setup against an OTLP HTTP collector on port `4318`:

```powershell
$env:OTEL_SERVICE_NAME="go-personal-finance-tracker"
$env:OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
$env:OTEL_EXPORTER_OTLP_INSECURE="true"
$env:OTEL_TRACES_SAMPLER_ARG="0.25"
```

If `OTEL_EXPORTER_OTLP_ENDPOINT` is left empty, the tracing hooks remain in place but the process uses a noop tracer provider.

## Quick API Walkthrough

Register:

```sh
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alan","email":"alan@example.com","password":"secure123"}'
```

Login:

```sh
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alan@example.com","password":"secure123"}'
```

Create a budget:

```sh
curl -X POST http://localhost:8080/api/v1/budgets \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"limit":500,"start_date":"2026-03-01T00:00:00Z","end_date":"2026-03-31T23:59:59Z"}'
```

Create a transaction:

```sh
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"type":"expense","amount":42.5,"category_id":1,"date":"2026-03-15T12:00:00Z","note":"Groceries"}'
```

List paginated data:

```sh
curl "http://localhost:8080/api/v1/budgets?page=1&page_size=20" -H "Authorization: Bearer <token>"
curl "http://localhost:8080/api/v1/transactions?page=1&page_size=20" -H "Authorization: Bearer <token>"
```

The versioned list endpoints now respond with a `data` array plus a `pagination` object. Legacy unversioned list endpoints remain array-shaped during the compatibility window.

## Testing

Run the full suite:

```sh
go test ./...
```

## Error Responses

The API uses typed application errors in the service layer and a centralized HTTP error responder in the transport layer.

- services return typed errors such as validation, unauthorized, not found, conflict, and internal
- controllers and middleware delegate error serialization to a shared responder instead of building ad hoc JSON bodies
- error payloads follow a consistent envelope so clients can rely on stable machine-readable codes

Example:

```json
{
  "error": {
    "code": "invalid_budget_limit",
    "message": "budget limit must be greater than zero"
  }
}
```

Successful responses are still resource-specific for now, while error responses are standardized across the API surface.

## Notes

- This project is intentionally scoped as an illustration repository rather than a production-complete finance platform.
- The legacy unversioned list endpoints still return model-shaped arrays for compatibility, while `/api/v1` list endpoints return paginated envelopes.
- Budget enforcement currently validates the transaction against the matching budget limit, but not yet against accumulated spending over a period.

## License

This project is licensed under the MIT License.
