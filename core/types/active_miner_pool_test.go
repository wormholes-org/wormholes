package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func TestSortByPledgeAmount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	addrs := make([]common.Address, 0)
	for i := 0; i < 5; i++ {
		priKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(priKey.PublicKey)
		addrs = append(addrs, addr)
	}

	activeMiners := new(ActiveMinerList)
	for i := 0; i < 5; i++ {
		activeMiners.ActiveMiners = append(activeMiners.ActiveMiners, &ActiveMiner{
			Address: addrs[i],
			Balance: big.NewInt(rand.Int63n(1000)),
			Height:  100,
		})
	}
	for i, v := range activeMiners.ActiveMiners {
		fmt.Println("i==", i, "v==", v.Balance, "addr==", v.Address.Hex())
	}
	fmt.Println("=============================")
	sorted, err := activeMiners.SortByPledgeAmount()
	if err != nil {
		t.Error(err.Error())
	}

	for i, v := range sorted.ActiveMiners {
		fmt.Println("i==", i, "v==", v.Balance, "addr==", v.Address.Hex())
	}

	fmt.Println("=============================")
	for i, v := range activeMiners.ActiveMiners {
		fmt.Println("i==", i, "v==", v.Balance, "addr==", v.Address.Hex())
	}
}


// 测试 7 + 4 方式选择validator
func TestSelectValidator(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	addrs := make([]common.Address, 0)
	for i := 0; i < 15; i++ {
		priKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(priKey.PublicKey)
		addrs = append(addrs, addr)
	}

	activeMiners := new(ActiveMinerList)
	for i := 0; i < 15; i++ {
		activeMiners.ActiveMiners = append(activeMiners.ActiveMiners, &ActiveMiner{
			Address: addrs[i],
			Balance: big.NewInt(rand.Int63n(1000)),
			Height:  100,
		})
	}

	sortedAddr, err := activeMiners.SortByPledgeAmount()
	for i, v := range sortedAddr.ActiveMiners {
		fmt.Println("fmt sortedAddr", "i", i, "addr==", v.Address, "balance", v.Balance)
	}
	if err != nil {
		t.Error(err.Error())
	}
	// Get the top 7 validators of the pledge amount
	fixedValidators := sortedAddr.ActiveMiners[:7]

	// Get 4 other validators besides the above 7
	random4Validators := activeMiners.ValidatorByDistanceAndWeight(activeMiners.ConvertToBigInt(sortedAddr.ActiveMiners[7:]), 4, randomHash())

	validators := make([]common.Address, 0)
	validators = append(validators, random4Validators...)
	for i, v := range fixedValidators {
		fmt.Println("append fixedValidators==", "i==", i, "v==", v.Address, "balance==", v.Balance)
		validators = append(validators, v.Address)
	}

	elevenValidator := new(ValidatorList)
	for i, addr := range validators {
		fmt.Println("append elevenValidator==", "i==", i, "addr==", addr)
		elevenValidator.AddValidator(addr, activeMiners.StakeBalance(addr), common.Address{})
	}

	for i, v := range elevenValidator.Validators {
		fmt.Println("bigInt==", v.Address().Hash().Big(), "i==", i, "addr==", v.Addr, "balance==", v.Balance, "proxy.addr", v.Proxy.Hex())
	}
}
