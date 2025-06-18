# Mito - Mitosis Chain Utilities

Mito is a command-line tool for interacting with the Mitosis blockchain, providing utilities for validator management, transaction handling, and EVM operations.

## Features

- **Unified Unit Handling**: Consistent support for wei, gwei, and MITO units across all commands
- **Online/Offline Modes**: Create transactions with or without network connectivity
- **Smart Configuration**: Automatic parameter resolution from config files and network
- **Transaction Management**: Create and send transactions with support for both signed and unsigned modes
- **Validator Operations**: Create, update, and manage validators
- **Collateral Management**: Deposit and withdraw collateral for validators
- **Wallet Management**: Create, import, and manage wallets with cast-compatible keystore support
- **Advanced Security**: Support for private keys and geth keyfiles with comprehensive validation
- **Intelligent Error Handling**: Clear error messages and automatic fallbacks

## Installation

```bash
cd cmd/mito
go build -o mito .
```

## Quick Start

```bash
# Configure defaults (one-time setup)
./mito config set-rpc https://rpc.dognet.mitosis.org
./mito config set-contract --validator-manager 0xECF7658978A03b3A35C2c5B33C449D74E8151Db0

# Create account with mito wallet or cast wallet
./mito wallet import my-val --priv-validator-key ~/.mitosisd/config/priv_validator_key.json

# Create a validator
./mito tx send validator create \
  --pubkey 0x1234... \
  --operator 0x5678... \
  --reward-manager 0x9abc... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"MyValidator"}' \
  --account my-val

# Check validator info
./mito query validator info --address 0x1234...
```

## Configuration Management

Mito supports both explicit flags and configuration files for seamless usage:

```bash
# Set default RPC URL
./mito config set-rpc https://rpc.dognet.mitosis.org

# Set ValidatorManager contract address
./mito config set-contract --validator-manager 0xECF7658978A03b3A35C2c5B33C449D74E8151Db0

# Set configuration for specific network
./mito config set-rpc https://custom-rpc.example.com --network testnet
./mito config set-contract --validator-manager 0x1234... --network testnet

# View current configuration
./mito config show
```

**Output Example:**
```
===== Current Configuration =====

[default]
rpc-url                                = https://rpc.dognet.mitosis.org
validator-manager-contract-address     = 0xECF7658978A03b3A35C2c5B33C449D74E8151Db0

[testnet]
rpc-url                                = https://testnet-rpc.example.com
validator-manager-contract-address     = 0x1234...

Config file location: /Users/user/.mito/config.json
```

Configuration is stored in `~/.mito/config.json` and automatically loaded for all commands. You can configure different networks and switch between them using the `--network` flag.

### Network Configuration

Mito supports multiple network configurations. The default network is called `default`, but you can create and use custom network configurations:

```bash
# Configure mainnet
./mito config set-rpc https://mainnet-rpc.example.com --network mainnet
./mito config set-contract --validator-manager 0x1234... --network mainnet

# Configure testnet
./mito config set-rpc https://testnet-rpc.example.com --network testnet
./mito config set-contract --validator-manager 0x1234... --network testnet

# Use specific network for commands
./mito query validator info --address 0x1234... --network testnet
./mito tx send validator create --pubkey 0x1234... --network mainnet
```

## Unit Handling

Mito provides unified unit handling across all monetary values:

### Supported Units
- **wei**: Base unit (1 wei)
- **gwei**: Giga-wei (10^9 wei) - **Default for gas prices and fees**
- **mito**: Native token unit (10^18 wei) - **Default for collateral amounts**

### Unit Examples
```bash
# Gas price examples (default: gwei)
--gas-price 20          # 20 gwei
--gas-price 20gwei      # 20 gwei
--gas-price 20000000000wei  # 20 gwei
--gas-price 0.00000002mito   # 20 gwei

# Contract fee examples (default: gwei)
--contract-fee 100      # 100 gwei
--contract-fee 100gwei  # 100 gwei
--contract-fee 0.0000001mito # 100 gwei

# Collateral examples (default: mito)
--amount 1.5            # 1.5 MITO
--initial-collateral 10 # 10 MITO
```

## Transaction Commands
`mito` supports several way to signing transaction.

```bash
# Signing with raw private key
./mito tx [create/send] [contract] [method]
  --private-key 0x1234...

# Signing with keyfile (interactive password input)
./mito tx [create/send] [contract] [method]
  --keyfile [keyfile-path]

# Signing with keyfile (password file)
./mito tx [create/send] [contract] [method]
  --keyfile [keyfile-path]
  --keyfile-password-file [keyfile-password-path]

# Signing with keyfile (insecure password)
./mito tx [create/send] [contract] [method]
  --keyfile [keyfile-path]
  --keyfile-password [keyfile-password-string]

# Signing with wallet (~/.mito/keystores)
./mito tx [create/send] [contract] [method]
  --account [account-name]

# Signing with wallet (change keystore directory)
./mito tx [create/send] [contract] [method]
  --account [account-name]
  --keystore-dir [keystore-directory-path]

# Signing with priv-validator-key.json
./mito tx [create/send] [contract] [method]
  --priv-validator-key [priv-validator-key-path]
```

### Create Transactions (tx create)
Creates transactions without sending them to the network:

```bash
# Create unsigned transaction
./mito tx create [contract] [method] \
  --unsigned \
  --nonce 10

# Save to file
./mito tx create [contract] [method] \
  [use signing method]
  --output ./transaction.json
```

### Send Transactions (tx send)
Creates, signs, and broadcasts transactions to the network:

```bash
# Send with keyfile
./mito tx send [contract] [method] \
  [use signing method]

# Skip confirmation prompt
./mito tx send collateral withdraw \
  [use signing method]
  --yes
```

### Contracts

#### Validator

Create validator
```bash
./mito tx [create/send] validator create \
  --pubkey 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
  --operator 0x1234567890123456789012345678901234567890 \
  --reward-manager 0x2345678901234567890123456789012345678901 \
  --commission-rate 5% \
  --initial-collateral 10.0 \
  --metadata '{"name":"MyValidator","description":"Professional validator"}' \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

Update Validator
```bash
# Update commission rate
./mito tx [create/send] validator update-reward-config \
  --validator 0x1111... \
  --commission-rate 3% \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want

# Update operator
./mito tx [create/send] validator update-operator \
  --validator 0x1111... \
  --operator 0x2222... \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want

# Update metadata
./mito tx [create/send] validator update-metadata \
  --validator 0x1111... \
  --metadata '{"name":"UpdatedValidator","version":"2.0"}' \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want

# Unjail validator
./mito tx [create/send] validator unjail \
  --validator 0x1111... \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

#### Collateral

Deposit Collateral
```bash
./mito tx [create/send] collateral deposit \
  --validator 0x1234567890123456789012345678901234567890 \
  --amount 1.5 \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

Withdraw Collateral
```bash
./mito tx [create/send] collateral withdraw \
  --validator 0x1234567890123456789012345678901234567890 \
  --amount 2.0 \
  --receiver 0x2345678901234567890123456789012345678901 \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

Set Permitted Owner
```bash
# Grant permission
./mito tx [create/send] collateral set-permitted-owner \
  --validator 0x1234567890123456789012345678901234567890 \
  --collateral-owner 0x2345678901234567890123456789012345678901 \
  --permitted \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want

# Revoke permission
./mito tx [create/send] collateral set-permitted-owner \
  --validator 0x1234567890123456789012345678901234567890 \
  --collateral-owner 0x2345678901234567890123456789012345678901 \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

Transfer Ownership
```bash
./mito tx [create/send] collateral transfer-ownership \
  --validator 0x1234567890123456789012345678901234567890 \
  --new-owner 0x3456789012345678901234567890123456789012 \
  --keyfile ./keystore/keyfile
  # or use any of signing method what you want
```

## Query Commands

### Query Validator Information
```bash
# Get detailed validator information
./mito query validator info --address 0x1111...

# Get detailed validator information with collateral owner limits
./mito query validator info --address 0x1111... --head 5
./mito query validator info --address 0x1111... --tail 3

# Get validator configuration
./mito query validator config

# Query with specific network
./mito query validator info --address 0x1111... --network testnet
```

### Query Contract Information
```bash
# Query current validator contracts from the ValidatorManager contract
./mito query contract validator

# Query with specific network
./mito query contract validator --network testnet
```

## Wallet Management

Mito provides built-in wallet management capabilities fully compatible with cast (Foundry) keystore format. You can use either mito wallet commands or cast wallet commands interchangeably.

### Create New Wallet
```bash
# Create new random keypair (display only)
./mito wallet new

# Create new keypair and save to default keystore (~/.mito/keystores) with account name
./mito wallet new my-validator
# Enter password when prompted

# Create new keypair and save to custom keystore directory
./mito wallet new my-validator --keystore-dir ~/.foundry/keystores
# Enter password when prompted

# Create with unsafe password (not recommended)
./mito wallet new my-validator --unsafe-password mypassword

# Create with custom keystore directory and unsafe password
./mito wallet new my-validator --keystore-dir ./custom-keystores --unsafe-password mypassword
```

### Generate Mnemonic
```bash
# Generate 12-word mnemonic (default)
./mito wallet new-mnemonic

# Generate 24-word mnemonic with multiple accounts
./mito wallet new-mnemonic --words 24 --accounts 3
```

### Import Existing Wallet
```bash
# Import from private key (interactive)
./mito wallet import my-validator --interactive

# Import from private key (non-interactive)
./mito wallet import my-validator --private-key 0x1234... --unsafe-password mypass

# Import from mnemonic
./mito wallet import my-validator \
  --mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" \
  --unsafe-password mypass

# Import with specific mnemonic index
./mito wallet import my-validator \
  --mnemonic "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about" \
  --mnemonic-index 1 \
  --unsafe-password mypass

# Import to custom directory
./mito wallet import my-validator \
  --private-key 0x1234... \
  --keystore-dir ./custom-keystores \
  --unsafe-password mypass

# Import with priv-validator-key.json
./mito wallet import my-validator \
  --priv-validator-key ~/.mito/config/priv-validator-key.json
```

### List Wallets
```bash
# List wallets in default directory (~/.mito/keystores)
./mito wallet list

# List wallets in custom directory
./mito wallet list --dir ./custom-keystores
```

### Export Private Key
```bash
# Export private key from default keystore
./mito wallet export my-validator
# Enter password when prompted

# Export from custom keystore directory
./mito wallet export my-validator --keystore-dir ./custom-keystores

# Export with unsafe password (not recommended)
./mito wallet export my-validator --unsafe-password mypassword
```

### Delete Wallet
```bash
# Delete wallet with confirmation prompt
./mito wallet delete my-validator

# Delete from custom keystore directory
./mito wallet delete my-validator --keystore-dir ./custom-keystores

# Delete without confirmation (skip prompt)
./mito wallet delete my-validator --yes
```

### Keystore Compatibility

All mito wallet commands create keystores compatible with cast (Foundry) and geth. The keystores are stored in `~/.mito/keystores` by default, following the same format as cast's `~/.foundry/keystores`.

```bash
# You can use cast commands with mito keystores
cast wallet list --dir ~/.mito/keystores

# And use mito commands with cast keystores
./mito wallet list --dir ~/.foundry/keystores

# Use mito keystores with mito transactions
./mito tx send validator create \
  --pubkey 0x1234... \
  --operator 0x5678... \
  --commission-rate 5% \
  --keyfile ~/.mito/keystores/my-validator
```

### Keyfile Compatibility

Mito is compatible with standard Ethereum keystore formats ([EIP-2335](https://eips.ethereum.org/EIPS/eip-2335)) used by both **cast** (Foundry) and **geth**. You can use either **mito wallet** or **cast wallet** for key management.

#### Key Management Options

Both mito and cast provide excellent key management capabilities:

```bash
# Option 1: Using mito wallet (built-in)
./mito wallet new my-validator-key
./mito wallet list

# Option 2: Using cast wallet (external)
cast wallet new my-validator-key
cast wallet import my-validator-key --interactive
cast wallet list

# Both create compatible keystores
./mito tx send validator create \
  --pubkey 0x1234... \
  --operator 0x1234... \
  --commission-rate 5% \
  --account my-validator-key
  # or --keyfile ~/.foundry/keystores/my-validator-key
```

## Advanced Usage

### Custom Gas Settings
```bash
./mito tx send validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --commission-rate 5% \
  --gas-limit 750000 \
  --gas-price 30gwei \
  --keyfile ./keystore/keyfile
```

### Custom Contract Fee
```bash
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --contract-fee 150gwei \
  --keyfile ./keystore/keyfile
```

### Network Override
```bash
# Override RPC URL for single command
./mito tx send validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --commission-rate 5% \
  --rpc-url https://custom-rpc.example.com \
  --keyfile ./keystore/keyfile

# Use specific network configuration
./mito tx send validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --commission-rate 5% \
  --network testnet \
  --keyfile ./keystore/keyfile
```
