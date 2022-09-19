// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package backend

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	lru "github.com/hashicorp/golang-lru"
)

const (
	NewBlockMsg = 0x07
	istanbulMsg = 0x11
)

var (
	// errDecodeFailed is returned when decode message fails
	errDecodeFailed = errors.New("fail to decode istanbul message")

	// errAddSelf is returned when add self peer
	errAddSelf = errors.New("unable to add self peer")

	// errNotValidator is returned when peer is not validator
	errNotValidator = errors.New("peer is not validator")
)

// Protocol implements consensus.Engine.Protocol
func (sb *Backend) Protocol() consensus.Protocol {
	return consensus.IstanbulProtocol
}

func (sb *Backend) decode(msg p2p.Msg) ([]byte, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errDecodeFailed
	}
	return data, istanbul.RLPHash(data), nil
}

func (sb *Backend) handlePeer(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	if sb.server.Self().ID() == p.ID() {
		return errAddSelf
	}
	if !sb.ibftEngine.IsValidator(crypto.PubkeyToAddress(*p.Node().Pubkey())) {
		return errNotValidator
	}
	if err := sb.peers.registerPeer(&peer{p, rw}); err != nil {
		return err
	}
	addr := crypto.PubkeyToAddress(*p.Node().Pubkey())
	for {
		if err := sb.handleMsg(addr, rw); err != nil {
			return err
		}
	}
}

func (sb *Backend) handleMsg(addr common.Address, rw p2p.MsgReadWriter) error {
	// Read the next message from the remote peer (in protoRW), and ensure it's fully consumed
	msg, err := rw.ReadMsg()
	if err != nil {
		return err
	}
	defer msg.Discard()

	// See if the consensus engine protocol can handle this message, e.g. istanbul will check for message is
	// istanbulMsg = 0x11, and NewBlockMsg = 0x07.
	handled, err := sb.HandleMsg(addr, msg)
	sb.logger.Debug("consensus|receiveMsg", "from", addr, "code", msg.Code, "handled", handled, "err", err)

	return err
}

// HandleMsg implements consensus.Handler.HandleMsg
func (sb *Backend) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if msg.Code == istanbulMsg {
		if !sb.coreStarted {
			sb.logger.Info("caver|HandleMsg|ErrStoppedEngine", "!sb.coreStarted", !sb.coreStarted)
			return true, nil
		}

		data, hash, err := sb.decode(msg)
		if err != nil {
			return true, errDecodeFailed
		}
		// Mark peer's message
		ms, ok := sb.recentMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
		} else {
			m, _ = lru.NewARC(inmemoryMessages)
			sb.recentMessages.Add(addr, m)
		}
		m.Add(hash, true)

		// Mark self known message
		if _, ok := sb.knownMessages.Get(hash); ok {
			return true, nil
		}
		sb.knownMessages.Add(hash, true)

		go sb.istanbulEventMux.Post(istanbul.MessageEvent{
			Code:    msg.Code,
			Payload: data,
		})
		return true, nil
	}
	//https://github.com/ConsenSys/quorum/pull/539
	//https://github.com/ConsenSys/quorum/issues/389
	if msg.Code == NewBlockMsg && sb.core != nil && sb.core.IsProposer() { // eth.NewBlockMsg: import cycle
		// this case is to safeguard the race of similar block which gets propagated from other node while this node is proposing
		// as p2p.Msg can only be decoded once (get EOF for any subsequence read), we need to make sure the payload is restored after we decode it
		sb.logger.Debug("BFT: received NewBlockMsg", "size", msg.Size, "payload.type", reflect.TypeOf(msg.Payload), "sender", addr)
		if reader, ok := msg.Payload.(*bytes.Reader); ok {
			payload, err := ioutil.ReadAll(reader)
			if err != nil {
				return true, err
			}
			reader.Reset(payload)       // ready to be decoded
			defer reader.Reset(payload) // restore so main eth/handler can decode
			var request struct {        // this has to be same as eth/protocol.go#newBlockData as we are reading NewBlockMsg
				Block *types.Block
				TD    *big.Int
			}
			if err := msg.Decode(&request); err != nil {
				sb.logger.Error("BFT: unable to decode the NewBlockMsg", "error", err)
				return false, nil
			}
			newRequestedBlock := request.Block
			if newRequestedBlock.Header().MixDigest == types.IstanbulDigest && sb.core.IsCurrentProposal(newRequestedBlock.Hash()) {
				sb.logger.Debug("BFT: block already proposed", "hash", newRequestedBlock.Hash(), "sender", addr)
				return true, nil
			}
		}
	}
	return false, nil
}

// SetBroadcaster implements consensus.Handler.SetBroadcaster
func (sb *Backend) SetBroadcaster(broadcaster consensus.Broadcaster) {
	sb.broadcaster = broadcaster
}

func (sb *Backend) NewChainHead() error {
	sb.coreMu.RLock()
	defer sb.coreMu.RUnlock()
	if !sb.coreStarted {
		sb.logger.Info("caver|NewChainHead|ErrStoppedEngine", "!sb.coreStarted", !sb.coreStarted)
		return istanbul.ErrStoppedEngine
	}
	go sb.istanbulEventMux.Post(istanbul.FinalCommittedEvent{})
	return nil
}
