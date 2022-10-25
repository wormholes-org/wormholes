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
	log.Info("nixMsg")
	if request.Version != "" {
		log.Info("msgVers", c.address.Hex(), request.Version)
	}
	curView := c.currentView()
	onlineProof := &istanbul.OnlineProof{
		View:       curView,
		Proposal:   request.Proposal,
		RandomHash: request.RandomHash,
		Version:    request.Version,
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

func (c *core) handleOnlineProof(msg *ibfttypes.Message, src istanbul.Validator) error {
	c.onlineProofsMu.Lock()
	defer c.onlineProofsMu.Unlock()
	var onlineProof *istanbul.OnlineProof
	err := msg.Decode(&onlineProof)
	if err != nil {
		return istanbulcommon.ErrFailedDecodeOnlineProof
	}

	c.acceptOnlineProof(msg, src)

	onlineValidators := c.onlineProofs[c.current.sequence.Uint64()]
	if onlineValidators == nil || (onlineValidators != nil && !onlineValidators.ExistAddress(src.Address())) {
		validator := types.NewOnlineValidator(c.current.sequence, src.Address(), onlineProof.RandomHash, onlineProof.Signature)
		onlineValidators.Validators = append(onlineValidators.Validators, validator)
	} else {
		return errors.New("This address is already online")
	}

	//c.onlineProofs[c.current.sequence.Uint64()] = tempMessageSet
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
