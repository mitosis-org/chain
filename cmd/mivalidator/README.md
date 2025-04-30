# Mitosis Validator CLI Tool

`mivalidator` is a command-line tool for managing validators in the Mitosis network. It allows validator operators to create validators, manage collateral, and update validator settings.

## Installation

To install the `mivalidator` tool, clone the repository and build the binary:

```bash
git clone https://github.com/mitosis-org/chain
cd chain
make build
```

The binary will be available at `build/mivalidator`.

## Configuration

All commands require a connection to an Ethereum RPC endpoint. You can specify this using the `--rpc-url` flag.

### RPC URL Configuration

The `--rpc-url` flag is required for all commands and specifies the Ethereum JSON-RPC endpoint to connect to. The tool uses this endpoint to interact with the blockchain.

### Authentication

The tool supports signing transactions using:

- Private key (using the `--private-key` flag)

Example:
```bash
mivalidator update-operator --validator 0x123... --rpc-url https://rpc.example.com --contract <contract-address> --private-key <your-private-key> --operator <new-operator-address>
```

> ⚠️ **Security Warning**: Providing private keys on the command line is not secure as they may be stored in your shell history. Consider using environment variables or other secure methods to provide private keys.

> Note: Read-only commands like `validator-info` do not require a private key.

## Available Commands

### Information Commands

#### validator-info

Retrieve detailed information about a validator from the ValidatorManager contract.

**Required flags:** `--validator`, `--rpc-url`, `--contract`

```bash
mivalidator validator-info --validator <validator-address> --rpc-url <rpc-url> --contract <contract-address> [--epoch <epoch-number>]
```

### Validator Creation

#### create-validator

Register a new validator in the ValidatorManager contract.

**Required flags:** `--pubkey`, `--operator`, `--reward-manager`, `--withdrawal-recipient`, `--initial-collateral`, `--commission-rate`, `--metadata`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator create-validator \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --pubkey <public-key-hex> \
  --operator <operator-address> \
  --reward-manager <reward-manager-address> \
  --withdrawal-recipient <withdrawal-recipient-address> \
  --commission-rate <percentage> \
  --initial-collateral <amount-in-MITO> \
  --metadata <json-string>
```

### Update Commands

#### update-operator

Update the operator address for an existing validator.

**Required flags:** `--validator`, `--operator`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator update-operator \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --operator <new-operator-address>
```

#### update-withdrawal-recipient

Update the withdrawal recipient for an existing validator. The withdrawal recipient is the address that receives the withdrawal of any ETH and validator rewards when a validator exits the system.

**Required flags:** `--validator`, `--recipient`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator update-withdrawal-recipient \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --recipient <new-recipient-address>
```

#### update-reward-manager

Update the reward manager address for an existing validator.

**Required flags:** `--validator`, `--reward-manager`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator update-reward-manager \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --reward-manager <new-reward-manager-address>
```

#### update-reward-config

Update the reward configuration for an existing validator. Currently, this allows updating the commission rate.

**Required flags:** `--validator`, `--commission-rate`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator update-reward-config \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --commission-rate <new-percentage>
```

#### update-metadata

Update the metadata for an existing validator.

**Required flags:** `--validator`, `--metadata`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator update-metadata \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --metadata <json-string>
```

### Collateral Management

#### deposit-collateral

Add more collateral to an existing validator.

**Required flags:** `--validator`, `--amount`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator deposit-collateral \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --amount <amount-in-MITO>
```

#### withdraw-collateral

Withdraw collateral from an existing validator. The withdrawn amount will be sent to the validator's withdrawal recipient address after a delay period.

**Required flags:** `--validator`, `--amount`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator withdraw-collateral \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address> \
  --amount <amount-in-MITO>
```

### Unjailing

#### unjail-validator

Unjail a validator that has been jailed due to downtime or other violations.

**Required flags:** `--validator`, `--rpc-url`, `--contract`, `--private-key`

```bash
mivalidator unjail-validator \
  --rpc-url <rpc-url> \
  --contract <contract-address> \
  --private-key <your-private-key> \
  --validator <validator-address>
```

## Common Options

All commands support the following flags:

```
--rpc-url         RPC URL for Ethereum client (default "http://localhost:8545")
--contract        ValidatorManager contract address (required)
```

For transaction-signing operations:

```
--private-key     Private key for signing transactions (required for state-changing operations)
--yes             Skip confirmation prompts (use with caution)
--nonce           Manually specify nonce for transaction (optional)
```
