package miner

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
	"time"
)

func (c *Certify) SendSignToOtherPeer(addr common.Address, height *big.Int) {
	log.Info("start SendSignToOtherPeer", "Address", addr.Hex(), "Height:", height)
	ques := &SignatureData{
		Address:   addr,
		Height:    height,
		Timestamp: uint64(time.Now().Unix()),
	}
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

	log.Info("Certify.GatherOtherPeerSignature >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	if c.miner.GetWorker().chain.CurrentHeader().Number.Cmp(height) >= 0 {
		return errors.New("GatherOtherPeerSignature: msg height < chain Number")
	}

	c.proofStatePool.ClearPrev(c.miner.GetWorker().chain.CurrentHeader().Number)

	if c.proofStatePool.proofs[height] == nil {
		ps := newProofState(validator, validator)
		ps.receiveValidatorsSum = big.NewInt(0)
		ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, c.stakers.StakeBalance(validator))
		ps.onlineValidator = make(OnlineValidator)
		ps.onlineValidator.Add(validator)
		ps.height = new(big.Int).Set(height)

		if c.self != validator {
			// add my own amount
			ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, c.stakers.StakeBalance(c.self))
			ps.onlineValidator.Add(c.self)
		}

		c.proofStatePool.proofs[height] = ps
		c.signatureResultCh <- height
		log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 1")
		return nil
	}

	curProofs := c.proofStatePool.proofs[height]
	if curProofs.onlineValidator.Has(validator) {
		return errors.New("GatherOtherPeerSignature: validator exist")
	}
	c.proofStatePool.proofs[height].onlineValidator.Add(validator)
	c.proofStatePool.proofs[height].receiveValidatorsSum = new(big.Int).Add(c.proofStatePool.proofs[height].receiveValidatorsSum, c.stakers.StakeBalance(validator))
	log.Info("Certify.GatherOtherPeerSignature", "receiveValidatorsSum", c.proofStatePool.proofs[height].receiveValidatorsSum, "heigh", height)
	c.signatureResultCh <- height
	log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 2")
	return nil
}
