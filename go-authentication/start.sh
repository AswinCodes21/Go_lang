#!/bin/bash

# Create .env file
cat > .env << EOL
PORT=8081
NATS_URL=nats://nats:4222
DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=admin@123
DB_NAME=golang_project
DB_PORT=5432
DB_SSLMODE=disable
JWT_SECRET=UlVwZFpYbGlzN2N3djd4b2lLMjV6OVF0QzM3TkFqQkY=
JWT_EXPIRATION_HOURS=24
EOL

# Function to check if NATS is ready
check_nats_ready() {
    local max_retries=30
    local retry_count=0
    
    echo "Waiting for NATS server to be ready..."
    
    while [ $retry_count -lt $max_retries ]; do
        # Check if the port is open
        if nc -z nats 4222; then
            echo "NATS server is ready!"
            return 0
        fi
        
        echo "NATS server not ready yet (attempt $((retry_count + 1))/$max_retries)..."
        retry_count=$((retry_count + 1))
        sleep 2
    done
    
    echo "Error: NATS server did not become ready in time"
    return 1
}

# Wait for NATS to be ready
if ! check_nats_ready; then
    echo "Failed to connect to NATS server. Exiting..."
    exit 1
fi

# Add a small delay to ensure NATS is fully ready
sleep 5

# Start the application
echo "Starting application..."
./main 