#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"

echo "EC_GENESIS_BLOCK_HASH_FILE: $EC_GENESIS_BLOCK_HASH_FILE"

echo "OWNER: $OWNER"

echo "GEN_VAL_MONIKER: $GEN_VAL_MONIKER"
echo "GEN_VAL_MNEMONIC (sha256): $(echo -n "$GEN_VAL_MNEMONIC" | sha256sum)"

echo "ARTIFACTS_DIR: $ARTIFACTS_DIR"
echo "----------------------"

EC_GENESIS_BLOCK_HASH=$(xargs < "$EC_GENESIS_BLOCK_HASH_FILE") # trim whitespace

GENESIS_FILE="$MITOSISD_HOME/config/genesis.json"
TEMP="$MITOSISD_HOME/config/genesis.json.tmp"

# Init for genesis validator
echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD init "$GEN_VAL_MONIKER" --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

# Setup genesis account for owner
$MITOSISD genesis add-genesis-account "$OWNER" 999999999000000ustake --keyring-backend test --home "$MITOSISD_HOME" # (1B - 1) * 1M ustake

# Setup genesis account for genesis validator
echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD keys add genval --recover --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME"
$MITOSISD genesis add-genesis-account genval 1000000ustake --keyring-backend test --home "$MITOSISD_HOME"

# Setup execution block hash on the genesis
hash=$(echo -n "$EC_GENESIS_BLOCK_HASH" | xxd -r -p | base64)
jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup additional modifications on the genesis
jq '.consensus.params.block.max_bytes = "-1"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.consensus.params.validator.pub_key_types = ["secp256k1"]' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.app_state.staking.params.unbonding_time = "1s"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.app_state.authority.owner = "'"$OWNER"'"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup gentx on the genesis
$MITOSISD genesis gentx genval 1000000ustake --chain-id "$MITOSISD_CHAIN_ID" --keyring-backend test --home "$MITOSISD_HOME"
$MITOSISD genesis collect-gentxs --home "$MITOSISD_HOME"

# artifacts
mkdir -p "$ARTIFACTS_DIR"
cp "$GENESIS_FILE" "$ARTIFACTS_DIR/genesis.json"
