package types

import (
	"github.com/ethereum/go-ethereum/rlp"
)

type EvilAction struct {
	Handled     bool
	EvilHeaders []*Header
}

func NewEvilAction(evilHeader *Header) *EvilAction {
	evilAction := EvilAction{
		Handled:     false,
		EvilHeaders: make([]*Header, 0),
	}
	evilAction.EvilHeaders = append(evilAction.EvilHeaders, evilHeader)
	return &evilAction
}

// @dev Check whether there is an evil header
func (ea *EvilAction) Exist(uncle *Header) bool {
	if uncle == nil {
		return false
	}

	for _, v := range ea.EvilHeaders {
		if v.Hash() == uncle.Hash() {
			return true
		}
	}
	return false
}

func EncodeEvilAction(ea EvilAction) ([]byte, error) {
	return rlp.EncodeToBytes(ea)
}

func DecodeEvilAction(data []byte) (EvilAction, error) {
	var ea EvilAction
	return ea, rlp.DecodeBytes(data, &ea)
}
