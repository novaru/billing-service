BINARY_NAME=billing-service
DRIVER=postgres
MIGRATIONS_DIR=./migrations

# loads .env file if it exists
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build:
	go build -o ./bin/$(BINARY_NAME) ./cmd/server

.PHONY: migrate
migrate:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DATABASE_URL)" up

.PHONY: rollback
rollback:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DATABASE_URL)" down

.PHONY: up
up:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DATABASE_URL)" up

.PHONY: reset
reset:
	goose -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DATABASE_URL)" reset

.PHONY: new
new:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make new name=add_users_table"; \
		exit 1; \
	fi
	goose -dir $(MIGRATIONS_DIR) create $(name) sql

