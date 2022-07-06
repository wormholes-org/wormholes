// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package trade

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// TradeMetaData contains all meta data concerning the Trade contract.
var TradeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_weth\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"BIDING\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"BIDINGBATCH\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"PRICING\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"PRICINGBATCH\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_toSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"biding1155\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_toSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"biding1155Batch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_royaltyRatio\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_minerSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_toSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"biding1155Mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_toSig\",\"type\":\"bytes\"}],\"name\":\"biding721\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_toSig\",\"type\":\"bytes\"}],\"name\":\"biding721Batch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"feeInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_fee\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"getRoyalty\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"getRoyaltyBatch\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_receivers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_royaltys\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_signer\",\"type\":\"address\"}],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"nonceAdd\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_fromSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"pricing1155\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_fromSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"pricing1155Batch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_royaltyRatio\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_minerSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_fromSig\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"pricing1155Mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_fromSig\",\"type\":\"bytes\"}],\"name\":\"pricing721\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"_fromSig\",\"type\":\"bytes\"}],\"name\":\"pricing721Batch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"_sig\",\"type\":\"bytes\"}],\"name\":\"recoverSigner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"_feeRatio\",\"type\":\"uint16\"}],\"name\":\"setFeeRatio\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeReceiver\",\"type\":\"address\"}],\"name\":\"setFeeReciver\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_weth\",\"type\":\"address\"}],\"name\":\"setWeth\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"signHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"weth\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// TradeABI is the input ABI used to generate the binding from.
// Deprecated: Use TradeMetaData.ABI instead.
var TradeABI = TradeMetaData.ABI

// Trade is an auto generated Go binding around an Ethereum contract.
type Trade struct {
	TradeCaller     // Read-only binding to the contract
	TradeTransactor // Write-only binding to the contract
	TradeFilterer   // Log filterer for contract events
}

// TradeCaller is an auto generated read-only Go binding around an Ethereum contract.
type TradeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TradeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TradeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TradeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TradeSession struct {
	Contract     *Trade            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TradeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TradeCallerSession struct {
	Contract *TradeCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TradeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TradeTransactorSession struct {
	Contract     *TradeTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TradeRaw is an auto generated low-level Go binding around an Ethereum contract.
type TradeRaw struct {
	Contract *Trade // Generic contract binding to access the raw methods on
}

// TradeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TradeCallerRaw struct {
	Contract *TradeCaller // Generic read-only contract binding to access the raw methods on
}

// TradeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TradeTransactorRaw struct {
	Contract *TradeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrade creates a new instance of Trade, bound to a specific deployed contract.
func NewTrade(address common.Address, backend bind.ContractBackend) (*Trade, error) {
	contract, err := bindTrade(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Trade{TradeCaller: TradeCaller{contract: contract}, TradeTransactor: TradeTransactor{contract: contract}, TradeFilterer: TradeFilterer{contract: contract}}, nil
}

// NewTradeCaller creates a new read-only instance of Trade, bound to a specific deployed contract.
func NewTradeCaller(address common.Address, caller bind.ContractCaller) (*TradeCaller, error) {
	contract, err := bindTrade(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TradeCaller{contract: contract}, nil
}

// NewTradeTransactor creates a new write-only instance of Trade, bound to a specific deployed contract.
func NewTradeTransactor(address common.Address, transactor bind.ContractTransactor) (*TradeTransactor, error) {
	contract, err := bindTrade(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TradeTransactor{contract: contract}, nil
}

// NewTradeFilterer creates a new log filterer instance of Trade, bound to a specific deployed contract.
func NewTradeFilterer(address common.Address, filterer bind.ContractFilterer) (*TradeFilterer, error) {
	contract, err := bindTrade(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TradeFilterer{contract: contract}, nil
}

// bindTrade binds a generic wrapper to an already deployed contract.
func bindTrade(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TradeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Trade *TradeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Trade.Contract.TradeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Trade *TradeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.Contract.TradeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Trade *TradeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Trade.Contract.TradeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Trade *TradeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Trade.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Trade *TradeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Trade *TradeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Trade.Contract.contract.Transact(opts, method, params...)
}

// FeeInfo is a free data retrieval call binding the contract method 0xe44f5c50.
//
// Solidity: function feeInfo(uint256 _price) view returns(address _receiver, uint256 _fee)
func (_Trade *TradeCaller) FeeInfo(opts *bind.CallOpts, _price *big.Int) (struct {
	Receiver common.Address
	Fee      *big.Int
}, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "feeInfo", _price)

	outstruct := new(struct {
		Receiver common.Address
		Fee      *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Fee = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// FeeInfo is a free data retrieval call binding the contract method 0xe44f5c50.
//
// Solidity: function feeInfo(uint256 _price) view returns(address _receiver, uint256 _fee)
func (_Trade *TradeSession) FeeInfo(_price *big.Int) (struct {
	Receiver common.Address
	Fee      *big.Int
}, error) {
	return _Trade.Contract.FeeInfo(&_Trade.CallOpts, _price)
}

// FeeInfo is a free data retrieval call binding the contract method 0xe44f5c50.
//
// Solidity: function feeInfo(uint256 _price) view returns(address _receiver, uint256 _fee)
func (_Trade *TradeCallerSession) FeeInfo(_price *big.Int) (struct {
	Receiver common.Address
	Fee      *big.Int
}, error) {
	return _Trade.Contract.FeeInfo(&_Trade.CallOpts, _price)
}

// GetRoyalty is a free data retrieval call binding the contract method 0xf533b802.
//
// Solidity: function getRoyalty(address _addr, uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Trade *TradeCaller) GetRoyalty(opts *bind.CallOpts, _addr common.Address, _id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "getRoyalty", _addr, _id, _price)

	outstruct := new(struct {
		Receiver common.Address
		Royalty  *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Royalty = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetRoyalty is a free data retrieval call binding the contract method 0xf533b802.
//
// Solidity: function getRoyalty(address _addr, uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Trade *TradeSession) GetRoyalty(_addr common.Address, _id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	return _Trade.Contract.GetRoyalty(&_Trade.CallOpts, _addr, _id, _price)
}

// GetRoyalty is a free data retrieval call binding the contract method 0xf533b802.
//
// Solidity: function getRoyalty(address _addr, uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Trade *TradeCallerSession) GetRoyalty(_addr common.Address, _id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	return _Trade.Contract.GetRoyalty(&_Trade.CallOpts, _addr, _id, _price)
}

// GetRoyaltyBatch is a free data retrieval call binding the contract method 0x2c428baf.
//
// Solidity: function getRoyaltyBatch(address _addr, uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Trade *TradeCaller) GetRoyaltyBatch(opts *bind.CallOpts, _addr common.Address, _ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "getRoyaltyBatch", _addr, _ids, _price)

	outstruct := new(struct {
		Receivers []common.Address
		Royaltys  []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receivers = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Royaltys = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetRoyaltyBatch is a free data retrieval call binding the contract method 0x2c428baf.
//
// Solidity: function getRoyaltyBatch(address _addr, uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Trade *TradeSession) GetRoyaltyBatch(_addr common.Address, _ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	return _Trade.Contract.GetRoyaltyBatch(&_Trade.CallOpts, _addr, _ids, _price)
}

// GetRoyaltyBatch is a free data retrieval call binding the contract method 0x2c428baf.
//
// Solidity: function getRoyaltyBatch(address _addr, uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Trade *TradeCallerSession) GetRoyaltyBatch(_addr common.Address, _ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	return _Trade.Contract.GetRoyaltyBatch(&_Trade.CallOpts, _addr, _ids, _price)
}

// Nonce is a free data retrieval call binding the contract method 0x863cf34a.
//
// Solidity: function nonce(address _addr, uint256 _id, address _signer) view returns(uint256)
func (_Trade *TradeCaller) Nonce(opts *bind.CallOpts, _addr common.Address, _id *big.Int, _signer common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "nonce", _addr, _id, _signer)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0x863cf34a.
//
// Solidity: function nonce(address _addr, uint256 _id, address _signer) view returns(uint256)
func (_Trade *TradeSession) Nonce(_addr common.Address, _id *big.Int, _signer common.Address) (*big.Int, error) {
	return _Trade.Contract.Nonce(&_Trade.CallOpts, _addr, _id, _signer)
}

// Nonce is a free data retrieval call binding the contract method 0x863cf34a.
//
// Solidity: function nonce(address _addr, uint256 _id, address _signer) view returns(uint256)
func (_Trade *TradeCallerSession) Nonce(_addr common.Address, _id *big.Int, _signer common.Address) (*big.Int, error) {
	return _Trade.Contract.Nonce(&_Trade.CallOpts, _addr, _id, _signer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Trade *TradeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Trade *TradeSession) Owner() (common.Address, error) {
	return _Trade.Contract.Owner(&_Trade.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Trade *TradeCallerSession) Owner() (common.Address, error) {
	return _Trade.Contract.Owner(&_Trade.CallOpts)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x2e295ec9.
//
// Solidity: function recoverSigner(bytes _data, bytes _sig) pure returns(address)
func (_Trade *TradeCaller) RecoverSigner(opts *bind.CallOpts, _data []byte, _sig []byte) (common.Address, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "recoverSigner", _data, _sig)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RecoverSigner is a free data retrieval call binding the contract method 0x2e295ec9.
//
// Solidity: function recoverSigner(bytes _data, bytes _sig) pure returns(address)
func (_Trade *TradeSession) RecoverSigner(_data []byte, _sig []byte) (common.Address, error) {
	return _Trade.Contract.RecoverSigner(&_Trade.CallOpts, _data, _sig)
}

// RecoverSigner is a free data retrieval call binding the contract method 0x2e295ec9.
//
// Solidity: function recoverSigner(bytes _data, bytes _sig) pure returns(address)
func (_Trade *TradeCallerSession) RecoverSigner(_data []byte, _sig []byte) (common.Address, error) {
	return _Trade.Contract.RecoverSigner(&_Trade.CallOpts, _data, _sig)
}

// SignHash is a free data retrieval call binding the contract method 0x1d43c6fd.
//
// Solidity: function signHash(bytes _data) pure returns(bytes32)
func (_Trade *TradeCaller) SignHash(opts *bind.CallOpts, _data []byte) ([32]byte, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "signHash", _data)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SignHash is a free data retrieval call binding the contract method 0x1d43c6fd.
//
// Solidity: function signHash(bytes _data) pure returns(bytes32)
func (_Trade *TradeSession) SignHash(_data []byte) ([32]byte, error) {
	return _Trade.Contract.SignHash(&_Trade.CallOpts, _data)
}

// SignHash is a free data retrieval call binding the contract method 0x1d43c6fd.
//
// Solidity: function signHash(bytes _data) pure returns(bytes32)
func (_Trade *TradeCallerSession) SignHash(_data []byte) ([32]byte, error) {
	return _Trade.Contract.SignHash(&_Trade.CallOpts, _data)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Trade *TradeCaller) Weth(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Trade.contract.Call(opts, &out, "weth")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Trade *TradeSession) Weth() (common.Address, error) {
	return _Trade.Contract.Weth(&_Trade.CallOpts)
}

// Weth is a free data retrieval call binding the contract method 0x3fc8cef3.
//
// Solidity: function weth() view returns(address)
func (_Trade *TradeCallerSession) Weth() (common.Address, error) {
	return _Trade.Contract.Weth(&_Trade.CallOpts)
}

// Biding1155 is a paid mutator transaction binding the contract method 0x5ce42591.
//
// Solidity: function biding1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactor) Biding1155(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "biding1155", _addr, _from, _to, _id, _amount, _price, _toSig, _data)
}

// Biding1155 is a paid mutator transaction binding the contract method 0x5ce42591.
//
// Solidity: function biding1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeSession) Biding1155(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _price, _toSig, _data)
}

// Biding1155 is a paid mutator transaction binding the contract method 0x5ce42591.
//
// Solidity: function biding1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactorSession) Biding1155(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _price, _toSig, _data)
}

// Biding1155Batch is a paid mutator transaction binding the contract method 0xe16465e2.
//
// Solidity: function biding1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactor) Biding1155Batch(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "biding1155Batch", _addr, _from, _to, _ids, _amounts, _price, _toSig, _data)
}

// Biding1155Batch is a paid mutator transaction binding the contract method 0xe16465e2.
//
// Solidity: function biding1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeSession) Biding1155Batch(_addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155Batch(&_Trade.TransactOpts, _addr, _from, _to, _ids, _amounts, _price, _toSig, _data)
}

// Biding1155Batch is a paid mutator transaction binding the contract method 0xe16465e2.
//
// Solidity: function biding1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, uint256 _price, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactorSession) Biding1155Batch(_addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _price *big.Int, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155Batch(&_Trade.TransactOpts, _addr, _from, _to, _ids, _amounts, _price, _toSig, _data)
}

// Biding1155Mint is a paid mutator transaction binding the contract method 0x118dd819.
//
// Solidity: function biding1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactor) Biding1155Mint(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "biding1155Mint", _addr, _from, _to, _id, _amount, _price, _royaltyRatio, _tokenURI, _minerSig, _toSig, _data)
}

// Biding1155Mint is a paid mutator transaction binding the contract method 0x118dd819.
//
// Solidity: function biding1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _toSig, bytes _data) returns()
func (_Trade *TradeSession) Biding1155Mint(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155Mint(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _price, _royaltyRatio, _tokenURI, _minerSig, _toSig, _data)
}

// Biding1155Mint is a paid mutator transaction binding the contract method 0x118dd819.
//
// Solidity: function biding1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint256 _price, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _toSig, bytes _data) returns()
func (_Trade *TradeTransactorSession) Biding1155Mint(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _price *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _toSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding1155Mint(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _price, _royaltyRatio, _tokenURI, _minerSig, _toSig, _data)
}

// Biding721 is a paid mutator transaction binding the contract method 0x2ffccd97.
//
// Solidity: function biding721(address _from, address _to, address _addr, uint256 _id, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeTransactor) Biding721(opts *bind.TransactOpts, _from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "biding721", _from, _to, _addr, _id, _price, _toSig)
}

// Biding721 is a paid mutator transaction binding the contract method 0x2ffccd97.
//
// Solidity: function biding721(address _from, address _to, address _addr, uint256 _id, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeSession) Biding721(_from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding721(&_Trade.TransactOpts, _from, _to, _addr, _id, _price, _toSig)
}

// Biding721 is a paid mutator transaction binding the contract method 0x2ffccd97.
//
// Solidity: function biding721(address _from, address _to, address _addr, uint256 _id, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeTransactorSession) Biding721(_from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding721(&_Trade.TransactOpts, _from, _to, _addr, _id, _price, _toSig)
}

// Biding721Batch is a paid mutator transaction binding the contract method 0x9550f7ee.
//
// Solidity: function biding721Batch(address _from, address _to, address _addr, uint256[] _ids, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeTransactor) Biding721Batch(opts *bind.TransactOpts, _from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "biding721Batch", _from, _to, _addr, _ids, _price, _toSig)
}

// Biding721Batch is a paid mutator transaction binding the contract method 0x9550f7ee.
//
// Solidity: function biding721Batch(address _from, address _to, address _addr, uint256[] _ids, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeSession) Biding721Batch(_from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding721Batch(&_Trade.TransactOpts, _from, _to, _addr, _ids, _price, _toSig)
}

// Biding721Batch is a paid mutator transaction binding the contract method 0x9550f7ee.
//
// Solidity: function biding721Batch(address _from, address _to, address _addr, uint256[] _ids, uint256 _price, bytes _toSig) returns()
func (_Trade *TradeTransactorSession) Biding721Batch(_from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _price *big.Int, _toSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Biding721Batch(&_Trade.TransactOpts, _from, _to, _addr, _ids, _price, _toSig)
}

// NonceAdd is a paid mutator transaction binding the contract method 0x0532450a.
//
// Solidity: function nonceAdd(address _addr, uint256 _id) returns()
func (_Trade *TradeTransactor) NonceAdd(opts *bind.TransactOpts, _addr common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "nonceAdd", _addr, _id)
}

// NonceAdd is a paid mutator transaction binding the contract method 0x0532450a.
//
// Solidity: function nonceAdd(address _addr, uint256 _id) returns()
func (_Trade *TradeSession) NonceAdd(_addr common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Trade.Contract.NonceAdd(&_Trade.TransactOpts, _addr, _id)
}

// NonceAdd is a paid mutator transaction binding the contract method 0x0532450a.
//
// Solidity: function nonceAdd(address _addr, uint256 _id) returns()
func (_Trade *TradeTransactorSession) NonceAdd(_addr common.Address, _id *big.Int) (*types.Transaction, error) {
	return _Trade.Contract.NonceAdd(&_Trade.TransactOpts, _addr, _id)
}

// Pricing1155 is a paid mutator transaction binding the contract method 0xcdf1ce44.
//
// Solidity: function pricing1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactor) Pricing1155(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pricing1155", _addr, _from, _to, _id, _amount, _fromSig, _data)
}

// Pricing1155 is a paid mutator transaction binding the contract method 0xcdf1ce44.
//
// Solidity: function pricing1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeSession) Pricing1155(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _fromSig, _data)
}

// Pricing1155 is a paid mutator transaction binding the contract method 0xcdf1ce44.
//
// Solidity: function pricing1155(address _addr, address _from, address _to, uint256 _id, uint256 _amount, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactorSession) Pricing1155(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _fromSig, _data)
}

// Pricing1155Batch is a paid mutator transaction binding the contract method 0xb9cae187.
//
// Solidity: function pricing1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactor) Pricing1155Batch(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pricing1155Batch", _addr, _from, _to, _ids, _amounts, _fromSig, _data)
}

// Pricing1155Batch is a paid mutator transaction binding the contract method 0xb9cae187.
//
// Solidity: function pricing1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeSession) Pricing1155Batch(_addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155Batch(&_Trade.TransactOpts, _addr, _from, _to, _ids, _amounts, _fromSig, _data)
}

// Pricing1155Batch is a paid mutator transaction binding the contract method 0xb9cae187.
//
// Solidity: function pricing1155Batch(address _addr, address _from, address _to, uint256[] _ids, uint256[] _amounts, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactorSession) Pricing1155Batch(_addr common.Address, _from common.Address, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155Batch(&_Trade.TransactOpts, _addr, _from, _to, _ids, _amounts, _fromSig, _data)
}

// Pricing1155Mint is a paid mutator transaction binding the contract method 0x668d713e.
//
// Solidity: function pricing1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactor) Pricing1155Mint(opts *bind.TransactOpts, _addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pricing1155Mint", _addr, _from, _to, _id, _amount, _royaltyRatio, _tokenURI, _minerSig, _fromSig, _data)
}

// Pricing1155Mint is a paid mutator transaction binding the contract method 0x668d713e.
//
// Solidity: function pricing1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeSession) Pricing1155Mint(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155Mint(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _royaltyRatio, _tokenURI, _minerSig, _fromSig, _data)
}

// Pricing1155Mint is a paid mutator transaction binding the contract method 0x668d713e.
//
// Solidity: function pricing1155Mint(address _addr, address _from, address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _minerSig, bytes _fromSig, bytes _data) payable returns()
func (_Trade *TradeTransactorSession) Pricing1155Mint(_addr common.Address, _from common.Address, _to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _minerSig []byte, _fromSig []byte, _data []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing1155Mint(&_Trade.TransactOpts, _addr, _from, _to, _id, _amount, _royaltyRatio, _tokenURI, _minerSig, _fromSig, _data)
}

// Pricing721 is a paid mutator transaction binding the contract method 0xfb8911ef.
//
// Solidity: function pricing721(address _from, address _to, address _addr, uint256 _id, bytes _fromSig) payable returns()
func (_Trade *TradeTransactor) Pricing721(opts *bind.TransactOpts, _from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pricing721", _from, _to, _addr, _id, _fromSig)
}

// Pricing721 is a paid mutator transaction binding the contract method 0xfb8911ef.
//
// Solidity: function pricing721(address _from, address _to, address _addr, uint256 _id, bytes _fromSig) payable returns()
func (_Trade *TradeSession) Pricing721(_from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing721(&_Trade.TransactOpts, _from, _to, _addr, _id, _fromSig)
}

// Pricing721 is a paid mutator transaction binding the contract method 0xfb8911ef.
//
// Solidity: function pricing721(address _from, address _to, address _addr, uint256 _id, bytes _fromSig) payable returns()
func (_Trade *TradeTransactorSession) Pricing721(_from common.Address, _to common.Address, _addr common.Address, _id *big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing721(&_Trade.TransactOpts, _from, _to, _addr, _id, _fromSig)
}

// Pricing721Batch is a paid mutator transaction binding the contract method 0xb27b23ec.
//
// Solidity: function pricing721Batch(address _from, address _to, address _addr, uint256[] _ids, bytes _fromSig) payable returns()
func (_Trade *TradeTransactor) Pricing721Batch(opts *bind.TransactOpts, _from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "pricing721Batch", _from, _to, _addr, _ids, _fromSig)
}

// Pricing721Batch is a paid mutator transaction binding the contract method 0xb27b23ec.
//
// Solidity: function pricing721Batch(address _from, address _to, address _addr, uint256[] _ids, bytes _fromSig) payable returns()
func (_Trade *TradeSession) Pricing721Batch(_from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing721Batch(&_Trade.TransactOpts, _from, _to, _addr, _ids, _fromSig)
}

// Pricing721Batch is a paid mutator transaction binding the contract method 0xb27b23ec.
//
// Solidity: function pricing721Batch(address _from, address _to, address _addr, uint256[] _ids, bytes _fromSig) payable returns()
func (_Trade *TradeTransactorSession) Pricing721Batch(_from common.Address, _to common.Address, _addr common.Address, _ids []*big.Int, _fromSig []byte) (*types.Transaction, error) {
	return _Trade.Contract.Pricing721Batch(&_Trade.TransactOpts, _from, _to, _addr, _ids, _fromSig)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Trade.Contract.RenounceOwnership(&_Trade.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Trade *TradeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Trade.Contract.RenounceOwnership(&_Trade.TransactOpts)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x9d05f2a4.
//
// Solidity: function setFeeRatio(uint16 _feeRatio) returns()
func (_Trade *TradeTransactor) SetFeeRatio(opts *bind.TransactOpts, _feeRatio uint16) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "setFeeRatio", _feeRatio)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x9d05f2a4.
//
// Solidity: function setFeeRatio(uint16 _feeRatio) returns()
func (_Trade *TradeSession) SetFeeRatio(_feeRatio uint16) (*types.Transaction, error) {
	return _Trade.Contract.SetFeeRatio(&_Trade.TransactOpts, _feeRatio)
}

// SetFeeRatio is a paid mutator transaction binding the contract method 0x9d05f2a4.
//
// Solidity: function setFeeRatio(uint16 _feeRatio) returns()
func (_Trade *TradeTransactorSession) SetFeeRatio(_feeRatio uint16) (*types.Transaction, error) {
	return _Trade.Contract.SetFeeRatio(&_Trade.TransactOpts, _feeRatio)
}

// SetFeeReciver is a paid mutator transaction binding the contract method 0x25cb10a5.
//
// Solidity: function setFeeReciver(address _feeReceiver) returns()
func (_Trade *TradeTransactor) SetFeeReciver(opts *bind.TransactOpts, _feeReceiver common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "setFeeReciver", _feeReceiver)
}

// SetFeeReciver is a paid mutator transaction binding the contract method 0x25cb10a5.
//
// Solidity: function setFeeReciver(address _feeReceiver) returns()
func (_Trade *TradeSession) SetFeeReciver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _Trade.Contract.SetFeeReciver(&_Trade.TransactOpts, _feeReceiver)
}

// SetFeeReciver is a paid mutator transaction binding the contract method 0x25cb10a5.
//
// Solidity: function setFeeReciver(address _feeReceiver) returns()
func (_Trade *TradeTransactorSession) SetFeeReciver(_feeReceiver common.Address) (*types.Transaction, error) {
	return _Trade.Contract.SetFeeReciver(&_Trade.TransactOpts, _feeReceiver)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _weth) returns()
func (_Trade *TradeTransactor) SetWeth(opts *bind.TransactOpts, _weth common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "setWeth", _weth)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _weth) returns()
func (_Trade *TradeSession) SetWeth(_weth common.Address) (*types.Transaction, error) {
	return _Trade.Contract.SetWeth(&_Trade.TransactOpts, _weth)
}

// SetWeth is a paid mutator transaction binding the contract method 0xb8d1452f.
//
// Solidity: function setWeth(address _weth) returns()
func (_Trade *TradeTransactorSession) SetWeth(_weth common.Address) (*types.Transaction, error) {
	return _Trade.Contract.SetWeth(&_Trade.TransactOpts, _weth)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Trade *TradeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Trade.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Trade *TradeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Trade.Contract.TransferOwnership(&_Trade.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Trade *TradeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Trade.Contract.TransferOwnership(&_Trade.TransactOpts, newOwner)
}

// TradeBIDINGIterator is returned from FilterBIDING and is used to iterate over the raw logs and unpacked data for BIDING events raised by the Trade contract.
type TradeBIDINGIterator struct {
	Event *TradeBIDING // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TradeBIDINGIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeBIDING)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TradeBIDING)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TradeBIDINGIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeBIDINGIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeBIDING represents a BIDING event raised by the Trade contract.
type TradeBIDING struct {
	From   common.Address
	To     common.Address
	Addr   common.Address
	Id     *big.Int
	Amount *big.Int
	Price  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterBIDING is a free log retrieval operation binding the contract event 0xf64cc560cd0c25fe33108f5a79cd3104e22f65e64b260be4b2072104012fa60a.
//
// Solidity: event BIDING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) FilterBIDING(opts *bind.FilterOpts) (*TradeBIDINGIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "BIDING")
	if err != nil {
		return nil, err
	}
	return &TradeBIDINGIterator{contract: _Trade.contract, event: "BIDING", logs: logs, sub: sub}, nil
}

// WatchBIDING is a free log subscription operation binding the contract event 0xf64cc560cd0c25fe33108f5a79cd3104e22f65e64b260be4b2072104012fa60a.
//
// Solidity: event BIDING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) WatchBIDING(opts *bind.WatchOpts, sink chan<- *TradeBIDING) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "BIDING")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeBIDING)
				if err := _Trade.contract.UnpackLog(event, "BIDING", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBIDING is a log parse operation binding the contract event 0xf64cc560cd0c25fe33108f5a79cd3104e22f65e64b260be4b2072104012fa60a.
//
// Solidity: event BIDING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) ParseBIDING(log types.Log) (*TradeBIDING, error) {
	event := new(TradeBIDING)
	if err := _Trade.contract.UnpackLog(event, "BIDING", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TradeBIDINGBATCHIterator is returned from FilterBIDINGBATCH and is used to iterate over the raw logs and unpacked data for BIDINGBATCH events raised by the Trade contract.
type TradeBIDINGBATCHIterator struct {
	Event *TradeBIDINGBATCH // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TradeBIDINGBATCHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeBIDINGBATCH)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TradeBIDINGBATCH)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TradeBIDINGBATCHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeBIDINGBATCHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeBIDINGBATCH represents a BIDINGBATCH event raised by the Trade contract.
type TradeBIDINGBATCH struct {
	From    common.Address
	To      common.Address
	Addr    common.Address
	Ids     []*big.Int
	Amounts []*big.Int
	Price   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBIDINGBATCH is a free log retrieval operation binding the contract event 0x931a303fc7adb73da26664e6aa7bcc8a7e8804bd43eacb06e2a552380a95bcb4.
//
// Solidity: event BIDINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) FilterBIDINGBATCH(opts *bind.FilterOpts) (*TradeBIDINGBATCHIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "BIDINGBATCH")
	if err != nil {
		return nil, err
	}
	return &TradeBIDINGBATCHIterator{contract: _Trade.contract, event: "BIDINGBATCH", logs: logs, sub: sub}, nil
}

// WatchBIDINGBATCH is a free log subscription operation binding the contract event 0x931a303fc7adb73da26664e6aa7bcc8a7e8804bd43eacb06e2a552380a95bcb4.
//
// Solidity: event BIDINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) WatchBIDINGBATCH(opts *bind.WatchOpts, sink chan<- *TradeBIDINGBATCH) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "BIDINGBATCH")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeBIDINGBATCH)
				if err := _Trade.contract.UnpackLog(event, "BIDINGBATCH", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBIDINGBATCH is a log parse operation binding the contract event 0x931a303fc7adb73da26664e6aa7bcc8a7e8804bd43eacb06e2a552380a95bcb4.
//
// Solidity: event BIDINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) ParseBIDINGBATCH(log types.Log) (*TradeBIDINGBATCH, error) {
	event := new(TradeBIDINGBATCH)
	if err := _Trade.contract.UnpackLog(event, "BIDINGBATCH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TradeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Trade contract.
type TradeOwnershipTransferredIterator struct {
	Event *TradeOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TradeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradeOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TradeOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TradeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradeOwnershipTransferred represents a OwnershipTransferred event raised by the Trade contract.
type TradeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Trade *TradeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TradeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Trade.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TradeOwnershipTransferredIterator{contract: _Trade.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Trade *TradeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TradeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Trade.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradeOwnershipTransferred)
				if err := _Trade.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Trade *TradeFilterer) ParseOwnershipTransferred(log types.Log) (*TradeOwnershipTransferred, error) {
	event := new(TradeOwnershipTransferred)
	if err := _Trade.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TradePRICINGIterator is returned from FilterPRICING and is used to iterate over the raw logs and unpacked data for PRICING events raised by the Trade contract.
type TradePRICINGIterator struct {
	Event *TradePRICING // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TradePRICINGIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradePRICING)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TradePRICING)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TradePRICINGIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradePRICINGIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradePRICING represents a PRICING event raised by the Trade contract.
type TradePRICING struct {
	From   common.Address
	To     common.Address
	Addr   common.Address
	Id     *big.Int
	Amount *big.Int
	Price  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPRICING is a free log retrieval operation binding the contract event 0xb14cef5ea4cbf14922e663d47af6b5458327dda3f30ba4ffd69e1d63a305ed2c.
//
// Solidity: event PRICING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) FilterPRICING(opts *bind.FilterOpts) (*TradePRICINGIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "PRICING")
	if err != nil {
		return nil, err
	}
	return &TradePRICINGIterator{contract: _Trade.contract, event: "PRICING", logs: logs, sub: sub}, nil
}

// WatchPRICING is a free log subscription operation binding the contract event 0xb14cef5ea4cbf14922e663d47af6b5458327dda3f30ba4ffd69e1d63a305ed2c.
//
// Solidity: event PRICING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) WatchPRICING(opts *bind.WatchOpts, sink chan<- *TradePRICING) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "PRICING")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradePRICING)
				if err := _Trade.contract.UnpackLog(event, "PRICING", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePRICING is a log parse operation binding the contract event 0xb14cef5ea4cbf14922e663d47af6b5458327dda3f30ba4ffd69e1d63a305ed2c.
//
// Solidity: event PRICING(address from, address to, address addr, uint256 id, uint256 amount, uint256 price)
func (_Trade *TradeFilterer) ParsePRICING(log types.Log) (*TradePRICING, error) {
	event := new(TradePRICING)
	if err := _Trade.contract.UnpackLog(event, "PRICING", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TradePRICINGBATCHIterator is returned from FilterPRICINGBATCH and is used to iterate over the raw logs and unpacked data for PRICINGBATCH events raised by the Trade contract.
type TradePRICINGBATCHIterator struct {
	Event *TradePRICINGBATCH // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TradePRICINGBATCHIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TradePRICINGBATCH)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TradePRICINGBATCH)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TradePRICINGBATCHIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TradePRICINGBATCHIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TradePRICINGBATCH represents a PRICINGBATCH event raised by the Trade contract.
type TradePRICINGBATCH struct {
	From    common.Address
	To      common.Address
	Addr    common.Address
	Ids     []*big.Int
	Amounts []*big.Int
	Price   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPRICINGBATCH is a free log retrieval operation binding the contract event 0x2997733f6482a8747bf9ae79a5de9c56a5cb8302f126bfe615de4316ae6dd8e8.
//
// Solidity: event PRICINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) FilterPRICINGBATCH(opts *bind.FilterOpts) (*TradePRICINGBATCHIterator, error) {

	logs, sub, err := _Trade.contract.FilterLogs(opts, "PRICINGBATCH")
	if err != nil {
		return nil, err
	}
	return &TradePRICINGBATCHIterator{contract: _Trade.contract, event: "PRICINGBATCH", logs: logs, sub: sub}, nil
}

// WatchPRICINGBATCH is a free log subscription operation binding the contract event 0x2997733f6482a8747bf9ae79a5de9c56a5cb8302f126bfe615de4316ae6dd8e8.
//
// Solidity: event PRICINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) WatchPRICINGBATCH(opts *bind.WatchOpts, sink chan<- *TradePRICINGBATCH) (event.Subscription, error) {

	logs, sub, err := _Trade.contract.WatchLogs(opts, "PRICINGBATCH")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TradePRICINGBATCH)
				if err := _Trade.contract.UnpackLog(event, "PRICINGBATCH", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePRICINGBATCH is a log parse operation binding the contract event 0x2997733f6482a8747bf9ae79a5de9c56a5cb8302f126bfe615de4316ae6dd8e8.
//
// Solidity: event PRICINGBATCH(address from, address to, address addr, uint256[] ids, uint256[] amounts, uint256 price)
func (_Trade *TradeFilterer) ParsePRICINGBATCH(log types.Log) (*TradePRICINGBATCH, error) {
	event := new(TradePRICINGBATCH)
	if err := _Trade.contract.UnpackLog(event, "PRICINGBATCH", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
