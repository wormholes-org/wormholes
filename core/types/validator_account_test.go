package types

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
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

func TestStability(t *testing.T) {
	//cannot have random addresses
	validatorList := prepareFixedValidator()
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	for i, v := range validatorList.Validators {
		fmt.Println("=====i====", i, "====weight====", v.Weight)
	}
	hash := common.HexToHash("0xba84ce252cfa4ac640ddbbec2e634d544cb4aacf1880b24126d9300b9a38534d")
	err, res := validatorList.CollectValidators(hash, 11)
	if err != nil {
		fmt.Println("error collect validators", err)
	}
	random11Amt := big.NewInt(0)
	for i, v := range res {
		fmt.Println("i===", i, "v====", v, "balance=====", validatorList.StakeBalance(v))
		random11Amt.Add(random11Amt, validatorList.StakeBalance(v))
	}
	fmt.Println("11StakeAmt====", random11Amt, "totalAmt", validatorList.TotalStakeBalance(), "hash", hash.Hex())
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

		//
		big.NewInt(17),
		big.NewInt(1),
		big.NewInt(1),
		big.NewInt(13),
		big.NewInt(13),
		big.NewInt(14),
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

		//
		common.HexToAddress("0xB6FD5851a8c1d9B1C22a210664Fbe7187C137582"),
		common.HexToAddress("0xa4E91908d98aC1b0F232B6873F0989cDE07c7C71"),
		common.HexToAddress("0x1778B78658dDb31a8F0b8ba80E8471225050c62d"),
		common.HexToAddress("0xC12703f9708eE5A5A704696Ea3Dcb0f1c784273a"),
		common.HexToAddress("0xb67Db7D7e97486b3f23369160363430c0e98dFC9"),
		common.HexToAddress("0x7d5412AeA1e796fC58b3A8Be12a2d853528007dC"),
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

func TestGetAddr(t *testing.T) {
	for i, address := range Get11Addr() {
		fmt.Println("i====", i, "====addr===", address)
	}
}

func Get11Addr() []common.Address {
	var addrs []common.Address
	for i := 0; i < 14; i++ {
		priKey, _ := crypto.GenerateKey()
		addrs = append(addrs, crypto.PubkeyToAddress(priKey.PublicKey))
	}
	return addrs
}

func TestRandomValidatorV2(t *testing.T) {
	validatorList := prepareFixedValidator()
	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	// This hash will calculate less than 15 validators
	// Stable calculation of a validator
	hash := common.HexToHash("0xd95b0361af0635770474311d57c543c90c609b41d18fa0af95583872c3e2ad6f")
	vals := validatorList.RandomValidatorV2(15, hash)
	for i, val := range vals {
		fmt.Println("====i====", i, "====val====", val)
	}
}

func TestProbability(t *testing.T) {
	var count int
	vl := MockValidatorList()
	for i := 0; i < 1000000; i++ {
		hash := randomHash()
		selectedValidators := vl.RandomValidatorV2(11, hash)

		// 1.The probability of an address appearing=========================start//
		//for _, v := range selectedValidators {
		//	if v == common.HexToAddress("0x53182359904a07925b3338260b3ee8CD7bAf1935") {
		//		count++
		//	}
		//}
		// The probability of an address appearing=========================end//

		//2. All the validators with 70,000 pledged amount are recorded in===============start//
		var specifiedValidator int
		for _, v := range selectedValidators {
			if vl.StakeBalance(v).Cmp(big.NewInt(800000)) < 0 {
				specifiedValidator++
			}
		}
		if specifiedValidator == 11 {
			count++
		}
		// All the validators with 70,000 pledged amount are recorded in===============end//
	}
	fmt.Println("====count===", count)
}

func MockValidatorList() *ValidatorList {
	var validators []*Validator
	stakeAmt := []*big.Int{
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),
		big.NewInt(700001),

		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
		big.NewInt(700000),
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

		//
		common.HexToAddress("0x53182359904a07925b3338260b3ee8CD7bAf1935"),
		common.HexToAddress("0xa01c49F5206DEB26cA7e8b6336E8D61F73b77ba8"),
		common.HexToAddress("0x007a585aC6607F550b6d60eCaEB88b88Ac479daD"),
		common.HexToAddress("0xEe7696c15fc6F3e72Db61867db46d95cC8C8A54F"),
		common.HexToAddress("0x4862713256F2D306029cA4e034292338aF1ef52f"),
		common.HexToAddress("0xc3a690d02FB7132378E7875a03E08Aff7b7Abb27"),
		common.HexToAddress("0xb107cCFd67fc281D9F00D1C41d168A948F0cC469"),
		common.HexToAddress("0x36B651f4CF1E7622783D6E534C667fEd02c559B9"),
		common.HexToAddress("0x7cb6bb45c1C379ADb1665826bcEDAB1c10FC520E"),
		common.HexToAddress("0xa8e9127D4F708aF1c8120dBdaEbb424A9dBb4245"),
		common.HexToAddress("0xA727FB70b907582Db1Ece0098b0784c5d6C128b3"),
		common.HexToAddress("0x8aaa8d0A8cB518b1979e4cF1b9E86607FbCFB9aF"),
		common.HexToAddress("0x84C2fE50aca4BeCF37ca8B240C6D71bA63FF3D93"),
		common.HexToAddress("0x5ff38AeFa2377e93190e05258302d9ffB18ADa0b"),
		common.HexToAddress("0xDc058854d613F2791064CE9Fa9E9Fb46CAE1D312"),
		common.HexToAddress("0x3a78bbE2Cd8863b9aC7A7e6c4906fCF64Ef2c484"),
	}

	for i := 0; i < len(addrs); i++ {
		validators = append(validators, NewValidator(addrs[i], stakeAmt[i], common.Address{}))
	}
	validatorList := NewValidatorList(validators)

	for _, vl := range validatorList.Validators {
		validatorList.CalculateAddressRange(vl.Addr, validatorList.StakeBalance(vl.Addr))
	}
	return validatorList
}

func TestRlpToHash(t *testing.T) {
	validator := NewValidator(RandomAddr(), big.NewInt(1000000), common.Address{})
	encValidator, _ := rlp.EncodeToBytes(validator)
	hash := common.BytesToHash(encValidator)
	t.Log("===hash===", hash.Hex())
}
