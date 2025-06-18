# Heavily inspired by Lighthouse: https://github.com/sigp/lighthouse/blob/stable/Makefile
# and Reth: https://github.com/paradigmxyz/reth/blob/main/Makefile
.DEFAULT_GOAL := help

VERSION := $(shell git describe --tags --always --dirty="-dev")

##@ Help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: v
v: ## Show the version
	@echo "Version: ${VERSION}"

##@ Build

.PHONY: clean
clean: package-clean ## Clean the build directory
	rm -rf build/

.PHONY: build-cli
build-cli: ## Build the CLI
	@mkdir -p ./build
	go build -trimpath -ldflags "-X github.com/flashbots/go-template/common.Version=${VERSION}" -v -o ./build/cli cmd/cli/main.go

.PHONY: build-httpserver
build-httpserver: ## Build the HTTP server
	@mkdir -p ./build
	go build -trimpath -ldflags "-X github.com/flashbots/go-template/common.Version=${VERSION}" -v -o ./build/httpserver cmd/httpserver/main.go

.PHONY: build
build: build-cli build-httpserver ## Build all binaries
	@echo "Binaries built in ./build/"

##@ Test & Development

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-race
test-race: ## Run tests with race detector
	go test -race ./...

.PHONY: lint
lint: ## Run linters
	gofmt -d -s .
	gofumpt -d -extra .
	go vet ./...
	staticcheck ./...
	golangci-lint run
	# nilaway ./...

.PHONY: fmt
fmt: ## Format the code
	gofmt -s -w .
	gci write .
	gofumpt -w -extra .
	go mod tidy

.PHONY: gofumpt
gofumpt: ## Run gofumpt
	gofumpt -l -w -extra .

.PHONY: lt
lt: lint test ## Run linters and tests

.PHONY: cover
cover: ## Run tests with coverage
	go test -coverprofile=/tmp/go-sim-lb.cover.tmp ./...
	go tool cover -func /tmp/go-sim-lb.cover.tmp
	unlink /tmp/go-sim-lb.cover.tmp

.PHONY: cover-html
cover-html: ## Run tests with coverage and open the HTML report
	go test -coverprofile=/tmp/go-sim-lb.cover.tmp ./...
	go tool cover -html=/tmp/go-sim-lb.cover.tmp
	unlink /tmp/go-sim-lb.cover.tmp

.PHONY: docker-cli
docker-cli: ## Build the CLI Docker image
	DOCKER_BUILDKIT=1 docker build \
		--platform linux/amd64 \
		--build-arg VERSION=${VERSION} \
		--file cli.dockerfile \
		--tag your-project \
	.

.PHONY: docker-httpserver
docker-httpserver: ## Build the HTTP server Docker image
	DOCKER_BUILDKIT=1 docker build \
		--platform linux/amd64 \
		--build-arg VERSION=${VERSION} \
		--file httpserver.dockerfile \
		--tag your-project \
	.

##@ Packaging

.PHONY: package-build
package-build: ## Build packages (without releasing)
	@echo "Building packages..."
	@goreleaser build --snapshot --clean
	@echo "✅ Packages built in dist/"

.PHONY: package-local
package-local: ## Build packages locally for testing
	@echo "Creating local release packages..."
	@goreleaser release --snapshot --clean
	@echo "✅ Release packages created in dist/"
	@echo "📦 Created packages:"
	@find dist/ -name "*.deb" -o -name "*.rpm" -o -name "*.tar.gz" | sort

.PHONY: package-test-reproducible
package-test-reproducible: ## Test reproducible builds
	@echo "🔄 Testing reproducible builds..."
	@mkdir -p ./test-reproducible
	@echo "  Building first version (with packages)..."
	@if goreleaser release --snapshot --clean >/dev/null 2>&1; then \
		echo "    ✅ First build completed"; \
		cp -r ./dist ./test-reproducible/build1; \
	else \
		echo "❌ First build failed"; \
		echo "Running with verbose output:"; \
		goreleaser release --snapshot --clean; \
		rm -rf ./test-reproducible; \
		exit 1; \
	fi
	@sleep 2
	@echo "  Building second version (with packages)..."
	@if goreleaser release --snapshot --clean >/dev/null 2>&1; then \
		echo "    ✅ Second build completed"; \
		cp -r ./dist ./test-reproducible/build2; \
	else \
		echo "❌ Second build failed"; \
		echo "Running with verbose output:"; \
		goreleaser release --snapshot --clean; \
		rm -rf ./test-reproducible; \
		exit 1; \
	fi
	@echo "  Comparing packages and binaries..."
	@BUILD1_DEBS=$$(find ./test-reproducible/build1 -name "*.deb" | wc -l); \
	BUILD2_DEBS=$$(find ./test-reproducible/build2 -name "*.deb" | wc -l); \
	BUILD1_BINS=$$(find ./test-reproducible/build1 -type f -name "go-template-*" | wc -l); \
	BUILD2_BINS=$$(find ./test-reproducible/build2 -type f -name "go-template-*" | wc -l); \
	echo "    Found $$BUILD1_DEBS .deb packages and $$BUILD1_BINS binaries in first build"; \
	echo "    Found $$BUILD2_DEBS .deb packages and $$BUILD2_BINS binaries in second build"; \
	if [ "$$BUILD1_DEBS" -eq 0 ] && [ "$$BUILD1_BINS" -eq 0 ]; then \
		echo "❌ No build artifacts found in first build"; \
		find ./test-reproducible/build1 -type f | head -10; \
		rm -rf ./test-reproducible; \
		exit 1; \
	fi
	@echo "  Comparing binary checksums..."
	@find ./test-reproducible/build1 -type f -name "go-template-*" -exec sha256sum {} \; | sed 's|./test-reproducible/build1/||' | sort > ./test-reproducible/checksums1_bins.txt
	@find ./test-reproducible/build2 -type f -name "go-template-*" -exec sha256sum {} \; | sed 's|./test-reproducible/build2/||' | sort > ./test-reproducible/checksums2_bins.txt
	@echo "  Comparing package checksums..."
	@find ./test-reproducible/build1 -name "*.deb" -exec sha256sum {} \; | sed 's|./test-reproducible/build1/||' | sort > ./test-reproducible/checksums1_debs.txt
	@find ./test-reproducible/build2 -name "*.deb" -exec sha256sum {} \; | sed 's|./test-reproducible/build2/||' | sort > ./test-reproducible/checksums2_debs.txt
	@if diff ./test-reproducible/checksums1_bins.txt ./test-reproducible/checksums2_bins.txt >/dev/null 2>&1; then \
		BINS_MATCH=true; \
	else \
		BINS_MATCH=false; \
	fi; \
	if diff ./test-reproducible/checksums1_debs.txt ./test-reproducible/checksums2_debs.txt >/dev/null 2>&1; then \
		DEBS_MATCH=true; \
	else \
		DEBS_MATCH=false; \
	fi; \
	if [ "$$BINS_MATCH" = "true" ] && [ "$$DEBS_MATCH" = "true" ]; then \
		echo "✅ Both binaries and packages are reproducible!"; \
	else \
		echo "❌ Builds are NOT reproducible!"; \
		if [ "$$BINS_MATCH" = "false" ]; then \
			echo "Binary differences:"; \
			diff ./test-reproducible/checksums1_bins.txt ./test-reproducible/checksums2_bins.txt || true; \
		fi; \
		if [ "$$DEBS_MATCH" = "false" ]; then \
			echo "Package differences:"; \
			diff ./test-reproducible/checksums1_debs.txt ./test-reproducible/checksums2_debs.txt || true; \
		fi; \
		rm -rf ./test-reproducible; \
		exit 1; \
	fi
	@rm -rf ./test-reproducible
	@echo "🎉 Reproducibility test passed"

.PHONY: package-install-local
package-install-local: package-local ## Install locally built package
	@echo "Installing local package..."
	@DEB_FILE=$$(find ./dist -name "*httpserver*.deb" | head -1); \
	if [ -n "$$DEB_FILE" ]; then \
		echo "Installing $$DEB_FILE"; \
		sudo dpkg -i "$$DEB_FILE" || sudo apt-get -f install -y; \
		echo "✅ Package installed successfully"; \
		echo "To start service: sudo systemctl start go-template-httpserver"; \
		echo "To check status: sudo systemctl status go-template-httpserver"; \
	else \
		echo "❌ No .deb file found in ./dist/"; \
		exit 1; \
	fi

.PHONY: package-uninstall
package-uninstall: ## Uninstall locally installed package
	@echo "Uninstalling go-template packages..."
	@if dpkg -l | grep -q go-template-httpserver; then \
		sudo systemctl stop go-template-httpserver || true; \
		sudo dpkg -r go-template-httpserver; \
		echo "✅ HTTP server package removed"; \
	fi
	@if dpkg -l | grep -q go-template-cli; then \
		sudo dpkg -r go-template-cli; \
		echo "✅ CLI package removed"; \
	fi

.PHONY: package-info
package-info: ## Show information about built packages
	@echo "📦 Package Information"
	@echo "====================="
	@for pkg in $$(find dist/ -name "*.deb" 2>/dev/null); do \
		echo "Package: $$pkg"; \
		echo "Size: $$(du -h "$$pkg" | cut -f1)"; \
		echo "Contents:"; \
		dpkg-deb --contents "$$pkg" | head -10; \
		echo "---"; \
	done

.PHONY: package-clean
package-clean: ## Clean packaging artifacts
	@echo "Cleaning packaging artifacts..."
	@rm -rf dist/
	@echo "✅ Packaging artifacts cleaned"

.PHONY: package-release
package-release: ## Create a release (requires git tag)
	@if [ "$$(git describe --exact-match --tags HEAD 2>/dev/null)" = "" ]; then \
		echo "❌ No git tag found. Create a tag first: git tag v1.0.0"; \
		exit 1; \
	fi
	@echo "🚀 Creating release for tag: $$(git describe --tags)"
	@goreleaser release --clean
	@echo "✅ Release created successfully"
