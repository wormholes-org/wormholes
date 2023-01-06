package miner

import (
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	mu                sync.Mutex
	self              common.Address
	eth               Backend
	otherMessages     *lru.ARCCache // the cache of peer's messages
	selfMessages      *lru.ARCCache // the cache of self messages
	cacheMessage      *lru.ARCCache
	eventMux          *event.TypeMux
	events            *event.TypeMuxSubscription
	stakers           *types.ValidatorList // all validator
	signatureResultCh chan *big.Int
	miner             Handler // Apply some of the capabilities of the parent class
	lock              sync.Mutex
	messageLock       sync.Mutex
	//receiveValidatorsSum *big.Int
	//validators           []common.Address
	voteIndex        uint64
	validatorsHeight []string
	proofStatePool   *ProofStatePool // Currently highly collected validators that have sent online proofs
	//msgHeight        *big.Int
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
	otherMsgs, _ := lru.NewARC(remotePeers)
	selfMsgs, _ := lru.NewARC(storeMsgs)
	cacheMsgs, _ := lru.NewARC(storeMsgs)
	certify := &Certify{
		self:              self,
		eth:               eth,
		eventMux:          new(event.TypeMux),
		otherMessages:     otherMsgs,
		selfMessages:      selfMsgs,
		cacheMessage:      cacheMsgs,
		miner:             handler,
		signatureResultCh: make(chan *big.Int),
		//receiveValidatorsSum: big.NewInt(0),
		//validators:           make([]common.Address, 0),
		voteIndex:        0,
		validatorsHeight: make([]string, 0),
		proofStatePool:   NewProofStatePool(),
		//msgHeight:        new(big.Int),
	}
	return certify
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

func (c *Certify) broadcast(msg *types.EmptyMsg) error {
	payload, err := c.signMessage(msg)
	if err != nil {
		log.Error("signMessage err", err)
		return err
	}
	// Broadcast payload
	//if err = c.Gossip(c.stakers, SendSignMsg, payload); err != nil {
	//	return err
	//}
	if miner, ok := c.miner.(*Miner); ok {
		miner.broadcaster.BroadcastEmptyBlockMsg(payload)
	}

	// send to self
	go c.eventMux.Post(types.EmptyMessageEvent{
		Code:    SendSignMsg,
		Payload: payload,
	})
	return nil
}

// Gossip Broadcast message to all stakers
//func (c *Certify) Gossip(valSet *types.ValidatorList, code uint64, payload []byte) error {
//	hash := istanbul.RLPHash(payload)
//	c.selfMessages.Add(hash, true)
//
//	targets := make(map[common.Address]bool)
//	for _, val := range valSet.Validators {
//		if val.Address() != c.Address() {
//			targets[val.Address()] = true
//		}
//	}
//	var ps map[common.Address]Peer
//	if miner, ok := c.miner.(*Miner); ok {
//		ps = miner.broadcaster.FindPeerSet(targets)
//	}
//	log.Info("certify gossip worker msg", "len", len(ps), "code", code)
//	for addr, p := range ps {
//		ms, ok := c.otherMessages.Get(addr)
//		var m *lru.ARCCache
//		if ok {
//			m, _ = ms.(*lru.ARCCache)
//			if _, k := m.Get(hash); k {
//				// This peer had this event, skip it
//				continue
//			}
//		} else {
//			m, _ = lru.NewARC(remotePeers)
//		}
//
//		m.Add(hash, true)
//		c.otherMessages.Add(addr, m)
//		go p.SendWorkerMsg(WorkerMsg, payload)
//	}
//	return nil
//}

//func (c *Certify) BroadcastEmptyBlockMsg(msg []byte) {
//	var ps map[common.Address]Peer
//	if miner, ok := c.miner.(*Miner); ok {
//		ps = miner.broadcaster.FindPeerSet(nil)
//	}
//
//	for _, p := range ps {
//		p.WriteQueueEmptyBlockMsg(msg)
//	}
//}

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

//func (c *Certify) Address() common.Address {
//	return c.self
//}

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

		if c.miner.GetWorker().isEmpty {
			//			log.Info("certify handleMsg post", "hash", hash)
			go c.eventMux.Post(types.EmptyMessageEvent{
				Code:    SendSignMsg,
				Payload: data,
			})
		} else {
			c.messageLock.Unlock()
			mc, ok := c.cacheMessage.Get(addr)
			var ml *lru.ARCCache
			if ok {
				ml, _ = mc.(*lru.ARCCache)
			} else {
				ml, _ = lru.NewARC(remotePeers)
				c.cacheMessage.Add(addr, ml)
			}
			ml.Add(hash, data)
			c.messageLock.Unlock()
			//log.Info("certify handleMsg cache", "hash", hash, "cache len", c.cacheMessage.Len())
		}
	}
	return false, nil
}

func (c *Certify) PostCacheMessage() {
	if c.cacheMessage.Len() <= 0 {
		return
	}

	cacheList := make([]interface{}, 0)
	for _, addr := range c.cacheMessage.Keys() {
		if ms, ok := c.cacheMessage.Get(addr); ok {
			cacheList = append(cacheList, ms)
		}
	}

	for _, ms := range cacheList {
		m, _ := ms.(*lru.ARCCache)
		if m.Len() <= 0 {
			continue
		}

		for _, hash := range m.Keys() {
			data, oks := m.Get(hash)
			if oks {
				m.Remove(hash)
				//log.Info("azh|repost", "hash", hash, "data", data)
				go c.eventMux.Post(types.EmptyMsg{
					Code: WorkerMsg,
					Msg:  data.([]byte),
				})
			} else {
				continue
			}
		}
	}
}

func (c *Certify) decode(msg p2p.Msg) ([]byte, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errDecodeFailed
	}
	return data, RLPHash(data), nil
}

func (c *Certify) handleEvents() {
	time.Sleep(time.Second * 10)
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
				//log.Info("Certify handle events")
				msg := new(types.EmptyMsg)
				if err := msg.FromPayload(ev.Payload); err != nil {
					log.Error("Certify Failed to decode message from payload", "err", err)
					break
				}
				sender, err := msg.RecoverAddress(ev.Payload)
				if err != nil {
					log.Error("Certify.handleEvents", "RecoverAddress error", err)
					break
				}

				var signature *types.SignatureData
				err = msg.Decode(&signature)
				if err != nil {
					log.Error("Certify.handleEvents", "msg.Decode error", err)
					break
				}

				_, err = Encode(signature)
				if err != nil {
					log.Error("Failed to encode", "subject", err)
					break
				}

				//c.msgHeight = signature.Height
				//log.Info("Certify.handleEvents", "msg.Code", msg.Code, "SendSignMsg", SendSignMsg, "Height", signature.Height)

				log.Info("azh|handleEvents", "self", c.self, "sender", sender, "vote", signature.Vote, "height", signature.Height)
				if msg.Code == SendSignMsg {
					//log.Info("Certify.handleEvents", "SendSignMsg", SendSignMsg, "msg.Address", msg.Address.Hex(),
					//	"signature.Address", signature.Address, "signature.Height", signature.Height, "signature.Timestamp", signature.Timestamp,
					//	"c.stakers number", len(c.stakers.Validators))
					//If the GatherOtherPeerSignature is ok, gossip message directly
					if err := c.GatherOtherPeerSignature(sender, signature.Vote, signature.Height, ev.Payload); err == nil {
						c.rebroadcast(c.self, ev.Payload)
					}
				}
			}
		}
	}
}
