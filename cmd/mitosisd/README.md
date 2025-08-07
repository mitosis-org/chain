# mitosisd

`mitosisd` is the main daemon for running a Mitosis Chain validator node. It handles consensus operations, EVM integration, and blockchain state management.

## Commands

### Node Management
```bash
# Initialize a new node
mitosisd init [moniker] --chain-id [chain-id]

# Start the node
mitosisd start

# Show node information
mitosisd comet show-node-id
mitosisd comet show-validator
```

### Genesis Configuration

#### Adding Smart Contracts to Genesis

Pre-deploy smart contracts directly in the genesis block using Foundry compilation artifacts.

**Prerequisites**
```bash
# Compile your contract
forge build path/to/YourContract.sol
```

**Usage**
```bash
mitosisd genesis add-contract [contract-address] [artifact-file] [flags]
```

**Examples**
```bash
# Basic usage
mitosisd genesis add-contract \
  0x1234567890123456789012345678901234567890 \
  out/SimpleStorage.sol/SimpleStorage.json

# With initial balance
mitosisd genesis add-contract \
  0x1234567890123456789012345678901234567890 \
  out/SimpleStorage.sol/SimpleStorage.json \
  --balance 1000000000000000000

# Use creation bytecode
mitosisd genesis add-contract \
  0x1234567890123456789012345678901234567890 \
  out/SimpleStorage.sol/SimpleStorage.json \
  --use-creation-code
```

**Flags**
- `--balance`: Initial ETH balance for the contract account (default: "0")
- `--use-creation-code`: Use creation bytecode instead of deployed bytecode

#### Other Genesis Commands
```bash
# Add genesis account
mitosisd genesis add-genesis-account [address] [coins]

# Add validator
mitosisd genesis add-validator [validator-info]

# Validate genesis file
mitosisd genesis validate
```

### Key Management
```bash
# Create new key
mitosisd keys add [key-name]

# List keys
mitosisd keys list

# Show key info
mitosisd keys show [key-name]
```

### Query Commands
```bash
# Query account
mitosisd query auth account [address]

# Query validator
mitosisd query evmvalidator validator [validator-id]

# Query block
mitosisd query block [height]
```

## Configuration

Configuration files are stored in `~/.mitosisd/config/`:
- `config.toml`: Node configuration
- `app.toml`: Application configuration  
- `genesis.json`: Cosmos SDK genesis
- `eth_genesis.json`: Ethereum genesis

## Home Directory

Use `--home` flag to specify custom home directory:
```bash
mitosisd --home /custom/path [command]
```