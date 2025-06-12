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
./mito wallet new my-validator
# or: cast wallet new my-validator

# Create a validator (online mode)
./mito tx send validator create \
  --pubkey 0x1234... \
  --operator 0x5678... \
  --reward-manager 0x9abc... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"MyValidator"}' \
  --keyfile ./my-validator

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

### Output Format
All monetary values are displayed in both MITO and Gwei for clarity:

```
Fee                        : 0.000000123456789 MITO (123.456789 Gwei)
Total Value                : 1.500000123456789 MITO (1500000123.456789 Gwei)
Gas Price                  : 0.000000025 MITO (25 Gwei)
```

## Online vs Offline Modes

### Online Mode (Default)
Requires RPC connection and automatically detects:
- Chain ID
- Gas price
- Nonce (for sending transactions)

```bash
# Online transaction (uses config defaults)
./mito tx send validator create \
  --pubkey 0x1234... \
  --operator 0x5678... \
  --commission-rate 5% \
  --keyfile ~/.foundry/keystores/my-validator
```

### Offline Mode
Create transactions without network connectivity:

```bash
# Offline unsigned transaction
./mito tx create validator create \
  --pubkey 0x1234... \
  --operator 0x5678... \
  --commission-rate 5% \
  --unsigned \
  --nonce 10 \
  --gas-price 20gwei \
  --chain-id 125883

# Offline signed transaction
./mito tx create collateral deposit \
  --validator 0x1234... \
  --amount 1.5 \
  --signed \
  --private-key 0xabc... \
  --nonce 15 \
  --gas-price 25gwei \
  --chain-id 125883
```

## Transaction Commands

### Create Transactions (tx create)
Creates transactions without sending them to the network:

```bash
# Create unsigned transaction
./mito tx create collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --unsigned \
  --nonce 10

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

# Save to file
./mito tx create collateral withdraw \
  --validator 0x1111... \
  --amount 0.5 \
  --receiver 0x4444... \
  --unsigned \
  --nonce 11 \
  --output ./transaction.json
```

### Send Transactions (tx send)
Creates, signs, and broadcasts transactions to the network:

```bash
# Send with keyfile
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --keyfile ./keystore/keyfile \
  --keyfile-password-file ./password.txt

# Send with private key
./mito tx send validator create \
  --pubkey 0x1111... \
  --operator 0x2222... \
  --reward-manager 0x3333... \
  --commission-rate 5% \
  --initial-collateral 1.5 \
  --metadata '{"name":"my-validator"}' \
  --private-key 0x1234...

# Skip confirmation prompt
./mito tx send collateral withdraw \
  --validator 0x1111... \
  --amount 0.5 \
  --receiver 0x4444... \
  --keyfile ./keystore/keyfile \
  --yes
```

## Validator Management

### Create Validator
```bash
./mito tx send validator create \
  --pubkey 0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef \
  --operator 0x1234567890123456789012345678901234567890 \
  --reward-manager 0x2345678901234567890123456789012345678901 \
  --commission-rate 5% \
  --initial-collateral 10.0 \
  --metadata '{"name":"MyValidator","description":"Professional validator"}' \
  --keyfile ./keystore/keyfile
```

### Update Validator
```bash
# Update commission rate
./mito tx send validator update-reward-config \
  --validator 0x1111... \
  --commission-rate 3% \
  --keyfile ./keystore/keyfile

# Update operator
./mito tx send validator update-operator \
  --validator 0x1111... \
  --operator 0x2222... \
  --keyfile ./keystore/keyfile

# Update metadata
./mito tx send validator update-metadata \
  --validator 0x1111... \
  --metadata '{"name":"UpdatedValidator","version":"2.0"}' \
  --keyfile ./keystore/keyfile

# Unjail validator
./mito tx send validator unjail \
  --validator 0x1111... \
  --keyfile ./keystore/keyfile
```

## Query Commands

### Query Validator Information
```bash
# Get detailed validator information
./mito query validator info --address 0x1111...

# Get detailed validator information with collateral owner limits
./mito query validator info --address 0x1111... --head 5 --tail 3

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

## Collateral Management

### Deposit Collateral
```bash
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 5.0 \
  --keyfile ./keystore/keyfile
```

### Withdraw Collateral
```bash
./mito tx send collateral withdraw \
  --validator 0x1111... \
  --amount 2.0 \
  --receiver 0x2222... \
  --keyfile ./keystore/keyfile
```

### Manage Collateral Permissions
```bash
# Grant permission
./mito tx send collateral set-permitted-owner \
  --validator 0x1111... \
  --collateral-owner 0x2222... \
  --permitted \
  --keyfile ./keystore/keyfile

# Revoke permission
./mito tx send collateral set-permitted-owner \
  --validator 0x1111... \
  --collateral-owner 0x2222... \
  --keyfile ./keystore/keyfile

# Transfer ownership
./mito tx send collateral transfer-ownership \
  --validator 0x1111... \
  --new-owner 0x3333... \
  --keyfile ./keystore/keyfile
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

## Security Features

### Signing Methods (Mutually Exclusive)
Choose exactly one signing method:

```bash
# Method 1: Keyfile (recommended)
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --keyfile ./keystore/keyfile

# Method 2: Keyfile with password file
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --keyfile ./keystore/keyfile \
  --keyfile-password-file ./password.txt

# Method 3: Private key
./mito tx send collateral deposit \
  --validator 0x1111... \
  --amount 1.5 \
  --private-key 0x1234...
```

### Keyfile Compatibility

Mito is compatible with standard Ethereum keystore formats ([EIP-2335](https://eips.ethereum.org/EIPS/eip-2335)) used by both **cast** (Foundry) and **geth**. You can use either **mito wallet** or **cast wallet** for key management.

#### Key Management Options

Both mito and cast provide excellent key management capabilities:

```bash
# Option 1: Using mito wallet (built-in)
./mito wallet new my-validator-key
./mito wallet import my-validator-key --interactive
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
  --keyfile ~/.mito/keystores/my-validator-key
  # or --keyfile ~/.foundry/keystores/my-validator-key
```

### Transaction Modes (Mutually Exclusive)
Choose exactly one transaction mode:

```bash
# Unsigned transaction (requires --nonce for offline mode)
./mito tx create validator create \
  --pubkey 0x1111... \
  --unsigned \
  --nonce 10

# Signed transaction (requires signing method)
./mito tx create validator create \
  --pubkey 0x1111... \
  --signed \
  --private-key 0x1234...
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

## Output Examples

### Transaction Creation
```
===== Deposit Collateral Transaction =====
Validator Address          : 0x1234567890123456789012345678901234567890
Collateral Amount          : 1.5 MITO
Fee                        : 0.000000123456789 MITO (123.456789 Gwei)
Total Value                : 1.500000123456789 MITO (1500000123.456789 Gwei)

🚨 IMPORTANT: Only permitted collateral owners can deposit collateral for a validator.
If your address is not a permitted collateral owner for this validator, the transaction will fail.

{"chainId":"0x7fff...","gas":"0x7a120","gasPrice":"0x4a817c800",...}
```

### Configuration Display
```
===== Current Configuration =====
RPC URL                      : https://rpc.dognet.mitosis.org
ValidatorManager Contract    : 0xECF7658978A03b3A35C2c5B33C449D74E8151Db0
Chain ID                     : (not set)
Config file location         : /Users/user/.mito/config.json
```

## Error Handling

Mito provides comprehensive error handling and validation:

### Missing Configuration
```bash
$ ./mito tx create validator create --pubkey 0x1111... --unsigned
Error: RPC connection required to get chain ID automatically.
Please provide --chain-id manually or set RPC URL with --rpc-url or 'mito config set-rpc'
```

### Flag Validation
```bash
$ ./mito tx create validator create --unsigned
Error: when using --unsigned, --nonce must be provided

$ ./mito tx send validator create --private-key 0x123... --keyfile ./key
Error: flags [private-key keyfile] are mutually exclusive
```

### Invalid Values
```bash
$ ./mito tx create collateral deposit --amount invalid
Error: invalid decimal format: invalid

$ ./mito tx create validator create --commission-rate 150%
Error: commission rate must be between 0% and 100%
```

## Command Reference

```
mito
├── tx
│   ├── create (offline transaction creation)
│   │   ├── validator
│   │   │   ├── create
│   │   │   ├── update-operator
│   │   │   ├── update-metadata
│   │   │   ├── update-reward-config
│   │   │   ├── update-reward-manager
│   │   │   └── unjail
│   │   └── collateral
│   │       ├── deposit
│   │       ├── withdraw
│   │       ├── set-permitted-owner
│   │       └── transfer-ownership
│   └── send (online transaction broadcasting)
│       ├── validator (same subcommands as create)
│       └── collateral (same subcommands as create)
├── query
│   ├── validator
│   │   ├── info
│   │   └── config
│   └── contract
│       └── validator
├── wallet
│   ├── new (create new random keypair)
│   ├── new-mnemonic (generate BIP39 mnemonic)
│   ├── import (import private key or mnemonic)
│   ├── list (list accounts in keystore)
│   ├── export (export private key from keystore)
│   └── delete (delete keystore file)
├── config
│   ├── set-rpc
│   ├── set-contract
│   └── show
└── version
```

## Development

The tool uses a modular architecture with:

- **Configuration Management**: `internal/config/`
- **Transaction Building**: `internal/tx/`
- **Output Formatting**: `internal/output/`
- **Unit Conversion**: `internal/units/`
- **Validation**: `internal/validation/`
- **Client Operations**: `internal/client/`
- **Utility Functions**: `internal/utils/`
- **Flag Management**: `internal/flags/`
- **Dependency Injection**: `internal/container/`
- **Command Structure**: `commands/tx/`, `commands/query/`, `commands/config/`, `commands/version/`

### Adding New Commands
1. Create command file in appropriate `commands/` directory
2. Add validation logic in `internal/validation/`
3. Add transaction logic in `internal/tx/`
4. Add output formatting in `internal/output/`
5. Register command in root command structure (`commands/root.go`)
6. Add any required client operations in `internal/client/`
