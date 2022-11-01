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
	"crypto/ecdsa"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	qbftengine "github.com/ethereum/go-ethereum/consensus/istanbul/qbft/engine"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type testerVote struct {
	validator string
	voted     string
	auth      bool
}

// testerAccountPool is a pool to maintain currently active tester accounts,
// mapped from textual names used in the tests below to actual Ethereum private
// keys capable of signing transactions.
type testerAccountPool struct {
	accounts map[string]*ecdsa.PrivateKey
}

func newTesterAccountPool() *testerAccountPool {
	return &testerAccountPool{
		accounts: make(map[string]*ecdsa.PrivateKey),
	}
}

func (ap *testerAccountPool) writeValidatorVote(header *types.Header, validator string, recipientAddress string, authorize bool) error {
	return qbftengine.ApplyHeaderQBFTExtra(
		header,
		qbftengine.WriteVote(ap.address(recipientAddress), authorize),
	)
}

func (ap *testerAccountPool) address(account string) common.Address {
	// Ensure we have a persistent key for the account
	if ap.accounts[account] == nil {
		ap.accounts[account], _ = crypto.GenerateKey()
	}
	// Resolve and return the Ethereum address
	return crypto.PubkeyToAddress(ap.accounts[account].PublicKey)
}

func TestSaveAndLoad(t *testing.T) {
	snap := &Snapshot{
		Epoch:  5,
		Number: 10,
		Hash:   common.HexToHash("1234567890"),
		Votes: []*Vote{
			{
				Validator: common.HexToAddress("1234567891"),
				Block:     15,
				Address:   common.HexToAddress("1234567892"),
				Authorize: false,
			},
		},
		Tally: map[common.Address]Tally{
			common.HexToAddress("1234567893"): {
				Authorize: false,
				Votes:     20,
			},
		},
		ValSet: validator.NewSet([]common.Address{
			common.HexToAddress("1234567894"),
			common.HexToAddress("1234567895"),
		}, istanbul.NewRoundRobinProposerPolicy()),
	}
	db := rawdb.NewMemoryDatabase()
	err := snap.store(db)
	if err != nil {
		t.Errorf("store snapshot failed: %v", err)
	}

	snap1, err := loadSnapshot(snap.Epoch, db, snap.Hash)
	if err != nil {
		t.Errorf("load snapshot failed: %v", err)
	}
	if snap.Epoch != snap1.Epoch {
		t.Errorf("epoch mismatch: have %v, want %v", snap1.Epoch, snap.Epoch)
	}
	if snap.Hash != snap1.Hash {
		t.Errorf("hash mismatch: have %v, want %v", snap1.Number, snap.Number)
	}
	if !reflect.DeepEqual(snap.Votes, snap.Votes) {
		t.Errorf("votes mismatch: have %v, want %v", snap1.Votes, snap.Votes)
	}
	if !reflect.DeepEqual(snap.Tally, snap.Tally) {
		t.Errorf("tally mismatch: have %v, want %v", snap1.Tally, snap.Tally)
	}
	if !reflect.DeepEqual(snap.ValSet, snap.ValSet) {
		t.Errorf("validator set mismatch: have %v, want %v", snap1.ValSet, snap.ValSet)
	}
}
