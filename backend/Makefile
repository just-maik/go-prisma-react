# Variables
BINARY_NAME=main

# Scripts
push:
	go run github.com/steebchen/prisma-client-go db push

build:
	make push
	go mod tidy
	go build -o ${BINARY_NAME} main.go

dev:
	make push
	go run main.go
