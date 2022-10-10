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
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) sendCommit() {
	sub := c.current.Subject()
	log.Info("ibftConsensus: sendCommit",
		"no", sub.View.Sequence.Uint64(),
		"round", sub.View.Round.String(),
		"hash", sub.Digest.Hex(),
		"self", c.Address().Hex())
	c.broadcastCommit(sub)
}

func (c *core) sendCommitForOldBlock(view *istanbul.View, digest common.Hash) {
	sub := &istanbul.Subject{
		View:   view,
		Digest: digest,
	}
	c.broadcastCommit(sub)
}

func (c *core) broadcastCommit(sub *istanbul.Subject) {
	logger := c.logger.New("state", c.state)

	encodedSubject, err := ibfttypes.Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}

	if c.IsProposer() {
		commits := c.current.Prepares.Values()
		encodedCommitSeals, errSeals := ibfttypes.Encode(commits)
		if errSeals != nil {
			logger.Error("Failed to encode", "commitseals", commits)
			return
		}
		c.broadcast(&ibfttypes.Message{
			Code:          ibfttypes.MsgCommit,
			Msg:           encodedSubject,
			CommittedSeal: encodedCommitSeals,
		})
	} else {
		c.broadcast(&ibfttypes.Message{
			Code: ibfttypes.MsgCommit,
			Msg:  encodedSubject,
		})
	}
}

func (c *core) handleCommit(msg *ibfttypes.Message, src istanbul.Validator) error {
	// Decode COMMIT message
	var commit *istanbul.Subject
	err := msg.Decode(&commit)

	if err != nil {
		log.Error("ibftConsensus: handleCommit Decodecommit  err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex())
		return istanbulcommon.ErrFailedDecodeCommit
	}
	var commitseals []*ibfttypes.Message
	if c.valSet.IsProposer(src.Address()) {
		err = msg.DecodeCommitSeals(&commitseals)
		if err != nil {
			log.Error("ibftConsensus: handleCommit DecodecommitSeals  err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex())
			return istanbulcommon.ErrFailedDecodeCommit
		}
	}

	log.Info("ibftConsensus: handleCommit info", "no", commit.View.Sequence,
		"round", commit.View.Round,
		"from", src.Address().Hex(),
		"hash", commit.Digest.Hex(),
		"self", c.Address().Hex())

	if err := c.checkMessage(ibfttypes.MsgCommit, commit.View); err != nil {
		log.Error("ibftConsensus: handleCommit checkMessage", "no", commit.View.Sequence,
			"round", commit.View.Round,
			"who", c.address.Hex(),
			"hash", commit.Digest.Hex(),
			"self", c.address.Hex(),
			"err", err.Error())
		return err
	}

	if err := c.verifyCommit(commit, src); err != nil {
		log.Error("ibftConsensus: handleCommit verifyCommit", "no", commit.View.Sequence, "round", commit.View.Round, "self", c.address.Hex(), "hash", commit.Digest.Hex(), "err", err.Error())
		return err
	}

	c.acceptCommit(msg, src)
	log.Info("ibftConsensus: handleCommit baseinfo", "no", commit.View.Sequence.Uint64(), "round", commit.View.Round, "from", src.Address().Hex(), "hash", commit.Digest.Hex(), "self", c.address.Hex())

	proposerCommited := false
	for _, v := range c.current.Commits.Values() {
		if c.valSet.IsProposer(v.Address) {
			proposerCommited = true
			break
		}
	}

	// Commit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	if c.current.Commits.Size() >= c.QuorumSize() && c.state.Cmp(ibfttypes.StateCommitted) < 0 && proposerCommited {
		// Still need to call LockHash here since state can skip Prepared state and jump directly to the Committed state.
		log.Info("ibftConsensus: handleCommit commit",
			"no", commit.View.Sequence,
			"round", commit.View.Round,
			"CommitsSize", c.current.Commits.Size(),
			"hash", commit.Digest.Hex(),
			"self", c.address.Hex(),
		)
		c.current.LockHash()
		c.commit()
	}

	return nil
}

// verifyCommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyCommit(commit *istanbul.Subject, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return istanbulcommon.ErrInconsistentSubject
	}

	return nil
}

func (c *core) acceptCommit(msg *ibfttypes.Message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the COMMIT message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	return nil
}

func GenTestExtra() ([]byte, error) {
	addr0 := common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963")
	addr1 := common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F")
	addr2 := common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD")
	addr3 := common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349")
	addr4 := common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326")
	addr5 := common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF")
	addr6 := common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6")
	valSet := validator.NewSet([]common.Address{addr0, addr1, addr2, addr3, addr4, addr5, addr6}, istanbul.NewRoundRobinProposerPolicy())

	ms := newMessageSet(valSet)

	view := &istanbul.View{
		Round:    new(big.Int),
		Sequence: new(big.Int),
	}
	pp := &istanbul.Preprepare{
		View:     view,
		Proposal: nil,
	}

	rawPP, err := rlp.EncodeToBytes(pp)
	if err != nil {
		return nil, nil
	}
	msg := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: valSet.GetProposer().Address(),
	}
	ms.Add(msg)
	msg0 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr0,
	}
	msg1 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr1,
	}
	msg2 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr2,
	}
	msg3 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr3,
	}
	msg4 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr4,
	}
	msg5 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr5,
	}
	msg6 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr6,
	}
	ms.Add(msg0)
	ms.Add(msg1)
	ms.Add(msg2)
	ms.Add(msg3)
	ms.Add(msg4)
	ms.Add(msg5)
	ms.Add(msg6)

	encodedCommitSeals, _ := ibfttypes.Encode(ms.Values())
	message := &ibfttypes.Message{
		Code:          ibfttypes.MsgCommit,
		Signature:     []byte{},
		CommittedSeal: encodedCommitSeals, // small hack
	}
	return rlp.EncodeToBytes(message)
}
