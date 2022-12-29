// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generative_project_data

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

// GenerativeProjectDataMetaData contains all meta data concerning the GenerativeProjectData contract.
var GenerativeProjectDataMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"_admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_generativeProjectAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"_paramAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdm\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeParamAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAddr\",\"type\":\"address\"}],\"name\":\"changeProjectAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"script\",\"type\":\"string\"}],\"name\":\"inflateScript\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"},{\"internalType\":\"enumInflate.ErrorCode\",\"name\":\"err\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"paramAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"generativeProjectAddr\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"}],\"name\":\"projectURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"}],\"name\":\"tokenBaseURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"}],\"name\":\"tokenHTML\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"projectId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"seed\",\"type\":\"bytes32\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"result\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// GenerativeProjectDataABI is the input ABI used to generate the binding from.
// Deprecated: Use GenerativeProjectDataMetaData.ABI instead.
var GenerativeProjectDataABI = GenerativeProjectDataMetaData.ABI

// GenerativeProjectData is an auto generated Go binding around an Ethereum contract.
type GenerativeProjectData struct {
	GenerativeProjectDataCaller     // Read-only binding to the contract
	GenerativeProjectDataTransactor // Write-only binding to the contract
	GenerativeProjectDataFilterer   // Log filterer for contract events
}

// GenerativeProjectDataCaller is an auto generated read-only Go binding around an Ethereum contract.
type GenerativeProjectDataCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectDataTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GenerativeProjectDataTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectDataFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GenerativeProjectDataFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GenerativeProjectDataSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GenerativeProjectDataSession struct {
	Contract     *GenerativeProjectData // Generic contract binding to set the session for
	CallOpts     bind.CallOpts          // Call options to use throughout this session
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// GenerativeProjectDataCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GenerativeProjectDataCallerSession struct {
	Contract *GenerativeProjectDataCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                // Call options to use throughout this session
}

// GenerativeProjectDataTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GenerativeProjectDataTransactorSession struct {
	Contract     *GenerativeProjectDataTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// GenerativeProjectDataRaw is an auto generated low-level Go binding around an Ethereum contract.
type GenerativeProjectDataRaw struct {
	Contract *GenerativeProjectData // Generic contract binding to access the raw methods on
}

// GenerativeProjectDataCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GenerativeProjectDataCallerRaw struct {
	Contract *GenerativeProjectDataCaller // Generic read-only contract binding to access the raw methods on
}

// GenerativeProjectDataTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GenerativeProjectDataTransactorRaw struct {
	Contract *GenerativeProjectDataTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGenerativeProjectData creates a new instance of GenerativeProjectData, bound to a specific deployed contract.
func NewGenerativeProjectData(address common.Address, backend bind.ContractBackend) (*GenerativeProjectData, error) {
	contract, err := bindGenerativeProjectData(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectData{GenerativeProjectDataCaller: GenerativeProjectDataCaller{contract: contract}, GenerativeProjectDataTransactor: GenerativeProjectDataTransactor{contract: contract}, GenerativeProjectDataFilterer: GenerativeProjectDataFilterer{contract: contract}}, nil
}

// NewGenerativeProjectDataCaller creates a new read-only instance of GenerativeProjectData, bound to a specific deployed contract.
func NewGenerativeProjectDataCaller(address common.Address, caller bind.ContractCaller) (*GenerativeProjectDataCaller, error) {
	contract, err := bindGenerativeProjectData(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectDataCaller{contract: contract}, nil
}

// NewGenerativeProjectDataTransactor creates a new write-only instance of GenerativeProjectData, bound to a specific deployed contract.
func NewGenerativeProjectDataTransactor(address common.Address, transactor bind.ContractTransactor) (*GenerativeProjectDataTransactor, error) {
	contract, err := bindGenerativeProjectData(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectDataTransactor{contract: contract}, nil
}

// NewGenerativeProjectDataFilterer creates a new log filterer instance of GenerativeProjectData, bound to a specific deployed contract.
func NewGenerativeProjectDataFilterer(address common.Address, filterer bind.ContractFilterer) (*GenerativeProjectDataFilterer, error) {
	contract, err := bindGenerativeProjectData(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectDataFilterer{contract: contract}, nil
}

// bindGenerativeProjectData binds a generic wrapper to an already deployed contract.
func bindGenerativeProjectData(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GenerativeProjectDataMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeProjectData *GenerativeProjectDataRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeProjectData.Contract.GenerativeProjectDataCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeProjectData *GenerativeProjectDataRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.GenerativeProjectDataTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeProjectData *GenerativeProjectDataRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.GenerativeProjectDataTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GenerativeProjectData *GenerativeProjectDataCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GenerativeProjectData.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GenerativeProjectData *GenerativeProjectDataTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GenerativeProjectData *GenerativeProjectDataTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataSession) Admin() (common.Address, error) {
	return _GenerativeProjectData.Contract.Admin(&_GenerativeProjectData.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0x01bc45c9.
//
// Solidity: function _admin() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) Admin() (common.Address, error) {
	return _GenerativeProjectData.Contract.Admin(&_GenerativeProjectData.CallOpts)
}

// GenerativeProjectAddr is a free data retrieval call binding the contract method 0x567fdcf3.
//
// Solidity: function _generativeProjectAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCaller) GenerativeProjectAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "_generativeProjectAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GenerativeProjectAddr is a free data retrieval call binding the contract method 0x567fdcf3.
//
// Solidity: function _generativeProjectAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataSession) GenerativeProjectAddr() (common.Address, error) {
	return _GenerativeProjectData.Contract.GenerativeProjectAddr(&_GenerativeProjectData.CallOpts)
}

// GenerativeProjectAddr is a free data retrieval call binding the contract method 0x567fdcf3.
//
// Solidity: function _generativeProjectAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) GenerativeProjectAddr() (common.Address, error) {
	return _GenerativeProjectData.Contract.GenerativeProjectAddr(&_GenerativeProjectData.CallOpts)
}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCaller) ParamAddr(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "_paramAddr")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataSession) ParamAddr() (common.Address, error) {
	return _GenerativeProjectData.Contract.ParamAddr(&_GenerativeProjectData.CallOpts)
}

// ParamAddr is a free data retrieval call binding the contract method 0xf4a290f7.
//
// Solidity: function _paramAddr() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) ParamAddr() (common.Address, error) {
	return _GenerativeProjectData.Contract.ParamAddr(&_GenerativeProjectData.CallOpts)
}

// InflateScript is a free data retrieval call binding the contract method 0x7af34af2.
//
// Solidity: function inflateScript(string script) view returns(string result, uint8 err)
func (_GenerativeProjectData *GenerativeProjectDataCaller) InflateScript(opts *bind.CallOpts, script string) (struct {
	Result string
	Err    uint8
}, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "inflateScript", script)

	outstruct := new(struct {
		Result string
		Err    uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Result = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Err = *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return *outstruct, err

}

// InflateScript is a free data retrieval call binding the contract method 0x7af34af2.
//
// Solidity: function inflateScript(string script) view returns(string result, uint8 err)
func (_GenerativeProjectData *GenerativeProjectDataSession) InflateScript(script string) (struct {
	Result string
	Err    uint8
}, error) {
	return _GenerativeProjectData.Contract.InflateScript(&_GenerativeProjectData.CallOpts, script)
}

// InflateScript is a free data retrieval call binding the contract method 0x7af34af2.
//
// Solidity: function inflateScript(string script) view returns(string result, uint8 err)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) InflateScript(script string) (struct {
	Result string
	Err    uint8
}, error) {
	return _GenerativeProjectData.Contract.InflateScript(&_GenerativeProjectData.CallOpts, script)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataSession) Owner() (common.Address, error) {
	return _GenerativeProjectData.Contract.Owner(&_GenerativeProjectData.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) Owner() (common.Address, error) {
	return _GenerativeProjectData.Contract.Owner(&_GenerativeProjectData.CallOpts)
}

// ProjectURI is a free data retrieval call binding the contract method 0x79ebc6f0.
//
// Solidity: function projectURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCaller) ProjectURI(opts *bind.CallOpts, projectId *big.Int) (string, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "projectURI", projectId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ProjectURI is a free data retrieval call binding the contract method 0x79ebc6f0.
//
// Solidity: function projectURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataSession) ProjectURI(projectId *big.Int) (string, error) {
	return _GenerativeProjectData.Contract.ProjectURI(&_GenerativeProjectData.CallOpts, projectId)
}

// ProjectURI is a free data retrieval call binding the contract method 0x79ebc6f0.
//
// Solidity: function projectURI(uint256 projectId) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) ProjectURI(projectId *big.Int) (string, error) {
	return _GenerativeProjectData.Contract.ProjectURI(&_GenerativeProjectData.CallOpts, projectId)
}

// TokenBaseURI is a free data retrieval call binding the contract method 0x46eaa748.
//
// Solidity: function tokenBaseURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCaller) TokenBaseURI(opts *bind.CallOpts, projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "tokenBaseURI", projectId, tokenId, seed)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenBaseURI is a free data retrieval call binding the contract method 0x46eaa748.
//
// Solidity: function tokenBaseURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataSession) TokenBaseURI(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenBaseURI(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// TokenBaseURI is a free data retrieval call binding the contract method 0x46eaa748.
//
// Solidity: function tokenBaseURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) TokenBaseURI(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenBaseURI(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// TokenHTML is a free data retrieval call binding the contract method 0x31c7280c.
//
// Solidity: function tokenHTML(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCaller) TokenHTML(opts *bind.CallOpts, projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "tokenHTML", projectId, tokenId, seed)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenHTML is a free data retrieval call binding the contract method 0x31c7280c.
//
// Solidity: function tokenHTML(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataSession) TokenHTML(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenHTML(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// TokenHTML is a free data retrieval call binding the contract method 0x31c7280c.
//
// Solidity: function tokenHTML(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) TokenHTML(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenHTML(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// TokenURI is a free data retrieval call binding the contract method 0x94139301.
//
// Solidity: function tokenURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCaller) TokenURI(opts *bind.CallOpts, projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	var out []interface{}
	err := _GenerativeProjectData.contract.Call(opts, &out, "tokenURI", projectId, tokenId, seed)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0x94139301.
//
// Solidity: function tokenURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataSession) TokenURI(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenURI(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// TokenURI is a free data retrieval call binding the contract method 0x94139301.
//
// Solidity: function tokenURI(uint256 projectId, uint256 tokenId, bytes32 seed) view returns(string result)
func (_GenerativeProjectData *GenerativeProjectDataCallerSession) TokenURI(projectId *big.Int, tokenId *big.Int, seed [32]byte) (string, error) {
	return _GenerativeProjectData.Contract.TokenURI(&_GenerativeProjectData.CallOpts, projectId, tokenId, seed)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "changeAdmin", newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeAdmin(&_GenerativeProjectData.TransactOpts, newAdm)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdm) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) ChangeAdmin(newAdm common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeAdmin(&_GenerativeProjectData.TransactOpts, newAdm)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) ChangeParamAddress(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "changeParamAddress", newAddr)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) ChangeParamAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeParamAddress(&_GenerativeProjectData.TransactOpts, newAddr)
}

// ChangeParamAddress is a paid mutator transaction binding the contract method 0x16a5041f.
//
// Solidity: function changeParamAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) ChangeParamAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeParamAddress(&_GenerativeProjectData.TransactOpts, newAddr)
}

// ChangeProjectAddress is a paid mutator transaction binding the contract method 0xe8082858.
//
// Solidity: function changeProjectAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) ChangeProjectAddress(opts *bind.TransactOpts, newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "changeProjectAddress", newAddr)
}

// ChangeProjectAddress is a paid mutator transaction binding the contract method 0xe8082858.
//
// Solidity: function changeProjectAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) ChangeProjectAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeProjectAddress(&_GenerativeProjectData.TransactOpts, newAddr)
}

// ChangeProjectAddress is a paid mutator transaction binding the contract method 0xe8082858.
//
// Solidity: function changeProjectAddress(address newAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) ChangeProjectAddress(newAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.ChangeProjectAddress(&_GenerativeProjectData.TransactOpts, newAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address admin, address paramAddr, address generativeProjectAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) Initialize(opts *bind.TransactOpts, admin common.Address, paramAddr common.Address, generativeProjectAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "initialize", admin, paramAddr, generativeProjectAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address admin, address paramAddr, address generativeProjectAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) Initialize(admin common.Address, paramAddr common.Address, generativeProjectAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.Initialize(&_GenerativeProjectData.TransactOpts, admin, paramAddr, generativeProjectAddr)
}

// Initialize is a paid mutator transaction binding the contract method 0xc0c53b8b.
//
// Solidity: function initialize(address admin, address paramAddr, address generativeProjectAddr) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) Initialize(admin common.Address, paramAddr common.Address, generativeProjectAddr common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.Initialize(&_GenerativeProjectData.TransactOpts, admin, paramAddr, generativeProjectAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.RenounceOwnership(&_GenerativeProjectData.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.RenounceOwnership(&_GenerativeProjectData.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectData *GenerativeProjectDataSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.TransferOwnership(&_GenerativeProjectData.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GenerativeProjectData *GenerativeProjectDataTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GenerativeProjectData.Contract.TransferOwnership(&_GenerativeProjectData.TransactOpts, newOwner)
}

// GenerativeProjectDataInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the GenerativeProjectData contract.
type GenerativeProjectDataInitializedIterator struct {
	Event *GenerativeProjectDataInitialized // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectDataInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectDataInitialized)
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
		it.Event = new(GenerativeProjectDataInitialized)
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
func (it *GenerativeProjectDataInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectDataInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectDataInitialized represents a Initialized event raised by the GenerativeProjectData contract.
type GenerativeProjectDataInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeProjectData *GenerativeProjectDataFilterer) FilterInitialized(opts *bind.FilterOpts) (*GenerativeProjectDataInitializedIterator, error) {

	logs, sub, err := _GenerativeProjectData.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectDataInitializedIterator{contract: _GenerativeProjectData.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_GenerativeProjectData *GenerativeProjectDataFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *GenerativeProjectDataInitialized) (event.Subscription, error) {

	logs, sub, err := _GenerativeProjectData.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectDataInitialized)
				if err := _GenerativeProjectData.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_GenerativeProjectData *GenerativeProjectDataFilterer) ParseInitialized(log types.Log) (*GenerativeProjectDataInitialized, error) {
	event := new(GenerativeProjectDataInitialized)
	if err := _GenerativeProjectData.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GenerativeProjectDataOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GenerativeProjectData contract.
type GenerativeProjectDataOwnershipTransferredIterator struct {
	Event *GenerativeProjectDataOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GenerativeProjectDataOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GenerativeProjectDataOwnershipTransferred)
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
		it.Event = new(GenerativeProjectDataOwnershipTransferred)
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
func (it *GenerativeProjectDataOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GenerativeProjectDataOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GenerativeProjectDataOwnershipTransferred represents a OwnershipTransferred event raised by the GenerativeProjectData contract.
type GenerativeProjectDataOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeProjectData *GenerativeProjectDataFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GenerativeProjectDataOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeProjectData.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GenerativeProjectDataOwnershipTransferredIterator{contract: _GenerativeProjectData.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GenerativeProjectData *GenerativeProjectDataFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GenerativeProjectDataOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GenerativeProjectData.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GenerativeProjectDataOwnershipTransferred)
				if err := _GenerativeProjectData.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GenerativeProjectData *GenerativeProjectDataFilterer) ParseOwnershipTransferred(log types.Log) (*GenerativeProjectDataOwnershipTransferred, error) {
	event := new(GenerativeProjectDataOwnershipTransferred)
	if err := _GenerativeProjectData.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
