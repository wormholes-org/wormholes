package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p"
)

type Broadcaster interface {
	//FindPeerSet(targets map[common.Address]bool) map[common.Address]Peer
	FindPeerSet() map[string]Peer
	EmptyResponse() chan string
}

type Handler interface {
	SetBroadcaster(Broadcaster)
	// HandleMsg handles a message from peer
	HandleMsg(address common.Address, data p2p.Msg) (bool, error)

	GetWorker() *worker
}

type Peer interface {
	SendWorkerMsg(msgCode uint64, data interface{}) error
	RequestEmptyMsg(msg []byte) int
}
