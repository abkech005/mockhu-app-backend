.PHONY: up down migrate migrate-down run

up:
	docker-compose up -d

down:
	docker-compose down

DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/mockhu_db?sslmode=disable
MIGRATE = migrate -path migrations -database "$(DATABASE_URL)"

migrate:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

run:
	go run cmd/api/main.go
