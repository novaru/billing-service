BINARY_NAME=billing-service

build:
	go build -o ./bin/$(BINARY_NAME) ./cmd/server

