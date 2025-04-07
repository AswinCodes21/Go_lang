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
    net-snmp \
    net-snmp-tools \
    procps

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/start.sh .

# Create necessary directories and set permissions
RUN mkdir -p /var/lib/net-snmp && \
    mkdir -p /var/run/net-snmp && \
    mkdir -p /var/log && \
    chown -R root:root /var/lib/net-snmp && \
    chown -R root:root /var/run/net-snmp && \
    chown -R root:root /var/log

# Create SNMP configuration
RUN echo "rocommunity public" > /etc/snmp/snmpd.conf && \
    echo "sysLocation Server Room" >> /etc/snmp/snmpd.conf && \
    echo "sysContact Admin <admin@example.com>" >> /etc/snmp/snmpd.conf && \
    echo "sysServices 72" >> /etc/snmp/snmpd.conf && \
    echo "agentAddress udp:161" >> /etc/snmp/snmpd.conf && \
    echo "master agentx" >> /etc/snmp/snmpd.conf && \
    echo "agentXSocket tcp:localhost:705" >> /etc/snmp/snmpd.conf && \
    echo "persistentDir /var/lib/net-snmp" >> /etc/snmp/snmpd.conf && \
    echo "createUser admin SHA admin AES" >> /etc/snmp/snmpd.conf

RUN chmod +x start.sh

EXPOSE 8081
EXPOSE 161/udp

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/health || exit 1

CMD ["sh", "./start.sh"]
