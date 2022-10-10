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

package core

import (
	"bytes"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	istanbulcommon "github.com/ethereum/go-ethereum/consensus/istanbul/common"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/consensus/istanbul/validator"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"testing"
)

func GetMessageSetWithPreprepare() (*messageSet, error) {
	addr0 := common.HexToAddress("0x4110E56ED25e21267FBeEf79244f47ada4e2E963")
	addr1 := common.HexToAddress("0x091DBBa95B26793515cc9aCB9bEb5124c479f27F")
	addr2 := common.HexToAddress("0x107837Ea83f8f06533DDd3fC39451Cd0AA8DA8BD")
	addr3 := common.HexToAddress("0x612DFa56DcA1F581Ed34b9c60Da86f1268Ab6349")
	addr4 := common.HexToAddress("0x84d84e6073A06B6e784241a9B13aA824AB455326")
	addr5 := common.HexToAddress("0x9e4d5C72569465270232ed7Af71981Ee82d08dBF")
	addr6 := common.HexToAddress("0xa270bBDFf450EbbC2d0413026De5545864a1b6d6")
	valSet := validator.NewSet([]common.Address{addr0, addr1, addr2, addr3, addr4, addr5, addr6}, istanbul.NewRoundRobinProposerPolicy())

	ms := newMessageSet(valSet)

	view := &istanbul.View{
		Round:    new(big.Int),
		Sequence: new(big.Int),
	}
	pp := &istanbul.Preprepare{
		View:     view,
		Proposal: makeBlock(1),
	}

	rawPP, err := rlp.EncodeToBytes(pp)
	if err != nil {
		return nil, errors.New("error mismatch: have %v, want nil")
	}
	msg := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: valSet.GetProposer().Address(),
	}
	ms.Add(msg)
	msg0 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr0,
	}
	msg1 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr1,
	}
	msg2 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr2,
	}
	msg3 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr3,
	}
	msg4 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr4,
	}
	msg5 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr5,
	}
	msg6 := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: addr6,
	}
	ms.Add(msg0)
	ms.Add(msg1)
	ms.Add(msg2)
	ms.Add(msg3)
	ms.Add(msg4)
	ms.Add(msg5)
	ms.Add(msg6)
	return ms, nil
}

func TestHandleCommitSeals(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

	proposal := newTestProposal()
	expectedSubject := &istanbul.Subject{
		View: &istanbul.View{
			Round:    big.NewInt(0),
			Sequence: proposal.Number(),
		},
		Digest: proposal.Hash(),
	}

	testCases := []struct {
		system      *testSystem
		expectedErr error
	}{
		{
			// normal case
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&istanbul.View{
							Round:    big.NewInt(0),
							Sequence: big.NewInt(1),
						},
						c.valSet,
					)

					if i == 0 {
						// replica 0 is the proposer
						c.state = ibfttypes.StatePrepared
					}
				}
				return sys
			}(),
			nil,
		},
		{
			// future message
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = ibfttypes.StatePreprepared
					} else {
						c.current = newTestRoundState(
							&istanbul.View{
								Round:    big.NewInt(2),
								Sequence: big.NewInt(3),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			istanbulcommon.ErrFutureMessage,
		},
		{
			// subject not match
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = ibfttypes.StatePreprepared
					} else {
						c.current = newTestRoundState(
							&istanbul.View{
								Round:    big.NewInt(0),
								Sequence: big.NewInt(0),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			istanbulcommon.ErrOldMessage,
		},
		{
			// jump state
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&istanbul.View{
							Round:    big.NewInt(0),
							Sequence: proposal.Number(),
						},
						c.valSet,
					)

					// only replica0 stays at ibfttypes.StatePreprepared
					// other replicas are at ibfttypes.StatePrepared
					if i != 0 {
						c.state = ibfttypes.StatePrepared
					} else {
						c.state = ibfttypes.StatePreprepared
					}
				}
				return sys
			}(),
			nil,
		},
		// TODO: double send message
	}

OUTER:
	for _, test := range testCases {
		test.system.Run(false)

		v0 := test.system.backends[0]
		r0 := v0.engine
		pps, _ := GetMessageSetWithPreprepare()
		encodedCommitSeals, _ := ibfttypes.Encode(pps.Values())
		for i, v := range test.system.backends {
			validator := r0.valSet.GetByIndex(uint64(i))
			m, _ := ibfttypes.Encode(v.engine.current.Subject())
			if err := r0.handleCommit(&ibfttypes.Message{
				Code:       ibfttypes.MsgCommit,
				Msg:        m,
				Address:    validator.Address(),
				Signature:  []byte{},
				Rewardlist: encodedCommitSeals, // small hack
			}, validator); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
				if r0.current.IsHashLocked() {
					t.Errorf("block should not be locked")
				}
				continue OUTER
			}
		}

		// prepared is normal case
		if r0.state != ibfttypes.StateCommitted {
			// There are not enough commit messages in core
			if r0.state != ibfttypes.StatePrepared {
				t.Errorf("state mismatch: have %v, want %v", r0.state, ibfttypes.StatePrepared)
			}
			if r0.current.Commits.Size() >= r0.QuorumSize() {
				t.Errorf("the size of commit messages should be less than %v", r0.QuorumSize())
			}
			if r0.current.IsHashLocked() {
				t.Errorf("block should not be locked")
			}
			continue
		}

		// core should have 2F+1 before Ceil2Nby3Block or Ceil(2N/3) prepare messages
		if r0.current.Commits.Size() < r0.QuorumSize() {
			t.Errorf("the size of commit messages should be larger than 2F+1 or Ceil(2N/3): size %v", r0.QuorumSize())
		}

		// check signatures large than F
		signedCount := 0
		committedSeals := v0.committedMsgs[0].committedSeals
		for _, validator := range r0.valSet.List() {
			for _, seal := range committedSeals {
				if bytes.Equal(validator.Address().Bytes(), seal[:common.AddressLength]) {
					signedCount++
					break
				}
			}
		}
		if signedCount <= r0.valSet.F() {
			t.Errorf("the expected signed count should be larger than %v, but got %v", r0.valSet.F(), signedCount)
		}
		if !r0.current.IsHashLocked() {
			t.Errorf("block should be locked")
		}
	}
}

func TestHandleCommit(t *testing.T) {
	N := uint64(4)
	F := uint64(1)

	proposal := newTestProposal()
	expectedSubject := &istanbul.Subject{
		View: &istanbul.View{
			Round:    big.NewInt(0),
			Sequence: proposal.Number(),
		},
		Digest: proposal.Hash(),
	}

	testCases := []struct {
		system      *testSystem
		expectedErr error
	}{
		{
			// normal case
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&istanbul.View{
							Round:    big.NewInt(0),
							Sequence: big.NewInt(1),
						},
						c.valSet,
					)

					if i == 0 {
						// replica 0 is the proposer
						c.state = ibfttypes.StatePrepared
					}
				}
				return sys
			}(),
			nil,
		},
		{
			// future message
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = ibfttypes.StatePreprepared
					} else {
						c.current = newTestRoundState(
							&istanbul.View{
								Round:    big.NewInt(2),
								Sequence: big.NewInt(3),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			istanbulcommon.ErrFutureMessage,
		},
		{
			// subject not match
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					if i == 0 {
						// replica 0 is the proposer
						c.current = newTestRoundState(
							expectedSubject.View,
							c.valSet,
						)
						c.state = ibfttypes.StatePreprepared
					} else {
						c.current = newTestRoundState(
							&istanbul.View{
								Round:    big.NewInt(0),
								Sequence: big.NewInt(0),
							},
							c.valSet,
						)
					}
				}
				return sys
			}(),
			istanbulcommon.ErrOldMessage,
		},
		{
			// jump state
			func() *testSystem {
				sys := NewTestSystemWithBackend(N, F)

				for i, backend := range sys.backends {
					c := backend.engine
					c.valSet = backend.peers
					c.current = newTestRoundState(
						&istanbul.View{
							Round:    big.NewInt(0),
							Sequence: proposal.Number(),
						},
						c.valSet,
					)

					// only replica0 stays at ibfttypes.StatePreprepared
					// other replicas are at ibfttypes.StatePrepared
					if i != 0 {
						c.state = ibfttypes.StatePrepared
					} else {
						c.state = ibfttypes.StatePreprepared
					}
				}
				return sys
			}(),
			nil,
		},
		// TODO: double send message
	}

OUTER:
	for _, test := range testCases {
		test.system.Run(false)

		v0 := test.system.backends[0]
		r0 := v0.engine

		for i, v := range test.system.backends {
			validator := r0.valSet.GetByIndex(uint64(i))
			m, _ := ibfttypes.Encode(v.engine.current.Subject())
			if err := r0.handleCommit(&ibfttypes.Message{
				Code:          ibfttypes.MsgCommit,
				Msg:           m,
				Address:       validator.Address(),
				Signature:     []byte{},
				CommittedSeal: validator.Address().Bytes(), // small hack
			}, validator); err != nil {
				if err != test.expectedErr {
					t.Errorf("error mismatch: have %v, want %v", err, test.expectedErr)
				}
				if r0.current.IsHashLocked() {
					t.Errorf("block should not be locked")
				}
				continue OUTER
			}
		}

		// prepared is normal case
		if r0.state != ibfttypes.StateCommitted {
			// There are not enough commit messages in core
			if r0.state != ibfttypes.StatePrepared {
				t.Errorf("state mismatch: have %v, want %v", r0.state, ibfttypes.StatePrepared)
			}
			if r0.current.Commits.Size() >= r0.QuorumSize() {
				t.Errorf("the size of commit messages should be less than %v", r0.QuorumSize())
			}
			if r0.current.IsHashLocked() {
				t.Errorf("block should not be locked")
			}
			continue
		}

		// core should have 2F+1 before Ceil2Nby3Block or Ceil(2N/3) prepare messages
		if r0.current.Commits.Size() < r0.QuorumSize() {
			t.Errorf("the size of commit messages should be larger than 2F+1 or Ceil(2N/3): size %v", r0.QuorumSize())
		}

		// check signatures large than F
		signedCount := 0
		committedSeals := v0.committedMsgs[0].committedSeals
		for _, validator := range r0.valSet.List() {
			for _, seal := range committedSeals {
				if bytes.Equal(validator.Address().Bytes(), seal[:common.AddressLength]) {
					signedCount++
					break
				}
			}
		}
		if signedCount <= r0.valSet.F() {
			t.Errorf("the expected signed count should be larger than %v, but got %v", r0.valSet.F(), signedCount)
		}
		if !r0.current.IsHashLocked() {
			t.Errorf("block should be locked")
		}
	}
}

// round is not checked for now
func TestVerifyCommit(t *testing.T) {
	// for log purpose
	privateKey, _ := crypto.GenerateKey()
	peer := validator.New(getPublicKeyAddress(privateKey))
	valSet := validator.NewSet([]common.Address{peer.Address()}, istanbul.NewRoundRobinProposerPolicy())

	sys := NewTestSystemWithBackend(uint64(1), uint64(0))

	testCases := []struct {
		expected   error
		commit     *istanbul.Subject
		roundState *roundState
	}{
		{
			// normal case
			expected: nil,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// old message
			expected: istanbulcommon.ErrInconsistentSubject,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// different digest
			expected: istanbulcommon.ErrInconsistentSubject,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				Digest: common.HexToHash("1234567890"),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// malicious package(lack of sequence)
			expected: istanbulcommon.ErrInconsistentSubject,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(0), Sequence: nil},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(1)},
				valSet,
			),
		},
		{
			// wrong prepare message with same sequence but different round
			expected: istanbulcommon.ErrInconsistentSubject,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(1), Sequence: big.NewInt(0)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
		{
			// wrong prepare message with same round but different sequence
			expected: istanbulcommon.ErrInconsistentSubject,
			commit: &istanbul.Subject{
				View:   &istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(1)},
				Digest: newTestProposal().Hash(),
			},
			roundState: newTestRoundState(
				&istanbul.View{Round: big.NewInt(0), Sequence: big.NewInt(0)},
				valSet,
			),
		},
	}
	for i, test := range testCases {
		c := sys.backends[0].engine
		c.current = test.roundState

		if err := c.verifyCommit(test.commit, peer); err != nil {
			if err != test.expected {
				t.Errorf("result %d: error mismatch: have %v, want %v", i, err, test.expected)
			}
		}
	}
}
