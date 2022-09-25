package miner

import (
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
)

func GenerateAddrs() (addrs []common.Address) {
	for i := 0; i < 11; i++ {
		prikey, _ := crypto.GenerateKey()
		addrs = append(addrs, crypto.PubkeyToAddress(prikey.PublicKey))
	}
	return addrs
}
