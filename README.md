# Product API

A small Go service implementing product creation and partial updates with Clean Architecture, PostgreSQL, OpenAPI documentation, and tests at domain, use-case, component, and repository-integration levels.

## Architecture

```text
HTTP handler -> use case -> repository interface -> PostgreSQL adapter
```

- `internal/domain`: entities, PATCH value objects, and business validation
- `internal/usecase`: application orchestration
- `internal/repository`: repository port and PostgreSQL adapter
- `internal/httpapi`: transport DTOs, handlers, router, and embedded OpenAPI spec
- `cmd/api`: composition root and manual dependency injection

Manual constructor injection keeps dependencies explicit without adding a DI framework.

## Start the service

Docker and Docker Compose are required.

```bash
docker compose up --build
```

The service is available at `http://localhost:8080`:

- Health check: `GET /health`
- Swagger UI: `GET /api-docs`
- OpenAPI document: `GET /api-docs/openapi.yaml`

To run against an existing PostgreSQL instance:

```bash
cp .env.example .env
set -a; source .env; set +a
go run ./cmd/api
```

Apply `migrations/001_create_products.sql` before starting the service outside Docker Compose.

## API examples

Create a product:

```bash
curl -i http://localhost:8080/product \
  -H 'Content-Type: application/json' \
  -d '{"name":"Keyboard","description":"Mechanical","price":1200,"sale_price":999}'
```

Patch only selected fields. An omitted field remains unchanged; `null` clears a nullable field:

```bash
curl -i -X PATCH http://localhost:8080/product/1 \
  -H 'Content-Type: application/json' \
  -d '{"description":null,"price":1100}'
```

`name` and `price` may be omitted but cannot be `null`. `description` and `sale_price` may be omitted or explicitly set to `null`.

## Tests

Run unit and component tests:

```bash
go test ./...
```

Run the PostgreSQL repository integration test while the Compose database is running:

```bash
docker compose up -d postgres
make test-integration
```

The test layers are:

- Domain unit tests: business validation
- Use-case unit tests: orchestration with a repository stub
- Component test: HTTP router through use case and an in-memory repository
- Repository integration test: PostgreSQL adapter against a real database
