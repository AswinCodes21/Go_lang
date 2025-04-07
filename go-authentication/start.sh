#!/bin/sh

# Wait for postgres to be ready
while ! nc -z postgres 5432; do
  echo "Waiting for postgres..."
  sleep 1
done

# Wait for NATS to be ready
while ! nc -z nats 4222; do
  echo "Waiting for NATS..."
  sleep 1
done

# Create log directory if it doesn't exist
mkdir -p /var/log

# Start the application
echo "Starting application..."
./main 