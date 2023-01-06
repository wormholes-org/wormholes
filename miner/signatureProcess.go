package miner

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"math/big"
)

func (c *Certify) AssembleAndStoreMessage(vote common.Address, height *big.Int) {
	log.Info("start SendSignToOtherPeer", "Address", vote.Hex(), "Height:", height)
	ques := &types.SignatureData{
		Vote:   vote,
		Height: height,
		//Timestamp: uint64(time.Now().Unix()),
	}
	encQues, err := Encode(ques)
	if err != nil {
		log.Error("Failed to encode", "subject", err)
		return
	}

	msg := &types.EmptyMsg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := c.signMessage(msg)
	if err != nil {
		log.Error("signMessage err", err)
		return
	}

	if _, ok := c.messageList.Load(payload); ok {
		return
	} else {
		c.messageList.Store(payload, true)
	}
}

//func (c *Certify) SendSignToOtherPeer(vote common.Address, height *big.Int) {
//	log.Info("start SendSignToOtherPeer", "Address", vote.Hex(), "Height:", height)
//	ques := &types.SignatureData{
//		Vote:   vote,
//		Height: height,
//		//Timestamp: uint64(time.Now().Unix()),
//	}
//	encQues, err := Encode(ques)
//	if err != nil {
//		log.Error("Failed to encode", "subject", err)
//		return
//	}
//	c.broadcast(&types.EmptyMsg{
//		Code: SendSignMsg,
//		Msg:  encQues,
//	})
//}

//func (c *Certify) GetSignedMessage(height *big.Int) ([]byte, error) {
//	ques := &types.SignatureData{
//		Vote:   c.self,
//		Height: height,
//		//Timestamp: uint64(time.Now().Unix()),
//	}
//	encQues, err := Encode(ques)
//	if err != nil {
//		log.Error("GetSignedMessage Failed to encode", "subject", err)
//		return nil, err
//	}
//
//	msg := &types.EmptyMsg{
//		Code: SendSignMsg,
//		Msg:  encQues,
//	}
//
//	payload, err := c.signMessage(msg)
//	if err != nil {
//		log.Error("GetSignedMessage signMessage err", err)
//		return nil, err
//	}
//
//	return payload, nil
//}

func (c *Certify) GatherOtherPeerSignature(addr, vote common.Address, height *big.Int, encQues []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	//log.Info("Certify.GatherOtherPeerSignature >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	if c.self != vote {
		log.Info("GatherOtherPeerSignature", "vote", vote)
		return nil
	}

	if c.miner.GetWorker().chain.CurrentHeader().Number.Cmp(height) >= 0 {
		return errors.New("GatherOtherPeerSignature: msg height < chain Number")
	}

	emptyAddress := common.Address{}
	validator := c.stakers.GetValidatorAddr(addr)
	if validator == emptyAddress {
		return errors.New("not a validator")
	}

	//log.Info("Certify.GatherOtherPeerSignature", "c.miner.GetWorker().chain.CurrentHeader().Number", c.miner.GetWorker().chain.CurrentHeader().Number,
	//	"height", height, "c.proofStatePool.proofs[height] == nil 1", c.proofStatePool.proofs[height.Uint64()] == nil)
	//c.proofStatePool.ClearPrev(c.miner.GetWorker().chain.CurrentHeader().Number)
	//log.Info("Certify.GatherOtherPeerSignature", "c.miner.GetWorker().chain.CurrentHeader().Number", c.miner.GetWorker().chain.CurrentHeader().Number,
	//	"height", height, "c.proofStatePool.proofs[height] == nil 2", c.proofStatePool.proofs[height.Uint64()] == nil)
	averageCoefficient, err := c.miner.GetWorker().GetAverageCoefficient() // need to divide 10
	if err != nil {
		return err
	}
	var weightBalance *big.Int
	//var coe uint8
	//var err error
	if _, ok := c.proofStatePool.proofs[height.Uint64()]; !ok {
		ps := newProofState(validator, validator)
		ps.receiveValidatorsSum = big.NewInt(0)
		//coe, err = c.miner.GetWorker().getValidatorCoefficient(validator)
		//if err != nil {
		//	return err
		//}
		//weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(validator), big.NewInt(int64(coe)))
		weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(validator), big.NewInt(int64(averageCoefficient)))
		weightBalance.Div(weightBalance, big.NewInt(10))
		ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, weightBalance)
		//		log.Info("Certify.GatherOtherPeerSignature", "validator", validator.Hex(),
		//			"balance", c.stakers.StakeBalance(validator), "average coe", averageCoefficient, "weightBalance", weightBalance,
		//			"receiveValidatorsSum", ps.receiveValidatorsSum, "height", height.Uint64())
		ps.onlineValidator = make(OnlineValidator)
		ps.onlineValidator.Add(validator)
		ps.height = new(big.Int).Set(height)
		ps.emptyBlockMessages = append(ps.emptyBlockMessages, encQues)

		//selfValidator := c.stakers.GetValidatorAddr(c.self)
		//if selfValidator != emptyAddrss && selfValidator != validator {
		//	// add my own amount
		//	//coe, err = c.miner.GetWorker().getValidatorCoefficient(c.self)
		//	//if err != nil {
		//	//	return err
		//	//}
		//	//weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(c.self), big.NewInt(int64(coe)))
		//	weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(selfValidator), big.NewInt(int64(averageCoefficient)))
		//	weightBalance.Div(weightBalance, big.NewInt(10))
		//	ps.receiveValidatorsSum = new(big.Int).Add(ps.receiveValidatorsSum, weightBalance)
		//	ps.onlineValidator.Add(selfValidator)
		//	selfSignedMessage, err := c.GetSignedMessage(new(big.Int).Set(height))
		//	if err != nil {
		//		return err
		//	}
		//	ps.emptyBlockMessages = append(ps.emptyBlockMessages, selfSignedMessage)
		//	log.Info("Certify.GatherOtherPeerSignature", "self", selfValidator.Hex(),
		//		"balance", c.stakers.StakeBalance(selfValidator), "average coe", averageCoefficient, "weightBalance", weightBalance,
		//		"receiveValidatorsSum", ps.receiveValidatorsSum, "height", height.Uint64())
		//}

		c.proofStatePool.proofs[height.Uint64()] = ps
		c.signatureResultCh <- height
		//log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 1")
		return nil
	}

	curProofs := c.proofStatePool.proofs[height.Uint64()]
	if curProofs.onlineValidator.Has(validator) {
		return errors.New("GatherOtherPeerSignature: validator exist")
	}
	c.proofStatePool.proofs[height.Uint64()].onlineValidator.Add(validator)
	c.proofStatePool.proofs[height.Uint64()].emptyBlockMessages = append(c.proofStatePool.proofs[height.Uint64()].emptyBlockMessages, encQues)
	//coe, err = c.miner.GetWorker().getValidatorCoefficient(validator)
	//if err != nil {
	//	return err
	//}
	//weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(validator), big.NewInt(int64(coe)))
	weightBalance = new(big.Int).Mul(c.stakers.StakeBalance(validator), big.NewInt(int64(averageCoefficient)))
	weightBalance.Div(weightBalance, big.NewInt(10))
	c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum = new(big.Int).Add(c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, weightBalance)
	//	log.Info("Certify.GatherOtherPeerSignature", "validator", validator.Hex(),
	//		"balance", c.stakers.StakeBalance(validator), "average coe", averageCoefficient, "weightBalance", weightBalance,
	//		"receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "height", height.Uint64())
	//log.Info("Certify.GatherOtherPeerSignature", "receiveValidatorsSum", c.proofStatePool.proofs[height.Uint64()].receiveValidatorsSum, "heigh", height)
	c.signatureResultCh <- height
	//log.Info("Certify.GatherOtherPeerSignature <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 2")
	return nil
}
