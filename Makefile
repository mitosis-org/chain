#!/usr/bin/make -f

#VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION := '0.0.1'

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

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILD_DIR)/

build-mitosisd:
	BINARY_NAME=mitosisd $(MAKE) build

build-midevtool:
	BINARY_NAME=midevtool $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	cd ${CURDIR}/cmd/$(BINARY_NAME) && go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

.PHONY: build

clean:
	rm -rf $(BUILD_DIR)/ artifacts/

.PHONY: clean

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
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
	$(golangci_lint_cmd) run --fix
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
$(TEST_TARGETS): run-tests

run-tests:
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
		ethereum/client-go:v1.15.5 init \
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
		ethereum/client-go:v1.15.5 \
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
		ghcr.io/paradigmxyz/reth:v1.3.4 init \
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
		ghcr.io/paradigmxyz/reth:v1.3.1 node \
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
		--log_level "info"

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
