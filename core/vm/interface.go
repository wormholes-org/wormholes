// Copyright 2016 The go-ethereum Authors
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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// StateDB is an EVM database for full state querying.
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(common.Address, common.Hash) common.Hash
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

	// Exist reports whether the given account exists in state.
	// Notably this should also return true for suicided accounts.
	Exist(common.Address) bool
	// Empty returns whether the given account is empty. Empty
	// is defined according to EIP161 (balance = nonce = code = 0).
	Empty(common.Address) bool

	PrepareAccessList(sender common.Address, dest *common.Address, precompiles []common.Address, txAccesses types.AccessList)
	AddressInAccessList(addr common.Address) bool
	SlotInAccessList(addr common.Address, slot common.Hash) (addressOk bool, slotOk bool)
	// AddAddressToAccessList adds the given address to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddAddressToAccessList(addr common.Address)
	// AddSlotToAccessList adds the given (address,slot) to the access list. This operation is safe to perform
	// even if the feature/fork is not active yet
	AddSlotToAccessList(addr common.Address, slot common.Hash)

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) error

	// *** modify to support nft transaction 20211215 begin ***
	ChangeNFTOwner(common.Address, common.Address, int, *big.Int)
	GetNFTOwner(common.Address) common.Address
	GetNFTOwner16(common.Address) common.Address
	// *** modify to support nft transaction 20211215 end ***
	//CreateNFTByOfficial([]common.Address, *big.Int)
	CreateNFTByUser(common.Address, common.Address, uint16, string) (common.Address, bool)
	ChangeApproveAddress(common.Address, common.Address)
	CancelApproveAddress(common.Address, common.Address)
	ChangeNFTApproveAddress(common.Address, common.Address)
	CancelNFTApproveAddress(common.Address, common.Address)
	ExchangeNFTToCurrency(common.Address, common.Address, *big.Int, int)
	PledgeToken(common.Address, *big.Int, common.Address, *big.Int) error
	GetPledgedTime(common.Address) *big.Int
	MinerConsign(common.Address, common.Address) error
	CancelPledgedToken(common.Address, *big.Int)
	OpenExchanger(common.Address, *big.Int, *big.Int, uint16, string, string, common.Address)
	CloseExchanger(common.Address, *big.Int)
	GetExchangerFlag(common.Address) bool
	GetOpenExchangerTime(common.Address) *big.Int
	GetFeeRate(common.Address) uint16
	GetExchangerName(common.Address) string
	GetExchangerURL(common.Address) string
	GetApproveAddress(common.Address) []common.Address
	//GetNFTBalance(common.Address) uint64
	GetNFTName(common.Address) string
	GetNFTSymbol(common.Address) string
	//GetNFTApproveAddress(common.Address) []common.Address
	GetNFTApproveAddress(common.Address) common.Address
	GetNFTMergeLevel(common.Address) uint8
	GetNFTCreator(common.Address) common.Address
	GetNFTRoyalty(common.Address) uint16
	GetNFTExchanger(common.Address) common.Address
	GetNFTMetaURL(common.Address) string
	IsExistNFT(common.Address) bool
	IsApproved(common.Address, common.Address) bool
	IsApprovedOne(common.Address, common.Address) bool
	IsApprovedForAll(common.Address, common.Address) bool
	GetPledgedBalance(common.Address) *big.Int
	InjectOfficialNFT(string, *big.Int, uint64, uint16, string)
	AddExchangerToken(common.Address, *big.Int)
	ModifyOpenExchangerTime(common.Address, *big.Int)
	SubExchangerToken(common.Address, *big.Int)
	SubExchangerBalance(common.Address, *big.Int)
	GetExchangerBalance(common.Address) *big.Int
	VoteOfficialNFT(*types.NominatedOfficialNFT, *big.Int) error
	ElectNominatedOfficialNFT(*big.Int)
	SubVoteWeight(common.Address, *big.Int)
	AddVoteWeight(common.Address, *big.Int)
	AddValidatorCoefficient(common.Address, uint8)
	SubValidatorCoefficient(common.Address, uint8)
	GetValidatorCoefficient(common.Address) uint8
	NextIndex() *big.Int
	//PledgeNFT(common.Address, *big.Int)
	//CancelPledgedNFT(common.Address)
	GetMergeNumber(common.Address) uint32
	//GetPledgedFlag(common.Address) bool
	//GetNFTPledgedBlockNumber(common.Address) *big.Int
	CalculateExchangeAmount(uint8, uint32) *big.Int
	GetExchangAmount(common.Address, *big.Int) *big.Int
	IsOfficialNFT(common.Address) bool
	GetOfficialMint() *big.Int
	GetUserMint() *big.Int
	ChangeSNFTAgentRecipient(common.Address, common.Address)
	ChangeSNFTNoMerge(common.Address, bool)
	GetDividendAddrs(common.Address) []common.Address
	SetDividendAddrs(common.Address, []common.Address)
}

// CallContext provides a basic interface for the EVM calling conventions. The EVM
// depends on this context being implemented for doing subcalls and initialising new EVM contracts.
type CallContext interface {
	// Call another contract
	Call(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Take another's contract code and execute within our own context
	CallCode(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
	// Same as CallCode except sender and value is propagated from parent to child scope
	DelegateCall(env *EVM, me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error)
	// Create a new contract
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error)
}
