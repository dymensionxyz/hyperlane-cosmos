package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	ismcli "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/client/cli"
	hookcli "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch/client/cli"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group core queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryMailboxes(),
		CmdQueryMailbox(),
		CmdQueryDelivered(),
		CmdQueryRecipientIsm(),
		CmdQueryVerifyDryRun(),
		ismcli.GetQueryCmd(),
		hookcli.GetQueryCmd(),
	)

	return cmd
}

func CmdQueryMailboxes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailboxes",
		Short: "Query all mailboxes",
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

			res, err := queryClient.Mailboxes(cmd.Context(), &types.QueryMailboxesRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "mailboxes")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryMailbox() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailbox [mailbox-id]",
		Short: "Query a mailbox by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Mailbox(cmd.Context(), &types.QueryMailboxRequest{
				Id: mailboxId.String(),
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

func CmdQueryDelivered() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delivered [mailbox-id] [message-id]",
		Short: "Query if a message has been delivered",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			messageId, err := util.DecodeHexAddress(args[1])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.Delivered(cmd.Context(), &types.QueryDeliveredRequest{
				Id:        mailboxId.String(),
				MessageId: messageId.String(),
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

func CmdQueryRecipientIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recipient-ism [app-id]",
		Short: "Query the recipient ISM ID for a registered application",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			appId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.RecipientIsm(cmd.Context(), &types.QueryRecipientIsmRequest{
				Recipient: appId.String(),
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

func CmdQueryVerifyDryRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-dry-run [ism-id] [message] [metadata]",
		Short: "Dry run verification of a message with an ISM",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			ismId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.VerifyDryRun(cmd.Context(), &types.QueryVerifyDryRunRequest{
				IsmId:    ismId.String(),
				Message:  args[1],
				Metadata: args[2],
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
