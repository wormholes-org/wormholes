package miner

import (
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type ProofStatePool struct {
	mu     sync.Mutex
	proofs map[*big.Int]*ProofState // Online attestation status of current altitude
}

func NewProofStatePool() *ProofStatePool {
	return &ProofStatePool{proofs: make(map[*big.Int]*ProofState)}
}

// ClearPrev Clear all proof data before this altitude
func (p *ProofStatePool) ClearPrev(height *big.Int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var removeHeight []*big.Int

	for k := range p.proofs {
		if k.Cmp(height) <= 0 {
			removeHeight = append(removeHeight, k)
		}
	}

	for _, v := range removeHeight {
		delete(p.proofs, v)
	}
}

func (psp *ProofStatePool) Put(height *big.Int, proposer, validator common.Address, vl *types.ValidatorList) bool {
	psp.mu.Lock()
	defer psp.mu.Unlock()
	for k, p := range psp.proofs {
		if k.Cmp(height) == 0 {
			// Proof data exists for this height
			if p.onlineValidator.Has(validator) {
				return false
			}
			p.onlineValidator.Add(validator)
			num, _ := strconv.Atoi(vl.StakeBalance(validator).String())
			p.receiveValidatorsSum += num
			p.count++
			return true
		}
	}
	// No proof data exists for this height
	ps := newProofState(proposer, validator)
	psp.proofs[height] = ps
	num, _ := strconv.Atoi(vl.StakeBalance(validator).String())
	ps.receiveValidatorsSum += num
	ps.count++
	return true
}

func (psp *ProofStatePool) GetProofCountByHeight(height *big.Int) int {
	psp.mu.Lock()
	defer psp.mu.Unlock()
	for h, v := range psp.proofs {
		if h.Cmp(height) == 0 {
			return v.count
		}
	}
	return -1
}

type ProofState struct {
	count                int // Represents the number of proofs collected
	receiveValidatorsSum int
	proposer             common.Address
	onlineValidator      OnlineValidator // The highly online validator of this block & reward addr
}

func newProofState(proposer, validator common.Address) *ProofState {
	vals := make(OnlineValidator)
	vals.Add(validator)
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
