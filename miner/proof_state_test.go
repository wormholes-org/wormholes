package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"testing"
)

func TestProofStatePoolPut(t *testing.T) {
	//proofStatePool := NewProofStatePool()
	//
	//addrs := GenerateAddrs()
	//for i := 0; i < 11; i++ {
	//	//proofStatePool.Put(big.NewInt(1), addrs[0], addrs[i])
	//}
	//
	//for k, v := range proofStatePool.proofs {
	//	fmt.Println("current height:", k.String(), "===current count: ", v.count)
	//}
}

func TestOnlineValidator_GetAllAddress(t *testing.T) {
	var vals OnlineValidator
	vals = make(OnlineValidator)
	vals[common.HexToAddress("0x1000000000000000000000000000000000000000")] = struct{}{}
	vals[common.HexToAddress("0x1000000000000000000000000000000000000001")] = struct{}{}
	vals[common.HexToAddress("0x1000000000000000000000000000000000000002")] = struct{}{}
	vals[common.HexToAddress("0x1000000000000000000000000000000000000003")] = struct{}{}
	vals[common.HexToAddress("0x1000000000000000000000000000000000000004")] = struct{}{}
	addrs := vals.GetAllAddress()
	for _, addr := range addrs {
		t.Log(addr.Hex())
	}
}
