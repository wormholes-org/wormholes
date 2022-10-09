package miner

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

func (w *worker) prepre0Loop() {
	for {
		select {
		case bonusAddrs := <-w.engine.GetPrepre0Ch():
			fmt.Println(bonusAddrs)
			//w.current.state.CreateNFTByOfficial16(*bonusAddrs, *bonusAddrs, 0)
			//header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
			//header.UncleHash = nilUncleHash
			//w.snapshotState.SnapshotCommits()
		}
	}
}

func (w *worker) RewardOnline(validators []common.Address) {
	w.current.state.RewardOnline(validators, w.current.header.Number)
}
