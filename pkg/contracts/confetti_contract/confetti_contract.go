// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package confetti_contract

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

// CONFETTIConfetti is an auto generated low-level Go binding around an user-defined struct.
type CONFETTIConfetti struct {
	ShapeCanon       string
	ShapeConfetti    string
	PalletteCanon    [4]string
	PalletteConfetti [2]string
}

// ConfettiContractMetaData contains all meta data concerning the ConfettiContract contract.
var ConfettiContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"OperatorNotAllowed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_algorithm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_counter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_fee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_limit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramsAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_tokenAddrErc721\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_uri\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"baseURI\",\"type\":\"string\"}],\"name\":\"changeBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newP\",\"type\":\"address\"}],\"name\":\"changeParam\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sweet\",\"type\":\"address\"}],\"name\":\"changeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getPaletteCanon\",\"outputs\":[{\"internalType\":\"string[4]\",\"name\":\"\",\"type\":\"string[4]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getPaletteConfetti\",\"outputs\":[{\"internalType\":\"string[2]\",\"name\":\"\",\"type\":\"string[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getParamValues\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"shapeCanon\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"shapeConfetti\",\"type\":\"string\"},{\"internalType\":\"string[4]\",\"name\":\"palletteCanon\",\"type\":\"string[4]\"},{\"internalType\":\"string[2]\",\"name\":\"palletteConfetti\",\"type\":\"string[2]\"}],\"internalType\":\"structCONFETTI.Confetti\",\"name\":\"confetti\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramsAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenIdGated\",\"type\":\"uint256\"}],\"name\":\"mintByToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_salePrice\",\"type\":\"uint256\"}],\"name\":\"royaltyInfo\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"algo\",\"type\":\"string\"}],\"name\":\"setAlgo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"setFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"setLimit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ConfettiContractABI is the input ABI used to generate the binding from.
// Deprecated: Use ConfettiContractMetaData.ABI instead.
var ConfettiContractABI = ConfettiContractMetaData.ABI

// ConfettiContract is an auto generated Go binding around an Ethereum contract.
type ConfettiContract struct {
	ConfettiContractCaller     // Read-only binding to the contract
	ConfettiContractTransactor // Write-only binding to the contract
	ConfettiContractFilterer   // Log filterer for contract events
}

// ConfettiContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConfettiContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfettiContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConfettiContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfettiContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConfettiContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConfettiContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConfettiContractSession struct {
	Contract     *ConfettiContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConfettiContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConfettiContractCallerSession struct {
	Contract *ConfettiContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ConfettiContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConfettiContractTransactorSession struct {
	Contract     *ConfettiContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ConfettiContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConfettiContractRaw struct {
	Contract *ConfettiContract // Generic contract binding to access the raw methods on
}

// ConfettiContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConfettiContractCallerRaw struct {
	Contract *ConfettiContractCaller // Generic read-only contract binding to access the raw methods on
}

// ConfettiContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConfettiContractTransactorRaw struct {
	Contract *ConfettiContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConfettiContract creates a new instance of ConfettiContract, bound to a specific deployed contract.
func NewConfettiContract(address common.Address, backend bind.ContractBackend) (*ConfettiContract, error) {
	contract, err := bindConfettiContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConfettiContract{ConfettiContractCaller: ConfettiContractCaller{contract: contract}, ConfettiContractTransactor: ConfettiContractTransactor{contract: contract}, ConfettiContractFilterer: ConfettiContractFilterer{contract: contract}}, nil
}

// NewConfettiContractCaller creates a new read-only instance of ConfettiContract, bound to a specific deployed contract.
func NewConfettiContractCaller(address common.Address, caller bind.ContractCaller) (*ConfettiContractCaller, error) {
	contract, err := bindConfettiContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractCaller{contract: contract}, nil
}

// NewConfettiContractTransactor creates a new write-only instance of ConfettiContract, bound to a specific deployed contract.
func NewConfettiContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ConfettiContractTransactor, error) {
	contract, err := bindConfettiContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractTransactor{contract: contract}, nil
}

// NewConfettiContractFilterer creates a new log filterer instance of ConfettiContract, bound to a specific deployed contract.
func NewConfettiContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ConfettiContractFilterer, error) {
	contract, err := bindConfettiContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractFilterer{contract: contract}, nil
}

// bindConfettiContract binds a generic wrapper to an already deployed contract.
func bindConfettiContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ConfettiContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfettiContract *ConfettiContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfettiContract.Contract.ConfettiContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfettiContract *ConfettiContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ConfettiContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfettiContract *ConfettiContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ConfettiContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConfettiContract *ConfettiContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConfettiContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConfettiContract *ConfettiContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConfettiContract *ConfettiContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConfettiContract.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_ConfettiContract *ConfettiContractCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_ConfettiContract *ConfettiContractSession) Admin() (common.Address, error) {
	return _ConfettiContract.Contract.Admin(&_ConfettiContract.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) Admin() (common.Address, error) {
	return _ConfettiContract.Contract.Admin(&_ConfettiContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_ConfettiContract *ConfettiContractCaller) Algorithm(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_algorithm")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_ConfettiContract *ConfettiContractSession) Algorithm() (string, error) {
	return _ConfettiContract.Contract.Algorithm(&_ConfettiContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0xa1bb1fcc.
//
// Solidity: function _algorithm() view returns(string)
func (_ConfettiContract *ConfettiContractCallerSession) Algorithm() (string, error) {
	return _ConfettiContract.Contract.Algorithm(&_ConfettiContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_ConfettiContract *ConfettiContractCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_counter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_ConfettiContract *ConfettiContractSession) Counter() (*big.Int, error) {
	return _ConfettiContract.Contract.Counter(&_ConfettiContract.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x7cd49fde.
//
// Solidity: function _counter() view returns(uint256)
func (_ConfettiContract *ConfettiContractCallerSession) Counter() (*big.Int, error) {
	return _ConfettiContract.Contract.Counter(&_ConfettiContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_ConfettiContract *ConfettiContractCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_ConfettiContract *ConfettiContractSession) Fee() (*big.Int, error) {
	return _ConfettiContract.Contract.Fee(&_ConfettiContract.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xc5b37c22.
//
// Solidity: function _fee() view returns(uint256)
func (_ConfettiContract *ConfettiContractCallerSession) Fee() (*big.Int, error) {
	return _ConfettiContract.Contract.Fee(&_ConfettiContract.CallOpts)
}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_ConfettiContract *ConfettiContractCaller) Limit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_limit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_ConfettiContract *ConfettiContractSession) Limit() (*big.Int, error) {
	return _ConfettiContract.Contract.Limit(&_ConfettiContract.CallOpts)
}

// Limit is a free data retrieval call binding the contract method 0xdf2fb92c.
//
// Solidity: function _limit() view returns(uint256)
func (_ConfettiContract *ConfettiContractCallerSession) Limit() (*big.Int, error) {
	return _ConfettiContract.Contract.Limit(&_ConfettiContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_ConfettiContract *ConfettiContractCaller) ParamsAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_paramsAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_ConfettiContract *ConfettiContractSession) ParamsAddress() (common.Address, error) {
	return _ConfettiContract.Contract.ParamsAddress(&_ConfettiContract.CallOpts)
}

// ParamsAddress is a free data retrieval call binding the contract method 0xadfc7dae.
//
// Solidity: function _paramsAddress() view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) ParamsAddress() (common.Address, error) {
	return _ConfettiContract.Contract.ParamsAddress(&_ConfettiContract.CallOpts)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_ConfettiContract *ConfettiContractCaller) TokenAddrErc721(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_tokenAddrErc721")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_ConfettiContract *ConfettiContractSession) TokenAddrErc721() (common.Address, error) {
	return _ConfettiContract.Contract.TokenAddrErc721(&_ConfettiContract.CallOpts)
}

// TokenAddrErc721 is a free data retrieval call binding the contract method 0xe653c650.
//
// Solidity: function _tokenAddrErc721() view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) TokenAddrErc721() (common.Address, error) {
	return _ConfettiContract.Contract.TokenAddrErc721(&_ConfettiContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_ConfettiContract *ConfettiContractCaller) Uri(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "_uri")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_ConfettiContract *ConfettiContractSession) Uri() (string, error) {
	return _ConfettiContract.Contract.Uri(&_ConfettiContract.CallOpts)
}

// Uri is a free data retrieval call binding the contract method 0x0dccc9ad.
//
// Solidity: function _uri() view returns(string)
func (_ConfettiContract *ConfettiContractCallerSession) Uri() (string, error) {
	return _ConfettiContract.Contract.Uri(&_ConfettiContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ConfettiContract *ConfettiContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ConfettiContract *ConfettiContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ConfettiContract.Contract.BalanceOf(&_ConfettiContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_ConfettiContract *ConfettiContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _ConfettiContract.Contract.BalanceOf(&_ConfettiContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ConfettiContract.Contract.GetApproved(&_ConfettiContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _ConfettiContract.Contract.GetApproved(&_ConfettiContract.CallOpts, tokenId)
}

// GetPaletteCanon is a free data retrieval call binding the contract method 0x7d5c0120.
//
// Solidity: function getPaletteCanon(uint256 id) view returns(string[4])
func (_ConfettiContract *ConfettiContractCaller) GetPaletteCanon(opts *bind.CallOpts, id *big.Int) ([4]string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "getPaletteCanon", id)

	if err != nil {
		return *new([4]string), err
	}

	out0 := *abi.ConvertType(out[0], new([4]string)).(*[4]string)

	return out0, err

}

// GetPaletteCanon is a free data retrieval call binding the contract method 0x7d5c0120.
//
// Solidity: function getPaletteCanon(uint256 id) view returns(string[4])
func (_ConfettiContract *ConfettiContractSession) GetPaletteCanon(id *big.Int) ([4]string, error) {
	return _ConfettiContract.Contract.GetPaletteCanon(&_ConfettiContract.CallOpts, id)
}

// GetPaletteCanon is a free data retrieval call binding the contract method 0x7d5c0120.
//
// Solidity: function getPaletteCanon(uint256 id) view returns(string[4])
func (_ConfettiContract *ConfettiContractCallerSession) GetPaletteCanon(id *big.Int) ([4]string, error) {
	return _ConfettiContract.Contract.GetPaletteCanon(&_ConfettiContract.CallOpts, id)
}

// GetPaletteConfetti is a free data retrieval call binding the contract method 0xcbfb504a.
//
// Solidity: function getPaletteConfetti(uint256 id) view returns(string[2])
func (_ConfettiContract *ConfettiContractCaller) GetPaletteConfetti(opts *bind.CallOpts, id *big.Int) ([2]string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "getPaletteConfetti", id)

	if err != nil {
		return *new([2]string), err
	}

	out0 := *abi.ConvertType(out[0], new([2]string)).(*[2]string)

	return out0, err

}

// GetPaletteConfetti is a free data retrieval call binding the contract method 0xcbfb504a.
//
// Solidity: function getPaletteConfetti(uint256 id) view returns(string[2])
func (_ConfettiContract *ConfettiContractSession) GetPaletteConfetti(id *big.Int) ([2]string, error) {
	return _ConfettiContract.Contract.GetPaletteConfetti(&_ConfettiContract.CallOpts, id)
}

// GetPaletteConfetti is a free data retrieval call binding the contract method 0xcbfb504a.
//
// Solidity: function getPaletteConfetti(uint256 id) view returns(string[2])
func (_ConfettiContract *ConfettiContractCallerSession) GetPaletteConfetti(id *big.Int) ([2]string, error) {
	return _ConfettiContract.Contract.GetPaletteConfetti(&_ConfettiContract.CallOpts, id)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string[4],string[2]) confetti)
func (_ConfettiContract *ConfettiContractCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) (CONFETTIConfetti, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new(CONFETTIConfetti), err
	}

	out0 := *abi.ConvertType(out[0], new(CONFETTIConfetti)).(*CONFETTIConfetti)

	return out0, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string[4],string[2]) confetti)
func (_ConfettiContract *ConfettiContractSession) GetParamValues(tokenId *big.Int) (CONFETTIConfetti, error) {
	return _ConfettiContract.Contract.GetParamValues(&_ConfettiContract.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns((string,string,string[4],string[2]) confetti)
func (_ConfettiContract *ConfettiContractCallerSession) GetParamValues(tokenId *big.Int) (CONFETTIConfetti, error) {
	return _ConfettiContract.Contract.GetParamValues(&_ConfettiContract.CallOpts, tokenId)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ConfettiContract *ConfettiContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ConfettiContract *ConfettiContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ConfettiContract.Contract.IsApprovedForAll(&_ConfettiContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_ConfettiContract *ConfettiContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _ConfettiContract.Contract.IsApprovedForAll(&_ConfettiContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfettiContract *ConfettiContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfettiContract *ConfettiContractSession) Name() (string, error) {
	return _ConfettiContract.Contract.Name(&_ConfettiContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ConfettiContract *ConfettiContractCallerSession) Name() (string, error) {
	return _ConfettiContract.Contract.Name(&_ConfettiContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfettiContract *ConfettiContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfettiContract *ConfettiContractSession) Owner() (common.Address, error) {
	return _ConfettiContract.Contract.Owner(&_ConfettiContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) Owner() (common.Address, error) {
	return _ConfettiContract.Contract.Owner(&_ConfettiContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ConfettiContract.Contract.OwnerOf(&_ConfettiContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_ConfettiContract *ConfettiContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _ConfettiContract.Contract.OwnerOf(&_ConfettiContract.CallOpts, tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ConfettiContract *ConfettiContractCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ConfettiContract *ConfettiContractSession) Paused() (bool, error) {
	return _ConfettiContract.Contract.Paused(&_ConfettiContract.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_ConfettiContract *ConfettiContractCallerSession) Paused() (bool, error) {
	return _ConfettiContract.Contract.Paused(&_ConfettiContract.CallOpts)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_ConfettiContract *ConfettiContractCaller) RoyaltyInfo(opts *bind.CallOpts, _tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "royaltyInfo", _tokenId, _salePrice)

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
func (_ConfettiContract *ConfettiContractSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _ConfettiContract.Contract.RoyaltyInfo(&_ConfettiContract.CallOpts, _tokenId, _salePrice)
}

// RoyaltyInfo is a free data retrieval call binding the contract method 0x2a55205a.
//
// Solidity: function royaltyInfo(uint256 _tokenId, uint256 _salePrice) view returns(address receiver, uint256 royaltyAmount)
func (_ConfettiContract *ConfettiContractCallerSession) RoyaltyInfo(_tokenId *big.Int, _salePrice *big.Int) (struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}, error) {
	return _ConfettiContract.Contract.RoyaltyInfo(&_ConfettiContract.CallOpts, _tokenId, _salePrice)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ConfettiContract *ConfettiContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ConfettiContract *ConfettiContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ConfettiContract.Contract.SupportsInterface(&_ConfettiContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_ConfettiContract *ConfettiContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _ConfettiContract.Contract.SupportsInterface(&_ConfettiContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfettiContract *ConfettiContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfettiContract *ConfettiContractSession) Symbol() (string, error) {
	return _ConfettiContract.Contract.Symbol(&_ConfettiContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ConfettiContract *ConfettiContractCallerSession) Symbol() (string, error) {
	return _ConfettiContract.Contract.Symbol(&_ConfettiContract.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ConfettiContract *ConfettiContractCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _ConfettiContract.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ConfettiContract *ConfettiContractSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ConfettiContract.Contract.TokenURI(&_ConfettiContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_ConfettiContract *ConfettiContractCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _ConfettiContract.Contract.TokenURI(&_ConfettiContract.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.Approve(&_ConfettiContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.Approve(&_ConfettiContract.TransactOpts, to, tokenId)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_ConfettiContract *ConfettiContractTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_ConfettiContract *ConfettiContractSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeAdmin(&_ConfettiContract.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeAdmin(&_ConfettiContract.TransactOpts, newAdm)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_ConfettiContract *ConfettiContractTransactor) ChangeBaseURI(opts *bind.TransactOpts, baseURI string) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "changeBaseURI", baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_ConfettiContract *ConfettiContractSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeBaseURI(&_ConfettiContract.TransactOpts, baseURI)
}

// ChangeBaseURI is a paid mutator transaction binding the contract method 0x39a0c6f9.
//
// Solidity: function changeBaseURI(string baseURI) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) ChangeBaseURI(baseURI string) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeBaseURI(&_ConfettiContract.TransactOpts, baseURI)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_ConfettiContract *ConfettiContractTransactor) ChangeParam(opts *bind.TransactOpts, newP common.Address) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "changeParam", newP)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_ConfettiContract *ConfettiContractSession) ChangeParam(newP common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeParam(&_ConfettiContract.TransactOpts, newP)
}

// ChangeParam is a paid mutator transaction binding the contract method 0x741149b1.
//
// Solidity: function changeParam(address newP) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) ChangeParam(newP common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeParam(&_ConfettiContract.TransactOpts, newP)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_ConfettiContract *ConfettiContractTransactor) ChangeToken(opts *bind.TransactOpts, sweet common.Address) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "changeToken", sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_ConfettiContract *ConfettiContractSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeToken(&_ConfettiContract.TransactOpts, sweet)
}

// ChangeToken is a paid mutator transaction binding the contract method 0x66829b16.
//
// Solidity: function changeToken(address sweet) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) ChangeToken(sweet common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.ChangeToken(&_ConfettiContract.TransactOpts, sweet)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_ConfettiContract *ConfettiContractTransactor) Initialize(opts *bind.TransactOpts, name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "initialize", name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_ConfettiContract *ConfettiContractSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.Initialize(&_ConfettiContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0x8f15b414.
//
// Solidity: function initialize(string name, string symbol, address admin, address paramsAddress) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Initialize(name string, symbol string, admin common.Address, paramsAddress common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.Initialize(&_ConfettiContract.TransactOpts, name, symbol, admin, paramsAddress)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_ConfettiContract *ConfettiContractTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_ConfettiContract *ConfettiContractSession) Mint() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Mint(&_ConfettiContract.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() payable returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Mint() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Mint(&_ConfettiContract.TransactOpts)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_ConfettiContract *ConfettiContractTransactor) MintByToken(opts *bind.TransactOpts, tokenIdGated *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "mintByToken", tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_ConfettiContract *ConfettiContractSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.MintByToken(&_ConfettiContract.TransactOpts, tokenIdGated)
}

// MintByToken is a paid mutator transaction binding the contract method 0x80d953dd.
//
// Solidity: function mintByToken(uint256 tokenIdGated) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) MintByToken(tokenIdGated *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.MintByToken(&_ConfettiContract.TransactOpts, tokenIdGated)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_ConfettiContract *ConfettiContractTransactor) OwnerMint(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "ownerMint", id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_ConfettiContract *ConfettiContractSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.OwnerMint(&_ConfettiContract.TransactOpts, id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.OwnerMint(&_ConfettiContract.TransactOpts, id)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ConfettiContract *ConfettiContractTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ConfettiContract *ConfettiContractSession) Pause() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Pause(&_ConfettiContract.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Pause() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Pause(&_ConfettiContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConfettiContract *ConfettiContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConfettiContract *ConfettiContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConfettiContract.Contract.RenounceOwnership(&_ConfettiContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ConfettiContract *ConfettiContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ConfettiContract.Contract.RenounceOwnership(&_ConfettiContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SafeTransferFrom(&_ConfettiContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SafeTransferFrom(&_ConfettiContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ConfettiContract *ConfettiContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ConfettiContract *ConfettiContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SafeTransferFrom0(&_ConfettiContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SafeTransferFrom0(&_ConfettiContract.TransactOpts, from, to, tokenId, data)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_ConfettiContract *ConfettiContractTransactor) SetAlgo(opts *bind.TransactOpts, algo string) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "setAlgo", algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_ConfettiContract *ConfettiContractSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetAlgo(&_ConfettiContract.TransactOpts, algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetAlgo(&_ConfettiContract.TransactOpts, algo)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ConfettiContract *ConfettiContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ConfettiContract *ConfettiContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetApprovalForAll(&_ConfettiContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetApprovalForAll(&_ConfettiContract.TransactOpts, operator, approved)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_ConfettiContract *ConfettiContractTransactor) SetFee(opts *bind.TransactOpts, fee *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "setFee", fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_ConfettiContract *ConfettiContractSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetFee(&_ConfettiContract.TransactOpts, fee)
}

// SetFee is a paid mutator transaction binding the contract method 0x69fe0e2d.
//
// Solidity: function setFee(uint256 fee) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SetFee(fee *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetFee(&_ConfettiContract.TransactOpts, fee)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_ConfettiContract *ConfettiContractTransactor) SetLimit(opts *bind.TransactOpts, limit *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "setLimit", limit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_ConfettiContract *ConfettiContractSession) SetLimit(limit *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetLimit(&_ConfettiContract.TransactOpts, limit)
}

// SetLimit is a paid mutator transaction binding the contract method 0x27ea6f2b.
//
// Solidity: function setLimit(uint256 limit) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) SetLimit(limit *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.SetLimit(&_ConfettiContract.TransactOpts, limit)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.TransferFrom(&_ConfettiContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _ConfettiContract.Contract.TransferFrom(&_ConfettiContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConfettiContract *ConfettiContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConfettiContract *ConfettiContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.TransferOwnership(&_ConfettiContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ConfettiContract *ConfettiContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ConfettiContract.Contract.TransferOwnership(&_ConfettiContract.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ConfettiContract *ConfettiContractTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ConfettiContract *ConfettiContractSession) Unpause() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Unpause(&_ConfettiContract.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Unpause() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Unpause(&_ConfettiContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ConfettiContract *ConfettiContractTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConfettiContract.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ConfettiContract *ConfettiContractSession) Withdraw() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Withdraw(&_ConfettiContract.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns()
func (_ConfettiContract *ConfettiContractTransactorSession) Withdraw() (*types.Transaction, error) {
	return _ConfettiContract.Contract.Withdraw(&_ConfettiContract.TransactOpts)
}

// ConfettiContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ConfettiContract contract.
type ConfettiContractApprovalIterator struct {
	Event *ConfettiContractApproval // Event containing the contract specifics and raw log

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
func (it *ConfettiContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractApproval)
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
		it.Event = new(ConfettiContractApproval)
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
func (it *ConfettiContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractApproval represents a Approval event raised by the ConfettiContract contract.
type ConfettiContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ConfettiContract *ConfettiContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*ConfettiContractApprovalIterator, error) {

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

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractApprovalIterator{contract: _ConfettiContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_ConfettiContract *ConfettiContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ConfettiContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractApproval)
				if err := _ConfettiContract.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseApproval(log types.Log) (*ConfettiContractApproval, error) {
	event := new(ConfettiContractApproval)
	if err := _ConfettiContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the ConfettiContract contract.
type ConfettiContractApprovalForAllIterator struct {
	Event *ConfettiContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *ConfettiContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractApprovalForAll)
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
		it.Event = new(ConfettiContractApprovalForAll)
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
func (it *ConfettiContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractApprovalForAll represents a ApprovalForAll event raised by the ConfettiContract contract.
type ConfettiContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ConfettiContract *ConfettiContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*ConfettiContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractApprovalForAllIterator{contract: _ConfettiContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_ConfettiContract *ConfettiContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *ConfettiContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractApprovalForAll)
				if err := _ConfettiContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseApprovalForAll(log types.Log) (*ConfettiContractApprovalForAll, error) {
	event := new(ConfettiContractApprovalForAll)
	if err := _ConfettiContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ConfettiContract contract.
type ConfettiContractInitializedIterator struct {
	Event *ConfettiContractInitialized // Event containing the contract specifics and raw log

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
func (it *ConfettiContractInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractInitialized)
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
		it.Event = new(ConfettiContractInitialized)
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
func (it *ConfettiContractInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractInitialized represents a Initialized event raised by the ConfettiContract contract.
type ConfettiContractInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ConfettiContract *ConfettiContractFilterer) FilterInitialized(opts *bind.FilterOpts) (*ConfettiContractInitializedIterator, error) {

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ConfettiContractInitializedIterator{contract: _ConfettiContract.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_ConfettiContract *ConfettiContractFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ConfettiContractInitialized) (event.Subscription, error) {

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractInitialized)
				if err := _ConfettiContract.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseInitialized(log types.Log) (*ConfettiContractInitialized, error) {
	event := new(ConfettiContractInitialized)
	if err := _ConfettiContract.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ConfettiContract contract.
type ConfettiContractOwnershipTransferredIterator struct {
	Event *ConfettiContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ConfettiContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractOwnershipTransferred)
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
		it.Event = new(ConfettiContractOwnershipTransferred)
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
func (it *ConfettiContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractOwnershipTransferred represents a OwnershipTransferred event raised by the ConfettiContract contract.
type ConfettiContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConfettiContract *ConfettiContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ConfettiContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractOwnershipTransferredIterator{contract: _ConfettiContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ConfettiContract *ConfettiContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ConfettiContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractOwnershipTransferred)
				if err := _ConfettiContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseOwnershipTransferred(log types.Log) (*ConfettiContractOwnershipTransferred, error) {
	event := new(ConfettiContractOwnershipTransferred)
	if err := _ConfettiContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the ConfettiContract contract.
type ConfettiContractPausedIterator struct {
	Event *ConfettiContractPaused // Event containing the contract specifics and raw log

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
func (it *ConfettiContractPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractPaused)
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
		it.Event = new(ConfettiContractPaused)
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
func (it *ConfettiContractPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractPaused represents a Paused event raised by the ConfettiContract contract.
type ConfettiContractPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ConfettiContract *ConfettiContractFilterer) FilterPaused(opts *bind.FilterOpts) (*ConfettiContractPausedIterator, error) {

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &ConfettiContractPausedIterator{contract: _ConfettiContract.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_ConfettiContract *ConfettiContractFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *ConfettiContractPaused) (event.Subscription, error) {

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractPaused)
				if err := _ConfettiContract.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParsePaused(log types.Log) (*ConfettiContractPaused, error) {
	event := new(ConfettiContractPaused)
	if err := _ConfettiContract.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ConfettiContract contract.
type ConfettiContractTransferIterator struct {
	Event *ConfettiContractTransfer // Event containing the contract specifics and raw log

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
func (it *ConfettiContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractTransfer)
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
		it.Event = new(ConfettiContractTransfer)
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
func (it *ConfettiContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractTransfer represents a Transfer event raised by the ConfettiContract contract.
type ConfettiContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ConfettiContract *ConfettiContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*ConfettiContractTransferIterator, error) {

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

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &ConfettiContractTransferIterator{contract: _ConfettiContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_ConfettiContract *ConfettiContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ConfettiContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractTransfer)
				if err := _ConfettiContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseTransfer(log types.Log) (*ConfettiContractTransfer, error) {
	event := new(ConfettiContractTransfer)
	if err := _ConfettiContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ConfettiContractUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the ConfettiContract contract.
type ConfettiContractUnpausedIterator struct {
	Event *ConfettiContractUnpaused // Event containing the contract specifics and raw log

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
func (it *ConfettiContractUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ConfettiContractUnpaused)
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
		it.Event = new(ConfettiContractUnpaused)
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
func (it *ConfettiContractUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ConfettiContractUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ConfettiContractUnpaused represents a Unpaused event raised by the ConfettiContract contract.
type ConfettiContractUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ConfettiContract *ConfettiContractFilterer) FilterUnpaused(opts *bind.FilterOpts) (*ConfettiContractUnpausedIterator, error) {

	logs, sub, err := _ConfettiContract.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &ConfettiContractUnpausedIterator{contract: _ConfettiContract.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_ConfettiContract *ConfettiContractFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *ConfettiContractUnpaused) (event.Subscription, error) {

	logs, sub, err := _ConfettiContract.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ConfettiContractUnpaused)
				if err := _ConfettiContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_ConfettiContract *ConfettiContractFilterer) ParseUnpaused(log types.Log) (*ConfettiContractUnpaused, error) {
	event := new(ConfettiContractUnpaused)
	if err := _ConfettiContract.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
