package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/core/01_interchain_security/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group interchain security queries under a subcommand
	cmd := &cobra.Command{
		Use:                        "ism",
		Short:                      fmt.Sprintf("Querying commands for the %s module", "ism"),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		CmdQueryIsms(),
		CmdQueryIsm(),
		CmdQueryAnnouncedStorageLocations(),
		CmdQueryLatestAnnouncedStorageLocation(),
	)

	return cmd
}

func CmdQueryIsms() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "isms",
		Short: "Query all ISMs",
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

			res, err := queryClient.Isms(cmd.Context(), &types.QueryIsmsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "isms")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryIsm() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ism [ism-id]",
		Short: "Query an ISM by ID",
		Args:  cobra.ExactArgs(1),
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
			res, err := queryClient.Ism(cmd.Context(), &types.QueryIsmRequest{
				Id: ismId.String(),
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

func CmdQueryAnnouncedStorageLocations() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "announced-storage-locations [mailbox-id] [validator-address]",
		Short: "Query announced storage locations for an ISM",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			mailbox, err := util.DecodeHexAddress(args[0]) // TODO: definitely hex?
			if err != nil {
				return err
			}

			validator, err := util.DecodeHexAddress(args[1]) // TODO: definitely hex?
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.AnnouncedStorageLocations(cmd.Context(), &types.QueryAnnouncedStorageLocationsRequest{
				MailboxId:        mailbox.String(),
				ValidatorAddress: validator.String(),
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, "announced-storage-locations")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdQueryLatestAnnouncedStorageLocation() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-announced-storage-location [ism-id]",
		Short: "Query the latest announced storage location for an ISM",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			mailbox, err := util.DecodeHexAddress(args[0]) // TODO: definitely hex?
			if err != nil {
				return err
			}

			validator, err := util.DecodeHexAddress(args[1]) // TODO: definitely hex?
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.LatestAnnouncedStorageLocation(cmd.Context(), &types.QueryLatestAnnouncedStorageLocationRequest{
				MailboxId:        mailbox.String(),
				ValidatorAddress: validator.String(),
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
