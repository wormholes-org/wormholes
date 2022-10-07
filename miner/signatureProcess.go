package miner

import (
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

	c.proofStatePool.ClearPrev(c.miner.GetWorker().chain.CurrentHeader().Number)
	
	for k, p := range c.proofStatePool.proofs {
		if k.Cmp(height) == 0 {
			// Proof data exists for this height
			if p.onlineValidator.Has(validator) {
				continue
			}
			p.onlineValidator.Add(validator)
			c.receiveValidatorsSum = new(big.Int).Add(c.stakers.StakeBalance(validator), c.receiveValidatorsSum)
			//log.Info("c.receiveValidatorsSum", "c.receiveValidatorsSum", c.receiveValidatorsSum)
			c.validators = append(c.validators, validator)
			c.signatureResultCh <- c.receiveValidatorsSum
			return nil
		} else {
			//log.Info("GatherOtherPeerSignature, not the same height, continue!")
			continue
		}
	}
	// No proof data exists for this height
	ps := newProofState(validator, validator)
	c.proofStatePool.proofs[height] = ps
	c.receiveValidatorsSum = new(big.Int).Add(c.stakers.StakeBalance(validator), c.receiveValidatorsSum)
	//log.Info("c.receiveValidatorsSum", "c.receiveValidatorsSum", c.receiveValidatorsSum)
	c.validators = append(c.validators, validator)

	//log.Info("c.validators", "c.validators", c.validators, "c.receiveValidatorsSum", c.receiveValidatorsSum)
	c.signatureResultCh <- c.receiveValidatorsSum
	return nil
}
