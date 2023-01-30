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
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
)

func (c *core) sendPreprepare(request *istanbul.Request) {
	logger := c.logger.New("state", c.state)
	// If I'm the proposer and I have the same sequence with the proposal

	log.Info("ibftConsensus: sendPreprepare",
		"no", request.Proposal.Number(), "no", c.current.sequence,
		"round", c.current.round.Uint64(), "isproposer", c.IsProposer(), "slef", c.address.Hex())
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 {
		curView := c.currentView()

		var proposal istanbul.Proposal
		if c.IsProposer() {
			proposal = request.Proposal
		} else {
			proposal = nil
		}
		preprepare, err := ibfttypes.Encode(&istanbul.Preprepare{
			View:     curView,
			Proposal: proposal,
		})

		consensusData := ConsensusData{
			Height: curView.Sequence.String(),
			Rounds: map[int64]RoundInfo{
				curView.Round.Int64(): {
					Method:     "sendPreprepare",
					Timestamp:  time.Now().UnixNano(),
					Sender:     c.address,
					Sequence:   curView.Sequence.Uint64(),
					Round:      curView.Round.Int64(),
					Hash:       request.Proposal.Hash(),
					Miner:      c.valSet.GetProposer().Address(),
					Error:      err,
					IsProposal: c.IsProposer(),
				},
			},
		}
		c.SaveData(consensusData)

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

func (c *core) simpleCheckPreprepare(msg *ibfttypes.Message, _ istanbul.Validator) error {
	var preprepare *istanbul.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return err
	}
	if preprepare.Proposal.Number().Uint64() < c.backend.CurrentNumber() {
		return istanbulcommon.ErrOldMessage
	}
	c.PutAddr(preprepare.View.Sequence.Uint64(), msg.Address)
	return nil
}

func (c *core) simpleCheckSubject(msg *ibfttypes.Message, _ istanbul.Validator) error {
	var rc *istanbul.Subject
	if err := msg.Decode(&rc); err != nil {
		return err
	}
	if rc.View.Sequence.Uint64() < c.backend.CurrentNumber() {
		return istanbulcommon.ErrOldMessage
	}
	c.PutAddr(rc.View.Sequence.Uint64(), msg.Address)
	return nil
}

func (c *core) handlePreprepare(msg *ibfttypes.Message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Decode PRE-PREPARE
	var preprepare *istanbul.Preprepare
	err := msg.Decode(&preprepare)
	if err != nil {
		return istanbulcommon.ErrFailedDecodePreprepare
	}

	// Save the online validator of the current sequence
	c.PutAddr(preprepare.View.Sequence.Uint64(), msg.Address)

	// Not the prepare message sent by the proposer, return nil directly
	if preprepare.Proposal == nil {
		return nil
	}

	roundInfo := RoundInfo{
		Method:     "handlePreprepare",
		Timestamp:  time.Now().UnixNano(),
		Sender:     src.Address(),
		Receiver:   c.address,
		Sequence:   preprepare.View.Sequence.Uint64(),
		Round:      preprepare.View.Round.Int64(),
		Hash:       preprepare.Proposal.Hash(),
		Miner:      c.valSet.GetProposer().Address(),
		Error:      err,
		IsProposal: c.IsProposer(),
	}

	consensusData := ConsensusData{
		Height: preprepare.View.Sequence.String(),
		Rounds: map[int64]RoundInfo{
			preprepare.View.Round.Int64(): roundInfo,
		},
	}
	c.SaveData(consensusData)

	log.Info("ibftConsensus: handlePreprepare",
		"no", preprepare.Proposal.Number().Uint64(),
		"round", preprepare.View.Round.String(),
		"from", src.Address().Hex(),
		"self", c.address.Hex(),
		"hash", preprepare.Proposal.Hash().Hex(),
		"state", c.state)

	// Ensure we have the same view with the PRE-PREPARE message
	// If it is old message, see if we need to broadcast COMMIT
	err = c.checkMessage(ibfttypes.MsgPreprepare, preprepare.View)
	roundInfo.Method = "handlePreprepare checkMessage"
	roundInfo.Timestamp = time.Now().UnixNano()
	roundInfo.Error = err
	consensusData.Rounds = map[int64]RoundInfo{
		preprepare.View.Round.Int64(): roundInfo,
	}
	c.SaveData(consensusData)
	if err != nil {
		log.Error("ibftConsensus: handlePreprepare checkMessage",
			"no", preprepare.Proposal.Number().Uint64(),
			"round", preprepare.View.Round.String(),
			"err", err.Error(),
			"hash", preprepare.Proposal.Hash().Hex(),
			"self", c.address.Hex())

		// remove this part on 2022-1002
		//if err == istanbulcommon.ErrOldMessage {
		//	// Get validator set for the given proposal
		//	if block, ok := preprepare.Proposal.(*types.Block); ok {
		//		if block.Header() == nil {
		//			log.Error("ibftConsensus: header is nil")
		//			return errors.New("ibftConsensus: header is nil")
		//		}
		//	} else {
		//		log.Error("ibftConsensus: block not ok")
		//		return errors.New("ibftConsensus: block not ok")
		//	}
		//	valSet := c.backend.ParentValidators(preprepare.Proposal).Copy()
		//	previousProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
		//	valSet.CalcProposer(previousProposer, preprepare.View.Round.Uint64())
		//	// Broadcast COMMIT if it is an existing block
		//	// 1. The proposer needs to be a proposer matches the given (Sequence + Round)
		//	// 2. The given block must exist
		//	if valSet.IsProposer(src.Address()) && c.backend.HasPropsal(preprepare.Proposal.Hash(), preprepare.Proposal.Number()) {
		//		log.Info("ibftConsensus: handlePreprepare sendCommitForOldBlock",
		//			"no", preprepare.Proposal.Number().String(),
		//			"round", preprepare.View.Round.String(),
		//			"hash", preprepare.Proposal.Hash().Hex(),
		//			"self", c.address.Hex())
		//		c.sendCommitForOldBlock(preprepare.View, preprepare.Proposal.Hash())
		//		return nil
		//	}
		//}
		return err
	}

	// Check if the message comes from current proposer
	if !c.valSet.IsProposer(src.Address()) {
		logger.Warn("Ignore preprepare messages from non-proposer", "no", preprepare.Proposal.Number().String(),
			"author", src.Address().Hex(), "round", preprepare.View.Round.String())
		return istanbulcommon.ErrNotFromProposer
	}
	// preProposer := c.backend.GetProposer(preprepare.Proposal.Number().Uint64() - 1)
	// if preProposer.String() == "0x0000000000000000000000000000000000000000" && preprepare.Proposal.Number().Uint64() > 1 {
	// 	log.Error("preProposer is empty block:", "no", preProposer.String())
	// 	return errors.New("preProposer is empty block")

	// } else {

	// }
	duration, err := c.backend.Verify(preprepare.Proposal)
	roundInfo.Method = "handlePreprepare verify"
	roundInfo.Timestamp = time.Now().UnixNano()
	roundInfo.Error = err
	consensusData.Rounds = map[int64]RoundInfo{
		preprepare.View.Round.Int64(): roundInfo,
	}
	c.SaveData(consensusData)
	if err != nil {
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
