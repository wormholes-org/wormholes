package backend

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"testing"
)

type node struct {
	server *p2p.Server
	peers  *peerSet
	key    *ecdsa.PrivateKey
	msgCh  chan p2p.Msg
	peerCh chan *p2p.Peer
}

func newNode(key *ecdsa.PrivateKey, listenAddr string, bootNode []*enode.Node) (*node, error) {
	peers := newPeerSet()
	msgCh := make(chan p2p.Msg, 10)
	peerCh := make(chan *p2p.Peer, 10)
	server := &p2p.Server{}
	server.MaxPeers = 160
	server.DialRatio = 2
	server.PrivateKey = key
	server.ListenAddr = listenAddr
	server.BootstrapNodes = bootNode

	server.Protocols = []p2p.Protocol{
		{
			Version: 100,
			Length:  22,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				if server.Self().ID() == p.ID() {
					return errAddSelf
				}
				addr := crypto.PubkeyToAddress(*p.Node().Pubkey())
				if err := peers.registerPeer(&peer{p, rw}); err != nil {
					return err
				}
				defer peers.unregisterPeer(addr)
				peerCh <- p
				for {
					msg, err := rw.ReadMsg()
					if err != nil {
						return err
					}
					msgCh <- msg
				}
			},
			NodeInfo: func() interface{} {
				return struct{}{}
			},
			PeerInfo: func(id enode.ID) interface{} {
				return struct{}{}
			},
		},
	}
	err := server.Start()
	if err != nil {
		return nil, err
	}
	return &node{
		server: server,
		peers:  peers,
		key:    key,
		msgCh:  msgCh,
		peerCh: peerCh,
	}, nil
}

func TestNetWork(t *testing.T) {
	bootKey, err := crypto.GenerateKey()
	if err != nil {
		t.Error()
		return
	}
	key1, err := crypto.GenerateKey()
	if err != nil {
		t.Error()
		return
	}
	key2, err := crypto.GenerateKey()
	if err != nil {
		t.Error()
		return
	}
	bootNode, err := newNode(bootKey, "127.0.0.1:62220", nil)
	if err != nil {
		t.Error()
		return
	}
	n1, err := newNode(key1, "127.0.0.1:62221", []*enode.Node{bootNode.server.Self()})
	if err != nil {
		t.Error()
		return
	}
	n2, err := newNode(key2, "127.0.0.1:62222", []*enode.Node{bootNode.server.Self()})
	if err != nil {
		t.Error()
		return
	}

	// add boot node
	<-n1.peerCh
	<-n2.peerCh
	// add other node
	<-n1.peerCh
	<-n2.peerCh

	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	if p := n1.peers.peer(addr2); p == nil {
		t.Error("not add peer n2")
	} else {
		p2p.Send(p.rw, 1, "hello")
	}
	if p := n2.peers.peer(addr1); p == nil {
		t.Error("not add peer n1")
	} else {
		p2p.Send(p.rw, 1, "hello")
	}
	msg1 := <-n2.msgCh
	msg2 := <-n1.msgCh
	t.Log(msg1.Code, msg1.String())
	t.Log(msg2.Code, msg2.String())
	if msg1.Code != 1 {
		t.Error("send msg code not equal")
	}
	if msg2.Code != 1 {
		t.Error("send msg code not equal")
	}
}
