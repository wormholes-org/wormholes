package backend

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

var (
	// errPeerSetClosed is returned if a peer is attempted to be added or removed
	// from the peer set after it has been terminated.
	errPeerSetClosed = errors.New("peerset closed")

	// errPeerAlreadyRegistered is returned if a peer is attempted to be added
	// to the peer set, but one with the same id already exists.
	errPeerAlreadyRegistered = errors.New("peer already registered")

	// errPeerNotRegistered is returned if a peer is attempted to be removed from
	// a peer set, but no peer with the given id exists.
	errPeerNotRegistered = errors.New("peer not registered")
)

type peer struct {
	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap
}

type peerSet struct {
	peers  map[common.Address]*peer
	lock   sync.RWMutex
	closed bool
}

func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[common.Address]*peer),
	}
}

func (ps *peerSet) registerPeer(peer *peer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errPeerSetClosed
	}
	id := crypto.PubkeyToAddress(*peer.Node().Pubkey())
	if _, ok := ps.peers[id]; ok {
		return errPeerAlreadyRegistered
	}

	ps.peers[id] = peer
	return nil
}

func (ps *peerSet) unregisterPeer(id common.Address) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	_, ok := ps.peers[id]
	if !ok {
		return errPeerNotRegistered
	}
	delete(ps.peers, id)
	return nil
}

func (ps *peerSet) peer(id common.Address) *peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

func (ps *peerSet) peersWithinList(list []istanbul.Validator) map[common.Address]*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	peers := make(map[common.Address]*peer)
	for _, val := range list {
		if p := ps.peers[val.Address()]; p != nil {
			peers[val.Address()] = p
		}
	}
	return peers
}

func (ps *peerSet) len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// close disconnects all peers.
func (ps *peerSet) close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
