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

func TestCalculateAddressRange(t *testing.T) {
	rand.Seed(time.Now().Unix())
	var validators []*Validator

	for i := 0; i < 100; i++ {
		validators = append(validators, NewValidator(RandomAddr(), big.NewInt(rand.Int63()), common.Address{}))
	}
	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	//---------------------------------------------//
	for _, v := range validatorList.Validators {
		fmt.Println("address---", v.Addr, "|---weight", v.Weight)
	}
}

func TestCalculateAddressRange2(t *testing.T) {
	var validators []*Validator

	stakeAmt := []*big.Int{
		big.NewInt(10),
		big.NewInt(10),
		big.NewInt(10),
		big.NewInt(10),
	}

	for i := 0; i < len(stakeAmt); i++ {
		validators = append(validators, NewValidator(RandomAddr(), stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}

	//---------------------------------------------//
	for _, v := range validatorList.Validators {
		fmt.Println("address---", v.Addr, "|---weight", v.Weight)
	}

	hash := randomHash()
	err, res := validatorList.CollectValidators(hash, 11)
	if err != nil {
		fmt.Println("error collect validators", err)
	}
	for i, v := range res {
		fmt.Println("i===", i, "v====", v, "balance=====", validatorList.StakeBalance(v))
	}
}

func TestCollectValidators(t *testing.T) {
	var count int
	for i := 0; i < 1000; i++ {
		//maxValue := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff").Hash().Big()
		validatorList := prepareFixedValidator()
		for _, vl := range validatorList.Validators {
			validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
		}
		//for _, v := range validatorList.Validators {
		//	fmt.Println("rangeList====", "addr===", v.Addr, "weight====", v.Weight, "maxValue====", maxValue)
		//}
		hash := randomHash()
		err, res := validatorList.CollectValidators(hash, 11)
		if err != nil {
			fmt.Println("error collect validators", err)
		}
		if len(res) == 11 {
			count++
		}
	}
	fmt.Println("------", count)

	//random11Amt := big.NewInt(0)
	//for i, v := range res {
	//	fmt.Println("i===", i, "v====", v, "balance=====", validatorList.StakeBalance(v))
	//	random11Amt.Add(random11Amt, validatorList.StakeBalance(v))
	//}
	//fmt.Println("11StakeAmt====", random11Amt, "totalAmt", validatorList.TotalStakeBalance(), "hash", hash.Hex())
}

func prepareFixedValidator() *ValidatorList {
	var validators []*Validator
	stakeAmt := []*big.Int{
		big.NewInt(14),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(15),
		big.NewInt(16),
		big.NewInt(17),
		big.NewInt(1),
		big.NewInt(7),
		big.NewInt(7),
		big.NewInt(5),
		big.NewInt(4),
		big.NewInt(3),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(15),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(13),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(12),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
		big.NewInt(15),
		big.NewInt(16),
		big.NewInt(17),
	}

	addrs := []common.Address{
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x31534d5C7b1eabb73425c2361661b878F4322f9D"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05"),
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),

		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
		common.HexToAddress(RandomAddr().Hex()),
	}

	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	return validatorList
}

func RandomAddr() common.Address {
	priKey, _ := crypto.GenerateKey()
	return crypto.PubkeyToAddress(priKey.PublicKey)
}
