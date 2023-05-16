package miner

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
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
	log.Info("azh|start to vote", "validators", len(c.stakers.Validators), "index", c.voteIndex, "round", c.round, "vote", vote, "height:", height)
	c.voteIndex++
	if c.voteIndex == len(c.stakers.Validators) {
		c.voteIndex = 0
		c.round++
	}

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
		c.requestEmpty <- payload
	}
	//log.Info("AssembleAndBroadcastMessage end")
}

func (c *Certify) GatherOtherPeerSignature(validator common.Address, height *big.Int, encQues []byte) error {
	var weightBalance *big.Int
	log.Info("GatherOtherPeerSignature", "c.proofStatePool", c.proofStatePool)
	if _, ok := c.proofStatePool.proofs[height.Uint64()]; !ok {
		_, proposerMessage := c.assembleMessage(height, c.self)
		ps := newProofState(c.self, proposerMessage, height)
		ps.receiveValidatorsSum = big.NewInt(0)
		//coe, err = c.miner.GetWorker().getValidatorCoefficient(validator)
		//if err != nil {
		//	return err
		//}
		//weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(int64(coe)))
		validatorBalance := c.stakers.StakeBalance(validator)
		weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(types.DEFAULT_VALIDATOR_COEFFICIENT))
		//weightBalance.Div(weightBalance, big.NewInt(10))
		ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, weightBalance)
		//log.Info("Certify.GatherOtherPeerSignature", "validator", validator.Hex(), "balance", validatorBalance, "average coe", averageCoefficient, "weightBalance", weightBalance, "receiveValidatorsSum", ps.receiveValidatorsSum, "height", height.Uint64())
		//ps.onlineValidator.Add(validator)
		//ps.height = new(big.Int).Set(height)
		ps.onlineValidator = append(ps.onlineValidator, validator)
		ps.emptyBlockMessages = append(ps.emptyBlockMessages, encQues)

		c.proofStatePool.proofs[height.Uint64()] = ps

		c.signatureResultCh <- VoteResult{
			height,
			ps.receiveValidatorsSum,
			ps.GetAllAddress(c.stakers),
			ps.GetAllEmptyMessage(),
		}
		//log.Info("GatherOtherPeerSignature", "height", height)
		//c.signatureResultCh <- height
		//log.Info("GatherOtherPeerSignature end", "height", height)
		//log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 1")
		return nil
	}

	curProofs := c.proofStatePool.proofs[height.Uint64()]
	if curProofs.onlineValidator.Has(validator) {
		return errors.New("GatherOtherPeerSignature: validator exist")
	}
	curProofs.onlineValidator = append(curProofs.onlineValidator, validator)
	curProofs.emptyBlockMessages = append(curProofs.emptyBlockMessages, encQues)
	//coe, err = c.miner.GetWorker().getValidatorCoefficient(validator)
	//if err != nil {
	//	return err
	//}
	//weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(int64(coe)))
	validatorBalance := c.stakers.StakeBalance(validator)
	weightBalance = new(big.Int).Mul(validatorBalance, big.NewInt(types.DEFAULT_VALIDATOR_COEFFICIENT))
	//weightBalance.Div(weightBalance, big.NewInt(10))
	curProofs.receiveValidatorsSum = new(big.Int).Add(curProofs.receiveValidatorsSum, weightBalance)
	c.signatureResultCh <- VoteResult{
		height,
		curProofs.receiveValidatorsSum,
		curProofs.GetAllAddress(c.stakers),
		curProofs.GetAllEmptyMessage(),
	}
	//log.Info("Certify.GatherOtherPeerSignature", "validator", validator.Hex(), "balance", validatorBalance, "average coe", averageCoefficient, "weightBalance", weightBalance, "receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "height", height.Uint64())
	//log.Info("Certify.GatherOtherPeerSignature", "receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "heigh", height)
	//c.signatureResultCh <- height
	//log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 2")
	return nil
}
