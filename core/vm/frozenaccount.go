package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

var FrozenAcconts []*types.FrozenAccount = []*types.FrozenAccount{
	&types.FrozenAccount{
		Account:      common.HexToAddress("0x20BE1E8038d8d780Abcc427dc10b4843A10de7e6"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-15 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 600),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xE7F1ec04Df5F1062898063B569Ead16172d67303"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-16 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 600),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0x6ADAA65a8D22cf4b262d610797676c191D25eF23"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-20 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 600),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xDf72e45b52bb55a52a366C6B8f7628dc12983445"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-14 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 1200),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xDf72e45b52bb55a52a366C6B8f7628dc12983445"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-15 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 1200),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xDf72e45b52bb55a52a366C6B8f7628dc12983445"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-16 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 1200),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xDf72e45b52bb55a52a366C6B8f7628dc12983445"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-17 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 1200),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xDf72e45b52bb55a52a366C6B8f7628dc12983445"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-09-19 02:55:00"),
		//UnfrozenTime: uint64(time.Now().Unix() + 1200),
	},
}

func getBig(num string) *big.Int {
	bigNum, _ := new(big.Int).SetString(num, 10)
	return bigNum
}

func getUnixTimestamp(t string) uint64 {
	tutc, _ := time.Parse("2006-01-02 15:04:05", t)
	timestamp := tutc.Unix()
	return uint64(timestamp)
}
