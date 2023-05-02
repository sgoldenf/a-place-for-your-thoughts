include .env

all: compose run_server

compose:
	docker-compose up --detach

migration:
	migrate -path migrations -database postgres://$(APP_DB_USER):$(APP_DB_PASSWORD)@localhost:5432/$(APP_DB)?sslmode=disable up

SERVER=server

$(SERVER):
	go build -o server ./cmd/web/

run_server: $(SERVER)
	./$(SERVER)

test:
	go test -v -cover ./...

clean:
	rm -f server

	