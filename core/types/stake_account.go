package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type Staker struct {
	Addr    common.Address
	Balance *big.Int
}

func (s *Staker) Address() common.Address {
	return s.Addr
}

func NewStaker(addr common.Address, balance *big.Int) *Staker {
	return &Staker{Addr: addr, Balance: balance}
}

type StakerList struct {
	Stakers []*Staker
}

type DBStaker struct {
	Staker
	DeleteFlag bool
}

type DBStakerList struct {
	DBStakers []*DBStaker
}

func (sl *StakerList) Len() int {
	return len(sl.Stakers)
}

func (sl *StakerList) Less(i, j int) bool {
	return sl.Stakers[i].Address().Hash().Big().Cmp(sl.Stakers[j].Address().Hash().Big()) < 0
}

func (sl *StakerList) Swap(i, j int) {
	sl.Stakers[i], sl.Stakers[j] = sl.Stakers[j], sl.Stakers[i]
}

func (sl *StakerList) AddStaker(addr common.Address, balance *big.Int) bool {
	for _, v := range sl.Stakers {
		if v.Address() == addr {
			v.Balance.Add(v.Balance, balance)
			sort.Sort(sl)
			return true
		}
	}
	sl.Stakers = append(sl.Stakers, NewStaker(addr, balance))
	sort.Sort(sl)
	return true
}

func (sl *StakerList) RemoveStaker(addr common.Address, balance *big.Int) bool {
	for i, v := range sl.Stakers {
		if v.Address() == addr {
			if v.Balance.Cmp(balance) > 0 {
				v.Balance.Sub(v.Balance, balance)
				sort.Sort(sl)
				return true
			} else {
				v.Balance.Sub(v.Balance, balance)
				sl.Stakers = append(sl.Stakers[:i], sl.Stakers[i+1:]...)
				sort.Sort(sl)
				return true
			}
		}
	}
	log.Warn("stake_account|RemoveStaker:address not exist", "addr", addr)
	return false
}

// Query K validators closest to random numbers based on distance and pledge amount
func (sl *StakerList) ValidatorByDistanceAndWeight(addr []*big.Int, k int, randomHash common.Hash) []common.Address {

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
			a := new(big.Int).Mul(sub1, sl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, sl.TotalStakeBalance())
			addrToWeightMap[v] = w
		} else {
			a := new(big.Int).Mul(sub2, sl.StakeBalance(common.BigToAddress(v)))
			w := new(big.Int).Div(a, sl.TotalStakeBalance())
			addrToWeightMap[v] = w
		}
	}

	// Sort by weight in descending order
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
func (sl *StakerList) TotalStakeBalance() *big.Int {
	var total = big.NewInt(0)
	for _, voter := range sl.Stakers {
		total.Add(total, voter.Balance)
	}
	return total
}

// Returns the amount of the staked node
func (sl *StakerList) StakeBalance(address common.Address) *big.Int {
	for _, st := range sl.Stakers {
		if st.Address() != address {
			continue
		}
		return st.Balance
	}
	return big.NewInt(0)
}

func (sl *StakerList) DeepCopy() *StakerList {
	tempStakerList := &StakerList{
		Stakers: make([]*Staker, 0, sl.Len()),
	}

	for _, staker := range sl.Stakers {
		tempStaker := Staker{
			Addr:    staker.Addr,
			Balance: new(big.Int).Set(staker.Balance),
		}
		tempStakerList.Stakers = append(tempStakerList.Stakers, &tempStaker)
	}

	return tempStakerList
}

func (sl *StakerList) selectAddress(rand *big.Int) (common.Address, error) {
	sum := big.NewInt(0)
	for _, staker := range sl.Stakers {
		sum.Add(sum, staker.Balance)
		if sum.Cmp(rand) >= 0 {
			return staker.Addr, nil
		}
	}

	return common.Address{}, errors.New("select address error")
}

func (sl *StakerList) SelectRandom4Address(num int, hash []byte) ([]common.Address, error) {
	var random4Address []common.Address
	tempStakers := sl.DeepCopy()
	//hsh256 := sha256.New()

	// if number of all stakers is less the selected staker number,
	// return all stakers
	if num > len(tempStakers.Stakers) {
		for _, s := range tempStakers.Stakers {
			random4Address = append(random4Address, s.Addr)
		}
	} else {
		for i := 0; i < num; i++ {
			total := tempStakers.TotalStakeBalance()
			//hash = hsh256.Sum(hash)
			hash = crypto.Keccak256(hash)
			mod := new(big.Int).Mod(new(big.Int).SetBytes(hash), total)
			address, err := tempStakers.selectAddress(mod)
			if err != nil {
				return nil, err
			}
			random4Address = append(random4Address, address)
			tempStakers.RemoveStaker(address, new(big.Int).Set(tempStakers.StakeBalance(address)))
		}
	}

	return random4Address, nil
}

func (sl *StakerList) GetByIndex(i uint64) *Staker {
	if i >= uint64(sl.Len()) {
		return &Staker{}
	}
	return sl.Stakers[i]
}

func rankByWordCount(wordFrequencies map[*big.Int]*big.Int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   *big.Int
	Value *big.Int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value.Cmp(p[j].Value) < 0 }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func HashToBigIntSilce(addrs []common.Address) (bigIntSlice []*big.Int) {
	for _, addr := range addrs {
		bigIntSlice = append(bigIntSlice, addr.Hash().Big())
	}
	return
}

// Address in ascending order
func SortAddr(addrs []common.Address) []*big.Int {
	bigIntSlice := make(BigIntSlice, 0)
	for _, addr := range addrs {
		bigIntSlice = append(bigIntSlice, addr.Hash().Big())
	}
	sort.Sort(bigIntSlice)
	return bigIntSlice
}

type BigIntSlice []*big.Int

func (m BigIntSlice) Len() int {
	return len(m)
}

func (m BigIntSlice) Less(i, j int) bool {
	return m[i].Cmp(m[j]) < 0
}

func (m BigIntSlice) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
