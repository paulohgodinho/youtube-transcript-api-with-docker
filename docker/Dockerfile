# Multi-stage Dockerfile for YouTube Transcript API

# Stage 1: Builder - Build the Go application
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY server/go.mod server/go.sum* ./

# Download dependencies
RUN go mod download

# Copy all server source files
COPY server/*.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o youtube-api .

# Stage 2: Runtime - Python Alpine image
FROM python:3.11-alpine

# Set working directory
WORKDIR /app

# Install system dependencies
RUN apk add --no-cache \
    ca-certificates \
    curl

# Install Python package
RUN pip install --no-cache-dir youtube-transcript-api==1.2.4

# Copy binary from builder stage
COPY --from=builder /build/youtube-api /app/youtube-api

# Make binary executable
RUN chmod +x /app/youtube-api

# Expose port 8080
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Default environment variables
ENV PYTHON_BIN=python3 \
    SERVER_PORT=8080 \
    REQUEST_TIMEOUT=30s

# Run the server
ENTRYPOINT ["/app/youtube-api"]
CMD ["-port", "8080"]
