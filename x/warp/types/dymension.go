package types

import (
	"errors"
	fmt "fmt"
	"math/big"
	"slices"
	"strings"

	hyperutil "github.com/bcp-innovations/hyperlane-cosmos/util"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type DymWarpMemoPayload struct {
	WarpPayload
	// see https://docs.hyperlane.xyz/docs/reference/libraries/message for message spec
	// this payload will go in the body field
	Memo []byte
}

func NewWarpMemoPayload(recipient []byte, amount big.Int, memo []byte) (DymWarpMemoPayload, error) {
	wpl, err := NewWarpPayload(recipient, amount)
	if err != nil {
		return DymWarpMemoPayload{}, err
	}

	return DymWarpMemoPayload{WarpPayload: wpl, Memo: memo}, nil
}

func (p DymWarpMemoPayload) Bytes() []byte {
	return slices.Concat(p.WarpPayload.Bytes(), p.Memo)
}

func ParseWarpMemoPayload(payload []byte) (DymWarpMemoPayload, error) {
	if len(payload) < 64 {
		return DymWarpMemoPayload{}, errors.New("payload is too short")
	}

	warp, err := ParseWarpPayload(payload[:64])
	if err != nil {
		return DymWarpMemoPayload{}, err
	}

	memo := payload[64:]

	return DymWarpMemoPayload{
		WarpPayload: warp,
		Memo:        memo,
	}, nil
}

func NewMsgDymRemoteTransfer(memo []byte, inner *MsgRemoteTransfer) *MsgDymRemoteTransfer {
	return &MsgDymRemoteTransfer{
		Signer: inner.Sender,
		Memo:   memo,
		Inner:  inner,
	}
}

func NewMsgDymCreateCollateralToken(inner *MsgCreateCollateralToken) *MsgDymCreateCollateralToken {
	return &MsgDymCreateCollateralToken{
		Signer: inner.Owner,
		Inner:  inner,
	}
}

func NewMsgDymCreateSyntheticToken(inner *MsgCreateSyntheticToken) *MsgDymCreateSyntheticToken {
	return &MsgDymCreateSyntheticToken{
		Signer: inner.Owner,
		Inner:  inner,
	}
}

// addr like dym166kyzqc2e0ewmwmv4vj68pzqp57tgts5lyawlc
// returns a value which can be passed to ethereum as recipient, e.g 0x000000000000000000000000d6ac41030acbf2edbb6cab25a384400d3cb42e14
func HexCosmosAddr(addr string) (string, error) {
	bz, err := sdk.GetFromBech32(addr, sdk.GetConfig().GetBech32AccountAddrPrefix())
	if err != nil {
		return "", fmt.Errorf("addr address from bech32: %w", err)
	}

	ret := hyperutil.EncodeEthHex(bz)
	ret = strings.TrimPrefix(ret, "0x")
	prefix := "0x000000000000000000000000"
	ret = prefix + ret

	return ret, nil
}
