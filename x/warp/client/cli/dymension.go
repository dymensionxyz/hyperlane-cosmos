package cli

import (
	"errors"
	"fmt"
	"strconv"

	"cosmossdk.io/math"
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

func CmdDymRemoteTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dym-transfer [token-id] [destination-domain] [recipient] [amount]",
		Short: "Send Hyperlane Token",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			tokenId, err := util.DecodeHexAddress(args[0])
			if err != nil {
				return err
			}

			destinationDomain, err := strconv.ParseUint(args[1], 10, 32)
			if err != nil {
				return err
			}

			recipient, err := util.DecodeHexAddress(args[2])
			if err != nil {
				return err
			}

			argAmount, ok := math.NewIntFromString(args[3])
			if !ok {
				return errors.New("invalid amount")
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			gasLimitInt, ok := math.NewIntFromString(gasLimit)
			if !ok {
				return errors.New("failed to convert `gasLimit` into math.Int")
			}

			maxFeeCoin, err := sdk.ParseCoinNormalized(maxFee)
			if err != nil {
				return err
			}

			var parsedHookId *util.HexAddress = nil
			if customHookId != "" {
				parsed, err := util.DecodeHexAddress(customHookId)
				if err != nil {
					return err
				}
				parsedHookId = &parsed
			}

			msg := types.NewMsgDymRemoteTransfer(nil, &types.MsgRemoteTransfer{
				TokenId:            tokenId,
				DestinationDomain:  uint32(destinationDomain),
				Sender:             clientCtx.GetFromAddress().String(),
				Recipient:          recipient,
				Amount:             argAmount,
				CustomHookId:       parsedHookId,
				GasLimit:           gasLimitInt,
				MaxFee:             maxFeeCoin,
				CustomHookMetadata: customHookMetadata,
			})

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	cmd.Flags().StringVar(&customHookId, "custom-hook-id", "", "custom DefaultHookId")
	cmd.Flags().StringVar(&customHookMetadata, "custom-hook-metadata", "", "custom hook metadata")

	cmd.Flags().StringVar(&gasLimit, "gas-limit", "0", "Overwrite InterchainGasPayment gas limit")

	cmd.Flags().StringVar(&maxFee, "max-hyperlane-fee", "0", "maximum Hyperlane InterchainGasPayment")

	return cmd
}
