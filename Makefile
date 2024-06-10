VERSION := $(shell git describe --tags --always --dirty="-dev")

.PHONY: all
all: clean build

.PHONY: v
v:
	@echo "Version: ${VERSION}"

.PHONY: clean
clean:
	rm -rf build/

.PHONY: build-cli
build-cli:
	@mkdir -p ./build
	go build -trimpath -ldflags "-X github.com/flashbots/go-template/common.Version=${VERSION}" -v -o ./build/cli cmd/cli/main.go

.PHONY: build-httpserver
build-httpserver:
	@mkdir -p ./build
	go build -trimpath -ldflags "-X github.com/flashbots/go-template/common.Version=${VERSION}" -v -o ./build/httpserver cmd/httpserver/main.go

.PHONY: test
test:
	go test ./...

.PHONY: test-race
test-race:
	go test -race ./...

.PHONY: lint
lint:
	gofmt -d -s .
	gofumpt -d -extra .
	go vet ./...
	staticcheck ./...
	golangci-lint run
	nilaway ./...

.PHONY: fmt
fmt:
	gofmt -s -w .
	gci write .
	gofumpt -w -extra .
	go mod tidy

.PHONY: gofumpt
gofumpt:
	gofumpt -l -w -extra .

.PHONY: lt
lt: lint test

.PHONY: cover
cover:
	go test -coverprofile=/tmp/go-sim-lb.cover.tmp ./...
	go tool cover -func /tmp/go-sim-lb.cover.tmp
	unlink /tmp/go-sim-lb.cover.tmp

.PHONY: cover-html
cover-html:
	go test -coverprofile=/tmp/go-sim-lb.cover.tmp ./...
	go tool cover -html=/tmp/go-sim-lb.cover.tmp
	unlink /tmp/go-sim-lb.cover.tmp

.PHONY: docker-cli
docker-cli:
	DOCKER_BUILDKIT=1 docker build \
		--platform linux/amd64 \
		--build-arg VERSION=${VERSION} \
		--file cli.dockerfile \
		--tag your-project \
	.

.PHONY: docker-httpserver
docker-httpserver:
	DOCKER_BUILDKIT=1 docker build \
		--platform linux/amd64 \
		--build-arg VERSION=${VERSION} \
		--file httpserver.dockerfile \
		--tag your-project \
	.
