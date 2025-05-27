.PHONY: all proto build test clean migrate-up migrate-down

all: proto build

# Generate protocol buffer code
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/auth/v1/auth.proto \
		api/proto/forum/v1/forum.proto

# Build the services
build:
	go build -o bin/auth-service ./cmd/auth-service
	go build -o bin/forum-service ./cmd/forum-service

# Run tests with coverage
test:
	go test -v -cover ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Database migrations
migrate-up:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/forum?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/forum?sslmode=disable" down

# Generate Swagger documentation
swagger:
	swag init -g cmd/auth-service/main.go -o docs/auth
	swag init -g cmd/forum-service/main.go -o docs/forum

# Run services
run-auth:
	go run cmd/auth-service/main.go

run-forum:
	go run cmd/forum-service/main.go 