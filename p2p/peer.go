// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	ErrShuttingDown = errors.New("shutting down")
)

const (
	baseProtocolVersion    = 5
	baseProtocolLength     = uint64(16)
	baseProtocolMaxMsgSize = 2 * 1024

	snappyProtocolVersion = 5

	pingInterval = 15 * time.Second
)

const (
	// devp2p message codes
	handshakeMsg = 0x00
	discMsg      = 0x01
	pingMsg      = 0x02
	pongMsg      = 0x03
)

// protoHandshake is the RLP structure of the protocol handshake.
type protoHandshake struct {
	Version    uint64
	Name       string
	Caps       []Cap
	ListenPort uint64
	ID         []byte // secp256k1 public key

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// PeerEventType is the type of peer events emitted by a p2p.Server
type PeerEventType string

const (
	// PeerEventTypeAdd is the type of event emitted when a peer is added
	// to a p2p.Server
	PeerEventTypeAdd PeerEventType = "add"

	// PeerEventTypeDrop is the type of event emitted when a peer is
	// dropped from a p2p.Server
	PeerEventTypeDrop PeerEventType = "drop"

	// PeerEventTypeMsgSend is the type of event emitted when a
	// message is successfully sent to a peer
	PeerEventTypeMsgSend PeerEventType = "msgsend"

	// PeerEventTypeMsgRecv is the type of event emitted when a
	// message is received from a peer
	PeerEventTypeMsgRecv PeerEventType = "msgrecv"
)

// PeerEvent is an event emitted when peers are either added or dropped from
// a p2p.Server or when a message is sent or received on a peer connection
type PeerEvent struct {
	Type          PeerEventType `json:"type"`
	Peer          enode.ID      `json:"peer"`
	Error         string        `json:"error,omitempty"`
	Protocol      string        `json:"protocol,omitempty"`
	MsgCode       *uint64       `json:"msg_code,omitempty"`
	MsgSize       *uint32       `json:"msg_size,omitempty"`
	LocalAddress  string        `json:"local,omitempty"`
	RemoteAddress string        `json:"remote,omitempty"`
}

// Peer represents a connected remote node.
type Peer struct {
	rw      *conn
	running map[string]*protoRW
	log     log.Logger
	created mclock.AbsTime

	wg       sync.WaitGroup
	protoErr chan error
	closed   chan struct{}
	disc     chan DiscReason

	// events receives message send / receive events if set
	events   *event.Feed
	testPipe *MsgPipeRW // for testing

	// Quorum
	EthPeerRegistered   chan struct{}
	EthPeerDisconnected chan struct{}

	Msgs *MsgList
}

// NewPeer returns a peer for testing purposes.
func NewPeer(id enode.ID, name string, caps []Cap) *Peer {
	pipe, _ := net.Pipe()
	node := enode.SignNull(new(enr.Record), id)
	conn := &conn{fd: pipe, transport: nil, node: node, caps: caps, name: name}
	peer := newPeer(log.Root(), conn, nil)
	close(peer.closed) // ensures Disconnect doesn't block
	return peer
}

// NewPeerPipe creates a peer for testing purposes.
// The message pipe given as the last parameter is closed when
// Disconnect is called on the peer.
func NewPeerPipe(id enode.ID, name string, caps []Cap, pipe *MsgPipeRW) *Peer {
	p := NewPeer(id, name, caps)
	p.testPipe = pipe
	return p
}

// ID returns the node's public key.
func (p *Peer) ID() enode.ID {
	return p.rw.node.ID()
}

// Node returns the peer's node descriptor.
func (p *Peer) Node() *enode.Node {
	return p.rw.node
}

// Name returns an abbreviated form of the name
func (p *Peer) Name() string {
	s := p.rw.name
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

// Fullname returns the node name that the remote node advertised.
func (p *Peer) Fullname() string {
	return p.rw.name
}

// Caps returns the capabilities (supported subprotocols) of the remote peer.
func (p *Peer) Caps() []Cap {
	// TODO: maybe return copy
	return p.rw.caps
}

// RunningCap returns true if the peer is actively connected using any of the
// enumerated versions of a specific protocol, meaning that at least one of the
// versions is supported by both this node and the peer p.
func (p *Peer) RunningCap(protocol string, versions []uint) bool {
	if proto, ok := p.running[protocol]; ok {
		for _, ver := range versions {
			if proto.Version == ver {
				return true
			}
		}
	}
	return false
}

// RemoteAddr returns the remote address of the network connection.
func (p *Peer) RemoteAddr() net.Addr {
	return p.rw.fd.RemoteAddr()
}

// LocalAddr returns the local address of the network connection.
func (p *Peer) LocalAddr() net.Addr {
	return p.rw.fd.LocalAddr()
}

// Disconnect terminates the peer connection with the given reason.
// It returns immediately and does not wait until the connection is closed.
func (p *Peer) Disconnect(reason DiscReason) {
	if p.testPipe != nil {
		p.testPipe.Close()
	}

	select {
	case p.disc <- reason:
	case <-p.closed:
	}
}

// String implements fmt.Stringer.
func (p *Peer) String() string {
	id := p.ID()
	return fmt.Sprintf("Peer %x %v", id[:8], p.RemoteAddr())
}

// Inbound returns true if the peer is an inbound connection
func (p *Peer) Inbound() bool {
	return p.rw.is(inboundConn)
}

func newPeer(log log.Logger, conn *conn, protocols []Protocol) *Peer {
	protomap := matchProtocols(protocols, conn.caps, conn)
	p := &Peer{
		rw:       conn,
		running:  protomap,
		created:  mclock.Now(),
		disc:     make(chan DiscReason),
		protoErr: make(chan error, len(protomap)+1), // protocols + pingLoop
		closed:   make(chan struct{}),
		log:      log.New("id", conn.node.ID(), "conn", conn.flags),
		Msgs: 	&MsgList {
					MuxCh: make(chan bool, 1),
					ExistMsg: make(chan struct{}, 1),
					Msg17s: make([]*Msg, 0),
					MsgNot17s: make([]*Msg, 0),
				},
	}
	return p
}

func (p *Peer) Log() log.Logger {
	return p.log
}

func (p *Peer) run() (remoteRequested bool, err error) {
	var (
		writeStart = make(chan struct{}, 1)
		writeErr   = make(chan error, 1)
		readErr    = make(chan error, 1)
		reason     DiscReason // sent to the peer
	)
	p.wg.Add(2)
	go p.readLoop(readErr)
	go p.pingLoop()

	// Start all protocol handlers.
	writeStart <- struct{}{}
	p.startProtocols(writeStart, writeErr)

	// Wait for an error or disconnect.
loop:
	for {
		select {
		case err = <-writeErr:
			// A write finished. Allow the next write to start if
			// there was no error.
			if err != nil {
				reason = DiscNetworkError
				log.Info("Peer|run()|disconnect TCP|writeErr", "reason", reason.String())
				break loop
			}
			writeStart <- struct{}{}
		case err = <-readErr:
			if r, ok := err.(DiscReason); ok {
				remoteRequested = true
				reason = r
			} else {
				reason = DiscNetworkError
			}
			log.Info("Peer|run()|disconnect TCP|readErr", "reason", reason.String())
			break loop
		case err = <-p.protoErr:
			reason = discReasonForError(err)
			log.Info("Peer|run()|disconnect TCP|p.protoErr", "reason", reason.String())
			break loop
		case err = <-p.disc:
			reason = discReasonForError(err)
			log.Info("Peer|run()|disconnect TCP|p.disc", "reason", reason.String())
			break loop
		}
	}

	close(p.closed)
	p.rw.close(reason)
	p.wg.Wait()
	return remoteRequested, err
}

func (p *Peer) pingLoop() {
	ping := time.NewTimer(pingInterval)
	defer p.wg.Done()
	defer ping.Stop()
	for {
		select {
		case <-ping.C:
			if err := SendItems(p.rw, pingMsg); err != nil {
				p.protoErr <- err
				return
			}
			ping.Reset(pingInterval)
		case <-p.closed:
			return
		}
	}
}

//func (p *Peer) readLoop(errc chan<- error) {
//	defer p.wg.Done()
//	for {
//		msg, err := p.rw.ReadMsg()
//		if msg.Code == 17{
//			log.Info("caver|readMsg|17", "msgCode", msg.Code, "err", err)
//		}
//		if err != nil {
//			errc <- err
//			return
//		}
//		msg.ReceivedAt = time.Now()
//		if err = p.handle(msg); err != nil {
//			if msg.Code ==17 {
//				log.Info("caver|handleMsg|17", "msgCode", msg.Code, "err", err)
//			}
//			errc <- err
//			return
//		}
//	}
//}


type MsgList struct {
	MuxCh chan bool
	ExistMsg chan struct{}
	Msg17s []*Msg
	MsgNot17s []*Msg
}

func (ml *MsgList) Add(msg *Msg, code uint64, ibftCode uint64) {
	ml.MuxCh<- true
	if code == 17 && (ibftCode == 0 || ibftCode == 3) {
		ml.Msg17s = append(ml.Msg17s, msg)
	} else {
		ml.MsgNot17s = append(ml.MsgNot17s, msg)
	}
	<-ml.MuxCh
}

func (ml *MsgList) Pop() *Msg {
	var msg *Msg

	ml.MuxCh<- true
	if len(ml.Msg17s) > 0 {
		msg = ml.Msg17s[0]
		if len(ml.Msg17s) > 1 {
			ml.Msg17s = ml.Msg17s[1:]
		} else {
			ml.Msg17s = ml.Msg17s[:0]
		}
	} else if len(ml.MsgNot17s) > 0 {
		msg = ml.MsgNot17s[0]
		if len(ml.MsgNot17s) > 1 {
			ml.MsgNot17s = ml.MsgNot17s[1:]
		} else {
			ml.MsgNot17s = ml.MsgNot17s[:0]
		}
	}
	<-ml.MuxCh

	return msg
}

func (ml *MsgList) NoticeExistMsg() {
	for {
		select {
		case <-time.After(100 * time.Millisecond):

			if len(ml.Msg17s) > 0 {
				for i := len(ml.Msg17s); i > 0; i-- {
					ml.ExistMsg <- struct{}{}
				}
			}
			if len(ml.MsgNot17s) > 0 {
				for i := len(ml.MsgNot17s); i > 0; i-- {
					ml.ExistMsg <- struct{}{}
				}
			}

		}
 	}
}



func (p *Peer) ReadMsg(errc chan<- error) {
	var code uint64
	var ibftCode uint64
	for {
		ibftCode = 9999
		msg, err := p.rw.ReadMsg()
		if msg.Code == 17{
			log.Info("caver|readMsg|17", "msgCode", msg.Code, "err", err)
		}
		if err != nil {
			errc <- err
			return
		}

		msg.ReceivedAt = time.Now()

		if msg.Code == 17 {
			code = 17
		} else {
			proto, err := p.getProto(msg.Code)
			if err == nil {
				code = msg.Code - proto.offset
				log.Info("Peer.ReadMsg()", "msg.Code", msg.Code, "proto.offset", proto.offset,
					"peer.id", p.ID().String())
				// for test
				tempMsg := Copy(msg).(Msg)
				tempPayload, _ := ioutil.ReadAll(msg.Payload)
				////fmt.Println(time.Now().UnixNano()/1e6, "ProtocolManager.handleMsg() msg.Code=", msg.Code, "tempPayload=", tempPayload)
				msg.Payload = bytes.NewReader(tempPayload)
				tempMsg.Payload = bytes.NewReader(tempPayload)
				//var tempdata []byte
				//tempMsg.Decode(&tempdata)
				////fmt.Println(time.Now().UnixNano()/1e6, "ProtocolManager.handleMsg() tempMsg.Code=", tempMsg.Code, "tempMsg.tempdata=", tempdata)
				if code == 17 {
					ibftMsg, _, err := Decode17Msg(tempMsg)
					if err != nil {
						log.Error("Peer.ReadMsg()", "Decode17Msg err", err)
					}
					ibftCode = ibftMsg.Code
				}
			}
		}

		p.Msgs.Add(&msg, code, ibftCode)

		log.Info("Peer.ReadMsg()", "msg.code", code, "p.Msgs.Msg17s lens", len(p.Msgs.Msg17s),
			"p.Msgs.MsgNot17s lens", len(p.Msgs.MsgNot17s), "peer.id", p.ID().String())
	}
}

//MsgPreprepare uint64 = iota
//MsgPrepare
//MsgCommit
//MsgRoundChange

type Message struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

func (m *Message) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

type Preprepare struct {
	View     *View
	Proposal Proposal
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Preprepare) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Proposal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Preprepare) DecodeRLP(s *rlp.Stream) error {
	var preprepare struct {
		View     *View
		Proposal *types.Block
	}

	if err := s.Decode(&preprepare); err != nil {
		return err
	}
	b.View, b.Proposal = preprepare.View, preprepare.Proposal

	return nil
}

type View struct {
	Round    *big.Int
	Sequence *big.Int
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (v *View) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{v.Round, v.Sequence})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (v *View) DecodeRLP(s *rlp.Stream) error {
	var view struct {
		Round    *big.Int
		Sequence *big.Int
	}

	if err := s.Decode(&view); err != nil {
		return err
	}
	v.Round, v.Sequence = view.Round, view.Sequence
	return nil
}

func (v *View) String() string {
	return fmt.Sprintf("{Round: %d, Sequence: %d}", v.Round.Uint64(), v.Sequence.Uint64())
}

// Cmp compares v and y and returns:
//   -1 if v <  y
//    0 if v == y
//   +1 if v >  y
func (v *View) Cmp(y *View) int {
	if v.Sequence.Cmp(y.Sequence) != 0 {
		return v.Sequence.Cmp(y.Sequence)
	}
	if v.Round.Cmp(y.Round) != 0 {
		return v.Round.Cmp(y.Round)
	}
	return 0
}

// Proposal supports retrieving height and serialized block to be used during Istanbul consensus.
type Proposal interface {
	// Number retrieves the sequence number of this proposal.
	Number() *big.Int

	// Hash retrieves the hash of this proposal.
	Hash() common.Hash

	EncodeRLP(w io.Writer) error

	DecodeRLP(s *rlp.Stream) error

	String() string
}

type Subject struct {
	View   *View
	Digest common.Hash
}

// EncodeRLP serializes b into the Ethereum RLP format.
func (b *Subject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{b.View, b.Digest})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (b *Subject) DecodeRLP(s *rlp.Stream) error {
	var subject struct {
		View   *View
		Digest common.Hash
	}

	if err := s.Decode(&subject); err != nil {
		return err
	}
	b.View, b.Digest = subject.View, subject.Digest
	return nil
}

func (b *Subject) String() string {
	return fmt.Sprintf("{View: %v, Digest: %v}", b.View, b.Digest.String())
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

func Decode17Msg(msg Msg) (*Message, common.Hash, error) {
	var data []byte
	if err := msg.Decode(&data); err != nil {
		return nil, common.Hash{}, errors.New("Decode17Msg decode p2p.Msg error")
	}

	msg17 := new(Message)
	err := rlp.DecodeBytes(data, &msg17)
	if err != nil {
		return nil, common.Hash{}, err
	}

	switch msg17.Code {
	case 0:
		var preprepare *Preprepare
		err := msg17.Decode(&preprepare)
		if err != nil {
			return nil, common.Hash{}, err
		}
		log.Info("Decode17Msg()", "ibfttypes.MsgPreprepare Proposal Number", preprepare.Proposal.Number().Uint64(),
			"Sequence", preprepare.View.Sequence.Uint64(), "Round", preprepare.View.Round.Uint64())
	case 1:

	case 2:

	case 3:
		var rc Subject
		if err := msg17.Decode(&rc); err != nil {
			return nil, common.Hash{}, err
		}
		log.Info("Decode17Msg()", "ibfttypes.MsgRoundChange Sequence", rc.View.Sequence.Uint64(), "Round", rc.View.Round.Uint64())

	default:
	}

	return msg17, RLPHash(data), nil
}


func (p *Peer) HandeMsg(errc chan<- error) {
	for {
		select {
		case <-p.Msgs.ExistMsg:
			msg := p.Msgs.Pop()
			if msg == nil {
				continue
			}
			if err := p.handle(*msg); err != nil {
				if msg.Code ==17 {
					log.Info("caver|handleMsg|17", "msgCode", msg.Code, "err", err)
				}
				errc <- err
				return
			}
		}
	}
}

func (p *Peer) readLoop(errc chan<- error) {
	defer p.wg.Done()
	go p.Msgs.NoticeExistMsg()
	go p.ReadMsg(errc)
	go p.HandeMsg(errc)
}

func (p *Peer) handle(msg Msg) error {
	switch {
	case msg.Code == pingMsg:
		msg.Discard()
		go SendItems(p.rw, pongMsg)
	case msg.Code == discMsg:
		var reason [1]DiscReason
		// This is the last message. We don't need to discard or
		// check errors because, the connection will be closed after it.
		rlp.Decode(msg.Payload, &reason)
		return reason[0]
	case msg.Code < baseProtocolLength:
		// ignore other base protocol messages
		return msg.Discard()
	default:
		// it's a subprotocol message
		proto, err := p.getProto(msg.Code)
		if err != nil {
			return fmt.Errorf("msg code out of range: %v", msg.Code)
		}
		if metrics.Enabled {
			m := fmt.Sprintf("%s/%s/%d/%#02x", ingressMeterName, proto.Name, proto.Version, msg.Code-proto.offset)
			metrics.GetOrRegisterMeter(m, nil).Mark(int64(msg.meterSize))
			metrics.GetOrRegisterMeter(m+"/packets", nil).Mark(1)
		}
		select {
		case proto.in <- msg:
			log.Info("Peer.handle()", "msg.Code", msg.Code, "proto.offset", proto.offset,
				"peer.id", p.ID().String())
			return nil
		case <-p.closed:
			return io.EOF
		}
	}
	return nil
}

func countMatchingProtocols(protocols []Protocol, caps []Cap) int {
	n := 0
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				n++
			}
		}
	}
	return n
}

// matchProtocols creates structures for matching named subprotocols.
func matchProtocols(protocols []Protocol, caps []Cap, rw MsgReadWriter) map[string]*protoRW {
	sort.Sort(capsByNameAndVersion(caps))
	offset := baseProtocolLength
	result := make(map[string]*protoRW)

outer:
	for _, cap := range caps {
		for _, proto := range protocols {
			if proto.Name == cap.Name && proto.Version == cap.Version {
				// If an old protocol version matched, revert it
				if old := result[cap.Name]; old != nil {
					offset -= old.Length
				}
				// Assign the new match
				result[cap.Name] = &protoRW{Protocol: proto, offset: offset, in: make(chan Msg), w: rw}
				offset += proto.Length

				continue outer
			}
		}
	}
	return result
}

func (p *Peer) startProtocols(writeStart <-chan struct{}, writeErr chan<- error) {
	p.wg.Add(len(p.running))
	for _, proto := range p.running {
		proto := proto
		proto.closed = p.closed
		proto.wstart = writeStart
		proto.werr = writeErr

		proto.Msgs = &MsgList {
		MuxCh: make(chan bool, 1),
		ExistMsg: make(chan struct{}, 1),
		Msg17s: make([]*Msg, 0),
		MsgNot17s: make([]*Msg, 0),
	}

		var rw MsgReadWriter = proto
		if p.events != nil {
			rw = newMsgEventer(rw, p.events, p.ID(), proto.Name, p.Info().Network.RemoteAddress, p.Info().Network.LocalAddress)
		}
		p.log.Trace(fmt.Sprintf("Starting protocol %s/%d", proto.Name, proto.Version))
		go func() {
			defer p.wg.Done()

			go proto.WriteSocket()
			go proto.Msgs.NoticeExistMsg()

			err := proto.Run(p, rw)
			if err == nil {
				p.log.Trace(fmt.Sprintf("Protocol %s/%d returned", proto.Name, proto.Version))
				err = errProtocolReturned
			} else if err != io.EOF {
				p.log.Trace(fmt.Sprintf("Protocol %s/%d failed", proto.Name, proto.Version), "err", err)
			}
			p.protoErr <- err
		}()
	}
}

// getProto finds the protocol responsible for handling
// the given message code.
func (p *Peer) getProto(code uint64) (*protoRW, error) {
	for _, proto := range p.running {
		if code >= proto.offset && code < proto.offset+proto.Length {
			return proto, nil
		}
	}
	return nil, newPeerError(errInvalidMsgCode, "%d", code)
}

type protoRW struct {
	Protocol
	in     chan Msg        // receives read messages
	closed <-chan struct{} // receives when peer is shutting down
	wstart <-chan struct{} // receives when write may start
	werr   chan<- error    // for write results
	offset uint64
	w      MsgWriter

	Msgs *MsgList
}

//func (rw *protoRW) WriteMsg(msg Msg) (err error) {
//	if msg.Code >= rw.Length {
//		return newPeerError(errInvalidMsgCode, "not handled")
//	}
//	msg.meterCap = rw.cap()
//	msg.meterCode = msg.Code
//
//	log.Info("protoRW.WriteMsg()", "msg.Code", msg.Code, "proto.offset", rw.offset)
//
//	msg.Code += rw.offset
//
//	select {
//	case <-rw.wstart:
//		err = rw.w.WriteMsg(msg)
//		// Report write status back to Peer.run. It will initiate
//		// shutdown if the error is non-nil and unblock the next write
//		// otherwise. The calling protocol code should exit for errors
//		// as well but we don't want to rely on that.
//		rw.werr <- err
//	case <-rw.closed:
//		err = ErrShuttingDown
//	}
//	return err
//}

func (rw *protoRW) WriteMsg(msg Msg) (err error) {
	var ibftCode uint64
	if msg.Code >= rw.Length {
		return newPeerError(errInvalidMsgCode, "not handled")
	}
	msg.meterCap = rw.cap()
	msg.meterCode = msg.Code

	log.Info("protoRW.WriteMsg()", "msg.Code", msg.Code, "proto.offset", rw.offset)

	msg.Code += rw.offset

	select {
	case <-rw.wstart:
		code := msg.Code - rw.offset
		// for test
		tempMsg := Copy(msg).(Msg)
		tempPayload, _ := ioutil.ReadAll(msg.Payload)
		////fmt.Println(time.Now().UnixNano()/1e6, "ProtocolManager.handleMsg() msg.Code=", msg.Code, "tempPayload=", tempPayload)
		msg.Payload = bytes.NewReader(tempPayload)
		tempMsg.Payload = bytes.NewReader(tempPayload)
		//var tempdata []byte
		//tempMsg.Decode(&tempdata)
		////fmt.Println(time.Now().UnixNano()/1e6, "ProtocolManager.handleMsg() tempMsg.Code=", tempMsg.Code, "tempMsg.tempdata=", tempdata)
		if code == 17 {
			ibftMsg, _, err := Decode17Msg(tempMsg)
			if err != nil {
				log.Error("protoRW.WriteMsg()", "Decode17Msg err", err)
			}
			ibftCode = ibftMsg.Code
		}
		rw.Msgs.Add(&msg, code, ibftCode)

	case <-rw.closed:
		err = ErrShuttingDown
	}
	return err
}

func (rw *protoRW) WriteSocket() (err error) {
	for {
		select {
		case <-rw.Msgs.ExistMsg:
			msg := rw.Msgs.Pop()
			if msg == nil {
				continue
			}

			err = rw.w.WriteMsg(*msg)
			// Report write status back to Peer.run. It will initiate
			// shutdown if the error is non-nil and unblock the next write
			// otherwise. The calling protocol code should exit for errors
			// as well but we don't want to rely on that.
			rw.werr <- err
		}
	}
}

func (rw *protoRW) ReadMsg() (Msg, error) {
	select {
	case msg := <-rw.in:
		msg.Code -= rw.offset
		log.Info("protoRW.ReadMsg()", "code", msg.Code)
		return msg, nil
	case <-rw.closed:
		return Msg{}, io.EOF
	}
}

// PeerInfo represents a short summary of the information known about a connected
// peer. Sub-protocol independent fields are contained and initialized here, with
// protocol specifics delegated to all connected sub-protocols.
type PeerInfo struct {
	ENR     string   `json:"enr,omitempty"` // Ethereum Node Record
	Enode   string   `json:"enode"`         // Node URL
	ID      string   `json:"id"`            // Unique node identifier
	Name    string   `json:"name"`          // Name of the node, including client type, version, OS, custom data
	Caps    []string `json:"caps"`          // Protocols advertised by this peer
	Network struct {
		LocalAddress  string `json:"localAddress"`  // Local endpoint of the TCP data connection
		RemoteAddress string `json:"remoteAddress"` // Remote endpoint of the TCP data connection
		Inbound       bool   `json:"inbound"`
		Trusted       bool   `json:"trusted"`
		Static        bool   `json:"static"`
	} `json:"network"`
	Protocols map[string]interface{} `json:"protocols"` // Sub-protocol specific metadata fields
}

// Info gathers and returns a collection of metadata known about a peer.
func (p *Peer) Info() *PeerInfo {
	// Gather the protocol capabilities
	var caps []string
	for _, cap := range p.Caps() {
		caps = append(caps, cap.String())
	}
	// Assemble the generic peer metadata
	info := &PeerInfo{
		Enode:     p.Node().URLv4(),
		ID:        p.ID().String(),
		Name:      p.Fullname(),
		Caps:      caps,
		Protocols: make(map[string]interface{}),
	}
	if p.Node().Seq() > 0 {
		info.ENR = p.Node().String()
	}
	info.Network.LocalAddress = p.LocalAddr().String()
	info.Network.RemoteAddress = p.RemoteAddr().String()
	info.Network.Inbound = p.rw.is(inboundConn)
	info.Network.Trusted = p.rw.is(trustedConn)
	info.Network.Static = p.rw.is(staticDialedConn)

	// Gather all the running protocol infos
	for _, proto := range p.running {
		protoInfo := interface{}("unknown")
		if query := proto.Protocol.PeerInfo; query != nil {
			if metadata := query(p.ID()); metadata != nil {
				protoInfo = metadata
			} else {
				protoInfo = "handshake"
			}
		}
		info.Protocols[proto.Name] = protoInfo
	}
	return info
}
