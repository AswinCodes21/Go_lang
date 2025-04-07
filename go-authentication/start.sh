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

# Start SNMP daemon with debug logging
echo "Starting SNMP daemon with debug logging..."
snmpd -f -DALL -Lf /var/log/snmpd.log -Lo &

# Wait for SNMP daemon to start
sleep 2

# Check if SNMP daemon is running
if ! pgrep snmpd > /dev/null; then
    echo "Error: SNMP daemon failed to start"
    echo "Checking SNMP log file..."
    cat /var/log/snmpd.log 2>/dev/null || echo "No SNMP log file found"
    echo "Checking process list..."
    ps aux
    exit 1
fi

echo "SNMP daemon started successfully"
echo "SNMP daemon process info:"
ps aux | grep snmpd

# Start the application
echo "Starting application..."
./main 