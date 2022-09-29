package miner

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	lru "github.com/hashicorp/golang-lru"
)

const (
	remotePeers = 2000 // Number of messages kept in consensus workers per round (11 * 2)
	storeMsgs   = 2500 // Number of messages stored by yourself
)

type Certify struct {
	mu                   sync.Mutex
	self                 common.Address
	eth                  Backend
	otherMessages        *lru.ARCCache // the cache of peer's messages
	selfMessages         *lru.ARCCache // the cache of self messages
	eventMux             *event.TypeMux
	events               *event.TypeMuxSubscription
	stakers              *types.ValidatorList // all validator
	signatureResultCh    chan *big.Int
	miner                Handler // Apply some of the capabilities of the parent class
	lock                 sync.Mutex
	receiveValidatorsSum *big.Int
	validators           []common.Address
	validatorsHeight     []string
	proofStatePool       *ProofStatePool // Currently highly collected validators that have sent online proofs
	msgHeight            *big.Int
}

func (c *Certify) Start() {
	c.subscribeEvents()
}
func (c *Certify) subscribeEvents() {
	c.events = c.eventMux.Subscribe(
		MessageEvent{},
	)
}

func NewCertify(self common.Address, eth Backend, handler Handler) *Certify {
	otherMsgs, _ := lru.NewARC(remotePeers)
	selfMsgs, _ := lru.NewARC(storeMsgs)
	certify := &Certify{
		self:                 self,
		eth:                  eth,
		eventMux:             new(event.TypeMux),
		otherMessages:        otherMsgs,
		selfMessages:         selfMsgs,
		miner:                handler,
		signatureResultCh:    make(chan *big.Int),
		receiveValidatorsSum: big.NewInt(0),
		validators:           make([]common.Address, 0),
		validatorsHeight:     make([]string, 0),
		proofStatePool:       NewProofStatePool(),
		msgHeight:            new(big.Int),
	}
	return certify
}

func (c *Certify) rebroadcast(from common.Address, payload []byte) error {
	// Broadcast payload
	if err := c.Gossip(c.stakers, SendSignMsg, payload); err != nil {
		return err
	}
	return nil
}

func (c *Certify) broadcast(from common.Address, msg *Msg) error {
	payload, err := c.signMessage(from, msg)
	if err != nil {
		log.Error("signMessage err", err)
		return err
	}
	// Broadcast payload
	if err = c.Gossip(c.stakers, SendSignMsg, payload); err != nil {
		return err
	}
	// send to self
	go c.eventMux.Post(msg)
	return nil
}

// Gossip Broadcast message to all stakers
func (c *Certify) Gossip(valSet *types.ValidatorList, code uint64, payload []byte) error {
	hash := istanbul.RLPHash(payload)
	c.selfMessages.Add(hash, true)

	targets := make(map[common.Address]bool)
	for _, val := range valSet.Validators {
		if val.Address() != c.Address() {
			targets[val.Address()] = true
		}
	}
	var ps map[common.Address]Peer
	if miner, ok := c.miner.(*Miner); ok {
		ps = miner.broadcaster.FindPeerSet(targets)
	}
	log.Info("certify gossip worker msg", "len", len(ps), "code", code)
	for addr, p := range ps {
		ms, ok := c.otherMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
			if _, k := m.Get(hash); k {
				// This peer had this event, skip it
				continue
			}
		} else {
			m, _ = lru.NewARC(remotePeers)
		}

		m.Add(hash, true)
		c.otherMessages.Add(addr, m)
		go p.SendWorkerMsg(WorkerMsg, payload)
	}
	return nil
}

func (c *Certify) signMessage(coinbase common.Address, msg *Msg) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = coinbase

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

func (c *Certify) Address() common.Address {
	return c.self
}

// HandleMsg handles a message from peer
func (c *Certify) HandleMsg(addr common.Address, msg p2p.Msg) (bool, error) {
	if msg.Code == WorkerMsg {
		data, hash, err := c.decode(msg)
		log.Info("certify handleMsg", "code", msg.Code, "payload", data)
		if err != nil {
			return true, err
		}
		// Mark peer's message
		ms, ok := c.otherMessages.Get(addr)
		var m *lru.ARCCache
		if ok {
			m, _ = ms.(*lru.ARCCache)
		} else {
			m, _ = lru.NewARC(remotePeers)
			c.otherMessages.Add(addr, m)
		}
		m.Add(hash, true)

		// Mark self known message
		if _, ok := c.selfMessages.Get(hash); ok {
			return true, nil
		}
		c.selfMessages.Add(hash, true)

		go c.eventMux.Post(MessageEvent{
			Code:    msg.Code,
			Payload: data,
		})
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
				return
			}
			// A real event arrived, process interesting content
			switch ev := event.Data.(type) {
			case MessageEvent:
				log.Info("Certify handle events")
				msg := new(Msg)
				if err := msg.FromPayload(ev.Payload); err != nil {
					log.Error("Certify Failed to decode message from payload", "err", err)
				}
				var signature *SignatureData
				msg.Decode(&signature)

				encQues, err := Encode(signature)
				if err != nil {
					log.Error("Failed to encode", "subject", err)
					return
				}

				c.msgHeight = signature.Height
				log.Info("signature", "Height", signature.Height)
				//if len(c.stakers.Validators) > 0 && c.stakers != nil {
				//	flag := c.stakers.GetByAddress(msg.Address)
				//	if flag == -1 {
				//		log.Error("Invalid address in message", "msg", msg)
				//		return
				//	}
				//}
				if msg.Code == SendSignMsg {
					log.Info("SendSignMsg", "SendSignMsg", c.stakers)
					//If the GatherOtherPeerSignature is ok, gossip message directly
					if err := c.GatherOtherPeerSignature(msg.Address, signature.Height, encQues); err == nil {
						c.rebroadcast(c.Address(), ev.Payload)})
					}
				}
			}
		}
	}
}
