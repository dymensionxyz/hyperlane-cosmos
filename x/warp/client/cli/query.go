package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group warp queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryTokens(),
		CmdQueryToken(),
		CmdQueryBridgedSupply(),
		CmdQueryRemoteRouters(),
		CmdQueryQuoteRemoteTransfer(),
	)

	return cmd
}

func CmdQueryTokens() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tokens",
		Short: "Query all tokens",
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

			res, err := queryClient.Tokens(cmd.Context(), &types.QueryTokensRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "tokens")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token [token-id]",
		Short: "Query a token by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Token(cmd.Context(), &types.QueryTokenRequest{
				Id: tokenId.String(),
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

func CmdQueryBridgedSupply() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bridged-supply [token-id]",
		Short: "Query the bridged supply for a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.BridgedSupply(cmd.Context(), &types.QueryBridgedSupplyRequest{
				Id: tokenId.String(),
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

func CmdQueryRemoteRouters() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remote-routers [token-id]",
		Short: "Query all remote routers for a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			res, err := queryClient.RemoteRouters(cmd.Context(), &types.QueryRemoteRoutersRequest{
				Id:         tokenId.String(),
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "remote-routers")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryQuoteRemoteTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quote-remote-transfer [token-id] [destination-domain]",
		Short: "Quote a remote transfer",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			// Validate that destination domain is a valid uint32
			_, err = strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.QuoteRemoteTransfer(cmd.Context(), &types.QueryQuoteRemoteTransferRequest{
				Id:                tokenId.String(),
				DestinationDomain: args[1],
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
