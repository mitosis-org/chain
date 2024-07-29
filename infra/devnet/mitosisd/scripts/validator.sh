#!/bin/bash

set -e
set -x

echo "----- Input Envs -----"
INFRA_DIR=${INFRA_DIR:-'./infra/devnet'}
DB_DIR=${DB_DIR:-'./tmp/devnet'}

MITOSIS_CHAIN_ID=${MITOSIS_CHAIN_ID:-'mitosis-devnet-1'}
MITOSISD=${MITOSISD:-'./build/mitosisd'}
MITOSISD_GENESIS_FILE=${MITOSISD_GENESIS_FILE:-'./tmp/devnet/register-to-ethos/artifacts/genesis-final.json'}
VAL_KEY_FILE=${VAL_KEY_FILE:-"$DB_DIR/mitosisd/config/priv_validator_key.json"}
EXECUTION_ENDPOINT=${EXECUTION_ENDPOINT:-'http://127.0.0.1:8551'}
PEER_RPC=${PEER_RPC:-''}
PEER_P2P=${PEER_P2P:-''}

echo "INFRA_DIR: $INFRA_DIR"
echo "DB_DIR: $DB_DIR"
echo "MITOSIS_CHAIN_ID: $MITOSIS_CHAIN_ID"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_GENESIS_FILE: $MITOSISD_GENESIS_FILE"
echo "VAL_KEY_FILE: $VAL_KEY_FILE"
echo "EXECUTION_ENDPOINT: $EXECUTION_ENDPOINT"
echo "PEER_RPC: $PEER_RPC"
echo "PEER_P2P: $PEER_P2P"
echo "----------------------"

MITOSIS_HOME="$DB_DIR/mitosisd"
MITOSIS_DENOM='thai'

EXECUTION_JWT_FILE="$INFRA_DIR/geth/config/common/jwt.hex"

if [ -d "$MITOSIS_HOME" ]; then
  echo "====================================================================="
  echo "[IMPORTANT] Home directory already exists. Skip initialization."
  echo "====================================================================="
else
  echo "[IMPORTANT] Initialize mitosisd home directory: $MITOSIS_HOME"
  $MITOSISD init taeguk --chain-id "$MITOSIS_CHAIN_ID" --default-denom "$MITOSIS_DENOM" --home "$MITOSIS_HOME"

  cp "$VAL_KEY_FILE" "$MITOSIS_HOME"/config/priv_validator_key.json
  cp "$MITOSISD_GENESIS_FILE" "$MITOSIS_HOME"/config/genesis.json

  app_toml="$MITOSIS_HOME"/config/app.toml
  config_toml="$MITOSIS_HOME"/config/config.toml

  # Setup app.toml
  sed -i.bak'' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.025thai"/' "$app_toml"
  sed -i.bak'' 's/validator-mode = false/validator-mode = true/' "$app_toml"
  sed -i.bak'' 's/mock = true/mock = false/' "$app_toml"
  sed -i.bak'' 's@endpoint = ""@endpoint = "'"$EXECUTION_ENDPOINT"'"@' "$app_toml"
  sed -i.bak'' 's@jwt-file = ""@jwt-file = "'"$EXECUTION_JWT_FILE"'"@' "$app_toml"

  # Setup config.toml
  # TODO(thai): octane wants nop but ethos don't want nop...
  #sed -i.bak'' 's/type = "flood"/type = "nop"/' "$config_toml"
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
$MITOSISD start --home "$MITOSIS_HOME" \
  $peer_param \
  --p2p.laddr=tcp://0.0.0.0:26656 \
  --rpc.laddr=tcp://0.0.0.0:26657 \
  --grpc.enable \
  --grpc.address=0.0.0.0:9090 \
  --api.enable \
  --api.address=tcp://0.0.0.0:1317 \
  --api.enabled-unsafe-cors \
