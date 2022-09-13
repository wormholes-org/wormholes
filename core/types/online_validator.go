package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type OnlineValidatorInfo struct {
	Height *big.Int
	Addrs  []common.Address
	Hashs  []common.Hash
}

func (o *OnlineValidatorInfo) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(o)
}

func (o *OnlineValidatorInfo) Decode(data []byte) error {
	return rlp.DecodeBytes(data, o)
}
