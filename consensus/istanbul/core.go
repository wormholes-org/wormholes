package istanbul

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Core interface {
	Start() error
	Stop() error
	IsProposer() bool

	// verify if a hash is the same as the proposed block in the current pending request
	//
	// this is useful when the engine is currently the proposer
	//
	// pending request is populated right at the preprepare stage so this would give us the earliest verification
	// to avoid any race condition of coming propagated blocks
	IsCurrentProposal(blockHash common.Hash) bool

	RoundInfo() []string

	ConsensusInfo() chan map[string]interface{}

	GetOnlineValidators() map[uint64]*types.OnlineValidatorList
	GetOnlineProofsMu() *sync.Mutex
}
