# mitosis

[![License](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![codecov](https://codecov.io/gh/mitosis-org/chain/graph/badge.svg?token=4Mkp1Ipjc3)](https://codecov.io/gh/mitosis-org/chain)
[![Security](https://github.com/mitosis-org/chain/actions/workflows/security.yml/badge.svg?branch=main)](https://github.com/mitosis-org/chain/actions/workflows/security.yml)
[![Quality Gate](https://github.com/mitosis-org/chain/actions/workflows/quality-gate.yml/badge.svg?branch=main)](https://github.com/mitosis-org/chain/actions/workflows/quality-gate.yml)

**Next-generation DeFi network enabling programmable liquidity across multiple protocols**

## What is Mitosis?

Mitosis is a Network for Programmable Liquidity that leads the next generation of DeFi through the tokenization of liquidity positions while ensuring seamless integration into the Mitosis Ecosystem. It transforms traditional DeFi positions into tokenized, programmable assets that can be seamlessly leveraged across multiple protocols. Mitosis allows liquidity providers to deposit once and seamlessly earn rewards across multiple protocols in the ecosystem, where your liquidity works harder and smarter.

## Goals

As a full Ethereum-compatible blockchain, Mitosis allows users to connect to the network and interact with smart contracts using familiar Ethereum tooling. Building a successful blockchain requires creating a high-quality implementation that is both secure and efficient, as well as being easy to use. It also requires building a strong community of contributors who can help support and improve the software.

More concretely, our goals are:

1. **Modularity**: Every component of Mitosis is built to be used as a library: well-tested, heavily documented and benchmarked. We envision that developers will import the node's crates, mix and match, and innovate on top of them.

2. **Performance**: Mitosis aims to be fast, leveraging the proven Cosmos SDK architecture with full EVM compatibility. We optimize for DeFi and cross-chain operations with minimal latency.

3. **Free for anyone to use any way they want**: Mitosis is free open source software, built for the community, by the community. By licensing the software under the GNU General Public License v3.0 (GPLv3) plus the Open Interoperability Requirement (OIR), we want developers to use it freely and promote interoperability.

4. **EVM Compatibility**: Full compatibility with Ethereum tooling, wallets, and infrastructure. Run your Ethereum dApps without changes while leveraging advanced Cosmos features.

5. **Developer Experience**: Unified development experience using Solidity for all logic, standard Ethereum tooling, and familiar wallet integration.

## How to run a chain in testing environment

We categorize testing development environments into:

- **Localnet** - For fast development and testing iterations in local environment. Runs a single validator for the mitosis chain.
- **Devnet** - For development and testing with complete form of components. Runs two validator nodes and a non-validator node for the mitosis chain.

### Chain IDs

- **Localnet**
  - Chain ID (EVM): `124899`
  - Chain ID (Cosmos SDK): `mitosis-localnet-1`
- **Devnet**
  - Chain ID (EVM): `124864`
  - Chain ID (Cosmos SDK): `mitosis-devnet-1`

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

### Contributing

If you want to contribute, or follow along with contributor discussion, you can use our GitHub discussions and issues.

- Our contributor guidelines can be found in [CONTRIBUTING.md](CONTRIBUTING.md).
- See our [Security Policy](SECURITY.md) for security-related contributions.

## Getting Help

If you have any questions:

- Open a discussion with your question, or
- Open an issue with the bug
- Check our documentation at [docs.mitosis.org](https://docs.mitosis.org)

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
