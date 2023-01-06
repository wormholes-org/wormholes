package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

const DEFAULT_VALIDATOR_COEFFICIENT = 70

type EmptyMessageEvent struct {
	Sender  common.Address
	Vote    common.Address
	Height  *big.Int
	Payload []byte
}

type EmptyMsg struct {
	Code      uint64
	Msg       []byte
	Address   common.Address
	Signature []byte
}

func (m *EmptyMsg) PayloadNoSig() ([]byte, error) {
	return rlp.EncodeToBytes(&EmptyMsg{
		Code:      m.Code,
		Msg:       m.Msg,
		Address:   m.Address,
		Signature: []byte{},
	})
}

func (m *EmptyMsg) Payload() ([]byte, error) {
	return rlp.EncodeToBytes(m)
}

func (m *EmptyMsg) FromPayload(b []byte) error {
	// Decode message
	err := rlp.DecodeBytes(b, &m)
	if err != nil {
		return err
	}
	return nil
}

func (m *EmptyMsg) Decode(val interface{}) error {
	return rlp.DecodeBytes(m.Msg, val)
}

func (m *EmptyMsg) RecoverAddress(b []byte) (common.Address, error) {
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

type OnlineZkQuestion struct {
	Height *big.Int // block height
}

type SignatureData struct {
	Vote   common.Address
	Height *big.Int
	//Timestamp uint64
}

type OnlineZkProof struct {
	Height *big.Int
	Proof  []byte
}
