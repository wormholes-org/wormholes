// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ethhelper

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

// CommoninterfaceMetaData contains all meta data concerning the Commoninterface contract.
var CommoninterfaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"TransferEvent721\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TransferEvent721Uri\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transferEvent1155\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// CommoninterfaceABI is the input ABI used to generate the binding from.
// Deprecated: Use CommoninterfaceMetaData.ABI instead.
var CommoninterfaceABI = CommoninterfaceMetaData.ABI

// Commoninterface is an auto generated Go binding around an Ethereum contract.
type Commoninterface struct {
	CommoninterfaceCaller     // Read-only binding to the contract
	CommoninterfaceTransactor // Write-only binding to the contract
	CommoninterfaceFilterer   // Log filterer for contract events
}

// CommoninterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type CommoninterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommoninterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CommoninterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommoninterfaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CommoninterfaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CommoninterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CommoninterfaceSession struct {
	Contract     *Commoninterface  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CommoninterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CommoninterfaceCallerSession struct {
	Contract *CommoninterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CommoninterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CommoninterfaceTransactorSession struct {
	Contract     *CommoninterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CommoninterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type CommoninterfaceRaw struct {
	Contract *Commoninterface // Generic contract binding to access the raw methods on
}

// CommoninterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CommoninterfaceCallerRaw struct {
	Contract *CommoninterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// CommoninterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CommoninterfaceTransactorRaw struct {
	Contract *CommoninterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCommoninterface creates a new instance of Commoninterface, bound to a specific deployed contract.
func NewCommoninterface(address common.Address, backend bind.ContractBackend) (*Commoninterface, error) {
	contract, err := bindCommoninterface(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Commoninterface{CommoninterfaceCaller: CommoninterfaceCaller{contract: contract}, CommoninterfaceTransactor: CommoninterfaceTransactor{contract: contract}, CommoninterfaceFilterer: CommoninterfaceFilterer{contract: contract}}, nil
}

// NewCommoninterfaceCaller creates a new read-only instance of Commoninterface, bound to a specific deployed contract.
func NewCommoninterfaceCaller(address common.Address, caller bind.ContractCaller) (*CommoninterfaceCaller, error) {
	contract, err := bindCommoninterface(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CommoninterfaceCaller{contract: contract}, nil
}

// NewCommoninterfaceTransactor creates a new write-only instance of Commoninterface, bound to a specific deployed contract.
func NewCommoninterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*CommoninterfaceTransactor, error) {
	contract, err := bindCommoninterface(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CommoninterfaceTransactor{contract: contract}, nil
}

// NewCommoninterfaceFilterer creates a new log filterer instance of Commoninterface, bound to a specific deployed contract.
func NewCommoninterfaceFilterer(address common.Address, filterer bind.ContractFilterer) (*CommoninterfaceFilterer, error) {
	contract, err := bindCommoninterface(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CommoninterfaceFilterer{contract: contract}, nil
}

// bindCommoninterface binds a generic wrapper to an already deployed contract.
func bindCommoninterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CommoninterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Commoninterface *CommoninterfaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Commoninterface.Contract.CommoninterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Commoninterface *CommoninterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Commoninterface.Contract.CommoninterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Commoninterface *CommoninterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Commoninterface.Contract.CommoninterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Commoninterface *CommoninterfaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Commoninterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Commoninterface *CommoninterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Commoninterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Commoninterface *CommoninterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Commoninterface.Contract.contract.Transact(opts, method, params...)
}

// TransferEvent721 is a free data retrieval call binding the contract method 0x2c9906dc.
//
// Solidity: function TransferEvent721() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCaller) TransferEvent721(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Commoninterface.contract.Call(opts, &out, "TransferEvent721")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TransferEvent721 is a free data retrieval call binding the contract method 0x2c9906dc.
//
// Solidity: function TransferEvent721() view returns(bytes32)
func (_Commoninterface *CommoninterfaceSession) TransferEvent721() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent721(&_Commoninterface.CallOpts)
}

// TransferEvent721 is a free data retrieval call binding the contract method 0x2c9906dc.
//
// Solidity: function TransferEvent721() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCallerSession) TransferEvent721() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent721(&_Commoninterface.CallOpts)
}

// TransferEvent721Uri is a free data retrieval call binding the contract method 0x0d282aab.
//
// Solidity: function TransferEvent721Uri() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCaller) TransferEvent721Uri(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Commoninterface.contract.Call(opts, &out, "TransferEvent721Uri")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TransferEvent721Uri is a free data retrieval call binding the contract method 0x0d282aab.
//
// Solidity: function TransferEvent721Uri() view returns(bytes32)
func (_Commoninterface *CommoninterfaceSession) TransferEvent721Uri() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent721Uri(&_Commoninterface.CallOpts)
}

// TransferEvent721Uri is a free data retrieval call binding the contract method 0x0d282aab.
//
// Solidity: function TransferEvent721Uri() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCallerSession) TransferEvent721Uri() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent721Uri(&_Commoninterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Commoninterface *CommoninterfaceCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Commoninterface.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Commoninterface *CommoninterfaceSession) Name() (string, error) {
	return _Commoninterface.Contract.Name(&_Commoninterface.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Commoninterface *CommoninterfaceCallerSession) Name() (string, error) {
	return _Commoninterface.Contract.Name(&_Commoninterface.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Commoninterface *CommoninterfaceCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _Commoninterface.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Commoninterface *CommoninterfaceSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Commoninterface.Contract.TokenURI(&_Commoninterface.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_Commoninterface *CommoninterfaceCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _Commoninterface.Contract.TokenURI(&_Commoninterface.CallOpts, tokenId)
}

// TransferEvent1155 is a free data retrieval call binding the contract method 0x8165f107.
//
// Solidity: function transferEvent1155() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCaller) TransferEvent1155(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Commoninterface.contract.Call(opts, &out, "transferEvent1155")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// TransferEvent1155 is a free data retrieval call binding the contract method 0x8165f107.
//
// Solidity: function transferEvent1155() view returns(bytes32)
func (_Commoninterface *CommoninterfaceSession) TransferEvent1155() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent1155(&_Commoninterface.CallOpts)
}

// TransferEvent1155 is a free data retrieval call binding the contract method 0x8165f107.
//
// Solidity: function transferEvent1155() view returns(bytes32)
func (_Commoninterface *CommoninterfaceCallerSession) TransferEvent1155() ([32]byte, error) {
	return _Commoninterface.Contract.TransferEvent1155(&_Commoninterface.CallOpts)
}
