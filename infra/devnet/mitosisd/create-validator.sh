#!/bin/bash

set -e

echo "----- Input Envs -----"
echo "MITOSISD: $MITOSISD"
echo "MIDEVTOOL: $MIDEVTOOL"
echo "MITOSISD_HOME: $MITOSISD_HOME"
echo "MITOSISD_CHAIN_ID: $MITOSISD_CHAIN_ID"
echo "EC_RPC_URL: $EC_RPC_URL"
echo "GOV_ENTRYPOINT: $GOV_ENTRYPOINT"
echo "VAL_ENTRYPOINT: $VAL_ENTRYPOINT"
echo "FUNDER_MNEMONIC (sha256): $(echo -n "$FUNDER_MNEMONIC" | sha256sum)"
echo "----------------------"

FUNDING_AMOUNT=100010ether # (100k + 10) MITO
INITIAL_COLLATERAL_AMOUNT=100000ether # 100k MITO

FUNDER_MNEMONIC_PLAIN=$(echo "$FUNDER_MNEMONIC" | base64 -d)
FUNDER_PRIVKEY=$(cast wallet private-key --mnemonic "$FUNDER_MNEMONIC_PLAIN")

VAL_PRIVKEY_FILE="$MITOSISD_HOME/config/priv_validator_key.json"
VAL_PUBKEY=0x$(jq -r ".pub_key.value" "$VAL_PRIVKEY_FILE" | base64 -d | xxd -p -c 1000 | tr '[:lower:]' '[:upper:]')
VAL_PRIVKEY=0x$(jq -r ".priv_key.value" "$VAL_PRIVKEY_FILE" | base64 -d | xxd -p -c 1000 | tr '[:lower:]' '[:upper:]')
VALIDATOR=$(cast wallet address --private-key "$VAL_PRIVKEY")

echo "Funding the validator: $VALIDATOR"

echo "Sending $FUNDING_AMOUNT from the funder to the validator..."
cast send "$VALIDATOR" \
  --value "$FUNDING_AMOUNT" \
  --mnemonic "$FUNDER_MNEMONIC_PLAIN" \
  --rpc-url "$EC_RPC_URL" \
  --gas-price 1gwei \
  --priority-gas-price 1gwei

sleep 2

echo "Updating the validator entrypoint contract address..."
$MIDEVTOOL governance execute \
  --entrypoint "$GOV_ENTRYPOINT" \
  --private-key "$FUNDER_PRIVKEY" \
  --msg '[{"@type":"/mitosis.evmvalidator.v1.MsgUpdateValidatorEntrypointContractAddr","authority":"mito1g86pactsvfrcglkvqzvdwkxhjshafu280q95p7","addr":"'"$VAL_ENTRYPOINT"'"}]' \
  --rpc "$EC_RPC_URL"

sleep 2

# Get validator pubkey in compressed format

echo "Registering the validator: $VALIDATOR (pubkey: $VAL_PUBKEY)..."
cast send "$VAL_ENTRYPOINT" "registerValidator(address, bytes, address)" \
    "$VALIDATOR" "$VAL_PUBKEY" "$VALIDATOR" \
    --value "$INITIAL_COLLATERAL_AMOUNT" \
    --private-key "$FUNDER_PRIVKEY" \
    --rpc-url "$EC_RPC_URL" \
    --gas-limit 2000000 \
    --gas-price 1gwei \
    --priority-gas-price 1gwei
