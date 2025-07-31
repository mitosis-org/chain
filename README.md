# Mitosis Chain

[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Security](https://github.com/mitosis-org/chain/actions/workflows/security.yml/badge.svg?branch=main)](https://github.com/mitosis-org/chain/actions/workflows/security.yml)
[![Quality Gate](https://github.com/mitosis-org/chain/actions/workflows/quality-gate.yml/badge.svg?branch=main)](https://github.com/mitosis-org/chain/actions/workflows/quality-gate.yml)

![Mitosis Chain Banner](assets/banner.png)

## What is Mitosis Chain?

Mitosis is a Network for Programmable Liquidity that leads the next generation of DeFi through the tokenization of liquidity positions while ensuring seamless integration into the Mitosis Ecosystem. It transforms traditional DeFi positions into tokenized, programmable assets that can be seamlessly leveraged across multiple protocols.

**Mitosis Chain** is the backbone blockchain infrastructure that powers the [Mitosis Protocol](https://github.com/mitosis-org/protocol) and Ecosystem. Built with a modular architecture separating consensus and execution layers, it delivers 100% EVM compatibility with instant finality through CometBFT and Cosmos SDK.

As a full Ethereum-compatible blockchain, Mitosis Chain allows users to connect to the network and interact with smart contracts using familiar Ethereum tooling. Building a successful blockchain requires creating a high-quality implementation that is both secure and efficient, as well as being easy to use. It also requires building a strong community of contributors who can help support and improve the software.

## How to run a chain in local environment

We categorize testing development environments into:

- **Localnet** - For fast development and testing iterations in local environment. Runs a single validator for the mitosis chain.
- **Devnet** - For development and testing with complete form of components. Runs two validator nodes and a non-validator node for the mitosis chain.

## Chain IDs

| Environment | EVM Chain ID | Cosmos SDK Chain ID    |
|-------------|--------------|------------------------|
| Localnet    | 124899       | mitosis-localnet-1    |
| Devnet      | 124864       | mitosis-devnet-1      |

### Localnet Setup

Localnet requires running both an execution client (`geth` or `reth`) and a consensus client (`mitosisd`).

**Prerequisites**
```bash
# Ensure submodules are fetched
git submodule update --init --recursive
```

**Setup and Run Execution Client**
```bash
# Initialize geth (removes old data if exists)
make setup-geth
# Alternative: make setup-reth

# Run geth
make run-geth
# Alternative: make run-reth
```

**Setup and Run Consensus Client**
```bash
# Initialize mitosisd (removes old data if exists)
make setup-mitosisd

# Run mitosisd
make run-mitosisd
```

**Deploy and Setup Consensus Entrypoint Contracts**
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

**Cleanup**
```bash
# Clean both clients (must clean both together)
make clean-geth    # or make clean-reth
make clean-mitosisd
```

### Devnet Setup

Devnet provides a more complete testing environment with multiple nodes.

**Build Docker Image**
```bash
make devnet-build
```

**Initialize and Start Devnet**
```bash
# Initialize the mitosis chain
make devnet-init

# Start all nodes
make devnet-up

# Verify nodes are running
docker logs mitosis-devnet-node-mitosisd-1
docker logs mitosis-devnet-node-reth-1

# Test RPC connectivity
cast block-number --rpc-url http://localhost:18545
# Or use curl:
curl -H "Content-Type: application/json" --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":124864}' http://localhost:18545
```

**Deploy Consensus Entrypoint Contracts (Required for Testing)**
```bash
# Run this in https://github.com/mitosis-org/protocol
# The deployed address would be:
# - ConsensusGovernanceEntrypoint:  0x06c9918ff483fd88C65dD02E788427cfF04545b9
# - ConsensusValidatorEntrypoint :  0x9866D79EF3e9c0c22Db2b55877013e13a60AD478
RPC_URL="http://127.0.0.1:18545" ./tools/deploy-consensus-entrypoints.sh
```

**Create a validator for testing**
```bash
make devnet-create-validator
```

**Devnet Management**
```bash
# Stop nodes (keeps data)
make devnet-down

# Complete cleanup (removes all data)
make devnet-clean
```

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

## Contributing

If you want to contribute, or follow along with contributor discussion, you can use our GitHub discussions and issues.

- Our contributor guidelines can be found in [CONTRIBUTING.md](CONTRIBUTING.md).
- See our [Security Policy](SECURITY.md) for security-related contributions.

## Getting Help

If you have any questions:

- Open a discussion with your question, or
- Open an issue with the bug
- Check our [developer docs](https://docs.mitosis.org/developers/)

## Security

See [SECURITY.md](SECURITY.md).

## Acknowledgements

Mitosis is built on the shoulders of giants. We would like to thank:

- **Cosmos SDK**: For providing the robust blockchain framework that powers our consensus layer
- **CometBFT**: For the battle-tested Byzantine fault-tolerant consensus mechanism
- **Octane**: For the EVM-Engine API integration that enables our modular architecture
- **Ethereum**: For the EVM specification and the vibrant ecosystem that inspired this project

## License

This project is licensed under the GNU General Public License v3.0 (GPLv3) **plus the Open Interoperability Requirement (OIR)**.

See the [LICENSE](LICENSE) file for details.
