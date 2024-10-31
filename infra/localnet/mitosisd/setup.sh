#!/bin/bash

set -e
set -x

echo "----- Input Envs -----"
echo "CHAIN_ID: $CHAIN_ID"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "GETH_INFRA_DIR: $GETH_INFRA_DIR"
echo "----------------------"

./build/mitosisd init validator --chain-id "$CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"
./build/mitosisd config set client chain-id "$CHAIN_ID" --home "$MITOSISD_HOME"
./build/mitosisd config set client keyring-backend test --home "$MITOSISD_HOME"

OWNER=$(./build/mitosisd keys add owner --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME" --output json | jq -r .address)
./build/mitosisd genesis add-genesis-account owner 80000000ustake --keyring-backend test --home "$MITOSISD_HOME"

./build/mitosisd keys add validator --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME"
./build/mitosisd genesis add-genesis-account validator 20000000ustake --keyring-backend test --home "$MITOSISD_HOME"

# modify genesis.json
GENESIS="$MITOSISD_HOME/config/genesis.json"
TEMP="$MITOSISD_HOME/config/genesis.json.tmp"

jq '.consensus.params.block.max_bytes = "-1"' "$GENESIS" > "$TEMP" && mv "$TEMP" "$GENESIS"
jq '.consensus.params.validator.pub_key_types = ["secp256k1"]' "$GENESIS" > "$TEMP" && mv "$TEMP" "$GENESIS"
hash=$(cast block --rpc-url http://127.0.0.1:8545 | grep hash | awk '{print $2}' | xxd -r -p | base64)
jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$GENESIS" > "$TEMP" && mv "$TEMP" "$GENESIS"
jq '.app_state.authority.owner = "'"$OWNER"'"' "$GENESIS" > "$TEMP" && mv "$TEMP" "$GENESIS"

# modify app.toml
sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001ustake"/' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@pruning = "default"@pruning = "nothing"@' "$MITOSISD_HOME"/config/app.toml # archiving mode
#sed -i.bak'' 's/mock = false/mock = true/' "$MITOSISD_HOME"/config/app.toml # Comment out this line to mock execution engine instead of using real geth.
sed -i.bak'' 's@endpoint = ""@endpoint = "http://127.0.0.1:8551"@' "$MITOSISD_HOME"/config/app.toml
sed -i.bak'' 's@jwt-file = ""@jwt-file = "'"$GETH_INFRA_DIR"'/jwt.hex"@' "$MITOSISD_HOME"/config/app.toml

# modify config.toml
sed -i.bak'' 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$MITOSISD_HOME"/config/config.toml
sed -i.bak'' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$MITOSISD_HOME"/config/config.toml

./build/mitosisd genesis gentx validator 10000000ustake --chain-id "$CHAIN_ID" --keyring-backend test --home "$MITOSISD_HOME"
./build/mitosisd genesis collect-gentxs --home "$MITOSISD_HOME"
