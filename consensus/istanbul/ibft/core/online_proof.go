package core

import (
	"errors"

	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"go.uber.org/zap/buffer"
)

func (c *core) sendOnlineProof(request *istanbul.OnlineProofRequest) {
	if request == nil {
		return
	}
	log.Info("ibftConsensus: sendOnlineProof",
		"no", request.Proposal.Number(),
		"sequence", c.current.Sequence(),
		"round", c.current.round,
		"self", c.address.Hex(),
	)

	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 {
		curView := c.currentView()
		onlineProof := &istanbul.OnlineProof{
			View:       curView,
			Proposal:   request.Proposal,
			RandomHash: request.RandomHash,
		}

		// sign data total byte: 8 + 32
		buffer := new(buffer.Buffer)
		buffer.AppendUint((curView.Sequence.Uint64()))
		for _, B := range request.RandomHash.Bytes() {
			buffer.AppendByte(B)
		}

		signature, err := c.backend.Sign((buffer.Bytes()))
		if err != nil {
			return
		}

		onlineProof.Signature = signature

		onlineProofEnc, err := ibfttypes.Encode(&onlineProof)

		if err != nil {
			log.Error("Failed to encode", "view", curView)
			return
		}

		log.Info("ibftConsensus: sendOnlineProof broadcast",
			"no", request.Proposal.Number(),
			"round", curView.Round,
			"self", c.address.Hex(),
			"hash", request.Proposal.Hash().Hex())

		c.broadcast(&ibfttypes.Message{
			Code: ibfttypes.MsgOnlineProof,
			Msg:  onlineProofEnc,
		})
	}
}

func (c *core) handleOnlineProof(msg *ibfttypes.Message, src istanbul.Validator) error {
	var onlineProof *istanbul.OnlineProof
	err := msg.Decode(&onlineProof)
	if err != nil {
		return istanbulcommon.ErrFailedDecodeOnlineProof
	}

	log.Info("ibftConsensus: handleOnlineProof",
		"no", onlineProof.Proposal.Number().Uint64(),
		"round", onlineProof.View.Round.String(),
		"from", src.Address().Hex(),
		"hash", onlineProof.Proposal.Hash().Hex(),
		"isproposer", c.IsProposer(),
		"self", c.address.Hex())
	// Ensure we have the same view with the ONLINE-PROOF message
	if err := c.checkMessage(ibfttypes.MsgOnlineProof, onlineProof.View); err != nil {
		log.Error("ibftConsensus: handleOnlineProof checkMessage", "no", onlineProof.Proposal.Number().Uint64(), "round", c.currentView().Round, "self", c.address.Hex(), "err", err.Error())
		return err
	}

	//TODOverify online proof
	currentProofs := c.onlineProofs[c.current.sequence.Uint64()]
	exist := false
	for _, msg := range currentProofs.messages {
		if msg.Address == src.Address() && onlineProof.Proposal.Number().Uint64() <= c.current.sequence.Uint64() {
			exist = true
			log.Warn("handleOnlineProof : This altitude has been certified online", "no", c.current.sequence.Uint64(), "src", src.Address())
			break
		}
	}
	if !exist {
		// Here is about to accept the ONLINE-PROOF msg
		c.acceptOnlineProof(msg, src)
	}

	// clear onlineProof memory
	for k, _ := range c.onlineProofs {
		if k < c.currentView().Sequence.Uint64() {
			log.Info("handleOnlineProof: clear onlineProof", "height", k, "currentView.no", c.currentView().Sequence)
			delete(c.onlineProofs, k)
		}
	}

	if len(c.onlineProofs) == 0 {
		log.Error("handleOnlineProof: len(onlineProof)==0", "height", c.currentView().Sequence, "round", c.currentView().Round)
		return errors.New("handleOnlineProof: len(onlineProof)==0")
	}

	log.Info("ibftConsensus: prepare to notify worker to commit", "no", c.currentView().Sequence,
		"round", c.currentView().Round, "len(OnlineProofs)", c.current.OnlineProofs.Size(),
		"self", c.address.Hex(), "state", c.state)
	if c.current.OnlineProofs.Size() >= c.QuorumSize() && c.state == ibfttypes.StateAcceptOnlineProofRequest { // Submit the collected online attestation data to the worker module
		onlineValidatorList := new(types.OnlineValidatorList)
		for _, v := range c.current.OnlineProofs.messages {
			err := v.Decode(&onlineProof)
			if err != nil {
				continue
			}
			validator := types.NewOnlineValidator(
				onlineProof.View.Sequence,
				v.Address,
				onlineProof.RandomHash,
				onlineProof.Signature)
			onlineValidatorList.Validators = append(onlineValidatorList.Validators, validator)
		}
		onlineValidatorList.Validators = append(onlineValidatorList.Validators)
		onlineValidatorList.Height = onlineProof.Proposal.Number()
		// Notify miners to submit blocks
		c.backend.NotifyWorkerToCommit(onlineValidatorList)

		log.Info("ibftConsensus: collected online proof",
			"no", c.currentView().Sequence,
			"round", c.currentView().Round,
			"size", c.current.OnlineProofs.Size(),
			"onlineproofs len", c.onlineProofs[c.current.sequence.Uint64()].Size(),
			"self", c.address.Hex())
		// Set state to StateAcceptRequest
		c.setState(ibfttypes.StateAcceptRequest)
	}
	// Stored in the temporary state of the engine, no data will be lost with the
	tempMessageSet := newMessageSet(c.valSet)
	for _, v := range c.current.OnlineProofs.Values() {
		tempMessageSet.Add(v)
	}
	c.onlineProofs[c.current.sequence.Uint64()] = tempMessageSet
	log.Info("ibftConsensus: onlineProofs",
		"no", c.currentView().Sequence,
		"round", c.currentView().Round,
		"size", c.current.OnlineProofs.Size(),
		"onlineproofs len", c.onlineProofs[c.current.sequence.Uint64()].Size(),
		"self", c.address.Hex(),
		"state", c.state)

	return nil
}

func (c *core) acceptOnlineProof(msg *ibfttypes.Message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the ONLINE-PROOF  message to current round state
	if err := c.current.OnlineProofs.Add(msg); err != nil {
		logger.Error("Failed to add ONLINE-PROOF message to round state",
			"no", c.currentView().Sequence,
			"round", c.currentView().Round, "msg", msg, "err", err)
		return err
	}

	return nil
}
