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
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemoryPeers      = 40
	inmemoryMessages   = 1024
)

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (sb *Backend) Author(header *types.Header) (common.Address, error) {
	return sb.EngineForBlockNumber(header.Number).Author(header)
}

// Signers extracts all the addresses who have signed the given header
// It will extract for each seal who signed it, regardless of if the seal is
// repeated
func (sb *Backend) Signers(header *types.Header) ([]common.Address, error) {
	return sb.EngineForBlockNumber(header.Number).Signers(header)
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *Backend) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil)
}

func (sb *Backend) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header) error {
	if header.Number.Cmp(big.NewInt(0)) == 0 {
		genesis := chain.GetHeaderByNumber(0)
		if err := sb.EngineForBlockNumber(big.NewInt(0)).VerifyHeader(chain, genesis, nil, nil); err != nil {
			sb.logger.Error("verifyHeader : invalid genesis block", "err", err, "hash", genesis.Hash())
			return err
		}
	} else {
		// Get the validatorset for this round
		istanbulExtra, err := types.ExtractIstanbulExtra(header)
		if err != nil {
			return istanbulcommon.ErrInvalidExtraDataFormat
		}
		validators := istanbulExtra.Validators
		valSet := validator.NewSet(validators, sb.config.ProposerPolicy)
		return sb.EngineForBlockNumber(header.Number).VerifyHeader(chain, header, parents, valSet)
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *Backend) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		errored := false
		for i, header := range headers {
			if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
				all, readeorrer := sb.chain.(*core.BlockChain).ReadValidatorPool(sb.chain.CurrentHeader())
				if readeorrer != nil {
					log.Info("VerifyHeadersAndcreaderr", "data:", readeorrer)
					continue
				}

				istanbulExtra, checkerr := types.ExtractIstanbulExtra(header)
				if checkerr != nil {
					log.Info("VerifyHeadersAndcheckerr", "data:", checkerr)
					continue
				}
				validators := istanbulExtra.Validators
				log.Info("VerifyHeadersAndcheckerr", "validators", validators)
				var total = big.NewInt(0)
				for _, v := range validators {
					balance := all.StakeBalance(v)
					total.Add(total, balance)
				}
				if total.Cmp(all.TargetSize()) < 0 {
					continue
				}

				var err error
				err = nil
				select {
				case <-abort:
					return
				case results <- err:
				}
				continue
			}
			var err error
			if errored {
				err = consensus.ErrUnknownAncestor
			} else {
				err = sb.verifyHeader(chain, header, headers[:i])
				if err != nil {
					log.Error("VerifyHeaders err", "err", err.Error())
				}
			}

			if err != nil {
				errored = true
			}

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (sb *Backend) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return sb.EngineForBlockNumber(block.Header().Number).VerifyUncles(chain, block)
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *Backend) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return istanbulcommon.ErrUnknownBlock
	}

	var valSet istanbul.ValidatorSet
	if c, ok := chain.(*core.BlockChain); ok {
		validatorList, err := c.Random11ValidatorFromPool(c.CurrentBlock().Header())
		if err != nil {
			log.Error("VerifySeal : invalid validator list", "no", c.CurrentBlock().Header(), "err", err)
			return err
		}
		for _, v := range validatorList.Validators {
			log.Info("Backend|VerifySeal", "height", c.CurrentBlock().Header().Number.Uint64(), "v", v)
		}

		valSet = validator.NewSet(validatorList.ConvertToAddress(), sb.config.ProposerPolicy)
	}

	return sb.EngineForBlockNumber(header.Number).VerifySeal(chain, header, valSet)
}

// PrepareForEmptyBlock initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) PrepareForEmptyBlock(chain consensus.ChainHeaderReader, header *types.Header) error {
	var valSet istanbul.ValidatorSet
	if c, ok := chain.(*core.BlockChain); ok {
		log.Info("Prepare", "header-no", header.Number.String(), "current-header", c.CurrentBlock().Header().Number.String())
		cHeader := c.CurrentBlock().Header()
		if cHeader == nil {
			return errors.New("prepare err: current header is nil")
		}
		validatorList, err := c.ReadValidatorPool(cHeader)
		if err != nil {
			log.Error("PrepareForEmptyBlock : err", "err", err)
			return err
		}
		valSet = validator.NewSet(validatorList.ConvertToAddress(), sb.config.ProposerPolicy)
	}

	err := sb.EngineForBlockNumber(header.Number).PrepareEmpty(chain, header, valSet)
	if err != nil {
		return err
	}
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *Backend) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	var valSet istanbul.ValidatorSet
	if c, ok := chain.(*core.BlockChain); ok {
		cBlk := c.CurrentBlock()
		if cBlk == nil {
			return errors.New("err prepare : current block is nil")
		}

		log.Info("Prepare : info", "header-no", header.Number.String(), "current-header", c.CurrentBlock().Header().Number)

		validatorList, err := c.Random11ValidatorFromPool(cBlk.Header())
		if err != nil {
			log.Error("Prepare: invalid validator list", "err", err, "no", cBlk.Header().Number)
			return err
		}

		for _, v := range validatorList.Validators {
			log.Info("Backend : Prepare", "height", cBlk.Number, "v", v)
		}
		valSet = validator.NewSet(validatorList.ConvertToAddress(), sb.config.ProposerPolicy)
	}

	err := sb.EngineForBlockNumber(header.Number).Prepare(chain, header, valSet)
	if err != nil {
		return err
	}

	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *Backend) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	sb.EngineForBlockNumber(header.Number).Finalize(chain, header, state, txs, uncles)
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (sb *Backend) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return sb.EngineForBlockNumber(header.Number).FinalizeAndAssemble(chain, header, state, txs, uncles, receipts)

}

// SealforEmptyBlock generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *Backend) SealforEmptyBlock(chain consensus.ChainHeaderReader, block *types.Block, validators []common.Address) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	var emptyBlock *types.Block
	header := block.Header()

	if sb.core == nil {
		return emptyBlock, errors.New("seal|ibft engine not active")
	}

	log.Info("caver|SealforEmptyBlock|enter", "sealNo", block.Number().String(), "is proposer", sb.core.IsProposer())

	//Get the validatorset for this round
	//istanbulExtra, err1 := types.ExtractIstanbulExtra(header)
	//if err1 != nil {
	//	log.Info("caver|seal|ExtractIstanbulExtra|Empty", "no", header.Number, "err", err1.Error())
	//	return emptyBlock, err1
	//}
	//valSet := validator.NewSet(istanbulExtra.Validators, sb.config.ProposerPolicy)
	valSet := validator.NewSet(validators, sb.config.ProposerPolicy)

	emptyBlock, err := sb.EngineForBlockNumber(header.Number).Seal(chain, block, valSet)
	if err != nil {
		log.Info("caver|SealforEmptyBlock|err", "sealNo", header.Number.String(), "err", err.Error())
		return emptyBlock, err
	}
	return emptyBlock, err
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *Backend) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Uint64() > 0 {
		log.Error("Seal : coinbase error", "err", "coinbase is 0")
		return errors.New("coinbase is 0")
	}

	if sb.core == nil {
		return errors.New("seal : ibft engine not active")
	}

	log.Info("seal : enter", "no", block.Number().String(), "is proposer", sb.core.IsProposer())

	//Get the validatorset for this round
	istanbulExtra, err1 := types.ExtractIstanbulExtra(header)
	if err1 != nil {
		log.Error("Seal : ExtractIstanbulExtra", "err", err1)
		return err1
	}

	valSet := validator.NewSet(istanbulExtra.Validators, sb.config.ProposerPolicy)

	block, err := sb.EngineForBlockNumber(header.Number).Seal(chain, block, valSet)
	if err != nil {
		return err
	}

	delay := time.Until(time.Unix(int64(block.Header().Time), 0))

	go func() {
		// wait for the timestamp of header, use this to adjust the block period
		select {
		case <-time.After(delay):
		case <-stop:
			results <- nil
			return
		}

		// get the proposed block hash and clear it if the seal() is completed.
		sb.sealMu.Lock()
		sb.proposedBlockHash = block.Hash()

		defer func() {
			sb.proposedBlockHash = common.Hash{}
			sb.sealMu.Unlock()
		}()

		log.Info("seal : post block into Istanbul engine", "no", block.NumberU64(),
			"hash", block.Hash())
		// post block into Istanbul engine
		go sb.EventMux().Post(istanbul.RequestEvent{
			Proposal: block,
		})
		proposerCommitData := new(types.ProposerBlock)
		for {

			select {
			case proposerCommit := <-sb.proposerCh:
				proposerCommitData = proposerCommit
			case enqueueBlock := <-sb.enqueueCh:
				if enqueueBlock != nil {
					enqueueBlock.ReceivedFrom = proposerCommitData
					log.Info("sb.enqueueCh", "round", proposerCommitData.Round, "sequence", proposerCommitData.Sequence)
					curBlock := new(types.Block)
					curBlock.ReceivedFrom = enqueueBlock
					results <- curBlock
					return
				}
			case result := <-sb.commitCh:
				// if the block hash and the hash from channel are the same,
				// return the result. Otherwise, keep waiting the next hash.
				if result != nil && block.Hash() == result.Hash() {
					//ProposerBlock := new(types.ProposerBlock)
					result.ReceivedFrom = proposerCommitData
					results <- result
					return
				}
			case <-stop:
				results <- nil
				return
			}
		}
	}()
	return nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *Backend) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "istanbul",
		Version:   "1.0",
		Service:   &API{chain: chain, backend: sb},
		Public:    true,
	}}
}

// Start implements consensus.Istanbul.Start
func (sb *Backend) Start(chain consensus.ChainHeaderReader, currentBlock func() *types.Block, hasBadBlock func(db ethdb.Reader, hash common.Hash) bool) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return istanbul.ErrStartedEngine
	}

	// clear previous data
	sb.proposedBlockHash = common.Hash{}
	if sb.commitCh != nil {
		close(sb.commitCh)
	}
	sb.commitCh = make(chan *types.Block, 1)

	if sb.notifyBlockCh != nil {
		close(sb.notifyBlockCh)
	}
	sb.notifyBlockCh = make(chan *types.OnlineValidatorList, 1)
	sb.chain = chain
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock

	// Check if qbft Consensus needs to be used after chain is set
	var err error
	if sb.IsQBFTConsensus() {
		err = sb.startQBFT()
	} else {
		err = sb.startIBFT()
	}

	if err != nil {
		return err
	}

	sb.coreStarted = true

	return nil
}

// Stop implements consensus.Istanbul.Stop
func (sb *Backend) Stop() error {
	log.Info("caver|Stop")
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		sb.logger.Info("caver|Stop|ErrStoppedEngine", "!sb.coreStarted", !sb.coreStarted)
		return istanbul.ErrStoppedEngine
	}
	if err := sb.stop(); err != nil {
		return err
	}
	sb.coreStarted = false

	return nil
}

func addrsToString(addrs []common.Address) []string {
	strs := make([]string, len(addrs))
	for i, addr := range addrs {
		strs[i] = addr.String()
	}
	return strs
}

func (sb *Backend) snapLogger(snap *Snapshot) log.Logger {
	return sb.logger.New(
		"snap.number", snap.Number,
		"snap.hash", snap.Hash.String(),
		"snap.epoch", snap.Epoch,
		"snap.validators", addrsToString(snap.validators()),
		"snap.votes", snap.Votes,
	)
}

func (sb *Backend) storeSnap(snap *Snapshot) error {
	logger := sb.snapLogger(snap)
	logger.Debug("BFT: store snapshot to database")
	if err := snap.store(sb.db); err != nil {
		logger.Error("BFT: failed to store snapshot to database", "err", err)
		return err
	}

	return nil
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (sb *Backend) snapshot(chain consensus.ChainHeaderReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)
	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := sb.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			sb.snapLogger(snap).Trace("BFT: loaded voting snapshot from cache")
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(sb.config.Epoch, sb.db, hash); err == nil {
				snap = s
				sb.snapLogger(snap).Trace("BFT: loaded voting snapshot from database")
				break
			}
		}

		// If we're at block zero, make a snapshot
		if number == 0 {
			genesis := chain.GetHeaderByNumber(0)
			if err := sb.EngineForBlockNumber(big.NewInt(0)).VerifyHeader(chain, genesis, nil, nil); err != nil {
				sb.logger.Error("BFT: invalid genesis block", "err", err)
				return nil, err
			}

			// Get the validators from genesis to create a snapshot
			validators, err := sb.EngineForBlockNumber(big.NewInt(0)).Validators(genesis)
			if err != nil {
				sb.logger.Error("BFT: invalid genesis block", "err", err)
				return nil, err
			}

			snap = newSnapshot(sb.config.Epoch, 0, genesis.Hash(), validator.NewSet(validators, sb.config.ProposerPolicy))
			if err := sb.storeSnap(snap); err != nil {
				return nil, err
			}
			break
		}

		// No snapshot for this header, gather the header and move backward
		var header *types.Header
		if len(parents) > 0 {
			// If we have explicit parents, pick from there (enforced)
			header = parents[len(parents)-1]
			if header.Hash() != hash || header.Number.Uint64() != number {
				return nil, consensus.ErrUnknownAncestor
			}
			parents = parents[:len(parents)-1]
		} else {
			// No explicit parents (or no more left), reach out to the database
			header = chain.GetHeader(hash, number)
			if header == nil {
				return nil, consensus.ErrUnknownAncestor
			}
		}

		headers = append(headers, header)
		number, hash = number-1, header.ParentHash
	}

	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}

	snap, err := sb.snapApply(snap, headers)
	if err != nil {
		return nil, err
	}
	sb.recents.Add(snap.Hash, snap)

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if err = sb.storeSnap(snap); err != nil {
			return nil, err
		}
	}

	return snap, err
}

// SealHash returns the hash of a block prior to it being sealed.
func (sb *Backend) SealHash(header *types.Header) common.Hash {
	return sb.EngineForBlockNumber(header.Number).SealHash(header)
}

func (sb *Backend) snapApply(snap *Snapshot, headers []*types.Header) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return snap, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, istanbulcommon.ErrInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != snap.Number+1 {
		return nil, istanbulcommon.ErrInvalidVotingChain
	}
	// Iterate through the headers and create a new snapshot
	snapCpy := snap.copy()

	for _, header := range headers {
		err := sb.snapApplyHeader(snapCpy, header)
		if err != nil {
			return nil, err
		}
	}
	snapCpy.Number += uint64(len(headers))
	snapCpy.Hash = headers[len(headers)-1].Hash()

	return snapCpy, nil
}

func (sb *Backend) snapApplyHeader(snap *Snapshot, header *types.Header) error {
	logger := sb.snapLogger(snap).New("header.number", header.Number.Uint64(), "header.hash", header.Hash().String())

	logger.Trace("BFT: apply header to voting snapshot")

	// Remove any votes on checkpoint blocks
	number := header.Number.Uint64()
	if number%snap.Epoch == 0 {
		snap.Votes = nil
		snap.Tally = make(map[common.Address]Tally)
	}

	// Resolve the authorization key and check against validators
	validator, err := sb.EngineForBlockNumber(header.Number).Author(header)
	if err != nil {
		logger.Error("BFT: invalid header author", "err", err)
		return err
	}

	logger = logger.New("header.author", validator)

	if _, v := snap.ValSet.GetByAddress(validator); v == nil {
		logger.Error("BFT: header author is not a validator")
		return istanbulcommon.ErrUnauthorized
	}

	// Read vote from header
	candidate, authorize, err := sb.EngineForBlockNumber(header.Number).ReadVote(header)
	if err != nil {
		logger.Error("BFT: invalid header vote", "err", err)
		return err
	}

	logger = logger.New("candidate", candidate.String(), "authorize", authorize)
	// Header authorized, discard any previous votes from the validator
	for i, vote := range snap.Votes {
		if vote.Validator == validator && vote.Address == candidate {
			logger.Trace("BFT: discard previous vote from tally", "old.authorize", vote.Authorize)
			// Uncast the vote from the cached tally
			snap.uncast(vote.Address, vote.Authorize)

			// Uncast the vote from the chronological list
			snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
			break // only one vote allowed
		}
	}

	logger.Debug("BFT: add vote to tally")
	if snap.cast(candidate, authorize) {
		snap.Votes = append(snap.Votes, &Vote{
			Validator: validator,
			Block:     number,
			Address:   candidate,
			Authorize: authorize,
		})
	}

	// If the vote passed, update the list of validators
	if tally := snap.Tally[candidate]; tally.Votes > snap.ValSet.Size()/2 {

		if tally.Authorize {
			logger.Info("BFT: reached majority to add validator")
			snap.ValSet.AddValidator(candidate)
		} else {
			logger.Info("BFT: reached majority to remove validator")
			snap.ValSet.RemoveValidator(candidate)

			// Discard any previous votes the deauthorized validator cast
			for i := 0; i < len(snap.Votes); i++ {
				if snap.Votes[i].Validator == candidate {
					// Uncast the vote from the cached tally
					snap.uncast(snap.Votes[i].Address, snap.Votes[i].Authorize)

					// Uncast the vote from the chronological list
					snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)

					i--
				}
			}
		}
		// Discard any previous votes around the just changed account
		for i := 0; i < len(snap.Votes); i++ {
			if snap.Votes[i].Address == candidate {
				snap.Votes = append(snap.Votes[:i], snap.Votes[i+1:]...)
				i--
			}
		}
		delete(snap.Tally, candidate)
	}
	return nil
}

func (sb *Backend) SealOnlineProofBlk(chain consensus.ChainHeaderReader, block *types.Block, notifyBlockCh chan *types.OnlineValidatorList, stop <-chan struct{}) error {
	timeout := time.NewTimer(60 * time.Second)
	defer timeout.Stop()
	log.Info("SealOnlineProofBlk : info", "height", block.Number())
	// Only this round of validators can send online proofs
	header := block.Header()

	if sb.core == nil {
		log.Error("SealOnlineProofBlk : sb.core", "err", errors.New("SealOnlineProofBlk : ibft engine not active !"), "no", header.Number.Uint64())
		return errors.New("SealOnlineProofBlk : ibft engine not active !")
	}

	//Get the validatorset for this round
	// istanbulExtra, err := types.ExtractIstanbulExtra(header)
	// if err != nil {
	// 	log.Error("SealOnlineProofBlk  : istanbulExtra", "err", err.Error())
	// 	return err
	// }

	// valSet := validator.NewSet(istanbulExtra.Validators, sb.config.ProposerPolicy)
	var valset istanbul.ValidatorSet
	if c, ok := chain.(*core.BlockChain); ok {
		log.Info("SealOnlineProofBlk : calculate valset")

		cBlk := c.CurrentBlock()
		validatorList, err := c.Random11ValidatorFromPool(cBlk.Header())
		if err != nil {
			log.Error("SealOnlineProofBlk : err", "err", err.Error(), "no", header.Number.Uint64())
			return err
		}
		log.Info("SealOnlineProofBlk : len", "len", len(validatorList.Validators), "no", header.Number.Uint64())

		valset = validator.NewSet(validatorList.ConvertToAddress(), sb.config.ProposerPolicy)
	}
	if _, v := valset.GetByAddress(sb.address); v == nil {
		log.Error("SealOnlineProofBlk  : ErrUnauthorized", "err", istanbulcommon.ErrUnauthorized, "no", header.Number.Uint64())
		return istanbulcommon.ErrUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		log.Error("SealOnlineProofBlk  : ErrUnknownAncestor", "err", consensus.ErrUnknownAncestor, "no", header.Number.Uint64())
		return consensus.ErrUnknownAncestor
	}

	// generate a local random number
	localTime := time.Now().Nanosecond()
	common.BigToHash(big.NewInt(int64(localTime)))
	log.Info("SealOnlineProofBlk : post OnlineProofEvent", "no", header.Number.Uint64())

	go func() {
		sb.EventMux().Post(istanbul.OnlineProofEvent{
			Proposal:   block,
			RandomHash: common.BigToHash(big.NewInt(int64(localTime))),
		})
		for {
			select {
			case onlineValidators := <-sb.notifyBlockCh:
				log.Info("SealOnlineProofBlk : onlineValidators", "no", block.NumberU64(), "info", onlineValidators)
				if onlineValidators != nil {
					notifyBlockCh <- onlineValidators
					return
				}
			case <-stop:
				log.Info("SealOnlineProofBlk : stop")
				return
			case <-timeout.C:
				log.Warn("SealOnlineProofBlk: Collect online proof timed out", "no", block.NumberU64())
				return
			}
		}
	}()
	return nil
}

func (sb *Backend) GossipOnlineProof(chain consensus.ChainHeaderReader, block *types.Block) error {
	log.Info("GossipOnlineProof : info", "height", block.Number())
	// Only this round of validators can send online proofs
	header := block.Header()

	if sb.core == nil {
		log.Error("GossipOnlineProof : sb.core", "err", errors.New("GossipOnlineProof : ibft engine not active !"), "no", header.Number.Uint64())
		return errors.New("GossipOnlineProof : ibft engine not active !")
	}

	// var valset istanbul.ValidatorSet
	// if c, ok := chain.(*core.BlockChain); ok {
	// 	log.Info("GossipOnlineProof : calculate valset")

	// 	cBlk := c.CurrentBlock()
	// 	validatorList, err := c.Random11ValidatorFromPool(cBlk.Header())
	// 	log.Info("GossipOnlineProof : len", "len", len(validatorList.Validators), "no", header.Number.Uint64())
	// 	if err != nil {
	// 		log.Error("GossipOnlineProof : err", "err", err.Error(), "no", header.Number.Uint64())
	// 		return err
	// 	}
	// 	valset = validator.NewSet(validatorList.ConvertToAddress(), sb.config.ProposerPolicy)
	// }
	// if _, v := valset.GetByAddress(sb.address); v == nil {
	// 	log.Error("GossipOnlineProof  : ErrUnauthorized", "err", istanbulcommon.ErrUnauthorized, "no", header.Number.Uint64())
	// 	return istanbulcommon.ErrUnauthorized
	// }

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		log.Error("GossipOnlineProof  : ErrUnknownAncestor", "err", consensus.ErrUnknownAncestor, "no", header.Number.Uint64())
		return consensus.ErrUnknownAncestor
	}

	// generate a local random number
	localTime := time.Now().Nanosecond()
	common.BigToHash(big.NewInt(int64(localTime)))
	log.Info("GossipOnlineProof : post OnlineProofEvent", "no", header.Number.Uint64())

	go sb.EventMux().Post(istanbul.OnlineProofEvent{
		Proposal:   block,
		RandomHash: common.BigToHash(big.NewInt(int64(localTime))),
		Version:    params.Version,
	})

	return nil
}

func (sb *Backend) FinalizeOnlineProofBlk(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	return sb.EngineForBlockNumber(header.Number).FinalizeOnlineProofBlk(chain, header, state, txs, uncles, receipts)
}
