#!/bin/bash

set -e
set -x

echo "----- Input Envs -----"
INFRA_DIR=${INFRA_DIR:-'./infra/devnet'}
DB_DIR=${DB_DIR:-'./tmp/devnet'}
ISOLATED_CONFIG_DIR=${ISOLATED_CONFIG_DIR:-'./infra/devnet/geth/config/val-1'}

echo "INFRA_DIR: $INFRA_DIR"
echo "DB_DIR: $DB_DIR"
echo "ISOLATED_CONFIG_DIR: $ISOLATED_CONFIG_DIR"
echo "----------------------"

COMMON_CONFIG_DIR="$INFRA_DIR/geth/config/common"
DATA_DIR="$DB_DIR/geth"

if [ -d "$DATA_DIR" ]; then
  echo "====================================================================="
  echo "[IMPORTANT] Data directory already exists. Skip initialization."
  echo "====================================================================="
else
  echo "[IMPORTANT] Initialize geth data directory: $DATA_DIR"
  geth init \
    --datadir "$DATA_DIR" \
    --db.engine pebble \
    --state.scheme=hash \
    "$COMMON_CONFIG_DIR/genesis.json"

  cp "$ISOLATED_CONFIG_DIR/nodekey" "$DATA_DIR/geth/nodekey"
  geth --datadir "$DATA_DIR" dumpconfig > "$DATA_DIR/config.toml"

  STATIC_NODES=$(cat "$ISOLATED_CONFIG_DIR/peers")
  sed -i.bak'' 's|StaticNodes = \[\]|StaticNodes = \['"$STATIC_NODES"'\]|' "$DATA_DIR/config.toml"
fi

# --db.engine pebble: pebble has a better performance than leveldb.
# --state.scheme=hash: we should use hash scheme when using archive mode. (Later, geth will support path scheme with `--gcmode archive`.)
# --syncmode full: we must use full sync mode because there is problem when using snap sync with Octane.
# --gcmode archive: we should use archive mode to support the full history of the data.
# --miner.recommit=500ms: it should be enough smaller than the block time.
geth --config "$DATA_DIR/config.toml" \
    --http \
    --http.addr 0.0.0.0 \
    --http.vhosts "*" \
    --http.corsdomain "*" \
    --http.api eth,net,web3,txpool,debug \
    --authrpc.addr 0.0.0.0 \
    --authrpc.jwtsecret "$COMMON_CONFIG_DIR/jwt.hex" \
    --authrpc.vhosts "*" \
    --datadir "$DATA_DIR" \
    --db.engine pebble \
    --state.scheme=hash \
    --syncmode full \
    --gcmode archive \
    --miner.recommit=500ms
