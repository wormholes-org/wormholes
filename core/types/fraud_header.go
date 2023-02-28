package types

import "github.com/ethereum/go-ethereum/rlp"

type FraudHeader struct {
	RemoteParentHeader *Header // Represents the parent header of the remote block
	LocalHeader        *Header // The node header representing the latest local block
}

func NewFraudHeader(rh, lh *Header) *FraudHeader {
	return &FraudHeader{
		RemoteParentHeader: rh,
		LocalHeader:        lh,
	}
}

func EncodeFraudHeader(fh FraudHeader) ([]byte, error) {
	return rlp.EncodeToBytes(fh)
}

func DecodeFraudHeader(data []byte) error {
	var fh FraudHeader
	return rlp.DecodeBytes(data, &fh)
}
