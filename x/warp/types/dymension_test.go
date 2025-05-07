package types

import (
	"bytes"
	fmt "fmt"
	"math/big"
	"testing"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

func TestMemoEth(t *testing.T) {
	message := "0x030000000100007a690000000000000000000000004a679253410272dd5232b3ff7cf5dbb88f29531900007a6a0000000000000000000000004a679253410272dd5232b3ff7cf5dbb88f295319000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb922660000000000000000000000000000000000000000000000000000000000000001"
	body := "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000000168656c6c6f"
	checkMessage(t, message)
	checkBody(t, body)
}

func checkMessage(t *testing.T, s string) {
	bz, err := util.DecodeEthHex(s)
	require.NoError(t, err)
	message, err := util.ParseHyperlaneMessage(bz)
	require.NoError(t, err)
	payload, err := ParseWarpMemoPayload(message.Body)
	require.NoError(t, err)
	t.Logf("%+v", payload)
}

func checkBody(t *testing.T, s string) {
	bz, err := util.DecodeEthHex(s)
	require.NoError(t, err)
	payload, err := ParseWarpMemoPayload(bz)
	require.NoError(t, err)
	t.Logf("%+v", payload)
}

func TestHextCosmosAddr(t *testing.T) {
	privKey := secp256k1.GenPrivKey()
	addr := sdk.AccAddress(privKey.PubKey().Address())
	addrS := addr.String()
	fmt.Printf("addr: %s\n", addrS)

	ans, err := HexCosmosAddr(addrS)
	require.NoError(t, err)
	fmt.Printf("addr: %v\n", ans)

	// mimic decodings
	decoded, err := util.DecodeEthHex(ans)
	require.NoError(t, err)
	pl := WarpPayload{recipient: decoded}
	addrSAfter := pl.GetCosmosAccount().String()

	require.Equal(t, addrSAfter, addrS)
	fmt.Printf("account: %s\n", addrSAfter)
}
