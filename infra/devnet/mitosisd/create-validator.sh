#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "MITOSISD_NODE_RPC: $MITOSISD_NODE_RPC"
echo "VALIDATOR_JSON_FILE: $VALIDATOR_JSON_FILE"

# optional if it needs to be funded
echo "NEED_FUNDING: $NEED_FUNDING" # if it is "yes", it will be funded by the funder
echo "FUNDER_MNEMONIC (sha256): $(echo -n "$FUNDER_MNEMONIC" | sha256sum)"
echo "FUNDING_AMOUNT: $FUNDING_AMOUNT" # e.g. 150000ustake
echo "----------------------"

VALIDATOR=$(mitosisd keys show validator --address --keyring-backend test --home "$MITOSISD_HOME")

if [ "$NEED_FUNDING" == "yes" ]; then
  echo "Funding the validator: $VALIDATOR"

  if $MITOSISD keys show funder --keyring-backend test --home "$MITOSISD_HOME" 2>/dev/null; then
    echo "The funder key already exists."
  else
    echo "Creating the funder key..."
    echo "$FUNDER_MNEMONIC" | base64 -d | $MITOSISD keys add funder --recover --algo "secp256k1" --keyring-backend test --home "$MITOSISD_HOME"
  fi

  echo "Sending $FUNDING_AMOUNT from the funder to the validator..."
  $MITOSISD --home "$MITOSISD_HOME" --keyring-backend test --chain-id "$MITOSISD_CHAIN_ID" \
    --fees 1000ustake --node "$MITOSISD_NODE_RPC" \
    tx bank send funder "$VALIDATOR" "$FUNDING_AMOUNT" --yes

  # wait for the transaction to be committed
  sleep 8
fi

echo "Creating the validator: $VALIDATOR"
$MITOSISD --home "$MITOSISD_HOME" --keyring-backend test --chain-id "$MITOSISD_CHAIN_ID" \
    --fees 1000ustake --node "$MITOSISD_NODE_RPC" \
    tx staking create-validator "$VALIDATOR_JSON_FILE" --from validator --yes
