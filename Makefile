.PHONY: run test test-integration compose-up compose-down

run:
	go run ./cmd/api

test:
	go test ./...

test-integration:
	TEST_DATABASE_URL=$${TEST_DATABASE_URL:-postgres://postgres:postgres@localhost:5432/products?sslmode=disable} go test -tags=integration ./internal/repository/postgres

compose-up:
	docker compose up --build

compose-down:
	docker compose down -v
