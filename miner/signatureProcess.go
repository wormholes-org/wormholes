package miner

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (c *Certify) AssembleAndBroadcastMessage(stakes *types.ValidatorList, height *big.Int) {
	//log.Info("AssembleAndBroadcastMessage", "validators len", len(c.stakers.Validators), "sender", c.addr, "vote index", c.voteIndex, "round", c.round)
	vote := stakes.Validators[c.voteIndex]
	var voteAddress common.Address
	if vote.Proxy == (common.Address{}) {
		voteAddress = vote.Addr
	} else {
		voteAddress = vote.Proxy
	}
	c.voteIndex++
	if c.voteIndex == len(stakes.Validators) {
		c.voteIndex = 0
		c.round++
	}

	log.Info("azh|start to vote", "validators", len(stakes.Validators), "vote", vote, "height:", height)
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
	log.Info("GatherOtherPeerSignature", "c.proofStatePool", c.proofStatePool)
	if proof, ok := c.proofStatePool.proofs[height.Uint64()]; !ok {
		proof = newProofState(height, nil, nil, false, c.self, nil, nil)
		proof.count = 0
		proof.onlineValidator = append(proof.onlineValidator, validator)
		proof.emptyBlockMessages = append(proof.emptyBlockMessages, encQues)
		return nil
	} else {
		if proof.onlineValidator.Has(validator) {
			return errors.New("GatherOtherPeerSignature: validator exist")
		}
		proof.onlineValidator = append(proof.onlineValidator, validator)
		proof.emptyBlockMessages = append(proof.emptyBlockMessages, encQues)
		if proof.validatorList != nil {
			proof.count++
			validatorBalance := proof.validatorList.StakeBalance(validator)
			weightBalance := new(big.Int).Mul(validatorBalance, big.NewInt(types.DEFAULT_VALIDATOR_COEFFICIENT))
			proof.receiveValidatorsSum = new(big.Int).Add(proof.receiveValidatorsSum, weightBalance)
		}
		if proof.receiveValidatorsSum.Cmp(proof.targetWeightBalance) > 0 && len(proof.onlineValidator) == proof.count {
			proof.empty = true
			c.signatureResultCh <- VoteResult{
				proof.height,
				proof.empty,
				proof.GetAllAddress(),
				proof.GetAllMessage(),
			}
		}
		return nil
	}
}
