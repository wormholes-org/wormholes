package ibftengine

import (
	"bytes"
	"encoding/binary"
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
		err := e.VerifyProposalHeader(chain, block.Header(), nil, validators)
		if err == nil || err == istanbulcommon.ErrEmptyCommittedSeals {
			// ignore errEmptyCommittedSeals error because we don't have the committed seals yet
			return 0, nil
		} else if err == consensus.ErrFutureBlock {
			return time.Until(time.Unix(int64(block.Header().Time), 0)), consensus.ErrFutureBlock
		}
		return 0, err
	}

}

func (e *Engine) VerifyProposalHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	return e.verifyProposalHeader(chain, header, parents, validators)
}
func (e *Engine) verifyProposalHeader(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
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
	if header.Difficulty == nil || header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
		return istanbulcommon.ErrInvalidDifficulty
	}

	return e.verifyProposalCascadingFields(chain, header, validators, parents)
}

func (e *Engine) verifyProposalCascadingFields(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet, parents []*types.Header) error {
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
	return nil
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
	if header.Difficulty == nil || header.Difficulty.Cmp(istanbulcommon.DefaultDifficulty) != 0 {
		return istanbulcommon.ErrInvalidDifficulty
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

	return e.verifyOnlineSeals(chain, header, parents, validators)
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

func (e *Engine) verifyOnlineProof(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	//TODO Parse out the online validator saved in uncles
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

func (e *Engine) verifyOnlineSeals(chain consensus.ChainHeaderReader, header *types.Header, parents []*types.Header, validators istanbul.ValidatorSet) error {
	number := header.Number.Uint64()

	if number == 0 {
		// We don't need to verify committed seals in the genesis block
		return nil
	}

	if header.Coinbase == (common.Address{}) && header.Number.Cmp(common.Big0) > 0 {
		log.Info("verifyOnlineSeals empty address", "no", header.Number.Uint64(), "coinbase", header.Coinbase.Hex())
		return nil
	}

	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}
	onlineSeal := extra.OnlineSeal

	log.Info("verifyOnlineSeals onlineSeal", "no", header.Number, "hash", header.Hash(), "coinbase", header.Coinbase, "onlineSeal", onlineSeal)
	var commitSeals *types.ProposerBlock
	var committers []common.Address
	err2 := rlp.DecodeBytes(onlineSeal, &commitSeals)
	if err2 != nil {
		log.Error("verifyOnlineSeals failed rlp decode OnlineSeal", "err2", err2.Error())
		return nil
	}
	if commitSeals != nil {
		log.Info("verifyOnlineSeals commitSeals info", "round", commitSeals.Round, "sequence", commitSeals.Sequence, "commit", commitSeals.Commit, "digies", commitSeals.Digest)

		var msgs []*ibfttypes.Message
		err := rlp.DecodeBytes(commitSeals.Commit, &msgs)
		if err != nil {
			log.Error("verifyOnlineSeals failed rlp decode commit msgs after", "err", err.Error())
		}
		for _, v := range msgs {
			committers = append(committers, v.Address)
		}
	}

	validatorsCpy := validators.Copy()
	// Check whether the committed seals are generated by validators
	validSeal := 0

	for _, addr := range validatorsCpy.List() {
		log.Info("verifyOnlineSeals validators", "no", header.Number, "addr", addr)
	}
	for _, addr := range committers {
		log.Info("verifyOnlineSeals committers", "no", header.Number, "addr", addr)
	}
	for _, addr := range committers {
		if validatorsCpy.RemoveValidator(addr) {
			validSeal++
			continue
		}
		return istanbulcommon.ErrInvalidCommittedSeals
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= validators.F() {
		log.Error("verifyOnlineSeals validSeal", "no", header.Number, "validSeal_len", validSeal, "validators.F()", validators.F())
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
	// var onlineValidators []common.Address
	// ibftCore := e.backend.GetCore()
	// if ibftCore != nil {
	// 	ibftCore.GetOnlineProofsMu().Lock()
	// 	vals := ibftCore.GetOnlineValidators()
	// 	if _, ok := vals[header.Number.Uint64()]; ok {
	// 		for _, v := range vals[header.Number.Uint64()].Validators {
	// 			onlineValidators = append(onlineValidators, v.Address)
	// 		}
	// 	}
	// 	ibftCore.GetOnlineProofsMu().Unlock()
	// }

	// for _, v := range onlineValidators {
	// 	log.Info("Prepare: onlineValidators", "no", header.Number, "onlineValidators", v)
	// }

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
		//validatorAddr []common.Address
		exchangerAddr []common.Address
		addrBigInt    []*big.Int
	)
	if c, ok := chain.(*core.BlockChain); ok {
		// reward to openExchangers
		//stakeList := c.ReadStakePool(c.CurrentBlock().Header())
		stakeList := c.GetStakerPool()
		for _, staker := range stakeList.Stakers {
			addrBigInt = append(addrBigInt, staker.Addr.Hash().Big())
		}
		benifitedStakers := stakeList.ValidatorByDistanceAndWeight(addrBigInt, 4, c.CurrentBlock().Header().Hash())
		exchangerAddr = append(exchangerAddr, benifitedStakers...)

		// reward to validators
		// validatorList, err := c.Random11ValidatorFromPool(c.CurrentBlock().Header())
		// if err != nil {
		// 	return err
		// }
		// benifitedValidators := c.RandomNValidatorFromEleven(6, validatorList, c.CurrentBlock().Header().Hash())
		// validatorAddr = append(validatorAddr, benifitedValidators...)

		// // reward to miner
		// validatorAddr = append(validatorAddr, header.Coinbase)

		// rewardToValidators, err := PickRewardValidators(c, onlineValidators)
		// if err != nil {
		// 	return err
		// }
		// validatorAddr = append(validatorAddr, rewardToValidators...)

		//new&update  at 20220523
		// validatorPool, err := c.ReadValidatorPool(c.CurrentBlock().Header())
		// if err != nil {
		// 	log.Error("Prepare : validator pool err", err, err)
		// 	return err
		// }
		// if validatorPool != nil && len(validatorPool.Validators) > 0 {
		// 	//k:proxy,v:validator
		// 	mp := make(map[string]*types.Validator, 0)
		// 	for _, v := range validatorPool.Validators {
		// 		if v.Proxy.String() != "0x0000000000000000000000000000000000000000" {
		// 			mp[v.Proxy.String()] = v
		// 		}
		// 	}

		// 	//If the reward address is on a proxy account, it will be restored to a pledge account
		// 	for index, a := range validatorAddr {
		// 		if v, ok := mp[a.String()]; ok {
		// 			validatorAddr[index] = v.Addr
		// 		}
		// 	}
		// }
	}
	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, validator.SortedAddresses(validators.List()), exchangerAddr, nil)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		header.Time = parent.Time + 120
	} else {
		header.Time = parent.Time + e.cfg.BlockPeriod
		if header.Time < uint64(time.Now().Unix()) {
			header.Time = uint64(time.Now().Unix())
		}
	}

	return nil
}

func (e *Engine) PrepareEmpty(chain consensus.ChainHeaderReader, header *types.Header, validators istanbul.ValidatorSet) error {
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
	// add validators in snapshot to extraData's validators section
	extra, err := prepareExtra(header, validator.SortedAddresses(validators.List()), nil, nil)
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp

	if header.Coinbase == common.HexToAddress("0x0000000000000000000000000000000000000000") && header.Number.Cmp(common.Big0) > 0 {
		header.Time = parent.Time + e.cfg.BlockPeriod
	} else {
		if header.Time < uint64(time.Now().Unix()) {
			header.Time = uint64(time.Now().Unix())
		}
	}

	return nil
}

func prepareExtra(header *types.Header, vals, exchangerAddr, validatorAddr []common.Address) ([]byte, error) {
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
	// // Prepare reward address
	// istanbulExtra, err := types.ExtractIstanbulExtra(header)
	// if err != nil {
	// 	return
	// }

	// if c, ok := chain.(*core.BlockChain); ok {
	// 	// empty block  reduce 0.1weight and normal block add 0.5weight
	// 	random11Validators, err := c.Random11ValidatorWithOutProxy(c.CurrentHeader())
	// 	if err != nil {
	// 		log.Error("Finalize : invalid validators", err.Error())
	// 		return
	// 	}
	// 	if header.Coinbase == (common.Address{}) {
	// 		// reduce 1 weight
	// 		for _, v := range random11Validators.Validators {
	// 			state.SubValidatorCoefficient(v.Address(), 10)
	// 		}
	// 	} else {
	// 		// add 2 weight
	// 		for _, v := range random11Validators.Validators {
	// 			state.AddValidatorCoefficient(v.Addr, 20)
	// 		}
	// 	}
	// }

	// log.Info("CreateNFTByOfficial16 start", "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// for _, addr := range istanbulExtra.ValidatorAddr {
	// 	log.Info("CreateNFTByOfficial16", "ValidatorAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// }
	// for _, addr := range istanbulExtra.ExchangerAddr {
	// 	log.Info("CreateNFTByOfficial16", "ExchangerAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// }
	// state.CreateNFTByOfficial16(istanbulExtra.ValidatorAddr, istanbulExtra.ExchangerAddr, header.Number)
	// log.Info("CreateNFTByOfficial16 end", "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())

	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (e *Engine) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// // Prepare reward address
	// istanbulExtra, err := types.ExtractIstanbulExtra(header)
	// if err != nil {
	// 	return nil, err
	// }

	// if c, ok := chain.(*core.BlockChain); ok {
	// 	// empty block  reduce 0.1weight and normal block add 0.5weight
	// 	random11Validators, err := c.Random11ValidatorWithOutProxy(c.CurrentHeader())
	// 	if err != nil {
	// 		log.Error("FinalizeAndAssemble : invalid validators", err.Error())
	// 		return nil, err
	// 	}
	// 	if header.Coinbase == (common.Address{}) {
	// 		// reduce 1 weight
	// 		for _, v := range random11Validators.Validators {
	// 			state.SubValidatorCoefficient(v.Address(), 10)
	// 		}
	// 	} else {
	// 		// add 2 weight
	// 		for _, v := range random11Validators.Validators {
	// 			state.AddValidatorCoefficient(v.Addr, 20)
	// 		}
	// 	}
	// }

	// log.Info("CreateNFTByOfficial16 start", "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// for _, addr := range istanbulExtra.ValidatorAddr {
	// 	log.Info("CreateNFTByOfficial16", "ValidatorAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// }

	// for _, addr := range istanbulExtra.ExchangerAddr {
	// 	log.Info("CreateNFTByOfficial16", "ExchangerAddr=", addr.Hex(), "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())
	// }
	// state.CreateNFTByOfficial16(nil, istanbulExtra.ExchangerAddr, header.Number)

	// log.Info("CreateNFTByOfficial16 end", "Coinbase=", header.Coinbase.Hex(), "height", header.Number.Uint64())

	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, new(trie.Trie)), nil
}

func (e *Engine) FinalizeOnlineProofBlk(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	state.AddBalance(common.HexToAddress("0x0000000000000000000000000000000000000001"), big.NewInt(1))

	/// No block rewards in Istanbul, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	//header.UncleHash = nilUncleHash

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, uncles, receipts, new(trie.Trie)), nil
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

// func WriteBenifitedAddr(h *types.Header, benifitedAddr []common.Address) error {
// 	if len(benifitedAddr) == 0 {
// 		return istanbulcommon.ErrInvalidBenifitedAddr
// 	}

// 	istanbulExtra, err := types.ExtractIstanbulExtra(h)
// 	if err != nil {
// 		return err
// 	}

// 	// istanbulExtra.BeneficiaryAddr = make([]common.Address, len(benifitedAddr))
// 	// copy(istanbulExtra.BeneficiaryAddr, benifitedAddr)

// 	payload, err := rlp.EncodeToBytes(&istanbulExtra)
// 	if err != nil {
// 		return err
// 	}

// 	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
// 	return nil
// }

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(ibfttypes.MsgCommit)})
	return buf.Bytes()
}

// helper func--------------------------------//

// func RecoverOnlineValidator(data []byte) ([]common.Address, error) {
// 	if len(data) < types.OnlineValidatorVanity {
// 		return []common.Address{}, errors.New("Online proof validator data is insufficient")
// 	}

// 	var addrs []common.Address
// 	for i := 32; i < types.OnlineValidatorVanity; i += 100 {
// 		//buffer := new(buffer.Buffer)
// 		data[32:]

// 	}

// }

func BytesToInt(bys []byte) int {
	bytebuff := bytes.NewBuffer(bys)
	var data int32
	binary.Read(bytebuff, binary.BigEndian, &data)
	return int(data)
}

func PickRewardValidators(chain *core.BlockChain, onlineValidators []common.Address) ([]common.Address, error) {
	header := chain.CurrentHeader()
	totalValidators, err := chain.ReadValidatorPool(header)
	if err != nil {
		log.Error("PickRewardValidators : invalid validator list", "err", err)
		return []common.Address{}, err
	}
	validator11 := make(map[common.Address]bool)
	elevenAddrs, err := chain.Random11ValidatorFromPool(header)
	if err != nil {
		log.Error("PickRewardValidators : invalid validator list", "err", err)
		return []common.Address{}, err
	}
	online7 := make(map[common.Address]bool)
	coinbase := chain.Coinbase()

	var reward7 []common.Address

	if err != nil {
		return nil, err
	}
	for _, v := range elevenAddrs.Validators {
		validator11[v.Addr] = true
	}
	for _, v := range onlineValidators {
		if _, ok := validator11[v]; ok {
			online7[v] = true
		}
	}

	for _, v := range onlineValidators {
		if _, ok := validator11[v]; !ok {
			delete(online7, v)
		}
	}
	final7 := make(map[common.Address]bool)
	if len(online7) >= 7 {
		for k, _ := range online7 {
			final7[k] = true
			if len(final7) >= 6 {
				if _, ok := validator11[coinbase]; ok {
					if _, ok1 := final7[coinbase]; !ok1 {
						final7[coinbase] = true
					}
				}
			}
			if len(final7) >= 7 {
				break
			}
		}
	}

	if len(final7) == 7 {
		for k, _ := range final7 {
			reward7 = append(reward7, k)
		}
		return reward7, nil
	} else if len(final7) > 7 {
		log.Error("reward7s > 7", "len", len(reward7))
	} else {
		log.Info("reward7s < 7", "len", len(reward7))
	}
	for _, v := range totalValidators.Validators {
		if len(reward7) >= 7 {
			break
		}
		if _, ok := validator11[v.Addr]; ok {
			reward7 = append(reward7, v.Addr)
		}
		if _, ok := validator11[v.Proxy]; ok {
			reward7 = append(reward7, v.Proxy)
		}
	}
	if len(reward7) != 7 {
		log.Error("reward7s != 7", "len", len(reward7))
	}
	log.Info("PickRewardValidators", "no", header.Number.Uint64()+1, "validators", reward7, "len", len(reward7))
	return reward7, nil
}
