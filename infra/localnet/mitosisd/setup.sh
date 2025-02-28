#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "EC_JWT_FILE: $EC_JWT_FILE"
echo "----------------------"

# banana omit eye gesture disagree fork zone cup promote plunge neither rug
GEN_VAL_MNEMONIC="YmFuYW5hIG9taXQgZXllIGdlc3R1cmUgZGlzYWdyZWUgZm9yayB6b25lIGN1cCBwcm9tb3RlIHBsdW5nZSBuZWl0aGVyIHJ1Zw=="
GEN_VAL_MONIKER="Mitosis Foundation"

echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD init "$GEN_VAL_MONIKER" --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

GENESIS_FILE="$MITOSISD_HOME/config/genesis.json"
TEMP="$MITOSISD_HOME/config/genesis.json.tmp"

# Setup execution block hash on the genesis
hash=$(cast block --rpc-url http://127.0.0.1:8545 | grep hash | awk '{print $2}' | xxd -r -p | base64)
jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup additional modifications on the genesis
jq '.consensus.params.block.max_bytes = "-1"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.consensus.params.validator.pub_key_types = ["secp256k1"]' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.app_state.staking.params.unbonding_time = "1s"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"
jq '.app_state.authority.owner = "'"$OWNER"'"' "$GENESIS_FILE" > "$TEMP" && mv "$TEMP" "$GENESIS_FILE"

# Setup app.toml
sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001ustake"/' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@pruning = "default"@pruning = "nothing"@' "$MITOSISD_HOME"/config/app.toml # archive node
#sed -i.bak'' 's/mock = false/mock = true/' "$MITOSISD_HOME"/config/app.toml # Comment out this line to mock execution engine instead of using real execution client.
sed -i.bak'' 's@endpoint = ""@endpoint = "http://127.0.0.1:8551"@' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@jwt-file = ""@jwt-file = "'"$EC_JWT_FILE"'"@' "$MITOSISD_HOME"/config/app.toml

# Setup config.toml
sed -i.bak'' 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$MITOSISD_HOME"/config/config.toml
sed -i.bak'' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$MITOSISD_HOME"/config/config.toml

# Get validator pubkey in compressed format for evmvalidator
VALIDATOR_PRIVKEY_FILE="$MITOSISD_HOME/config/priv_validator_key.json"
PUBKEY_UNCOMPRESSED=$(jq -r ".pub_key.value" "$VALIDATOR_PRIVKEY_FILE" | base64 -d | xxd -p -c 1000)
# Convert from uncompressed 65-byte to compressed 33-byte pubkey format
# For secp256k1, uncompressed starts with '04', compressed starts with '02' or '03'
# We'll use a temporary Python script for conversion
python3 -c "
import binascii
import sys

# Read uncompressed pubkey
uncompressed_hex = '$PUBKEY_UNCOMPRESSED'
if uncompressed_hex.startswith('04'):
    uncompressed_bytes = binascii.unhexlify(uncompressed_hex)
    # Get X and Y coordinates
    x, y = uncompressed_bytes[1:33], uncompressed_bytes[33:65]
    # Choose prefix based on Y being even or odd
    prefix = '02' if (y[-1] % 2 == 0) else '03'
    # Create compressed pubkey
    compressed_hex = prefix + binascii.hexlify(x).decode('ascii')
    print(compressed_hex)
else:
    # Already compressed or invalid
    print(uncompressed_hex)
" > compressed_pubkey.txt
COMPRESSED_PUBKEY=$(cat compressed_pubkey.txt)
rm compressed_pubkey.txt

# Add validator to evmvalidator genesis state
# Parameters: pubkey, collateral (gwei), extra_voting_power, jailed
$MITOSISD add-genesis-validator "$COMPRESSED_PUBKEY" 1000000000000000000 0 false --home "$MITOSISD_HOME"

# Comment out if you need to collect gentxs
#$MITOSISD genesis collect-gentxs --home "$MITOSISD_HOME"
