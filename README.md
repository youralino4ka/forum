# Forum with Admin Panel

A microservices-based forum application with WebSocket chat and admin panel.

## Architecture

The application consists of two microservices:
1. Auth Service - Handles user authentication and management
2. Forum Service - Manages forum functionality and WebSocket chat

### Technology Stack
- Language: Go
- Database: PostgreSQL
- Migrations: golang-migrate
- Communication: gRPC + WebSockets
- Testing: Go testing with mocks
- Documentation: Swagger (swaggo)
- Logging: Zap

## Project Structure
```
.
├── api/                    # Proto files and generated gRPC code
├── cmd/                    # Application entry points
│   ├── auth-service/       # Auth service main
│   └── forum-service/      # Forum service main
├── configs/                # Configuration files
├── docs/                   # Documentation and Swagger files
├── internal/               # Internal packages
│   ├── auth/              # Auth service implementation
│   ├── forum/             # Forum service implementation
│   └── common/            # Shared code between services
├── migrations/            # Database migrations
├── scripts/              # Utility scripts
└── pkg/                 # Public packages that can be imported
```

## Prerequisites
- Go 1.21 or higher
- PostgreSQL
- Docker (optional)

## Setup and Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/forum.git
cd forum
```

2. Install dependencies
```bash
go mod download
```

3. Set up the database
```bash
# Create PostgreSQL database
createdb forum

# Run migrations
migrate -path migrations -database "postgresql://localhost:5432/forum?sslmode=disable" up
```

4. Start the services
```bash
# Start Auth Service
go run cmd/auth-service/main.go

# Start Forum Service
go run cmd/forum-service/main.go
```

## Testing
Run tests with coverage:
```bash
go test -v -cover ./...
```

## API Documentation
Swagger documentation is available at:
- Auth Service: http://localhost:8081/swagger/
- Forum Service: http://localhost:8082/swagger/

## Features
- User authentication and authorization
- Forum discussions and threads
- Real-time WebSocket chat
- Admin panel for user management and content moderation
- Automatic message cleanup for chat 