// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package candy_contract

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

// CandyContractMetaData contains all meta data concerning the CandyContract contract.
var CandyContractMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"algorithm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getPalette\",\"outputs\":[{\"internalType\":\"string[4]\",\"name\":\"\",\"type\":\"string[4]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getParamValues\",\"outputs\":[{\"internalType\":\"string[4]\",\"name\":\"\",\"type\":\"string[4]\"},{\"internalType\":\"string\",\"name\":\"shape\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"size\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"surface\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getShape\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getSize\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"getSurface\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"ownerMint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"algo\",\"type\":\"string\"}],\"name\":\"setAlgo\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// CandyContractABI is the input ABI used to generate the binding from.
// Deprecated: Use CandyContractMetaData.ABI instead.
var CandyContractABI = CandyContractMetaData.ABI

// CandyContract is an auto generated Go binding around an Ethereum contract.
type CandyContract struct {
	CandyContractCaller     // Read-only binding to the contract
	CandyContractTransactor // Write-only binding to the contract
	CandyContractFilterer   // Log filterer for contract events
}

// CandyContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type CandyContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CandyContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CandyContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CandyContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CandyContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CandyContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CandyContractSession struct {
	Contract     *CandyContract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CandyContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CandyContractCallerSession struct {
	Contract *CandyContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// CandyContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CandyContractTransactorSession struct {
	Contract     *CandyContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CandyContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type CandyContractRaw struct {
	Contract *CandyContract // Generic contract binding to access the raw methods on
}

// CandyContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CandyContractCallerRaw struct {
	Contract *CandyContractCaller // Generic read-only contract binding to access the raw methods on
}

// CandyContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CandyContractTransactorRaw struct {
	Contract *CandyContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCandyContract creates a new instance of CandyContract, bound to a specific deployed contract.
func NewCandyContract(address common.Address, backend bind.ContractBackend) (*CandyContract, error) {
	contract, err := bindCandyContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CandyContract{CandyContractCaller: CandyContractCaller{contract: contract}, CandyContractTransactor: CandyContractTransactor{contract: contract}, CandyContractFilterer: CandyContractFilterer{contract: contract}}, nil
}

// NewCandyContractCaller creates a new read-only instance of CandyContract, bound to a specific deployed contract.
func NewCandyContractCaller(address common.Address, caller bind.ContractCaller) (*CandyContractCaller, error) {
	contract, err := bindCandyContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CandyContractCaller{contract: contract}, nil
}

// NewCandyContractTransactor creates a new write-only instance of CandyContract, bound to a specific deployed contract.
func NewCandyContractTransactor(address common.Address, transactor bind.ContractTransactor) (*CandyContractTransactor, error) {
	contract, err := bindCandyContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CandyContractTransactor{contract: contract}, nil
}

// NewCandyContractFilterer creates a new log filterer instance of CandyContract, bound to a specific deployed contract.
func NewCandyContractFilterer(address common.Address, filterer bind.ContractFilterer) (*CandyContractFilterer, error) {
	contract, err := bindCandyContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CandyContractFilterer{contract: contract}, nil
}

// bindCandyContract binds a generic wrapper to an already deployed contract.
func bindCandyContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CandyContractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CandyContract *CandyContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CandyContract.Contract.CandyContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CandyContract *CandyContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CandyContract.Contract.CandyContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CandyContract *CandyContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CandyContract.Contract.CandyContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CandyContract *CandyContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CandyContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CandyContract *CandyContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CandyContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CandyContract *CandyContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CandyContract.Contract.contract.Transact(opts, method, params...)
}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_CandyContract *CandyContractCaller) Algorithm(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "algorithm")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_CandyContract *CandyContractSession) Algorithm() (string, error) {
	return _CandyContract.Contract.Algorithm(&_CandyContract.CallOpts)
}

// Algorithm is a free data retrieval call binding the contract method 0x2e58ca46.
//
// Solidity: function algorithm() view returns(string)
func (_CandyContract *CandyContractCallerSession) Algorithm() (string, error) {
	return _CandyContract.Contract.Algorithm(&_CandyContract.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CandyContract *CandyContractCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CandyContract *CandyContractSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CandyContract.Contract.BalanceOf(&_CandyContract.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_CandyContract *CandyContractCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _CandyContract.Contract.BalanceOf(&_CandyContract.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CandyContract.Contract.GetApproved(&_CandyContract.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _CandyContract.Contract.GetApproved(&_CandyContract.CallOpts, tokenId)
}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_CandyContract *CandyContractCaller) GetPalette(opts *bind.CallOpts, id *big.Int) ([4]string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getPalette", id)

	if err != nil {
		return *new([4]string), err
	}

	out0 := *abi.ConvertType(out[0], new([4]string)).(*[4]string)

	return out0, err

}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_CandyContract *CandyContractSession) GetPalette(id *big.Int) ([4]string, error) {
	return _CandyContract.Contract.GetPalette(&_CandyContract.CallOpts, id)
}

// GetPalette is a free data retrieval call binding the contract method 0x505e570a.
//
// Solidity: function getPalette(uint256 id) view returns(string[4])
func (_CandyContract *CandyContractCallerSession) GetPalette(id *big.Int) ([4]string, error) {
	return _CandyContract.Contract.GetPalette(&_CandyContract.CallOpts, id)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_CandyContract *CandyContractCaller) GetParamValues(opts *bind.CallOpts, tokenId *big.Int) ([4]string, string, string, string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getParamValues", tokenId)

	if err != nil {
		return *new([4]string), *new(string), *new(string), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new([4]string)).(*[4]string)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)
	out2 := *abi.ConvertType(out[2], new(string)).(*string)
	out3 := *abi.ConvertType(out[3], new(string)).(*string)

	return out0, out1, out2, out3, err

}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_CandyContract *CandyContractSession) GetParamValues(tokenId *big.Int) ([4]string, string, string, string, error) {
	return _CandyContract.Contract.GetParamValues(&_CandyContract.CallOpts, tokenId)
}

// GetParamValues is a free data retrieval call binding the contract method 0xc0dca239.
//
// Solidity: function getParamValues(uint256 tokenId) view returns(string[4], string shape, string size, string surface)
func (_CandyContract *CandyContractCallerSession) GetParamValues(tokenId *big.Int) ([4]string, string, string, string, error) {
	return _CandyContract.Contract.GetParamValues(&_CandyContract.CallOpts, tokenId)
}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_CandyContract *CandyContractCaller) GetShape(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getShape", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_CandyContract *CandyContractSession) GetShape(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetShape(&_CandyContract.CallOpts, id)
}

// GetShape is a free data retrieval call binding the contract method 0x22ced721.
//
// Solidity: function getShape(uint256 id) view returns(string)
func (_CandyContract *CandyContractCallerSession) GetShape(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetShape(&_CandyContract.CallOpts, id)
}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_CandyContract *CandyContractCaller) GetSize(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getSize", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_CandyContract *CandyContractSession) GetSize(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetSize(&_CandyContract.CallOpts, id)
}

// GetSize is a free data retrieval call binding the contract method 0x023c23db.
//
// Solidity: function getSize(uint256 id) view returns(string)
func (_CandyContract *CandyContractCallerSession) GetSize(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetSize(&_CandyContract.CallOpts, id)
}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_CandyContract *CandyContractCaller) GetSurface(opts *bind.CallOpts, id *big.Int) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "getSurface", id)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_CandyContract *CandyContractSession) GetSurface(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetSurface(&_CandyContract.CallOpts, id)
}

// GetSurface is a free data retrieval call binding the contract method 0x11ad8047.
//
// Solidity: function getSurface(uint256 id) view returns(string)
func (_CandyContract *CandyContractCallerSession) GetSurface(id *big.Int) (string, error) {
	return _CandyContract.Contract.GetSurface(&_CandyContract.CallOpts, id)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CandyContract *CandyContractCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CandyContract *CandyContractSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CandyContract.Contract.IsApprovedForAll(&_CandyContract.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_CandyContract *CandyContractCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _CandyContract.Contract.IsApprovedForAll(&_CandyContract.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CandyContract *CandyContractCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CandyContract *CandyContractSession) Name() (string, error) {
	return _CandyContract.Contract.Name(&_CandyContract.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_CandyContract *CandyContractCallerSession) Name() (string, error) {
	return _CandyContract.Contract.Name(&_CandyContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CandyContract *CandyContractCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CandyContract *CandyContractSession) Owner() (common.Address, error) {
	return _CandyContract.Contract.Owner(&_CandyContract.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_CandyContract *CandyContractCallerSession) Owner() (common.Address, error) {
	return _CandyContract.Contract.Owner(&_CandyContract.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CandyContract.Contract.OwnerOf(&_CandyContract.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_CandyContract *CandyContractCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _CandyContract.Contract.OwnerOf(&_CandyContract.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CandyContract *CandyContractCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CandyContract *CandyContractSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CandyContract.Contract.SupportsInterface(&_CandyContract.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_CandyContract *CandyContractCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _CandyContract.Contract.SupportsInterface(&_CandyContract.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CandyContract *CandyContractCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CandyContract *CandyContractSession) Symbol() (string, error) {
	return _CandyContract.Contract.Symbol(&_CandyContract.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_CandyContract *CandyContractCallerSession) Symbol() (string, error) {
	return _CandyContract.Contract.Symbol(&_CandyContract.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CandyContract *CandyContractCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _CandyContract.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CandyContract *CandyContractSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CandyContract.Contract.TokenURI(&_CandyContract.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_CandyContract *CandyContractCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _CandyContract.Contract.TokenURI(&_CandyContract.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.Approve(&_CandyContract.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.Approve(&_CandyContract.TransactOpts, to, tokenId)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_CandyContract *CandyContractTransactor) Mint(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "mint")
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_CandyContract *CandyContractSession) Mint() (*types.Transaction, error) {
	return _CandyContract.Contract.Mint(&_CandyContract.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0x1249c58b.
//
// Solidity: function mint() returns()
func (_CandyContract *CandyContractTransactorSession) Mint() (*types.Transaction, error) {
	return _CandyContract.Contract.Mint(&_CandyContract.TransactOpts)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_CandyContract *CandyContractTransactor) OwnerMint(opts *bind.TransactOpts, id *big.Int) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "ownerMint", id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_CandyContract *CandyContractSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.OwnerMint(&_CandyContract.TransactOpts, id)
}

// OwnerMint is a paid mutator transaction binding the contract method 0xf19e75d4.
//
// Solidity: function ownerMint(uint256 id) returns()
func (_CandyContract *CandyContractTransactorSession) OwnerMint(id *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.OwnerMint(&_CandyContract.TransactOpts, id)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CandyContract *CandyContractTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CandyContract *CandyContractSession) RenounceOwnership() (*types.Transaction, error) {
	return _CandyContract.Contract.RenounceOwnership(&_CandyContract.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_CandyContract *CandyContractTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _CandyContract.Contract.RenounceOwnership(&_CandyContract.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.SafeTransferFrom(&_CandyContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.SafeTransferFrom(&_CandyContract.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CandyContract *CandyContractTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CandyContract *CandyContractSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CandyContract.Contract.SafeTransferFrom0(&_CandyContract.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_CandyContract *CandyContractTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _CandyContract.Contract.SafeTransferFrom0(&_CandyContract.TransactOpts, from, to, tokenId, data)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_CandyContract *CandyContractTransactor) SetAlgo(opts *bind.TransactOpts, algo string) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "setAlgo", algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_CandyContract *CandyContractSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _CandyContract.Contract.SetAlgo(&_CandyContract.TransactOpts, algo)
}

// SetAlgo is a paid mutator transaction binding the contract method 0x0764da1d.
//
// Solidity: function setAlgo(string algo) returns()
func (_CandyContract *CandyContractTransactorSession) SetAlgo(algo string) (*types.Transaction, error) {
	return _CandyContract.Contract.SetAlgo(&_CandyContract.TransactOpts, algo)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CandyContract *CandyContractTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CandyContract *CandyContractSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CandyContract.Contract.SetApprovalForAll(&_CandyContract.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_CandyContract *CandyContractTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _CandyContract.Contract.SetApprovalForAll(&_CandyContract.TransactOpts, operator, approved)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.TransferFrom(&_CandyContract.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_CandyContract *CandyContractTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _CandyContract.Contract.TransferFrom(&_CandyContract.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CandyContract *CandyContractTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _CandyContract.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CandyContract *CandyContractSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CandyContract.Contract.TransferOwnership(&_CandyContract.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_CandyContract *CandyContractTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _CandyContract.Contract.TransferOwnership(&_CandyContract.TransactOpts, newOwner)
}

// CandyContractApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the CandyContract contract.
type CandyContractApprovalIterator struct {
	Event *CandyContractApproval // Event containing the contract specifics and raw log

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
func (it *CandyContractApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CandyContractApproval)
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
		it.Event = new(CandyContractApproval)
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
func (it *CandyContractApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CandyContractApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CandyContractApproval represents a Approval event raised by the CandyContract contract.
type CandyContractApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CandyContract *CandyContractFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*CandyContractApprovalIterator, error) {

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

	logs, sub, err := _CandyContract.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CandyContractApprovalIterator{contract: _CandyContract.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_CandyContract *CandyContractFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *CandyContractApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _CandyContract.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CandyContractApproval)
				if err := _CandyContract.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_CandyContract *CandyContractFilterer) ParseApproval(log types.Log) (*CandyContractApproval, error) {
	event := new(CandyContractApproval)
	if err := _CandyContract.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CandyContractApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the CandyContract contract.
type CandyContractApprovalForAllIterator struct {
	Event *CandyContractApprovalForAll // Event containing the contract specifics and raw log

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
func (it *CandyContractApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CandyContractApprovalForAll)
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
		it.Event = new(CandyContractApprovalForAll)
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
func (it *CandyContractApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CandyContractApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CandyContractApprovalForAll represents a ApprovalForAll event raised by the CandyContract contract.
type CandyContractApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CandyContract *CandyContractFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*CandyContractApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CandyContract.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &CandyContractApprovalForAllIterator{contract: _CandyContract.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_CandyContract *CandyContractFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *CandyContractApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _CandyContract.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CandyContractApprovalForAll)
				if err := _CandyContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
func (_CandyContract *CandyContractFilterer) ParseApprovalForAll(log types.Log) (*CandyContractApprovalForAll, error) {
	event := new(CandyContractApprovalForAll)
	if err := _CandyContract.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CandyContractOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the CandyContract contract.
type CandyContractOwnershipTransferredIterator struct {
	Event *CandyContractOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *CandyContractOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CandyContractOwnershipTransferred)
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
		it.Event = new(CandyContractOwnershipTransferred)
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
func (it *CandyContractOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CandyContractOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CandyContractOwnershipTransferred represents a OwnershipTransferred event raised by the CandyContract contract.
type CandyContractOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CandyContract *CandyContractFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*CandyContractOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CandyContract.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &CandyContractOwnershipTransferredIterator{contract: _CandyContract.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_CandyContract *CandyContractFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *CandyContractOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _CandyContract.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CandyContractOwnershipTransferred)
				if err := _CandyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_CandyContract *CandyContractFilterer) ParseOwnershipTransferred(log types.Log) (*CandyContractOwnershipTransferred, error) {
	event := new(CandyContractOwnershipTransferred)
	if err := _CandyContract.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CandyContractTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the CandyContract contract.
type CandyContractTransferIterator struct {
	Event *CandyContractTransfer // Event containing the contract specifics and raw log

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
func (it *CandyContractTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CandyContractTransfer)
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
		it.Event = new(CandyContractTransfer)
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
func (it *CandyContractTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CandyContractTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CandyContractTransfer represents a Transfer event raised by the CandyContract contract.
type CandyContractTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CandyContract *CandyContractFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*CandyContractTransferIterator, error) {

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

	logs, sub, err := _CandyContract.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &CandyContractTransferIterator{contract: _CandyContract.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_CandyContract *CandyContractFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *CandyContractTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _CandyContract.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CandyContractTransfer)
				if err := _CandyContract.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_CandyContract *CandyContractFilterer) ParseTransfer(log types.Log) (*CandyContractTransfer, error) {
	event := new(CandyContractTransfer)
	if err := _CandyContract.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
