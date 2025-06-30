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


## For Users

### Installation

```bash
# Install latest release
curl -sSL https://raw.githubusercontent.com/mitosis-org/chain/main/scripts/install.sh | bash

# Or download from releases
wget https://github.com/mitosis-org/chain/releases/latest/download/mitosisd-linux-amd64
chmod +x mitosisd-linux-amd64 && mv mitosisd-linux-amd64 /usr/local/bin/mitosisd
```

### Running a Node

```bash
# Initialize node
mitosisd init my-node --chain-id mitosis-localnet-1

# Start node
mitosisd start
```

## For Developers

### Prerequisites

- Go 1.24+
- Git

### Building from Source

First, clone the repository:

```bash
git clone https://github.com/mitosis-org/chain
cd chain
```

Next, build the binary:

```bash
make build
```

### Running Tests

```bash
make test
```

### Development Environment

For detailed development environment setup including localnet and devnet configurations, please see our [Contributing Guide](CONTRIBUTING.md#-environment-setup).

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
