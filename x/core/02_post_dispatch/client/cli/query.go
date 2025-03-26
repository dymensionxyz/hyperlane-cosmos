package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group post dispatch queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "hooks",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "hooks"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryIgps(),
		CmdQueryIgp(),
		CmdQueryDestinationGasConfigs(),
		CmdQueryQuoteGasPayment(),
	)

	return cmd
}

func CmdQueryIgps() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "igps",
		Short: "Query all IGPs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.Igps(cmd.Context(), &types.QueryIgpsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "igps")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryIgp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "igp [igp-id]",
		Short: "Query an IGP by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Igp(cmd.Context(), &types.QueryIgpRequest{
				Id: igpId.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryDestinationGasConfigs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destination-gas-configs [igp-id]",
		Short: "Query destination gas configs for an IGP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.DestinationGasConfigs(cmd.Context(), &types.QueryDestinationGasConfigsRequest{
				Id:         igpId.String(),
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "destination-gas-configs")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryQuoteGasPayment() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote-gas-payment [igp-id] [destination-domain] [gas-amount]",
		Short: "Quote gas payment for a destination",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			igpId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			destinationDomain := args[1]
			gas := args[2]

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QuoteGasPayment(cmd.Context(), &types.QueryQuoteGasPaymentRequest{
				IgpId:             igpId.String(),
				DestinationDomain: destinationDomain,
				GasLimit:          gas,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
