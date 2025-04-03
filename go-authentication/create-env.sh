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
EOF

# Make the script executable
chmod +x /root/create-env.sh 