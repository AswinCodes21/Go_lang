# Build stage
FROM golang:1.21-alpine AS builder

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
    netcat-openbsd

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/start.sh .

RUN chmod +x start.sh

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

CMD ["sh", "./start.sh"]
