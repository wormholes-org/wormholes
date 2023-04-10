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
	storeMsgs = 32768 // Number of messages stored by yourself
)

type Certify struct {
	mu                sync.Mutex
	self              common.Address
	eth               Backend
	recentMessages    map[string]*EmptyPeerInfo // the cache of peer's messages
	selfMessages      *lru.ARCCache             // the cache of self messages
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
	requestEmpty      chan []byte
	status            int
	purge             chan struct{}
}

type VoteResult struct {
	Height           *big.Int
	ReceiveSum       *big.Int
	OnlineValidators []common.Address
	EmptyMessages    [][]byte
}

type EmptyPeerInfo struct {
	Peer    Peer
	Message [][]byte
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
	selfMsgs, _ := lru.NewARC(storeMsgs)
	certify := &Certify{
		self:              self,
		eth:               eth,
		eventMux:          new(event.TypeMux),
		selfMessages:      selfMsgs,
		recentMessages:    make(map[string]*EmptyPeerInfo),
		miner:             handler,
		signatureResultCh: make(chan VoteResult),
		voteIndex:         0,
		validatorsHeight:  make([]string, 0),
		proofStatePool:    NewProofStatePool(),
		requestEmpty:      make(chan []byte, 1000),
		purge:             make(chan struct{}, 1),
	}
	return certify
}

func (c *Certify) RequestEmptyMessage() {
	miner, ok := c.miner.(*Miner)
	if !ok {
		return
	}

	emptyResponse := miner.broadcaster.EmptyResponse()

	var haveDone bool
	mixCh, maxCh := 5, 15
	delPeer := make(map[string]struct{})

	broad := func() {
		if len(c.recentMessages)*2/3 < maxCh {
			maxCh = len(c.recentMessages) * 2 / 3
		}

		if len(c.recentMessages)/3 < mixCh {
			mixCh = len(c.recentMessages) / 3
		}

		peerStatus := miner.broadcaster.PeerStatus()
		if c.status == 0 {
			if len(peerStatus) > maxCh {
				c.status = 1
			}
		} else {
			if len(peerStatus) < mixCh {
				c.status = 0
			}
		}

		for id, info := range c.recentMessages {
			if haveDone {
				peerStatus = miner.broadcaster.PeerStatus()
				if c.status == 0 {
					if len(peerStatus) > maxCh {
						c.status = 1
					}
				} else {
					if len(peerStatus) < mixCh {
						c.status = 0
					}
				}
			}

			haveDone = false

			_, ok := peerStatus[id]
			if ok {
				continue
			}

			//log.Info("azh|cache message", "max", maxCh, "mix", mixCh, "count", count, "status", c.status)
			index := 0
			if c.status == 0 {
				for _, msg := range info.Message {
					if err := info.Peer.RequestEmptyMsg(msg); err != nil {
						break
					} else {
						index++
					}
				}

				if index < len(info.Message)-1 {
					info.Message = info.Message[index:]
				} else {
					delPeer[id] = struct{}{}
				}

				haveDone = true
			} else {
				break
			}
		}

		if len(delPeer) > 0 {
			for id, _ := range delPeer {
				delete(c.recentMessages, id)
			}
		}
	}

	for {

		select {
		case msg := <-c.requestEmpty:
			ps := miner.broadcaster.FindPeerSet()

			for id, p := range ps {
				if info, ok := c.recentMessages[id]; ok {
					info.Message = append(info.Message, msg)
				} else {
					messages := make([][]byte, 0)
					messages = append(messages, msg)
					infos := &EmptyPeerInfo{
						Peer:    p,
						Message: messages,
					}
					c.recentMessages[id] = infos
				}
			}

			broad()

		case id := <-emptyResponse:
			if _, ok := c.recentMessages[id]; ok {
				delete(c.recentMessages, id)
			}

			if _, ok := delPeer[id]; ok {
				delete(delPeer, id)
			}

		case <-c.purge:
			c.recentMessages = make(map[string]*EmptyPeerInfo)
			delPeer = make(map[string]struct{})
		}
	}
}

//func (c *Certify) rebroadcast(height uint64, payload []byte) error {
//	// Broadcast payload
//	//if err := c.Gossip(c.stakers, SendSignMsg, payload); err != nil {
//	//	return err
//	//}
//	if miner, ok := c.miner.(*Miner); ok {
//		miner.broadcaster.BroadcastEmptyBlockMsg(payload)
//	}
//
//	return nil
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
		if currentHeight.Cmp(new(big.Int).Sub(signature.Height, big.NewInt(1))) != 0 {
			//return true, errors.New("GatherOtherPeerSignature: msg height < chain Number")
			return true, nil
		}

		sender, err := msg.RecoverAddress(data)
		if err != nil {
			log.Error("Certify.handleEvents", "RecoverAddress error", err)
			return true, err
		}

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

		r, _ := c.selfMessages.Get(sender)
		_, ok1 := r.(*lru.ARCCache).Get(hash)
		log.Info("azh|handle empty", "add", ok1, "msg hash", hash)

		if c.stakers == nil {
			return true, nil
		}

		if c.stakers.GetValidatorAddr(sender) == (common.Address{}) {
			return true, xerrors.New("Certify.handleEvents the vote is not a miner")
		}

		log.Info("azh|emptyMessage", "height", signature.Height, "from", sender, "vote", signature.Vote, "round", signature.Round)
		if c.self == signature.Vote {
			emptyMsg := types.EmptyMessageEvent{
				Sender:  sender,
				Height:  signature.Height,
				Payload: data,
			}
			go c.eventMux.Post(emptyMsg)
			return true, nil
		} else {
			c.requestEmpty <- data
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
