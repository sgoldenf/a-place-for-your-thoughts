include .env

all: run_postgres migration run_server

compose:
	docker-compose up --detach

migration:
	migrate -path migrations -database postgres://$(APP_DB_USER):$(APP_DB_PASSWORD)@localhost:$(DB_PORT)/$(APP_DB)?sslmode=disable up

run_server:
	go run ./cmd/web/
	