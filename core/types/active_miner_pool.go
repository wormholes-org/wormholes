package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sort"
)

type ActiveMiner struct {
	Address common.Address
	Balance *big.Int
	Height  uint64
}

func NewActiveMiner(addr common.Address, balance *big.Int, height uint64) *ActiveMiner {
	return &ActiveMiner{
		Address: addr,
		Balance: balance,
		Height:  height,
	}
}

type ActiveMinerList struct {
	ActiveMiners []*ActiveMiner
}

func (l *ActiveMinerList) Len() int {
	return len(l.ActiveMiners)
}

func (l *ActiveMinerList) Less(i, j int) bool {
	return l.ActiveMiners[i].Address.Hash().Big().Cmp(l.ActiveMiners[j].Address.Hash().Big()) < 0
}

func (l *ActiveMinerList) Swap(i, j int) {
	l.ActiveMiners[i], l.ActiveMiners[j] = l.ActiveMiners[j], l.ActiveMiners[i]
}

func (l *ActiveMinerList) AddAndUpdateActiveAddr(addr common.Address, balance *big.Int, height uint64) bool {
	for _, v := range l.ActiveMiners {
		if v.Address == addr {
			v.Height = height
			sort.Sort(l)
			return true
		}
	}
	l.ActiveMiners = append(l.ActiveMiners, NewActiveMiner(addr, balance, height))
	sort.Sort(l)
	return true
}

func (l *ActiveMinerList) RemoveActiveAddr(addr common.Address) (bool, error) {
	if res, index := l.GetIndex(addr); res {
		l.ActiveMiners = append(l.ActiveMiners[:index], l.ActiveMiners[index+1:]...)
		return true, nil
	}
	return false, errors.New("RemoveActiveAddr: The address was not found")
}

func (l *ActiveMinerList) GetIndex(addr common.Address) (bool, int) {
	for i, v := range l.ActiveMiners {
		if v.Address == addr {
			return true, i
		}
	}
	return false, -1
}

func (l *ActiveMinerList) ConvertToBigInt(validators []*ActiveMiner) (bigIntSlice []*big.Int) {
	if len(validators) == 0 {
		for _, m := range l.ActiveMiners {
			bigIntSlice = append(bigIntSlice, m.Address.Hash().Big())
		}
		return
	}
	for _, m := range validators {
		bigIntSlice = append(bigIntSlice, m.Address.Hash().Big())
	}
	return
}

// Query K validators closest to random numbers based on distance and pledge amount
// addr must be sorted !
func (l *ActiveMinerList) ValidatorByDistanceAndWeight(addr []*big.Int, k int, randomHash common.Hash) []common.Address {
	maxValue := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff").Hash().Big()

	addrToWeightMap := make(map[*big.Int]*big.Int, 0)

	r1 := randomHash[12:]
	x := common.BytesToAddress(r1).Hash().Big()

	for _, v := range addr {
		sub1 := big.NewInt(0)
		sub2 := big.NewInt(0)

		sub1 = sub1.Sub(v, x)
		sub1 = sub1.Abs(sub1)
		sub2 = sub2.Sub(maxValue, sub1)
		if sub1.Cmp(sub2) < 0 {
			a := new(big.Int).Mul(sub1, l.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, l.TotalStakeBalance())
			addrToWeightMap[v] = w
		} else {
			a := new(big.Int).Mul(sub2, l.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, l.TotalStakeBalance())
			addrToWeightMap[v] = w
		}
	}

	sortMap := rankByWordCount(addrToWeightMap)
	res := make([]common.Address, 0)

	for i := 0; i < sortMap.Len(); i++ {
		if i < k {
			res = append(res, common.BigToAddress(sortMap[i].Key))
		} else {
			break
		}
	}
	return res
}

// Calculate the total amount of the stake account
func (l *ActiveMinerList) TotalStakeBalance() *big.Int {
	var total = big.NewInt(0)
	for _, v := range l.ActiveMiners {
		total.Add(total, v.Balance)
	}
	return total
}

// Returns the amount of the staked node
func (l *ActiveMinerList) StakeBalance(address common.Address) *big.Int {
	for _, v := range l.ActiveMiners {
		if v.Address.Hex() != address.Hex() {
			continue
		}
		return v.Balance
	}
	return big.NewInt(0)
}

func (l *ActiveMinerList) SortByPledgeAmount() (*ActiveMinerList, error) {
	activeMinersCopy := new(ActiveMinerList)
	for _, m := range l.ActiveMiners {
		activeMinersCopy.ActiveMiners = append(activeMinersCopy.ActiveMiners, m)
	}

	sortedActiveMiners, err := InsertSort(activeMinersCopy)
	if err != nil {
		return nil, err
	}
	return sortedActiveMiners, nil
}

func InsertSort(arr *ActiveMinerList) (*ActiveMinerList, error) {
	if len(arr.ActiveMiners) <= 1 {
		return nil, errors.New("active miners less than one")
	}
	for i := 1; i < len(arr.ActiveMiners); i++ {
		current := arr.ActiveMiners[i]
		before := i - 1
		for before >= 0 && current.Balance.Cmp(arr.ActiveMiners[before].Balance) < 0 {
			arr.ActiveMiners[before+1] = arr.ActiveMiners[before]
			before--
		}
		arr.ActiveMiners[before+1] = current
	}
	return arr, nil
}

func (al *ActiveMinerList) GetByAddress(addr common.Address) int {
	for i, v := range al.ActiveMiners {
		if v.Address == addr {
			return i
		}
	}
	return -1
}

func (al *ActiveMinerList) GetByIndex(i uint64) ActiveMiner {
	if i >= uint64(al.Len()) {
		return ActiveMiner{}
	}
	return *al.ActiveMiners[i]
}
