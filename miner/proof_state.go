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

func (psp *ProofStatePool) Put(height *big.Int, proposer, validator common.Address, vl *types.ValidatorList) bool {
	psp.mu.Lock()
	defer psp.mu.Unlock()
	for k, p := range psp.proofs {
		if k == height.Uint64() {
			// Proof data exists for this height
			if p.onlineValidator.Has(validator) {
				return false
			}
			p.onlineValidator = append(p.onlineValidator, validator)
			p.receiveValidatorsSum = new(big.Int).Add(p.receiveValidatorsSum, vl.StakeBalance(validator))
			p.count++
			return true
		}
	}
	// No proof data exists for this height
	ps := newProofState(height, nil, nil, false, proposer, []byte{}, [][]byte{})
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
	validatorList        *types.ValidatorList
	targetWeightBalance  *big.Int
	empty                bool
	receiveValidatorsSum *big.Int
	proposer             common.Address
	proposerMessage      []byte
	onlineValidator      OnlineValidator // The highly online validator of this block & reward addr
	onlineMessages       [][]byte
	emptyBlockMessages   [][]byte
}

func newProofState(height, targetWeightBalance *big.Int, validatorList *types.ValidatorList, empty bool, proposer common.Address, proposerMessage []byte, onlineMessage [][]byte) *ProofState {
	vals := make(OnlineValidator, 0)
	emptyMessage := make([][]byte, 0)
	return &ProofState{
		count:               0,
		height:              height,
		validatorList:       validatorList,
		empty:               empty,
		targetWeightBalance: targetWeightBalance,
		proposer:            proposer,
		proposerMessage:     proposerMessage,
		onlineValidator:     vals,
		onlineMessages:      onlineMessage,
		emptyBlockMessages:  emptyMessage,
	}
}

func (ps ProofState) GetAllAddress(validators *types.ValidatorList) []common.Address {
	addrs := make([]common.Address, len(ps.onlineValidator)+1)
	addrs[0] = validators.GetValidatorAddr(ps.proposer)
	for i, val := range ps.onlineValidator {
		addrs[i+1] = validators.GetValidatorAddr(val)
	}

	return addrs
}

func (ps ProofState) GetAllMessage() [][]byte {
	emptyMessages := make([][]byte, 1)
	emptyMessages[0] = ps.proposerMessage
	return append(emptyMessages, ps.emptyBlockMessages...)
}

type OnlineValidator []common.Address

func (ov OnlineValidator) Has(addr common.Address) bool {
	for _, val := range ov {
		if val == addr {
			return true
		}
	}
	return false
}

//func (ov OnlineValidator) Add(addr common.Address) {
//	ov[addr] = struct{}{}
//}
//
//func (ov OnlineValidator) Delete(addr common.Address) {
//	delete(ov, addr)
//}
