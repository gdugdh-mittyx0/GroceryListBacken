.PHONY: all

all: run

run:
	go run cmd/server/main.go

swag:
	swag init -g cmd/server/main.go -o ./swag
