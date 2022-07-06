// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package nft1155

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

// Nft1155MetaData contains all meta data concerning the Nft1155 contract.
var Nft1155MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_uri\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"enumLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"ApprovalLevel\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"ratio\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"Royalty\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"values\",\"type\":\"uint256[]\"}],\"name\":\"TransferBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"TransferSingle\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"value\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"URI\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"accounts\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"}],\"name\":\"balanceOfBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"info\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint16\",\"name\":\"_ratio\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"}],\"name\":\"infoBatch\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_receivers\",\"type\":\"address[]\"},{\"internalType\":\"uint16[]\",\"name\":\"_ratios\",\"type\":\"uint16[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_account\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_miner\",\"type\":\"address\"}],\"name\":\"isMiner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_super\",\"type\":\"address\"}],\"name\":\"isSuper\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_royaltyRatio\",\"type\":\"uint16\"},{\"internalType\":\"string\",\"name\":\"_tokenURI\",\"type\":\"string\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"_royaltyRatios\",\"type\":\"uint16[]\"},{\"internalType\":\"string[]\",\"name\":\"_tokenURIs\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"mintBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"_receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_royalty\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"_price\",\"type\":\"uint256\"}],\"name\":\"royaltyInfoBatch\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_receivers\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_royaltys\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256[]\",\"name\":\"ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeBatchTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"},{\"internalType\":\"enumLevel\",\"name\":\"_level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"_ratio\",\"type\":\"uint16\"}],\"name\":\"setRoyalty\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_ids\",\"type\":\"uint256[]\"},{\"internalType\":\"uint16[]\",\"name\":\"_ratios\",\"type\":\"uint16[]\"}],\"name\":\"setRoyaltyBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// Nft1155ABI is the input ABI used to generate the binding from.
// Deprecated: Use Nft1155MetaData.ABI instead.
var Nft1155ABI = Nft1155MetaData.ABI

// Nft1155 is an auto generated Go binding around an Ethereum contract.
type Nft1155 struct {
	Nft1155Caller     // Read-only binding to the contract
	Nft1155Transactor // Write-only binding to the contract
	Nft1155Filterer   // Log filterer for contract events
}

// Nft1155Caller is an auto generated read-only Go binding around an Ethereum contract.
type Nft1155Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Nft1155Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Nft1155Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Nft1155Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Nft1155Session struct {
	Contract     *Nft1155          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Nft1155CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Nft1155CallerSession struct {
	Contract *Nft1155Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Nft1155TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Nft1155TransactorSession struct {
	Contract     *Nft1155Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Nft1155Raw is an auto generated low-level Go binding around an Ethereum contract.
type Nft1155Raw struct {
	Contract *Nft1155 // Generic contract binding to access the raw methods on
}

// Nft1155CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Nft1155CallerRaw struct {
	Contract *Nft1155Caller // Generic read-only contract binding to access the raw methods on
}

// Nft1155TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Nft1155TransactorRaw struct {
	Contract *Nft1155Transactor // Generic write-only contract binding to access the raw methods on
}

// NewNft1155 creates a new instance of Nft1155, bound to a specific deployed contract.
func NewNft1155(address common.Address, backend bind.ContractBackend) (*Nft1155, error) {
	contract, err := bindNft1155(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Nft1155{Nft1155Caller: Nft1155Caller{contract: contract}, Nft1155Transactor: Nft1155Transactor{contract: contract}, Nft1155Filterer: Nft1155Filterer{contract: contract}}, nil
}

// NewNft1155Caller creates a new read-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Caller(address common.Address, caller bind.ContractCaller) (*Nft1155Caller, error) {
	contract, err := bindNft1155(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Caller{contract: contract}, nil
}

// NewNft1155Transactor creates a new write-only instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Transactor(address common.Address, transactor bind.ContractTransactor) (*Nft1155Transactor, error) {
	contract, err := bindNft1155(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Nft1155Transactor{contract: contract}, nil
}

// NewNft1155Filterer creates a new log filterer instance of Nft1155, bound to a specific deployed contract.
func NewNft1155Filterer(address common.Address, filterer bind.ContractFilterer) (*Nft1155Filterer, error) {
	contract, err := bindNft1155(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Nft1155Filterer{contract: contract}, nil
}

// bindNft1155 binds a generic wrapper to an already deployed contract.
func bindNft1155(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Nft1155ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.Nft1155Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.Nft1155Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Nft1155 *Nft1155CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Nft1155.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Nft1155 *Nft1155TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Nft1155 *Nft1155TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Nft1155.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Caller) BalanceOf(opts *bind.CallOpts, account common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOf", account, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155Session) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address account, uint256 id) view returns(uint256)
func (_Nft1155 *Nft1155CallerSession) BalanceOf(account common.Address, id *big.Int) (*big.Int, error) {
	return _Nft1155.Contract.BalanceOf(&_Nft1155.CallOpts, account, id)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Caller) BalanceOfBatch(opts *bind.CallOpts, accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "balanceOfBatch", accounts, ids)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155Session) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// BalanceOfBatch is a free data retrieval call binding the contract method 0x4e1273f4.
//
// Solidity: function balanceOfBatch(address[] accounts, uint256[] ids) view returns(uint256[])
func (_Nft1155 *Nft1155CallerSession) BalanceOfBatch(accounts []common.Address, ids []*big.Int) ([]*big.Int, error) {
	return _Nft1155.Contract.BalanceOfBatch(&_Nft1155.CallOpts, accounts, ids)
}

// Info is a free data retrieval call binding the contract method 0x2e340599.
//
// Solidity: function info(uint256 _id) view returns(address _receiver, uint16 _ratio)
func (_Nft1155 *Nft1155Caller) Info(opts *bind.CallOpts, _id *big.Int) (struct {
	Receiver common.Address
	Ratio    uint16
}, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "info", _id)

	outstruct := new(struct {
		Receiver common.Address
		Ratio    uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Ratio = *abi.ConvertType(out[1], new(uint16)).(*uint16)

	return *outstruct, err

}

// Info is a free data retrieval call binding the contract method 0x2e340599.
//
// Solidity: function info(uint256 _id) view returns(address _receiver, uint16 _ratio)
func (_Nft1155 *Nft1155Session) Info(_id *big.Int) (struct {
	Receiver common.Address
	Ratio    uint16
}, error) {
	return _Nft1155.Contract.Info(&_Nft1155.CallOpts, _id)
}

// Info is a free data retrieval call binding the contract method 0x2e340599.
//
// Solidity: function info(uint256 _id) view returns(address _receiver, uint16 _ratio)
func (_Nft1155 *Nft1155CallerSession) Info(_id *big.Int) (struct {
	Receiver common.Address
	Ratio    uint16
}, error) {
	return _Nft1155.Contract.Info(&_Nft1155.CallOpts, _id)
}

// InfoBatch is a free data retrieval call binding the contract method 0xd9e16218.
//
// Solidity: function infoBatch(uint256[] _ids) view returns(address[] _receivers, uint16[] _ratios)
func (_Nft1155 *Nft1155Caller) InfoBatch(opts *bind.CallOpts, _ids []*big.Int) (struct {
	Receivers []common.Address
	Ratios    []uint16
}, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "infoBatch", _ids)

	outstruct := new(struct {
		Receivers []common.Address
		Ratios    []uint16
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receivers = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.Ratios = *abi.ConvertType(out[1], new([]uint16)).(*[]uint16)

	return *outstruct, err

}

// InfoBatch is a free data retrieval call binding the contract method 0xd9e16218.
//
// Solidity: function infoBatch(uint256[] _ids) view returns(address[] _receivers, uint16[] _ratios)
func (_Nft1155 *Nft1155Session) InfoBatch(_ids []*big.Int) (struct {
	Receivers []common.Address
	Ratios    []uint16
}, error) {
	return _Nft1155.Contract.InfoBatch(&_Nft1155.CallOpts, _ids)
}

// InfoBatch is a free data retrieval call binding the contract method 0xd9e16218.
//
// Solidity: function infoBatch(uint256[] _ids) view returns(address[] _receivers, uint16[] _ratios)
func (_Nft1155 *Nft1155CallerSession) InfoBatch(_ids []*big.Int) (struct {
	Receivers []common.Address
	Ratios    []uint16
}, error) {
	return _Nft1155.Contract.InfoBatch(&_Nft1155.CallOpts, _ids)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _account, address _operator) view returns(bool)
func (_Nft1155 *Nft1155Caller) IsApprovedForAll(opts *bind.CallOpts, _account common.Address, _operator common.Address) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "isApprovedForAll", _account, _operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _account, address _operator) view returns(bool)
func (_Nft1155 *Nft1155Session) IsApprovedForAll(_account common.Address, _operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, _account, _operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address _account, address _operator) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) IsApprovedForAll(_account common.Address, _operator common.Address) (bool, error) {
	return _Nft1155.Contract.IsApprovedForAll(&_Nft1155.CallOpts, _account, _operator)
}

// IsMiner is a free data retrieval call binding the contract method 0x701b70ac.
//
// Solidity: function isMiner(address _miner) view returns(bool)
func (_Nft1155 *Nft1155Caller) IsMiner(opts *bind.CallOpts, _miner common.Address) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "isMiner", _miner)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsMiner is a free data retrieval call binding the contract method 0x701b70ac.
//
// Solidity: function isMiner(address _miner) view returns(bool)
func (_Nft1155 *Nft1155Session) IsMiner(_miner common.Address) (bool, error) {
	return _Nft1155.Contract.IsMiner(&_Nft1155.CallOpts, _miner)
}

// IsMiner is a free data retrieval call binding the contract method 0x701b70ac.
//
// Solidity: function isMiner(address _miner) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) IsMiner(_miner common.Address) (bool, error) {
	return _Nft1155.Contract.IsMiner(&_Nft1155.CallOpts, _miner)
}

// IsSuper is a free data retrieval call binding the contract method 0xc3c8df99.
//
// Solidity: function isSuper(address _super) view returns(bool)
func (_Nft1155 *Nft1155Caller) IsSuper(opts *bind.CallOpts, _super common.Address) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "isSuper", _super)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsSuper is a free data retrieval call binding the contract method 0xc3c8df99.
//
// Solidity: function isSuper(address _super) view returns(bool)
func (_Nft1155 *Nft1155Session) IsSuper(_super common.Address) (bool, error) {
	return _Nft1155.Contract.IsSuper(&_Nft1155.CallOpts, _super)
}

// IsSuper is a free data retrieval call binding the contract method 0xc3c8df99.
//
// Solidity: function isSuper(address _super) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) IsSuper(_super common.Address) (bool, error) {
	return _Nft1155.Contract.IsSuper(&_Nft1155.CallOpts, _super)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nft1155 *Nft1155Caller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nft1155 *Nft1155Session) Owner() (common.Address, error) {
	return _Nft1155.Contract.Owner(&_Nft1155.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Nft1155 *Nft1155CallerSession) Owner() (common.Address, error) {
	return _Nft1155.Contract.Owner(&_Nft1155.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Nft1155 *Nft1155Caller) RoyaltyInfo(opts *bind.CallOpts, _id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "royaltyInfo", _id, _price)

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

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Nft1155 *Nft1155Session) RoyaltyInfo(_id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	return _Nft1155.Contract.RoyaltyInfo(&_Nft1155.CallOpts, _id, _price)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _id, uint256 _price) view returns(address _receiver, uint256 _royalty)
func (_Nft1155 *Nft1155CallerSession) RoyaltyInfo(_id *big.Int, _price *big.Int) (struct {
	Receiver common.Address
	Royalty  *big.Int
}, error) {
	return _Nft1155.Contract.RoyaltyInfo(&_Nft1155.CallOpts, _id, _price)
}

// RoyaltyInfoBatch is a free data retrieval call binding the contract method 0xb645901c.
//
// Solidity: function royaltyInfoBatch(uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Nft1155 *Nft1155Caller) RoyaltyInfoBatch(opts *bind.CallOpts, _ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "royaltyInfoBatch", _ids, _price)

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

// RoyaltyInfoBatch is a free data retrieval call binding the contract method 0xb645901c.
//
// Solidity: function royaltyInfoBatch(uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Nft1155 *Nft1155Session) RoyaltyInfoBatch(_ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	return _Nft1155.Contract.RoyaltyInfoBatch(&_Nft1155.CallOpts, _ids, _price)
}

// RoyaltyInfoBatch is a free data retrieval call binding the contract method 0xb645901c.
//
// Solidity: function royaltyInfoBatch(uint256[] _ids, uint256 _price) view returns(address[] _receivers, uint256[] _royaltys)
func (_Nft1155 *Nft1155CallerSession) RoyaltyInfoBatch(_ids []*big.Int, _price *big.Int) (struct {
	Receivers []common.Address
	Royaltys  []*big.Int
}, error) {
	return _Nft1155.Contract.RoyaltyInfoBatch(&_Nft1155.CallOpts, _ids, _price)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_Nft1155 *Nft1155CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _Nft1155.Contract.SupportsInterface(&_Nft1155.CallOpts, interfaceId)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 _id) view returns(string)
func (_Nft1155 *Nft1155Caller) Uri(opts *bind.CallOpts, _id *big.Int) (string, error) {
	var out []interface{}
	err := _Nft1155.contract.Call(opts, &out, "uri", _id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 _id) view returns(string)
func (_Nft1155 *Nft1155Session) Uri(_id *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, _id)
}

// Uri is a free data retrieval call binding the contract method 0x0e89341c.
//
// Solidity: function uri(uint256 _id) view returns(string)
func (_Nft1155 *Nft1155CallerSession) Uri(_id *big.Int) (string, error) {
	return _Nft1155.Contract.Uri(&_Nft1155.CallOpts, _id)
}

// Mint is a paid mutator transaction binding the contract method 0x8031ae9b.
//
// Solidity: function mint(address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _data) returns()
func (_Nft1155 *Nft1155Transactor) Mint(opts *bind.TransactOpts, _to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mint", _to, _id, _amount, _royaltyRatio, _tokenURI, _data)
}

// Mint is a paid mutator transaction binding the contract method 0x8031ae9b.
//
// Solidity: function mint(address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _data) returns()
func (_Nft1155 *Nft1155Session) Mint(_to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, _to, _id, _amount, _royaltyRatio, _tokenURI, _data)
}

// Mint is a paid mutator transaction binding the contract method 0x8031ae9b.
//
// Solidity: function mint(address _to, uint256 _id, uint256 _amount, uint16 _royaltyRatio, string _tokenURI, bytes _data) returns()
func (_Nft1155 *Nft1155TransactorSession) Mint(_to common.Address, _id *big.Int, _amount *big.Int, _royaltyRatio uint16, _tokenURI string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.Mint(&_Nft1155.TransactOpts, _to, _id, _amount, _royaltyRatio, _tokenURI, _data)
}

// MintBatch is a paid mutator transaction binding the contract method 0x526858a1.
//
// Solidity: function mintBatch(address _to, uint256[] _ids, uint256[] _amounts, uint16[] _royaltyRatios, string[] _tokenURIs, bytes _data) returns()
func (_Nft1155 *Nft1155Transactor) MintBatch(opts *bind.TransactOpts, _to common.Address, _ids []*big.Int, _amounts []*big.Int, _royaltyRatios []uint16, _tokenURIs []string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "mintBatch", _to, _ids, _amounts, _royaltyRatios, _tokenURIs, _data)
}

// MintBatch is a paid mutator transaction binding the contract method 0x526858a1.
//
// Solidity: function mintBatch(address _to, uint256[] _ids, uint256[] _amounts, uint16[] _royaltyRatios, string[] _tokenURIs, bytes _data) returns()
func (_Nft1155 *Nft1155Session) MintBatch(_to common.Address, _ids []*big.Int, _amounts []*big.Int, _royaltyRatios []uint16, _tokenURIs []string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.MintBatch(&_Nft1155.TransactOpts, _to, _ids, _amounts, _royaltyRatios, _tokenURIs, _data)
}

// MintBatch is a paid mutator transaction binding the contract method 0x526858a1.
//
// Solidity: function mintBatch(address _to, uint256[] _ids, uint256[] _amounts, uint16[] _royaltyRatios, string[] _tokenURIs, bytes _data) returns()
func (_Nft1155 *Nft1155TransactorSession) MintBatch(_to common.Address, _ids []*big.Int, _amounts []*big.Int, _royaltyRatios []uint16, _tokenURIs []string, _data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.MintBatch(&_Nft1155.TransactOpts, _to, _ids, _amounts, _royaltyRatios, _tokenURIs, _data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nft1155 *Nft1155Transactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nft1155 *Nft1155Session) RenounceOwnership() (*types.Transaction, error) {
	return _Nft1155.Contract.RenounceOwnership(&_Nft1155.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Nft1155 *Nft1155TransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Nft1155.Contract.RenounceOwnership(&_Nft1155.TransactOpts)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeBatchTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeBatchTransferFrom", from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeBatchTransferFrom is a paid mutator transaction binding the contract method 0x2eb2c2d6.
//
// Solidity: function safeBatchTransferFrom(address from, address to, uint256[] ids, uint256[] amounts, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeBatchTransferFrom(from common.Address, to common.Address, ids []*big.Int, amounts []*big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeBatchTransferFrom(&_Nft1155.TransactOpts, from, to, ids, amounts, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155Transactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "safeTransferFrom", from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155Session) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, amount, data)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0xf242432a.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 id, uint256 amount, bytes data) returns()
func (_Nft1155 *Nft1155TransactorSession) SafeTransferFrom(from common.Address, to common.Address, id *big.Int, amount *big.Int, data []byte) (*types.Transaction, error) {
	return _Nft1155.Contract.SafeTransferFrom(&_Nft1155.TransactOpts, from, to, id, amount, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Transactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155Session) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_Nft1155 *Nft1155TransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _Nft1155.Contract.SetApprovalForAll(&_Nft1155.TransactOpts, operator, approved)
}

// SetLevel is a paid mutator transaction binding the contract method 0xe44a3169.
//
// Solidity: function setLevel(address _addr, uint8 _level) returns()
func (_Nft1155 *Nft1155Transactor) SetLevel(opts *bind.TransactOpts, _addr common.Address, _level uint8) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setLevel", _addr, _level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xe44a3169.
//
// Solidity: function setLevel(address _addr, uint8 _level) returns()
func (_Nft1155 *Nft1155Session) SetLevel(_addr common.Address, _level uint8) (*types.Transaction, error) {
	return _Nft1155.Contract.SetLevel(&_Nft1155.TransactOpts, _addr, _level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xe44a3169.
//
// Solidity: function setLevel(address _addr, uint8 _level) returns()
func (_Nft1155 *Nft1155TransactorSession) SetLevel(_addr common.Address, _level uint8) (*types.Transaction, error) {
	return _Nft1155.Contract.SetLevel(&_Nft1155.TransactOpts, _addr, _level)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x11e3d8d7.
//
// Solidity: function setRoyalty(uint256 _id, uint16 _ratio) returns()
func (_Nft1155 *Nft1155Transactor) SetRoyalty(opts *bind.TransactOpts, _id *big.Int, _ratio uint16) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setRoyalty", _id, _ratio)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x11e3d8d7.
//
// Solidity: function setRoyalty(uint256 _id, uint16 _ratio) returns()
func (_Nft1155 *Nft1155Session) SetRoyalty(_id *big.Int, _ratio uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.SetRoyalty(&_Nft1155.TransactOpts, _id, _ratio)
}

// SetRoyalty is a paid mutator transaction binding the contract method 0x11e3d8d7.
//
// Solidity: function setRoyalty(uint256 _id, uint16 _ratio) returns()
func (_Nft1155 *Nft1155TransactorSession) SetRoyalty(_id *big.Int, _ratio uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.SetRoyalty(&_Nft1155.TransactOpts, _id, _ratio)
}

// SetRoyaltyBatch is a paid mutator transaction binding the contract method 0x0d801835.
//
// Solidity: function setRoyaltyBatch(uint256[] _ids, uint16[] _ratios) returns()
func (_Nft1155 *Nft1155Transactor) SetRoyaltyBatch(opts *bind.TransactOpts, _ids []*big.Int, _ratios []uint16) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "setRoyaltyBatch", _ids, _ratios)
}

// SetRoyaltyBatch is a paid mutator transaction binding the contract method 0x0d801835.
//
// Solidity: function setRoyaltyBatch(uint256[] _ids, uint16[] _ratios) returns()
func (_Nft1155 *Nft1155Session) SetRoyaltyBatch(_ids []*big.Int, _ratios []uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.SetRoyaltyBatch(&_Nft1155.TransactOpts, _ids, _ratios)
}

// SetRoyaltyBatch is a paid mutator transaction binding the contract method 0x0d801835.
//
// Solidity: function setRoyaltyBatch(uint256[] _ids, uint16[] _ratios) returns()
func (_Nft1155 *Nft1155TransactorSession) SetRoyaltyBatch(_ids []*big.Int, _ratios []uint16) (*types.Transaction, error) {
	return _Nft1155.Contract.SetRoyaltyBatch(&_Nft1155.TransactOpts, _ids, _ratios)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nft1155 *Nft1155Transactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Nft1155.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nft1155 *Nft1155Session) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nft1155.Contract.TransferOwnership(&_Nft1155.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Nft1155 *Nft1155TransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Nft1155.Contract.TransferOwnership(&_Nft1155.TransactOpts, newOwner)
}

// Nft1155ApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Nft1155 contract.
type Nft1155ApprovalForAllIterator struct {
	Event *Nft1155ApprovalForAll // Event containing the contract specifics and raw log

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
func (it *Nft1155ApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ApprovalForAll)
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
		it.Event = new(Nft1155ApprovalForAll)
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
func (it *Nft1155ApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ApprovalForAll represents a ApprovalForAll event raised by the Nft1155 contract.
type Nft1155ApprovalForAll struct {
	Account  common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) FilterApprovalForAll(opts *bind.FilterOpts, account []common.Address, operator []common.Address) (*Nft1155ApprovalForAllIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155ApprovalForAllIterator{contract: _Nft1155.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *Nft1155ApprovalForAll, account []common.Address, operator []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ApprovalForAll", accountRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ApprovalForAll)
				if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed account, address indexed operator, bool approved)
func (_Nft1155 *Nft1155Filterer) ParseApprovalForAll(log types.Log) (*Nft1155ApprovalForAll, error) {
	event := new(Nft1155ApprovalForAll)
	if err := _Nft1155.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155ApprovalLevelIterator is returned from FilterApprovalLevel and is used to iterate over the raw logs and unpacked data for ApprovalLevel events raised by the Nft1155 contract.
type Nft1155ApprovalLevelIterator struct {
	Event *Nft1155ApprovalLevel // Event containing the contract specifics and raw log

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
func (it *Nft1155ApprovalLevelIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155ApprovalLevel)
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
		it.Event = new(Nft1155ApprovalLevel)
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
func (it *Nft1155ApprovalLevelIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155ApprovalLevelIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155ApprovalLevel represents a ApprovalLevel event raised by the Nft1155 contract.
type Nft1155ApprovalLevel struct {
	Addr  common.Address
	Level uint8
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterApprovalLevel is a free log retrieval operation binding the contract event 0xe3db383ed6be73b83dd8d70de0d7a91a57662effa97e0e262f3e67a5ce81dfae.
//
// Solidity: event ApprovalLevel(address indexed addr, uint8 level)
func (_Nft1155 *Nft1155Filterer) FilterApprovalLevel(opts *bind.FilterOpts, addr []common.Address) (*Nft1155ApprovalLevelIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "ApprovalLevel", addrRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155ApprovalLevelIterator{contract: _Nft1155.contract, event: "ApprovalLevel", logs: logs, sub: sub}, nil
}

// WatchApprovalLevel is a free log subscription operation binding the contract event 0xe3db383ed6be73b83dd8d70de0d7a91a57662effa97e0e262f3e67a5ce81dfae.
//
// Solidity: event ApprovalLevel(address indexed addr, uint8 level)
func (_Nft1155 *Nft1155Filterer) WatchApprovalLevel(opts *bind.WatchOpts, sink chan<- *Nft1155ApprovalLevel, addr []common.Address) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "ApprovalLevel", addrRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155ApprovalLevel)
				if err := _Nft1155.contract.UnpackLog(event, "ApprovalLevel", log); err != nil {
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

// ParseApprovalLevel is a log parse operation binding the contract event 0xe3db383ed6be73b83dd8d70de0d7a91a57662effa97e0e262f3e67a5ce81dfae.
//
// Solidity: event ApprovalLevel(address indexed addr, uint8 level)
func (_Nft1155 *Nft1155Filterer) ParseApprovalLevel(log types.Log) (*Nft1155ApprovalLevel, error) {
	event := new(Nft1155ApprovalLevel)
	if err := _Nft1155.contract.UnpackLog(event, "ApprovalLevel", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155OwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Nft1155 contract.
type Nft1155OwnershipTransferredIterator struct {
	Event *Nft1155OwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Nft1155OwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155OwnershipTransferred)
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
		it.Event = new(Nft1155OwnershipTransferred)
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
func (it *Nft1155OwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155OwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155OwnershipTransferred represents a OwnershipTransferred event raised by the Nft1155 contract.
type Nft1155OwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nft1155 *Nft1155Filterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Nft1155OwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155OwnershipTransferredIterator{contract: _Nft1155.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Nft1155 *Nft1155Filterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Nft1155OwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155OwnershipTransferred)
				if err := _Nft1155.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Nft1155 *Nft1155Filterer) ParseOwnershipTransferred(log types.Log) (*Nft1155OwnershipTransferred, error) {
	event := new(Nft1155OwnershipTransferred)
	if err := _Nft1155.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155RoyaltyIterator is returned from FilterRoyalty and is used to iterate over the raw logs and unpacked data for Royalty events raised by the Nft1155 contract.
type Nft1155RoyaltyIterator struct {
	Event *Nft1155Royalty // Event containing the contract specifics and raw log

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
func (it *Nft1155RoyaltyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155Royalty)
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
		it.Event = new(Nft1155Royalty)
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
func (it *Nft1155RoyaltyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155RoyaltyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155Royalty represents a Royalty event raised by the Nft1155 contract.
type Nft1155Royalty struct {
	Id       *big.Int
	Ratio    uint16
	Receiver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterRoyalty is a free log retrieval operation binding the contract event 0x611d12c0f8b2d9f4cfb23a30f560228db53e712dfcd34bb5b239e702efc2d22f.
//
// Solidity: event Royalty(uint256 id, uint16 ratio, address receiver)
func (_Nft1155 *Nft1155Filterer) FilterRoyalty(opts *bind.FilterOpts) (*Nft1155RoyaltyIterator, error) {

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "Royalty")
	if err != nil {
		return nil, err
	}
	return &Nft1155RoyaltyIterator{contract: _Nft1155.contract, event: "Royalty", logs: logs, sub: sub}, nil
}

// WatchRoyalty is a free log subscription operation binding the contract event 0x611d12c0f8b2d9f4cfb23a30f560228db53e712dfcd34bb5b239e702efc2d22f.
//
// Solidity: event Royalty(uint256 id, uint16 ratio, address receiver)
func (_Nft1155 *Nft1155Filterer) WatchRoyalty(opts *bind.WatchOpts, sink chan<- *Nft1155Royalty) (event.Subscription, error) {

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "Royalty")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155Royalty)
				if err := _Nft1155.contract.UnpackLog(event, "Royalty", log); err != nil {
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

// ParseRoyalty is a log parse operation binding the contract event 0x611d12c0f8b2d9f4cfb23a30f560228db53e712dfcd34bb5b239e702efc2d22f.
//
// Solidity: event Royalty(uint256 id, uint16 ratio, address receiver)
func (_Nft1155 *Nft1155Filterer) ParseRoyalty(log types.Log) (*Nft1155Royalty, error) {
	event := new(Nft1155Royalty)
	if err := _Nft1155.contract.UnpackLog(event, "Royalty", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferBatchIterator is returned from FilterTransferBatch and is used to iterate over the raw logs and unpacked data for TransferBatch events raised by the Nft1155 contract.
type Nft1155TransferBatchIterator struct {
	Event *Nft1155TransferBatch // Event containing the contract specifics and raw log

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
func (it *Nft1155TransferBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferBatch)
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
		it.Event = new(Nft1155TransferBatch)
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
func (it *Nft1155TransferBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferBatch represents a TransferBatch event raised by the Nft1155 contract.
type Nft1155TransferBatch struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Ids      []*big.Int
	Values   []*big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferBatch is a free log retrieval operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) FilterTransferBatch(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferBatchIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferBatchIterator{contract: _Nft1155.contract, event: "TransferBatch", logs: logs, sub: sub}, nil
}

// WatchTransferBatch is a free log subscription operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) WatchTransferBatch(opts *bind.WatchOpts, sink chan<- *Nft1155TransferBatch, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferBatch", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferBatch)
				if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
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

// ParseTransferBatch is a log parse operation binding the contract event 0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb.
//
// Solidity: event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] values)
func (_Nft1155 *Nft1155Filterer) ParseTransferBatch(log types.Log) (*Nft1155TransferBatch, error) {
	event := new(Nft1155TransferBatch)
	if err := _Nft1155.contract.UnpackLog(event, "TransferBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155TransferSingleIterator is returned from FilterTransferSingle and is used to iterate over the raw logs and unpacked data for TransferSingle events raised by the Nft1155 contract.
type Nft1155TransferSingleIterator struct {
	Event *Nft1155TransferSingle // Event containing the contract specifics and raw log

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
func (it *Nft1155TransferSingleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155TransferSingle)
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
		it.Event = new(Nft1155TransferSingle)
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
func (it *Nft1155TransferSingleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155TransferSingleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155TransferSingle represents a TransferSingle event raised by the Nft1155 contract.
type Nft1155TransferSingle struct {
	Operator common.Address
	From     common.Address
	To       common.Address
	Id       *big.Int
	Value    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferSingle is a free log retrieval operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) FilterTransferSingle(opts *bind.FilterOpts, operator []common.Address, from []common.Address, to []common.Address) (*Nft1155TransferSingleIterator, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155TransferSingleIterator{contract: _Nft1155.contract, event: "TransferSingle", logs: logs, sub: sub}, nil
}

// WatchTransferSingle is a free log subscription operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) WatchTransferSingle(opts *bind.WatchOpts, sink chan<- *Nft1155TransferSingle, operator []common.Address, from []common.Address, to []common.Address) (event.Subscription, error) {

	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "TransferSingle", operatorRule, fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155TransferSingle)
				if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
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

// ParseTransferSingle is a log parse operation binding the contract event 0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62.
//
// Solidity: event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 value)
func (_Nft1155 *Nft1155Filterer) ParseTransferSingle(log types.Log) (*Nft1155TransferSingle, error) {
	event := new(Nft1155TransferSingle)
	if err := _Nft1155.contract.UnpackLog(event, "TransferSingle", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Nft1155URIIterator is returned from FilterURI and is used to iterate over the raw logs and unpacked data for URI events raised by the Nft1155 contract.
type Nft1155URIIterator struct {
	Event *Nft1155URI // Event containing the contract specifics and raw log

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
func (it *Nft1155URIIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Nft1155URI)
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
		it.Event = new(Nft1155URI)
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
func (it *Nft1155URIIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Nft1155URIIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Nft1155URI represents a URI event raised by the Nft1155 contract.
type Nft1155URI struct {
	Value string
	Id    *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterURI is a free log retrieval operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) FilterURI(opts *bind.FilterOpts, id []*big.Int) (*Nft1155URIIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.FilterLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return &Nft1155URIIterator{contract: _Nft1155.contract, event: "URI", logs: logs, sub: sub}, nil
}

// WatchURI is a free log subscription operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) WatchURI(opts *bind.WatchOpts, sink chan<- *Nft1155URI, id []*big.Int) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Nft1155.contract.WatchLogs(opts, "URI", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Nft1155URI)
				if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
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

// ParseURI is a log parse operation binding the contract event 0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b.
//
// Solidity: event URI(string value, uint256 indexed id)
func (_Nft1155 *Nft1155Filterer) ParseURI(log types.Log) (*Nft1155URI, error) {
	event := new(Nft1155URI)
	if err := _Nft1155.contract.UnpackLog(event, "URI", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
