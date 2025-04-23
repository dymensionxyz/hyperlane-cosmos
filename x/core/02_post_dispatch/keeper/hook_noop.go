package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dymensionxyz/hyperlane-cosmos/util"
	"github.com/dymensionxyz/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

type NoopHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = NoopHookHandler{}

func (i NoopHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := i.k.noopHooks.Has(ctx, hookId.GetInternalId())
	if err != nil || !has {
		return false, errors.Wrapf(types.ErrHookDoesNotExistOrIsNotRegistered, "%s", hookId.String())
	}
	return has, nil
}

func (i NoopHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_UNUSED
}

func (i NoopHookHandler) PostDispatch(ctx context.Context, _, hookId util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage, _ sdk.Coins) (sdk.Coins, error) {
	has, err := i.k.noopHooks.Has(ctx, hookId.GetInternalId())
	if err != nil || !has {
		return nil, errors.Wrapf(types.ErrHookDoesNotExistOrIsNotRegistered, "%s", hookId.String())
	}

	return sdk.NewCoins(), nil
}

func (i NoopHookHandler) QuoteDispatch(_ context.Context, _, _ util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage) (sdk.Coins, error) {
	return sdk.NewCoins(), nil
}
