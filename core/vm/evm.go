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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

// emptyCodeHash is used by create to ensure deployment is disallowed to already
// deployed contract addresses (relevant after the account abstraction).
var emptyCodeHash = crypto.Keccak256Hash(nil)

// const CancelDayPledgedInterval = 720 * 24 // blockNumber of per hour * 24h
const CancelDayPledgedInterval = 5 // blockNumber of per hour * 24h
// const CancelPledgedInterval = 365 * 720 * 24 // day * blockNumber of per hour * 24h
const CancelPledgedInterval = 3 * 24 // for test
// const CloseExchangerInterval = 365 * 720 * 24 // day * blockNumber of per hour * 24h
const CloseExchangerInterval = 3 * 24 // for test
// const CancelNFTPledgedInterval = 365 * 720 * 24 // day * blockNumber of per hour * 24h
const CancelNFTPledgedInterval = 3 * 24 // for test
const VALIDATOR_COEFFICIENT = 70

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
	TransferNFTFunc func(StateDB, string, common.Address, *big.Int) error
	//CreateNFTByOfficialFunc     func(StateDB, []common.Address, *big.Int)
	CreateNFTByUserFunc         func(StateDB, common.Address, common.Address, uint16, string) (common.Address, bool)
	ChangeApproveAddressFunc    func(StateDB, common.Address, common.Address)
	CancelApproveAddressFunc    func(StateDB, common.Address, common.Address)
	ChangeNFTApproveAddressFunc func(StateDB, common.Address, common.Address)
	CancelNFTApproveAddressFunc func(StateDB, common.Address, common.Address)
	ExchangeNFTToCurrencyFunc   func(StateDB, common.Address, string, *big.Int) error
	PledgeTokenFunc             func(StateDB, common.Address, *big.Int, *types.Wormholes, *big.Int) error
	StakerPledgeFunc            func(StateDB, common.Address, common.Address, *big.Int, *big.Int, *types.Wormholes) error
	GetPledgedTimeFunc          func(StateDB, common.Address, common.Address) *big.Int
	GetStakerPledgedFunc        func(StateDB, common.Address, common.Address) *types.StakerExtension
	MinerConsignFunc            func(StateDB, common.Address, *types.Wormholes) error
	MinerBecomeFunc             func(StateDB, common.Address, *types.Wormholes) error
	CancelPledgedTokenFunc      func(StateDB, common.Address, *big.Int)
	CancelStakerPledgeFunc      func(StateDB, common.Address, common.Address, *big.Int, *big.Int)
	OpenExchangerFunc           func(StateDB, common.Address, *big.Int, *big.Int, uint16, string, string, string)
	CloseExchangerFunc          func(StateDB, common.Address, *big.Int)
	GetExchangerFlagFunc        func(StateDB, common.Address) bool
	GetOpenExchangerTimeFunc    func(StateDB, common.Address) *big.Int
	GetFeeRateFunc              func(StateDB, common.Address) uint16
	GetExchangerNameFunc        func(StateDB, common.Address) string
	GetExchangerURLFunc         func(StateDB, common.Address) string
	GetApproveAddressFunc       func(StateDB, common.Address) []common.Address
	//GetNFTBalanceFunc           func(StateDB, common.Address) uint64
	GetNFTNameFunc   func(StateDB, common.Address) string
	GetNFTSymbolFunc func(StateDB, common.Address) string
	//GetNFTApproveAddressFunc func(StateDB, common.Address) []common.Address
	GetNFTApproveAddressFunc               func(StateDB, common.Address) common.Address
	GetNFTMergeLevelFunc                   func(StateDB, common.Address) uint8
	GetNFTCreatorFunc                      func(StateDB, common.Address) common.Address
	GetNFTRoyaltyFunc                      func(StateDB, common.Address) uint16
	GetNFTExchangerFunc                    func(StateDB, common.Address) common.Address
	GetNFTMetaURLFunc                      func(StateDB, common.Address) string
	IsExistNFTFunc                         func(StateDB, common.Address) bool
	IsApprovedFunc                         func(StateDB, common.Address, common.Address) bool
	IsApprovedOneFunc                      func(StateDB, common.Address, common.Address) bool
	IsApprovedForAllFunc                   func(StateDB, common.Address, common.Address) bool
	VerifyPledgedBalanceFunc               func(StateDB, common.Address, *big.Int) bool
	VerifyStakerPledgedBalanceFunc         func(StateDB, common.Address, common.Address, *big.Int) bool
	InjectOfficialNFTFunc                  func(StateDB, string, *big.Int, uint64, uint16, string)
	BuyNFTBySellerOrExchangerFunc          func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByBuyerFunc                      func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByBuyerFunc               func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByExchangerFunc           func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByApproveExchangerFunc           func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BatchBuyNFTByApproveExchangerFunc      func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyAndMintNFTByApprovedExchangerFunc   func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	BuyNFTByExchangerFunc                  func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	AddExchangerTokenFunc                  func(StateDB, common.Address, *big.Int)
	ModifyOpenExchangerTimeFunc            func(StateDB, common.Address, *big.Int)
	SubExchangerTokenFunc                  func(StateDB, common.Address, *big.Int)
	SubExchangerBalanceFunc                func(StateDB, common.Address, *big.Int)
	VerifyExchangerBalanceFunc             func(StateDB, common.Address, *big.Int) bool
	GetNftAddressAndLevelFunc              func(string) (common.Address, int, error)
	VoteOfficialNFTFunc                    func(StateDB, *types.NominatedOfficialNFT, *big.Int) error
	ElectNominatedOfficialNFTFunc          func(StateDB, *big.Int)
	NextIndexFunc                          func(db StateDB) *big.Int
	VoteOfficialNFTByApprovedExchangerFunc func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	//ChangeRewardFlagFunc                   func(StateDB, common.Address, uint8)
	//PledgeNFTFunc                   func(StateDB, common.Address, *big.Int)
	//CancelPledgedNFTFunc            func(StateDB, common.Address)
	GetMergeNumberFunc func(StateDB, common.Address) uint32
	//GetPledgedFlagFunc              func(StateDB, common.Address) bool
	//GetNFTPledgedBlockNumberFunc    func(StateDB, common.Address) *big.Int
	RecoverValidatorCoefficientFunc           func(StateDB, common.Address) error
	BatchForcedSaleSNFTByApproveExchangerFunc func(StateDB, *big.Int, common.Address, common.Address, *types.Wormholes, *big.Int) error
	ChangeSnftRecipientFunc                   func(StateDB, common.Address, string)
	ChangeSNFTNoMergeFunc                     func(StateDB, common.Address, bool)
	GetDividendFunc                           func(StateDB, common.Address) error
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
	//CreateNFTByOfficial                CreateNFTByOfficialFunc
	CreateNFTByUser         CreateNFTByUserFunc
	ChangeApproveAddress    ChangeApproveAddressFunc
	CancelApproveAddress    CancelApproveAddressFunc
	ChangeNFTApproveAddress ChangeNFTApproveAddressFunc
	CancelNFTApproveAddress CancelNFTApproveAddressFunc
	ExchangeNFTToCurrency   ExchangeNFTToCurrencyFunc
	PledgeToken             PledgeTokenFunc
	StakerPledge            StakerPledgeFunc
	GetPledgedTime          GetPledgedTimeFunc
	GetStakerPledged        GetStakerPledgedFunc
	MinerConsign            MinerConsignFunc
	MinerBecome             MinerBecomeFunc
	CancelPledgedToken      CancelPledgedTokenFunc
	CancelStakerPledge      CancelStakerPledgeFunc
	OpenExchanger           OpenExchangerFunc
	CloseExchanger          CloseExchangerFunc
	GetExchangerFlag        GetExchangerFlagFunc
	GetOpenExchangerTime    GetOpenExchangerTimeFunc
	GetFeeRate              GetFeeRateFunc
	GetExchangerName        GetExchangerNameFunc
	GetExchangerURL         GetExchangerURLFunc
	GetApproveAddress       GetApproveAddressFunc
	//GetNFTBalance                      GetNFTBalanceFunc
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
	VerifyStakerPledgedBalance         VerifyStakerPledgedBalanceFunc
	InjectOfficialNFT                  InjectOfficialNFTFunc
	BuyNFTBySellerOrExchanger          BuyNFTBySellerOrExchangerFunc
	BuyNFTByBuyer                      BuyNFTByBuyerFunc
	BuyAndMintNFTByBuyer               BuyAndMintNFTByBuyerFunc
	BuyAndMintNFTByExchanger           BuyAndMintNFTByExchangerFunc
	BuyNFTByApproveExchanger           BuyNFTByApproveExchangerFunc
	BatchBuyNFTByApproveExchanger      BatchBuyNFTByApproveExchangerFunc
	BuyAndMintNFTByApprovedExchanger   BuyAndMintNFTByApprovedExchangerFunc
	BuyNFTByExchanger                  BuyNFTByExchangerFunc
	AddExchangerToken                  AddExchangerTokenFunc
	ModifyOpenExchangerTime            ModifyOpenExchangerTimeFunc
	SubExchangerToken                  SubExchangerTokenFunc
	SubExchangerBalance                SubExchangerBalanceFunc
	VerifyExchangerBalance             VerifyExchangerBalanceFunc
	GetNftAddressAndLevel              GetNftAddressAndLevelFunc
	VoteOfficialNFT                    VoteOfficialNFTFunc
	ElectNominatedOfficialNFT          ElectNominatedOfficialNFTFunc
	NextIndex                          NextIndexFunc
	VoteOfficialNFTByApprovedExchanger VoteOfficialNFTByApprovedExchangerFunc
	//ChangeRewardFlag                   ChangeRewardFlagFunc
	//PledgeNFT                   PledgeNFTFunc
	//CancelPledgedNFT            CancelPledgedNFTFunc
	GetMergeNumber GetMergeNumberFunc
	//GetPledgedFlag              GetPledgedFlagFunc
	//GetNFTPledgedBlockNumber    GetNFTPledgedBlockNumberFunc
	RecoverValidatorCoefficient           RecoverValidatorCoefficientFunc
	BatchForcedSaleSNFTByApproveExchanger BatchForcedSaleSNFTByApproveExchangerFunc
	ChangeSnftRecipient                   ChangeSnftRecipientFunc
	ChangeSNFTNoMerge                     ChangeSNFTNoMergeFunc
	GetDividend                           GetDividendFunc
	// Block information

	ParentHeader *types.Header

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

// hashMsg return the hash of plain msg
func hashMsg(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}

// recoverAddress recover the address from sig
func RecoverAddress(msg string, sigStr string) (common.Address, error) {
	if !strings.HasPrefix(sigStr, "0x") &&
		!strings.HasPrefix(sigStr, "0X") {
		return common.Address{}, fmt.Errorf("signature must be started with 0x or 0X")
	}
	sigData, err := hexutil.Decode(sigStr)
	if err != nil {
		return common.Address{}, err
	}
	if len(sigData) != 65 {
		return common.Address{}, fmt.Errorf("signature must be 65 bytes long")
	}
	if sigData[64] != 27 && sigData[64] != 28 {
		return common.Address{}, fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sigData[64] -= 27
	hash, _ := hashMsg([]byte(msg))
	//fmt.Println("sigdebug hash=", hexutil.Encode(hash))
	rpk, err := crypto.SigToPub(hash, sigData)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*rpk), nil
}

func GetSnftAddrs(db StateDB, nftParentAddress string, addr common.Address) []common.Address {
	var nftAddrs []common.Address
	emptyAddress := common.Address{}
	if strings.HasPrefix(nftParentAddress, "0x") ||
		strings.HasPrefix(nftParentAddress, "0X") {
		nftParentAddress = string([]byte(nftParentAddress)[2:])
	}

	if len(nftParentAddress) != 39 {
		return nftAddrs
	}

	addrInt := big.NewInt(0)
	addrInt.SetString(nftParentAddress, 16)
	addrInt.Lsh(addrInt, 4)

	// 3. retrieve all the sibling leaf nodes of nftAddr
	siblingInt := big.NewInt(0)
	//nftAddrSLen := len(nftAddrS)
	for i := 0; i < 16; i++ {
		// 4. convert bigInt to common.Address, and then get Account from the trie.
		siblingInt.Add(addrInt, big.NewInt(int64(i)))
		//siblingAddr := common.BigToAddress(siblingInt)
		siblingAddrS := hex.EncodeToString(siblingInt.Bytes())
		siblingAddrSLen := len(siblingAddrS)
		var prefix0 string
		for i := 0; i < 40-siblingAddrSLen; i++ {
			prefix0 = prefix0 + "0"
		}
		siblingAddrS = prefix0 + siblingAddrS
		siblingAddr := common.HexToAddress(siblingAddrS)
		//fmt.Println("siblingAddr=", siblingAddr.String())

		siblingOwner := db.GetNFTOwner16(siblingAddr)
		if siblingOwner != emptyAddress &&
			siblingOwner != addr {
			nftAddrs = append(nftAddrs, siblingAddr)
		}
	}

	return nftAddrs
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

	// Fail if we're trying to transfer more than the available balance
	if nftTransaction {
		switch wormholes.Type {
		case 10:
			//pledgedBalance := evm.StateDB.GetStakerPledgedBalance(caller.Address(), addr)
			//if pledgedBalance.Cmp(value) != 0 {
			//	baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
			//	Erb1000 := big.NewInt(700)
			//	Erb1000.Mul(Erb1000, baseErb)
			//	if value.Sign() > 0 && !evm.Context.VerifyStakerPledgedBalance(evm.StateDB, caller.Address(), addr, new(big.Int).Add(value, Erb1000)) {
			//		return nil, gas, ErrInsufficientBalance
			//	}
			//}

		case 14:
			// recover buyer address
			msgText := wormholes.Buyer.Amount +
				wormholes.Buyer.NFTAddress +
				wormholes.Buyer.Exchanger +
				wormholes.Buyer.BlockNumber +
				wormholes.Buyer.Seller
			buyer, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
			if err != nil {
				return nil, gas, err
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}

		case 17:
			// recover buyer address
			msgText := wormholes.Buyer.Amount +
				wormholes.Buyer.Exchanger +
				wormholes.Buyer.BlockNumber +
				wormholes.Buyer.Seller
			buyer, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
			if err != nil {
				return nil, gas, err
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}
		case 18:
			// recover buyer address
			msgText := wormholes.Buyer.Amount +
				wormholes.Buyer.NFTAddress +
				wormholes.Buyer.Exchanger +
				wormholes.Buyer.BlockNumber +
				wormholes.Buyer.Seller
			buyer, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
			if err != nil {
				return nil, gas, err
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}
		case 19:
			// recover buyer address
			msgText := wormholes.Buyer.Amount +
				wormholes.Buyer.Exchanger +
				wormholes.Buyer.BlockNumber +
				wormholes.Buyer.Seller
			buyer, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
			if err != nil {
				return nil, gas, err
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}
		case 20:
			// recover buyer address
			msgText := wormholes.Buyer.Amount +
				wormholes.Buyer.NFTAddress +
				wormholes.Buyer.Exchanger +
				wormholes.Buyer.BlockNumber +
				wormholes.Buyer.Seller
			buyer, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
			if err != nil {
				return nil, gas, err
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}
		case 22:
			baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
			Erb100 := big.NewInt(700)
			Erb100.Mul(Erb100, baseErb)
			if value.Sign() > 0 && !evm.Context.VerifyExchangerBalance(evm.StateDB, caller.Address(), new(big.Int).Add(value, Erb100)) {
				return nil, gas, ErrInsufficientBalance
			}
		//case 24:
		case 27:
			// recover buyer address
			emptyAddress := common.Address{}
			var buyer common.Address
			if len(wormholes.BuyerAuth.Exchanger) > 0 &&
				len(wormholes.BuyerAuth.BlockNumber) > 0 &&
				len(wormholes.BuyerAuth.Sig) > 0 {
				buyer, err = RecoverAddress(wormholes.BuyerAuth.Exchanger+wormholes.BuyerAuth.BlockNumber, wormholes.BuyerAuth.Sig)
				if err != nil {
					return nil, gas, err
				}
			}
			if buyer == emptyAddress {
				msgText := wormholes.Buyer.Amount +
					wormholes.Buyer.NFTAddress +
					wormholes.Buyer.Exchanger +
					wormholes.Buyer.BlockNumber +
					wormholes.Buyer.Seller
				buyerApproved, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
				if err != nil {
					return nil, gas, err
				}
				buyer = buyerApproved
			}
			if value.Sign() > 0 && !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}
		case 28:
			// recover buyer address
			emptyAddress := common.Address{}
			var buyer common.Address
			if len(wormholes.BuyerAuth.Exchanger) > 0 &&
				len(wormholes.BuyerAuth.BlockNumber) > 0 &&
				len(wormholes.BuyerAuth.Sig) > 0 {
				buyer, err = RecoverAddress(wormholes.BuyerAuth.Exchanger+wormholes.BuyerAuth.BlockNumber, wormholes.BuyerAuth.Sig)
				if err != nil {
					return nil, gas, err
				}
			}
			if buyer == emptyAddress {
				msgText := wormholes.Buyer.Amount +
					wormholes.Buyer.NFTAddress +
					wormholes.Buyer.Exchanger +
					wormholes.Buyer.BlockNumber +
					wormholes.Buyer.Seller
				buyerApproved, err := RecoverAddress(msgText, wormholes.Buyer.Sig)
				if err != nil {
					return nil, gas, err
				}
				buyer = buyerApproved
			}

			nftAddress, _, err := evm.Context.GetNftAddressAndLevel(wormholes.Buyer.NFTAddress)
			if err != nil {
				return nil, gas, err
			}
			initamount := evm.StateDB.CalculateExchangeAmount(1, 1)
			amount := evm.StateDB.GetExchangAmount(nftAddress, initamount)

			snftAddrs := GetSnftAddrs(evm.StateDB, wormholes.Buyer.NFTAddress, buyer)
			snftNum := len(snftAddrs)
			value := new(big.Int).Mul(big.NewInt(int64(snftNum)), amount)
			if !evm.Context.CanTransfer(evm.StateDB, buyer, value) {
				return nil, gas, ErrInsufficientBalance
			}

		default:
			if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
				return nil, gas, ErrInsufficientBalance
			}
		}
	} else {
		if value.Sign() != 0 && !evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			return nil, gas, ErrInsufficientBalance
		}
	}

	snapshot := evm.StateDB.Snapshot()
	p, isPrecompile := evm.precompile(addr)

	//if !evm.StateDB.Exist(addr) && !nftTransaction {
	if !evm.StateDB.Exist(addr) {
		//if !isPrecompile && evm.chainRules.IsEIP158 && value.Sign() == 0 {
		//	// Calling a non existing account, don't do anything, but ping the tracer
		//	if evm.Config.Debug && evm.depth == 0 {
		//		evm.Config.Tracer.CaptureStart(evm, caller.Address(), addr, false, input, gas, value)
		//		evm.Config.Tracer.CaptureEnd(ret, 0, 0, nil)
		//	}
		//	return nil, gas, nil
		//}
		evm.StateDB.CreateAccount(addr)
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

func (evm *EVM) IsContractAddress(address common.Address) bool {
	emptyHash := common.Hash{}
	codeHash := evm.StateDB.GetCodeHash(address)
	if codeHash != emptyHash {
		return true
	}
	return false
}

func (evm *EVM) TransferNFTByContract(caller ContractRef, input []byte, gas uint64) (ret []byte, overGas uint64, err error) {
	prefix := "21eceff7"
	//strInput := "b88d4fde" +
	//	"000000000000000000000000d9145cce52d386f254917e481eb44e9943f39138" +
	//	"000000000000000000000000d9145cce52d386f254917e481eb44e9943f39138" +
	//	"0000000000000000000000000000000000000000000000000000000000000003" +
	//	"0000000000000000000000000000000000000000000000000000000000000008" +
	//	"0000000000000000000000000000000000000000000000000000000000000002" +
	//	"0102000000000000000000000000000000000000000000000000000000000000"
	strInput := hex.EncodeToString(input)
	constData1 := "0000000000000000000000000000000000000000000000000000000000000008"
	//0x21eceff70000000000000000000000005b38da6a701c568545dcfcb03fcb875f56beddc40000000000000000000000000000000000000000000000000000000000000001

	if len(input) != 138 {
		return nil, gas, errors.New("input len error")
	}
	if !strings.HasPrefix(strInput, prefix) {
		return nil, gas, errors.New("input data error")
	}
	fromBytes := input[16:36]
	toBytes := input[48:68]
	nftAddressBytes := input[80:100]

	data1 := input[100:132]
	data2 := input[132:164]
	//data3 := input[164:196]

	if hex.EncodeToString(data1) != constData1 {
		return nil, gas, errors.New("input format error")
	}
	bigData3Len, _ := new(big.Int).SetString(hex.EncodeToString(data2), 16)
	if bigData3Len.Uint64() > 32 {
		return nil, gas, errors.New("input data error")
	}

	from := common.BytesToAddress(fromBytes)
	to := common.BytesToAddress(toBytes)

	bigNftAddr := new(big.Int).SetBytes(nftAddressBytes)
	bigSnft, _ := new(big.Int).SetString("8000000000000000000000000000000000000", 16)
	var strNftAddress string
	strNftAddress = "0x"
	if bigNftAddr.Cmp(bigSnft) >= 0 { // snft
		strNftAddress = strNftAddress + bigNftAddr.Text(16)
	} else {
		strNftAddress = strNftAddress + hex.EncodeToString(nftAddressBytes)
	}
	//nftAddress := common.BytesToAddress(nftAddressBytes)

	if evm.Context.VerifyNFTOwner(evm.StateDB, strNftAddress, from) {
		err := evm.Context.TransferNFT(evm.StateDB, strNftAddress, to, evm.Context.BlockNumber)
		if err != nil {
			return nil, gas, err
		}
	}

	//if evm.IsContractAddress(to) {
	//	erc721 := OnERC721Received(caller.Address(), from, strNftAddress)
	//	ret, overGas, err = evm.Call(AccountRef(SnftVirtualContractAddress), to, erc721, gas, big.NewInt(0))
	//}

	return ret, overGas, err
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
		log.Error("HandleNFT() format error", "wormholes.Type", wormholes.Type, "error", formatErr, "blocknumber", evm.Context.BlockNumber.Uint64())
		return nil, gas, formatErr
	}

	switch wormholes.Type {
	case 0: // create nft by user
		log.Info("HandleNFT(), CreateNFTByUser>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if wormholes.Royalty <= 0 {
			log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type,
				"error", ErrRoyaltyNotMoreThan0, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrRoyaltyNotMoreThan0
		}
		if wormholes.Royalty >= 10000 {
			log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type,
				"error", ErrRoyaltyNotLessthan10000, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrRoyaltyNotLessthan10000
		}

		exchanger := common.Address{}
		if len(wormholes.Exchanger) > 2 {
			if !strings.HasPrefix(wormholes.Exchanger, "0x") &&
				!strings.HasPrefix(wormholes.Exchanger, "0X") {
				log.Error("HandleNFT(), CreateNFTByUser(), exchanger format error",
					"wormholes.Exchanger", wormholes.Exchanger, "blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, ErrExchangerFormat
			}
			exchanger = common.HexToAddress(wormholes.Exchanger)

			exchangerFlag := evm.Context.GetExchangerFlag(evm.StateDB, exchanger)
			if exchangerFlag != true {
				log.Error("HandleNFT(), CreateNFTByUser", "wormholes.Type", wormholes.Type,
					"error", ErrNotExchanger, "blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, ErrNotExchanger
			}
		}

		evm.Context.CreateNFTByUser(evm.StateDB,
			exchanger,
			addr,
			wormholes.Royalty,
			wormholes.MetaURL)
		log.Info("HandleNFT(), CreateNFTByUser<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

	case 1: //transfer nft
		//nftAddress, _, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		//if err != nil {
		//	log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type, "error",
		//		err, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, err
		//}
		if evm.Context.VerifyNFTOwner(evm.StateDB, wormholes.NFTAddress, caller.Address()) {
			//if !evm.StateDB.Exist(nftAddress) {
			//	evm.StateDB.CreateAccount(nftAddress)
			//}
			log.Info("HandleNFT(), TransferNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
			err := evm.Context.TransferNFT(evm.StateDB, wormholes.NFTAddress, addr, evm.Context.BlockNumber)
			if err != nil {
				log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type,
					"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, err
			}
			log.Info("HandleNFT(), TransferNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
		} else {
			log.Error("HandleNFT(), TransferNFT", "wormholes.Type", wormholes.Type,
				"error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotOwner
		}

	case 2: //approve a nft's authority
		log.Info("HandleNFT(), ChangeNFTApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if !evm.Context.GetExchangerFlag(evm.StateDB, addr) {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type,
				"error", ErrNotExchanger, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotExchanger
		}
		nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		if int(level2) != level1 {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"input nft level", level1, "real nft level", level2, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotExistNft
		}
		if IsOfficialNFT(nftAddress) {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"error", ErrNotAllowedOfficialNFT, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotAllowedOfficialNFT
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		log.Info("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type,
			"wormholes.NFTAddress", wormholes.NFTAddress, "approvedaddress", addr.String(), "nftOwner", nftOwner.String(),
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), ChangeNFTApproveAddress", "wormholes.Type", wormholes.Type,
				"error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotOwner
		}
		evm.Context.ChangeNFTApproveAddress(
			evm.StateDB,
			nftAddress,
			addr)
		log.Info("HandleNFT(), ChangeNFTApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 3:
		log.Info("HandleNFT(), CancelNFTApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			log.Error("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		if int(level2) != level1 {
			log.Error("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"input nft level", level1, "real nft level", level2, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotExistNft
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		log.Info("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type,
			"wormholes.NFTAddress", wormholes.NFTAddress, "approvedaddress", addr.String(), "nftOwner", nftOwner.String(),
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), CancelNFTApproveAddress", "wormholes.Type", wormholes.Type,
				"error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotOwner
		}
		evm.Context.CancelNFTApproveAddress(
			evm.StateDB,
			nftAddress,
			addr)
		log.Info("HandleNFT(), CancelNFTApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 4: //approve all nft's authority
		if !evm.Context.GetExchangerFlag(evm.StateDB, addr) {
			log.Error("HandleNFT(), ChangeApproveAddress", "wormholes.Type", wormholes.Type,
				"error", ErrNotExchanger, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotExchanger
		}
		log.Info("HandleNFT(), ChangeApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		evm.Context.ChangeApproveAddress(
			evm.StateDB,
			caller.Address(),
			addr)
		log.Info("HandleNFT(), ChangeApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 5:
		log.Info("HandleNFT(), CancelApproveAddress>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		evm.Context.CancelApproveAddress(
			evm.StateDB,
			caller.Address(),
			addr)
		log.Info("HandleNFT(), CancelApproveAddress<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 6: //NFT exchange
		log.Info("HandleNFT(), ExchangeNFTToCurrency>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		if err != nil {
			return nil, gas, err
		}
		if !IsOfficialNFT(nftAddress) {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type,
				"nft address", wormholes.NFTAddress, "error", ErrNotMintByOfficial, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotMintByOfficial
		}
		nftOwner := evm.StateDB.GetNFTOwner16(nftAddress)
		if nftOwner != caller.Address() {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"nft owner", nftOwner, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotOwner
		}
		level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		if int(level2) != level1 {
			log.Error("HandleNFT(), ExchangeNFTToCurrency", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
				"input nft level", level1, "real nft level", level2, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrNotExistNft
		}
		//pledgedFlag := evm.Context.GetPledgedFlag(evm.StateDB, nftAddress)
		//if pledgedFlag {
		//	return nil, gas, ErrHasBeenPledged
		//}
		evm.Context.ExchangeNFTToCurrency(
			evm.StateDB,
			caller.Address(),
			wormholes.NFTAddress,
			evm.Context.BlockNumber)
		log.Info("HandleNFT(), ExchangeNFTToCurrency<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 7: //NFT pledge
		//log.Info("HandleNFT(), PledgeNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type,
		//	"blocknumber", evm.Context.BlockNumber.Uint64())
		//if !strings.HasPrefix(wormholes.NFTAddress, "0x") &&
		//	!strings.HasPrefix(wormholes.NFTAddress, "0X") {
		//	return nil, gas, ErrStartIndex
		//}
		//nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		//if err != nil {
		//	return nil, gas, err
		//}
		//if !IsOfficialNFT(nftAddress) {
		//	log.Error("HandleNFT(), PledgeNFT", "wormholes.Type", wormholes.Type,
		//		"nft address", wormholes.NFTAddress, "error", ErrNotMintByOfficial, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotMintByOfficial
		//}
		//level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		//if int(level2) != level1 {
		//	log.Error("HandleNFT(), PledgeNFT", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
		//		"input nft level", level1, "real nft level", level2, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotExistNft
		//}
		//if level2 < 1 {
		//	log.Error("HandleNFT(), PledgeNFT", "wormholes.Type", wormholes.Type,
		//		"error", ErrNotMergedSNFT, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotMergedSNFT
		//}
		//if !evm.Context.VerifyNFTOwner(evm.StateDB, wormholes.NFTAddress, caller.Address()) {
		//	return nil, gas, ErrNotOwner
		//}
		//if evm.Context.GetPledgedFlag(evm.StateDB, nftAddress) {
		//	return nil, gas, ErrRepeatedPledge
		//}
		//evm.Context.PledgeNFT(evm.StateDB, nftAddress, evm.Context.BlockNumber)
		//log.Info("HandleNFT(), PledgeNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type,
		//	"blocknumber", evm.Context.BlockNumber.Uint64())
	case 8: //cancel nft pledge
		//log.Info("HandleNFT(), CancelPledgedNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type,
		//	"blocknumber", evm.Context.BlockNumber.Uint64())
		//if !strings.HasPrefix(wormholes.NFTAddress, "0x") &&
		//	!strings.HasPrefix(wormholes.NFTAddress, "0X") {
		//	return nil, gas, ErrStartIndex
		//}
		//nftAddress, level1, err := evm.Context.GetNftAddressAndLevel(wormholes.NFTAddress)
		//if err != nil {
		//	return nil, gas, err
		//}
		//if !IsOfficialNFT(nftAddress) {
		//	log.Error("HandleNFT(), CancelPledgedNFT", "wormholes.Type", wormholes.Type,
		//		"nft address", wormholes.NFTAddress, "error", ErrNotMintByOfficial, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotMintByOfficial
		//}
		//level2 := evm.StateDB.GetNFTMergeLevel(nftAddress)
		//if int(level2) != level1 {
		//	log.Error("HandleNFT(), CancelPledgedNFT", "wormholes.Type", wormholes.Type, "nft address", wormholes.NFTAddress,
		//		"input nft level", level1, "real nft level", level2, "error", ErrNotOwner, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotExistNft
		//}
		//if level2 < 1 {
		//	log.Error("HandleNFT(), CancelPledgedNFT", "wormholes.Type", wormholes.Type,
		//		"error", ErrNotMergedSNFT, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrNotMergedSNFT
		//}
		//if !evm.Context.VerifyNFTOwner(evm.StateDB, wormholes.NFTAddress, caller.Address()) {
		//	return nil, gas, ErrNotOwner
		//}
		//nftPledgedTime := evm.Context.GetNFTPledgedBlockNumber(evm.StateDB, nftAddress)
		//if big.NewInt(CancelNFTPledgedInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, nftPledgedTime)) > 0 {
		//	log.Error("HandleNFT(), CancelPledgedNFT", "wormholes.Type", wormholes.Type,
		//		"error", ErrTooCloseToCancel, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrTooCloseToCancel
		//}
		//if !evm.Context.GetPledgedFlag(evm.StateDB, nftAddress) {
		//	return nil, gas, ErrNotPledge
		//}
		//evm.Context.CancelPledgedNFT(evm.StateDB, nftAddress)
		//log.Info("HandleNFT(), CancelPledgedNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type,
		//	"blocknumber", evm.Context.BlockNumber.Uint64())

	//case 9: // pledge token
	//	var firstTime bool = false
	//	baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
	//	Erb100000 := big.NewInt(70000)
	//	Erb100000.Mul(Erb100000, baseErb)
	//
	//	if !evm.Context.VerifyPledgedBalance(evm.StateDB, caller.Address(), Erb100000) {
	//		//if this account has not pledged
	//		if value.Cmp(Erb100000) < 0 {
	//			log.Error("HandleNFT(), PledgeToken", "wormholes.Type", wormholes.Type,
	//				"error", ErrNotMoreThan100000ERB, "blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, ErrNotMoreThan100000ERB
	//		}
	//		firstTime = true
	//	}
	//
	//	currentBlockNumber := new(big.Int).Set(evm.Context.BlockNumber)
	//	pledgedBalance := evm.StateDB.GetPledgedBalance(caller.Address())
	//	// if append pledgebalance, reset pledgedblocknumber
	//	if pledgedBalance != nil && pledgedBalance.Cmp(big.NewInt(0)) > 0 {
	//		pledgedBlockNumber := evm.StateDB.GetPledgedTime(caller.Address())
	//		height, err := UnstakingHeight(pledgedBalance, value, pledgedBlockNumber.Uint64(), currentBlockNumber.Uint64(), CancelPledgedInterval)
	//		if err != nil {
	//			return nil, gas, err
	//		}
	//		bigHeight := new(big.Int).SetUint64(height)
	//		bigCancelPledgedInterval := new(big.Int).SetUint64(CancelPledgedInterval)
	//		currentBlockNumber = new(big.Int).Add(currentBlockNumber, bigHeight)
	//		currentBlockNumber = new(big.Int).Sub(currentBlockNumber, bigCancelPledgedInterval)
	//	}
	//
	//	log.Info("HandleNFT()", "PledgeToken.req", wormholes, "blocknumber", evm.Context.BlockNumber.Uint64())
	//	if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
	//		log.Info("HandleNFT(), PledgeToken>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//		err := evm.Context.PledgeToken(evm.StateDB, caller.Address(), value, &wormholes, currentBlockNumber)
	//		if err != nil {
	//			log.Error("HandleNFT(), PledgeToken", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, err
	//		}
	//		if firstTime {
	//			evm.StateDB.AddValidatorCoefficient(caller.Address(), VALIDATOR_COEFFICIENT)
	//		}
	//		log.Info("HandleNFT(), PledgeToken<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//	} else {
	//		log.Error("HandleNFT(), PledgeToken", "wormholes.Type", wormholes.Type,
	//			"error", ErrInsufficientBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrInsufficientBalance
	//	}

	case 9: //staker token
		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		Erb100 := big.NewInt(700)
		Erb100.Mul(Erb100, baseErb)
		stakerpledged := evm.Context.GetStakerPledged(evm.StateDB, caller.Address(), addr)

		if stakerpledged.Balance == big0 {
			if value.Cmp(Erb100) < 0 {
				log.Error("HandleNFT(), StakerPledge", "wormholes.Type", wormholes.Type,
					"error", ErrNotMoreThan100ERB, "blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, ErrNotMoreThan100ERB
			}
		}

		currentBlockNumber := new(big.Int).Set(evm.Context.BlockNumber)

		log.Info("HandleNFT()", "StakerPledge.req", wormholes, "blocknumber", evm.Context.BlockNumber.Uint64())
		if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			log.Info("HandleNFT(), StakerPledge>>>>>>>>>>", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())

			err := evm.Context.StakerPledge(evm.StateDB, caller.Address(), addr, value, currentBlockNumber, &wormholes)
			if err != nil {
				log.Error("HandleNFT(), StakerPledge", "wormholes.Type", wormholes.Type,
					"blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, err
			}
			log.Info("HandleNFT(), StakerPledge<<<<<<<<<<", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
		} else {
			log.Error("HandleNFT(), StakerPledge", "wormholes.Type", wormholes.Type,
				"error", ErrInsufficientBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrInsufficientBalance
		}

	//case 10: // cancel pledge of token
	//	log.Info("HandleNFT(), CancelPledgedToken>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	//	pledgedTime := evm.Context.GetPledgedTime(evm.StateDB, caller.Address())
	//	if big.NewInt(CancelPledgedInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, pledgedTime)) > 0 {
	//		log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type,
	//			"error", ErrTooCloseToCancel, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrTooCloseToCancel
	//	}
	//
	//	pledgedBalance := evm.StateDB.GetPledgedBalance(caller.Address())
	//	if pledgedBalance.Cmp(value) == 0 {
	//		// cancel pledged balance
	//		log.Info("HandleNFT(), CancelPledgedToken, cancel all", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//		evm.Context.CancelPledgedToken(evm.StateDB, caller.Address(), value)
	//		coe := evm.StateDB.GetValidatorCoefficient(caller.Address())
	//		evm.StateDB.SubValidatorCoefficient(caller.Address(), coe)
	//
	//	} else {
	//		// cancel partial pledged balance
	//		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
	//		Erb100000 := big.NewInt(70000)
	//		Erb100000.Mul(Erb100000, baseErb)
	//
	//		if evm.Context.VerifyPledgedBalance(evm.StateDB, caller.Address(), new(big.Int).Add(Erb100000, value)) {
	//			log.Info("HandleNFT(), CancelPledgedToken, cancel partial", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//			evm.Context.CancelPledgedToken(evm.StateDB, caller.Address(), value)
	//		} else {
	//			log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type,
	//				"error", ErrInsufficientPledgedBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, ErrInsufficientPledgedBalance
	//		}
	//	}
	//	log.Info("HandleNFT(), CancelPledgedToken<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	case 10: // cancel pledge of token
		log.Info("HandleNFT(), CancelPledgedToken>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		stakerpledged := evm.Context.GetStakerPledged(evm.StateDB, caller.Address(), addr)
		//pledgedTime := stakerpledged.BlockNumber
		//if big.NewInt(CancelDayPledgedInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, pledgedTime)) > 0 {
		//	log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type,
		//		"error", ErrTooCloseToCancel, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrTooCloseToCancel
		//}

		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		Erb100 := big.NewInt(700)
		Erb100.Mul(Erb100, baseErb)
		pledgedBalance := stakerpledged.Balance

		if pledgedBalance.Cmp(value) != 0 {
			if Erb100.Cmp(new(big.Int).Sub(pledgedBalance, value)) > 0 {
				log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type,
					"error", "the after revocation is less than 700ERB", "blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, errors.New("the after revocation is less than 700ERB")
			}
		}
		if big.NewInt(CancelDayPledgedInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, stakerpledged.BlockNumber)) <= 0 {
			log.Info("HandleNFT(), CancelPledgedToken, cancel all", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
			Erb100000 := big.NewInt(70000)
			Erb100000.Mul(Erb100000, baseErb)
			if !evm.Context.VerifyPledgedBalance(evm.StateDB, addr, new(big.Int).Add(Erb100000, value)) {
				log.Info("HandleNFT(), CancelPledgedToken, cancel partial", "wormholes.Type", wormholes.Type,
					"blocknumber", evm.Context.BlockNumber.Uint64())
				//coe := evm.StateDB.GetValidatorCoefficient(addr)
				evm.StateDB.RemoveValidatorCoefficient(addr)
			}
			evm.Context.CancelStakerPledge(evm.StateDB, caller.Address(), addr, value, evm.Context.BlockNumber)
		} else {
			log.Error("HandleNFT(), CancelPledgedToken", "wormholes.Type", wormholes.Type,
				"error", ErrTooCloseToCancel, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrTooCloseToCancel
		}
		log.Info("HandleNFT(), CancelPledgedToken<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 12: // become miner
		baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
		Erb100000 := big.NewInt(70000)
		Erb100000.Mul(Erb100000, baseErb)

		if evm.Context.VerifyPledgedBalance(evm.StateDB, caller.Address(), Erb100000) {
			//if this account has not pledged
			log.Info("HandleNFT()", "MinerBecome.req", wormholes, "blocknumber", evm.Context.BlockNumber.Uint64())
			//if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
			log.Info("HandleNFT(), Start|MinerBecome>>>>>>>>>>", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
			err := evm.Context.MinerBecome(evm.StateDB, caller.Address(), &wormholes)
			if err != nil {
				log.Error("HandleNFT(), End|MinerBecome<<<<<<<<<<", "wormholes.Type", wormholes.Type,
					"blocknumber", evm.Context.BlockNumber.Uint64())
				return nil, gas, err
			}
			evm.StateDB.AddValidatorCoefficient(caller.Address(), VALIDATOR_COEFFICIENT)

			log.Info("HandleNFT(), End|MinerBecome<<<<<<<<<<", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())

		} else {
			log.Error("HandleNFT(), MinerBecome", "wormholes.Type", wormholes.Type,
				"error", ErrInsufficientPledgedBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, ErrInsufficientPledgedBalance
		}
	//case 11: //open exchanger
	//	log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	//	// value must be greater than or equal to 100 ERB
	//	unitErb, _ := new(big.Int).SetString("1000000000000000000", 10)
	//	if value.Cmp(new(big.Int).Mul(big.NewInt(700), unitErb)) < 0 {
	//		log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//			"error", ErrNotMoreThan100ERB, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrNotMoreThan100ERB
	//	}
	//	if wormholes.FeeRate <= 0 {
	//		log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//			"error", ErrFeeRateNotMoreThan0, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrFeeRateNotMoreThan0
	//	}
	//	if wormholes.FeeRate >= 10000 {
	//		log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//			"error", ErrFeeRateNotLessThan10000, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrFeeRateNotLessThan10000
	//	}
	//
	//	exchangerFlag := evm.Context.GetExchangerFlag(evm.StateDB, addr)
	//	if caller.Address() == addr {
	//		if exchangerFlag == true {
	//			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//				"error", ErrReopenExchanger, "blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, ErrReopenExchanger
	//		}
	//
	//		if evm.Context.CanTransfer(evm.StateDB, addr, value) {
	//			log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//			evm.Context.OpenExchanger(
	//				evm.StateDB,
	//				addr,
	//				value,
	//				evm.Context.BlockNumber,
	//				wormholes.FeeRate,
	//				wormholes.Name,
	//				wormholes.Url,
	//				wormholes.ProxyAddress)
	//			log.Info("HandleNFT(), OpenExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//
	//		} else {
	//			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//				"error", ErrInsufficientBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, ErrInsufficientBalance
	//		}
	//
	//	} else {
	//		// caller give addr a exchanger as a gift, if addr has been a exchanger, add value to exchangerbalance.
	//		if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
	//			evm.Context.Transfer(evm.StateDB, caller.Address(), addr, value)
	//		} else {
	//			log.Error("HandleNFT(), OpenExchanger", "wormholes.Type", wormholes.Type,
	//				"error", ErrInsufficientBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//			return nil, gas, ErrInsufficientBalance
	//		}
	//
	//		if exchangerFlag == true {
	//			evm.Context.AddExchangerToken(evm.StateDB, addr, value)
	//		} else {
	//			log.Info("HandleNFT(), OpenExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//			evm.Context.OpenExchanger(
	//				evm.StateDB,
	//				addr,
	//				value,
	//				evm.Context.BlockNumber,
	//				wormholes.FeeRate,
	//				wormholes.Name,
	//				wormholes.Url,
	//				wormholes.ProxyAddress)
	//			log.Info("HandleNFT(), OpenExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//				"blocknumber", evm.Context.BlockNumber.Uint64())
	//		}
	//	}

	//case 12: //close exchanger
	//	log.Info("HandleNFT(), CloseExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	//	openExchangerTime := evm.Context.GetOpenExchangerTime(evm.StateDB, caller.Address())
	//	if big.NewInt(CloseExchangerInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, openExchangerTime)) > 0 {
	//		log.Error("HandleNFT(), CloseExchanger", "wormholes.Type", wormholes.Type, "error", ErrTooCloseWithOpenExchanger)
	//		return nil, gas, ErrTooCloseWithOpenExchanger
	//	}
	//	evm.Context.CloseExchanger(evm.StateDB, caller.Address(), evm.Context.BlockNumber)
	//	log.Info("HandleNFT(), CloseExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	//	//evm.StateDB.CloseExchanger(caller.Address(), evm.Context.BlockNumber)
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
		log.Info("HandleNFT(), BuyNFTBySellerOrExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		log.Info("HandleNFT(), BuyNFTBySellerOrExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if err != nil {
			log.Error("HandleNFT(), BuyNFTBySellerOrExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
	case 15:
		log.Info("HandleNFT(), BuyNFTByBuyer>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		log.Info("HandleNFT(), BuyNFTByBuyer<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByBuyer", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
	case 16:
		log.Info("HandleNFT(), BuyAndMintNFTByBuyer>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		log.Info("HandleNFT(), BuyAndMintNFTByBuyer<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByBuyer", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
	case 17:
		log.Info("HandleNFT(), BuyAndMintNFTByExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), BuyAndMintNFTByExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 18:
		log.Info("HandleNFT(), BuyNFTByApproveExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByApproveExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), BuyNFTByApproveExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 19:
		log.Info("HandleNFT(), BuyAndMintNFTByApprovedExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		if err != nil {
			log.Error("HandleNFT(), BuyAndMintNFTByApprovedExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), BuyAndMintNFTByApprovedExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 20:
		log.Info("HandleNFT(), BuyNFTByExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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
		if err != nil {
			log.Error("HandleNFT(), BuyNFTByExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), BuyNFTByExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	//case 21:
	//	if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
	//		log.Info("HandleNFT(), AddExchangerToken>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//		currentBlockNumber := new(big.Int).Set(evm.Context.BlockNumber)
	//		exchangerBalance := evm.StateDB.GetExchangerBalance(caller.Address())
	//		if exchangerBalance != nil && exchangerBalance.Cmp(big.NewInt(0)) > 0 {
	//			openExchangerBlockNumber := evm.Context.GetOpenExchangerTime(evm.StateDB, caller.Address())
	//			height, err := UnstakingHeight(exchangerBalance, value, openExchangerBlockNumber.Uint64(), currentBlockNumber.Uint64(), CloseExchangerInterval)
	//			if err != nil {
	//				return nil, gas, err
	//			}
	//			bigHeight := new(big.Int).SetUint64(height)
	//			bigCloseExchangerInterval := new(big.Int).SetUint64(CloseExchangerInterval)
	//			currentBlockNumber = new(big.Int).Add(currentBlockNumber, bigHeight)
	//			currentBlockNumber = new(big.Int).Sub(currentBlockNumber, bigCloseExchangerInterval)
	//		}
	//		evm.Context.ModifyOpenExchangerTime(evm.StateDB, caller.Address(), currentBlockNumber)
	//		evm.Context.AddExchangerToken(evm.StateDB, caller.Address(), value)
	//		log.Info("HandleNFT(), AddExchangerToken<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//	} else {
	//		log.Error("HandleNFT(), AddExchangerToken", "wormholes.Type", wormholes.Type,
	//			"error", ErrInsufficientBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrInsufficientBalance
	//	}
	//case 22:
	//	openExchangerTime := evm.Context.GetOpenExchangerTime(evm.StateDB, caller.Address())
	//	if big.NewInt(CloseExchangerInterval).Cmp(new(big.Int).Sub(evm.Context.BlockNumber, openExchangerTime)) > 0 {
	//		log.Error("HandleNFT(), SubExchangerToken", "wormholes.Type", wormholes.Type,
	//			"error", ErrTooCloseForWithdraw, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrTooCloseForWithdraw
	//	}
	//	baseErb, _ := new(big.Int).SetString("1000000000000000000", 10)
	//	Erb100 := big.NewInt(700)
	//	Erb100.Mul(Erb100, baseErb)
	//	if evm.Context.VerifyExchangerBalance(evm.StateDB, caller.Address(), new(big.Int).Add(value, Erb100)) {
	//		log.Info("HandleNFT(), SubExchangerToken>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//		evm.Context.SubExchangerToken(evm.StateDB, caller.Address(), value)
	//		log.Info("HandleNFT(), SubExchangerToken<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//			"blocknumber", evm.Context.BlockNumber.Uint64())
	//	} else {
	//		log.Error("HandleNFT(), SubExchangerToken", "wormholes.Type", wormholes.Type,
	//			"error", ErrInsufficientExchangerBalance, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, ErrInsufficientExchangerBalance
	//	}
	case 31:
		//MinerConsign
		log.Info("HandleNFT()", "MinerConsign.req", wormholes, "blocknumber", evm.Context.BlockNumber.Uint64())
		//if evm.Context.CanTransfer(evm.StateDB, caller.Address(), value) {
		log.Info("HandleNFT(), Start|MinerConsign>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		err := evm.Context.MinerConsign(evm.StateDB, caller.Address(), &wormholes)
		if err != nil {
			log.Error("HandleNFT(), End|MinerConsign<<<<<<<<<<", "wormholes.Type", wormholes.Type,
				"blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), End|MinerConsign<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		//} else {
		//	log.Error("HandleNFT(), MinerConsign error", "wormholes.Type", wormholes.Type, "error", ErrInsufficientBalance)
		//	return nil, gas, ErrInsufficientBalance
		//}

	case 23:
		log.Info("HandleNFT(), VoteOfficialNFT>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		//if !strings.HasPrefix(wormholes.StartIndex, "0x") &&
		//	!strings.HasPrefix(wormholes.StartIndex, "0X") {
		//	return nil, gas, ErrStartIndex
		//}
		//startIndex, _ := new(big.Int).SetString(wormholes.StartIndex[2:], 16)
		startIndex := evm.StateDB.NextIndex()
		var number uint64 = 4096
		var royalty uint16 = 1000 // default 10%
		//if wormholes.Royalty <= 0 {
		//	log.Error("HandleNFT(), VoteOfficialNFT", "wormholes.Type", wormholes.Type,
		//		"error", ErrRoyaltyNotMoreThan0, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrRoyaltyNotMoreThan0
		//}
		//if wormholes.Royalty >= 10000 {
		//	log.Error("HandleNFT(), VoteOfficialNFT", "wormholes.Type", wormholes.Type,
		//		"error", ErrRoyaltyNotLessthan10000, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrRoyaltyNotLessthan10000
		//}
		var dir = wormholes.Dir
		if len(dir) <= 0 {
			dir = types.DefaultDir
		}
		var creator = wormholes.Creator
		if len(creator) <= 0 {
			creator = caller.Address().Hex()
		}
		nominatedNFT := types.NominatedOfficialNFT{
			InjectedOfficialNFT: types.InjectedOfficialNFT{
				Dir:        dir,
				StartIndex: startIndex,
				//Number: wormholes.Number,
				Number: number,
				//Royalty: wormholes.Royalty,
				Royalty: royalty,
				Creator: creator,
				Address: caller.Address(),
			},
		}
		err := evm.Context.VoteOfficialNFT(evm.StateDB, &nominatedNFT, evm.Context.BlockNumber)
		if err != nil {
			log.Error("HandleNFT(), VoteOfficialNFT", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), VoteOfficialNFT<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

	case 24:
		log.Info("HandleNFT(), VoteOfficialNFTByApprovedExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

		//if wormholes.Royalty <= 0 {
		//	log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type,
		//		"error", ErrRoyaltyNotMoreThan0, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrRoyaltyNotMoreThan0
		//}
		//if wormholes.Royalty >= 10000 {
		//	log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type,
		//		"error", ErrRoyaltyNotLessthan10000, "blocknumber", evm.Context.BlockNumber.Uint64())
		//	return nil, gas, ErrRoyaltyNotLessthan10000
		//}

		err := evm.Context.VoteOfficialNFTByApprovedExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		if err != nil {
			log.Error("HandleNFT(), VoteOfficialNFTByApprovedExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}

		log.Info("HandleNFT(), VoteOfficialNFTByApprovedExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

	case 25:
		log.Info("HandleNFT(), ChangeSnftRecipient>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

		evm.Context.ChangeSnftRecipient(evm.StateDB, caller.Address(), wormholes.ProxyAddress)

		log.Info("HandleNFT(), ChangeSnftRecipient<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 26:
		log.Info("HandleNFT(), RecoverValidatorCoefficient>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		err := evm.Context.RecoverValidatorCoefficient(evm.StateDB, caller.Address())
		if err != nil {
			return nil, gas, err
		}
		log.Info("HandleNFT(), RecoverValidatorCoefficient<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 27:
		log.Info("HandleNFT(), BatchBuyNFTByApproveExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		if value.Cmp(big.NewInt(0)) <= 0 {
			return nil, gas, ErrTransAmount
		}
		err := evm.Context.BatchBuyNFTByApproveExchanger(
			evm.StateDB,
			evm.Context.BlockNumber,
			caller.Address(),
			addr,
			&wormholes,
			value)
		if err != nil {
			log.Error("HandleNFT(), BatchBuyNFTByApproveExchanger", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), BatchBuyNFTByApproveExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	//case 28:
	//	log.Info("HandleNFT(), BatchForcedSaleSNFTByApproveExchanger>>>>>>>>>>", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	//	err := evm.Context.BatchForcedSaleSNFTByApproveExchanger(
	//		evm.StateDB,
	//		evm.Context.BlockNumber,
	//		caller.Address(),
	//		addr,
	//		&wormholes,
	//		value)
	//	if err != nil {
	//		log.Error("HandleNFT(), BatchForcedSaleSNFTByApproveExchanger", "wormholes.Type", wormholes.Type,
	//			"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
	//		return nil, gas, err
	//	}
	//	log.Info("HandleNFT(), BatchForcedSaleSNFTByApproveExchanger<<<<<<<<<<", "wormholes.Type", wormholes.Type,
	//		"blocknumber", evm.Context.BlockNumber.Uint64())
	case 29:
		log.Info("HandleNFT(), GetDividend>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
		err := evm.Context.GetDividend(evm.StateDB, caller.Address())
		if err != nil {
			log.Error("HandleNFT(), GetDividend", "wormholes.Type", wormholes.Type,
				"error", err, "blocknumber", evm.Context.BlockNumber.Uint64())
			return nil, gas, err
		}
		log.Info("HandleNFT(), GetDividend<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	case 30:
		log.Info("HandleNFT(), ChangeSNFTNoMerge>>>>>>>>>>", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())

		evm.Context.ChangeSNFTNoMerge(evm.StateDB, caller.Address(), wormholes.NoAutoMerge)

		log.Info("HandleNFT(), ChangeSNFTNoMerge<<<<<<<<<<", "wormholes.Type", wormholes.Type,
			"blocknumber", evm.Context.BlockNumber.Uint64())
	default:
		log.Error("HandleNFT()", "wormholes.Type", wormholes.Type, "error", ErrNotExistNFTType,
			"blocknumber", evm.Context.BlockNumber.Uint64())
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

// reference:https://github.com/wormholes-org/wormholes/issues/9
func unstakingHeight(stakedAmt *big.Int, appendAmt *big.Int, sno uint64, cno uint64, lockedNo uint64) uint64 {
	var curRemainingTime uint64

	if sno+lockedNo > cno {
		curRemainingTime = sno + lockedNo - cno
	}

	total := big.NewFloat(0).Add(new(big.Float).SetInt(stakedAmt), new(big.Float).SetInt(appendAmt))
	h1 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(stakedAmt), total), new(big.Float).SetInt(big.NewInt(int64(curRemainingTime))))
	h2 := big.NewFloat(0).Mul(big.NewFloat(0).Quo(new(big.Float).SetInt(appendAmt), total), new(big.Float).SetInt(big.NewInt(int64(lockedNo))))
	delayHeight, _ := big.NewFloat(0).Add(h1, h2).Uint64()
	if delayHeight < lockedNo/2 {
		delayHeight = lockedNo / 2
	}
	return delayHeight
}

func checkParams(stakedAmt *big.Int, appendAmt *big.Int, sno uint64, cno uint64) (bool, error) {
	if stakedAmt.Cmp(big.NewInt(0)) == 0 || appendAmt.Cmp(big.NewInt(0)) == 0 {
		return false, errors.New("illegal amount")
	}
	if cno == 0 {
		return false, errors.New("illegal height")
	}
	if sno > cno {
		return false, errors.New("the current height is less than the starting height of the pledge")
	}
	return true, nil
}
