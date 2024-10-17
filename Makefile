build:
	@go build -o bin/book-rental-api cmd/main.go

run: build
	@./bin/book-rental-api

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down
