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

# Wait for SNMP to be ready
while ! nc -zu snmp 161; do
  echo "Waiting for SNMP..."
  sleep 1
done

# Start the application
./main 