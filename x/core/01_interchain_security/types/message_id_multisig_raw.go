package types

import (
	"context"
	"fmt"
	"slices"

	"github.com/bcp-innovations/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &MessageIdMultisigISMRaw{}

func (m *MessageIdMultisigISMRaw) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

func (m *MessageIdMultisigISMRaw) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG_RAW
}

// both rawMetadata and message are UNTRUSTED
func (m *MessageIdMultisigISMRaw) Verify(_ context.Context, rawMetadata []byte, message util.HyperlaneMessage) (bool, error) {
	// get validator sigs
	metadata, err := NewMessageIdMultisigRawMetadata(rawMetadata)
	if err != nil {
		return false, err
	}

	digest := metadata.Digest(&message)

	// recover pubkeys from sigs and check they are in the validators list
	// make sure a threshold signed the digest
	return VerifyMultisig(m.Validators, m.Threshold, metadata.Signatures, digest)
}

func (m *MessageIdMultisigISMRaw) GetThreshold() uint32 {
	return m.Threshold
}

func (m *MessageIdMultisigISMRaw) GetValidators() []string {
	return m.Validators
}

func (m *MessageIdMultisigISMRaw) Validate() error {
	return ValidateNewMultisig(m)
}

type MessageIdMultisigRawMetadata struct {
	SignatureCount uint32
	Signatures     [][]byte
}

// NewMessageIdMultisigMetadata validates and creates a new metadata object
func NewMessageIdMultisigRawMetadata(metadata []byte) (MessageIdMultisigRawMetadata, error) {
	/*
	 * Format of metadata:
	 * [  68:????] Validator signatures (length := threshold * 65)
	 */
	// originMerkleTreeOffset := 0
	signaturesOffset := 68
	signatureLength := 65

	if len(metadata) < signaturesOffset {
		return MessageIdMultisigRawMetadata{}, fmt.Errorf("invalid metadata length: got %v, expected at least %v bytes", len(metadata), signaturesOffset)
	}

	signaturesLen := len(metadata) - signaturesOffset
	signatureCount := uint32(signaturesLen / signatureLength)

	if signaturesLen%signatureLength != 0 {
		return MessageIdMultisigRawMetadata{}, fmt.Errorf("invalid signatures length in metadata")
	}

	var signatures [][]byte
	for i := 0; i < int(signatureCount); i++ {
		start := signaturesOffset + (i * signatureLength)
		sig := make([]byte, signatureLength)
		copy(sig, metadata[start:start+signatureLength])
		signatures = append(signatures, sig)
	}

	return MessageIdMultisigRawMetadata{
		SignatureCount: uint32(signaturesLen / signaturesOffset),
		Signatures:     signatures,
	}, nil
}

func (m *MessageIdMultisigRawMetadata) Bytes() []byte {
	var signaturesBytes []byte
	for _, sig := range m.Signatures {
		signaturesBytes = append(signaturesBytes, sig...)
	}

	// blank bytes where signatures are in the official version
	empty := [68]byte{}

	return slices.Concat(
		empty[:],
		signaturesBytes,
	)
}

func (m *MessageIdMultisigRawMetadata) Digest(message *util.HyperlaneMessage) [32]byte {
	// we only care about message id
	return checkpointDigest(
		message.Origin,
		[32]byte{},
		[32]byte{},
		0,
		message.Id(),
	)
}
