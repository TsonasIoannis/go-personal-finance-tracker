# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first for caching dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o personal-finance-tracker cmd/main.go

# Stage 2: Run
FROM scratch

WORKDIR /root/

# Copy built binary from builder stage
COPY --from=builder /app/personal-finance-tracker .

# Set environment variables (optional)
ENV PORT=8080

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./personal-finance-tracker"]