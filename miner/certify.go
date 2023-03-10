package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	lru "github.com/hashicorp/golang-lru"
	"golang.org/x/xerrors"
	"math/big"
	"sync"
)

const (
	//remotePeers = 2000 // Number of messages kept in consensus workers per round (11 * 2)
	storeMsgs = 2500 // Number of messages stored by yourself
)

type Certify struct {
	mu   sync.Mutex
	self common.Address
	eth  Backend
	//otherMessages     *lru.ARCCache // the cache of peer's messages
	selfMessages      *lru.ARCCache // the cache of self messages
	eventMux          *event.TypeMux
	events            *event.TypeMuxSubscription
	stakers           *types.ValidatorList // all validator
	signatureResultCh chan VoteResult
	miner             Handler // Apply some of the capabilities of the parent class
	lock              sync.Mutex
	round             uint64
	voteIndex         int
	validatorsHeight  []string
	proofStatePool    *ProofStatePool // Currently highly collected validators that have sent online proofs
}

func (c *Certify) Start() {
	c.subscribeEvents()
}

func (c *Certify) subscribeEvents() {
	c.events = c.eventMux.Subscribe(
		types.EmptyMessageEvent{},
	)
}

func NewCertify(self common.Address, eth Backend, handler Handler) *Certify {
	//otherMsgs, _ := lru.NewARC(remotePeers)
	selfMsgs, _ := lru.NewARC(storeMsgs)
	certify := &Certify{
		self:              self,
		eth:               eth,
		eventMux:          new(event.TypeMux),
		selfMessages:      selfMsgs,
		miner:             handler,
		signatureResultCh: make(chan VoteResult),
		voteIndex:         0,
		validatorsHeight:  make([]string, 0),
		proofStatePool:    NewProofStatePool(),
	}
	return certify
}

type VoteResult struct {
	Height           *big.Int
	ReceiveSum       *big.Int
	OnlineValidators []common.Address
	EmptyMessages    [][]byte
}

func (c *Certify) rebroadcast(from common.Address, payload []byte) error {
	// Broadcast payload
	//if err := c.Gossip(c.stakers, SendSignMsg, payload); err != nil {
	//	return err
	//}
	if miner, ok := c.miner.(*Miner); ok {
		miner.broadcaster.BroadcastEmptyBlockMsg(payload)
	}

	return nil
}

func (c *Certify) signMessage(msg *types.EmptyMsg) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.self

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *Certify) sign(data []byte) ([]byte, error) {
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, c.eth.GetNodeKey())
}

func (c *Certify) assembleMessage(height *big.Int, vote common.Address) (error, []byte) {
	ques := &types.SignatureData{
		Vote:   vote,
		Height: height,
		Round:  c.round,
	}
	encQues, err := Encode(ques)
	if err != nil {
		//log.Error("Failed to encode", "subject", err)
		return err, nil
	}

	msg := &types.EmptyMsg{
		Code: SendSignMsg,
		Msg:  encQues,
	}

	payload, err := c.signMessage(msg)
	if err != nil {
		//log.Error("signMessage err", err)
		return err, nil
	}
	return nil, payload
}

//func (c *Certify) Address() common.Address {
//	return c.self
//}

// HandleMsg handles a message from peer
func (c *Certify) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code == WorkerMsg {
		data, hash, err := c.decode(msg)
		//log.Info("certify handleMsg", "code", msg.Code, "payload", data)
		if err != nil {
			return true, err
		}

		msg := new(types.EmptyMsg)
		if err := msg.FromPayload(data); err != nil {
			log.Error("Certify Failed to decode message from payload", "err", err)
			return true, err
		}

		var signature *types.SignatureData
		err = msg.Decode(&signature)
		if err != nil {
			log.Error("Certify.handleEvents", "msg.Decode error", err)
			return true, err
		}

		currentHeight := c.miner.GetWorker().chain.CurrentHeader().Number
		//deal only current or more than two block message
		if currentHeight.Cmp(new(big.Int).Sub(signature.Height, big.NewInt(1))) > 0 ||
			currentHeight.Cmp(new(big.Int).Sub(signature.Height, big.NewInt(3))) < 0{
			return true, nil
		}

		sender, err := msg.RecoverAddress(data)
		if err != nil {
			log.Error("Certify.handleEvents", "RecoverAddress error", err)
			return true, err
		}

		if c.round != signature.Round{
			return true, nil
		}

		log.Info("azh|emptyMessage", "height", signature.Height, "from", sender, "vote", signature.Vote, "round", signature.Round)

		c.rebroadcast(addr, data)

		if currentHeight.Cmp(new(big.Int).Sub(signature.Height, big.NewInt(1))) < 0 {
			return true, nil
		}

		if c.stakers == nil {
			return true, nil
		}

		if c.stakers.GetValidatorAddr(sender) == (common.Address{}) {
			if addr == sender {
				return true, xerrors.New("Certify.handleEvents the vote is not a miner")
			}
			return true, nil
		}

		if c.self == signature.Vote {
			ms, ok := c.selfMessages.Get(sender)
			var m *lru.ARCCache
			if ok {
				m, _ = ms.(*lru.ARCCache)
				if _, ok := m.Get(hash); ok {
					return true, nil
				}
			} else {
				m, _ = lru.NewARC(storeMsgs)
				c.selfMessages.Add(sender, m)
			}
			m.Add(hash, true)

			emptyMsg := types.EmptyMessageEvent{
				Sender:  sender,
				Height:  signature.Height,
				Payload: data,
			}
			go c.eventMux.Post(emptyMsg)
			return true, nil
		}
	}
	return false, nil
}

func (c *Certify) decode(msg p2p.Msg) ([]byte, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errDecodeFailed
	}
	return data, RLPHash(data), nil
}

func (c *Certify) handleEvents() {
	log.Info("Certify handle events start")
	for {
		select {
		case event, ok := <-c.events.Chan():
			if !ok {
				continue
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case types.EmptyMessageEvent:
				log.Info("handleEvents", "sender", ev.Sender, "height", ev.Height)
				c.GatherOtherPeerSignature(ev.Sender, ev.Height, ev.Payload)
			}
		}
	}
}
