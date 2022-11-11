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
	"errors"
	"reflect"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"

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
		if c.current.Commits.Size() >= c.QuorumSize() {
			encodedCommitSeals, _ := ibfttypes.Encode(c.current.Commits.Values())
			c.broadcast(&ibfttypes.Message{
				Code:               ibfttypes.MsgCommit,
				Msg:                encodedSubject,
				ProposerCommitSeal: encodedCommitSeals,
				FinaleBlock:        c.finaleBlock,
			})
		}
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
	// Commit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.

	if c.current.Commits.Size() >= c.QuorumSize() && c.IsProposer() {
		if c.commitHeight < commit.View.Sequence.Uint64() {
			if c.state.Cmp(ibfttypes.StateCommitted) < 0 {
				// Still need to call LockHash here since state can skip Prepared state and jump directly to the Committed state.
				log.Info("ibftConsensus: handleCommit proposer commit",
					"no", commit.View.Sequence,
					"round", commit.View.Round,
					"CommitsSize", c.current.Commits.Size(),
					"hash", commit.Digest.Hex(),
					"self", c.address.Hex(),
				)
				c.commitHeight = commit.View.Sequence.Uint64()
				c.commitMsg = *msg
				curBlock := new(types.ProposerBlock)
				curBlock.Sequence = commit.View.Sequence
				curBlock.Round = commit.View.Round
				curBlock.Digest = commit.Digest

				//curBlock.Commit = c.current.Commits.Values()
				curBlock.Commit, _ = ibfttypes.Encode(c.current.Commits.Values())

				// get online validators
				var onlineValidators []common.Address
				for _, v := range c.current.Commits.Values() {
					onlineValidators = append(onlineValidators, v.Address)
				}

				// assert block
				blk := c.backend.GetProposerBlock()

				// assert task state
				state := c.backend.GetProposerState()
				scopy := state.Copy()
				deepBlk := deepCopyBlock(blk)
				deepBlk, _, err = c.restructureBlock(deepBlk, scopy, curBlock)
				if err != nil {
					return err
				}
				deepBlk, err = deepBlk.UpdateBlockSig(c.backend.GetPirvateKey())
				if err != nil {
					return err
				}
				c.finaleBlock, err = ibfttypes.Encode(deepBlk)
				if err != nil {
					return err
				}
				c.backend.GetProposerCh() <- curBlock
				c.sendCommit()
				c.current.LockHash()
				c.commit()
				return nil
			} else {
				log.Error("ibftConsensus: handleCommit proposer commit > StateCommitted err",
					"no", commit.View.Sequence,
					"round", commit.View.Round,
					"CommitsSize", c.current.Commits.Size(),
					"hash", commit.Digest.Hex(),
					"self", c.address.Hex(),
				)
				return istanbulcommon.ErrProposerCommitted
			}
		} else {
			log.Error("ibftConsensus: handleCommit ErrProposerCommitted err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex(), "height", c.commitHeight)
			return istanbulcommon.ErrProposerCommitted
		}
	}

	var commitseals []*ibfttypes.Message
	if c.valSet.IsProposer(src.Address()) {
		err = msg.DecodeCommitlist(&commitseals)
		if err != nil {
			log.Error("ibftConsensus: handleCommit DecodeRewardSeals err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex(), "err", err.Error())
			return istanbulcommon.ErrFailedDecodeCommit
		} else {
			log.Info("ibftConsensus: handleCommit DecodeRewardSeals ok")
		}
	} else {
		log.Error("ibftConsensus: handleCommit Decodecommit  ErrNotFromProposer err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex())
		return istanbulcommon.ErrNotFromProposer
	}

	if len(commitseals) >= c.QuorumSize() && c.state.Cmp(ibfttypes.StateCommitted) < 0 {
		// Still need to call LockHash here since state can skip Prepared state and jump directly to the Committed state.
		curBlock := new(types.ProposerBlock)
		curBlock.Sequence = commit.View.Sequence
		curBlock.Round = commit.View.Round
		curBlock.Digest = commit.Digest
		curBlock.Commit, _ = ibfttypes.Encode(commitseals)
		c.backend.GetProposerCh() <- curBlock
		c.current.LockHash()
		c.finaleBlock = msg.FinaleBlock
		var fBlock *types.Block
		err = msg.DecodeFinalBlock(&fBlock)
		if err != nil {
			return err
		}
		log.Info("ibftConsensus: handleCommit commit",
			"no", commit.View.Sequence,
			"round", commit.View.Round,
			"CommitsSize", c.current.Commits.Size(),
			"hash", commit.Digest.Hex(),
			"self", c.address.Hex(),
			"finalBlock", msg.FinaleBlock,
		)
		c.backend.SetFinalBlock(fBlock)
		c.commit()
	} else {
		log.Error("ibftConsensus: handleCommit len(commitseals) < c.QuorumSize() err", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.Address().Hex())
		return istanbulcommon.ErrSmallThenQuorumSize
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

func (c *core) restructureBlock(block *types.Block, stateDB *state.StateDB, proposerBlock *types.ProposerBlock) (*types.Block, *state.StateDB, error) {
	log.Info("enter restructureBlock", "no", block.Number(), "hash", block.Hash())
	if block == nil || stateDB == nil {
		log.Error("block is nil or statedb is nil ")
		return nil, nil, errors.New("block is nil or statedb is nil")
	}

	var (
		onlineAddrs     []common.Address
		tempProposerBlk *types.ProposerBlock
		msgs            []*ibfttypes.Message
	)

	tempProposerBlk = proposerBlock
	err := rlp.DecodeBytes(proposerBlock.Commit, &msgs)
	if err != nil {
		log.Error("failed rlp decode commit msgs", "err", err.Error())
		return nil, nil, err
	}
	for _, v := range msgs {
		log.Info("online validator", "addr", v.Address.Hex())
		onlineAddrs = append(onlineAddrs, v.Address)
		if len(onlineAddrs) == 7 {
			break
		}
	}

	// get staker
	istanbulExtra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		log.Error("failed extract IstanbulExtra", "err", err.Error())
		return nil, nil, err
	}
	// Online data stored in extra
	payload, err := rlp.EncodeToBytes(tempProposerBlk)
	if err != nil {
		log.Error("failed rlp encode onlineSeal", "err", err.Error())
		return nil, nil, err
	}

	istanbulExtra.OnlineSeal = payload
	extraPayload, err := rlp.EncodeToBytes(istanbulExtra)
	if err != nil {
		log.Error("failed rlp encode istanbulExtra", "err", err)
		return nil, nil, err
	}

	blk := deepCopyBlock(block)

	// set extra
	extra := blk.Header().Extra
	extra = append(extra[:types.IstanbulExtraVanity], extraPayload...)
	blk.SetExtra(extra)

	// set root
	root := stateDB.IntermediateRoot(true)
	blk.SetRoot(root)

	// // rewards
	state := stateDB.Copy()
	// state.CreateNFTByOfficial16(onlineAddrs, istanbulExtra.ExchangerAddr, block.Number())
	return blk, state, nil
}

func deepCopyBlock(block *types.Block) *types.Block {
	if block == nil {
		return nil
	}

	h := types.CopyHeader(block.Header())
	return types.NewBlock(h, block.Transactions(), block.Uncles(), nil, new(trie.Trie))
}
