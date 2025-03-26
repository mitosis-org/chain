# Mitosis Chain

## Architecture

### Overview

![architecture.png](assets/architecture.png)

The Mitosis Chain employs a modular architecture that separates execution from consensus.

- The execution layer is fully EVM-compatible, enabling unmodified Ethereum execution clients to process transactions, manage state, and execute smart contracts.
- The consensus layer is built upon the Cosmos SDK and utilizes CometBFT for consensus.
- The two layers communicate with each other using [Engine API](https://hackmd.io/@danielrachi/engine_api). The consensus layer utilizes [Octane](https://github.com/omni-network/omni/tree/main/octane) for Engine API implementation.

Most of our logic exists on the execution layer, while the consensus layer is kept thin by having only minimal code and responsibilities for consensus.

### Validator & Governance System

In most Cosmos SDK-based chains, validator and governance systems are built using `x/staking`, `x/slashing`, `x/distribution`, and `x/gov` modules provided by Cosmos SDK.
However, we don't use all of them. We implement most of the necessary logic as smart contracts in the EVM (execution layer). \
For example:
- A user stakes and delegates $MITO to a validator on EVM.
- An operator creates and operates a validator on EVM.
- Staking rewards are distributed on EVM.
- For governance, users cast votes and proposals are executed on EVM.

The contracts manage all user flows and serve as the source of truth.
Some information from the contracts is delivered to the consensus layer through EVM logs on `ConsensusValidatorEntrypoint` and `ConsensusGovernanceEntrypoint` contracts.
When an EVM block has been created and finalized on the consensus layer, the consensus layer parses and processes the EVM logs in the block.

### Core Modules

#### `x/evmengine` (forked from Octane)

This module communicates with an execution client through Engine API and wraps an EVM block into one transaction in the consensus layer.
Note that there can only be one transaction wrapped from an EVM block, and other types of transactions are prohibited in a consensus block. \
This module is forked from [Octane](https://github.com/omni-network/omni/tree/main/octane). There are limited changes from the upstream for seamless integration with `x/evmvalidator` and `x/evmgov`.

#### `x/evmvalidator`

This module manages a validator set on the consensus layer. Note that there are no features such as delegation and reward distribution because those features are implemented in contracts on EVM.
Most states are managed in the EVM contracts, and this module simply applies validator set changes and consensus voting power updates delivered from `ConsensusValidatorEntrypoint`. \
We could say this module is a lightweight version of `x/staking` that has EVM contracts as the source of truth.
This module also implements some parts of interfaces of `x/staking` to integrate with `x/slashing` and `x/evidence`.

#### `x/evmgov`

This module provides arbitrary message execution from EVM. Governance on EVM can trigger arbitrary message execution against consensus layer modules. It can be used for cases such as module parameter changes.
It is very lightweight compared to `x/gov` because there are no concepts such as proposals and voting power.
These concepts are implemented in EVM contracts, and this module simply executes arbitrary messages delivered from `ConsensusGovernanceEntrypoint`.

## Setup

We categorize environments for setup into:
- **Localnet**
  - Localnet is for fast development and testing iterations in local environment.
  - It runs a single validator for the mitosis chain.
- **Devnet**
  - Devnet is for development and testing with complete form of the components.
  - It runs two validator nodes and a non-validator node for the mitosis chain.

**Chain IDs**
- **Localnet**
  - Chain ID (EVM): `124899`
  - Chain ID (Cosmos SDK): `mitosis-localnet-1`
- **Devnet**
  - Chain ID (EVM): `124864`
  - Chain ID (Cosmos SDK): `mitosis-devnet-1`

### Localnet

You should run an execution client (`geth` or `reth`) and an consensus client (`mitosisd`) both.


Pre-requisites

Make sure you have fetched the GitHub submodules:

```sh
git submodule update --init --recursive
```

Setup and run an execution client (`geth` or `reth`):

```bash
# It initializes geth. If there was already initialized, it remove all old data and re-initialize it.
make setup-geth # or `make setup-reth`

# Note that it just tries to use existing data instead of setting up geth automatically.
# You should run `setup-geth` if you haven't initialized geth yet or want to reset it.
make run-geth # or `make setup-reth`
```

Setup and run an consensus client (`mitosisd`):
```bash
# It initializes mitosisd. If there was already initialized, it remove all old data and re-initialize it.
make setup-mitosisd

# Note that it just tries to use existing data instead of setting up mitosisd automatically.
# You should run `setup-mitosisd` if you haven't initialized mitosisd yet or want to reset it.
make run-mitosisd
```

Remove all data (reset) of `geth` and `mitosisd`:
```bash
# Note that it won't be working as expected if you clean up only one of geth and mitosisd.
make clean-geth # or `make clean-reth`
make clean-mitosisd
```

Deploy & Setup consensus entrypoint contracts:
```bash
# Run this command in https://github.com/mitosis-org/protocol to deploy the consensus entrypoint contracts.
# The deployed address would be:
# - ConsensusGovernanceEntrypoint:  0x06c9918ff483fd88C65dD02E788427cfF04545b9
# - ConsensusValidatorEntrypoint :  0x9866D79EF3e9c0c22Db2b55877013e13a60AD478
./tools/deploy-consensus-entrypoints.sh

# Note that `ConsensusGovernanceEntrypoint` address is managed in `app.toml`:
#   [evmgov]
#   entrypoint = "0x06c9918ff483fd88C65dD02E788427cfF04545b9"

# Update `ConsensusValidatorEntrypoint` address in x/evmvalidator module.
./build/midevtool governance execute \
  --entrypoint 0x06c9918ff483fd88C65dD02E788427cfF04545b9 \
  --private-key 0x5a496832ac0d7a484e6996301a5511dbc3b723d037bc61261ecaf425bd6a5b37 \
  --msg '[{"@type":"/mitosis.evmvalidator.v1.MsgUpdateValidatorEntrypointContractAddr","authority":"mito1g86pactsvfrcglkvqzvdwkxhjshafu280q95p7","addr":"0x9866D79EF3e9c0c22Db2b55877013e13a60AD478"}]'
```

### Devnet

#### Build a dockerfile for mitosisd
```bash
make devnet-build
```

#### Setup and run the mitosis chain
```bash
# Init a mitosis chain.
# It prepares a genesis file for the mitosis chain.
make devnet-init

# Run nodes for a mitosis chain.
make devnet-up

# Check the status of the nodes.
docker logs mitosis-devnet-node-mitosisd-1
docker logs mitosis-devnet-node-reth-1

# Check geth rpc is working properly.
cast block-number --rpc-url http://localhost:18545
# You can use `curl` instead if `cast` is not installed.
curl -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":124864}' http://localhost:18545
```

#### Deploy consensus entrypoint contracts:
```bash
# Run this command in https://github.com/mitosis-org/protocol to deploy the consensus entrypoint contracts.
# The deployed address would be:
# - ConsensusGovernanceEntrypoint:  0x06c9918ff483fd88C65dD02E788427cfF04545b9
# - ConsensusValidatorEntrypoint :  0x9866D79EF3e9c0c22Db2b55877013e13a60AD478
RPC_URL="http://127.0.0.1:18545" ./tools/deploy-consensus-entrypoints.sh
```

#### Create a validator for the mitosis chain
```bash
# It creates a validator for the `subval` node.
# Note that it setup `ConsensusValidatorEntrypoint` address automatically.
make devnet-create-validator
```

#### Stop the mitosis chain and reset the data
```bash
# It stops the nodes.
# You can start them again with keeping existing data through `make devnet-up`.
make devnet-down

# It stops all nodes and removes all data of the nodes.
# It also removes the initialization data which is created by `make devnet-init`.
make devnet-clean
```
