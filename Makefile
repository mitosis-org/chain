#!/usr/bin/make -f

#VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
VERSION := '0.0.1'

COMMIT := $(shell git log -1 --format='%H')

BUILD_DIR ?= $(CURDIR)/build
MITOSISD_HOME = $(CURDIR)/tmp/localnet
LEDGER_ENABLED ?= true

CHAIN_ID = 'mitosis-devnet-1'

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

setup-localnet: build
	@rm -rf $(MITOSISD_HOME)
	@./build/mitosisd init localnet --chain-id $(CHAIN_ID) --default-denom thai --home $(MITOSISD_HOME)
	@./build/mitosisd config set client chain-id $(CHAIN_ID) --home $(MITOSISD_HOME)
	@./build/mitosisd config set client keyring-backend test --home $(MITOSISD_HOME)
#	@./build/mitosisd keys add validator --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd keys add alice --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd keys add bob --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd genesis add-genesis-account validator 10000000000000000000000000stake --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd genesis add-genesis-account alice 1000000000000000000stake --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd genesis add-genesis-account bob 1000000000000000000stake --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd genesis gentx validator 1000000000stake --chain-id $(CHAIN_ID) --home $(MITOSISD_HOME) --keyring-backend test
#	@./build/mitosisd genesis collect-gentxs --home $(MITOSISD_HOME)
	@sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.025thai"/' $(MITOSISD_HOME)/config/app.toml

run-localnet: setup-localnet
	@(./build/mitosisd start \
		--consensus.create_empty_blocks "false" \
		--p2p.pex "false" \
		--home $(MITOSISD_HOME) \
		--api.enabled-unsafe-cors \
		--api.enable \
		--log_level "debug")


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
