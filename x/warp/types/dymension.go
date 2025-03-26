package types

import (
	"errors"
	"math/big"
	"slices"
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

	// ~~~ GPT ~~~
	// Trim leading zeros from recipient bytes
	// recipient := warp.Recipient()
	// for len(recipient) > 0 && recipient[0] == 0 {
	// 	recipient = recipient[1:]
	// }

	// warp.recipient = recipient
	// ~~~~~~~~~~~
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
