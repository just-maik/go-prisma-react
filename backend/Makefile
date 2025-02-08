# Variables
BINARY_NAME=backend
GO_VERSION=1.23
MAIN_FILE=main.go
TEST_DIR=./...
SWAGGER_DIR=./docs
PRISMA_DIR=./prisma

# Ensure correct Go version
.PHONY: check-go-version
check-go-version:
	@go version | grep -q "go${GO_VERSION}" || (echo "Error: Required Go version is ${GO_VERSION}" && exit 1)

# Database operations
.PHONY: db-push
db-push:
	go run github.com/steebchen/prisma-client-go db push

.PHONY: prisma-generate
prisma-generate:
	go run github.com/steebchen/prisma-client-go generate

# Build operations
.PHONY: build
build: check-go-version db-push
	go mod tidy
	go build -o ${BINARY_NAME} ${MAIN_FILE}

# Development
.PHONY: dev
dev: check-go-version db-push swagger
	go run ${MAIN_FILE}

# Testing
.PHONY: test
test: check-go-version
	go test -v ${TEST_DIR}

.PHONY: test-coverage
test-coverage: check-go-version
	go test -v -coverprofile=coverage.out ${TEST_DIR}
	go tool cover -html=coverage.out

# Documentation
.PHONY: swagger
swagger:
	swag init -g ${MAIN_FILE} -d . -o ${SWAGGER_DIR}

# Cleanup
.PHONY: clean
clean:
	rm -f ${BINARY_NAME}
	rm -f coverage.out
	rm -rf ${SWAGGER_DIR}/docs.go ${SWAGGER_DIR}/swagger.json ${SWAGGER_DIR}/swagger.yaml
	rm -rf ${PRISMA_DIR}/db

# Run integration tests
.PHONY: test-integration
test-integration: check-go-version
	go test -v ./cmd/test/...

# All-in-one development setup
.PHONY: setup
setup: check-go-version prisma-generate swagger
	go mod tidy
	go mod download

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build            - Build the application"
	@echo "  clean            - Remove built binary and generated files"
	@echo "  db-push          - Push database schema changes"
	@echo "  dev              - Run the application in development mode"
	@echo "  prisma-generate  - Generate Prisma client"
	@echo "  setup           - Initial project setup"
	@echo "  swagger         - Generate Swagger documentation"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  test-integration - Run integration tests"
