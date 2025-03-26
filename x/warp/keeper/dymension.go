package keeper

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/bcp-innovations/hyperlane-cosmos/x/warp/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ util.HyperlaneApp = &DymensionHandler{}

type DymensionHandler struct {
	*Keeper
	hook DymHook
}

type DymHookArgs struct {
	MailboxId util.HexAddress

	// original unmdified Message
	Message util.HyperlaneMessage

	// original unmodified Memo
	Memo []byte

	// who was credited
	Account sdk.AccAddress

	// how much was credited
	Coins sdk.Coins
}

type DymHook interface {
	Handle(ctx context.Context, args DymHookArgs) error
}

func NewDymensionHandler(k *Keeper, hook DymHook) *DymensionHandler {
	ret := &DymensionHandler{k, hook}
	ret.RegisterModulesDymension()
	return ret
}

// must be called after new keeper
func (k *DymensionHandler) RegisterModulesDymension() {
	k.GetCoreKeeper().AppRouter().RegisterModule(uint8(types.HYP_TOKEN_TYPE_COLLATERAL_MEMO), k)
	k.GetCoreKeeper().AppRouter().RegisterModule(uint8(types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO), k)
}

// this wasn't exported in upstream
func (k *DymensionHandler) GetCoreKeeper() types.CoreKeeper {
	return k.coreKeeper
}

// originally copied from warp keeper :: Handle
func (k *DymensionHandler) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	token, err := k.HypTokens.Get(ctx, message.Recipient.GetInternalId())
	if err != nil {
		return err
	}

	payloadMemo, err := types.ParseWarpMemoPayload(message.Body)
	if err != nil {
		return err
	}
	payload := payloadMemo.WarpPayload

	if token.OriginMailbox != mailboxId {
		return fmt.Errorf("invalid origin mailbox address")
	}

	remoteRouter, err := k.EnrolledRouters.Get(ctx, collections.Join(message.Recipient.GetInternalId(), message.Origin))
	if err != nil {
		return fmt.Errorf("no enrolled router found for origin %d", message.Origin)
	}

	if message.Sender.String() != strings.ToLower(remoteRouter.ReceiverContract) {
		return fmt.Errorf("invalid receiver contract")
	}

	// Check token type
	err = nil
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL_MEMO {
		if !slices.Contains(k.enabledTokens, int32(types.HYP_TOKEN_TYPE_COLLATERAL_MEMO)) {
			return fmt.Errorf("module disabled collateral tokens")
		}
		err = k.RemoteReceiveCollateral(ctx, token, payload)
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO {
		if !slices.Contains(k.enabledTokens, int32(types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO)) {
			return fmt.Errorf("module disabled synthetic tokens")
		}
		err = k.RemoteReceiveSynthetic(ctx, token, payload)
	} else {
		panic("inconsistent store")
	}

	var account sdk.AccAddress
	var coins sdk.Coins
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL_MEMO {
		account, coins = k.AccountAndCoinsCollat(payload, token)
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO {
		account, coins = k.AccountAndCoinsSynth(payload, token)
	}

	k.hook.Handle(ctx, DymHookArgs{
		MailboxId: mailboxId,
		Message:   message,
		Memo:      payloadMemo.Memo,
		Account:   account,
		Coins:     coins,
	})

	return err
}

func (k *DymensionHandler) AccountAndCoinsCollat(payload types.WarpPayload, token types.HypToken) (sdk.AccAddress, sdk.Coins) {
	account := sdk.AccAddress(payload.Recipient()[12:32])

	amount := math.NewIntFromBigInt(payload.Amount())

	return account, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount))
}

func (k *DymensionHandler) AccountAndCoinsSynth(payload types.WarpPayload, token types.HypToken) (sdk.AccAddress, sdk.Coins) {
	account := payload.GetCosmosAccount()

	amount := math.NewIntFromBigInt(payload.Amount())

	return account, sdk.NewCoins(sdk.NewCoin(token.OriginDenom, amount))
}

// just a slight mod to the upstream method
func (ms msgServer) DymCreateSyntheticToken(ctx context.Context, wrapped *types.MsgDymCreateSyntheticToken) (*types.MsgDymCreateSyntheticTokenResponse, error) {
	msg := wrapped.Inner
	tType := types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO
	if !slices.Contains(ms.k.enabledTokens, int32(tType)) {
		return nil, fmt.Errorf("module disabled synthetic tokens")
	}

	has, err := ms.k.coreKeeper.MailboxIdExists(ctx, msg.OriginMailbox)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", msg.OriginMailbox.String())
	}

	tokenId, err := ms.k.coreKeeper.AppRouter().GetNextSequence(ctx, uint8(tType))
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId,
		Owner:         msg.Owner,
		TokenType:     tType,
		OriginMailbox: msg.OriginMailbox,
		OriginDenom:   fmt.Sprintf("hyperlane/%s", tokenId.String()),
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), newToken); err != nil {
		return nil, err
	}

	return &types.MsgDymCreateSyntheticTokenResponse{Inner: &types.MsgCreateSyntheticTokenResponse{Id: tokenId}}, nil
}

// just a slight mod to the upstream method
func (ms msgServer) DymCreateCollateralToken(ctx context.Context, wrapped *types.MsgDymCreateCollateralToken) (*types.MsgDymCreateCollateralTokenResponse, error) {

	msg := wrapped.Inner
	tType := types.HYP_TOKEN_TYPE_COLLATERAL_MEMO

	if !slices.Contains(ms.k.enabledTokens, int32(tType)) {
		return nil, fmt.Errorf("module disabled collateral tokens")
	}

	err := sdk.ValidateDenom(msg.OriginDenom)
	if err != nil {
		return nil, fmt.Errorf("origin denom %s is invalid", msg.OriginDenom)
	}

	has, err := ms.k.coreKeeper.MailboxIdExists(ctx, msg.OriginMailbox)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("failed to find mailbox with id: %s", msg.OriginMailbox.String())
	}

	tokenId, err := ms.k.coreKeeper.AppRouter().GetNextSequence(ctx, uint8(tType))
	if err != nil {
		return nil, err
	}

	newToken := types.HypToken{
		Id:            tokenId,
		Owner:         msg.Owner,
		TokenType:     tType,
		OriginMailbox: msg.OriginMailbox,
		OriginDenom:   msg.OriginDenom,
	}

	if err = ms.k.HypTokens.Set(ctx, tokenId.GetInternalId(), newToken); err != nil {
		return nil, err
	}
	return &types.MsgDymCreateCollateralTokenResponse{Inner: &types.MsgCreateCollateralTokenResponse{Id: tokenId}}, nil
}

func (ms msgServer) DymRemoteTransfer(ctx context.Context, wrapped *types.MsgDymRemoteTransfer) (*types.MsgDymRemoteTransferResponse, error) {
	goCtx := sdk.UnwrapSDKContext(ctx)
	msg := wrapped.Inner

	token, err := ms.k.HypTokens.Get(ctx, msg.TokenId.GetInternalId())
	if err != nil {
		return nil, fmt.Errorf("failed to find token with id: %s", msg.TokenId.String())
	}

	customHookMetadata, err := util.DecodeEthHex(msg.CustomHookMetadata)
	if err != nil {
		return nil, fmt.Errorf("invalid custom hook metadata")
	}

	var messageResultId util.HexAddress
	if token.TokenType == types.HYP_TOKEN_TYPE_COLLATERAL_MEMO {
		// NOTE: sending Memo from Cosmos not yet supported (not needed for our MVP, since we don't multihop on other chains)
		result, err := ms.k.RemoteTransferCollateral(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.CustomHookId, msg.GasLimit, msg.MaxFee, customHookMetadata)
		if err != nil {
			return nil, err
		}
		messageResultId = result
	} else if token.TokenType == types.HYP_TOKEN_TYPE_SYNTHETIC_MEMO {
		// NOTE: sending Memo from Cosmos not yet supported (not needed for our MVP, since we don't multihop on other chains)
		result, err := ms.k.RemoteTransferSynthetic(goCtx, token, msg.Sender, msg.DestinationDomain, msg.Recipient, msg.Amount, msg.CustomHookId, msg.GasLimit, msg.MaxFee, customHookMetadata)
		if err != nil {
			return nil, err
		}
		messageResultId = result
	} else {
		return nil, errors.New("invalid token type")
	}

	return &types.MsgDymRemoteTransferResponse{
		Inner: &types.MsgRemoteTransferResponse{
			MessageId: messageResultId,
		},
	}, nil
}
