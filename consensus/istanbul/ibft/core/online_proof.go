package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) sendOnlineProof(request *istanbul.OnlineProofRequest) {
	if request == nil {
		return
	}
	log.Info("sendOnlineProof", "height", request.Proposal.Number())
	if c.current.Sequence().Cmp(request.Proposal.Number()) == 0 {
		curView := c.currentView()
		onlineProof, err := ibfttypes.Encode(&istanbul.OnlineProof{
			View:       curView,
			Proposal:   request.Proposal,
			RandomHash: request.RandomHash,
		})
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
			Msg:  onlineProof,
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
		"hash", onlineProof.Proposal.Hash().Hex())
	// Ensure we have the same view with the ONLINE-PROOF message
	if err := c.checkMessage(ibfttypes.MsgOnlineProof, onlineProof.View); err != nil {
		log.Error("handleOnlineProof : checkMessage", "no", onlineProof.Proposal.Number().Uint64(), "err", err.Error(), "hash", onlineProof.Proposal.Hash().Hex())
		return err
	}

	//TODO  verify online proof

	// Here is about to accept the ONLINE-PROOF msg
	c.acceptOnlineProof(msg, src)

	if c.current.OnlineProofs.Size() >= c.QuorumSize() && c.state == ibfttypes.StateAcceptOnlineProofRequest {
		// Submit the collected online attestation data to the worker module
		var (
			addrs       []common.Address
			hashs       []common.Hash
			onlineProof *istanbul.OnlineProof
		)
		for _, v := range c.current.OnlineProofs.messages {
			err := v.Decode(&onlineProof)
			if err != nil {
				continue
			}
			addrs = append(addrs, v.Address)
			hashs = append(hashs, onlineProof.RandomHash)
		}

		onlineValidators := &types.OnlineValidatorInfo{
			Height: c.current.sequence,
			Addrs:  addrs,
			Hashs:  hashs,
		}

		// Notify miners to submit blocks
		c.backend.NotifyWorkerToCommit(onlineValidators)

		// Set state to StateAcceptRequest
		c.setState(ibfttypes.StateAcceptRequest)
	}

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
