# go-template

[![Goreport status](https://goreportcard.com/badge/github.com/flashbots/go-template)](https://goreportcard.com/report/github.com/flashbots/go-template)
[![Test status](https://github.com/flashbots/go-template/workflows/Checks/badge.svg?branch=main)](https://github.com/flashbots/go-template/actions?query=workflow%3A%22Checks%22)

Toolbox and building blocks for new Go projects, to get started quickly and right-footed!

* [`Makefile`](https://github.com/flashbots/go-template/blob/main/Makefile) with `lint`, `test`, `build`, `fmt` and more
* Linting with `gofmt`, `gofumpt`, `go vet`, `staticcheck` and `golangci-lint`
* Logging setup using the [slog logger](https://pkg.go.dev/golang.org/x/exp/slog) (with debug and json logging options)
* [GitHub Workflows](.github/workflows/) for linting and testing, as well as releasing and publishing Docker images
* Entry files for [CLI](/cmd/cli/main.go) and [HTTP server](/cmd/httpserver/main.go)
* Webserver with
  * Graceful shutdown, implementing `livez`, `readyz` and draining API handlers
  * Prometheus metrics
  * Using https://pkg.go.dev/github.com/go-chi/chi/v5 for routing
  * [Urfave](https://cli.urfave.org/) for cli args
* https://github.com/uber-go/nilaway
* See also:
  * Public project setup: https://github.com/flashbots/flashbots-repository-template
  * Repository for common Go utilities: https://github.com/flashbots/go-utils

Pick and choose whatever is useful to you! Don't feel the need to use everything, or even to follow this structure.

---

## Getting started

**Build CLI**

```bash
make build-cli
```

**Build HTTP server**

```bash
make build-httpserver
```

**Install dev dependencies**

```bash
go install mvdan.cc/gofumpt@v0.4.0
go install honnef.co/go/tools/cmd/staticcheck@v0.4.2
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3
go install go.uber.org/nilaway/cmd/nilaway@v0.0.0-20240821220108-c91e71c080b7
```

**Lint, test, format**

```bash
make lint
make test
make fmt
```
