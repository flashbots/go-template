This file provides guidance to LLMs when working with code in this repository.

## Build Commands

```bash
make build-cli        # Build CLI binary to ./build/cli
make build-httpserver # Build HTTP server binary to ./build/httpserver
make build            # Build all binaries
```

## Test Commands

```bash
make test             # Run all tests
make test-race        # Run tests with race detector
go test ./... -run TestName  # Run a single test

# Database tests require a running Postgres instance:
docker run -d --name postgres-test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres postgres
RUN_DB_TESTS=1 make test
```

## Lint and Format

```bash
make lint   # Run all linters (gofmt, gofumpt, go vet, staticcheck, golangci-lint)
make fmt    # Format code (gofmt, gci, gofumpt, go mod tidy)
make lt     # Run both lint and test
```

## Architecture

This is a Go project template with two entry points:

- **cmd/cli/main.go** - CLI application entry point using urfave/cli
- **cmd/httpserver/main.go** - HTTP server entry point with graceful shutdown

### Key Packages

- **httpserver/** - HTTP server with chi router, includes `/livez`, `/readyz`, `/drain`, `/undrain` endpoints, metrics middleware, and optional pprof
- **database/** - Postgres database layer using sqlx with sql-migrate for migrations
- **database/migrations/** - In-memory migrations registered in `Migrations` variable
- **metrics/** - VictoriaMetrics-based Prometheus metrics with HTTP middleware
- **common/** - Shared utilities including structured logging setup (slog-based via httplog)

### HTTP Server Pattern

The server runs two HTTP servers: main API (default :8080) and metrics (default :8090). Supports graceful shutdown with configurable drain duration for load balancer compatibility.

### Database Migrations

Migrations are defined as Go code in `database/migrations/` and registered in `migrations.Migrations`. They run automatically on database connection unless `DB_DONT_APPLY_SCHEMA` env var is set.
