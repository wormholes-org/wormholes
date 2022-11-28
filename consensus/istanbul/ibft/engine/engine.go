package ibftengine

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"golang.org/x/crypto/sha3"
)

var (
	nilUncleHash  = types.CalcUncleHash(nil)                 // Always Keccak256(RLP([])) as uncles are meaningless outside of PoW.
	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.
)

type SignerFn func(data []byte) ([]byte, error)

type Engine struct {
	cfg *istanbul.Config

	signer  common.Address // Ethereum address of the signing key
	sign    SignerFn       // Signer function to authorize hashes with
	backend istanbul.Backend
}

func NewEngine(cfg *istanbul.Config, signer common.Address, sign SignerFn, backend istanbul.Backend) *Engine {
	return &Engine{
		cfg:     cfg,
		signer:  signer,
		sign:    sign,
		backend: backend,
	}
}

func (e *Engine) Author(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, err
	}

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		return common.HexToAddress("0x0000000000000000000000000000000000000000"), nil
	} else {
		addr, err := istanbul.GetSignatureAddress(sigHash(header).Bytes(), extra.Seal)
		if err != nil {
			return addr, err
		}
		return addr, nil
	}

}

func (e *Engine) CommitHeader(header *types.Header, seals [][]byte, round *big.Int) error {
	// Append seals into extra-data
	return writeCommittedSeals(header, seals)
}

func (e *Engine) VerifyBlockProposal(chain consensus.ChainHeaderReader, block *types.Block, validators istanbul.ValidatorSet) (time.Duration, error) {
	// check block body
	txnHash := types.DeriveSha(block.Transactions(), new(trie.Trie))
	if txnHash != block.Header().TxHash {
		return 0, istanbulcommon.ErrMismatchTxhashes
	}

	// uncleHash := types.CalcUncleHash(block.Uncles())
	// if uncleHash != nilUncleHash {
	// 	return 0, istanbulcommon.ErrInvalidUncleHash
	// }

	if block.Coinbase() == common.HexToAddress("0x0000000000000000000000000000000000000000") && block.Number().Cmp(common.Big0) > 0 {
		return 0, istanbulcommon.ErrEmptyBlock
	} else {
		// verify the header of proposed block
		err := e.VerifyHeader(chain, block.Header(), nil, validators)
		if err == nil || err == istanbulcommon.ErrEmptyCommittedSeals {
			// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
			return 0, nil
		} else if err == consensus.ErrFutureBlock {
			return time.Until(time.Unix(int64(block.Header().Time), 0)), consensus.ErrFutureBlock
		}
		return 0, err
	}

}

func (e *Engine) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	return e.verifyHeader(chain, header, parents, validators)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (e *Engine) verifyHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	if header.Number == nil {
		return istanbulcommon.ErrUnknownBlock
	}

	// Don't waste time checking blocks from the future (adjusting for allowed threshold)
	adjustedTimeNow := time.Now().Add(time.Duration(e.cfg.AllowedFutureBlockTime) * time.Second).Unix()
	if header.Time > uint64(adjustedTimeNow) {
		return consensus.ErrFutureBlock
	}

	if _, err := types.ExtractIstanbulExtra(header); err != nil {
		return istanbulcommon.ErrInvalidExtraDataFormat
	}

	if header.Nonce != (istanbulcommon.EmptyBlockNonce) && !bytes.Equal(header.Nonce[:], nonceAuthVote) && !bytes.Equal(header.Nonce[:], nonceDropVote) {
		return istanbulcommon.ErrInvalidNonce
	}

	// Ensure that the mix digest is zero as we don't have fork protection currently
	if header.MixDigest != types.IstanbulDigest {
		return istanbulcommon.ErrInvalidMixDigest
	}

	// Ensure that the block doesn't contain any uncles which are meaningless in Istanbul
	// if header.UncleHash != nilUncleHash {
	// 	return istanbulcommon.ErrInvalidUncleHash
	// }

	// Ensure that the block's difficulty is meaningful (may not be correct at this point)
	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		if header.Difficulty == nil || header.Difficulty.Cmp(big.NewInt(24)) != 0 {
			return istanbulcommon.ErrInvalidDifficulty
		}
	} else {
		if header.Difficulty == nil || header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
			return istanbulcommon.ErrInvalidDifficulty
		}
	}

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		//return istanbulcommon.ErrEmptyBlock
		return e.verifyCascadingFields(chain, header, validators, parents)
	}

	return e.verifyCascadingFields(chain, header, validators, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (e *Engine) verifyCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Check parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	// Ensure that the block's parent has right number and hash
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}

	// Ensure that the block's timestamp isn't too close to it's parent
	if parent.Time+e.cfg.BlockPeriod > header.Time {
		return istanbulcommon.ErrInvalidTimestamp
	}

	// Verify signer
	if err := e.verifySigner(chain, header, parents, validators); err != nil {
		return err
	}

	return e.verifyCommittedSeals(chain, header, parents, validators)
}

func (e *Engine) verifySigner(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return istanbulcommon.ErrUnknownBlock
	}

	// Resolve the authorization key and check against signers
	signer, err := e.Author(header)

	if err != nil {
		return err
	}

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		return nil
	} else {
		// Signer should be in the validator set of previous block's extraData.
		if _, v := validators.GetByAddress(signer); v == nil {
			log.Error("cavar|verifySigner-signer", "no", header.Number.Text(10), "header", header.Hash().Hex(), "sign", signer.Hex())
			for _, addr := range validators.List() {
				log.Error("cavar|verifySigner-val", "no", header.Number.Text(10), "header", header.Hash().Hex(), "val-addr", addr.Address().Hex())
			}
			return istanbulcommon.ErrUnauthorized
		}
	}

	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (e *Engine) verifyCommittedSeals(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	number := header.Number.Uint64()

	if number == 0 {
		// We don't need to verify committed seals in the genesis block
		return nil
	}

	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}
	committedSeal := extra.CommittedSeal

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		return nil
	} else {
		// The length of Committed seals should be larger than 0
		if len(committedSeal) == 0 {
			return istanbulcommon.ErrEmptyCommittedSeals
		}
	}

	validatorsCpy := validators.Copy()

	// Check whether the committed seals are generated by validators
	validSeal := 0
	committers, err := e.Signers(header)
	if err != nil {
		return err
	}

	for _, addr := range committers {
		if validatorsCpy.RemoveValidator(addr) {
			validSeal++
			continue
		}
		log.Error("caver|verifyCommittedSeals|committers", "no", header.Number, "header", header.Hash(), "addr", addr)
		for _, addr := range validators.List() {
			log.Error("caver|verifyCommittedSeals|validatorset", "no", header.Number.Text(10), "addr", addr)
		}
		for _, addr := range committers {
			log.Info("caver|verifyCommittedSeals|committedseals", "no", header.Number, "addr", addr)
		}
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= validators.F() {
		log.Error("caver|verifyCommittedSeals|validSeal", "no", header.Number.Text(10), "validSeal_len", validSeal, "validators.F()", validators.F())
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	return nil
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (e *Engine) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// if len(block.Uncles()) > 0 {
	// 	return istanbulcommon.ErrInvalidUncleHash
	// }
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (e *Engine) VerifySeal(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet) error {

	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return istanbulcommon.ErrUnknownBlock
	}

	// ensure that the difficulty equals to istanbulcommon.DefaultDifficulty
	if header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
		return istanbulcommon.ErrInvalidDifficulty
	}

	return e.verifySigner(chain, header, nil, validators)
}

func (e *Engine) Prepare(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet) error {
	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") &&
		header.Number.Cmp(common.Big0) > 0 {
		return errors.New("not a normal block")
	}
	header.Nonce = istanbulcommon.EmptyBlockNonce
	header.MixDigest = types.IstanbulDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// use the same difficulty for all blocks
	header.Difficulty = istanbulcommon.DefaultDifficulty
	var (
		validatorAddr []common.Address
		exchangerAddr []common.Address
		addrBigInt    []*big.Int
		rewardSeals   [][]byte
	)
	if c, ok := chain.(*core.BlockChain); ok {
		if header.Number.Uint64() == 1 {
			// Block 1 does not issue any rewards
			validatorAddr = make([]common.Address, 0)
		} else {
			// quorum Size
			random11Validators, err := c.Random11ValidatorWithOutProxy(c.CurrentHeader())
			if err != nil {
				log.Error("Prepare : invalid validators", err.Error())
				return errors.New("Prepare: invalid validators")
			}
			quorumSize := e.QuorumSize(random11Validators.Len())
			if quorumSize == 0 {
				log.Error("Prepare invalid quorum size", "no", header.Number, "size", quorumSize)
				return errors.New("invalid quorum size")
			}
			log.Info("Prepare quorum size", "no", header.Number, "size", quorumSize)
			// Get the header of the last normal block
			preHeader, err := getPreHash(chain, header)
			if err != nil {
				log.Error("Prepare get preHash err", "err", err, "no", header.Number, "hash", header.Hash().Hex())
				return err
			}
			log.Info("Prepare getPreHash ok", "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
			commiters, err := e.Signers(preHeader)
			if err != nil {
				log.Error("Prepare commit seal err", "err", err.Error(), "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				return err
			}
			if len(commiters) < quorumSize {
				log.Error("Prepare commiters len less than 7", "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				return errors.New("Prepare commiters len less than 7")
			}
			for _, v := range commiters {
				if len(validatorAddr) == quorumSize {
					break
				}
				// reward to onlineValidtors
				validatorAddr = append(validatorAddr, v)
			}
			for _, v := range validatorAddr {
				log.Info("Prepare: onlineValidator", "addr", v.Hex(), "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
			}
			// copy commitSeals to rewardSeals
			rewardSeals, err = e.copyCommitSeals(preHeader)
			if err != nil {
				log.Error("copy commitSeals err", "err", err, "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				return err
			}
		}

		// reward to openExchangers
		stakeList := c.GetStakerPool()
		for _, staker := range stakeList.Stakers {
			addrBigInt = append(addrBigInt, staker.Addr.Hash().Big())
		}
		benifitedStakers := stakeList.ValidatorByDistanceAndWeight(addrBigInt, 4, c.CurrentBlock().Header().Hash())
		exchangerAddr = append(exchangerAddr, benifitedStakers...)

		//new&update  at 20220523
		validatorPool, err := c.ReadValidatorPool(c.CurrentBlock().Header())
		if err != nil {
			log.Error("Prepare : validator pool err", err, err)
			return err
		}
		if validatorPool != nil && len(validatorPool.Validators) > 0 {
			//k:proxy,v:validator
			mp := make(map[string]*types.Validator, 0)
			for _, v := range validatorPool.Validators {
				if v.Proxy.String() != "0x0000000000000000000000000000000000000000" {
					mp[v.Proxy.String()] = v
				}
			}

			//If the reward address is on a proxy account, it will be restored to a pledge account
			for index, a := range validatorAddr {
				if v, ok := mp[a.String()]; ok {
					validatorAddr[index] = v.Addr
				}
			}
		}
	}

	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, validator.SortedAddresses(validators.List()), exchangerAddr, validatorAddr, rewardSeals)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	header.Time = parent.Time + e.cfg.BlockPeriod
	if header.Time < uint64(time.Now().Unix()) {
		header.Time = uint64(time.Now().Unix())
	}

	return nil
}

// copy commit seal to reward seals
func (e *Engine) copyCommitSeals(header *types.Header) ([][]byte, error) {
	// extract istanbul extra
	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	rewardSeals := make([][]byte, len(extra.CommittedSeal))
	for i, v := range extra.CommittedSeal {
		rewardSeals[i] = make([]byte, types.IstanbulExtraSeal)
		copy(rewardSeals[i][:], v[:])
	}
	return rewardSeals, nil
}

func getPreHash(chain consensus.ChainHeaderReader, header *types.Header) (*types.Header, error) {
	preHeader := chain.GetHeaderByHash(header.ParentHash)
	if preHeader == nil {
		return nil, errors.New("getPreHash : invalid preHeader")
	}
	if preHeader.Coinbase == (common.Address{}) {
		preHeader, err := getPreHash(chain, preHeader)
		if err != nil {
			return nil, err
		}
		return preHeader, nil
	}
	return preHeader, nil
}

func (e *Engine) PrepareEmpty(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet) error {

	if header.Coinbase != common.HexToAddress("0x0000000000000000000000000000000000000000") {
		return errors.New("not a empty block")
	}

	header.Nonce = istanbulcommon.EmptyBlockNonce
	header.MixDigest = types.IstanbulDigest

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same difficulty for all blocks

	// modification on 20221102 start
	//header.Difficulty = istanbulcommon.DefaultDifficulty
	header.Difficulty = big.NewInt(24)
	// modification on 20221102 end

	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, validator.SortedAddresses(validators.List()), nil, nil, nil)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp

	if header.Number.Cmp(common.Big0) > 0 {
		header.Time = parent.Time + 120
		if header.Time < uint64(time.Now().Unix()) {
			header.Time = uint64(time.Now().Unix())
		}
	} else {
		header.Time = parent.Time + e.cfg.BlockPeriod
		if header.Time < uint64(time.Now().Unix()) {
			header.Time = uint64(time.Now().Unix())
		}
	}

	return nil
}

func prepareExtra(header *types.Header, vals, exchangerAddr, validatorAddr []common.Address, rewardSeals [][]byte) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.IstanbulExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.IstanbulExtraVanity])

	ist := &types.IstanbulExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
		ExchangerAddr: exchangerAddr,
		ValidatorAddr: validatorAddr,
		RewardSeal:    rewardSeals,
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (e *Engine) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header) {
	if c, ok := chain.(*core.BlockChain); ok {
		istanbulExtra, err := types.ExtractIstanbulExtra(header)
		if err != nil {
			return
		}

		// 1 Adjust the weight of the previous block
		// pick 7 validator from rewardSeals
		var validatorAddr []common.Address
		if header.Number.Uint64() == 1 {
			// Block 1 does not issue any rewards
			validatorAddr = make([]common.Address, 0)
		} else {
			prePreHeader := c.GetHeaderByHash(c.CurrentHeader().ParentHash)
			random11Validators, err := c.Random11ValidatorWithOutProxy(prePreHeader)
			if err != nil {
				log.Error("Finalize : invalid validators", err.Error())
				return
			}
			if c.CurrentHeader().Coinbase == (common.Address{}) {
				// reduce 1 weight
				for _, v := range random11Validators.Validators {
					state.SubValidatorCoefficient(v.Address(), 10)
				}
			} else {
				validatorPool, err := c.ReadValidatorPool(c.CurrentBlock().Header())
				if err != nil {
					log.Error("Finalize : validator pool err", err, err)
					return
				}
				// quorum Size
				quorumSize := e.QuorumSize(random11Validators.Len())
				if quorumSize == 0 {
					log.Error("Finalize invalid quorum size", "no", header.Number, "size", quorumSize)
					return
				}
				log.Info("Finalize quorum size", "no", header.Number, "size", quorumSize)
				// Get the header of the last normal block
				preHeader, err := getPreHash(chain, header)
				if err != nil {
					log.Error("Finalize get preHash err", "err", err, "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
					return
				}
				log.Info("Finalize getPreHash ok", "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				// decode rewards
				// preHeader + currentRewadSeal
				rewarders, err := e.RecoverRewards(preHeader, istanbulExtra.RewardSeal)
				if err != nil {
					log.Error("Finalize rewarders err", "err", err.Error(), "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
					return
				}

				var rewardersWithoutProxy []common.Address
				for _, v := range istanbulExtra.ValidatorAddr {
					rewardersWithoutProxy = append(rewardersWithoutProxy, v)
				}
				if validatorPool != nil && len(validatorPool.Validators) > 0 {
					//k:proxy,v:validator
					mp := make(map[string]*types.Validator, 0)
					for _, v := range validatorPool.Validators {
						if v.Proxy.String() != "0x0000000000000000000000000000000000000000" {
							mp[v.Proxy.String()] = v
						}
					}
					//If the reward address is on a proxy account, it will be restored to a pledge account
					for index, a := range rewardersWithoutProxy {
						if v, ok := mp[a.String()]; ok {
							rewardersWithoutProxy[index] = v.Addr
						}
					}
				}

				for _, v := range rewardersWithoutProxy {
					log.Info("Finalize: onlineValidator", "addr", v.Hex(), "len", len(rewarders), "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
					// add 2 weight
					state.AddValidatorCoefficient(v, 20)
				}
				if len(rewarders) < quorumSize {
					log.Error("Finalize commiters len less than 7", "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
					return
				}
				for _, v := range rewarders {
					if len(validatorAddr) == quorumSize {
						break
					}
					// reward to onlineValidtors
					validatorAddr = append(validatorAddr, v)
				}

				if validatorPool != nil && len(validatorPool.Validators) > 0 {
					//k:proxy,v:validator
					mp := make(map[string]*types.Validator, 0)
					for _, v := range validatorPool.Validators {
						if v.Proxy.String() != "0x0000000000000000000000000000000000000000" {
							mp[v.Proxy.String()] = v
						}
					}
					//If the reward address is on a proxy account, it will be restored to a pledge account
					for index, a := range validatorAddr {
						if v, ok := mp[a.String()]; ok {
							validatorAddr[index] = v.Addr
						}
					}
				}
			}
		}
		for _, addr := range validatorAddr {
			log.Info("Finalize : CreateNFTByOfficial16", "ValidatorAddr=", addr.Hex(), "Coinbase", header.Coinbase.Hex(), "no", header.Number.Uint64())
		}
		for _, addr := range istanbulExtra.ExchangerAddr {
			log.Info("Finalize : CreateNFTByOfficial16", "ExchangerAddr=", addr.Hex(), "Coinbase", header.Coinbase.Hex(), "no", header.Number.Uint64())
		}

		state.CreateNFTByOfficial16(validatorAddr, istanbulExtra.ExchangerAddr, header.Number)

		/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
		header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
		header.UncleHash = nilUncleHash

	}
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (e *Engine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// Prepare reward address
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}

	// 1 Adjust the weight of the previous block
	if c, ok := chain.(*core.BlockChain); ok {
		if header.Number.Uint64() == 1 {
			log.Info("Finishing assembling first block")
		} else {
			prePreHeader := c.GetHeaderByHash(c.CurrentHeader().ParentHash)
			random11Validators, err := c.Random11ValidatorWithOutProxy(prePreHeader)
			if err != nil {
				log.Error("FinalizeAndAssemble : invalid validators", err.Error())
				return nil, err
			}
			// previous is an empty block
			if c.CurrentHeader().Coinbase == (common.Address{}) {
				// reduce 1 weight
				for _, v := range random11Validators.Validators {
					state.SubValidatorCoefficient(v.Address(), 10)
				}
			} else {
				validatorPool, err := c.ReadValidatorPool(c.CurrentBlock().Header())
				if err != nil {
					log.Error("FinalizeAndAssemble : validator pool err", err, err)
					return nil, err
				}
				// // previous is a normal block
				// // take out the rewardseal in the header directly, this is the data of the last normal block
				// // Get the header of the last normal block
				// preHeader, err := getPreHash(chain, header)
				// if err != nil {
				// 	log.Error("FinalizeAndAssemble get preHash err", "err", err, "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				// 	return nil, err
				// }
				// log.Info("FinalizeAndAssemble getPreHash ok", "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				// // decode rewards
				// // preHeader + currentRewadSeal
				// rewarders, err := e.RecoverRewards(preHeader, istanbulExtra.RewardSeal)
				// if err != nil {
				// 	log.Error("FinalizeAndAssemble rewarders err", "err", err.Error(), "preHeader", preHeader.Number, "preHash", preHeader.Hash().Hex(), "no", header.Number, "hash", header.Hash().Hex())
				// 	return nil, err
				// }
				var rewardersWithoutProxy []common.Address
				for _, v := range istanbulExtra.ValidatorAddr {
					rewardersWithoutProxy = append(rewardersWithoutProxy, v)
				}
				if validatorPool != nil && len(validatorPool.Validators) > 0 {
					//k:proxy,v:validator
					mp := make(map[string]*types.Validator, 0)
					for _, v := range validatorPool.Validators {
						if v.Proxy.String() != "0x0000000000000000000000000000000000000000" {
							mp[v.Proxy.String()] = v
						}
					}
					//If the reward address is on a proxy account, it will be restored to a pledge account
					for index, a := range rewardersWithoutProxy {
						if v, ok := mp[a.String()]; ok {
							rewardersWithoutProxy[index] = v.Addr
						}
					}
				}
				for _, v := range rewardersWithoutProxy {
					log.Info("FinalizeAndAssemble: onlineValidator", "addr", v.Hex(), "len", len(rewardersWithoutProxy), "no", header.Number, "hash", header.Hash().Hex())
					// add 2 weight
					state.AddValidatorCoefficient(v, 20)
				}
			}
		}
	}

	for _, addr := range istanbulExtra.ValidatorAddr {
		log.Info("FinalizeAndAssemble : CreateNFTByOfficial16", "ValidatorAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "no", header.Number.Uint64())
	}

	for _, addr := range istanbulExtra.ExchangerAddr {
		log.Info("FinalizeAndAssemble : CreateNFTByOfficial16", "ExchangerAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "no", header.Number.Uint64())
	}
	state.CreateNFTByOfficial16(istanbulExtra.ValidatorAddr, istanbulExtra.ExchangerAddr, header.Number)

	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, new(trie.Trie)), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (e *Engine) Seal(chain consensus.ChainHeaderReader, block *types.Block, validators istanbul.ValidatorSet) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		log.Info("caver|empty|Seal")
		parent := chain.GetHeader(header.ParentHash, number-1)
		if parent == nil {
			return block, consensus.ErrUnknownAncestor
		}

		return e.updateBlock(parent, block)
	} else {
		if _, v := validators.GetByAddress(e.signer); v == nil {
			log.Info("caver|seal|unauthorized", "no", number, "e.signer", e.signer.Hex())
			return block, istanbulcommon.ErrUnauthorized
		}
		parent := chain.GetHeader(header.ParentHash, number-1)
		if parent == nil {
			return block, consensus.ErrUnknownAncestor
		}

		return e.updateBlock(parent, block)
	}

}

// update timestamp and signature of the block based on its number of transactions
func (e *Engine) updateBlock(parent *types.Header, block *types.Block) (*types.Block, error) {
	// sign the hash
	header := block.Header()
	seal, err := e.sign(sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.IstanbulExtraSeal != 0 {
		return istanbulcommon.ErrInvalidSignature
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

func (e *Engine) SealHash(header *types.Header) common.Hash {
	return sigHash(header)
}

func (e *Engine) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return new(big.Int)
}

func (e *Engine) Validators(header *types.Header) ([]common.Address, error) {
	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}

	return extra.Validators, nil
}

func (e *Engine) Signers(header *types.Header) ([]common.Address, error) {
	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return []common.Address{}, err
	}
	committedSeal := extra.CommittedSeal
	proposalSeal := PrepareCommittedSeal(header.Hash())

	var addrs []common.Address
	// 1. Get committed seals from current header
	for _, seal := range committedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := istanbulcommon.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			return nil, istanbulcommon.ErrInvalidSignature
		}
		addrs = append(addrs, addr)
	}

	return addrs, nil
}

func (e *Engine) RecoverRewards(header *types.Header, rewardSeal [][]byte) ([]common.Address, error) {
	// extra, err := types.ExtractIstanbulExtra(header)
	// if err != nil {
	// 	return []common.Address{}, err
	// }
	// rewardSeal := extra.RewardSeal
	proposalSeal := PrepareCommittedSeal(header.Hash())

	var addrs []common.Address
	// 1. Get committed seals from current header
	for _, seal := range rewardSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := istanbulcommon.GetSignatureAddress(proposalSeal, seal)
		if err != nil {
			return nil, istanbulcommon.ErrInvalidSignature
		}
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (e *Engine) Address() common.Address {
	return e.signer
}

func (e *Engine) WriteVote(header *types.Header, candidate common.Address, authorize bool) error {
	header.Coinbase = candidate
	if authorize {
		copy(header.Nonce[:], nonceAuthVote)
	} else {
		copy(header.Nonce[:], nonceDropVote)
	}

	return nil
}

func (e *Engine) ReadVote(header *types.Header) (candidate common.Address, authorize bool, err error) {
	switch {
	case bytes.Equal(header.Nonce[:], nonceAuthVote):
		authorize = true
	case bytes.Equal(header.Nonce[:], nonceDropVote):
		authorize = false
	default:
		return common.Address{}, false, istanbulcommon.ErrInvalidVote
	}

	return header.Coinbase, authorize, nil
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.IstanbulExtraSeal {
			return istanbulcommon.ErrInvalidCommittedSeals
		}
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(istanbulExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(ibfttypes.MsgCommit)})
	return buf.Bytes()
}

// helper func--------------------------------//

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func (e *Engine) QuorumSize(valSize int) int {
	return 2*(int(math.Ceil(float64(valSize)/3))-1) + 1
}
