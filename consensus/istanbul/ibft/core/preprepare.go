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
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
)

const (
	PreprepareStep1 uint64 = iota
	PreprepareStep2
)

var csssStat = PreprepareStep1                     //consensus state mark PreprepareStep1 = 0, PreprepareStep2 = 1
var randSeedMessages *messageSet = new(messageSet) // collected random data message

func (c *core) sendPreprepare(request *istanbul.Request) {
	if csssStat == PreprepareStep1 {
		c.sendPreprepareStep1(request)
		return
	}
	if csssStat == PreprepareStep2 {
		c.sendPreprepareStep2(request)
		return
	}
}

func (c *core) sendPreprepareStep1(request *istanbul.Request) {
	log.Info("ibftConsensus: sendPreprepareStep1 [csss]",
		"no", request.Proposal.Number(), "no", c.current.sequence,
		"round", c.current.round.Uint64(), "isproposer", c.IsProposer(), "slef", c.address.Hex())
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 {
		curView := c.currentView()
		log.Info("start collect random seed [csss]", "no", request.Proposal.Number(),
			"round", curView.Round,
			"author", c.address.Hex(),
			"hash", request.Proposal.Hash().Hex(),
			"self", c.address.Hex())
		if c.IsProposer() { //start collect random seed
			randSeedMessages = new(messageSet) //TODO: early clear
			c.broadcast(&ibfttypes.Message{
				Code: ibfttypes.MsgPreprepare,
				Msg:  []byte{},
			})
		} else { //send random seed
			//TODO: check necessary for send future random data
			c.BroadcastLocalRandomData()
		}
	}
}

func (c *core) BroadcastLocalRandomData() {
	csssStat = PreprepareStep2
	//TODO generate & send random Seed
	c.broadcast(&ibfttypes.Message{
		Code: ibfttypes.MsgPrepare,
		Msg:  c.localRandomBytes(),
	})
}

// local random address
func (c *core) localRandomBytes() []byte {
	rndSeed := time.Now().UnixNano()
	rand.Seed(rndSeed)
	bigSeed := big.NewInt(rand.Int63())
	//rndDat = common.BigToHash(bigSeed)
	return bigSeed.Bytes()
}

func (c *core) sendPreprepareStep2(request *istanbul.Request) {
	logger := c.logger.New("state", c.state)
	// If I'm the proposer and I have the same sequence with the proposal

	log.Info("ibftConsensus: sendPreprepareStep2",
		"no", request.Proposal.Number(), "no", c.current.sequence,
		"round", c.current.round.Uint64(), "isproposer", c.IsProposer(), "slef", c.address.Hex())
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 && c.IsProposer() {
		curView := c.currentView()
		preprepare, err := ibfttypes.Encode(&istanbul.Preprepare{
			View:     curView,
			Proposal: request.Proposal,
		})
		if err != nil {
			logger.Error("Failed to encode", "view", curView)
			return
		}

		log.Info("ibftConsensus: broadcast", "no", request.Proposal.Number(),
			"round", curView.Round,
			"author", c.address.Hex(),
			"hash", request.Proposal.Hash().Hex(),
			"self", c.address.Hex())
		c.broadcast(&ibfttypes.Message{
			Code: ibfttypes.MsgPreprepare,
			Msg:  preprepare,
		})
	}
}

func (c *core) handlePreprepare(msg *ibfttypes.Message, src istanbul.Validator) error {
	if csssStat == 0 {
		return c.handlePreprepareStep1(msg, src)
	} else {
		return c.handlePreprepareStep2(msg, src)
	}
	return nil
}

func (c *core) handlePreprepareStep1(msg *ibfttypes.Message, src istanbul.Validator) error {
	//logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare *istanbul.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return istanbulcommon.ErrFailedDecodePreprepare
	}

	log.Info("ibftConsensus: handlePreprepare",
		"no", preprepare.Proposal.Number().Uint64(),
		"round", preprepare.View.Round.String(),
		"from", src.Address().Hex(),
		"self", c.address.Hex(),
		"hash", preprepare.Proposal.Hash().Hex(),
		"state", c.state)
	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(ibfttypes.MsgPreprepare, preprepare.View); err != nil {
		log.Error("ibftConsensus: handlePreprepare checkMessage",
			"no", preprepare.Proposal.Number().Uint64(),
			"round", preprepare.View.Round.String(),
			"err", err.Error(),
			"hash", preprepare.Proposal.Hash().Hex(),
			"self", c.address.Hex())
		if err == istanbulcommon.ErrOldMessage {
			// Get validator set for the given proposal
			valSet := c.backend.ParentValidators(preprepare.Proposal).Copy()
			previousProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, preprepare.View.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(preprepare.Proposal.Hash(), preprepare.Proposal.Number()) {
				log.Info("ibftConsensus: handlePreprepare sendCommitForOldBlock",
					"no", preprepare.Proposal.Number().String(),
					"round", preprepare.View.Round.String(),
					"hash", preprepare.Proposal.Hash().Hex(),
					"self", c.address.Hex())
				c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash())
				return nil
			}
		}
		return err
	}

	if c.IsProposer() { //is proposer
		//TODO: collect random seed
		//TODO: check is message from validator, same round, same height
		randSeedMessages.Add(msg)
		if randSeedMessages.Size() >= c.QuorumSize() {
			csssStat = PreprepareStep2
			//TODO: Ready To PreprepareStep2
			c.AssambleNewBlockWithRandomData()
		}

		//// If it is locked, it can only process on the locked block.
		//// Passing verifyPrepare and checkMessage implies it is processing on the locked block since it was verified in the Preprepared state.
		//if err := c.verifyPrepare(prepare, src); err != nil {
		//	log.Info("ibftConsensus: handlePrepare verifyPrepare",
		//		"no", prepare.View.Sequence,
		//		"round", prepare.View.Round.String(),
		//		"hash", prepare.Digest.Hex(),
		//		"self", c.address.Hex(),
		//		"err", err)
		//	return err
		//}
		//
		//c.acceptPrepare(msg, src)
		//// Change to Prepared state if we've received enough PREPARE messages or it is locked
		//// and we are in earlier state before Prepared state.
		//if ((c.current.IsHashLocked() && prepare.Digest == c.current.GetLockedHash()) || c.current.GetPrepareOrCommitSize() >= c.QuorumSize()) &&
		//	c.state.Cmp(ibfttypes.StatePrepared) < 0 {
		//	c.current.LockHash()
		//	c.setState(ibfttypes.StatePrepared)
		//	log.Info("ibftConsensus: handlePrepare sendCommit",
		//		"no", prepare.View.Sequence,
		//		"round", prepare.View.Round,
		//		"prepare+commitSize", c.current.GetPrepareOrCommitSize(),
		//		"hash", prepare.Digest.Hex(),
		//		"self", c.address.Hex(),
		//	)
		//	c.sendCommit()
		//}
		//
		//return nil
	} else {
		if !c.valSet.IsProposer(src.Address()) {
			//TODO: start send random seed
			c.BroadcastLocalRandomData()
		}
	}
	return nil
}

func (c *core) AssambleNewBlockWithRandomData() {
	//TODO: Add Random Data to Block Extra
	var blk = c.current.Preprepare.Proposal
	var randSeedData = randSeedMessages.Encode()
	blk.SetExtra(randSeedData)
	log.Info(string(randSeedData))
	//打印随机数
	//分配奖励
	//TODO: Start preprepare round
}

func (c *core) handlePreprepareStep2(msg *ibfttypes.Message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare *istanbul.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return istanbulcommon.ErrFailedDecodePreprepare
	}

	log.Info("ibftConsensus: handlePreprepare",
		"no", preprepare.Proposal.Number().Uint64(),
		"round", preprepare.View.Round.String(),
		"from", src.Address().Hex(),
		"self", c.address.Hex(),
		"hash", preprepare.Proposal.Hash().Hex(),
		"state", c.state)
	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	if err := c.checkMessage(ibfttypes.MsgPreprepare, preprepare.View); err != nil {
		log.Error("ibftConsensus: handlePreprepare checkMessage",
			"no", preprepare.Proposal.Number().Uint64(),
			"round", preprepare.View.Round.String(),
			"err", err.Error(),
			"hash", preprepare.Proposal.Hash().Hex(),
			"self", c.address.Hex())
		if err == istanbulcommon.ErrOldMessage {
			// Get validator set for the given proposal
			valSet := c.backend.ParentValidators(preprepare.Proposal).Copy()
			previousProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
			valSet.CalcProposer(previousProposer, preprepare.View.Round.Uint64())
			// Broadcast COMMIT if it is an existing block
			// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
			// 2. The given block must exist
			if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(preprepare.Proposal.Hash(), preprepare.Proposal.Number()) {
				log.Info("ibftConsensus: handlePreprepare sendCommitForOldBlock",
					"no", preprepare.Proposal.Number().String(),
					"round", preprepare.View.Round.String(),
					"hash", preprepare.Proposal.Hash().Hex(),
					"self", c.address.Hex())
				c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash())
				return nil
			}
		}
		return err
	}

	// Check if the message comes from current proposer
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore preprepare messages from non-proposer", "no", preprepare.Proposal.Number().String(),
			"author", src.Address().Hex(), "round", preprepare.View.Round.String())
		return istanbulcommon.ErrNotFromProposer
	}

	preProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
	log.Info("preProposer:", preProposer.String())
	if duration, err := c.backend.Verify(preprepare.Proposal); err != nil {
		// if it's a future block, we will handle it again after the duration
		if err == consensus.ErrFutureBlock {
			logger.Info("Proposed block will be handled in the future", "err", err, "duration", duration)
			c.stopFuturePreprepareTimer()
			c.futurePreprepareTimer = time.AfterFunc(duration, func() {
				c.sendEvent(backlogEvent{
					src: src,
					msg: msg,
				})
			})
		} else {
			logger.Warn("Failed to verify proposal", "err", err, "duration", duration)
			log.Info("caver|handlePreprepare|sendNextRoundChange1", "no", preprepare.Proposal.Number().String(),
				"round", preprepare.View.Round.String(), "is proposer", strconv.FormatBool(c.IsProposer()))
			c.sendNextRoundChange()
		}
		return err
	}

	// Here is about to accept the PRE-PREPARE
	if c.state == ibfttypes.StateAcceptRequest {
		// Send ROUND CHANGE if the locked proposal and the received proposal are different
		if c.current.IsHashLocked() {
			if preprepare.Proposal.Hash() == c.current.GetLockedHash() {
				log.Info("ibftConsensus: preprepare.Proposal.Hash() == c.current.GetLockedHash()",
					"no", preprepare.Proposal.Number(),
					"round", c.currentView().Round,
					"self", c.address.Hex())
				// Broadcast COMMIT and enters Prepared state directly
				c.acceptPreprepare(preprepare)
				c.setState(ibfttypes.StatePrepared)
				c.sendCommit()
			} else {
				log.Info("ibftConsensus: handlePreprepare sendNextRoundChange2", "no", preprepare.Proposal.Number().String(),
					"round", preprepare.View.Round.String(), "isProposer", strconv.FormatBool(c.IsProposer()))
				// Send round change
				c.sendNextRoundChange()
			}
		} else {
			// Either
			//   1. the locked proposal and the received proposal match
			//   2. we have no locked proposal
			c.acceptPreprepare(preprepare)
			c.setState(ibfttypes.StatePreprepared)
			log.Info("ibftConsensus: handlePreprepare sendPrepare",
				"no", preprepare.View.Sequence,
				"round", preprepare.View.Round,
				"self", c.address.Hex(),
			)
			c.sendPrepare()
		}
	}
	return nil
}

func (c *core) acceptPreprepare(preprepare *istanbul.Preprepare) {
	c.consensusTimestamp = time.Now()
	c.current.SetPreprepare(preprepare)
}
