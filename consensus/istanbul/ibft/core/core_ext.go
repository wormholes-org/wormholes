package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *core) RandomSelectBonus(allocBonus []common.Address) {
	//TODO: Find Finalize alloc Bonus
	prepre0Ch <- &allocBonus
	//TODO wait channel return block
	tmpBlk := <-prepre1Ch
	fmt.Println(tmpBlk)
	//get validators
	//state.CreateNFTByOfficial16(allocBonus, istanbulExtra.ExchangerAddr, 0)
	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	//header.Root = state.IntermediateRoot(chain.Config().IsEIP158(0))
	//set validators
}

func (c *core) GetPrepre0Ch() chan *[]common.Address {
	return prepre0Ch
}

func (c *core) GetPrepre1Ch() chan *types.Block {
	return prepre1Ch
}
