package miner

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type ProofStatePool struct {
	mu     sync.Mutex
	proofs map[uint64]*ProofState // Online attestation status of current altitude
}

func NewProofStatePool() *ProofStatePool {
	return &ProofStatePool{proofs: make(map[uint64]*ProofState)}
}

// ClearPrev Clear all proof data before this altitude
func (p *ProofStatePool) ClearPrev(height *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var removeHeight []uint64

	for k := range p.proofs {
		if k <= height.Uint64() {
			removeHeight = append(removeHeight, k)
		}
	}

	for _, v := range removeHeight {
		delete(p.proofs, v)
	}
}

func (p *ProofStatePool) SetTargetWeightBalance(height *big.Int, weight *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.proofs[height.Uint64()]; !ok {
		ps := newProofState(common.Address{}, common.Address{})
		ps.targetWeightBalance = new(big.Int).Set(weight)
		p.proofs[height.Uint64()] = ps
	} else {
		p.proofs[height.Uint64()].targetWeightBalance = new(big.Int).Set(weight)
	}
}

func (p *ProofStatePool) SetValidatorList(height *big.Int, list *types.ValidatorList) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.proofs[height.Uint64()]; !ok {
		ps := newProofState(common.Address{}, common.Address{})
		ps.validatorList = list
		p.proofs[height.Uint64()] = ps
	} else {
		p.proofs[height.Uint64()].validatorList = list
	}
}

func (p *ProofStatePool) UpdateNextIndex(height *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if ps, ok := p.proofs[height.Uint64()]; ok {
		if ps.validatorList == nil || len(ps.validatorList.Validators) <= 0 {
			return
		}
		ps.nextIndex = (ps.nextIndex + 1) % len(ps.validatorList.Validators)
	}
}

func (psp *ProofStatePool) Put(height *big.Int, proposer, validator common.Address, vl *types.ValidatorList) bool {
	psp.mu.Lock()
	defer psp.mu.Unlock()
	for k, p := range psp.proofs {
		if k == height.Uint64() {
			// Proof data exists for this height
			if p.onlineValidator.Has(validator) {
				return false
			}
			p.onlineValidator.Add(validator)
			p.receiveValidatorsSum = new(big.Int).Add(p.receiveValidatorsSum, vl.StakeBalance(validator))
			p.count++
			return true
		}
	}
	// No proof data exists for this height
	ps := newProofState(proposer, validator)
	psp.proofs[height.Uint64()] = ps
	ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, vl.StakeBalance(validator))
	ps.count++
	return true
}

func (psp *ProofStatePool) GetProofCountByHeight(height *big.Int) int {
	psp.mu.Lock()
	defer psp.mu.Unlock()
	for h, v := range psp.proofs {
		if h == height.Uint64() {
			return v.count
		}
	}
	return -1
}

type ProofState struct {
	count                int // Represents the number of proofs collected
	height               *big.Int
	receiveValidatorsSum *big.Int
	proposer             common.Address
	onlineValidator      OnlineValidator // The highly online validator of this block & reward addr
	emptyBlockMessages   [][]byte

	validatorList       *types.ValidatorList
	targetWeightBalance *big.Int
	nextIndex           int
}

func newProofState(proposer, validator common.Address) *ProofState {
	emptyAddress := common.Address{}
	vals := make(OnlineValidator)
	if validator != emptyAddress {
		vals.Add(validator)
	}

	return &ProofState{count: 0, proposer: proposer, onlineValidator: vals}
}

type OnlineValidator map[common.Address]struct{}

func (ov OnlineValidator) Has(addr common.Address) bool {
	_, ok := ov[addr]
	return ok
}

func (ov OnlineValidator) Add(addr common.Address) {
	ov[addr] = struct{}{}
}

func (ov OnlineValidator) Delete(addr common.Address) {
	delete(ov, addr)
}

func (ov OnlineValidator) GetAllAddress() []common.Address {
	var addrs []common.Address
	for address, _ := range ov {
		addrs = append(addrs, address)
	}
	return addrs
}
