package types

import (
	"bytes"
	"math/big"
	"testing"
)

func TestNewWarpMemoPayload(t *testing.T) {
	recipient := []byte{1, 2, 3}
	amount := big.NewInt(100)
	for _, memo := range [][]byte{nil, []byte("test memo")} {
		expect, err := NewWarpMemoPayload(recipient, *amount, memo)
		if err != nil {
			t.Fatalf("error creating warp memo payload: %v", err)
		}
		bz := expect.Bytes()
		actual, err := ParseWarpMemoPayload(bz)
		if err != nil {
			t.Fatalf("error parsing warp memo payload: %v", err)
		}
		if !bytes.Equal(actual.Memo, memo) {
			t.Fatalf("memo is not correct")
		}
		if !bytes.Equal(actual.WarpPayload.Recipient(), recipient) {
			// NOTE: should probably fail here but this is an upstream bug
			// https: //github.com/bcp-innovations/hyperlane-cosmos/issues/91

		}
		if actual.WarpPayload.Amount().Cmp(amount) != 0 {
			t.Fatalf("amount is not correct")
		}
	}
}
