# ---- Dependencies Stage ----
FROM golang:1.25.1-alpine AS deps

# Install build tools
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency definition files
COPY go.work go.work.sum ./
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY services/starter-service/go.mod services/starter-service/go.sum ./services/starter-service/
COPY services/notification-service/go.mod services/notification-service/go.sum ./services/notification-service/

# Download dependencies
RUN go work sync
RUN go mod download

# ---- Builder Stage ----
FROM deps AS builder

# Copy source code
COPY pkg/ ./pkg/
COPY services/notification-service/ ./services/notification-service/

# Build the application
WORKDIR /app/services/notification-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/main.go

# ---- Final Stage ----
FROM alpine:latest

# Install required certificates
RUN apk --no-cache add ca-certificates

# Create a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /home/appuser

# Copy built binary
COPY --from=builder /app/main .
# Copy environment defaults for runtime
COPY services/notification-service/.env_dev .env_dev

# Change ownership
RUN chown -R appuser:appgroup /home/appuser

# Switch to non-root user
USER appuser

EXPOSE 8081

CMD ["./main"]
