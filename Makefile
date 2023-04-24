all: run_postgres migration run_server

run_postgres:
	docker-compose up --detach

migration:
	 migrate -path migrations -database 'postgres://sgoldenf:sgoldenf@localhost:5432/blog?sslmode=disable' up

run_server:
	go run ./cmd/web/
	