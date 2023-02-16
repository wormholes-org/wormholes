package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type OnlineMsg struct {
	Addr common.Address
	Msg  *BlockMessage
}

type BlockMessage struct {
	Code          uint64
	Msg           []byte
	Address       common.Address
	Signature     []byte
	CommittedSeal []byte
}

func (m *BlockMessage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, []interface{}{m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal})
}

// DecodeRLP implements rlp.Decoder, and load the consensus fields from a RLP stream.
func (m *BlockMessage) DecodeRLP(s *rlp.Stream) error {
	var msg struct {
		Code          uint64
		Msg           []byte
		Address       common.Address
		Signature     []byte
		CommittedSeal []byte
	}

	if err := s.Decode(&msg); err != nil {
		return err
	}
	m.Code, m.Msg, m.Address, m.Signature, m.CommittedSeal = msg.Code, msg.Msg, msg.Address, msg.Signature, msg.CommittedSeal
	return nil
}

func (m *BlockMessage) FromPayload(b []byte) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}
	return nil
}

// ==============================================
//
// define the functions that needs to be provided for core.
//func (m *BlockMessage) FromPayload(b []byte, validateFn func([]byte, []byte) (common.Address, error)) error {
//	// Decode message
//	err := rlp.DecodeBytes(b, &m)
//	if err != nil {
//		return err
//	}
//
//	// Validate message (on a message without Signature)
//	if validateFn != nil {
//		var payload []byte
//		payload, err = m.PayloadNoSig()
//		if err != nil {
//			return err
//		}
//
//		signerAdd, err := validateFn(payload, m.Signature)
//		if err != nil {
//			return err
//		}
//		if !bytes.Equal(signerAdd.Bytes(), m.Address.Bytes()) {
//			return errors.New("message not signed by the sender")
//		}
//	}
//	return nil
//}

func (m *BlockMessage) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *BlockMessage) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&BlockMessage{
		Code:          m.Code,
		Msg:           m.Msg,
		Address:       m.Address,
		Signature:     []byte{},
		CommittedSeal: m.CommittedSeal,
	})
}

func (m *BlockMessage) RecoverAddress(b []byte) (common.Address, error) {
	err := m.FromPayload(b)
	if err != nil {
		return common.Address{}, err
	}
	data, err := m.PayloadNoSig()
	if err != nil {
		return common.Address{}, err
	}
	hash := crypto.Keccak256(data)
	pub, err := crypto.SigToPub(hash, m.Signature)
	if err != nil {
		return common.Address{}, err
	}
	address := crypto.PubkeyToAddress(*pub)

	return address, nil
}

func (m *BlockMessage) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *BlockMessage) String() string {
	return fmt.Sprintf("{Code: %v, Address: %v}", m.Code, m.Address.String())
}

// ==============================================
//
// helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}
