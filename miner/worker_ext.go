package miner

import (
	"fmt"
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
