# go-template

[![Goreport status](https://goreportcard.com/badge/github.com/flashbots/go-template)](https://goreportcard.com/report/github.com/flashbots/go-template)
[![Test status](https://github.com/flashbots/go-template/actions/workflows/checks.yml/badge.svg?branch=main)](https://github.com/flashbots/go-template/actions?query=workflow%3A%22Checks%22)

Toolbox and building blocks for new Go projects, to get started quickly and right-footed!

Pick and choose whatever is useful to you! Don't feel the need to use everything, or even to follow this structure.

## What's Included

This template provides two entry points:

- **CLI application** ([`cmd/cli/main.go`](/cmd/cli/main.go)) - Command-line tool using [urfave/cli](https://cli.urfave.org/)
- **HTTP server** ([`cmd/httpserver/main.go`](/cmd/httpserver/main.go)) - Web server with graceful shutdown, health checks, and metrics

### Features

- [`Makefile`](https://github.com/flashbots/go-template/blob/main/Makefile) with `lint`, `test`, `build`, `fmt` and more
- Linting with `gofmt`, `gofumpt`, `go vet`, `staticcheck` and `golangci-lint`
- Logging setup using the [slog logger](https://pkg.go.dev/golang.org/x/exp/slog) (with debug and json logging options)
- [GitHub Workflows](.github/workflows/) for linting and testing, as well as releasing and publishing Docker images
- Webserver with graceful shutdown, health probes, and Prometheus metrics
- Postgres database with migrations

---

## Quick Start

**Build and run the HTTP server:**

```bash
make build-httpserver
./build/httpserver --listen-addr 127.0.0.1:8080 --metrics-addr 127.0.0.1:8090
```

**Build and run the CLI:**

```bash
make build-cli
./build/cli
```

---

## Project Structure

| Directory              | Description                                                |
| ---------------------- | ---------------------------------------------------------- |
| `cmd/cli/`             | CLI application entry point (urfave/cli)                   |
| `cmd/httpserver/`      | HTTP server entry point                                    |
| `httpserver/`          | HTTP server implementation (chi router, graceful shutdown) |
| `database/`            | Postgres database layer using sqlx                         |
| `database/migrations/` | Database migrations (run automatically on connection)      |
| `metrics/`             | Prometheus metrics (VictoriaMetrics-based)                 |
| `common/`              | Shared utilities (structured logging)                      |

---

## HTTP Server Endpoints

The server runs two HTTP servers: main API (default `:8080`) and metrics (default `:8090`).

| Endpoint   | Port | Description                                        |
| ---------- | ---- | -------------------------------------------------- |
| `/api`     | 8080 | Main API endpoint                                  |
| `/livez`   | 8080 | Liveness probe for health checks                   |
| `/readyz`  | 8080 | Readiness probe for health checks                  |
| `/drain`   | 8080 | Enable drain mode (for graceful shutdown)          |
| `/undrain` | 8080 | Disable drain mode                                 |
| `/debug/*` | 8080 | pprof debug endpoints (when `--pprof` flag is set) |
| `/metrics` | 8090 | Prometheus metrics                                 |

### CLI Flags

| Flag              | Default          | Description                           |
| ----------------- | ---------------- | ------------------------------------- |
| `--listen-addr`   | `127.0.0.1:8080` | Address for API server                |
| `--metrics-addr`  | `127.0.0.1:8090` | Address for Prometheus metrics        |
| `--log-json`      | `false`          | Log in JSON format                    |
| `--log-debug`     | `false`          | Enable debug logging                  |
| `--log-uid`       | `false`          | Add UUID to all log messages          |
| `--log-service`   | `your-project`   | Service name in logs                  |
| `--pprof`         | `false`          | Enable pprof debug endpoint           |
| `--drain-seconds` | `45`             | Seconds to wait in drain HTTP request |

---

## Development

### Build Commands

```bash
make build-cli        # Build CLI binary to ./build/cli
make build-httpserver # Build HTTP server binary to ./build/httpserver
make build            # Build all binaries
```

### Lint, Test, Format

```bash
make lint   # Run all linters (gofmt, gofumpt, go vet, staticcheck, golangci-lint)
make test   # Run all tests
make fmt    # Format code (gofmt, gci, gofumpt, go mod tidy)
make lt     # Run both lint and test
```

### Install Dev Dependencies

```bash
go install mvdan.cc/gofumpt@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install go.uber.org/nilaway/cmd/nilaway@latest
go install github.com/daixiang0/gci@latest
```

### Database Tests

Database tests require a running Postgres instance and the `RUN_DB_TESTS` environment variable:

```bash
# Start the database
docker run -d --name postgres-test -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres postgres

# Run the tests
RUN_DB_TESTS=1 make test

# Stop the database
docker rm -f postgres-test
```

---

## Environment Variables

| Variable               | Description                                            |
| ---------------------- | ------------------------------------------------------ |
| `RUN_DB_TESTS`         | Set to `1` to run database integration tests           |
| `DB_DONT_APPLY_SCHEMA` | Set to skip automatic migration on database connection |

---

## Related Resources

- [Flashbots Repository Template](https://github.com/flashbots/flashbots-repository-template) - Public project setup
- [go-utils](https://github.com/flashbots/go-utils) - Common Go utilities
- [goperf.dev](https://goperf.dev) - Advanced Golang knowledge, tips & tricks