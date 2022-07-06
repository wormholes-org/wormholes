// Copyright 2017 The go-ethereum Authors
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

package backend

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSign(t *testing.T) {
	b := newBackend()
	defer b.Stop()
	data := []byte("Here is a string....")
	sig, err := b.Sign(data)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	//Check signature recover
	hashData := crypto.Keccak256(data)
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])
	if signer != getAddress() {
		t.Errorf("address mismatch: have %v, want %s", signer.Hex(), getAddress().Hex())
	}
}

func TestCheckSignature(t *testing.T) {
	key, _ := generatePrivateKey()
	data := []byte("Here is a string....")
	hashData := crypto.Keccak256(data)
	sig, _ := crypto.Sign(hashData, key)
	b := newBackend()
	defer b.Stop()
	a := getAddress()
	err := b.CheckSignature(data, a, sig)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	a = getInvalidAddress()
	err = b.CheckSignature(data, a, sig)
	if err != istanbulcommon.ErrInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, istanbulcommon.ErrInvalidSignature)
	}
}

func TestCheckValidatorSignature(t *testing.T) {
	vset, keys := newTestValidatorSet(5)

	// 1. Positive test: sign with validator's key should succeed
	data := []byte("dummy data")
	hashData := crypto.Keccak256(data)
	for i, k := range keys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		// CheckValidatorSignature should succeed
		addr, err := istanbul.CheckValidatorSignature(vset, data, sig)
		if err != nil {
			t.Errorf("error mismatch: have %v, want nil", err)
		}
		validator := vset.GetByIndex(uint64(i))
		if addr != validator.Address() {
			t.Errorf("validator address mismatch: have %v, want %v", addr, validator.Address())
		}
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	// Sign
	sig, err := crypto.Sign(hashData, key)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := istanbul.CheckValidatorSignature(vset, data, sig)
	if err != istanbul.ErrUnauthorizedAddress {
		t.Errorf("error mismatch: have %v, want %v", err, istanbul.ErrUnauthorizedAddress)
	}
	emptyAddr := common.Address{}
	if addr != emptyAddr {
		t.Errorf("address mismatch: have %v, want %v", addr, emptyAddr)
	}
}

func TestCommit(t *testing.T) {
	backend := newBackend()
	defer backend.Stop()

	commitCh := make(chan *types.Block)
	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	testCases := []struct {
		expectedErr       error
		expectedSignature [][]byte
		expectedBlock     func() *types.Block
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-1)...)},
			func() *types.Block {
				chain, engine := newBlockChain(1, big.NewInt(0))
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				return updateQBFTBlock(block, engine.Address())
			},
		},
		{
			// invalid signature
			istanbulcommon.ErrInvalidCommittedSeals,
			nil,
			func() *types.Block {
				chain, engine := newBlockChain(1, big.NewInt(0))
				block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
				return updateQBFTBlock(block, engine.Address())
			},
		},
	}

	for _, test := range testCases {
		expBlock := test.expectedBlock()
		go func() {
			result := <-backend.commitCh
			commitCh <- result
		}()

		backend.proposedBlockHash = expBlock.Hash()
		if err := backend.Commit(expBlock, test.expectedSignature, big.NewInt(0)); err != nil {
			if err != test.expectedErr {
				t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
			}
		}

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				if result.Hash() != expBlock.Hash() {
					t.Errorf("hash mismatch: have %v, want %v", result.Hash(), expBlock.Hash())
				}
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
	}
}

func TestGetProposer(t *testing.T) {
	chain, engine := newBlockChain(1, big.NewInt(0))
	defer engine.Stop()
	block := makeBlock(chain, engine, chain.Genesis())
	chain.InsertChain(types.Blocks{block})
	expected := engine.GetProposer(1)
	actual := engine.Address()
	if actual != expected {
		t.Errorf("proposer mismatch: have %v, want %v", actual.Hex(), expected.Hex())
	}
}

// TestQBFTTransitionDeadlock test whether a deadlock occurs when testQBFTBlock is set to 1
// This was fixed as part of commit 2a8310663ecafc0233758ca7883676bf568e926e
func TestQBFTTransitionDeadlock(t *testing.T) {
	timeout := time.After(1 * time.Minute)
	done := make(chan bool)
	go func() {
		chain, engine := newBlockChain(1, big.NewInt(1))
		defer engine.Stop()
		// Create an insert a new block into the chain.
		block := makeBlock(chain, engine, chain.Genesis())
		_, err := chain.InsertChain(types.Blocks{block})
		if err != nil {
			t.Errorf("Error inserting block: %v", err)
		}

		if err = engine.NewChainHead(); err != nil {
			t.Errorf("Error posting NewChainHead Event: %v", err)
		}

		if !engine.IsQBFTConsensus() {
			t.Errorf("IsQBFTConsensus() should return true after block insertion")
		}
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Deadlock occurred during IBFT to QBFT transition")
	case <-done:
	}
}

func TestIsQBFTConsensus(t *testing.T) {
	chain, engine := newBlockChain(1, big.NewInt(2))
	defer engine.Stop()
	qbftConsensus := engine.IsQBFTConsensus()
	if qbftConsensus {
		t.Errorf("IsQBFTConsensus() should return false")
	}

	// Create an insert a new block into the chain.
	block := makeBlock(chain, engine, chain.Genesis())
	_, err := chain.InsertChain(types.Blocks{block})
	if err != nil {
		t.Errorf("Error inserting block: %v", err)
	}

	if err = engine.NewChainHead(); err != nil {
		t.Errorf("Error posting NewChainHead Event: %v", err)
	}

	secondBlock := makeBlock(chain, engine, block)
	_, err = chain.InsertChain(types.Blocks{secondBlock})
	if err != nil {
		t.Errorf("Error inserting block: %v", err)
	}

	qbftConsensus = engine.IsQBFTConsensus()
	if !qbftConsensus {
		t.Errorf("IsQBFTConsensus() should return true after block insertion")
	}
}

/**
 * SimpleBackend
 * Private key: bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1
 * Public key: 04a2bfb0f7da9e1b9c0c64e14f87e8fb82eb0144e97c25fe3a977a921041a50976984d18257d2495e7bfd3d4b280220217f429287d25ecdf2b0d7c0f7aae9aa624
 * Address: 0x70524d664ffe731100208a0154e556f9bb679ae6
 */
func getAddress() common.Address {
	return common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
}

func getInvalidAddress() common.Address {
	return common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")
}

func generatePrivateKey() (*ecdsa.PrivateKey, error) {
	key := "bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1"
	return crypto.HexToECDSA(key)
}

func newTestValidatorSet(n int) (istanbul.ValidatorSet, []*ecdsa.PrivateKey) {
	// generate validators
	keys := make(Keys, n)
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		keys[i] = privateKey
		addrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
	}
	vset := validator.NewSet(addrs, istanbul.NewRoundRobinProposerPolicy())
	sort.Sort(keys) //Keys need to be sorted by its public key address
	return vset, keys
}

type Keys []*ecdsa.PrivateKey

func (slice Keys) Len() int {
	return len(slice)
}

func (slice Keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice Keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func newBackend() (b *Backend) {
	_, b = newBlockChain(1, big.NewInt(0))
	key, _ := generatePrivateKey()
	b.privateKey = key
	return
}

func TestGetExtraData(t *testing.T){
	address := []string{
		"0x4991Dd8c3307b00cB3d022496591B03BC4a606f4",
		"0xb54E590418aba304590880aB2190fB49Ed6e09fb",
		"0x505DDba1c6F04193BEc16efdFac2153cd8129166",
		"0xaa688d41E4851Be887f7048Be05545BE36c5Ba4E",
		"0x8203dC4bdedd0b75E06c19A24cF04A0Dc7b65ff3",
		"0xa9dCa9B08E6073027FBd39646D6f34C7c72B89c5",
		"0x859b13333fC4BeBD773384f967C483aF6F950c95",
		"0x03fD4f1000bd771065046D048559B49e630e483C",
		"0x426a5cB3c4d162854136c233Ad6387D8a9C67A85",
		"0x25328C16AE3EdA0a06B7352D724811c93d47Ff8a",
		"0x3FDE68a0D0597aFb92F642bD5b1cC65293e2B919",
		"0xD8a52B98cD120e709027234098b5B9302e422A3D",
		"0x8bfD0f755b09698a0b462Df4146A498864e5fD44",
		"0x8dD1d68f9F5B55cEefF5AbB7999067ae80023B02",
		"0x1FD75AB3e74FE70Db59FCEe6E6635086E3EACc47",
	}
	Replace, _ := istExtraEncodeString(address)
	fmt.Println(Replace)
}

func istExtraEncodeString(vlds []string) (string, error) {
	validators := []common.Address{}
	for _, v := range vlds {
		validators = append(validators, common.HexToAddress(v))
	}
	return istExtraEncode("0x00", validators)
}

func istExtraEncode(vanity string, validators []common.Address) (string, error) {
	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return "", err
	}

	if len(newVanity) < types.IstanbulExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    validators,
		Seal:          make([]byte, types.IstanbulExtraSeal),
		CommittedSeal: [][]byte{},
		BeneficiaryAddr: []common.Address{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(newVanity, payload...)), nil
}
