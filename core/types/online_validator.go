package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type OnlineValidator struct {
	Height    *big.Int
	Address   common.Address
	Hash      common.Hash
	Signature []byte
}

func NewOnlineValidator(height *big.Int, addr common.Address, hash common.Hash, sig []byte) *OnlineValidator {
	return &OnlineValidator{Height: height, Address: addr, Hash: hash, Signature: sig}
}

type OnlineValidatorList struct {
	Height     *big.Int
	Validators []*OnlineValidator
}

func (ol *OnlineValidatorList) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(ol)
}

func (ol *OnlineValidatorList) Decode(data []byte) error {
	return rlp.DecodeBytes(data, ol)
}

func (ol *OnlineValidatorList) Size() int {
	if ol.Validators != nil {
		return len(ol.Validators)
	}
	return 0
}
