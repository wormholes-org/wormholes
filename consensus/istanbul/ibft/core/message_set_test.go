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
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/istanbul"
	ibfttypes "github.com/ethereum/go-ethereum/consensus/istanbul/ibft/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func TestMessageSetWithPreprepare(t *testing.T) {
	valSet := newTestValidatorSet(4)

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
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	msg := &ibfttypes.Message{
		Code:    ibfttypes.MsgPreprepare,
		Msg:     rawPP,
		Address: valSet.GetProposer().Address(),
	}

	err = ms.Add(msg)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	err = ms.Add(msg)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}

func TestMessageSetWithSubject(t *testing.T) {
	valSet := newTestValidatorSet(4)

	ms := newMessageSet(valSet)

	view := &istanbul.View{
		Round:    new(big.Int),
		Sequence: new(big.Int),
	}

	sub := &istanbul.Subject{
		View:   view,
		Digest: common.HexToHash("1234567890"),
	}

	rawSub, err := rlp.EncodeToBytes(sub)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	msg := &ibfttypes.Message{
		Code:    ibfttypes.MsgPrepare,
		Msg:     rawSub,
		Address: valSet.GetProposer().Address(),
	}

	err = ms.Add(msg)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	err = ms.Add(msg)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if ms.Size() != 1 {
		t.Errorf("the size of message set mismatch: have %v, want 1", ms.Size())
	}
}

func TestExt(t *testing.T) {
	valSet := newTestValidatorSet(7)
	ms := newMessageSet(valSet)
	msg := &ibfttypes.Message{
		Code:    ibfttypes.MsgPrepare,
		Msg:     common.LocalRandomBytes(),
		Address: valSet.GetProposer().Address(),
	}
	time.Sleep(1)
	var err = ms.Add(msg)
	if err != nil {
		t.Error("add fail")
		return
	}
	msg = &ibfttypes.Message{
		Code:    ibfttypes.MsgPrepare,
		Msg:     common.LocalRandomBytes(),
		Address: valSet.GetByIndex(5).Address(),
	}
	time.Sleep(1)

	if err = ms.Add(msg); err != nil {
		t.Error(err)
	}

	msg = &ibfttypes.Message{
		Code:    ibfttypes.MsgPrepare,
		Msg:     []byte{},
		Address: valSet.GetByIndex(6).Address(),
	}
	time.Sleep(1)

	if err = ms.Add(msg); err != nil {
		t.Error(err)
	}

	msg = &ibfttypes.Message{
		Code:    ibfttypes.MsgPrepare,
		Msg:     common.LocalRandomBytes(),
		Address: common.HexToAddress("0x0"),
	}
	time.Sleep(1)

	err = ms.Add(msg)
	if err != nil {
		t.Error(err)
	}
	addrs := ms.GetAddrs()
	t.Log(addrs)
	rnd := ms.CalcRandSeed()
	t.Log(rnd)
	code := ms.Encode()
	t.Log(code)
}
