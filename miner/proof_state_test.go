package miner

import (
	"fmt"
	"math/big"
	"testing"
)

func TestProofStatePoolPut(t *testing.T) {
	proofStatePool := NewProofStatePool()

	addrs := GenerateAddrs()
	for i := 0; i < 11; i++ {
		proofStatePool.Put(big.NewInt(1), addrs[0], addrs[i])
	}

	for k, v := range proofStatePool.proofs {
		fmt.Println("current height:", k.String(), "===current count: ", v.count)
	}
}
