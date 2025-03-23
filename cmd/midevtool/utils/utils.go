package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetEthClient creates and returns an Ethereum client
func GetEthClient(rpcURL string) (*ethclient.Client, error) {
	return ethclient.Dial(rpcURL)
}

// GetPrivateKey converts a hex string to an ECDSA private key
func GetPrivateKey(key string) *ecdsa.PrivateKey {
	// Remove 0x prefix if present
	key = strings.TrimPrefix(key, "0x")

	privKey := ethcrypto.ToECDSAUnsafe(common.Hex2Bytes(key))
	return privKey
}

// CreateTransactOpts creates transaction options for a contract call
func CreateTransactOpts(client *ethclient.Client, privKey *ecdsa.PrivateKey, value *big.Int) *bind.TransactOpts {
	addr := common.BytesToAddress(ethcrypto.PubkeyToAddress(privKey.PublicKey).Bytes())

	nonce, err := client.PendingNonceAt(context.Background(), addr)
	if err != nil {
		panic(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(privKey, chainID)
	if err != nil {
		panic(err)
	}

	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.Value = value

	return opts
}

// ParseValue parses a string into a big.Int value
func ParseValue(value string) (*big.Int, error) {
	valueInt := big.NewInt(0)
	if value != "" {
		var success bool
		valueInt, success = big.NewInt(0).SetString(value, 10)
		if !success {
			return nil, fmt.Errorf("invalid format")
		}
	}
	return valueInt, nil
}
