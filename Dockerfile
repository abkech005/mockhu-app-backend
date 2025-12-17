# Multi-stage Dockerfile for Mockhu Backend

# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mockhu-api cmd/api/main.go

# Stage 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/mockhu-api .

# Copy migrations
COPY --from=builder /app/migrations ./migrations

# Create storage directories
RUN mkdir -p storage/avatars storage/posts storage/messages

# Expose port
EXPOSE 8085

# Run the application
CMD ["./mockhu-api"]
