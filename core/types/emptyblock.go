package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const DEFAULT_VALIDATOR_COEFFICIENT = 70

type EmptyMessageEvent struct {
	Sender  common.Address
	Height  *big.Int
	Payload []byte
}

type SignatureData struct {
	Vote   common.Address
	Height *big.Int
	//Timestamp uint64
	Round uint64
}
