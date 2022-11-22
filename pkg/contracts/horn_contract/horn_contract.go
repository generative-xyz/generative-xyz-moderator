// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package horn_contract

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
	_ = abi.ConvertType
)

// HORNSHorn is an auto generated low-level Go binding around an user-defined struct.
type HORNSHorn struct {
	Nation       string
	PalletTop    string
	PalletBottom string
}

// HornContractMetaData contains all meta data concerning the HornContract contract.
var HornContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_algorithm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_limit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_tokenAddrErc721\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"changeBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newP\",\"type\":\"address\"}],\"name\":\"changeParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sweet\",\"type\":\"address\"}],\"name\":\"changeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getPaletteBottom\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getPaletteTop\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getParamValues\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"nation\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"palletTop\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"palletBottom\",\"type\":\"string\"}],\"internalType\":\"structHORNS.Horn\",\"name\":\"horn\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenIdGated\",\"type\":\"uint256\"}],\"name\":\"mintByToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"algo\",\"type\":\"string\"}],\"name\":\"setAlgo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"setLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// HornContractABI is the input ABI used to generate the binding from.
// Deprecated: Use HornContractMetaData.ABI instead.
var HornContractABI = HornContractMetaData.ABI

// HornContract is an auto generated Go binding around an Ethereum contract.
type HornContract struct {
	HornContractCaller     // Read-only binding to the contract
	HornContractTransactor // Write-only binding to the contract
	HornContractFilterer   // Log filterer for contract events
}

// HornContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type HornContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HornContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type HornContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HornContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type HornContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// HornContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type HornContractSession struct {
	Contract     *HornContract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// HornContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type HornContractCallerSession struct {
	Contract *HornContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// HornContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type HornContractTransactorSession struct {
	Contract     *HornContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// HornContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type HornContractRaw struct {
	Contract *HornContract // Generic contract binding to access the raw methods on
}

// HornContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type HornContractCallerRaw struct {
	Contract *HornContractCaller // Generic read-only contract binding to access the raw methods on
}

// HornContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type HornContractTransactorRaw struct {
	Contract *HornContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewHornContract creates a new instance of HornContract, bound to a specific deployed contract.
func NewHornContract(address common.Address, backend bind.ContractBackend) (*HornContract, error) {
	contract, err := bindHornContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &HornContract{HornContractCaller: HornContractCaller{contract: contract}, HornContractTransactor: HornContractTransactor{contract: contract}, HornContractFilterer: HornContractFilterer{contract: contract}}, nil
}

// NewHornContractCaller creates a new read-only instance of HornContract, bound to a specific deployed contract.
func NewHornContractCaller(address common.Address, caller bind.ContractCaller) (*HornContractCaller, error) {
	contract, err := bindHornContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &HornContractCaller{contract: contract}, nil
}

// NewHornContractTransactor creates a new write-only instance of HornContract, bound to a specific deployed contract.
func NewHornContractTransactor(address common.Address, transactor bind.ContractTransactor) (*HornContractTransactor, error) {
	contract, err := bindHornContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &HornContractTransactor{contract: contract}, nil
}

// NewHornContractFilterer creates a new log filterer instance of HornContract, bound to a specific deployed contract.
func NewHornContractFilterer(address common.Address, filterer bind.ContractFilterer) (*HornContractFilterer, error) {
	contract, err := bindHornContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &HornContractFilterer{contract: contract}, nil
}

// bindHornContract binds a generic wrapper to an already deployed contract.
func bindHornContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := HornContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HornContract *HornContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HornContract.Contract.HornContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HornContract *HornContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.Contract.HornContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HornContract *HornContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HornContract.Contract.HornContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_HornContract *HornContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _HornContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_HornContract *HornContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_HornContract *HornContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _HornContract.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_HornContract *HornContractCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_HornContract *HornContractSession) Admin() (common.Address, error) {
	return _HornContract.Contract.Admin(&_HornContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_HornContract *HornContractCallerSession) Admin() (common.Address, error) {
	return _HornContract.Contract.Admin(&_HornContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_HornContract *HornContractCaller) Algorithm(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_algorithm")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_HornContract *HornContractSession) Algorithm() (string, error) {
	return _HornContract.Contract.Algorithm(&_HornContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_HornContract *HornContractCallerSession) Algorithm() (string, error) {
	return _HornContract.Contract.Algorithm(&_HornContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_HornContract *HornContractCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_HornContract *HornContractSession) Counter() (*big.Int, error) {
	return _HornContract.Contract.Counter(&_HornContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_HornContract *HornContractCallerSession) Counter() (*big.Int, error) {
	return _HornContract.Contract.Counter(&_HornContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_HornContract *HornContractCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_HornContract *HornContractSession) Fee() (*big.Int, error) {
	return _HornContract.Contract.Fee(&_HornContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_HornContract *HornContractCallerSession) Fee() (*big.Int, error) {
	return _HornContract.Contract.Fee(&_HornContract.CallOpts)
}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_HornContract *HornContractCaller) Limit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_limit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_HornContract *HornContractSession) Limit() (*big.Int, error) {
	return _HornContract.Contract.Limit(&_HornContract.CallOpts)
}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_HornContract *HornContractCallerSession) Limit() (*big.Int, error) {
	return _HornContract.Contract.Limit(&_HornContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_HornContract *HornContractCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_HornContract *HornContractSession) ParamsAddress() (common.Address, error) {
	return _HornContract.Contract.ParamsAddress(&_HornContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_HornContract *HornContractCallerSession) ParamsAddress() (common.Address, error) {
	return _HornContract.Contract.ParamsAddress(&_HornContract.CallOpts)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_HornContract *HornContractCaller) TokenAddrErc721(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_tokenAddrErc721")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_HornContract *HornContractSession) TokenAddrErc721() (common.Address, error) {
	return _HornContract.Contract.TokenAddrErc721(&_HornContract.CallOpts)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_HornContract *HornContractCallerSession) TokenAddrErc721() (common.Address, error) {
	return _HornContract.Contract.TokenAddrErc721(&_HornContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_HornContract *HornContractCaller) Uri(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "_uri")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_HornContract *HornContractSession) Uri() (string, error) {
	return _HornContract.Contract.Uri(&_HornContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_HornContract *HornContractCallerSession) Uri() (string, error) {
	return _HornContract.Contract.Uri(&_HornContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_HornContract *HornContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_HornContract *HornContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _HornContract.Contract.BalanceOf(&_HornContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_HornContract *HornContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _HornContract.Contract.BalanceOf(&_HornContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_HornContract *HornContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_HornContract *HornContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _HornContract.Contract.GetApproved(&_HornContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_HornContract *HornContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _HornContract.Contract.GetApproved(&_HornContract.CallOpts, tokenId)
}

// GetPaletteBottom is a free data retrieval call binding the contract method 0x6a2c3995.
//
// Solidity: function getPaletteBottom(uint256 id) view returns(string)
func (_HornContract *HornContractCaller) GetPaletteBottom(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "getPaletteBottom", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetPaletteBottom is a free data retrieval call binding the contract method 0x6a2c3995.
//
// Solidity: function getPaletteBottom(uint256 id) view returns(string)
func (_HornContract *HornContractSession) GetPaletteBottom(id *big.Int) (string, error) {
	return _HornContract.Contract.GetPaletteBottom(&_HornContract.CallOpts, id)
}

// GetPaletteBottom is a free data retrieval call binding the contract method 0x6a2c3995.
//
// Solidity: function getPaletteBottom(uint256 id) view returns(string)
func (_HornContract *HornContractCallerSession) GetPaletteBottom(id *big.Int) (string, error) {
	return _HornContract.Contract.GetPaletteBottom(&_HornContract.CallOpts, id)
}

// GetPaletteTop is a free data retrieval call binding the contract method 0x5d36b3ea.
//
// Solidity: function getPaletteTop(uint256 id) view returns(string)
func (_HornContract *HornContractCaller) GetPaletteTop(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "getPaletteTop", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetPaletteTop is a free data retrieval call binding the contract method 0x5d36b3ea.
//
// Solidity: function getPaletteTop(uint256 id) view returns(string)
func (_HornContract *HornContractSession) GetPaletteTop(id *big.Int) (string, error) {
	return _HornContract.Contract.GetPaletteTop(&_HornContract.CallOpts, id)
}

// GetPaletteTop is a free data retrieval call binding the contract method 0x5d36b3ea.
//
// Solidity: function getPaletteTop(uint256 id) view returns(string)
func (_HornContract *HornContractCallerSession) GetPaletteTop(id *big.Int) (string, error) {
	return _HornContract.Contract.GetPaletteTop(&_HornContract.CallOpts, id)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string) horn)
func (_HornContract *HornContractCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) (HORNSHorn, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new(HORNSHorn), err
	}

	out0 := *abi.ConvertType(out[0], new(HORNSHorn)).(*HORNSHorn)

	return out0, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string) horn)
func (_HornContract *HornContractSession) GetParamValues(tokenId *big.Int) (HORNSHorn, error) {
	return _HornContract.Contract.GetParamValues(&_HornContract.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string) horn)
func (_HornContract *HornContractCallerSession) GetParamValues(tokenId *big.Int) (HORNSHorn, error) {
	return _HornContract.Contract.GetParamValues(&_HornContract.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_HornContract *HornContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_HornContract *HornContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _HornContract.Contract.IsApprovedForAll(&_HornContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_HornContract *HornContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _HornContract.Contract.IsApprovedForAll(&_HornContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_HornContract *HornContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_HornContract *HornContractSession) Name() (string, error) {
	return _HornContract.Contract.Name(&_HornContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_HornContract *HornContractCallerSession) Name() (string, error) {
	return _HornContract.Contract.Name(&_HornContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_HornContract *HornContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_HornContract *HornContractSession) Owner() (common.Address, error) {
	return _HornContract.Contract.Owner(&_HornContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_HornContract *HornContractCallerSession) Owner() (common.Address, error) {
	return _HornContract.Contract.Owner(&_HornContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_HornContract *HornContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_HornContract *HornContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _HornContract.Contract.OwnerOf(&_HornContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_HornContract *HornContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _HornContract.Contract.OwnerOf(&_HornContract.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_HornContract *HornContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_HornContract *HornContractSession) Paused() (bool, error) {
	return _HornContract.Contract.Paused(&_HornContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_HornContract *HornContractCallerSession) Paused() (bool, error) {
	return _HornContract.Contract.Paused(&_HornContract.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_HornContract *HornContractCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

	outstruct := new(struct {
		Receiver      common.Address
		RoyaltyAmount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Receiver = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.RoyaltyAmount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_HornContract *HornContractSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _HornContract.Contract.RoyaltyInfo(&_HornContract.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_HornContract *HornContractCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _HornContract.Contract.RoyaltyInfo(&_HornContract.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HornContract *HornContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HornContract *HornContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HornContract.Contract.SupportsInterface(&_HornContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_HornContract *HornContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _HornContract.Contract.SupportsInterface(&_HornContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_HornContract *HornContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_HornContract *HornContractSession) Symbol() (string, error) {
	return _HornContract.Contract.Symbol(&_HornContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_HornContract *HornContractCallerSession) Symbol() (string, error) {
	return _HornContract.Contract.Symbol(&_HornContract.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_HornContract *HornContractCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _HornContract.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_HornContract *HornContractSession) TokenURI(tokenId *big.Int) (string, error) {
	return _HornContract.Contract.TokenURI(&_HornContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_HornContract *HornContractCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _HornContract.Contract.TokenURI(&_HornContract.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_HornContract *HornContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.Approve(&_HornContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.Approve(&_HornContract.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_HornContract *HornContractTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_HornContract *HornContractSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeAdmin(&_HornContract.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_HornContract *HornContractTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeAdmin(&_HornContract.TransactOpts, newAdm)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_HornContract *HornContractTransactor) ChangeBaseURI(opts *bind.TransactOpts, baseURI string) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "changeBaseURI", baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_HornContract *HornContractSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeBaseURI(&_HornContract.TransactOpts, baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_HornContract *HornContractTransactorSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeBaseURI(&_HornContract.TransactOpts, baseURI)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_HornContract *HornContractTransactor) ChangeParam(opts *bind.TransactOpts, newP common.Address) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "changeParam", newP)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_HornContract *HornContractSession) ChangeParam(newP common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeParam(&_HornContract.TransactOpts, newP)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_HornContract *HornContractTransactorSession) ChangeParam(newP common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeParam(&_HornContract.TransactOpts, newP)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_HornContract *HornContractTransactor) ChangeToken(opts *bind.TransactOpts, sweet common.Address) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "changeToken", sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_HornContract *HornContractSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeToken(&_HornContract.TransactOpts, sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_HornContract *HornContractTransactorSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.ChangeToken(&_HornContract.TransactOpts, sweet)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_HornContract *HornContractTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "initialize", name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_HornContract *HornContractSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.Initialize(&_HornContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_HornContract *HornContractTransactorSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.Initialize(&_HornContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_HornContract *HornContractTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_HornContract *HornContractSession) Mint() (*types.Transaction, error) {
	return _HornContract.Contract.Mint(&_HornContract.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_HornContract *HornContractTransactorSession) Mint() (*types.Transaction, error) {
	return _HornContract.Contract.Mint(&_HornContract.TransactOpts)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_HornContract *HornContractTransactor) MintByToken(opts *bind.TransactOpts, tokenIdGated *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "mintByToken", tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_HornContract *HornContractSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.MintByToken(&_HornContract.TransactOpts, tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_HornContract *HornContractTransactorSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.MintByToken(&_HornContract.TransactOpts, tokenIdGated)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_HornContract *HornContractTransactor) OwnerMint(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "ownerMint", id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_HornContract *HornContractSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.OwnerMint(&_HornContract.TransactOpts, id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_HornContract *HornContractTransactorSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.OwnerMint(&_HornContract.TransactOpts, id)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_HornContract *HornContractTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_HornContract *HornContractSession) Pause() (*types.Transaction, error) {
	return _HornContract.Contract.Pause(&_HornContract.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_HornContract *HornContractTransactorSession) Pause() (*types.Transaction, error) {
	return _HornContract.Contract.Pause(&_HornContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_HornContract *HornContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_HornContract *HornContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _HornContract.Contract.RenounceOwnership(&_HornContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_HornContract *HornContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _HornContract.Contract.RenounceOwnership(&_HornContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SafeTransferFrom(&_HornContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SafeTransferFrom(&_HornContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_HornContract *HornContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_HornContract *HornContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _HornContract.Contract.SafeTransferFrom0(&_HornContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_HornContract *HornContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _HornContract.Contract.SafeTransferFrom0(&_HornContract.TransactOpts, from, to, tokenId, data)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_HornContract *HornContractTransactor) SetAlgo(opts *bind.TransactOpts, algo string) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "setAlgo", algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_HornContract *HornContractSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _HornContract.Contract.SetAlgo(&_HornContract.TransactOpts, algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_HornContract *HornContractTransactorSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _HornContract.Contract.SetAlgo(&_HornContract.TransactOpts, algo)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_HornContract *HornContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_HornContract *HornContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _HornContract.Contract.SetApprovalForAll(&_HornContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_HornContract *HornContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _HornContract.Contract.SetApprovalForAll(&_HornContract.TransactOpts, operator, approved)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_HornContract *HornContractTransactor) SetFee(opts *bind.TransactOpts, fee *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "setFee", fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_HornContract *HornContractSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SetFee(&_HornContract.TransactOpts, fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_HornContract *HornContractTransactorSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SetFee(&_HornContract.TransactOpts, fee)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_HornContract *HornContractTransactor) SetLimit(opts *bind.TransactOpts, limit *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "setLimit", limit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_HornContract *HornContractSession) SetLimit(limit *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SetLimit(&_HornContract.TransactOpts, limit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_HornContract *HornContractTransactorSession) SetLimit(limit *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.SetLimit(&_HornContract.TransactOpts, limit)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.TransferFrom(&_HornContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_HornContract *HornContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _HornContract.Contract.TransferFrom(&_HornContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_HornContract *HornContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_HornContract *HornContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.TransferOwnership(&_HornContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_HornContract *HornContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _HornContract.Contract.TransferOwnership(&_HornContract.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_HornContract *HornContractTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_HornContract *HornContractSession) Unpause() (*types.Transaction, error) {
	return _HornContract.Contract.Unpause(&_HornContract.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_HornContract *HornContractTransactorSession) Unpause() (*types.Transaction, error) {
	return _HornContract.Contract.Unpause(&_HornContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_HornContract *HornContractTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _HornContract.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_HornContract *HornContractSession) Withdraw() (*types.Transaction, error) {
	return _HornContract.Contract.Withdraw(&_HornContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_HornContract *HornContractTransactorSession) Withdraw() (*types.Transaction, error) {
	return _HornContract.Contract.Withdraw(&_HornContract.TransactOpts)
}

// HornContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the HornContract contract.
type HornContractApprovalIterator struct {
	Event *HornContractApproval // Event containing the contract specifics and raw log

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
func (it *HornContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractApproval)
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
		it.Event = new(HornContractApproval)
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
func (it *HornContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractApproval represents a Approval event raised by the HornContract contract.
type HornContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*HornContractApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &HornContractApprovalIterator{contract: _HornContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *HornContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractApproval)
				if err := _HornContract.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) ParseApproval(log types.Log) (*HornContractApproval, error) {
	event := new(HornContractApproval)
	if err := _HornContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the HornContract contract.
type HornContractApprovalForAllIterator struct {
	Event *HornContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *HornContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractApprovalForAll)
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
		it.Event = new(HornContractApprovalForAll)
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
func (it *HornContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractApprovalForAll represents a ApprovalForAll event raised by the HornContract contract.
type HornContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_HornContract *HornContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*HornContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &HornContractApprovalForAllIterator{contract: _HornContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_HornContract *HornContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *HornContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractApprovalForAll)
				if err := _HornContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_HornContract *HornContractFilterer) ParseApprovalForAll(log types.Log) (*HornContractApprovalForAll, error) {
	event := new(HornContractApprovalForAll)
	if err := _HornContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the HornContract contract.
type HornContractInitializedIterator struct {
	Event *HornContractInitialized // Event containing the contract specifics and raw log

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
func (it *HornContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractInitialized)
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
		it.Event = new(HornContractInitialized)
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
func (it *HornContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractInitialized represents a Initialized event raised by the HornContract contract.
type HornContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_HornContract *HornContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*HornContractInitializedIterator, error) {

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &HornContractInitializedIterator{contract: _HornContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_HornContract *HornContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *HornContractInitialized) (event.Subscription, error) {

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractInitialized)
				if err := _HornContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_HornContract *HornContractFilterer) ParseInitialized(log types.Log) (*HornContractInitialized, error) {
	event := new(HornContractInitialized)
	if err := _HornContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the HornContract contract.
type HornContractOwnershipTransferredIterator struct {
	Event *HornContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *HornContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractOwnershipTransferred)
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
		it.Event = new(HornContractOwnershipTransferred)
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
func (it *HornContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractOwnershipTransferred represents a OwnershipTransferred event raised by the HornContract contract.
type HornContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_HornContract *HornContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*HornContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &HornContractOwnershipTransferredIterator{contract: _HornContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_HornContract *HornContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *HornContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractOwnershipTransferred)
				if err := _HornContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_HornContract *HornContractFilterer) ParseOwnershipTransferred(log types.Log) (*HornContractOwnershipTransferred, error) {
	event := new(HornContractOwnershipTransferred)
	if err := _HornContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the HornContract contract.
type HornContractPausedIterator struct {
	Event *HornContractPaused // Event containing the contract specifics and raw log

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
func (it *HornContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractPaused)
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
		it.Event = new(HornContractPaused)
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
func (it *HornContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractPaused represents a Paused event raised by the HornContract contract.
type HornContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_HornContract *HornContractFilterer) FilterPaused(opts *bind.FilterOpts) (*HornContractPausedIterator, error) {

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &HornContractPausedIterator{contract: _HornContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_HornContract *HornContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *HornContractPaused) (event.Subscription, error) {

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractPaused)
				if err := _HornContract.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_HornContract *HornContractFilterer) ParsePaused(log types.Log) (*HornContractPaused, error) {
	event := new(HornContractPaused)
	if err := _HornContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the HornContract contract.
type HornContractTransferIterator struct {
	Event *HornContractTransfer // Event containing the contract specifics and raw log

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
func (it *HornContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractTransfer)
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
		it.Event = new(HornContractTransfer)
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
func (it *HornContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractTransfer represents a Transfer event raised by the HornContract contract.
type HornContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*HornContractTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &HornContractTransferIterator{contract: _HornContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *HornContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractTransfer)
				if err := _HornContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_HornContract *HornContractFilterer) ParseTransfer(log types.Log) (*HornContractTransfer, error) {
	event := new(HornContractTransfer)
	if err := _HornContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// HornContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the HornContract contract.
type HornContractUnpausedIterator struct {
	Event *HornContractUnpaused // Event containing the contract specifics and raw log

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
func (it *HornContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(HornContractUnpaused)
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
		it.Event = new(HornContractUnpaused)
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
func (it *HornContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *HornContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// HornContractUnpaused represents a Unpaused event raised by the HornContract contract.
type HornContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_HornContract *HornContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*HornContractUnpausedIterator, error) {

	logs, sub, err := _HornContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &HornContractUnpausedIterator{contract: _HornContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_HornContract *HornContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *HornContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _HornContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(HornContractUnpaused)
				if err := _HornContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_HornContract *HornContractFilterer) ParseUnpaused(log types.Log) (*HornContractUnpaused, error) {
	event := new(HornContractUnpaused)
	if err := _HornContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
