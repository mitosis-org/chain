package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/omni-network/omni/lib/errors"
	"strconv"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/spf13/cobra"
)

// GetGenesisValidatorCmd returns a command to add a genesis validator to the evmvalidator module
func GetGenesisValidatorCmd(defaultNodeHome string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-genesis-validator [pubkey] [collateral] [extra-voting-power] [jailed]",
		Short: "Add a genesis validator to the evmvalidator module",
		Long: `Add a genesis validator to the evmvalidator module.
Example:
$ mitosisd add-genesis-validator 03a98478cf8213c7fea5a328d89675b5b544fb0c677893690b88473aa3aac0f3ec 1000000000000000000 0 false
`,
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			cdc := clientCtx.Codec

			// Parse pubkey (hex format)
			pubkeyHex := args[0]
			if len(pubkeyHex) < 2 {
				return fmt.Errorf("pubkey too short")
			}

			// Remove 0x prefix if present
			if pubkeyHex[:2] == "0x" {
				pubkeyHex = pubkeyHex[2:]
			}

			pubkey, err := hex.DecodeString(pubkeyHex)
			if err != nil {
				return errors.Wrap(err, "failed to decode pubkey string")
			}

			valAddr, err := types.PubkeyToEthAddress(pubkey)
			if err != nil {
				return errors.Wrap(err, "failed to convert pubkey to eth address")
			}

			// Parse collateral
			collateral, ok := math.NewIntFromString(args[1])
			if !ok {
				return fmt.Errorf("invalid collateral amount: %s", args[1])
			}

			// Parse extra voting power
			extraVotingPower, ok := math.NewIntFromString(args[2])
			if !ok {
				return fmt.Errorf("invalid extra voting power: %s", args[2])
			}

			// Parse jailed status
			jailed, err := strconv.ParseBool(args[3])
			if err != nil {
				return errors.Wrap(err, "failed to parse jailed status")
			}

			// Create validator
			validator := types.Validator{
				Addr:             valAddr,
				Pubkey:           pubkey,
				Collateral:       collateral,
				ExtraVotingPower: extraVotingPower,
				VotingPower:      math.ZeroInt(), // Will be computed during InitGenesis
				Jailed:           jailed,
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
