include .env
MIGRATIONS_DIR := ./migrations

.PHONY: migrate-create
migration: 
	@echo "Creating new migration..."
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) \
	-seq $(name)

.PHONY: migrate-up
migrate-up: 
	@echo "Applying migrations..."
	@migrate -path=./migrations \
	-database=$(PG_ADDR) up

.PHONY: migrate-down
migrate-down: 
	@echo "Reverting last migration..."
	@migrate -path=./migrations \
	-database=$(PG_ADDR) down