package types

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

type Validator struct {
	Addr    common.Address
	Balance *big.Int
	Proxy   common.Address
	Weight  []*big.Int
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

func NewValidatorList(validators []*Validator) *ValidatorList {
	validatorList := new(ValidatorList)
	for _, v := range validators {
		validatorList.AddValidator(v.Addr, v.Balance, v.Proxy)
	}
	return validatorList
}

func (vl *ValidatorList) Len() int {
	return len(vl.Validators)
}

// Less Sort by pledge amount in descending order
func (vl *ValidatorList) Less(i, j int) bool {
	return vl.Validators[i].Balance.Cmp(vl.Validators[j].Balance) > 0
}

func (vl *ValidatorList) Swap(i, j int) {
	vl.Validators[i], vl.Validators[j] = vl.Validators[j], vl.Validators[i]
}

// AddValidator Sort by distance in ascending order
func (vl *ValidatorList) AddValidator(addr common.Address, balance *big.Int, proxy common.Address) bool {
	empty := common.Address{}
	for _, v := range vl.Validators {
		if v.Address() == addr /*&& v.Proxy.String() == "0x0000000000000000000000000000000000000000" */ {
			// Usage scenarios: pledge, additional pledge, delegation
			v.Balance.Add(v.Balance, balance)
			if proxy.Hex() != empty.Hex() {
				v.Proxy = proxy
			}
			sort.Sort(vl)
			return true
		}
	}
	vl.Validators = append(vl.Validators, NewValidator(addr, balance, proxy))
	sort.Sort(vl)
	return true
}

func (vl *ValidatorList) RemoveValidator(addr common.Address, balance *big.Int) bool {
	for i, v := range vl.Validators {
		if v.Address() == addr {
			if v.Balance.Cmp(balance) > 0 {
				v.Balance.Sub(v.Balance, balance)
				sort.Sort(vl)
				//vl.CalculateAddressRange(addr, balance)
				return true
			} else if v.Balance.Cmp(balance) == 0 {
				v.Balance.Sub(v.Balance, balance)
				vl.Validators = append(vl.Validators[:i], vl.Validators[i+1:]...)
				//vl.CalculateAddressRange(addr, balance)
				return true
			}
			vl.Validators = append(vl.Validators[:i], vl.Validators[i+1:]...)
			//vl.CalculateAddressRange(addr, balance)
			return true
		}
	}
	return false
}

func (vl *ValidatorList) CalculateAddressRange(address common.Address, stakeAmt *big.Int) {
	addrNo := address.Hash().Big()
	totalAmt := vl.TotalStakeBalance()
	minValue := big.NewInt(0)
	maxValue := common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff").Hash().Big()

	// stakeAmt / totalStakeAmt * maxValue * 7
	rangeLength := big.NewInt(0).Div(big.NewInt(0).Mul(big.NewInt(7), big.NewInt(0).Mul(stakeAmt, maxValue)), totalAmt)

	if rangeLength.Cmp(maxValue) > 0 {
		addrRange := []*big.Int{minValue, maxValue}
		vl.AddAddrRange(address, addrRange)
	}

	if big.NewInt(0).Add(addrNo, rangeLength).Cmp(maxValue) < 0 {
		addrRange := []*big.Int{addrNo, big.NewInt(0).Add(addrNo, rangeLength)}
		vl.AddAddrRange(address, addrRange)
	} else {
		modValue := big.NewInt(0).Mod(big.NewInt(0).Add(addrNo, rangeLength), maxValue)
		addrRange := []*big.Int{addrNo, maxValue, minValue, modValue}
		vl.AddAddrRange(address, addrRange)
	}

	//if rangeLenth.Cmp(maxValue) > 0 {
	//	addrRange := []*big.Int{minValue, maxValue}
	//	vl.AddAddrRange(address, addrRange)
	//}
	//if addrNo.Cmp(rangeLenth) > 0 && big.NewInt(0).Add(addrNo, rangeLenth).Cmp(maxValue) < 0 {
	//	addrRange := []*big.Int{big.NewInt(0).Sub(addrNo, rangeLenth), big.NewInt(0).Add(addrNo, rangeLenth)}
	//	vl.AddAddrRange(address, addrRange)
	//}
	//
	//if addrNo.Cmp(rangeLenth) < 0 && big.NewInt(0).Add(addrNo, rangeLenth).Cmp(maxValue) < 0 {
	//	addrRange := []*big.Int{minValue, big.NewInt(0).Add(addrNo, rangeLenth),
	//		big.NewInt(0).Add(big.NewInt(0).Sub(maxValue, rangeLenth), addrNo), maxValue}
	//	vl.AddAddrRange(address, addrRange)
	//}
	//
	//if addrNo.Cmp(rangeLenth) > 0 && big.NewInt(0).Add(addrNo, rangeLenth).Cmp(maxValue) > 0 {
	//	addrRange := []*big.Int{big.NewInt(0).Sub(addrNo, rangeLenth), maxValue, minValue, big.NewInt(0).Sub(big.NewInt(0).Add(addrNo, rangeLenth), maxValue)}
	//	vl.AddAddrRange(address, addrRange)
	//}

	return
}

// ValidatorByDistanceAndWeight Query K validators closest to random numbers based on distance and pledge amount
func (vl *ValidatorList) ValidatorByDistanceAndWeight(addr []*big.Int, k int, randomHash common.Hash) []common.Address {
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
			a := new(big.Int).Mul(sub1, vl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, vl.TotalStakeBalance())
			addrToWeightMap[v] = w
		} else {
			a := new(big.Int).Mul(sub2, vl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, vl.TotalStakeBalance())
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

func (vl *ValidatorList) RandomValidatorV2(k int, randomHash common.Hash) []common.Address {
	err, validators := vl.CollectValidators(randomHash, k)
	if err != nil {
		return []common.Address{}
	}

	// Make up for less than K
	diffCount := k - len(validators)
	for _, v := range vl.Validators {
		flg := false
		for _, vv := range validators {
			if vv == v.Addr {
				flg = true
				break
			}
		}
		if !flg && diffCount > 0 {
			validators = append(validators, v.Addr)
			diffCount--
		}
	}
	return validators
}

// CollectValidators Collect the k validators closest to the drop point
func (vl *ValidatorList) CollectValidators(randomHash common.Hash, k int) (error, []common.Address) {
	r1 := randomHash[12:]
	point := common.BytesToAddress(r1).Hash().Big()

	var validators []common.Address
	var count int
	for _, v := range vl.Validators {
		if count == k {
			break
		}
		if v.Weight == nil {
			continue
		}
		if len(v.Weight) == 2 {
			if point.Cmp(v.Weight[0]) > 0 && point.Cmp(v.Weight[1]) < 0 {
				validators = append(validators, v.Addr)
				count++
			}
		}
		if len(v.Weight) == 4 {
			if (point.Cmp(v.Weight[0]) > 0 && point.Cmp(v.Weight[1]) < 0) ||
				(point.Cmp(v.Weight[2]) > 0 && point.Cmp(v.Weight[3]) < 0) {
				validators = append(validators, v.Addr)
				count++
			}
		}
	}
	return nil, validators
}

// TotalStakeBalance Calculate the total amount of the stake account
func (vl *ValidatorList) TotalStakeBalance() *big.Int {
	var total = big.NewInt(0)
	for _, voter := range vl.Validators {
		total.Add(total, voter.Balance)
	}
	return total
}

// StakeBalance Returns the amount of the staked node
func (vl *ValidatorList) StakeBalance(address common.Address) *big.Int {
	for _, st := range vl.Validators {
		if st.Address().Hex() != address.Hex() && st.Proxy.Hex() != address.Hex() {
			continue
		}
		return st.Balance
	}
	return big.NewInt(0)
}

func (vl *ValidatorList) ConvertToAddress() (addrs []common.Address) {
	for _, validator := range vl.Validators {
		addrs = append(addrs, validator.Addr)
	}
	return
}

func (vl *ValidatorList) ConvertToBigInt(validators []*Validator) (bigIntSlice []*big.Int) {
	if len(validators) == 0 {
		return
	}
	for _, m := range validators {
		bigIntSlice = append(bigIntSlice, m.Addr.Hash().Big())
	}
	return
}

func (vl *ValidatorList) Exist(addr common.Address) bool {
	for _, v := range vl.Validators {
		if v.Addr == addr || v.Proxy == addr {
			return true
		}
	}
	return false
}

func (vl *ValidatorList) ExistProxy(addr common.Address) bool {
	emptyAddr := common.Address{}
	for _, v := range vl.Validators {
		if v.Addr == addr && v.Proxy != emptyAddr {
			return true
		}
	}
	return false
}

func (vl *ValidatorList) GetProxy(delegate common.Address) (common.Address, bool) {
	for _, v := range vl.Validators {
		if v.Addr == delegate {
			return v.Proxy, true
		}
	}
	return common.Address{}, false
}

func (vl *ValidatorList) ExistAdderRange(addr common.Address) bool {
	for _, v := range vl.Validators {
		if (v.Addr == addr || v.Proxy == addr) && v.Weight != nil {
			return true
		}
	}
	return false
}

func (vl *ValidatorList) AddAddrRange(addr common.Address, weight []*big.Int) {
	for _, v := range vl.Validators {
		if v.Addr == addr || v.Proxy == addr {
			v.Weight = weight
		}
	}
}

func (vl *ValidatorList) GetByAddress(addr common.Address) int {
	for i, v := range vl.Validators {
		if v.Addr == addr || v.Proxy == addr {
			return i
		}
	}
	return -1
}

func (vl *ValidatorList) GetByIndex(i uint64) Validator {
	if i >= uint64(vl.Len()) {
		return Validator{}
	}
	return *vl.Validators[i]
}

func (vl *ValidatorList) GetValidatorByAddr(addr common.Address) *Validator {
	for _, v := range vl.Validators {
		if v.Addr == addr {
			return v
		}
		continue
	}
	return &Validator{}
}
