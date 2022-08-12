package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

func TestInjectedOfficialNFTList(t *testing.T) {
	injected1 := &InjectedOfficialNFT{
		Dir:        "/ipfs/test1111",
		StartIndex: new(big.Int).SetInt64(0),
		Number:     65536,
		Royalty:    100,
		Creator:    "0xB7987546EA03f4167e1F424C89C094BebbC112A6",
	}
	injected2 := &InjectedOfficialNFT{
		Dir:        "/ipfs/test2222",
		StartIndex: new(big.Int).SetInt64(65536),
		Number:     131072,
		Royalty:    100,
		Creator:    "0xB7987546EA03f4167e1F424C89C094BebbC112A6",
	}

	var injectedList InjectedOfficialNFTList
	injectedList.InjectedOfficialNFTs = append(injectedList.InjectedOfficialNFTs, injected1)
	injectedList.InjectedOfficialNFTs = append(injectedList.InjectedOfficialNFTs, injected2)

	address := common.HexToAddress("0x8000000000000000000000000000000000010000")

	injectedList.GetInjectedInfo(address)

}

func TestSNFTExchangeList(t *testing.T) {
	snftExchange := &SNFTExchange{
		InjectedInfo: InjectedInfo{
			MetalUrl: "/ipfs/test1111",
			Royalty:  100,
			Creator:  "0xB7987546EA03f4167e1F424C89C094BebbC112A6",
		},
		NFTAddress:         common.HexToAddress("0x8000000000000000000000000000000000000100"),
		MergeLevel:         2,
		CurrentMintAddress: common.HexToAddress("0x8000000000000000000000000000000000000100"),
		BlockNumber:        new(big.Int).SetUint64(1),
	}
	var SNFTExchangePool SNFTExchangeList

	SNFTExchangePool.SNFTExchanges = append(SNFTExchangePool.SNFTExchanges, snftExchange)
	var num int = 0
	for {
		address, _, b := SNFTExchangePool.PopAddress(new(big.Int).SetInt64(2))
		if !b {
			t.Log("num= ", num)
			break
		}
		num++
		t.Log(address.Hex())
	}

}
