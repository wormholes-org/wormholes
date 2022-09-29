package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type FrozenAccount struct {
	Account      common.Address
	Amount       *big.Int
	UnfrozenTime uint64
}
