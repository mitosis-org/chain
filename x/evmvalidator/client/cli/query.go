package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/ethereum/go-ethereum/common"
	mitotypes "github.com/mitosis-org/chain/types"
	"github.com/mitosis-org/chain/x/evmvalidator/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the query commands for x/evmvalidator
func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "evmvalidator",
		Short:                      "Querying commands for the evmvalidator module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryValidator(),
		GetCmdQueryValidatorByConsAddr(),
		GetCmdQueryValidators(),
		GetCmdQueryWithdrawal(),
		GetCmdQueryWithdrawals(),
		GetCmdQueryWithdrawalsByValidator(),
	)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the current evmvalidator parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(cmd.Context(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidator implements the query validator command.
func GetCmdQueryValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator [address]",
		Short: "Query a validator",
		Long:  `Query details about an individual validator by providing its Ethereum address.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse validator address
			addrStr := args[0]
			if !common.IsHexAddress(addrStr) {
				return fmt.Errorf("invalid validator address format: %s", addrStr)
			}
			valAddr := mitotypes.EthAddress(common.HexToAddress(addrStr))

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Validator(cmd.Context(), &types.QueryValidatorRequest{ValAddr: valAddr.Bytes()})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidatorByConsAddr implements the query validator by consensus address command.
func GetCmdQueryValidatorByConsAddr() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-by-cons [cons-address]",
		Short: "Query a validator by consensus address",
		Long:  `Query details about an individual validator by providing its consensus address.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse consensus address
			consAddr := args[0]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.ValidatorByConsAddr(cmd.Context(), &types.QueryValidatorByConsAddrRequest{ConsAddr: consAddr})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryValidators implements the query all validators command.
func GetCmdQueryValidators() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validators",
		Short: "Query for all validators",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Validators(cmd.Context(), &types.QueryValidatorsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "validators")

	return cmd
}

// GetCmdQueryWithdrawal implements the query withdrawal command.
func GetCmdQueryWithdrawal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawal [id]",
		Short: "Query a withdrawal by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse withdrawal ID
			var id uint64
			if _, err := fmt.Sscanf(args[0], "%d", &id); err != nil {
				return fmt.Errorf("invalid withdrawal ID: %s", args[0])
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Withdrawal(cmd.Context(), &types.QueryWithdrawalRequest{Id: id})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryWithdrawals implements the query all withdrawals command.
func GetCmdQueryWithdrawals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawals",
		Short: "Query for all withdrawals",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Withdrawals(cmd.Context(), &types.QueryWithdrawalsRequest{Pagination: pageReq})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "withdrawals")

	return cmd
}

// GetCmdQueryWithdrawalsByValidator implements the query withdrawals by validator command.
func GetCmdQueryWithdrawalsByValidator() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdrawals-by-validator [validator-address]",
		Short: "Query withdrawals for a specific validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			// Parse validator address
			addrStr := args[0]
			if !common.IsHexAddress(addrStr) {
				return fmt.Errorf("invalid validator address format: %s", addrStr)
			}
			valAddr := mitotypes.EthAddress(common.HexToAddress(addrStr))

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.WithdrawalsByValidator(cmd.Context(), &types.QueryWithdrawalsByValidatorRequest{
				ValAddr:    valAddr.Bytes(),
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "withdrawals")

	return cmd
}
