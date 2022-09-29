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

func randomHash() common.Hash {
	rand.Seed(time.Now().Local().UnixMicro())
	var hash common.Hash
	if n, err := rand.Read(hash[:]); n != common.HashLength || err != nil {
		panic(err)
	}
	return hash
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
