# go-template

[![Test status](https://github.com/flashbots/go-template/workflows/Checks/badge.svg)](https://github.com/flashbots/go-template/actions?query=workflow%3A%22Checks%22)

Starting point for new Go projects:

* Entry files:
  * CLI utility: [`cmd/cli/main.go`](https://github.com/flashbots/go-template/blob/main/cmd/cli/main.go)
  * HTTP server: [`cmd/httpserver/main.go`](https://github.com/flashbots/go-template/blob/main/cmd/httpserver/main.go)
* Logging setup using the Zap logger (with debug and json logging options)
* Linting (with lint, go vet and staticcheck) & tests
* GitHub Workflow for linting and testing
* [`Makefile`](https://github.com/flashbots/go-template/blob/main/Makefile)
* Setup for building and publishing Docker images

For public projects also take a look at https://github.com/flashbots/flashbots-repository-template

We also have a repository for common Go utilities: https://github.com/flashbots/go-utils

---

## Getting started

### Install dependencies

```bash
go mod download -x
```

### Test

```bash
make test
make test-race
make lint
```

### Build CLI

```bash
make build-cli
```

### Build HTTP server

```bash
make build-httpserver
```
