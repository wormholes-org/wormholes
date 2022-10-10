package miner

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

const (
	WorkerMsg          = 0x12
	SendSignMsg uint64 = iota
)

type MessageEvent struct {
	Code    uint64
	Payload []byte
}

type Msg struct {
	Code      uint64
	Msg       []byte
	Address   common.Address
	Signature []byte
}

func (m *Msg) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&Msg{
		Code:      m.Code,
		Msg:       m.Msg,
		Address:   m.Address,
		Signature: []byte{},
	})
}

func (m *Msg) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *Msg) FromPayload(b []byte) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Msg) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

type OnlineZkQuestion struct {
	Height *big.Int // block height
}

type SignatureData struct {
	Address common.Address
	Height  *big.Int
}

type OnlineZkProof struct {
	Height *big.Int
	Proof  []byte
}

//-----------------------helper functions

func Encode(val interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(val)
}

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}
