package miner

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (c *Certify) SendSignToOtherPeer(addr common.Address, height *big.Int) {
	log.Info("start SendSignToOtherPeer", "Address", addr.Hex(), "Height:", height)
	ques := &SignatureData{Address: addr, Height: height}
	encQues, err := Encode(ques)
	if err != nil {
		log.Error("Failed to encode", "subject", err)
		return
	}
	c.broadcast(c.Address(), &Msg{
		Code: SendSignMsg,
		Msg:  encQues,
	})
}

func (c *Certify) GatherOtherPeerSignature(validator common.Address, height *big.Int, encQues []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.miner.GetWorker().chain.CurrentHeader().Number.Cmp(height) >= 0 {
		return errors.New("GatherOtherPeerSignature: msg height < chain Number")
	}

	c.proofStatePool.ClearPrev(c.miner.GetWorker().chain.CurrentHeader().Number)

	if c.proofStatePool.proofs[height] == nil {
		ps := newProofState(validator, validator)
		c.proofStatePool.proofs[height] = ps
		c.receiveValidatorsSum = new(big.Int).Add(c.stakers.StakeBalance(validator), c.receiveValidatorsSum)
		c.validators = append(c.validators, validator)
		c.signatureResultCh <- c.receiveValidatorsSum
		return nil
	}

	curProofs := c.proofStatePool.proofs[height]
	if curProofs.onlineValidator.Has(validator) {
		return errors.New("GatherOtherPeerSignature: validator exist")
	}
	curProofs.onlineValidator.Add(validator)
	c.receiveValidatorsSum = new(big.Int).Add(c.stakers.StakeBalance(validator), c.receiveValidatorsSum)
	//log.Info("c.receiveValidatorsSum", "c.receiveValidatorsSum", c.receiveValidatorsSum)
	c.validators = append(c.validators, validator)
	c.signatureResultCh <- c.receiveValidatorsSum
	return nil
}
