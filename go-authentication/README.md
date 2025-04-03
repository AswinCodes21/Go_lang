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
