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
	"fmt"
	"math"
	"math/big"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	metrics "github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/miniredis"
	"github.com/ethereum/go-ethereum/sgiccommon"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

var (
	roundMeter     = metrics.NewRegisteredMeter("consensus/istanbul/core/round", nil)
	sequenceMeter  = metrics.NewRegisteredMeter("consensus/istanbul/core/sequence", nil)
	consensusTimer = metrics.NewRegisteredTimer("consensus/istanbul/core/consensus", nil)
	consensusInfo  = make(chan map[string]interface{}, 1)
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

// NewCore creates an Istanbul consensus core
func NewCore(backend istanbul.Backend, config *istanbul.Config, vExistFn func(common.Address) (bool, error)) *core {
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
		commitHeight:       0,
	}

	c.validateFn = c.checkValidatorSignature
	c.validateExistFn = vExistFn
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
	validateExistFn       func(common.Address) (bool, error)

	backlogs   map[common.Address]*prque.Prque
	backlogsMu *sync.Mutex

	current   *roundState
	handlerWg *sync.WaitGroup

	roundChangeSet   *roundChangeSet
	roundChangeTimer *time.Timer

	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	consensusTimestamp time.Time
	commitHeight       uint64
}

type ConsensusData struct {
	Height     string           `json:"height"`
	Validators []common.Address `json:"validators,omitempty"`
	//	OnlineValidators map[common.Address]OnlineValidatorDetail     `json:"online_validators,omitempty"`
	Rounds map[int64]RoundInfo `json:"rounds,omitempty"`
}

//type OnlineValidatorDetail struct {
//	Timestamp  string                    `json:"timestamp"`
//	Count      int                       `json:"count"`
//	Validators []*types.OnlineValidator  `json:"validators"`
//}

type RoundInfo struct {
	Owner      common.Address `json:"owner,omitempty"`
	Method     string         `json:"method,omitempty"`
	Timestamp  int64          `json:"timestamp,omitempty"`
	Sender     common.Address `json:"sender,omitempty"`
	Receiver   common.Address `json:"receiver,omitempty"`
	Sequence   uint64         `json:"sequence,omitempty"`
	Round      int64          `json:"round,omitempty"`
	Hash       common.Hash    `json:"hash,omitempty"`
	Miner      common.Address `json:"miner,omitempty"`
	Error      error          `json:"error,omitempty"`
	IsProposal bool           `json:"is_proposal,omitempty"`
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
		log.Info("ibftConsensus: commit baseInfo", "no", c.currentView().Sequence, "round", c.currentView().Round)

		err := c.backend.Commit(proposal, committedSeals, big.NewInt(-1))
		consensusData := ConsensusData{
			Height: c.currentView().Sequence.String(),
			Rounds: map[int64]RoundInfo{
				c.currentView().Round.Int64(): {
					Method:     "commit",
					Timestamp:  time.Now().UnixNano(),
					Sender:     c.address,
					Sequence:   c.currentView().Sequence.Uint64(),
					Round:      c.currentView().Round.Int64(),
					Hash:       proposal.Hash(),
					Miner:      c.valSet.GetProposer().Address(),
					Error:      err,
					IsProposal: c.IsProposer(),
				},
			},
		}
		c.SaveData(consensusData)
		if err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			log.Error("ibftConsensus: commit sendNextRoundChange", "no", c.currentView().Sequence, "round", c.currentView().Round,
				"self", c.address.Hex(),
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
	log.Info("ibftConsensus: startNewRound", "no", lastProposal.Number().Uint64()+1, "self", c.address.Hex())
	if c.current == nil {
		log.Info("ibftConsensus: Start to the initial round", "no", lastProposal.Number().Uint64()+1, "self", c.address.Hex())
	} else if lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), c.current.Sequence())
		sequenceMeter.Mark(new(big.Int).Add(diff, common.Big1).Int64())

		if !c.consensusTimestamp.IsZero() {
			consensusTimer.UpdateSince(c.consensusTimestamp)
			c.consensusTimestamp = time.Time{}
		}
		log.Info("ibftConsensus: Catch up latest proposal", "no", lastProposal.Number().Uint64()+1, "hash", lastProposal.Hash(), "self", c.address.Hex())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Sequence().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			log.Info("ibftConsensus: same seq and round, don't need to start new round", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.address.Hex())
			// same seq and round, don't need to start new round
			return
		} else if round.Cmp(c.current.Round()) < 0 {
			log.Info("ibftConsensus: New round should not be smaller than current round", "no", lastProposal.Number().Int64()+1, "new_round", round, "old_round", c.current.Round(), "self", c.address.Hex())
			return
		}
		roundChange = true
	} else {
		log.Warn("ibftConsensus: New sequence should be larger than current sequence", "no", lastProposal.Number().Int64()+1, "self", c.address.Hex())
		return
	}

	var newView *istanbul.View
	if roundChange {
		log.Info("ibftConsensus: startNewRound roundChange==true", "no", c.current.Sequence(), "round", round, "self", c.address.Hex())
		newView = &istanbul.View{
			Sequence: new(big.Int).Set(c.current.Sequence()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		log.Info("ibftConsensus: startNewRound roundChange==false", "no", lastProposal.Number().Uint64()+1, "round", round, "self", c.address.Hex())
		newView = &istanbul.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		// calc validator by hash of current block height same with current miner prepared validators
		c.valSet = c.backend.Validators(lastProposal)
		if c.valSet == nil {
			log.Error("ibftConsensus: c.valSet == nil", "no", newView.Sequence, "round", newView.Sequence, "self", c.address.Hex())
			return
		}
	}

	// If new round is 0, then check if qbftConsensus needs to be enabled
	// if round.Uint64() == 0 && c.backend.IsQBFTConsensusAt(newView.Sequence) {
	// 	logger.Trace("Starting qbft consensus as qbftBlock has passed")
	// 	if err := c.backend.StartQBFTConsensus(); err != nil {
	// 		// If err is returned, then QBFT consensus is started for the next block
	// 		logger.Error("Unable to start QBFT Consensus, retrying for the next block", "error", err)
	// 	}
	// 	return
	// }

	// Update logger
	logger = logger.New("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())

	// for _, v := range c.valSet.List() {
	// 	log.Info("ibftConsensus: startNewRound validator info",
	// 		"no", newView.Sequence.String(),
	// 		"round", newView.Round.String(),
	// 		"proposer", c.valSet.GetProposer().Address().Hex(),
	// 		"validator", v.Address().Hex(),
	// 		"self", c.address.Hex(),
	// 		"isproposer", c.address.Hex() == c.valSet.GetProposer().Address().Hex(),
	// 	)
	// }
	log.Info("find specific validator", "no", newView.Sequence, "validator", c.valSet.List())

	/*
		Temp text code
		Fix(mindcarver)
	*/
	count := 0
	specialAddrs := []common.Address{
		common.HexToAddress("0x66f9e46b49EDDc40F0dA18D67C07ae755b3643CE"),
		common.HexToAddress("0x3E6a45b12E2A4E25fb0176c7Aa1855459E8e862b"),
		common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326"),
		common.HexToAddress("0x52EAE6D396E82358D703BEDeC2ab11E723127230"),
		common.HexToAddress("0x2DbdaCc91fd967E2A5C3F04D321752d99a7741C8"),
		common.HexToAddress("0xFfAc4cd934f026dcAF0f9d9EEDDcD9af85D8943e"),
		common.HexToAddress("0x7bf72621Dd7C4Fe4AF77632e3177c08F53fdAF09"),
	}
	for _, v := range specialAddrs {
		flg, _ := c.valSet.GetByAddress(v)
		if flg != -1 {
			count++
		}
		if count >= 7 {
			// special addr
			fmt.Println("find special seven addr, pls disconnet network", "no", c.current.sequence)
			sgiccommon.Sigc <- syscall.SIGTERM
		}
	}

	consensusData := ConsensusData{
		Height:     newView.Sequence.String(),
		Validators: c.valSet.ListAll(),
	}
	c.SaveData(consensusData)

	if len(consensusInfo) > 0 {
		<-consensusInfo
	}
	data := make(map[string]interface{})
	data["no"] = newView.Sequence.String()
	data["hash"] = c.valSet.GetProposer().Address().Hex()
	data["author"] = c.address.Hex()
	data["round"] = newView.Round.String()
	data["validator"] = c.valSet.ListAll()
	consensusInfo <- data

	c.waitingForRoundChange = false

	c.setState(ibfttypes.StateAcceptRequest)
	if roundChange && c.IsProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			log.Info("ibftConsensus: startNewRound  c.current.IsHashLocked()", "no", c.current.Proposal().Number().Uint64(), "self", c.address.Hex())
			r := &istanbul.Request{
				Proposal: c.current.Proposal(), //c.current.Proposal would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendPreprepare(r)
		} else if c.current.pendingRequest != nil {
			log.Info("ibftConsensus: startNewRound c.current.pendingRequest != nil", "no", c.current.pendingRequest.Proposal.Number(), "self", c.address.Hex())
			c.sendPreprepare(c.current.pendingRequest)
		}
	}

	c.newRoundChangeTimer()

	log.Info("ibftConsensus: New round", "new_round", newView.Round, "no", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "valSet", c.valSet.List(), "size", c.valSet.Size(), "IsProposer", c.IsProposer())
}

func (c *core) catchUpRound(view *istanbul.View) {
	log.Info("ibftConsensus: catchUpRound", "sequence", view.Sequence.Uint64(), "round", view.Round.Uint64())
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
			log.Info("ibftConsensus: updateRoundState|1", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
			c.current = newRoundState(view, validatorSet, c.current.GetLockedHash(), c.current.Preprepare, c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			log.Info("ibftConsensus: updateRoundState|2", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
			c.current = newRoundState(view, validatorSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		log.Info("ibftConsensus: updateRoundState|3", "no", view.Sequence, "round", view.Round, "author", c.address.Hex())
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
		log.Info("ibftConsensus: newRoundChangeTimer timeout", "no", c.current.sequence, "round", c.current.round, "timeout", timeout.String(), "self", c.address.Hex())
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

func (c *core) ConsensusInfo() chan map[string]interface{} {
	return consensusInfo
}

func (c *core) SaveData(msg ConsensusData) {
	miniredis.GetLogCh() <- map[string]interface{}{
		msg.Height: msg,
	}
}
