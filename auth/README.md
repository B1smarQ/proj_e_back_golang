# Authentication Microservice

A clean architecture authentication microservice built with Go, using JWT for authentication and MongoDB for data storage.

## Features

- User registration and login
- JWT-based authentication
- Token validation
- MongoDB integration
- Clean architecture implementation
- HTTP API with Gorilla Mux

## Prerequisites

- Go 1.21 or higher
- MongoDB
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Create a `.env` file in the root directory with the following variables:
   ```
   MONGO_URI=mongodb://localhost:27017
   JWT_KEY=your-secret-key-change-this-in-production
   PORT=8080
   ```

## Running the Service

```bash
go run cmd/api/main.go
```

The service will start on port 8080 by default.

## API Endpoints

### Register User
```http
POST /auth/register
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
    "email": "user@example.com",
    "password": "password123"
}
```

### Validate Token
```http
GET /auth/validate
Authorization: Bearer <your-jwt-token>
```

## Project Structure

```
auth/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   │   └── user.go
│   ├── usecase/
│   │   └── user_usecase.go
│   ├── repository/
│   │   └── mongo/
│   │       └── user_repository.go
│   └── delivery/
│       └── http/
│           └── handler.go
├── pkg/
├── .env
├── go.mod
└── README.md
```

## Clean Architecture

The project follows clean architecture principles:

1. **Domain Layer**: Contains business entities and interfaces
2. **Use Case Layer**: Implements business logic
3. **Interface Layer**: Handles HTTP requests and responses
4. **Infrastructure Layer**: Implements database operations

## Security Considerations

- JWT tokens expire after 24 hours
- Passwords are hashed using bcrypt
- MongoDB connection uses default security settings (configure as needed)
- Environment variables for sensitive data

## License

MIT
```
