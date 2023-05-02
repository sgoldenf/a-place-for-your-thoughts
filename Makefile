include .env

all: run_postgres migration run_server

compose:
	docker-compose up --detach

migration:
	migrate -path migrations -database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(APP_DB)?sslmode=disable up

run_server:
	go run ./cmd/web/
	