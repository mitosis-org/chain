# Mito - Mitosis Chain Utilities

Mito is a command-line tool for interacting with the Mitosis blockchain, providing utilities for validator management, transaction handling, and EVM operations.

## Features

- **Transaction Management**: Create and send transactions with support for both signed and unsigned modes
- **Validator Operations**: Create, update, and manage validators
- **Collateral Management**: Deposit and withdraw collateral for validators
- **Configuration Management**: Store and manage RPC URLs and contract addresses
- **Security**: Support for both private keys and geth keyfiles with mutual exclusivity validation
- **Offline Mode**: Create transactions without RPC connection

## Installation

```bash
cd cmd/mito
go build -o mito .
```

## Configuration

Set up default RPC URL and contract address to avoid repeating them in every command:

```bash
# Set RPC URL
./mito config set-rpc http://localhost:8545

# Set ValidatorManager contract address
./mito config set-contract 0x1234567890123456789012345678901234567890

# Show current configuration
./mito config show
```

Configuration is stored in `~/.mito/config.json`.

## Usage

### Transaction Commands

#### Create Transactions (Offline)

```bash
# Create unsigned transaction
./mito tx create collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --unsigned

# Create signed transaction
./mito tx create validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --reward-manager 0x3333... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"my-validator"}' \
  --signed \
  --private-key 0x1234...
```

#### Send Transactions (Online)

```bash
# Send collateral deposit transaction
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --private-key 0x1234... \
  --rpc http://localhost:8545

# Send validator creation transaction
./mito tx send validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --reward-manager 0x3333... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"my-validator"}' \
  --keyfile /path/to/keyfile
```

### Validator Commands

```bash
# Get validator information
./mito validator info --validator-address 0x1111...

# Update validator commission rate
./mito tx send validator update-reward-config \
  --validator 0x1111... \
  --commission-rate 5% \
  --keyfile /path/to/keyfile

# Update validator reward manager
./mito tx send validator update-reward-manager \
  --validator 0x1111... \
  --reward-manager 0x2222... \
  --keyfile /path/to/keyfile
```

### Collateral Management Commands

```bash
# Set permitted collateral owner
./mito tx send collateral set-permitted-owner \
  --validator 0x1111... \
  --collateral-owner 0x2222... \
  --permitted \
  --keyfile /path/to/keyfile

# Revoke collateral owner permission
./mito tx send collateral set-permitted-owner \
  --validator 0x1111... \
  --collateral-owner 0x2222... \
  --keyfile /path/to/keyfile

# Transfer collateral ownership
./mito tx send collateral transfer-ownership \
  --validator 0x1111... \
  --new-owner 0x3333... \
  --keyfile /path/to/keyfile
```

### Security Features

#### Signing Methods (Mutually Exclusive)

Choose one of the following signing methods:

- `--private-key`: Provide private key directly (hex format)
- `--keyfile`: Use geth keyfile (more secure)

```bash
# Using private key
./mito tx send collateral deposit --private-key 0x1234...

# Using keyfile (will prompt for password)
./mito tx send collateral deposit --keyfile /path/to/keyfile

# Using keyfile with password file
./mito tx send collateral deposit \
  --keyfile /path/to/keyfile \
  --keyfile-password-file /path/to/password.txt

# ERROR: Cannot use both
./mito tx send collateral deposit --private-key 0x1234... --keyfile /path/to/keyfile
# Error: flags [private-key keyfile] are mutually exclusive
```

## Command Structure

```
mito
├── tx
│   ├── send
│   │   ├── validator (create, update-operator, update-metadata, update-reward-config, update-reward-manager, unjail)
│   │   └── collateral (deposit, withdraw, set-permitted-owner, transfer-ownership)
│   └── create
│       ├── validator (create, update-operator, update-metadata, update-reward-config, update-reward-manager, unjail)
│       └── collateral (deposit, withdraw, set-permitted-owner, transfer-ownership)
├── validator
│   └── info
└── config
    ├── set-rpc
    ├── set-contract
    └── show
```

## Examples

### Complete Validator Setup

```bash
# 1. Configure defaults
./mito config set-rpc http://localhost:8545
./mito config set-contract 0x1234567890123456789012345678901234567890

# 2. Create validator
./mito tx send validator create \
  --pubkey 0x1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111 \
  --operator 0x2222222222222222222222222222222222222222 \
  --reward-manager 0x3333333333333333333333333333333333333333 \
  --commission-rate 5% \
  --initial-collateral 10.0 \
  --metadata '{"name":"MyValidator","description":"My awesome validator"}' \
  --keyfile /path/to/keyfile

# 3. Check validator info
./mito validator info --validator-address 0x1111111111111111111111111111111111111111

# 4. Deposit additional collateral
./mito tx send collateral deposit \
  --validator 0x1111111111111111111111111111111111111111 \
  --amount 5.0 \
  --keyfile /path/to/keyfile

# 5. Update commission rate
./mito tx send validator update-reward-config \
  --validator 0x1111111111111111111111111111111111111111 \
  --commission-rate 3% \
  --keyfile /path/to/keyfile

# 6. Set permitted collateral owner
./mito tx send collateral set-permitted-owner \
  --validator 0x1111111111111111111111111111111111111111 \
  --collateral-owner 0x4444444444444444444444444444444444444444 \
  --permitted \
  --keyfile /path/to/keyfile
```

### Offline Transaction Creation

```bash
# Create unsigned transaction for later signing
./mito tx create validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --reward-manager 0x3333... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"test"}' \
  --unsigned \
  --output validator-create.json

# Create signed transaction offline
./mito tx create collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --signed \
  --private-key 0x1234... \
  --offline \
  --chain-id 1337 \
  --gas-limit 500000 \
  --gas-price 20000000000 \
  --nonce 42 \
  --fee 0.1
```

## Error Handling

The tool provides clear error messages and usage instructions:

- **Missing required flags**: Shows which flags are required
- **Mutually exclusive flags**: Explains which flags cannot be used together
- **Invalid addresses**: Validates Ethereum address format
- **Configuration issues**: Guides users to set up missing configuration

## Development

To add new commands or modify existing ones, see the source code in the `cmd/` directory. The tool uses a modular structure with:

- `commands.go`: Command hierarchy definition
- `common.go`: Shared utilities and flag validation
- `config.go`: Configuration management
- `tx_*.go`: Transaction-specific implementations
- `validator_*.go`: Validator-specific implementations 