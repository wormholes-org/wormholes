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

package core

import (
	"bytes"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	metrics "github.com/ethereum/go-ethereum/metrics"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	roundMeter     = metrics.NewRegisteredMeter("consensus/istanbul/core/round", nil)
	sequenceMeter  = metrics.NewRegisteredMeter("consensus/istanbul/core/sequence", nil)
	consensusTimer = metrics.NewRegisteredTimer("consensus/istanbul/core/consensus", nil)
)

// New creates an Istanbul consensus core
func New(backend istanbul.Backend, config *istanbul.Config) *core {
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              ibfttypes.StateAcceptRequest,
		handlerWg:          new(sync.WaitGroup),
		logger:             log.New("address", backend.Address()),
		backend:            backend,
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},
	}

	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *istanbul.Config
	address common.Address
	state   ibfttypes.State
	logger  log.Logger

	backend               istanbul.Backend
	events                *event.TypeMuxSubscription
	finalCommittedSub     *event.TypeMuxSubscription
	timeoutSub            *event.TypeMuxSubscription
	futurePreprepareTimer *time.Timer

	valSet                istanbul.ValidatorSet
	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error)

	backlogs   map[common.Address]*prque.Prque
	backlogsMu *sync.Mutex

	current   *roundState
	handlerWg *sync.WaitGroup

	roundChangeSet   *roundChangeSet
	roundChangeTimer *time.Timer

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	consensusTimestamp time.Time
}

func (c *core) finalizeMessage(msg *ibfttypes.Message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.Address()

	// Assign the CommittedSeal if it's a COMMIT message and proposal is not nil
	if msg.Code == ibfttypes.MsgCommit && c.current.Proposal() != nil {
		msg.CommittedSeal = []byte{}
		seal := PrepareCommittedSeal(c.current.Proposal().Hash())
		// Add proof of consensus
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(msg *ibfttypes.Message) {
	logger := c.logger.New("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(c.valSet, msg.Code, payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) currentView() *istanbul.View {
	return &istanbul.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) IsProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
}

func (c *core) IsCurrentProposal(blockHash common.Hash) bool {
	return c.current != nil && c.current.pendingRequest != nil && c.current.pendingRequest.Proposal.Hash() == blockHash
}

func (c *core) commit() {
	c.setState(ibfttypes.StateCommitted)

	proposal := c.current.Proposal()
	if proposal != nil {
		committedSeals := make([][]byte, c.current.Commits.Size())
		for i, v := range c.current.Commits.Values() {
			committedSeals[i] = make([]byte, types.IstanbulExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}
		log.Info("carver|commit|backend", "no", proposal.Number().Uint64())
		if err := c.backend.Commit(proposal, committedSeals, big.NewInt(-1)); err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			log.Error("carver|commit|sendNextRoundChange", "no", proposal.Number().Uint64(),
				"hash", proposal.Hash().Hex(), "err", err.Error())
			c.sendNextRoundChange()
			return
		}
	}
}

// startNewRound starts a new round. if round equals to 0, it means to starts a new sequence
func (c *core) startNewRound(round *big.Int) {
	var logger log.Logger
	if c.current == nil {
		logger = c.logger.New("old_round", -1, "old_seq", 0)
	} else {
		logger = c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence())
	}

	logger.Trace("Start new ibft round")

	roundChange := false
	// Try to get last proposal
	lastProposal, lastProposer := c.backend.LastProposal()
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), c.current.Sequence())
		sequenceMeter.Mark(new(big.Int).Add(diff, common.Big1).Int64())

		if !c.consensusTimestamp.IsZero() {
			consensusTimer.UpdateSince(c.consensusTimestamp)
			c.consensusTimestamp = time.Time{}
		}
		logger.Trace("Catch up latest proposal", "number", lastProposal.Number().Uint64(), "hash", lastProposal.Hash())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Sequence().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			// same seq and round, don't need to start new round
			return
		} else if round.Cmp(c.current.Round()) < 0 {
			logger.Warn("New round should not be smaller than current round", "seq", lastProposal.Number().Int64(), "new_round", round, "old_round", c.current.Round())
			return
		}
		roundChange = true
	} else {
		logger.Warn("New sequence should be larger than current sequence", "new_seq", lastProposal.Number().Int64())
		return
	}

	var newView *istanbul.View
	if roundChange {
		log.Info("caver|startNewRound|roundChange=true", "currentNo", lastProposal.Number().Uint64(), "round", round)
		newView = &istanbul.View{
			Sequence: new(big.Int).Set(c.current.Sequence()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		log.Info("caver|startNewRound|roundChange=false", "currentNo", lastProposal.Number().Uint64(), "round", round)
		newView = &istanbul.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		// 以当前链最新高度的哈希计算validator 与当前矿工正在prepare执行计算的validator是一致的
		c.valSet = c.backend.Validators(lastProposal)
		if c.valSet == nil {
			log.Error("startNewRound err : c.valSet == nil")
			return
		}
	}

	// If new round is 0, then check if qbftConsensus needs to be enabled
	if round.Uint64() == 0 && c.backend.IsQBFTConsensusAt(newView.Sequence) {
		logger.Trace("Starting qbft consensus as qbftBlock has passed")
		if err := c.backend.StartQBFTConsensus(); err != nil {
			// If err is returned, then QBFT consensus is started for the next block
			logger.Error("Unable to start QBFT Consensus, retrying for the next block", "error", err)
		}
		return
	}

	// Update logger
	logger = logger.New("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())

	log.Info("startNewRound|proposer", "c.valSet.List()", len(c.valSet.List()))
	for _, v := range c.valSet.List() {
		log.Info("startNewRound|proposer", "valSet", v.String())
	}

	log.Info("caver|startNewRound|proposer", "proposer ", c.valSet.GetProposer())
	log.Info("caver|startNewRound|proposer", "no", newView.Sequence.String(),
		"round", newView.Round.String(), "proposer", c.valSet.GetProposer().Address().String())

	// temp print validator
	for i, v := range c.valSet.List() {
		log.Info("caver|startNewRound|validator", "no", newView.Sequence.String(),
			"round", newView.Round.String(), "i", i, "addr", v.Address().String())
	}

	c.waitingForRoundChange = false
	// 在这个阶段如果本地已经发出一个request事件，顺便把这个事件发出去，让其他节点处理，自己就进入到preprepare阶段
	c.setState(ibfttypes.StateAcceptRequest)
	if roundChange && c.IsProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			log.Info("caver|c.current.IsHashLocked()", "currentProposal", c.current.Proposal().Number().Uint64())
			r := &istanbul.Request{
				Proposal: c.current.Proposal(), //c.current.Proposal would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendPreprepare(r)
		} else if c.current.pendingRequest != nil {
			log.Info("caver|c.current.pendingRequest != nil", "currentPendingRequest", c.current.pendingRequest.Proposal.Number().Uint64())
			c.sendPreprepare(c.current.pendingRequest)
		}
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
}

func (c *core) catchUpRound(view *istanbul.View) {

	log.Info("caver|catchUpRound", "sequence", view.Sequence.Uint64(), "round", view.Round.Uint64())
	logger := c.logger.New("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())

	if view.Round.Cmp(c.current.Round()) > 0 {
		roundMeter.Mark(new(big.Int).Sub(view.Round, c.current.Round()).Int64())
	}
	c.waitingForRoundChange = true

	// Need to keep block locked for round catching up
	c.updateRoundState(view, c.valSet, true)
	c.roundChangeSet.Clear(view.Round)
	c.newRoundChangeTimer()

	logger.Trace("Catch up round", "new_round", view.Round, "new_seq", view.Sequence, "new_proposer", c.valSet)
}

// updateRoundState updates round state by checking if locking block is necessary
func (c *core) updateRoundState(view *istanbul.View, validatorSet istanbul.ValidatorSet, roundChange bool) {
	// Lock only if both roundChange is true and it is locked
	if roundChange && c.current != nil {
		if c.current.IsHashLocked() {
			log.Info("carver|updateRoundState|1", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
			c.current = newRoundState(view, validatorSet, c.current.GetLockedHash(), c.current.Preprepare, c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			log.Info("carver|updateRoundState|2", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
			c.current = newRoundState(view, validatorSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		log.Info("carver|updateRoundState|3", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
		c.current = newRoundState(view, validatorSet, common.Hash{}, nil, nil, c.backend.HasBadProposal)
	}
}

func (c *core) setState(state ibfttypes.State) {
	if c.state != state {
		c.state = state
	}
	if state == ibfttypes.StateAcceptRequest {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) Address() common.Address {
	return c.address
}

func (c *core) stopFuturePreprepareTimer() {
	if c.futurePreprepareTimer != nil {
		c.futurePreprepareTimer.Stop()
	}
}

func (c *core) stopTimer() {
	c.stopFuturePreprepareTimer()
	if c.roundChangeTimer != nil {
		c.roundChangeTimer.Stop()
	}
}

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	// set timeout based on the round number
	timeout := time.Duration(c.config.RequestTimeout) * time.Millisecond
	round := c.current.Round().Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
		log.Info("newRoundChangeTimer : timeout", "no", c.current.sequence, "round", c.current.round, "timeout", timeout.String())
	}
	c.roundChangeTimer = time.AfterFunc(timeout, func() {
		c.sendEvent(timeoutEvent{})
	})
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return istanbul.CheckValidatorSignature(c.valSet, data, sig)
}

func (c *core) QuorumSize() int {
	return (2 * c.valSet.F()) + 1
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(ibfttypes.MsgCommit)})
	return buf.Bytes()
}

func (c *core) RoundInfo() (roundInfo []string) {
	rs := c.current
	roundInfo = append(roundInfo, rs.round.String())
	roundInfo = append(roundInfo, rs.sequence.String())
	return
}
