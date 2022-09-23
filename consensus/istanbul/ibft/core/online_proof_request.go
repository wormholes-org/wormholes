package core

import (
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/log"
)

func (c *core) handleOnlineProofRequest(request *istanbul.OnlineProofRequest) error {
	log.Info("ibftConsensus: handleOnlineProofRequest", "no", c.currentView().Sequence, "round", c.currentView().Round, "state", c.state, "self", c.address.Hex())
	if err := c.checkOnlineProofRequestMsg(request); err != nil {
		if err == istanbulcommon.ErrInvalidMessage {
			log.Error("ibftConsensus: invalid online proof request", "no", c.currentView().Sequence, "round", c.currentView().Round, "state", c.state, "self", c.address.Hex())
			return err
		}
		log.Warn("ibftConsensus: unexpected online proof request", "err", err, "no", request.Proposal.Number(), "round", c.currentView().Round, "hash", request.Proposal.Hash(), "self", c.address.Hex())
		return err
	}

	c.current.pendingOnlineProofRequest = request
	if c.state == ibfttypes.StateAcceptOnlineProofRequest {
		log.Info("ibftConsensus: handleOnlineProofRequest sendOnlineProof", "no", request.Proposal.Number(), "round", c.currentView().Round, "hash", request.Proposal.Hash(), "self", c.address.Hex())
		c.sendOnlineProof(request)
	}
	return nil
}

func (c *core) checkOnlineProofRequestMsg(request *istanbul.OnlineProofRequest) error {
	if request == nil || request.Proposal == nil {
		return istanbulcommon.ErrInvalidMessage
	}

	if c := c.current.sequence.Cmp(request.Proposal.Number()); c > 0 {
		return istanbulcommon.ErrOldMessage
	} else if c < 0 {
		return istanbulcommon.ErrFutureMessage
	} else {
		return nil
	}
}

func (c *core) storeOnlineProofRequestMsg(request *istanbul.OnlineProofRequest) {
	logger := c.logger.New("state", c.state)

	logger.Trace("Store future online proof request", "number", request.Proposal.Number(), "hash", request.Proposal.Hash())

	c.pendindingOnlineProofRequestsMu.Lock()
	defer c.pendindingOnlineProofRequestsMu.Unlock()

	c.pendingOnlineProofRequests.Push(request, float32(-request.Proposal.Number().Int64()))
}

func (c *core) processPendingOnlineProofRequests() {
	log.Info("ibftConsensus: processPendingOnlineProofRequests start", "no", c.currentView().Sequence, "round", c.currentView().Round, "self", c.address.Hex())
	c.pendindingOnlineProofRequestsMu.Lock()
	defer c.pendindingOnlineProofRequestsMu.Unlock()

	for !(c.pendingOnlineProofRequests.Empty()) {
		m, prio := c.pendingOnlineProofRequests.Pop()
		r, ok := m.(*istanbul.OnlineProofRequest)
		if !ok {
			c.logger.Warn("Malformed online proof request, skip", "msg", m)
			continue
		}
		// Push back if it's a future message
		err := c.checkOnlineProofRequestMsg(r)
		if err != nil {
			if err == istanbulcommon.ErrFutureMessage {
				c.logger.Trace("Stop processing online proof request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())
				c.pendingOnlineProofRequests.Push(m, prio)
				break
			}
			c.logger.Trace("Skip the pending online proof request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash(), "err", err)
			continue
		}
		c.logger.Trace("Post pending online proof  request", "number", r.Proposal.Number(), "hash", r.Proposal.Hash())

		go c.sendEvent(istanbul.OnlineProofEvent{
			Proposal:   r.Proposal,
			RandomHash: r.RandomHash,
		})
	}
	log.Info("ibftConsensus: processPendingOnlineProofRequests end", "no", c.currentView().Sequence)
}
