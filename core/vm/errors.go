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

package vm

import (
	"errors"
	"fmt"
)

// List evm execution errors
var (
	ErrOutOfGas                     = errors.New("out of gas")
	ErrCodeStoreOutOfGas            = errors.New("contract creation code storage out of gas")
	ErrDepth                        = errors.New("max call depth exceeded")
	ErrInsufficientBalance          = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision     = errors.New("contract address collision")
	ErrExecutionReverted            = errors.New("execution reverted")
	ErrMaxCodeSizeExceeded          = errors.New("max code size exceeded")
	ErrInvalidJump                  = errors.New("invalid jump destination")
	ErrWriteProtection              = errors.New("write protection")
	ErrReturnDataOutOfBounds        = errors.New("return data out of bounds")
	ErrGasUintOverflow              = errors.New("gas uint64 overflow")
	ErrInvalidCode                  = errors.New("invalid code: must not begin with 0xef")
	ErrNotOwner                     = errors.New("not the owner of the nft")
	ErrNotExistNFTType              = errors.New("not exist nft type")
	ErrInsufficientPledgedBalance   = errors.New("not sufficient pledged balance")
	ErrStartIndex                   = errors.New("StartIndex is not string of 0x!")
	ErrNotExchanger                 = errors.New("not exchanger")
	ErrWormholesFormat              = errors.New("wormholes format error, can't unmarshal")
	ErrInsufficientExchangerBalance = errors.New("insufficient exchanger balance for transfer")
	ErrNotMoreThan100ERB            = errors.New("not more than 200 ERB")
	ErrTooCloseWithOpenExchanger    = errors.New("too close with openexchanger")
	ErrTooCloseForWithdraw          = errors.New("too close for Withdraw")
	ErrTooCloseToCancel             = errors.New("too close to cancel")
	ErrRoyaltyNotMoreThan0          = errors.New("royalty not more than 0")
	ErrRoyaltyNotLessthan10000      = errors.New("royalty not less than 10000")
	ErrFeeRateNotMoreThan0          = errors.New("feerate not more than 0")
	ErrFeeRateNotLessThan10000      = errors.New("feerate not less than 10000")
	ErrNotMintByOfficial            = errors.New("not mint by official")
	ErrTransAmount                  = errors.New("trans amount must be larger than 0")
	ErrNotMoreThan100000ERB         = errors.New("not more than 50000 ERB")
	ErrNotAllowedOfficialNFT        = errors.New("not allowed approve official nft")
	ErrExchangerFormat              = errors.New("exchanger format error")
	ErrNftLevel                     = errors.New("input nft level error")
	ErrMinerProxy                   = errors.New("cannot delegate repeatedly")
	ErrRepeatedPledge               = errors.New("no repeated pledge")
	ErrReopenExchanger              = errors.New("reopen exchanger")
	ErrNotPledge                    = errors.New("not pledge")
)

// ErrStackUnderflow wraps an evm error when the items on the stack less
// than the minimal requirement.
type ErrStackUnderflow struct {
	stackLen int
	required int
}

func (e *ErrStackUnderflow) Error() string {
	return fmt.Sprintf("stack underflow (%d <=> %d)", e.stackLen, e.required)
}

// ErrStackOverflow wraps an evm error when the items on the stack exceeds
// the maximum allowance.
type ErrStackOverflow struct {
	stackLen int
	limit    int
}

func (e *ErrStackOverflow) Error() string {
	return fmt.Sprintf("stack limit reached %d (%d)", e.stackLen, e.limit)
}

// ErrInvalidOpCode wraps an evm error when an invalid opcode is encountered.
type ErrInvalidOpCode struct {
	opcode OpCode
}

func (e *ErrInvalidOpCode) Error() string { return fmt.Sprintf("invalid opcode: %s", e.opcode) }
