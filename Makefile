#!/usr/bin/make -f

VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')

COMMIT := $(shell git log -1 --format='%H')

BUILD_DIR ?= $(CURDIR)/build
LEDGER_ENABLED ?= true

# ********** Golang configs **********

CMTVERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
export GO111MODULE = on

GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

# ********** process build tags **********

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# ********** process linker flags **********

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=mitosis \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=mitosisd \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(CMTVERSION)

# DB backend selection
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
ifeq (badgerdb,$(findstring badgerdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += badgerdb
endif
# handle rocksdb
ifeq (rocksdb,$(findstring rocksdb,$(COSMOS_BUILD_OPTIONS)))
  CGO_ENABLED=1
  build_tags += rocksdb
endif
# handle boltdb
ifeq (boltdb,$(findstring boltdb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += boltdb
endif

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)

ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags := $(strip $(ldflags))

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# Check for debug option
ifeq (debug,$(findstring debug,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

###############################################################################
###                                  Build                                  ###
###############################################################################

# Detect current OS and architecture
HOST_OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
HOST_ARCH := $(shell uname -m)

# Normalize architecture names
ifeq ($(HOST_ARCH),x86_64)
  HOST_ARCH := amd64
endif
ifeq ($(HOST_ARCH),aarch64)
  HOST_ARCH := arm64
endif

# Convert darwin to match Go's GOOS
ifeq ($(HOST_OS),darwin)
  HOST_GOOS := darwin
else
  HOST_GOOS := $(HOST_OS)
endif

# Default to host OS and architecture
TARGET_OS ?= $(HOST_GOOS)
TARGET_ARCH ?= $(HOST_ARCH)

# Validate target OS and architecture
VALID_OS := darwin linux
VALID_ARCH := amd64 arm64

ifeq ($(filter $(TARGET_OS),$(VALID_OS)),)
  $(error Invalid TARGET_OS: $(TARGET_OS). Valid options: $(VALID_OS))
endif

ifeq ($(filter $(TARGET_ARCH),$(VALID_ARCH)),)
  $(error Invalid TARGET_ARCH: $(TARGET_ARCH). Valid options: $(VALID_ARCH))
endif

# Set cross-compilation environment variables
export GOOS = $(TARGET_OS)
export GOARCH = $(TARGET_ARCH)

# Adjust build output directory for cross-compilation
BUILD_SUFFIX := $(TARGET_OS)-$(TARGET_ARCH)
BUILD_OUTPUT_DIR := $(BUILD_DIR)/$(BUILD_SUFFIX)

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_OUTPUT_DIR)/

# Build targets for specific binaries (current architecture)
build-mitosisd:
	BINARY_NAME=mitosisd $(MAKE) build

build-midevtool:
	BINARY_NAME=midevtool $(MAKE) build

build-mito:
	BINARY_NAME=mito $(MAKE) build

# Build targets for cross-compilation
# Usage: make cross-build-mitosisd darwin arm64
cross-build-mitosisd:
	$(MAKE) build-mitosisd TARGET_OS=$(word 2,$(MAKECMDGOALS)) TARGET_ARCH=$(word 3,$(MAKECMDGOALS))

cross-build-midevtool:
	$(MAKE) build-midevtool TARGET_OS=$(word 2,$(MAKECMDGOALS)) TARGET_ARCH=$(word 3,$(MAKECMDGOALS))

cross-build-mito:
	$(MAKE) build-mito TARGET_OS=$(word 2,$(MAKECMDGOALS)) TARGET_ARCH=$(word 3,$(MAKECMDGOALS))

# Build all binaries for current architecture
build-all:
	$(MAKE) build-mitosisd
	$(MAKE) build-midevtool
	$(MAKE) build-mito

# Build all binaries for all supported architectures
cross-build-all:
	@echo "Building for all supported architectures..."
	$(MAKE) build-all TARGET_OS=darwin TARGET_ARCH=amd64
	$(MAKE) build-all TARGET_OS=darwin TARGET_ARCH=arm64
	$(MAKE) build-all TARGET_OS=linux TARGET_ARCH=amd64
	$(MAKE) build-all TARGET_OS=linux TARGET_ARCH=arm64
	@echo "Cross-compilation complete! Binaries are in $(BUILD_DIR)/"

$(BUILD_TARGETS): go.sum $(BUILD_OUTPUT_DIR)/
	@echo "Building $(BINARY_NAME) for $(TARGET_OS)/$(TARGET_ARCH)..."
	cd ${CURDIR}/cmd/$(BINARY_NAME) && GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILD_OUTPUT_DIR)/:
	mkdir -p $(BUILD_OUTPUT_DIR)/

$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

# Allow make to accept additional arguments as fake targets for cross-compilation
.PHONY: build build-all cross-build-all build-mitosisd build-midevtool build-mito cross-build-mitosisd cross-build-midevtool cross-build-mito install

clean:
	rm -rf $(BUILD_DIR)/ artifacts/

# Show build help
build-help:
	@echo "Build Commands:"
	@echo "  Current architecture builds:"
	@echo "    make build-mitosisd     - Build mitosisd for current platform ($(HOST_GOOS)/$(HOST_ARCH))"
	@echo "    make build-midevtool    - Build midevtool for current platform"
	@echo "    make build-mito         - Build mito for current platform"
	@echo "    make build-all          - Build all binaries for current platform"
	@echo ""
	@echo "  Cross-compilation builds:"
	@echo "    make cross-build-mitosisd darwin amd64    - Build mitosisd for macOS x64"
	@echo "    make cross-build-mitosisd darwin arm64    - Build mitosisd for macOS ARM64"
	@echo "    make cross-build-mitosisd linux amd64     - Build mitosisd for Linux x64"
	@echo "    make cross-build-mitosisd linux arm64     - Build mitosisd for Linux ARM64"
	@echo "    (same pattern for cross-build-midevtool and cross-build-mito)"
	@echo ""
	@echo "  Bulk builds:"
	@echo "    make cross-build-all    - Build all binaries for all supported architectures"
	@echo ""
	@echo "  Alternative using environment variables:"
	@echo "    TARGET_OS=darwin TARGET_ARCH=arm64 make build-mitosisd"
	@echo ""
	@echo "  Output directory: $(BUILD_DIR)/{os}-{arch}/"
	@echo "  Supported OS: darwin, linux"
	@echo "  Supported ARCH: amd64, arm64"

.PHONY: clean build-help

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=$$(go env GOPATH)/bin/golangci-lint
golangci_version=v1.64.8

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m

lint-fix:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --fix --out-format=tab --issues-exit-code=0

format:
	@go install mvdan.cc/gofumpt@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./tests/mocks/*" -not -path "./bindings/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go" | xargs gofumpt -w -l
	@$(golangci_lint_cmd) run --fix
.PHONY: format

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################
PACKAGES_UNIT=$(shell go list ./... | grep -v -e '/tests/e2e')
PACKAGES_E2E=$(shell cd tests/e2e && go list ./... | grep '/e2e')
TEST_PACKAGES=./...
TEST_TARGETS := test-unit test-e2e

test-unit: ARGS=-timeout=5m -tags='norace'
test-unit: TEST_PACKAGES=$(PACKAGES_UNIT)
test-e2e: ARGS=-timeout=25m -v
test-e2e: TEST_PACKAGES=$(PACKAGES_E2E)
$(TEST_TARGETS): test

test:
ifneq (,$(shell which tparse 2>/dev/null))
	@echo "--> Running tests"
	@go test -mod=readonly -json $(ARGS) $(TEST_PACKAGES) | tparse
else
	@echo "--> Running tests"
	@go test -mod=readonly $(ARGS) $(TEST_PACKAGES)
endif

###############################################################################
###                           Tests & Simulation                            ###
###############################################################################
bufgen: ## Generates protobufs using buf generate.
	@./scripts/protocgen.sh

###############################################################################
###                                Localnet                                 ###
###############################################################################

MITOSISD_HOME = $(CURDIR)/tmp/localnet/mitosisd
MITOSISD_CHAIN_ID = 'mitosis-localnet-1'
MITOSISD_INFRA_DIR = $(CURDIR)/infra/localnet/mitosisd
EC_INFRA_DIR = $(CURDIR)/infra/localnet/ec
GETH_DATA_DIR = $(CURDIR)/tmp/localnet/geth
RETH_DATA_DIR = $(CURDIR)/tmp/localnet/reth
GOV_ENTRYPOINT = '0x06c9918ff483fd88C65dD02E788427cfF04545b9'

clean-geth:
	rm -rf $(GETH_DATA_DIR)

setup-geth: clean-geth
	docker run --rm \
		-v $(EC_INFRA_DIR):/infra \
		-v $(GETH_DATA_DIR):/data \
		ethereum/client-go:v1.15.11 init \
			--datadir /data \
			--db.engine pebble \
			--state.scheme=hash \
			/infra/genesis.json

run-geth:
	docker run --rm \
		-p 30303:30303 \
		-p 8545:8545 \
		-p 8551:8551 \
		-v $(EC_INFRA_DIR):/infra \
		-v $(GETH_DATA_DIR):/data \
		ethereum/client-go:v1.15.11 \
			--datadir /data \
			--http \
			--http.addr 0.0.0.0 \
			--http.vhosts "*" \
			--http.api eth,net,web3,txpool,rpc,debug \
			--authrpc.addr 0.0.0.0 \
			--authrpc.vhosts "*" \
			--authrpc.jwtsecret /infra/jwt.hex \
			--db.engine pebble \
			--state.scheme=hash \
			--syncmode full \
			--gcmode archive \
			--miner.recommit=500ms

clean-reth:
	rm -rf $(RETH_DATA_DIR)

setup-reth: clean-reth
	docker run --rm \
		-v $(EC_INFRA_DIR):/infra \
		-v $(RETH_DATA_DIR):/data \
		ghcr.io/paradigmxyz/reth:v1.3.12 init \
			--datadir /data \
			--chain /infra/genesis.json

run-reth:
	docker run --rm \
		-p 30303:30303 \
		-p 30303:30303/udp \
		-p 8545:8545 \
		-p 8551:8551 \
		-p 9001:9001 \
		-v $(EC_INFRA_DIR):/infra \
		-v $(RETH_DATA_DIR):/data \
		ghcr.io/paradigmxyz/reth:v1.3.12 node \
			--datadir /data \
			--chain /infra/genesis.json \
			--http \
			--http.addr 0.0.0.0 \
			--http.api eth,net,web3,txpool,rpc,debug,trace \
			--authrpc.addr 0.0.0.0 \
			--authrpc.jwtsecret /infra/jwt.hex \
			--metrics 0.0.0.0:9001 \
			--builder.interval 30ms \
			--builder.deadline 1

clean-mitosisd:
	rm -rf $(MITOSISD_HOME)

setup-mitosisd: build clean-mitosisd
	MITOSISD=./build/mitosisd \
		MITOSISD_HOME=$(MITOSISD_HOME) \
		MITOSISD_CHAIN_ID=$(MITOSISD_CHAIN_ID) \
		EC_JWT_FILE=$(EC_INFRA_DIR)/jwt.hex \
		GOV_ENTRYPOINT=$(GOV_ENTRYPOINT) \
		$(MITOSISD_INFRA_DIR)/setup.sh

run-mitosisd:
	./build/mitosisd start \
		--home $(MITOSISD_HOME) \
		--p2p.laddr=tcp://0.0.0.0:26656 \
		--rpc.laddr=tcp://0.0.0.0:26657 \
		--grpc.enable \
		--grpc.address=0.0.0.0:9090 \
		--api.enable \
		--api.address=tcp://0.0.0.0:1317 \
		--api.enabled-unsafe-cors \
		--log_level "abci-wrapper:debug,x/evmengine:debug,x/evmvalidator:debug,x/evmgov:debug,*:info"

###############################################################################
###                                  Devnet                                 ###
###############################################################################

devnet-clean:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile '*' down
	rm -rf ./tmp/devnet

devnet-build:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile '*' build

devnet-init:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile init up -d

devnet-up:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile validator --profile node up -d

devnet-down:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile validator --profile node down

devnet-create-validator:
	docker compose --project-directory ./ -f ./infra/devnet/docker-compose.devnet.yml -p mitosis-devnet \
		--profile create-validator up -d
