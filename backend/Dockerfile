# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Create non-root user for building
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Generate go.sum file with proper checksums
RUN go mod tidy

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata && \
    rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S appuser -u 1001

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env file if it exists
COPY --from=builder /app/.env* ./

# Set proper permissions
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
CMD ["./main"]

# Metadata
LABEL maintainer="CMDB Lite Team" \
      version="1.0.0" \
      description="Backend API for CMDB Lite" \
      org.opencontainers.image.title="CMDB Lite Backend" \
      org.opencontainers.image.description="Backend API for CMDB Lite" \
      org.opencontainers.image.version="1.0.0" \
      org.opencontainers.image.authors="CMDB Lite Team" \
      org.opencontainers.image.vendor="CMDB Lite"
