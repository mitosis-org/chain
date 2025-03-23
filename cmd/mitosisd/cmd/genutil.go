package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cometbft/cometbft/crypto/secp256k1"

	cfg "github.com/cometbft/cometbft/config"
	cmtos "github.com/cometbft/cometbft/libs/os"
	"github.com/cometbft/cometbft/p2p"
	"github.com/cometbft/cometbft/privval"
	"github.com/cosmos/go-bip39"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

// InitializeNodeValidatorFiles creates private validator and p2p configuration files.
func InitializeNodeValidatorFiles(config *cfg.Config) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	return InitializeNodeValidatorFilesFromMnemonic(config, "")
}

// InitializeNodeValidatorFilesFromMnemonic creates private validator and p2p configuration files using the given mnemonic.
// If no valid mnemonic is given, a random one will be used instead.
func InitializeNodeValidatorFilesFromMnemonic(config *cfg.Config, mnemonic string) (nodeID string, valPubKey cryptotypes.PubKey, err error) {
	if len(mnemonic) > 0 && !bip39.IsMnemonicValid(mnemonic) {
		return "", nil, fmt.Errorf("invalid mnemonic")
	}
	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return "", nil, err
	}

	nodeID = string(nodeKey.ID())

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := os.MkdirAll(filepath.Dir(pvKeyFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvKeyFile), err)
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := os.MkdirAll(filepath.Dir(pvStateFile), 0o777); err != nil {
		return "", nil, fmt.Errorf("could not create directory %q: %w", filepath.Dir(pvStateFile), err)
	}

	var filePV *privval.FilePV
	if len(mnemonic) == 0 {
		if cmtos.FileExists(pvKeyFile) {
			filePV = privval.LoadFilePV(pvKeyFile, pvStateFile)
		} else {
			privKey := secp256k1.GenPrivKey()
			filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
			filePV.Save()
		}
	} else {
		privKey := secp256k1.GenPrivKeySecp256k1([]byte(mnemonic))
		filePV = privval.NewFilePV(privKey, pvKeyFile, pvStateFile)
		filePV.Save()
	}

	tmValPubKey, err := filePV.GetPubKey()
	if err != nil {
		return "", nil, err
	}

	valPubKey, err = cryptocodec.FromCmtPubKeyInterface(tmValPubKey)
	if err != nil {
		return "", nil, err
	}

	return nodeID, valPubKey, nil
}
