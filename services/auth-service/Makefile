BIN_DIR=bin
SERVER_BIN=$(BIN_DIR)/auth-service
PROTO_DIR=./proto
USER_SERVICE_PROTO=$(PROTO_DIR)/user-service.proto
AUTH_SERVICE_PROTO=$(PROTO_DIR)/auth-service.proto
OUT_DIR=./proto

IMAGE_NAME=auth-service
CONTAINER_NAME=auth-service-container

.PHONY: all proto build run test docker-build docker-run clean

all: proto build

# Generate Go code from proto files
proto:
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(USER_SERVICE_PROTO)
	protoc --go_out=$(OUT_DIR) --go-grpc_out=$(OUT_DIR) $(AUTH_SERVICE_PROTO)


# Build the Go server
build:
	mkdir -p $(BIN_DIR)
	go build -o $(SERVER_BIN) ./cmd/auth-service

# Run server locally
run:
	go run ./cmd/auth-service

# Run tests with coverage
test:
	go test -v ./... -cover

# Clean build files
clean:
	rm -rf $(BIN_DIR)
	rm -f $(PROTO_DIR)/*.pb.go
