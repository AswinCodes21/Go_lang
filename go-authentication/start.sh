#!/bin/sh

# Create .env file with environment variables
cat > /root/.env << EOF
# Server Config
PORT=8081

# DB Config 
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin@123
DB_NAME=golang_project 
DB_SSLMODE=disable

# JWT Secret Key
JWT_SECRET=UlVwZFpYbGlzN2N3djd4b2lLMjV6OVF0QzM3TkFqQkY=
JWT_EXPIRATION_HOURS=24

# NATS Config
NATS_URL=nats://nats:4222
EOF

# Wait for NATS to be ready
echo "Waiting for NATS server to be ready..."
max_retries=30
retry_count=0
while [ $retry_count -lt $max_retries ]; do
  if wget --no-verbose --tries=1 --spider http://nats:8222/healthz; then
    echo "NATS server is ready!"
    break
  fi
  echo "NATS server not ready yet, waiting... ($(($retry_count + 1))/$max_retries)"
  retry_count=$((retry_count + 1))
  sleep 2
done

if [ $retry_count -eq $max_retries ]; then
  echo "NATS server did not become ready in time, but continuing anyway..."
fi

# Run the main application
exec ./main 