# Repository Guidelines

## Project Structure & Module Organization
- Source layout: `cmd/cli` and `cmd/httpserver` are entrypoints; shared packages live in `common/`, `httpserver/`, `database/`, and `metrics/`.
- Tests: co-located `*_test.go` files (e.g., `httpserver/handler_test.go`, `database/database_test.go`).
- Artifacts: binaries output to `build/` via Makefile targets. Dockerfiles: `cli.dockerfile`, `httpserver.dockerfile`.
- Go version: 1.24 (see `go.mod`). Module path: `github.com/flashbots/go-template`.

## Build, Test, and Development Commands
- `make build`: builds CLI and HTTP server into `build/`.
- `make build-cli` / `make build-httpserver`: build individual binaries.
- `make lint`: run formatters and linters (`gofmt`, `gofumpt`, `go vet`, `staticcheck`, `golangci-lint`).
- `make test` / `make test-race`: run tests (optionally with race detector).
- `make fmt`: apply formatting (`gofmt`, `gci`, `gofumpt`) and `go mod tidy`.
- `make cover` / `make cover-html`: coverage summary / HTML report.
- Docker: `make docker-cli`, `make docker-httpserver` build images using the respective Dockerfiles.

## Coding Style & Naming Conventions
- Formatting: run `make fmt` before committing. Go files must pass `golangci-lint` (config in `.golangci.yaml`).
- Style: idiomatic Go; exported identifiers use CamelCase; packages lower-case short names; errors returned, not panicked.
- JSON tags: prefer snake_case (configured via `tagliatelle`).
- Logging: use `common.SetupLogger` (slog) and structured fields; respect `--log-json` and `--log-debug` flags.

## Testing Guidelines
- Framework: standard `testing` with `testify/require` for assertions.
- Run: `make test` (or `go test ./...`).
- Database tests: gated by `RUN_DB_TESTS=1` and a Postgres instance. Example:
  - `docker run -d --name postgres-test -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres`
  - `RUN_DB_TESTS=1 make test`
- Coverage: aim to keep or increase coverage; use `make cover`.

## Commit & Pull Request Guidelines
- Commits: concise, imperative mood; scope prefixes like `ci:`, `docs:`, `fix:`, `feat:` when helpful. Reference issues/PRs.
- PRs: include a clear description, linked issues, test plan (commands run), and any config notes. Ensure `make fmt lint test` pass locally.

## Security & Configuration Tips
- Configuration: prefer flags/env via `common.GetEnv` and CLI options; avoid hardcoding secrets.
- HTTP server: `/debug` (pprof) is opt-in; do not expose publicly. Metrics on `/metrics` (Prometheus format).
