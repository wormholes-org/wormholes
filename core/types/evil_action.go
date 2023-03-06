package types

import "github.com/ethereum/go-ethereum/rlp"

type EvilAction struct {
	Handled     bool
	EvilHeaders []*Header
}

func NewEvilAction(evilHeader *Header) *EvilAction {
	evilAction := EvilAction{
		Handled:     false,
		EvilHeaders: make([]*Header, 6),
	}
	evilAction.EvilHeaders = append(evilAction.EvilHeaders, evilHeader)
	return &evilAction
}

func EncodeEvilAction(ea EvilAction) ([]byte, error) {
	return rlp.EncodeToBytes(ea)
}

func DecodeEvilAction(data []byte) error {
	var ea EvilAction
	return rlp.DecodeBytes(data, &ea)
}
