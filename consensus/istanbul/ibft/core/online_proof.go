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
	log.Info("sendOnlineProof",
		"no", request.Proposal.Number(),
		"sequence", c.current.Sequence(),
		"round", c.current.round,
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

		log.Info("sendOnlineProof : broadcast",
			"no", request.Proposal.Number(),
			"round", curView.Round,
			"author", c.address.Hex(),
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

	log.Info("handleOnlineProof",
		"no", onlineProof.Proposal.Number().Uint64(),
		"round", onlineProof.View.Round.String(),
		"from", src.Address().Hex(),
		"hash", onlineProof.Proposal.Hash().Hex(),
		"isproposer", c.IsProposer())
	// Ensure we have the same view with the ONLINE-PROOF message
	if err := c.checkMessage(ibfttypes.MsgOnlineProof, onlineProof.View); err != nil {
		log.Error("handleOnlineProof : checkMessage", "no", onlineProof.Proposal.Number().Uint64(), "err", err.Error(), "hash", onlineProof.Proposal.Hash().Hex())
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
		log.Info("handleOnlineProof : prepare to notify worker to commit")
		// Notify miners to submit blocks
		c.backend.NotifyWorkerToCommit(onlineValidatorList)

		// Set state to StateAcceptRequest
		c.setState(ibfttypes.StateAcceptRequest)
	}
	// Stored in the temporary state of the engine, no data will be lost with the
	tempMessageSet := newMessageSet(c.valSet)
	for _, v := range c.current.OnlineProofs.Values() {
		tempMessageSet.Add(v)
	}
	c.onlineProofs[c.current.sequence.Uint64()] = tempMessageSet
	log.Info("handleOnlineProof : QuorumSize", "height", onlineProof.Proposal.Number().Uint64(), "size", c.current.OnlineProofs.Size(), "onlineproofs len", c.onlineProofs[c.current.sequence.Uint64()].Size())

	return nil
}

func (c *core) acceptOnlineProof(msg *ibfttypes.Message, src istanbul.Validator) error {
	logger := c.logger.New("from", src, "state", c.state)

	// Add the ONLINE-PROOF  message to current round state
	if err := c.current.OnlineProofs.Add(msg); err != nil {
		logger.Error("Failed to add ONLINE-PROOF message to round state", "msg", msg, "err", err)
		return err
	}

	return nil
}
