package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	pdmodule "github.com/bcp-innovations/hyperlane-cosmos/x/core/02_post_dispatch"
	ism "github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/types"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
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
		ism.GetQueryCmd(),
		pdmodule.GetQueryCmd(),
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
				MailboxId: mailboxId,
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
				MailboxId: mailboxId,
				MessageId: messageId,
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
				AppId: appId,
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
