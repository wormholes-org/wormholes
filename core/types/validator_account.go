package types

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"sort"
)

type Validator struct {
	Addr    common.Address
	Balance *big.Int
	Proxy   common.Address
}

func (v *Validator) Address() common.Address {
	return v.Addr
}

func NewValidator(addr common.Address, balance *big.Int, proxy common.Address) *Validator {
	return &Validator{Addr: addr, Balance: balance, Proxy: proxy}
}

type ValidatorList struct {
	Validators []*Validator
}

func NewValidatorList(validators []Validator) *ValidatorList {
	var validatorList *ValidatorList
	for i := 0; i < len(validators); i++ {
		validatorList.Validators = append(validatorList.Validators, &validators[i])
	}
	return validatorList
}

func (vl *ValidatorList) Len() int {
	return len(vl.Validators)
}

func (vl *ValidatorList) Less(i, j int) bool {
	return vl.Validators[i].Address().Hash().Big().Cmp(vl.Validators[j].Address().Hash().Big()) < 0
}

func (vl *ValidatorList) Swap(i, j int) {
	vl.Validators[i], vl.Validators[j] = vl.Validators[j], vl.Validators[i]
}

// Sort by distance in ascending order
func (sl *ValidatorList) AddValidator(addr common.Address, balance *big.Int, proxy common.Address) bool {
	empty := common.Address{}
	for _, v := range sl.Validators {
		if v.Address() == addr /*&& v.Proxy.String() == "0x0000000000000000000000000000000000000000" */ {
			// Usage scenarios: pledge, additional pledge, delegation
			v.Balance.Add(v.Balance, balance)
			if proxy.Hex() != empty.Hex() {
				v.Proxy = proxy
			}
			sort.Sort(sl)
			return true
		}
	}
	sl.Validators = append(sl.Validators, NewValidator(addr, balance, proxy))
	sort.Sort(sl)
	return true
}

func (sl *ValidatorList) RemoveValidator(addr common.Address, balance *big.Int) bool {
	for i, v := range sl.Validators {
		if v.Address() == addr {
			if v.Balance.Cmp(balance) > 0 {
				v.Balance.Sub(v.Balance, balance)
				sort.Sort(sl)
				return true
			} else if v.Balance.Cmp(balance) == 0 {
				v.Balance.Sub(v.Balance, balance)
				sl.Validators = append(sl.Validators[:i], sl.Validators[i+1:]...)
				return true
			}
			sl.Validators = append(sl.Validators[:i], sl.Validators[i+1:]...)
			return true
		}
	}
	return false
}

// Query K validators closest to random numbers based on distance and pledge amount
func (sl *ValidatorList) ValidatorByDistanceAndWeight(addr []*big.Int, k int, randomHash common.Hash) []common.Address {
	// The maximum value of address to big Int
	maxValue := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff").Hash().Big()

	// Record the weight corresponding to the address
	addrToWeightMap := make(map[*big.Int]*big.Int, 0)

	// Hash to 160-bit address
	r1 := randomHash[12:]
	x := common.BytesToAddress(r1).Hash().Big()

	for _, v := range addr {
		sub1 := big.NewInt(0)
		sub2 := big.NewInt(0)

		// The obtained sub1 and sub2 are two distance values, which need to be taken from the smallest
		sub1 = sub1.Sub(v, x)
		sub1 = sub1.Abs(sub1)
		sub2 = sub2.Sub(maxValue, sub1)
		if sub1.Cmp(sub2) < 0 {
			a := new(big.Int).Mul(sub1, sl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, sl.TotalStakeBalance())
			addrToWeightMap[v] = w
		} else {
			a := new(big.Int).Mul(sub2, sl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, sl.TotalStakeBalance())
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
func (sl *ValidatorList) TotalStakeBalance() *big.Int {
	var total = big.NewInt(0)
	for _, voter := range sl.Validators {
		total.Add(total, voter.Balance)
	}
	return total
}

// Returns the amount of the staked node
func (sl *ValidatorList) StakeBalance(address common.Address) *big.Int {
	for _, st := range sl.Validators {
		if st.Address().Hex() != address.Hex() && st.Proxy.Hex() != address.Hex() {
			continue
		}
		return st.Balance
	}
	return big.NewInt(0)
}

func (sl *ValidatorList) ConvertToAddress() (addrs []common.Address) {
	for _, validator := range sl.Validators {
		addrs = append(addrs, validator.Addr)
	}
	return
}

func (sl *ValidatorList) ConvertToBigInt() (bigIntSlice []*big.Int) {
	for _, validator := range sl.Validators {
		bigIntSlice = append(bigIntSlice, validator.Addr.Hash().Big())
	}
	return
}

func (sl *ValidatorList) Exist(addr common.Address) bool {
	for _, v := range sl.Validators {
		if v.Addr == addr || v.Proxy == addr {
			return true
		}
	}
	return false
}

func (sl *ValidatorList) ExistProxy(addr common.Address) bool {
	emptyAddr := common.Address{}
	for _, v := range sl.Validators {
		if v.Addr == addr && v.Proxy != emptyAddr {
			return true
		}
	}
	return false
}
