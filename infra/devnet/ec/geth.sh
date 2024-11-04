#!/bin/sh

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

if [ "$MODE" = "archive" ]; then # archive node
  # --state.scheme=hash: we should use hash scheme when using archive mode. (Later, geth will support path scheme with `--gcmode archive`.)
  STATE_SCHEME="hash"
  # --gcmode archive: we should use archive mode to support the full history of the data.
  GC_MODE_PARAMS="--gcmode archive"
elif [ "$MODE" = "full" ]; then # full node
  STATE_SCHEME="path"
  GC_MODE_PARAMS=""
else
  echo "Invalid MODE: $MODE"
  exit 1
fi

if find "$DATA_DIR" -mindepth 1 -maxdepth 1 | read _; then
  echo "====================================================================="
  echo "[IMPORTANT] Data directory already exists. Skip initialization."
  echo "====================================================================="
else
  echo "[IMPORTANT] Initialize geth data directory: $DATA_DIR"
  geth init \
    --datadir "$DATA_DIR" \
    --db.engine pebble \
    --state.scheme=$STATE_SCHEME \
    "$GENESIS_FILE"

  cp "$NODE_KEY_FILE" "$DATA_DIR/geth/nodekey"
  geth --datadir "$DATA_DIR" dumpconfig > "$DATA_DIR/config.toml"

  STATIC_NODES=$(cat "$PEERS_FILE")
  sed -i.bak'' 's|StaticNodes = \[\]|StaticNodes = \['"$STATIC_NODES"'\]|' "$DATA_DIR/config.toml"
fi

# shellcheck disable=SC2086
# --db.engine pebble: pebble has a better performance than leveldb.
# --syncmode full: we must use full sync mode because there is problem when using snap sync with Octane.
# --miner.recommit=500ms: it should be enough smaller than the block time.
geth --config "$DATA_DIR/config.toml" \
    --datadir "$DATA_DIR" \
    --http \
    --http.addr 0.0.0.0 \
    --http.vhosts "*" \
    --http.corsdomain "*" \
    --http.api eth,net,web3,txpool,rpc,debug \
    --authrpc.addr 0.0.0.0 \
    --authrpc.vhosts "*" \
    --authrpc.jwtsecret "$JWT_FILE" \
    --db.engine pebble \
    --state.scheme=$STATE_SCHEME \
    --syncmode full \
    $GC_MODE_PARAMS \
    --miner.recommit=500ms
