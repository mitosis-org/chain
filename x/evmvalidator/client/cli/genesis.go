package cli

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/ethereum/go-ethereum/common"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/omni-network/omni/lib/errors"
	"github.com/spf13/cobra"
)

// parsePubkey parses pubkey from either hex or base64 format
func parsePubkey(pubkeyStr string) ([]byte, error) {
	// Try hex format first (with or without 0x prefix)
	pubkeyStr = strings.TrimPrefix(pubkeyStr, "0x")

	// Check if it looks like hex (only contains hex characters)
	if len(pubkeyStr) == 66 && isHexString(pubkeyStr) { // 33 bytes * 2 = 66 hex chars
		return hex.DecodeString(pubkeyStr)
	}

	// Try base64 format
	pubkey, err := base64.StdEncoding.DecodeString(pubkeyStr)
	if err == nil && len(pubkey) == 33 { // secp256k1 compressed pubkey is 33 bytes
		return pubkey, nil
	}

	// If base64 failed, try hex anyway
	pubkey, hexErr := hex.DecodeString(pubkeyStr)
	if hexErr == nil && len(pubkey) == 33 {
		return pubkey, nil
	}

	return nil, fmt.Errorf("invalid pubkey format: must be 33-byte compressed secp256k1 key in hex or base64 format")
}

// isHexString checks if a string contains only hex characters
func isHexString(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

// GetGenesisValidatorCmd returns a command to add a genesis validator to the evmvalidator module
func GetGenesisValidatorCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-validator [pubkey] [collateral-owner] [collateral] [extra-voting-power] [jailed]",
		Short: "Add a genesis validator to the evmvalidator module",
		Long: `Add a genesis validator to the evmvalidator module.

The pubkey can be provided in either format:
- Hex format: 03a98478cf8213c7fea5a328d89675b5b544fb0c677893690b88473aa3aac0f3ec
- Base64 format (Tendermint): AxB91wLslhi5+SijCPT7Zxmsb04h1me96N3iKRstM3XX

Examples:
$ mitosisd genesis add-validator 03a98478cf8213c7fea5a328d89675b5b544fb0c677893690b88473aa3aac0f3ec 0x0123456789abcdef0123456789abcdef01234567 1000000000000000000 0 false
$ mitosisd genesis add-validator AxB91wLslhi5+SijCPT7Zxmsb04h1me96N3iKRstM3XX 0x0123456789abcdef0123456789abcdef01234567 1000000000000000000 0 false
`,
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			// Parse pubkey (supports both hex and base64 format)
			pubkey, err := parsePubkey(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to parse pubkey")
			}

			valAddr, err := types.PubkeyToEthAddress(pubkey)
			if err != nil {
				return errors.Wrap(err, "failed to convert pubkey to eth address")
			}

			// Parse collateral owner address
			collateralOwner := mitotypes.EthAddress(common.HexToAddress(args[1]))

			// Parse collateral
			collateral := math.NewUintFromString(args[2])

			// Parse extra voting power
			extraVotingPower := math.NewUintFromString(args[3])

			// Parse jailed status
			jailed, err := strconv.ParseBool(args[4])
			if err != nil {
				return errors.Wrap(err, "failed to parse jailed status")
			}

			// Create validator
			validator := types.Validator{
				Addr:             valAddr,
				Pubkey:           pubkey,
				Collateral:       collateral,
				CollateralShares: math.NewUint(0), // not used in InitGenesis
				ExtraVotingPower: extraVotingPower,
				VotingPower:      0, // Will be computed during InitGenesis
				Jailed:           jailed,
				Bonded:           false,
			}

			//nolint:forbidigo
			{
				println("------------ Validator Info ------------")
				println("evm address :", valAddr.String())
				println("acc address :", sdk.AccAddress(valAddr.Bytes()).String())
				println("val address :", sdk.ValAddress(valAddr.Bytes()).String())
				println("cons address:", validator.MustConsAddr().String())
				println()
				println("pubkey (hex)   :", fmt.Sprintf("%X", pubkey))
				println("pubkey (base64):", base64.StdEncoding.EncodeToString(pubkey))
				println("----------------------------------------")
			}

			// Get the server context
			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			// Get genesis
			genFile := config.GenesisFile()
			appState, genDoc, err := genutiltypes.GenesisStateFromGenFile(genFile)
			if err != nil {
				return fmt.Errorf("failed to unmarshal genesis state: %w", err)
			}

			// Get evmvalidator genesis state
			var evmvalGenState types.GenesisState
			if appState[types.ModuleName] != nil {
				cdc.MustUnmarshalJSON(appState[types.ModuleName], &evmvalGenState)
			} else {
				evmvalGenState = *types.DefaultGenesisState()
			}

			// Add validator to genesis state
			evmvalGenState.Validators = append(evmvalGenState.Validators, validator)

			// Create collateral ownership record for the validator
			ownership := types.CollateralOwnership{
				ValAddr:        valAddr,
				Owner:          collateralOwner,
				Shares:         math.NewUint(0), // not used in InitGenesis
				CreationHeight: 0,               // not used in InitGenesis
			}
			evmvalGenState.CollateralOwnerships = append(evmvalGenState.CollateralOwnerships, ownership)

			// Validate the modified genesis state
			if err := evmvalGenState.Validate(); err != nil {
				return fmt.Errorf("failed to validate evmvalidator genesis state: %w", err)
			}

			// Marshal the evmvalidator genesis state
			appState[types.ModuleName] = cdc.MustMarshalJSON(&evmvalGenState)

			// Save the app state to genesis file
			appStateJSON, err := json.MarshalIndent(appState, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal application genesis state: %w", err)
			}

			genDoc.AppState = appStateJSON
			return genutil.ExportGenesisFile(genDoc, genFile)
		},
	}

	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The application home directory")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
