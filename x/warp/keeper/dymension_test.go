package keeper_test

import (
	"context"
	"math/big"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	i "github.com/dymensionxyz/hyperlane-cosmos/tests/integration"
	"github.com/dymensionxyz/hyperlane-cosmos/util"
	coreTypes "github.com/dymensionxyz/hyperlane-cosmos/x/core/types"
	"github.com/dymensionxyz/hyperlane-cosmos/x/warp/keeper"
	"github.com/dymensionxyz/hyperlane-cosmos/x/warp/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("dymension.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var owner i.TestValidatorAddress
	var sender i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChainWithEnabledTokens([]int32{int32(types.HYP_TOKEN_TYPE_COLLATERAL_MEMO)})
		owner = i.GenerateTestValidatorAddress("Owner")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(owner.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("MsgRemoteTransfer && MsgRemoteReceiveCollateral (Memo) (dummy hook) (valid) (Collateral)", func() {

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ DYMENSION CHANGE ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// Need to do this at the start to register the app router modules
		h := &DummyHook{}
		dymHandler := keeper.NewDymensionHandler(&s.App().WarpKeeper)
		dymHandler.SetHook(h)
		dymHandler.RegisterDymensionTokens()

		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, mailboxId, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL_MEMO)
		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)
		// Act
		_, err = s.RunTx(&types.MsgDymRemoteTransfer{
			Inner: &types.MsgRemoteTransfer{
				Sender:            sender.Address,
				TokenId:           tokenId,
				DestinationDomain: remoteRouter.ReceiverDomain,
				Recipient:         receiverAddress,
				Amount:            amount,
				CustomHookId:      &igpId,
				GasLimit:          math.ZeroInt(),
				MaxFee:            maxFee,
			},
			Memo: []byte("test memo"),
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Sub(amount.Add(maxFee.Amount))))

		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		memo := []byte("test memo")

		warpPayload, err := types.NewWarpMemoPayload(warpRecipient, *big.NewInt(amount.Int64()), memo)
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     coreTypes.MESSAGE_VERSION,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		senderBalance = s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ DYMENSION CHANGE ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// don't expect the balance to have been transferred, but do ensure the hook was called
		Expect(h.called).To(BeTrue())
		Expect(h.gotMemo).To(Equal(memo))
	})

	It("MsgRemoteTransfer && MsgRemoteReceiveCollateral (Memo) (default hook) (valid) (Collateral)", func() {

		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~ DYMENSION CHANGE ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
		// Need to do this at the start to register the app router modules
		dymHandler := keeper.NewDymensionHandler(&s.App().WarpKeeper)
		h := &keeper.DymDefaultHook{
			DymensionHandler: dymHandler,
		}
		dymHandler.SetHook(h)
		dymHandler.RegisterDymensionTokens()

		// Arrange
		receiverAddress, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		remoteRouter := types.RemoteRouter{
			ReceiverDomain:   1,
			ReceiverContract: "0x934b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591958e865def0",
			Gas:              math.NewInt(50000),
		}

		amount := math.NewInt(100)
		maxFee := sdk.NewCoin(denom, math.NewInt(250000))

		tokenId, mailboxId, igpId, _ := createToken(s, &remoteRouter, owner.Address, sender.Address, types.HYP_TOKEN_TYPE_COLLATERAL_MEMO)
		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderBalance := s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)
		// Act
		_, err = s.RunTx(&types.MsgDymRemoteTransfer{
			Inner: &types.MsgRemoteTransfer{
				Sender:            sender.Address,
				TokenId:           tokenId,
				DestinationDomain: remoteRouter.ReceiverDomain,
				Recipient:         receiverAddress,
				Amount:            amount,
				CustomHookId:      &igpId,
				GasLimit:          math.ZeroInt(),
				MaxFee:            maxFee,
			},
			Memo: []byte("test memo"),
		})
		Expect(err).To(BeNil())

		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Sub(amount.Add(maxFee.Amount))))

		receiverContract, err := util.DecodeHexAddress(remoteRouter.ReceiverContract)
		Expect(err).To(BeNil())

		warpRecipient, err := sdk.GetFromBech32(sender.Address, "hyp")
		Expect(err).To(BeNil())

		memo := []byte("test memo")

		warpPayload, err := types.NewWarpMemoPayload(warpRecipient, *big.NewInt(amount.Int64()), memo)
		Expect(err).To(BeNil())

		message := util.HyperlaneMessage{
			Version:     3,
			Nonce:       1,
			Origin:      remoteRouter.ReceiverDomain,
			Sender:      receiverContract,
			Destination: 0,
			Recipient:   tokenId,
			Body:        warpPayload.Bytes(),
		}

		senderBalance = s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom)

		_, err = s.RunTx(&coreTypes.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())
		Expect(s.App().BankKeeper.GetBalance(s.Ctx(), sender.AccAddress, denom).Amount).To(Equal(senderBalance.Amount.Add(amount)))
	})
})

type DummyHook struct {
	called  bool
	gotMemo []byte
}

func (h *DummyHook) OnHyperlaneMessage(ctx context.Context, args keeper.OnHyperlaneMessageArgs) error {
	h.called = true
	h.gotMemo = args.Memo
	return nil
}

func DymGetAltnernateTokenID(ctx sdk.Context, coreKeeper types.CoreKeeper, tokenType uint8) (util.HexAddress, error) {
	return coreKeeper.AppRouter().GetNextSequence(ctx, uint8(tokenType))
}
