package miner

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (c *Certify) AssembleAndBroadcastMessage(height *big.Int) {
	//log.Info("AssembleAndBroadcastMessage", "validators len", len(c.stakers.Validators), "sender", c.addr, "vote index", c.voteIndex, "round", c.round)
	vote := c.stakers.Validators[c.voteIndex]
	var voteAddress common.Address
	if vote.Proxy == (common.Address{}) {
		voteAddress = vote.Addr
	} else {
		voteAddress = vote.Proxy
	}
	c.voteIndex++
	if c.voteIndex == len(c.stakers.Validators) {
		c.voteIndex = 0
		c.round++
	}

	log.Info("azh|start to vote", "validators", len(c.stakers.Validators), "vote", vote, "height:", height)
	err, payload := c.assembleMessage(height, voteAddress)
	if err != nil {
		return
	}

	if voteAddress == c.self {
		currentBlock := c.miner.GetWorker().eth.BlockChain().CurrentBlock()
		c.miner.GetWorker().mux.Post(core.NewMinedBlockEvent{Block: currentBlock})

		emptyMsg := types.EmptyMessageEvent{
			Sender:  c.self,
			Height:  height,
			Payload: payload,
		}
		go c.eventMux.Post(emptyMsg)
	} else {
		if miner, ok := c.miner.(*Miner); ok {
			miner.broadcaster.BroadcastEmptyBlockMsg(payload)
		}
	}
	//log.Info("AssembleAndBroadcastMessage end")
}

func (c *Certify) GatherOtherPeerSignature(validator common.Address, height *big.Int, encQues []byte) error {
	var weightBalance *big.Int
	log.Info("GatherOtherPeerSignature", "c.proofStatePool", c.proofStatePool)
	curProofs := c.proofStatePool.proofs[height.Uint64()]
	if curProofs.onlineValidator.Has(validator) {
		return errors.New("GatherOtherPeerSignature: validator exist")
	}
	c.proofStatePool.proofs[height.Uint64()].onlineValidator = append(c.proofStatePool.proofs[height.Uint64()].onlineValidator, validator)
	c.proofStatePool.proofs[height.Uint64()].emptyBlockMessages = append(c.proofStatePool.proofs[height.Uint64()].emptyBlockMessages, encQues)
	//coe, err = c.miner.GetWorker().getValidatorCoefficient(validator)
	//if err != nil {
	//	return err
	//}
	//weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(int64(coe)))
	validatorBalance := c.stakers.StakeBalance(validator)
	weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(types.DEFAULT_VALIDATOR_COEFFICIENT))
	//weightBalance.Div(weightBalance, big.NewInt(10))
	c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum = new(big.Int).Add(c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, weightBalance)
	//log.Info("Certify.GatherOtherPeerSignature", "validator", validator.Hex(), "balance", validatorBalance, "average coe", averageCoefficient, "weightBalance", weightBalance, "receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "height", height.Uint64())
	//log.Info("Certify.GatherOtherPeerSignature", "receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "heigh", height)
	c.signatureResultCh <- height
	//log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 2")
	return nil
}
