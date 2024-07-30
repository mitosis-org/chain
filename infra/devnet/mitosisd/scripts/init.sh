#!/bin/bash

set -e
set -x

echo "----- Input Envs -----"
INFRA_DIR=${INFRA_DIR:-'./infra/devnet'}
DB_DIR=${DB_DIR:-'./tmp/devnet'}

MITOSIS_CHAIN_ID=${MITOSIS_CHAIN_ID:-'mitosis-devnet-1'}
MITOSISD=${MITOSISD:-'./build/mitosisd'}
ETHOSD=${ETHOSD:-'./ethos/ethos-chain/build/ethosd'}
ETHOS_RPC=${ETHOS_RPC:-'http://127.0.0.1:26657'}

echo "INFRA_DIR: $INFRA_DIR"
echo "DB_DIR: $DB_DIR"
echo "MITOSIS_CHAIN_ID: $MITOSIS_CHAIN_ID"
echo "MITOSISD: $MITOSISD"
echo "ETHOSD: $ETHOSD"
echo "ETHOS_RPC: $ETHOS_RPC"
echo "----------------------"

# mitosis variables
MITOSIS_HOME="$DB_DIR/mitosisd"
MITOSIS_DENOM='thai'
MITOSIS_GEN_ACC_MNEMONIC='Ymxvb2QgYW1hemluZyBwYXNzIGxpbWl0IGFsbG93IGZsb29yIGZpbmFsIHByYWN0aWNlIGNoaWVmIHRyaWFsIG9ibGlnZSBob29kIGRyaXAgcGFsbSBwcm9ncmFtIGZsdXNoIG1pbGxpb24gZm9sZCBvcmFuZ2UgZGFyaW5nIHN3YXAgZmx5IHJlc2N1ZSBsaW1i'

EXECUTION_GENESIS_BLOCK_HASH=$(xargs < "$INFRA_DIR/geth/config/common/genesis-block-hash.txt") # trim whitespace

# ethos variables
ETHOS_HOME="$DB_DIR/ethosd"
ETHOS_CHAIN_ID='ethos'
ETHOS_DENOM='stake'
ETHOS_ADMIN_MNEMONIC='Ymxvb2QgYW1hemluZyBwYXNzIGxpbWl0IGFsbG93IGZsb29yIGZpbmFsIHByYWN0aWNlIGNoaWVmIHRyaWFsIG9ibGlnZSBob29kIGRyaXAgcGFsbSBwcm9ncmFtIGZsdXNoIG1pbGxpb24gZm9sZCBvcmFuZ2UgZGFyaW5nIHN3YXAgZmx5IHJlc2N1ZSBsaW1i'

# output variables
ARTIFACTS_PATH="$DB_DIR/artifacts"

mitosis_initial_genesis_file="$MITOSIS_HOME/config/genesis.json"
mitosis_final_genesis_file="$MITOSIS_HOME/config/genesis-final.json"

function create_initial_genesis_for_mitosis() {
  # Generate genesis file for mitosisd
  $MITOSISD init taeguk --chain-id "$MITOSIS_CHAIN_ID" --default-denom "$MITOSIS_DENOM" --home "$MITOSIS_HOME"

  # Setup genesis account
  echo "$MITOSIS_GEN_ACC_MNEMONIC" | base64 -d | $MITOSISD keys add gen --recover --keyring-backend test --home "$MITOSIS_HOME"
  $MITOSISD genesis add-genesis-account gen 1000000000000000000000000"$MITOSIS_DENOM" --keyring-backend test --home "$MITOSIS_HOME"

  # Setup execution block hash
  hash=$(echo -n "$EXECUTION_GENESIS_BLOCK_HASH" | xxd -r -p | base64)
  jq --arg hash "$hash" '.app_state.evmengine.execution_block_hash = $hash' "$mitosis_initial_genesis_file" > tmp.json && mv tmp.json "$mitosis_initial_genesis_file"

  # Setup additional genesis modifications
  jq '.consensus.params.block.max_bytes = "-1"' "$mitosis_initial_genesis_file" > tmp.json && mv tmp.json "$mitosis_initial_genesis_file"

  echo "--- initial genesis for mitosis ---"
  cat "$mitosis_initial_genesis_file"
  echo "-----------------------------------"
}

function add_mitosis_chain_to_ethos() {
  genesis_hash=$(sha256sum "$mitosis_initial_genesis_file" | awk '{ print $1 }')
  binary_hash=$(sha256sum "$MITOSISD" | awk '{ print $1 }')
  # test gdate command exists
  if command -v gdate &> /dev/null; then
    # for macos or unix systems
    spawn_time=$(gdate -u --date 'now + 1 minutes' +"%Y-%m-%dT%H:%M:%SZ")
  else
    # for general linux systems
    spawn_time=$(date -u -D +1m +"%Y-%m-%dT%H:%M:%SZ")
  fi

  cat > proposal.json <<EOF
    {
        "title": "Add mitosis chain as consumer",
        "summary": "Add mitosis chain as consumer",
        "chain_id": "$MITOSIS_CHAIN_ID",
        "initial_height": {
            "revision_number": 1,
            "revision_height": 1
        },
        "genesis_hash": "$genesis_hash",
        "binary_hash": "$binary_hash",
        "spawn_time": "$spawn_time",
        "consumer_redistribution_fraction": "0.5",
        "blocks_per_distribution_transmission": 1000,
        "distribution_transmission_channel": "channel-1",
        "historical_entries": 500,
        "ccv_timeout_period": 2419200000000000,
        "transfer_timeout_period": 1800000000000,
        "unbonding_period": 1728000000000000,
        "deposit": "1000$PROVIDER_DENOM"
    }
EOF

  echo "--- proposal.json ---"
  cat proposal.json
  echo "---------------------"

  echo "Submitting proposal to ethos chain..."
  echo "$ETHOS_ADMIN_MNEMONIC" | base64 -d | $ETHOSD keys add admin --recover --keyring-backend test --home "$ETHOS_HOME"
  $ETHOSD tx provider add-consumer-chain proposal.json --from admin --chain-id "$ETHOS_CHAIN_ID" --keyring-backend test --home "$ETHOS_HOME" --node "$ETHOS_RPC" --gas-prices 1"$ETHOS_DENOM" -y

  echo "Waiting for proposal to be accepted..."
  sleep 10
}

function create_final_genesis_for_mitosis() {
  $ETHOSD query provider consumer-genesis "$MITOSIS_CHAIN_ID" -o json --node "$ETHOS_RPC" > ccv-state.json

  echo "--- ccv-state.json ---"
  cat ccv-state.json
  echo "----------------------"

  cp "$mitosis_initial_genesis_file" "$mitosis_final_genesis_file"
  jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' "$mitosis_final_genesis_file" ccv-state.json > tmp.json && mv tmp.json "$mitosis_final_genesis_file"

  echo "--- final genesis for mitosis ---"
  cat "$mitosis_final_genesis_file"
  echo "---------------------------------"
}

echo "Creating initial genesis for mitosis..."
create_initial_genesis_for_mitosis

echo "Adding mitosis chain to ethos chain..."
add_mitosis_chain_to_ethos

echo "Creating final genesis for mitosis..."
create_final_genesis_for_mitosis

echo "Copying artifacts to $ARTIFACTS_PATH..."
mkdir -p "$ARTIFACTS_PATH"
cp "$mitosis_initial_genesis_file" "$ARTIFACTS_PATH/genesis-initial.json"
cp "$mitosis_final_genesis_file" "$ARTIFACTS_PATH/genesis-final.json"
cp proposal.json "$ARTIFACTS_PATH/proposal.json"
