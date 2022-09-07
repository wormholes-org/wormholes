package vm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"time"
)

var FrozenAcconts []*types.FrozenAccount = []*types.FrozenAccount{
	&types.FrozenAccount{
		Account:      common.HexToAddress("0xC65F08C9Dfceb0988631B175E293Af5666535CF0"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-12-01"),
	},
	&types.FrozenAccount{
		Account:      common.HexToAddress("0x8C1931096C17f32FF6a1eFabe642422995a5013B"),
		Amount:       getBig("10000000000000000000000"),
		UnfrozenTime: getUnixTimestamp("2022-12-01"),
	},
}

func getBig(num string) *big.Int {
	bigNum, _ := new(big.Int).SetString(num, 10)
	return bigNum
}

func getUnixTimestamp(t string) uint64 {
	tutc, _ := time.Parse("2006-01-02", t)
	timestamp := tutc.Unix()
	return uint64(timestamp)
}
