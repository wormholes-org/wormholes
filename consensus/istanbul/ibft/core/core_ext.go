package core

import (
	"github.com/ethereum/go-ethereum/common"
)

func (c *core) RandomSelectBonus(allocBonus []common.Address, seed int64) {
	//get validators
	//state.CreateNFTByOfficial16(allocBonus, istanbulExtra.ExchangerAddr, 0)
	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	//header.Root = state.IntermediateRoot(chain.Config().IsEIP158(0))
	//set validators
}
