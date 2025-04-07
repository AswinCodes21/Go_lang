# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    curl \
    netcat-openbsd \
    procps

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/start.sh .

# Create necessary directories
RUN mkdir -p /var/log

RUN chmod +x start.sh

EXPOSE 8081
EXPOSE 104/udp

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

CMD ["sh", "./start.sh"]
