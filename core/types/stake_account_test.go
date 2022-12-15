package types

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestBeValidatorProbaBilitity(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x9A1711a10e3d5baA4e0cE970dF6E33DC50EF0992"),
		common.HexToAddress("0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c"),
		common.HexToAddress("0xEdfC22E9CfB4e24815C3a12e81bF10caB9cE4D26"),
		common.HexToAddress("0x085ABc35ed85d26C2795b64C6fFb89B68aB1c479"),
		common.HexToAddress("0xb31b41E5EF219fB0CC9935Ad914158cf8970DB44"),
	}
	rand.Seed(time.Now().UnixNano())
	var sl StakerList
	for i := int64(0); i < int64(len(addrs)); i++ {
		sl.AddStaker(addrs[i], big.NewInt(int64(100*(rand.Intn(122)))))
	}
	addrsToBigSlice := SortAddr(addrs)
	for _, bigIntSlice := range addrsToBigSlice {
		fmt.Println("bigIntSlice:", bigIntSlice, "addr==", common.BigToAddress(bigIntSlice).String())
	}
	fmt.Println("=========================")

	//randomValue, _ := new(big.Int).SetString("b31b41E5EF219fB0CC9935Ad914158cf8970DB44", 16)
	countMap := make(map[string]int, 12)
	fmt.Println("time.now==before", time.Now())

	for i := 0; i < 10000; i++ {
		randomValue := randomHash()
		res := sl.ValidatorByDistanceAndWeight(addrsToBigSlice, 2, randomValue)

		for _, addr := range res {
			countMap[addr.Hex()] += 1
		}
	}
	fmt.Println("time.now==after", time.Now())

	fmt.Println("final Total")
	for addr, count := range countMap {
		fmt.Println("addr:=", addr, "bigint", common.HexToAddress(addr).Hash().Big().String(), "==count:", count, "balance", sl.StakeBalance(common.HexToAddress(addr)))
	}
}

func TestAddStakerSortByAddrDescend(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x085ABc35ed85d26C2795b64C6fFb89B68aB1c479"),
		common.HexToAddress("0xb31b41E5EF219fB0CC9935Ad914158cf8970DB44"),
		common.HexToAddress("0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c"),
		common.HexToAddress("0x9A1711a10e3d5baA4e0cE970dF6E33DC50EF0992"),
		common.HexToAddress("0xEdfC22E9CfB4e24815C3a12e81bF10caB9cE4D26"),
	}
	var sl StakerList
	for i := int64(0); i < int64(len(addrs)); i++ {
		sl.AddStaker(addrs[i], big.NewInt(10000*(i+1)))
	}
	for i, staker := range sl.Stakers {
		fmt.Println("i", i, "addr", staker.Addr.Hex(), "balance", staker.Balance.String(), "bigInt", staker.Addr.Hash().Big().String())
	}
}

func TestValidatorByWeightAndDistance(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x085ABc35ed85d26C2795b64C6fFb89B68aB1c479"),
		common.HexToAddress("0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c"),
		common.HexToAddress("0x9A1711a10e3d5baA4e0cE970dF6E33DC50EF0992"),
		common.HexToAddress("0xb31b41E5EF219fB0CC9935Ad914158cf8970DB44"),
		common.HexToAddress("0xEdfC22E9CfB4e24815C3a12e81bF10caB9cE4D26"),
	}
	var sl StakerList
	for i := int64(0); i < int64(len(addrs)); i++ {
		s := NewStaker(addrs[i], big.NewInt(10000*(i+1)))
		sl.Stakers = append(sl.Stakers, s)
	}
	addrsToBigSlice := SortAddr(addrs)
	for _, bigIntSlice := range addrsToBigSlice {
		fmt.Println("bigIntSlice:", bigIntSlice, "addr==", common.BigToAddress(bigIntSlice).String())
	}
	fmt.Println("=========================")

	addrsToBigSlice = SortAddr(addrs)
	for _, bigIntSlice := range addrsToBigSlice {
		fmt.Println("bigIntSlice:", bigIntSlice, "addr==", common.BigToAddress(bigIntSlice).String())
	}
	randomValue := randomHash()
	res := sl.ValidatorByDistanceAndWeight(addrsToBigSlice, 5, randomValue)
	for _, addr := range res {
		fmt.Println("res", addr.String())
	}
}

// Batch generation of private keys and related addresses
func TestBatchGenPriKeyAndAddr(t *testing.T) {
	for i := 0; i < 7; i++ {
		priKey, _ := crypto.GenerateKey()
		hexPriKey := common.Bytes2Hex(crypto.FromECDSA(priKey))
		addr := crypto.PubkeyToAddress(priKey.PublicKey)
		fmt.Println("i=", i, "hex=", hexPriKey, "addr=", addr.Hex())
	}
}

func TestGenerateSevenValidator(t *testing.T) {
	var addrs []common.Address

	for i := 0; i < 7; i++ {
		priKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(priKey.PublicKey)
		addrs = append(addrs, addr)
	}

	var sl StakerList
	for i := int64(0); i < int64(len(addrs)); i++ {
		s := NewStaker(addrs[i], big.NewInt(10000*(i+1)))
		sl.Stakers = append(sl.Stakers, s)
	}
	addrsToBigSlice := SortAddr(addrs)
	for _, bigIntSlice := range addrsToBigSlice {
		fmt.Println("bigIntSlice:", bigIntSlice, "addr==", common.BigToAddress(bigIntSlice).String())
	}
	fmt.Println("=========================")

	addrsToBigSlice = SortAddr(addrs)
	for _, bigIntSlice := range addrsToBigSlice {
		fmt.Println("bigIntSlice:", bigIntSlice, "addr==", common.BigToAddress(bigIntSlice).String())
	}
	randomValue := randomHash()
	res := sl.ValidatorByDistanceAndWeight(addrsToBigSlice, 11, randomValue)
	for _, addr := range res {
		fmt.Println("res", addr.String())
	}
}

func TestPrintAddr(t *testing.T) {
	pri, _ := crypto.HexToECDSA("7b2546a5d4e658d079c6b2755c6d7495edd01a686fddae010830e9c93b23e398")
	addr := crypto.PubkeyToAddress(pri.PublicKey)
	fmt.Println("addr=", addr.Hex())
}

func BenchmarkValidatorByDistanceAndWeight(b *testing.B) {
	var addrs []common.Address
	for i := 0; i < 100; i++ {
		priKey, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(priKey.PublicKey)
		addrs = append(addrs, addr)
	}

	var sl StakerList
	for i := int64(0); i < int64(len(addrs)); i++ {
		s := NewStaker(addrs[i], big.NewInt(10000*(i+1)))
		sl.Stakers = append(sl.Stakers, s)
	}

	b.StartTimer()
	addrsToBigSlice := SortAddr(addrs)
	for i := 0; i < b.N; i++ {
		randomValue := randomHash()
		_ = sl.ValidatorByDistanceAndWeight(addrsToBigSlice, 11, randomValue)
	}
	b.StopTimer()
}

func TestLowPledgeProbability(t *testing.T) {
	var count int
	addrs := MockStakersData()
	addrsToBigSlice := SortAddr(addrs)
	sl := MockStakers()
	for i := 0; i < 100000; i++ {
		rv := randomHash()
		res := sl.ValidatorByDistanceAndWeight(addrsToBigSlice, 4, rv)
		for _, v := range res {
			if v == common.HexToAddress("0xAc17A48d782DE6985D471F6EEcB1023C08f0CB05") ||
				v == common.HexToAddress("0xa337175D23bB0ee4f879e38B73973CDe829b9d29") ||
				v == common.HexToAddress("0x9baeC4D975B0Dd5baa904200Ea11727eDD593be6") ||
				v == common.HexToAddress("0x77c33e2951d2851D9D648161BEfAC8f4C60D1181") ||
				v == common.HexToAddress("0x5a7f652BC51Fb99747fe2641A8CDd6CbFE51b201") ||
				v == common.HexToAddress("0x2c2909db351764D92ecc313a0D8baF72735C5165") {
				count++
				break
			}
		}
	}
	fmt.Println("final res", "count==", count)
}

func MockStakersData() []common.Address {
	return []common.Address{
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0xE2FA892CC5CC268a0cC1d924EC907C796351C645"),
		common.HexToAddress("0xdb33217fE3F74bD41c550B06B624E23ab7f55d05"),
		common.HexToAddress("0xc067825f4B7a53Bb9f2Daf72fF22C8EE39736afF"),
		common.HexToAddress("0xbbaE84E9879F908700c6ef5D15e928Abfb556a21"),
		common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6"),
		common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),
		common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963"),
		common.HexToAddress("0x31534d5C7b1eabb73425c2361661b878F4322f9D"),
		common.HexToAddress("0x20cb28AE861c322A9A86b4F9e36Ad6977930fA05"),
		common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD"),
		common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F"),

		//
		common.HexToAddress("0xAc17A48d782DE6985D471F6EEcB1023C08f0CB05"),
		common.HexToAddress("0xa337175D23bB0ee4f879e38B73973CDe829b9d29"),
		common.HexToAddress("0x9baeC4D975B0Dd5baa904200Ea11727eDD593be6"),
		common.HexToAddress("0x77c33e2951d2851D9D648161BEfAC8f4C60D1181"),
		common.HexToAddress("0x5a7f652BC51Fb99747fe2641A8CDd6CbFE51b201"),
		common.HexToAddress("0x2c2909db351764D92ecc313a0D8baF72735C5165"),
	}
}

func MockStakers() *StakerList {
	addrs := MockStakersData()

	c := 700
	c2 := 70000
	stakeAmt := []*big.Int{
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),
		big.NewInt(int64(c)),

		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
		big.NewInt(int64(c2)),
	}

	var stakers []*Staker
	for i := 0; i < len(addrs); i++ {
		stakers = append(stakers, NewStaker(addrs[i], stakeAmt[i]))
	}
	stakerList := &StakerList{Stakers: stakers}
	return stakerList
}
