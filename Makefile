#!/usr/bin/make -f

#VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION := '0.0.1'

COMMIT := $(shell git log -1 --format='%H')

BUILD_DIR ?= $(CURDIR)/build
MITOSISD_HOME = $(CURDIR)/tmp/localnet/mitosisd
LEDGER_ENABLED ?= true

CHAIN_ID = 'mitosis-localnet-1'

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

$(BUILD_TARGETS): go.sum $(BUILD_DIR)/
	cd ${CURDIR}/cmd/$(BINARY_NAME) && go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILD_DIR)/:
	mkdir -p $(BUILD_DIR)/

.PHONY: build

clean:
	rm -rf $(BUILD_DIR)/ artifacts/

.PHONY: clean

###############################################################################
###                                  Docker                                 ###
###############################################################################

docker-mitosisd-build:
	@echo "Building docker image for mitosisd"
	docker build -f mitosisd.Dockerfile -t mitosisd .

###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
golangci_version=v1.52.2

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
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" -not -path "./client/docs/statik/statik.go" -not -path "./tests/mocks/*" -not -name "*.pb.go" -not -name "*.pb.gw.go" -not -name "*.pulsar.go" -not -path "./crypto/keys/secp256k1/*" | xargs gofumpt -w -l
	$(golangci_lint_cmd) run --fix
.PHONY: format

###############################################################################
###                                Localnet                                 ###
###############################################################################

ETH_INFRA_DIR = ./infra/localnet/geth
ETH_DATA_DIR = ./tmp/localnet/geth

clean-geth:
	rm -rf $(MITOSISD_HOME)

setup-geth:
	rm -rf $(ETH_DATA_DIR)
	docker run --rm \
		-v $(ETH_INFRA_DIR):/infra \
		-v $(ETH_DATA_DIR):/.data \
		ethereum/client-go init \
		--datadir /.data \
		/infra/genesis.json

run-geth:
	docker run --rm \
		-p 30303:30303 \
		-p 8545:8545 \
		-p 8551:8551 \
		-v $(ETH_INFRA_DIR):/infra \
		-v $(ETH_DATA_DIR):/.data \
		ethereum/client-go \
		--http \
		--http.addr 0.0.0.0 \
		--http.api eth,net \
		--authrpc.addr 0.0.0.0 \
		--authrpc.jwtsecret /infra/jwt.hex \
		--authrpc.vhosts "*" \
		--datadir /.data \
		--miner.recommit=500ms

setup-localnet: build
	rm -rf $(MITOSISD_HOME)

	./build/mitosisd init localnet --chain-id $(CHAIN_ID) --default-denom thai --home $(MITOSISD_HOME)
	./build/mitosisd config set client chain-id $(CHAIN_ID) --home $(MITOSISD_HOME)
	./build/mitosisd config set client keyring-backend test --home $(MITOSISD_HOME)
	./build/mitosisd keys add olivia --keyring-backend test --home $(MITOSISD_HOME)
	./build/mitosisd genesis add-genesis-account olivia 10000000000000000000000000thai --keyring-backend test --home $(MITOSISD_HOME)

	cp ./infra/localnet/mitosisd/priv_validator_key.json $(MITOSISD_HOME)/config/
	jq --argfile ccv ./infra/localnet/mitosisd/ccv-state.json '.app_state.ccvconsumer = $$ccv' $(MITOSISD_HOME)/config/genesis.json > $(MITOSISD_HOME)/config/genesis.json.tmp && cp $(MITOSISD_HOME)/config/genesis.json.tmp $(MITOSISD_HOME)/config/genesis.json
	jq --arg hash `cast block --rpc-url http://127.0.0.1:8545 | grep hash | awk '{print $$2}' | xxd -r -p | base64` \
	 '.app_state.evmengine.execution_block_hash = $$hash' $(MITOSISD_HOME)/config/genesis.json > $(MITOSISD_HOME)/config/genesis.json.tmp && cp $(MITOSISD_HOME)/config/genesis.json.tmp $(MITOSISD_HOME)/config/genesis.json
	jq '.consensus.params.block.max_bytes = "-1"' $(MITOSISD_HOME)/config/genesis.json > $(MITOSISD_HOME)/config/genesis.json.tmp && cp $(MITOSISD_HOME)/config/genesis.json.tmp $(MITOSISD_HOME)/config/genesis.json

	sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.025thai"/' $(MITOSISD_HOME)/config/app.toml
	#sed -i.bak'' 's/validator-mode = false/validator-mode = true/' $(MITOSISD_HOME)/config/app.toml
	sed -i.bak'' 's/mock = true/mock = false/' $(MITOSISD_HOME)/config/app.toml
	sed -i.bak'' 's@endpoint = ""@endpoint = "http://127.0.0.1:8551"@' $(MITOSISD_HOME)/config/app.toml
	sed -i.bak'' 's@jwt-file = ""@jwt-file = "'$(ETH_INFRA_DIR)'/jwt.hex"@' $(MITOSISD_HOME)/config/app.toml

	sed -i.bak'' 's/type = "flood"/type = "nop"/' $(MITOSISD_HOME)/config/config.toml

clean-localnet:
	rm -rf $(MITOSISD_HOME)

run-localnet:
	./build/mitosisd start \
		--home $(MITOSISD_HOME) \
		--api.enabled-unsafe-cors \
		--api.enable \
		--log_level "debug"


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
