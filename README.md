# go-template

[![Goreport status](https://goreportcard.com/badge/github.com/flashbots/go-template)](https://goreportcard.com/report/github.com/flashbots/go-template)
[![Test status](https://github.com/flashbots/go-template/workflows/Checks/badge.svg?branch=main)](https://github.com/flashbots/go-template/actions?query=workflow%3A%22Checks%22)

Toolbox for new Go projects!

* [`Makefile`](https://github.com/flashbots/go-template/blob/main/Makefile) with `lint`, `test`, `build`, `fmt` and more
* Linting with `gofmt`, `gofumpt`, `go vet`, `staticcheck` and `golangci-lint`
* Logging setup using the Zap logger (with debug and json logging options)
* [GitHub Workflows](.github/workflows/) for linting and testing, as well as releasing and publishing Docker images
* Entry files for [CLI](/cmd/cli/main.go) and [HTTP server](/cmd/httpserver/main.go)
* Webserver with graceful shutdown, implementing `livez`, `readyz` and draining API handlers
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
go install mvdan.cc/gofumpt@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/daixiang0/gci@latest
```

**Lint, test, format**

```bash
make lint
make test
make fmt
```
