#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"

echo "EC_GENESIS_BLOCK_HASH_FILE: $EC_GENESIS_BLOCK_HASH_FILE"

echo "GEN_VAL_MNEMONIC (sha256): $(echo -n "$GEN_VAL_MNEMONIC" | sha256sum)"

echo "ARTIFACTS_DIR: $ARTIFACTS_DIR"
echo "----------------------"

EC_GENESIS_BLOCK_HASH=$(xargs < "$EC_GENESIS_BLOCK_HASH_FILE") # trim whitespace

GENESIS_FILE="$MITOSISD_HOME/config/genesis.json"
TEMP="$MITOSISD_HOME/config/genesis.json.tmp"

# Init for genesis validator
echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD init tmp --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

# Setup execution block hash on the genesis
hash=$(echo -n "$EC_GENESIS_BLOCK_HASH" | xxd -r -p | base64)
jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup additional modifications on the genesis
jq '.consensus.params.block.max_bytes = "-1"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.consensus.params.validator.pub_key_types = ["secp256k1"]' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Get validator pubkey in compressed format for a genesis validator
VAL_PRIVKEY_FILE="$MITOSISD_HOME/config/priv_validator_key.json"
COMPRESSED_PUBKEY=$(jq -r ".pub_key.value" "$VAL_PRIVKEY_FILE" | base64 -d | xxd -p -c 1000 | tr '[:lower:]' '[:upper:]')

# Add validator to evmvalidator genesis state
# Parameters: pubkey, collateral (gwei), extra_voting_power, jailed
$MITOSISD add-genesis-validator "$COMPRESSED_PUBKEY" 1000000000000000 0 false --home "$MITOSISD_HOME" # 1M MITO as collateral

# artifacts
mkdir -p "$ARTIFACTS_DIR"
cp "$GENESIS_FILE" "$ARTIFACTS_DIR/genesis.json"
