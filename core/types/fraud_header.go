package types

import "github.com/ethereum/go-ethereum/rlp"

type FraudHeader struct {
	RemoteHeader *Header // Represents the header of the remote block
	LocalHeader  *Header // The node header representing the latest local block
}

func NewFraudHeader(rh, lh *Header) *FraudHeader {
	return &FraudHeader{
		RemoteHeader: rh,
		LocalHeader:  lh,
	}
}

func EncodeFraudHeader(fh FraudHeader) ([]byte, error) {
	return rlp.EncodeToBytes(fh)
}

func DecodeFraudHeader(data []byte, fh FraudHeader) (FraudHeader, error) {
	return fh, rlp.DecodeBytes(data, &fh)
}
