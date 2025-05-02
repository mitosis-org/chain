#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "EC_JWT_FILE: $EC_JWT_FILE"
echo "GOV_ENTRYPOINT: $GOV_ENTRYPOINT"
echo "----------------------"

# banana omit eye gesture disagree fork zone cup promote plunge neither rug
GEN_VAL_MNEMONIC="YmFuYW5hIG9taXQgZXllIGdlc3R1cmUgZGlzYWdyZWUgZm9yayB6b25lIGN1cCBwcm9tb3RlIHBsdW5nZSBuZWl0aGVyIHJ1Zw=="

echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD init validator --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

GENESIS_FILE="$MITOSISD_HOME/config/genesis.json"
TEMP="$MITOSISD_HOME/config/genesis.json.tmp"

# Setup execution block hash on the genesis
hash=$(cast block --rpc-url http://127.0.0.1:8545 | grep hash | awk '{print $2}' | xxd -r -p | base64)
jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup additional modifications on the genesis
jq '.consensus.params.block.max_bytes = "-1"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.consensus.params.validator.pub_key_types = ["secp256k1"]' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup app.toml
sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001ustake"/' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@pruning = "default"@pruning = "nothing"@' "$MITOSISD_HOME"/config/app.toml # archive node
#sed -i.bak'' 's/mock = false/mock = true/' "$MITOSISD_HOME"/config/app.toml # Comment out this line to mock execution engine instead of using real execution client.
sed -i.bak'' 's@endpoint = "http://127.0.0.1:8551"@endpoint = "http://127.0.0.1:8551"@' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@jwt-file = ""@jwt-file = "'"$EC_JWT_FILE"'"@' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@entrypoint = "0x0000000000000000000000000000000000000000"@entrypoint = "'"$GOV_ENTRYPOINT"'"@' "$MITOSISD_HOME"/config/app.toml

# Setup config.toml
sed -i.bak'' 's/type = "flood"/type = "nop"/' "$MITOSISD_HOME"/config/config.toml # we don't use mempool in consensus layer
sed -i.bak'' 's/broadcast = true/broadcast = false/' "$MITOSISD_HOME"/config/config.toml # we don't use mempool in consensus layer
sed -i.bak'' 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$MITOSISD_HOME"/config/config.toml

# Get validator pubkey in compressed format for a genesis validator
VAL_PRIVKEY_FILE="$MITOSISD_HOME/config/priv_validator_key.json"
COMPRESSED_PUBKEY=$(jq -r ".pub_key.value" "$VAL_PRIVKEY_FILE" | base64 -d | xxd -p -c 1000 | tr '[:lower:]' '[:upper:]')
VAL_ADDR="0x$(echo -n "$COMPRESSED_PUBKEY" | xxd -r -p | sha256sum | head -c 40)"

# Add validator to evmvalidator genesis state
# Parameters: pubkey, collateral_owner, collateral (gwei), extra_voting_power, jailed
$MITOSISD add-genesis-validator "$COMPRESSED_PUBKEY" "$VAL_ADDR" 1000000000000000 0 false --home "$MITOSISD_HOME" # 1M MITO as collateral

# Comment out if you need to collect gentxs
#$MITOSISD genesis collect-gentxs --home "$MITOSISD_HOME"
