#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "EC_JWT_FILE: $EC_JWT_FILE"
echo "----------------------"

# sponsor worry verify hobby armed physical olympic find speak wink develop blush
OWNER_MNEMONIC="c3BvbnNvciB3b3JyeSB2ZXJpZnkgaG9iYnkgYXJtZWQgcGh5c2ljYWwgb2x5bXBpYyBmaW5kIHNwZWFrIHdpbmsgZGV2ZWxvcCBibHVzaA=="
# banana omit eye gesture disagree fork zone cup promote plunge neither rug
GEN_VAL_MNEMONIC="YmFuYW5hIG9taXQgZXllIGdlc3R1cmUgZGlzYWdyZWUgZm9yayB6b25lIGN1cCBwcm9tb3RlIHBsdW5nZSBuZWl0aGVyIHJ1Zw=="
GEN_VAL_MONIKER="Mitosis Foundation"

echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD init "$GEN_VAL_MONIKER" --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

#sleep 10

# Setup genesis account for owner
OWNER=$(echo "$OWNER_MNEMONIC" | base64 -d | $MITOSISD keys add owner --recover --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME" --output json | jq -r .address)
$MITOSISD genesis add-genesis-account owner 999999999000000ustake --keyring-backend test --home "$MITOSISD_HOME" # (1B - 1) * 1M ustake

# Setup genesis account for genesis validator
echo "$GEN_VAL_MNEMONIC" | base64 -d | $MITOSISD keys add genval --recover --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME"
$MITOSISD genesis add-genesis-account genval 1000000ustake --keyring-backend test --home "$MITOSISD_HOME"

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

# Setup gentx on the genesis
$MITOSISD genesis gentx genval 1000000ustake --chain-id "$MITOSISD_CHAIN_ID" --keyring-backend test --home "$MITOSISD_HOME"
$MITOSISD genesis collect-gentxs --home "$MITOSISD_HOME"
