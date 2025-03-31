package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
)

func CmdDymCreateCollateralToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dym-create-collateral-token [origin-mailbox] [origin-denom]",
		Short: "Create a Hyperlane Collateral Token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgDymCreateCollateralToken{
				Signer: clientCtx.GetFromAddress().String(),
				Inner: &types.MsgCreateCollateralToken{
					Owner:         clientCtx.GetFromAddress().String(),
					OriginMailbox: mailboxId,
					OriginDenom:   args[1],
				},
			}

			_, err = sdk.AccAddressFromBech32(msg.Signer)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Signer))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&ismId, "ism-id", "", "ISM ID; if not specified, DefaultISM is used")

	return cmd
}

func CmdDymCreateSyntheticToken() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dym-create-synthetic-token [origin-mailbox]",
		Short: "Create a Hyperlane Synthetic Token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			mailboxId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			msg := types.MsgDymCreateSyntheticToken{
				Signer: clientCtx.GetFromAddress().String(),
				Inner: &types.MsgCreateSyntheticToken{
					Owner:         clientCtx.GetFromAddress().String(),
					OriginMailbox: mailboxId,
				},
			}

			_, err = sdk.AccAddressFromBech32(msg.Signer)
			if err != nil {
				panic(fmt.Errorf("invalid owner address (%s)", msg.Signer))
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&ismId, "ism-id", "", "ISM ID; if not specified, DefaultISM is used")

	return cmd
}
