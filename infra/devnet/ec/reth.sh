#!/bin/bash

set -e

echo "----- Input Envs -----"
# mandatory
echo "MODE: $MODE" # archive, full
echo "DATA_DIR: $DATA_DIR"
echo "JWT_FILE: $JWT_FILE"

# only for initialization
echo "GENESIS_FILE: $GENESIS_FILE"
echo "NODE_KEY_FILE: $NODE_KEY_FILE"
echo "PEERS_FILE: $PEERS_FILE"
echo "----------------------"

if [ "$MODE" == "archive" ]; then # archive node
  FULL_NODE_PARAMS=""
elif [ "$MODE" == "full" ]; then # full node
  FULL_NODE_PARAMS="--full"
else
  echo "Invalid MODE: $MODE"
  exit 1
fi

if find "$DATA_DIR" -mindepth 1 -maxdepth 1 | read _; then
  echo "====================================================================="
  echo "[IMPORTANT] Data directory already exists. Skip initialization."
  echo "====================================================================="
else
  echo "[IMPORTANT] Initialize reth data directory: $DATA_DIR"
  reth init \
    --datadir "$DATA_DIR" \
    --chain "$GENESIS_FILE"

  cp "$NODE_KEY_FILE" "$DATA_DIR/discovery-secret"

  trusted_nodes=$(cat "$PEERS_FILE")
  sed -i.bak'' 's|trusted_nodes = \[\]|trusted_nodes = \['"$trusted_nodes"'\]|' "$DATA_DIR/reth.toml"
fi

# Workaround to resolve the issue that a peer count is zero at startup.
# Just wait a few seconds for docker containers for peers to be started
sleep 3

# shellcheck disable=SC2086
reth node \
    --datadir "$DATA_DIR" \
    --chain "$GENESIS_FILE" \
    --http \
    --http.addr 0.0.0.0 \
    --http.api eth,net,web3,txpool,rpc,debug,trace,admin \
    --authrpc.addr 0.0.0.0 \
    --authrpc.jwtsecret "$JWT_FILE" \
    --metrics 0.0.0.0:9001 \
    --builder.interval 30ms \
    --builder.deadline 1 \
    --engine.legacy \
    $FULL_NODE_PARAMS
