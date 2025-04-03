# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Final stage
FROM alpine:latest

# Install necessary certificates, curl, and wget for healthcheck
RUN apk --no-cache add ca-certificates curl wget dos2unix

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy any necessary config files
COPY --from=builder /app/config ./config

# Copy the start script and make it executable
COPY start.sh .
# Convert Windows line endings to Unix and make executable
RUN dos2unix start.sh && chmod +x start.sh

# Expose port (using the port from .env)
EXPOSE 8081

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

# Command to run the application
CMD ["./start.sh"]
