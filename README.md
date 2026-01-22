# go-template

[![Goreport status](https://goreportcard.com/badge/github.com/flashbots/go-template)](https://goreportcard.com/report/github.com/flashbots/go-template)
[![Test status](https://github.com/flashbots/go-template/actions/workflows/checks.yml/badge.svg?branch=main)](https://github.com/flashbots/go-template/actions?query=workflow%3A%22Checks%22)

Toolbox and building blocks for new Go projects, to get started quickly and right-footed!

## Features

- [`Makefile`](https://github.com/flashbots/go-template/blob/main/Makefile) with `lint`, `test`, `build`, `fmt` and more
- Linting with `gofmt`, `gofumpt`, `go vet`, `staticcheck` and `golangci-lint`
- Logging setup using the [slog logger](https://pkg.go.dev/log/slog) (via [httplog](https://github.com/go-chi/httplog))
- [GitHub Workflows](.github/workflows/) for linting and testing, as well as releasing and publishing Docker images
- Entry files for [CLI](/cmd/cli/main.go) and [HTTP server](/cmd/httpserver/main.go)
- HTTP server with:
  - Graceful shutdown with configurable drain duration
  - Health check endpoints (`/livez`, `/readyz`, `/drain`, `/undrain`)
  - Prometheus metrics (via [VictoriaMetrics](https://github.com/VictoriaMetrics/metrics))
  - [chi](https://github.com/go-chi/chi) router
  - [urfave/cli](https://cli.urfave.org/) for CLI argument parsing
  - Optional pprof debug endpoints
- Postgres database layer with:
  - [sqlx](https://github.com/jmoiron/sqlx) for database operations
  - [sql-migrate](https://github.com/rubenv/sql-migrate) for migrations
  - In-memory migration definitions
- Static analysis with [nilaway](https://github.com/uber-go/nilaway)

See also:
- Public project setup: https://github.com/flashbots/flashbots-repository-template
- Repository for common Go utilities: https://github.com/flashbots/go-utils
- Advanced Go performance tips: https://goperf.dev

Pick and choose whatever is useful to you! Don't feel the need to use everything, or even to follow this structure.

---

## Project Structure

```
.
├── cmd/
│   ├── cli/              # CLI application entry point
│   └── httpserver/       # HTTP server entry point
├── common/               # Shared utilities (logging, version)
├── database/             # Postgres database layer
│   ├── migrations/       # Database migrations (in-memory)
│   └── vars/             # Database table names
├── httpserver/           # HTTP server implementation
└── metrics/              # Prometheus metrics and middleware
```

---

## Getting Started

### Build

```bash
# Build CLI
make build-cli

# Build HTTP server
make build-httpserver

# Build all binaries
make build
```

Binaries are output to `./build/`.

### Run

```bash
# Run the HTTP server
./build/httpserver --help

# Run the CLI
./build/cli --help
```

### HTTP Server Options

| Flag | Default | Description |
|------|---------|-------------|
| `--listen-addr` | `127.0.0.1:8080` | Address for the API server |
| `--metrics-addr` | `127.0.0.1:8090` | Address for Prometheus metrics |
| `--log-json` | `false` | Output logs in JSON format |
| `--log-debug` | `false` | Enable debug logging |
| `--log-uid` | `false` | Add UUID to all log messages |
| `--log-service` | `your-project` | Service name in logs |
| `--pprof` | `false` | Enable pprof debug endpoint at `/debug` |
| `--drain-seconds` | `45` | Seconds to wait during drain |

### HTTP Endpoints

| Endpoint | Description |
|----------|-------------|
| `/api` | Main API endpoint |
| `/livez` | Liveness probe (always returns 200 if server is running) |
| `/readyz` | Readiness probe (returns 200 when ready to serve traffic) |
| `/drain` | Initiate graceful drain (for load balancer compatibility) |
| `/undrain` | Cancel drain and resume serving traffic |
| `/debug/*` | pprof endpoints (when `--pprof` is enabled) |
| `/metrics` | Prometheus metrics (on metrics server port) |

---

## Development

### Install Dev Dependencies

```bash
go install mvdan.cc/gofumpt@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install go.uber.org/nilaway/cmd/nilaway@latest
go install github.com/daixiang0/gci@latest
```

### Lint, Test, Format

```bash
# Run all linters
make lint

# Run tests
make test

# Run tests with race detector
make test-race

# Format code
make fmt

# Run both lint and test
make lt

# Run tests with coverage
make cover
```

### Database Tests

Database tests require a running Postgres instance and will only run when `RUN_DB_TESTS=1` is set:

```bash
# Start Postgres
docker run -d --name postgres-test \
  -p 5432:5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=postgres \
  postgres

# Run tests including database tests
RUN_DB_TESTS=1 make test

# Stop Postgres
docker rm -f postgres-test
```

### Environment Variables

| Variable | Description |
|----------|-------------|
| `RUN_DB_TESTS` | Set to `1` to enable database tests |
| `DB_DONT_APPLY_SCHEMA` | Set to skip automatic migration on database connect |

---

## Docker

```bash
# Build CLI Docker image
make docker-cli

# Build HTTP server Docker image
make docker-httpserver
```

---

## Architecture Notes

### HTTP Server

The server runs two HTTP servers simultaneously:
- **Main API server** (default `:8080`) - Serves API endpoints and health checks
- **Metrics server** (default `:8090`) - Serves Prometheus metrics

This separation allows metrics to be scraped without affecting API traffic and enables different network policies for each.

### Graceful Shutdown

The server implements graceful shutdown with the following sequence:
1. Receive termination signal (SIGTERM/SIGINT)
2. Stop accepting new connections
3. Wait for in-flight requests to complete (configurable timeout)
4. Shutdown metrics server
5. Exit

The `/drain` and `/undrain` endpoints allow manual control over readiness, useful for load balancer integration during deployments.

### Database Migrations

Migrations are defined as Go code in `database/migrations/` and are automatically applied on database connection unless `DB_DONT_APPLY_SCHEMA` is set. This approach keeps migrations version-controlled alongside the code.

---

## License

See [LICENSE](LICENSE) file for details.
