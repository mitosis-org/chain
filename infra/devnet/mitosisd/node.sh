#!/bin/bash

set -e

echo "----- Input Envs -----"
##### mandatory
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "PEER_RPC: $PEER_RPC" # empty means a node without peer
echo "PEER_P2P: $PEER_P2P"

##### only for initialization
echo "MODE: $MODE" # archive, full
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "GENESIS_FILE: $GENESIS_FILE"
echo "EC_ENDPOINT: $EC_ENDPOINT"
echo "EC_JWT_FILE: $EC_JWT_FILE"
echo "VAL_MNEMONIC (sha256): $(echo -n "$VAL_MNEMONIC" | sha256sum)"
echo "----------------------"

if find "$MITOSISD_HOME" -mindepth 1 -maxdepth 1 | read; then
  echo "====================================================================="
  echo "[IMPORTANT] Home directory already exists. Skip initialization."
  echo "====================================================================="
else
  echo "[IMPORTANT] Initialize mitosisd home directory: $MITOSISD_HOME"

  if [ -n "$VAL_MNEMONIC" ]; then
    echo "It's a Validator Node."
    echo "$VAL_MNEMONIC" | base64 -d | $MITOSISD init validator --recover --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"

    echo "$VAL_MNEMONIC" | base64 -d | $MITOSISD keys add validator --recover --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME"
  else
    echo "It's a Non-Validator Node."
    $MITOSISD init validator --chain-id "$MITOSISD_CHAIN_ID" --default-denom ustake --home "$MITOSISD_HOME"
  fi

  cp "$GENESIS_FILE" "$MITOSISD_HOME"/config/genesis.json

  app_toml="$MITOSISD_HOME"/config/app.toml
  config_toml="$MITOSISD_HOME"/config/config.toml

  # Setup app.toml
  sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001ustake"/' "$app_toml"
  sed -i.bak'' 's@endpoint = ""@endpoint = "'"$EC_ENDPOINT"'"@' "$app_toml"
  sed -i.bak'' 's@jwt-file = ""@jwt-file = "'"$EC_JWT_FILE"'"@' "$app_toml"

  if [ "$MODE" == "archive" ]; then # archive node
    sed -i.bak'' 's@pruning = "default"@pruning = "nothing"@' "$app_toml"
  elif [ "$MODE" == "full" ]; then # full node
    sed -i.bak'' 's@pruning = "default"@pruning = "default"@' "$app_toml"
  else
    echo "Invalid MODE: $MODE"
    exit 1
  fi

  # Setup config.toml
  sed -i.bak'' 's/timeout_commit = "5s"/timeout_commit = "1s"/' "$config_toml"
  sed -i.bak'' 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["*"\]/' "$config_toml"
fi

# Wait for peer
peer_param=""
if [ -n "$PEER_RPC" ]; then
  until curl -s "$PEER_RPC" 1>/dev/null 2>&1; do
    echo "Waiting for peer: $PEER_RPC"
    sleep 3
  done
  node_id=$(curl -s "$PEER_RPC"/status | jq -r .result.node_info.id)
  peer_param=--p2p.persistent_peers=$node_id@$PEER_P2P
fi

# Start mitosisd
# shellcheck disable=SC2086
$MITOSISD start --home "$MITOSISD_HOME" \
  $peer_param \
  --p2p.laddr=tcp://0.0.0.0:26656 \
  --rpc.laddr=tcp://0.0.0.0:26657 \
  --grpc.enable \
  --grpc.address=0.0.0.0:9090 \
  --api.enable \
  --api.address=tcp://0.0.0.0:1317 \
  --api.enabled-unsafe-cors \
  --log_level "info"
