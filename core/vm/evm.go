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
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

//const CancelPledgedInterval = 365 * 720 * 24	// day * blockNumber of per hour * 24h
const CancelPledgedInterval = 3 * 24 // for test

type (
	// CanTransferFunc is the signature of a transfer guard function
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	// TransferFunc is the signature of a transfer function
	TransferFunc func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the n'th block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash

	// VerifyNFTOwnerFunc is to judge whether the owner own the nft
	VerifyNFTOwnerFunc func(StateDB, string, common.Address) bool
	// TransferNFTFunc is the signature of a TransferNFT function
	TransferNFTFunc             func(StateDB, string, common.Address) error
	CreateNFTByOfficialFunc     func(StateDB, []common.Address, *big.Int)
	CreateNFTByUserFunc         func(StateDB, common.Address, common.Address, uint32, string) (common.Address, bool)
	ChangeApproveAddressFunc    func(StateDB, common.Address, common.Address)
	CancelApproveAddressFunc    func(StateDB, common.Address, common.Address)
	ChangeNFTApproveAddressFunc func(StateDB, common.Address, common.Address)
	CancelNFTApproveAddressFunc func(StateDB, common.Address, common.Address)
	ExchangeNFTToCurrencyFunc   func(StateDB, common.Address, string, *big.Int) error
	PledgeTokenFunc             func(StateDB, common.Address, *big.Int, *types.Wormholes, *big.Int) error
	GetPledgedTimeFunc          func(StateDB, common.Address) *big.Int
	MinerConsignFunc            func(StateDB, common.Address, *types.Wormholes) error
	CancelPledgedTokenFunc      func(StateDB, common.Address, *big.Int)
	OpenExchangerFunc           func(StateDB, common.Address, *big.Int, *big.Int, uint32, string, string)
	CloseExchangerFunc          func(StateDB, common.Address, *big.Int)
	GetExchangerFlagFunc        func(StateDB, common.Address) bool
	GetOpenExchangerTimeFunc    func(StateDB, common.Address) *big.Int
	GetFeeRateFunc              func(StateDB, common.Address) uint32
	GetExchangerNameFunc        func(StateDB, common.Address) string
	GetExchangerURLFunc         func(StateDB, common.Address) string
	GetApproveAddressFunc       func(StateDB, common.Address) []common.Address
	GetNFTBalanceFunc           func(StateDB, common.Address) uint64
	GetNFTNameFunc              func(StateDB, common.Address) string
	GetNFTSymbolFunc            func(StateDB, common.Address) string
	//GetNFTApproveAddressFunc func(StateDB, common.Address) []common.Address
	GetNFTApproveAddressFunc               func(StateDB, common.Address) common.Address
	GetNFTMergeLevelFunc                   func(StateDB, common.Address) uint8
	GetNFTCreatorFunc                      func(StateDB, common.Address) common.Address
	GetNFTRoyaltyFunc                      func(StateDB, common.Address) uint32
	GetNFTExchangerFunc                    func(StateDB, common.Address) common.Address
	GetNFTMetaURLFunc                      func(StateDB, common.Address) string
	IsExistNFTFunc                         func(StateDB, common.Address) bool
	IsApprovedFunc                         func(StateDB, common.Address, common.Address) bool
	IsApprovedOneFunc                      func(StateDB, common.Address, common.Address) bool
	IsApprovedForAllFunc                   func(StateDB, common.Address, common.Address) bool
	VerifyPledgedBalanceFunc               func(StateDB, common.Address, *big.Int) bool
	InjectOfficialNFTFunc                  func(StateDB, string, *big.Int, uint64, uint32, string)
	BuyNFTBySellerOrExchangerFunc          func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByBuyerFunc                      func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByBuyerFunc               func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByExchangerFunc           func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByApproveExchangerFunc           func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByApprovedExchangerFunc   func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByExchangerFunc                  func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	AddExchangerTokenFunc                  func(StateDB, common.Address, *big.Int)
	SubExchangerTokenFunc                  func(StateDB, common.Address, *big.Int)
	SubExchangerBalanceFunc                func(StateDB, common.Address, *big.Int)
	VerifyExchangerBalanceFunc             func(StateDB, common.Address, *big.Int) bool
	GetNftAddressAndLevelFunc              func(string) (common.Address, int, error)
	VoteOfficialNFTFunc                    func(StateDB, *types.NominatedOfficialNFT)
	ElectNominatedOfficialNFTFunc          func(StateDB)
	NextIndexFunc                          func(db StateDB) *big.Int
	AddOrUpdateActiveMinerFunc             func(StateDB, common.Address, *big.Int, uint64)
	VoteOfficialNFTByApprovedExchangerFunc func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	ChangeRewardFlagFunc                   func(StateDB, common.Address, uint8)
)

func (evm *EVM) precompile(addr common.Address) (PrecompiledContract, bool) {
	var precompiles map[common.Address]PrecompiledContract
	switch {
	case evm.chainRules.IsBerlin:
		precompiles = PrecompiledContractsBerlin
	case evm.chainRules.IsIstanbul:
		precompiles = PrecompiledContractsIstanbul
	case evm.chainRules.IsByzantium:
		precompiles = PrecompiledContractsByzantium
	default:
		precompiles = PrecompiledContractsHomestead
	}
	p, ok := precompiles[addr]
	return p, ok
}

// BlockContext provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type BlockContext struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// *** modify to support nft transaction 20211215 begin ***
	// VerifyNFTOwner is to judge whether the owner own the nft
	VerifyNFTOwner VerifyNFTOwnerFunc
	// TransferNFT transfers NFT from one owner to the other
	TransferNFT TransferNFTFunc
	// *** modify to support nft transaction 20211215 end ***
	CreateNFTByOfficial                CreateNFTByOfficialFunc
	CreateNFTByUser                    CreateNFTByUserFunc
	ChangeApproveAddress               ChangeApproveAddressFunc
	CancelApproveAddress               CancelApproveAddressFunc
	ChangeNFTApproveAddress            ChangeNFTApproveAddressFunc
	CancelNFTApproveAddress            CancelNFTApproveAddressFunc
	ExchangeNFTToCurrency              ExchangeNFTToCurrencyFunc
	PledgeToken                        PledgeTokenFunc
	GetPledgedTime                     GetPledgedTimeFunc
	MinerConsign                       MinerConsignFunc
	CancelPledgedToken                 CancelPledgedTokenFunc
	OpenExchanger                      OpenExchangerFunc
	CloseExchanger                     CloseExchangerFunc
	GetExchangerFlag                   GetExchangerFlagFunc
	GetOpenExchangerTime               GetOpenExchangerTimeFunc
	GetFeeRate                         GetFeeRateFunc
	GetExchangerName                   GetExchangerNameFunc
	GetExchangerURL                    GetExchangerURLFunc
	GetApproveAddress                  GetApproveAddressFunc
	GetNFTBalance                      GetNFTBalanceFunc
	GetNFTName                         GetNFTNameFunc
	GetNFTSymbol                       GetNFTSymbolFunc
	GetNFTApproveAddress               GetNFTApproveAddressFunc
	GetNFTMergeLevel                   GetNFTMergeLevelFunc
	GetNFTCreator                      GetNFTCreatorFunc
	GetNFTRoyalty                      GetNFTRoyaltyFunc
	GetNFTExchanger                    GetNFTExchangerFunc
	GetNFTMetaURL                      GetNFTMetaURLFunc
	IsExistNFT                         IsExistNFTFunc
	IsApproved                         IsApprovedFunc
	IsApprovedOne                      IsApprovedOneFunc
	IsApprovedForAll                   IsApprovedForAllFunc
	VerifyPledgedBalance               VerifyPledgedBalanceFunc
	InjectOfficialNFT                  InjectOfficialNFTFunc
	BuyNFTBySellerOrExchanger          BuyNFTBySellerOrExchangerFunc
	BuyNFTByBuyer                      BuyNFTByBuyerFunc
	BuyAndMintNFTByBuyer               BuyAndMintNFTByBuyerFunc
	BuyAndMintNFTByExchanger           BuyAndMintNFTByExchangerFunc
	BuyNFTByApproveExchanger           BuyNFTByApproveExchangerFunc
	BuyAndMintNFTByApprovedExchanger   BuyAndMintNFTByApprovedExchangerFunc
	BuyNFTByExchanger                  BuyNFTByExchangerFunc
	AddExchangerToken                  AddExchangerTokenFunc
	SubExchangerToken                  SubExchangerTokenFunc
	SubExchangerBalance                SubExchangerBalanceFunc
	VerifyExchangerBalance             VerifyExchangerBalanceFunc
	GetNftAddressAndLevel              GetNftAddressAndLevelFunc
	VoteOfficialNFT                    VoteOfficialNFTFunc
	ElectNominatedOfficialNFT          ElectNominatedOfficialNFTFunc
	NextIndex                          NextIndexFunc
	AddOrUpdateActiveMiner             AddOrUpdateActiveMinerFunc
	VoteOfficialNFTByApprovedExchanger VoteOfficialNFTByApprovedExchangerFunc
	ChangeRewardFlag                   ChangeRewardFlagFunc

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    uint64         // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
	BaseFee     *big.Int       // Provides information for BASEFEE
}

// TxContext provides the EVM with information about a transaction.
// All fields can change between transactions.
type TxContext struct {
	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE
}

// EVM is the Ethereum Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context BlockContext
	TxContext
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	Config Config
	// global (to this context) ethereum virtual machine
	// used throughout the execution of the tx.
	interpreter *EVMInterpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	// callGasTemp holds the gas available for the current call. This is needed because the
	// available gas is calculated in gasCall* according to the 63/64 rule and later
	// applied in opCall*.
	callGasTemp uint64
}

// *** modify to support nft transaction 20211215 begin ***

// NewEVM returns a new EVM. The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(blockCtx BlockContext, txCtx TxContext, statedb StateDB, chainConfig *params.ChainConfig, config Config) *EVM {
	evm := &EVM{
		Context:     blockCtx,
		TxContext:   txCtx,
		StateDB:     statedb,
		Config:      config,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(blockCtx.BlockNumber),
	}
	evm.interpreter = NewEVMInterpreter(evm, config)
	return evm
}

// Reset resets the EVM with a new transaction context.Reset
// This is not threadsafe and should only be done very cautiously.
func (evm *EVM) Reset(txCtx TxContext, statedb StateDB) {
	evm.TxContext = txCtx
	evm.StateDB = statedb
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// Cancelled returns true if Cancel has been called
func (evm *EVM) Cancelled() bool {
	return atomic.LoadInt32(&evm.abort) == 1
}

// Interpreter returns the current interpreter
func (evm *EVM) Interpreter() *EVMInterpreter {
	return evm.interpreter
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	var nftTransaction bool = false
	var wormholes types.Wormholes
	if evm.Config.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	snapshot := evm.StateDB.Snapshot()
	p, isPrecompile := evm.precompile(addr)

	if !evm.StateDB.Exist(addr) && nftTransaction {
		if !isPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.Config.Debug && evm.depth == 0 {
				evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
				evm.Config.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, gas, nil
		}
		evm.StateDB.CreateAccount(addr)
	}
	//fmt.Println("input=", string(input))
	//fmt.Println("caller.Address=", caller.Address().String())
	// *** modify to support nft transaction 20211215 begin ***
	if len(input) > 10 {
		if string(input[:10]) == "wormholes:" {
			jsonErr := json.Unmarshal(input[10:], &wormholes)
			if jsonErr == nil {
				nftTransaction = true
			} else {
				log.Error("EVM.Call(), wormholes unmarshal error", "jsonErr", jsonErr,
					"wormholes", string(input))
				return nil, gas, ErrWormholesFormat
			}
		}
	}
	log.Info("EVM.Call()", "nftTransaction", nftTransaction)
	if nftTransaction {
		log.Info("EVM.Call()", "nftTransaction", nftTransaction, "wormholes.Type", wormholes.Type)
		ret, gas, err = evm.HandleNFT(caller, addr, wormholes, gas, value)
		if err != nil {
			return ret, gas, err
		}
	} else {
		evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
	}
	// *** modify to support nft transaction 20211215 end ***

	// Capture the tracer start/end events in debug mode
	if evm.Config.Debug && evm.depth == 0 {
		evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
		defer func(startGas uint64, startTime time.Time) { // Lazy evaluation of the parameters
			evm.Config.Tracer.CaptureEnd(ret, startGas-gas, time.Since(startTime), err)
		}(gas, time.Now())
	}

	if isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		code := evm.StateDB.GetCode(addr)
		if len(code) == 0 {
			ret, err = nil, nil // gas is unchanged
		} else {
			addrCopy := addr
			// If the account has no code, we can abort here
			// The depth-check is already done, and precompiles handled above
			contract := NewContract(caller, AccountRef(addrCopy), value, gas)
			contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), code)
			ret, err = evm.interpreter.Run(contract, input, false)
			gas = contract.Gas
		}
	}
	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
		// TODO: consider clearing up unused snapshots:
		//} else {
		//	evm.StateDB.DiscardSnapshot(snapshot)
	}
	return ret, gas, err
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	if evm.Config.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	// Note although it's noop to transfer X ether to caller itself. But
	// if caller doesn't have enough balance, it would be an error to allow
	// over-charging itself. So the check here is necessary.
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}
	var snapshot = evm.StateDB.Snapshot()

	// It is allowed to call precompiles, even via delegatecall
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract(caller, AccountRef(caller.Address()), value, gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = evm.interpreter.Run(contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.Config.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	var snapshot = evm.StateDB.Snapshot()

	// It is allowed to call precompiles, even via delegatecall
	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		addrCopy := addr
		// Initialise a new contract and make initialise the delegate values
		contract := NewContract(caller, AccountRef(caller.Address()), nil, gas).AsDelegate()
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		ret, err = evm.interpreter.Run(contract, input, false)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error) {
	if evm.Config.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// We take a snapshot here. This is a bit counter-intuitive, and could probably be skipped.
	// However, even a staticcall is considered a 'touch'. On mainnet, static calls were introduced
	// after all empty accounts were deleted, so this is not required. However, if we omit this,
	// then certain tests start failing; stRevertTest/RevertPrecompiledTouchExactOOG.json.
	// We could change this, but for now it's left for legacy reasons
	var snapshot = evm.StateDB.Snapshot()

	// We do an AddBalance of zero here, just in order to trigger a touch.
	// This doesn't matter on Mainnet, where all empties are gone at the time of Byzantium,
	// but is the correct thing to do and matters on other networks, in tests, and potential
	// future scenarios
	evm.StateDB.AddBalance(addr, big0)

	if p, isPrecompile := evm.precompile(addr); isPrecompile {
		ret, gas, err = RunPrecompiledContract(p, input, gas)
	} else {
		// At this point, we use a copy of address. If we don't, the go compiler will
		// leak the 'contract' to the outer scope, and make allocation for 'contract'
		// even if the actual execution ends on RunPrecompiled above.
		addrCopy := addr
		// Initialise a new contract and set the code that is to be used by the EVM.
		// The contract is a scoped environment for this execution context only.
		contract := NewContract(caller, AccountRef(addrCopy), new(big.Int), gas)
		contract.SetCallCode(&addrCopy, evm.StateDB.GetCodeHash(addrCopy), evm.StateDB.GetCode(addrCopy))
		// When an error was returned by the EVM or when setting the creation code
		// above we revert to the snapshot and consume any gas remaining. Additionally
		// when we're in Homestead this also counts for code storage gas errors.
		ret, err = evm.interpreter.Run(contract, input, true)
		gas = contract.Gas
	}
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			gas = 0
		}
	}
	return ret, gas, err
}

type codeAndHash struct {
	code []byte
	hash common.Hash
}

func (c *codeAndHash) Hash() common.Hash {
	if c.hash == (common.Hash{}) {
		c.hash = crypto.Keccak256Hash(c.code)
	}
	return c.hash
}

// create creates a new contract using code as deployment code.
func (evm *EVM) create(caller ContractRef, codeAndHash *codeAndHash, gas uint64, value *big.Int, address common.Address) ([]byte, common.Address, uint64, error) {
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gas, ErrDepth
	}
	if !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.Address{}, gas, ErrInsufficientBalance
	}
	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)
	// We add this to the access list _before_ taking a snapshot. Even if the creation fails,
	// the access-list change should not be rolled back
	if evm.chainRules.IsBerlin {
		evm.StateDB.AddAddressToAccessList(address)
	}
	// Ensure there's no existing contract already at the designated address
	contractHash := evm.StateDB.GetCodeHash(address)
	if evm.StateDB.GetNonce(address) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(address)
	if evm.chainRules.IsEIP158 {
		evm.StateDB.SetNonce(address, 1)
	}
	evm.Context.Transfer(evm.StateDB, caller.Address(), address, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := NewContract(caller, AccountRef(address), value, gas)
	contract.SetCodeOptionalHash(&address, codeAndHash)

	if evm.Config.NoRecursion && evm.depth > 0 {
		return nil, address, gas, nil
	}

	if evm.Config.Debug && evm.depth == 0 {
		evm.Config.Tracer.CaptureStart(evm, caller.Address(), address, true, codeAndHash.code, gas, value)
	}
	start := time.Now()

	ret, err := evm.interpreter.Run(contract, nil, false)

	// Check whether the max code size has been exceeded, assign err if the case.
	if err == nil && evm.chainRules.IsEIP158 && len(ret) > params.MaxCodeSize {
		err = ErrMaxCodeSizeExceeded
	}

	// Reject code starting with 0xEF if EIP-3541 is enabled.
	if err == nil && len(ret) >= 1 && ret[0] == 0xEF && evm.chainRules.IsLondon {
		err = ErrInvalidCode
	}

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
		if contract.UseGas(createDataGas) {
			evm.StateDB.SetCode(address, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if err != nil && (evm.chainRules.IsHomestead || err != ErrCodeStoreOutOfGas) {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != ErrExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	if evm.Config.Debug && evm.depth == 0 {
		evm.Config.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	}
	return ret, address, contract.Gas, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gas uint64, value *big.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	contractAddr = crypto.CreateAddress(caller.Address(), evm.StateDB.GetNonce(caller.Address()))
	return evm.create(caller, &codeAndHash{code: code}, gas, value, contractAddr)
}

// Create2 creates a new contract using code as deployment code.
//
// The different between Create2 with Create is Create2 uses sha3(0xff ++ msg.sender ++ salt ++ sha3(init_code))[12:]
// instead of the usual sender-and-nonce-hash as the address where the contract is initialized at.
func (evm *EVM) Create2(caller ContractRef, code []byte, gas uint64, endowment *big.Int, salt *uint256.Int) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	codeAndHash := &codeAndHash{code: code}
	contractAddr = crypto.CreateAddress2(caller.Address(), salt.Bytes32(), codeAndHash.Hash().Bytes())
	return evm.create(caller, codeAndHash, gas, endowment, contractAddr)
}

// ChainConfig returns the environment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig { return evm.chainConfig }

func (evm *EVM) HandleNFT(
	caller ContractRef,
	addr common.Address,
	wormholes types.Wormholes,
	gas uint64,
	value *big.Int) (ret []byte, leftOverGas uint64, err error) {

	formatErr := wormholes.CheckFormat()
	if formatErr != nil {
		log.Error("HandleNFT() format error", "wormholes.Type", wormholes.Type, "error", formatErr)
		return nil, gas, formatErr
	}

	switch wormholes.Type {
	case 0: // create nft by user
		log.Info("HandleNFT(), CreateNFTByUser>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if wormholes.Royalty <= 0 {
			log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotMoreThan0)
			return nil, gas, ErrRoyaltyNotMoreThan0
		}
		if wormholes.Royalty >= 10000 {
			log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotLessthan10000)
			return nil, gas, ErrRoyaltyNotLessthan10000
		}

		exchanger := common.Address{}
		if len(wormholes.Exchanger) > 2 {
			if !strings.HasPrefix(wormholes.Exchanger, "0x") &&
				!strings.HasPrefix(wormholes.Exchanger, "0X") {
				log.Error("HandleNFT(), CreateNFTByUser(), exchanger format error",
					"wormholes.Exchanger", wormholes.Exchanger)
				return nil, gas, ErrExchangerFormat
			}
			exchanger = common.HexToAddress(wormholes.Exchanger)

			exchangerFlag := evm.Context.GetExchangerFlag(evm.StateDB, exchanger)
			if exchangerFlag != true {
				log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type, "error", ErrNotExchanger)
				return nil, gas, ErrNotExchanger
			}
		}

		evm.Context.CreateNFTByUser(evm.StateDB,
			exchanger,
			addr,
			wormholes.Royalty,
			wormholes.MetaURL)
		log.Info("HandleNFT(), CreateNFTByUser<<<<<<<<<<", "wormholes.Type", wormholes.Type)

	case 1: //transfer nft
		nftAddress, _, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
		if evm.Context.VerifyNFTOwner(evm.StateDB, wormholes.NFTAddress, caller.Address()) {
			if !evm.StateDB.Exist(nftAddress) {
				evm.StateDB.CreateAccount(nftAddress)
			}
			log.Info("HandleNFT(), TransferNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type)
			err := evm.Context.TransferNFT(evm.StateDB, wormholes.NFTAddress, addr)
			if err != nil {
				log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type, "error", err)
				return nil, gas, err
			}
			log.Info("HandleNFT(), TransferNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		} else {
			log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type, "error", ErrNotOwner)
			return nil, gas, ErrNotOwner
		}

	case 2: //approve a nft's authority
		log.Info("HandleNFT(), ChangeNFTApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if !evm.Context.GetExchangerFlag(evm.StateDB, addr) {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "error", ErrNotExchanger)
			return nil, gas, ErrNotExchanger
		}
		nftAddress, _, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
		if IsOfficialNFT(nftAddress) {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress, "error", ErrNotAllowedOfficialNFT)
			return nil, gas, ErrNotAllowedOfficialNFT
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		log.Info("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type,
			"wormholes.NFTAddress", wormholes.NFTAddress, "approvedaddress", addr.String(), "nftOwner", nftOwner.String())
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "error", ErrNotOwner)
			return nil, gas, ErrNotOwner
		}
		evm.Context.ChangeNFTApproveAddress(
			evm.StateDB,
			nftAddress,
			addr)
		log.Info("HandleNFT(), ChangeNFTApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 3:
		log.Info("HandleNFT(), CancelNFTApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		nftAddress, _, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			log.Error("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		log.Info("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type,
			"wormholes.NFTAddress", wormholes.NFTAddress, "approvedaddress", addr.String(), "nftOwner", nftOwner.String())
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type, "error", ErrNotOwner)
			return nil, gas, ErrNotOwner
		}
		evm.Context.CancelNFTApproveAddress(
			evm.StateDB,
			nftAddress,
			addr)
		log.Info("HandleNFT(), CancelNFTApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 4: //approve all nft's authority
		if !evm.Context.GetExchangerFlag(evm.StateDB, addr) {
			log.Error("HandleNFT(), ChangeApproveAddress", "wormholes.Type", wormholes.Type, "error", ErrNotExchanger)
			return nil, gas, ErrNotExchanger
		}
		log.Info("HandleNFT(), ChangeApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		evm.Context.ChangeApproveAddress(
			evm.StateDB,
			caller.Address(),
			addr)
		log.Info("HandleNFT(), ChangeApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 5:
		log.Info("HandleNFT(), CancelApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		evm.Context.CancelApproveAddress(
			evm.StateDB,
			caller.Address(),
			addr)
		log.Info("HandleNFT(), CancelApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 6: //NFT exchange
		log.Info("HandleNFT(), ExchangeNFTToCurrency>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			return nil, gas, err
		}
		if !IsOfficialNFT(nftAddress) {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress, "error", ErrNotMintByOfficial)
			return nil, gas, ErrNotMintByOfficial
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"nft owner", nftOwner, "error", ErrNotOwner)
			return nil, gas, ErrNotOwner
		}
		level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		if int(level2) < level1 {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"input nft level", level1, "real nft level", level2, "error", ErrNotOwner)
			return nil, gas, ErrNftLevel
		}
		evm.Context.ExchangeNFTToCurrency(
			evm.StateDB,
			caller.Address(),
			wormholes.NFTAddress,
			evm.Context.BlockNumber)
		log.Info("HandleNFT(), ExchangeNFTToCurrency<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 7: //NFT pledge

	case 8: //cancel nft pledge

	case 9: // pledge token
		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		Erb100000 := big.NewInt(100000)
		Erb100000.Mul(Erb100000, baseErb)
		if !evm.Context.VerifyPledgedBalance(evm.StateDB, caller.Address(), Erb100000) {
			//if this account has not pledged
			if value.Cmp(Erb100000) < 0 {
				log.Error("HandleNFT(), PledgeToken", "wormholes.Type", wormholes.Type, "error", ErrNotMoreThan100000ERB)
				return nil, gas, ErrNotMoreThan100000ERB
			}
		}

		currentBlockNumber := new(big.Int).Set(evm.Context.BlockNumber)
		pledgedBalance := evm.StateDB.GetPledgedBalance(caller.Address())
		// if append pledgebalance, reset pledgedblocknumber
		if pledgedBalance.Cmp(big.NewInt(0)) > 0 {
			pledgedBlockNumber := evm.StateDB.GetPledgedTime(caller.Address())
			height, err := UnstakingHeight(pledgedBalance, value, pledgedBlockNumber.Uint64(), currentBlockNumber.Uint64(), CancelPledgedInterval)
			if err != nil {
				return nil, gas, err
			}
			bigHeight := new(big.Int).SetUint64(height)
			bigCancelPledgedInterval := new(big.Int).SetUint64(CancelPledgedInterval)
			currentBlockNumber = new(big.Int).Add(currentBlockNumber, bigHeight)
			currentBlockNumber = new(big.Int).Sub(currentBlockNumber, bigCancelPledgedInterval)
		}

		log.Info("HandleNFT()", "PledgeToken.req", wormholes)
		if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			log.Info("HandleNFT(), PledgeToken>>>>>>>>>>", "wormholes.Type", wormholes.Type)
			err := evm.Context.PledgeToken(evm.StateDB, caller.Address(), value, &wormholes, currentBlockNumber)
			if err != nil {
				log.Info("HandleNFT(), PledgeToken<<<<<<<<<<", "wormholes.Type", wormholes.Type)
				return nil, gas, err
			}
			log.Info("HandleNFT(), PledgeToken<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		} else {
			log.Error("HandleNFT(), PledgeToken", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
			return nil, gas, ErrInsufficientBalance
		}

	case 10: // cancel pledge of token
		log.Info("HandleNFT(), CancelPledgedToken>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		pledgedTime := evm.Context.GetPledgedTime(evm.StateDB, caller.Address())
		if big.NewInt(CancelPledgedInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, pledgedTime)) > 0 {
			log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type, "error", ErrTooCloseToCancel)
			return nil, gas, ErrTooCloseToCancel
		}

		pledgedBalance := evm.StateDB.GetPledgedBalance(caller.Address())
		if pledgedBalance.Cmp(value) == 0 {
			// cancel pledged balance
			log.Info("HandleNFT(), CancelPledgedToken, cancel all", "wormholes.Type", wormholes.Type)
			evm.Context.CancelPledgedToken(evm.StateDB, caller.Address(), value)

		} else {
			// cancel partial pledged balance
			baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
			Erb100000 := big.NewInt(100000)
			Erb100000.Mul(Erb100000, baseErb)

			if evm.Context.VerifyPledgedBalance(evm.StateDB, caller.Address(), new(big.Int).Add(Erb100000, value)) {
				log.Info("HandleNFT(), CancelPledgedToken, cancel partial", "wormholes.Type", wormholes.Type)
				evm.Context.CancelPledgedToken(evm.StateDB, caller.Address(), value)
			} else {
				log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type, "error", ErrInsufficientPledgedBalance)
				return nil, gas, ErrInsufficientPledgedBalance
			}
		}
		log.Info("HandleNFT(), CancelPledgedToken<<<<<<<<<<", "wormholes.Type", wormholes.Type)

	case 11: //open exchanger
		log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		// value must be greater than or equal to 100 ERB
		unitErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		if value.Cmp(new(big.Int).Mul(big.NewInt(100), unitErb)) < 0 {
			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrNotMoreThan100ERB)
			return nil, gas, ErrNotMoreThan100ERB
		}
		if wormholes.FeeRate <= 0 {
			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrFeeRateNotMoreThan0)
			return nil, gas, ErrFeeRateNotMoreThan0
		}
		if wormholes.FeeRate >= 10000 {
			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrFeeRateNotLessThan10000)
			return nil, gas, ErrFeeRateNotLessThan10000
		}

		exchangerFlag := evm.Context.GetExchangerFlag(evm.StateDB, addr)
		if caller.Address() == addr {
			if exchangerFlag == true {
				log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrReopenExchanger)
				return nil, gas, ErrReopenExchanger
			}

			if evm.Context.CanTransfer(evm.StateDB, addr, value) {
				log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
				evm.Context.OpenExchanger(
					evm.StateDB,
					addr,
					value,
					evm.Context.BlockNumber,
					wormholes.FeeRate,
					wormholes.Name,
					wormholes.Url)
				log.Info("HandleNFT(), OpenExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)

			} else {
				log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
				return nil, gas, ErrInsufficientBalance
			}

		} else {
			// caller give addr a exchanger as a gift, if addr has been a exchanger, add value to exchangerbalance.
			if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
				evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
			} else {
				log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
				return nil, gas, ErrInsufficientBalance
			}

			if exchangerFlag == true {
				evm.Context.AddExchangerToken(evm.StateDB, addr, value)
			} else {
				log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
				evm.Context.OpenExchanger(
					evm.StateDB,
					addr,
					value,
					evm.Context.BlockNumber,
					wormholes.FeeRate,
					wormholes.Name,
					wormholes.Url)
				log.Info("HandleNFT(), OpenExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
			}
		}

	case 12: //close exchanger
		log.Info("HandleNFT(), CloseExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		//const CloseExchangerInterval = 180 * 720 * 24	// day * blockNumber of per hour * 24h
		const CloseExchangerInterval = 3 * 24 // for test
		openExchangerTime := evm.Context.GetOpenExchangerTime(evm.StateDB, caller.Address())
		if big.NewInt(CloseExchangerInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, openExchangerTime)) > 0 {
			log.Error("HandleNFT(), CloseExchanger", "wormholes.Type", wormholes.Type, "error", ErrTooCloseWithOpenExchanger)
			return nil, gas, ErrTooCloseWithOpenExchanger
		}
		evm.Context.CloseExchanger(evm.StateDB, caller.Address(), evm.Context.BlockNumber)
		log.Info("HandleNFT(), CloseExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		//evm.StateDB.CloseExchanger(caller.Address(), evm.Context.BlockNumber)
	//case 13:
	//	if !strings.HasPrefix(wormholes.StartIndex, "0x") &&
	//		!strings.HasPrefix(wormholes.StartIndex, "0X") {
	//		return nil, gas, ErrStartIndex
	//	}
	//	startIndex, _ := new(big.Int).SetString(wormholes.StartIndex[2:], 16)
	//	log.Info("HandleNFT(), InjectOfficialNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type)
	//	evm.Context.InjectOfficialNFT(evm.StateDB,
	//		wormholes.Dir,
	//		startIndex,
	//		wormholes.Number,
	//		wormholes.Royalty,
	//		wormholes.Creator)
	//	log.Info("HandleNFT(), InjectOfficialNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 14:
		log.Info("HandleNFT(), BuyNFTBySellerOrExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyNFTBySellerOrExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyNFTBySellerOrExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyNFTBySellerOrExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 15:
		log.Info("HandleNFT(), BuyNFTByBuyer>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyNFTByBuyer(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyNFTByBuyer<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByBuyer", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 16:
		log.Info("HandleNFT(), BuyAndMintNFTByBuyer>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyAndMintNFTByBuyer(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyAndMintNFTByBuyer<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByBuyer", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 17:
		log.Info("HandleNFT(), BuyAndMintNFTByExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyAndMintNFTByExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyAndMintNFTByExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 18:
		log.Info("HandleNFT(), BuyNFTByApproveExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyNFTByApproveExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyNFTByApproveExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByApproveExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 19:
		log.Info("HandleNFT(), BuyAndMintNFTByApprovedExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyAndMintNFTByApprovedExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyAndMintNFTByApprovedExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByApprovedExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 20:
		log.Info("HandleNFT(), BuyNFTByExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BuyNFTByExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		log.Info("HandleNFT(), BuyNFTByExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}
	case 21:
		if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			log.Info("HandleNFT(), AddExchangerToken>>>>>>>>>>", "wormholes.Type", wormholes.Type)
			evm.Context.AddExchangerToken(evm.StateDB, caller.Address(), value)
			log.Info("HandleNFT(), AddExchangerToken<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		} else {
			log.Error("HandleNFT(), AddExchangerToken", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
			return nil, gas, ErrInsufficientBalance
		}
	case 22:
		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		Erb100 := big.NewInt(100)
		Erb100.Mul(Erb100, baseErb)
		if evm.Context.VerifyExchangerBalance(evm.StateDB, caller.Address(), new(big.Int).Add(value, Erb100)) {
			log.Info("HandleNFT(), SubExchangerToken>>>>>>>>>>", "wormholes.Type", wormholes.Type)
			evm.Context.SubExchangerToken(evm.StateDB, caller.Address(), value)
			log.Info("HandleNFT(), SubExchangerToken<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		} else {
			log.Error("HandleNFT(), SubExchangerToken", "wormholes.Type", wormholes.Type, "error", ErrInsufficientExchangerBalance)
			return nil, gas, ErrInsufficientExchangerBalance
		}
	case 30:
		log.Info("HandleNFT(), SendLivenessTx>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			log.Info("HandleNFT(), Start|LivenessTx>>>>>>>>>>", "wormholes.Type", wormholes.Type)
			// Online transaction execution transfer
			evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
			// Add the online address to the active Miners Pool after the online transaction executes the transfer
			evm.Context.AddOrUpdateActiveMiner(evm.StateDB, caller.Address(), value, evm.Context.BlockNumber.Uint64())
			log.Info("HandleNFT(), End|LivenessTx<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		} else {
			log.Error("HandleNFT(), SendLivenessTx error", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
			return nil, gas, ErrInsufficientBalance
		}
	case 31:
		//MinerConsign
		log.Info("HandleNFT()", "MinerConsign.req", wormholes)
		//if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		log.Info("HandleNFT(), Start|MinerConsign>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		err := evm.Context.MinerConsign(evm.StateDB, caller.Address(), &wormholes)
		if err != nil {
			log.Info("HandleNFT(), End|MinerConsign<<<<<<<<<<", "wormholes.Type", wormholes.Type)
			return nil, gas, err
		}
		log.Info("HandleNFT(), End|MinerConsign<<<<<<<<<<", "wormholes.Type", wormholes.Type)
		//} else {
		//	log.Error("HandleNFT(), MinerConsign error", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
		//	return nil, gas, ErrInsufficientBalance
		//}
	case 23:
		log.Info("HandleNFT(), VoteOfficialNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		//if !strings.HasPrefix(wormholes.StartIndex, "0x") &&
		//	!strings.HasPrefix(wormholes.StartIndex, "0X") {
		//	return nil, gas, ErrStartIndex
		//}
		//startIndex, _ := new(big.Int).SetString(wormholes.StartIndex[2:], 16)
		startIndex := evm.StateDB.NextIndex()
		var number uint64 = 65536
		if wormholes.Royalty <= 0 {
			log.Error("HandleNFT(), VoteOfficialNFT", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotMoreThan0)
			return nil, gas, ErrRoyaltyNotMoreThan0
		}
		if wormholes.Royalty >= 10000 {
			log.Error("HandleNFT(), VoteOfficialNFT", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotLessthan10000)
			return nil, gas, ErrRoyaltyNotLessthan10000
		}

		nominatedNFT := types.NominatedOfficialNFT{
			InjectedOfficialNFT: types.InjectedOfficialNFT{
				Dir:        wormholes.Dir,
				StartIndex: startIndex,
				//Number: wormholes.Number,
				Number:  number,
				Royalty: wormholes.Royalty,
				Creator: wormholes.Creator,
				Address: caller.Address(),
			},
		}
		evm.Context.VoteOfficialNFT(evm.StateDB, &nominatedNFT)
		log.Info("HandleNFT(), VoteOfficialNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type)

	case 24:
		log.Info("HandleNFT(), VoteOfficialNFTByApprovedExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type)

		if wormholes.Royalty <= 0 {
			log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotMoreThan0)
			return nil, gas, ErrRoyaltyNotMoreThan0
		}
		if wormholes.Royalty >= 10000 {
			log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type, "error", ErrRoyaltyNotLessthan10000)
			return nil, gas, ErrRoyaltyNotLessthan10000
		}

		err := evm.Context.VoteOfficialNFTByApprovedExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		if err != nil {
			log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type, "error", err)
			return nil, gas, err
		}

		log.Info("HandleNFT(), VoteOfficialNFTByApprovedExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type)
	case 25:
		log.Info("HandleNFT(), ChangeRewardFlag>>>>>>>>>>", "wormholes.Type", wormholes.Type)
		evm.Context.ChangeRewardFlag(
			evm.StateDB,
			caller.Address(),
			wormholes.RewardFlag)
		log.Info("HandleNFT(), ChangeRewardFlag<<<<<<<<<<", "wormholes.Type", wormholes.Type)

	default:
		log.Error("HandleNFT()", "wormholes.Type", wormholes.Type, "error", ErrNotExistNFTType)
		return nil, gas, ErrNotExistNFTType
	}

	return nil, gas, nil
}

// IsOfficialNFT return true if nft address is created by official
func IsOfficialNFT(nftAddress common.Address) bool {
	maskByte := byte(128)
	nftByte := nftAddress[0]
	result := maskByte & nftByte
	if result == 128 {
		return true
	}
	return false
}

// UnstakingHeight @title    UnstakingHeight
// @description   UnstakingHeight Returns the height at which stakers can get their stake back
// @auth      mindcarver        2022/08/01
// @param     stakedAmt        *big.Int   the total amount of the current pledge
// @param     appendAmt        *big.Int   additional amount
// @param     sno        		uint64    starting height
// @param     cno       	    uint64    current height
// @param     lockedNo          uint64    lock time (counted in blocks)
// @return                      uint64     delay the amount of time (in blocks) that can be unstakes
func UnstakingHeight(stakedAmt, appendAmt *big.Int, sno, cno, lockedNo uint64) (uint64, error) {
	_, err := checkParams(stakedAmt, appendAmt, sno, cno)
	if err != nil {
		return 0, err
	}
	return unstakingHeight(stakedAmt, appendAmt, sno, cno, lockedNo), nil
}

func unstakingHeight(stakedAmt *big.Int, appendAmt *big.Int, sno uint64, cno uint64, lockedNo uint64) uint64 {
	var curRemainingTime uint64

	if sno+lockedNo > cno {
		curRemainingTime = sno + lockedNo - cno
	}

	total := big.NewFloat(0).Add(new(big.Float).SetInt(stakedAmt), new(big.Float).SetInt(appendAmt))
	h1 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(stakedAmt), total), new(big.Float).SetInt(big.NewInt(int64(curRemainingTime))))
	h2 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(appendAmt), total), new(big.Float).SetInt(big.NewInt(int64(lockedNo))))
	delayHeight, _ := big.NewFloat(0).Add(h1, h2).Uint64()
	return delayHeight
}

func checkParams(stakedAmt *big.Int, appendAmt *big.Int, sno uint64, cno uint64) (bool, error) {
	if stakedAmt.Cmp(big.NewInt(0)) == 0 || appendAmt.Cmp(big.NewInt(0)) == 0 {
		return false, errors.New("illegal amount")
	}
	if sno == 0 || cno == 0 {
		return false, errors.New("illegal height")
	}
	if sno > cno {
		return false, errors.New("the current height is less than the starting height of the pledge")
	}
	return true, nil
}
