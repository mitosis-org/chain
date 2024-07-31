# Mitosis

Mitosis is an Ecosystem-Owned Liquidity (EOL) layer1 blockchain that facilitates newly created modular blockchains to capture TVL and attract users through the Mitosis governance process.

## Setup

We categorize environments for setup into:
- **Localnet**
  - Localnet is for fast development and testing iterations in local environment.
  - It doesn't run full dependencies such as ethos chain, but mitosis chain only. So you can't test integration with Ethos.
  - It runs a single mitosis chain validator, which consists of `geth` and `mitosisd`.
- **Devnet**
  - Devnet is for development and testing with complete form of the components.
  - It runs not only a mitosis chain, but also external components: an ethos chain, an eigenlayer, and a IBC relayer.

**Chain IDs**
- **Localnet**
  - Chain ID (EVM): `25560`
  - Chain ID (Cosmos SDK): `mitosis-localnet-1`
- **Devnet**
  - Chain ID (EVM): `25559`
  - Chain ID (Cosmos SDK): `mitosis-devnet-1`

### Localnet

You should run `geth` and `mitosisd` both.


Pre-requisites

Make sure you have fetched the GitHub submodules:

```sh
git submodule update --init --recursive
```

Setup and run `geth`:

```bash
# It initializes geth. If there was already initialized, it remove all old data and re-initialize it.
make setup-geth

# Note that it just tries to use existing data instead of setting up geth automatically.
# You should run `setup-geth` if you haven't initialized geth yet or want to reset it.
make run-geth
```

Setup and run `mitosisd`:
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
make clean-geth
make clean-mitosisd
```

### Devnet

See [Devnet](infra/devnet/README.md)
