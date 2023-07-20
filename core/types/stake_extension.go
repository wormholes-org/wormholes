package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

type StakersExtensionList struct {
	StakerExtensions []*StakerExtension
}
type StakerExtension struct {
	Addr        common.Address
	Balance     *big.Int
	BlockNumber *big.Int
}

func (sl *StakersExtensionList) AddStakerPledge(addr common.Address, balance *big.Int, blocknumber *big.Int) bool {
	for _, v := range sl.StakerExtensions {
		if v.Addr == addr {
			v.Balance.Add(v.Balance, balance)
			v.BlockNumber = blocknumber
			return true
		}
	}
	sl.StakerExtensions = append(sl.StakerExtensions, &StakerExtension{Addr: addr, Balance: balance, BlockNumber: blocknumber})
	return true
}

func (sl *StakersExtensionList) RemoveStakerPledge(addr common.Address, balance *big.Int) bool {
	for i, v := range sl.StakerExtensions {
		if v.Addr == addr {
			if v.Balance.Cmp(balance) > 0 {
				v.Balance.Sub(v.Balance, balance)
				return true
			} else {
				v.Balance.Sub(v.Balance, balance)
				sl.StakerExtensions = append(sl.StakerExtensions[:i], sl.StakerExtensions[i+1:]...)
				return true
			}
		}
	}
	return false
}
func NewStakerPledge(addr common.Address, balance *big.Int, blocknumber *big.Int) *StakerExtension {
	return &StakerExtension{Addr: addr, Balance: balance, BlockNumber: blocknumber}
}

func (sl *StakersExtensionList) DeepCopy() *StakersExtensionList {
	tempStakerList := &StakersExtensionList{
		StakerExtensions: make([]*StakerExtension, 0, len(sl.StakerExtensions)),
	}
	for _, staker := range sl.StakerExtensions {
		tempStaker := StakerExtension{
			Addr:        staker.Addr,
			Balance:     new(big.Int).Set(staker.Balance),
			BlockNumber: new(big.Int).Set(staker.BlockNumber),
		}
		tempStakerList.StakerExtensions = append(tempStakerList.StakerExtensions, &tempStaker)
	}
	return tempStakerList
}

func (sl *StakersExtensionList) IsExist(addr common.Address) bool {
	if sl == nil ||
		len(sl.StakerExtensions) == 0 {
		return false
	}

	for _, staker := range sl.StakerExtensions {
		if staker.Addr == addr {
			return true
		}
	}

	return false
}

func (sl *StakersExtensionList) GetBalance(addr common.Address) *big.Int {
	if sl == nil || sl.StakerExtensions == nil {
		return big.NewInt(0)
	}

	for _, staker := range sl.StakerExtensions {
		if staker.Addr == addr {
			return staker.Balance
		}
	}

	return big.NewInt(0)
}

func (sl *StakersExtensionList) GetLen() int {
	if sl == nil || sl.StakerExtensions == nil {
		return 0
	}

	return len(sl.StakerExtensions)
}
