package core

import (
	"errors"
	"fmt"
	"hash"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

type testHasher struct {
	hasher hash.Hash
}

func newHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

func (h *testHasher) Reset() {
	h.hasher.Reset()
}

func (h *testHasher) Update(key, val []byte) {
	h.hasher.Write(key)
	h.hasher.Write(val)
}

func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}

func TestExtractOnlineProof(t *testing.T) {
	curView := &istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(1)}
	hasher := newHasher()
	proposal := types.NewBlock(&types.Header{Number: big.NewInt(int64(0))}, nil, nil, nil, hasher)

	onlineProof := &istanbul.OnlineProof{
		View:       curView,
		Proposal:   proposal,
		RandomHash: common.HexToHash("0x00000000000000000000"),
	}
	onlineProofEnc, _ := ibfttypes.Encode(onlineProof)
	onlineProof.Signature, _ = Sign(onlineProofEnc)
	payload, _ := rlp.EncodeToBytes(onlineProof)

	// decode process

	var onlineProofDec *istanbul.OnlineProof

	err := rlp.DecodeBytes(payload, &onlineProofDec)
	if err != nil {
		fmt.Println("DecodeBytes err", "=======err:   ", err)
	}
	err2 := CheckSignature(onlineProofEnc, common.HexToAddress("0x44d952db5dfB4CBb54443554F4bB9cbeBee2194c"), onlineProofDec.Signature)
	if err2 != nil {
		fmt.Println("TestExtractOnlineProof err", "=======err:   ", err2)
	} else {
		fmt.Println("TestExtractOnlineProof success")
	}
}

func Sign(data []byte) ([]byte, error) {
	priHex := "f616c4d20311a2e73c67ef334630f834b7fb42304a1d4448fb2058e9940ecc0a"
	priKey, _ := crypto.HexToECDSA(priHex)
	hashData := crypto.Keccak256(data)
	return crypto.Sign(hashData, priKey)
}

func CheckSignature(data []byte, address common.Address, sig []byte) error {
	signer, err := istanbul.GetSignatureAddress(data, sig)
	if err != nil {
		return err
	}
	if signer != address {
		return errors.New("ErrInvalidSignature")
	}

	return nil
}
