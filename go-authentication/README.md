# go-authentication
 .env
│   .gitignore
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│   README.md
│
├───cmd
│       main.go
│
├───config
│       config.go
│
├───db
│       db.go
│       migrate.go
│
├───internal
│   ├───delivery
│   │       auth_handler.go
│   │       chat_handler.go
│   │       middleware.go
│   │       websocket_handler.go
│   │       
│   ├───domain
│   │       chat.go
│   │       user.go
│   │
│   ├───repository
│   │       chat_repo.go
│   │       user_repo.go
│   │
│   ├───routes
│   │       routes.go
│   │
│   └───usecase
│           auth_usecase.go
│           chat_usecase.go
│
├───pkg
│       jwt.go
│       logger.go
│       websocket.go
│
└───tests
# Go Authentication and Chat Application

A Go-based authentication and real-time chat application with PostgreSQL database and NATS messaging system.

## Features

- User Authentication (Signup/Login)
- JWT-based Authorization
- Real-time Chat using WebSocket
- Message History
- Conversation Management
- NATS for Message Queue
- PostgreSQL for Data Storage

## Prerequisites

- Go 1.16 or higher
- Docker and Docker Compose
- PostgreSQL
- NATS Server

## Project Structure

```
go-authentication/
├── cmd/
│   └── main.go
├── internal/
│   ├── domain/
│   ├── repository/
│   ├── usecase/
│   ├── delivery/
│   └── services/
├── pkg/
├── tests/
├── Dockerfile
├── docker-compose.yml
├── start.sh
└── README.md
```

## Setup and Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd go-authentication
```

2. Build and run using Docker:
```bash
docker-compose down --rmi all
docker-compose up --build
```

The application will be available at `http://localhost:8081`

## API Endpoints

### Authentication

1. **Signup**
   - Endpoint: `POST /signup`
   - Request Body:
     ```json
     {
         "name": "string",
         "email": "string",
         "password": "string"
     }
     ```

2. **Login**
   - Endpoint: `POST /login`
   - Request Body:
     ```json
     {
         "email": "string",
         "password": "string"
     }
     ```
   - Returns JWT token for authentication

### Messages

1. **Send Message**
   - Endpoint: `POST /messages/send`
   - Headers: `Authorization: Bearer <token>`
   - Request Body:
     ```json
     {
         "to": "integer",
         "content": "string"
     }
     ```

2. **Get Messages**
   - Endpoint: `GET /messages/:user_id`
   - Headers: `Authorization: Bearer <token>`
   - Query Parameters:
     - `limit`: number of messages to return (default: 50)
     - `offset`: number of messages to skip (default: 0)

### Chat

1. **Get Chat Messages**
   - Endpoint: `GET /chat/messages/:user_id`
   - Headers: `Authorization: Bearer <token>`

2. **Get Conversations**
   - Endpoint: `GET /chat/conversations`
   - Headers: `Authorization: Bearer <token>`

3. **WebSocket Connection**
   - Endpoint: `ws://localhost:8081/chat/ws`
   - Headers: `Authorization: Bearer <token>`

## Environment Variables

Create a `.env` file with the following variables:

```env
# Server Config
PORT=8081

# DB Config
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin@123
DB_NAME=golang_project
DB_SSLMODE=disable

# JWT Config
JWT_SECRET=your-secret-key
JWT_EXPIRATION_HOURS=24

# NATS Config
NATS_URL=nats://nats:4222
```

## Testing

Run tests using:
```bash
go test ./...
```

## Docker Configuration

The application uses three services:
1. **app**: The main Go application
2. **postgres**: PostgreSQL database
3. **nats**: NATS messaging server

## Error Handling

The application includes comprehensive error handling for:
- Authentication failures
- Database operations
- Message sending/receiving
- WebSocket connections
- Input validation

## Security Features

- JWT-based authentication
- Password hashing
- Input validation
- Protected routes
- Secure WebSocket connections

## Logging

The application logs:
- Authentication attempts
- Message operations
- Database operations
- WebSocket connections
- Error conditions

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.