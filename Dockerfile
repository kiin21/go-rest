# ---------- Build stage ----------
FROM golang:1.25.1-alpine AS builder

# Install git and timezone data
RUN apk add --no-cache git ca-certificates tzdata build-base

# Set working directory
WORKDIR /app

# Copy go.mod & go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source
COPY . .

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
RUN $(go env GOPATH)/bin/swag init \
    -g cmd/api/main.go \
    -o docs \
    --parseDependency --parseInternal --parseDepth 10

# Build application binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# ---------- Runtime stage ----------
FROM alpine:latest

# Install CA certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary & docs from builder
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/.env_dev .env_dev

EXPOSE 3000

CMD ["./main"]
