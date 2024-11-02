# Mitosis

Mitosis is an Ecosystem-Owned Liquidity (EOL) layer1 blockchain that facilitates newly created modular blockchains to capture TVL and attract users through the Mitosis governance process.

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

#### Create a validator for the mitosis chain
```bash
# It creates a validator for the `subval` node.
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
