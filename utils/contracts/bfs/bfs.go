// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bfs

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

// BfsMetaData contains all meta data concerning the Bfs contract.
var BfsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"FileExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"bfsId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"chunks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"dataStorage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"}],\"name\":\"getId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"}],\"name\":\"load\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BfsABI is the input ABI used to generate the binding from.
// Deprecated: Use BfsMetaData.ABI instead.
var BfsABI = BfsMetaData.ABI

// Bfs is an auto generated Go binding around an Ethereum contract.
type Bfs struct {
	BfsCaller     // Read-only binding to the contract
	BfsTransactor // Write-only binding to the contract
	BfsFilterer   // Log filterer for contract events
}

// BfsCaller is an auto generated read-only Go binding around an Ethereum contract.
type BfsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BfsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BfsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BfsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BfsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BfsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BfsSession struct {
	Contract     *Bfs              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BfsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BfsCallerSession struct {
	Contract *BfsCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BfsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BfsTransactorSession struct {
	Contract     *BfsTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BfsRaw is an auto generated low-level Go binding around an Ethereum contract.
type BfsRaw struct {
	Contract *Bfs // Generic contract binding to access the raw methods on
}

// BfsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BfsCallerRaw struct {
	Contract *BfsCaller // Generic read-only contract binding to access the raw methods on
}

// BfsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BfsTransactorRaw struct {
	Contract *BfsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBfs creates a new instance of Bfs, bound to a specific deployed contract.
func NewBfs(address common.Address, backend bind.ContractBackend) (*Bfs, error) {
	contract, err := bindBfs(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Bfs{BfsCaller: BfsCaller{contract: contract}, BfsTransactor: BfsTransactor{contract: contract}, BfsFilterer: BfsFilterer{contract: contract}}, nil
}

// NewBfsCaller creates a new read-only instance of Bfs, bound to a specific deployed contract.
func NewBfsCaller(address common.Address, caller bind.ContractCaller) (*BfsCaller, error) {
	contract, err := bindBfs(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BfsCaller{contract: contract}, nil
}

// NewBfsTransactor creates a new write-only instance of Bfs, bound to a specific deployed contract.
func NewBfsTransactor(address common.Address, transactor bind.ContractTransactor) (*BfsTransactor, error) {
	contract, err := bindBfs(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BfsTransactor{contract: contract}, nil
}

// NewBfsFilterer creates a new log filterer instance of Bfs, bound to a specific deployed contract.
func NewBfsFilterer(address common.Address, filterer bind.ContractFilterer) (*BfsFilterer, error) {
	contract, err := bindBfs(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BfsFilterer{contract: contract}, nil
}

// bindBfs binds a generic wrapper to an already deployed contract.
func bindBfs(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BfsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bfs *BfsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bfs.Contract.BfsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bfs *BfsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bfs.Contract.BfsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bfs *BfsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bfs.Contract.BfsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Bfs *BfsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Bfs.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Bfs *BfsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bfs.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Bfs *BfsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Bfs.Contract.contract.Transact(opts, method, params...)
}

// BfsId is a free data retrieval call binding the contract method 0xe0f64a38.
//
// Solidity: function bfsId(address , string ) view returns(uint256)
func (_Bfs *BfsCaller) BfsId(opts *bind.CallOpts, arg0 common.Address, arg1 string) (*big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "bfsId", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BfsId is a free data retrieval call binding the contract method 0xe0f64a38.
//
// Solidity: function bfsId(address , string ) view returns(uint256)
func (_Bfs *BfsSession) BfsId(arg0 common.Address, arg1 string) (*big.Int, error) {
	return _Bfs.Contract.BfsId(&_Bfs.CallOpts, arg0, arg1)
}

// BfsId is a free data retrieval call binding the contract method 0xe0f64a38.
//
// Solidity: function bfsId(address , string ) view returns(uint256)
func (_Bfs *BfsCallerSession) BfsId(arg0 common.Address, arg1 string) (*big.Int, error) {
	return _Bfs.Contract.BfsId(&_Bfs.CallOpts, arg0, arg1)
}

// Chunks is a free data retrieval call binding the contract method 0x85273a90.
//
// Solidity: function chunks(address , string ) view returns(uint256)
func (_Bfs *BfsCaller) Chunks(opts *bind.CallOpts, arg0 common.Address, arg1 string) (*big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "chunks", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Chunks is a free data retrieval call binding the contract method 0x85273a90.
//
// Solidity: function chunks(address , string ) view returns(uint256)
func (_Bfs *BfsSession) Chunks(arg0 common.Address, arg1 string) (*big.Int, error) {
	return _Bfs.Contract.Chunks(&_Bfs.CallOpts, arg0, arg1)
}

// Chunks is a free data retrieval call binding the contract method 0x85273a90.
//
// Solidity: function chunks(address , string ) view returns(uint256)
func (_Bfs *BfsCallerSession) Chunks(arg0 common.Address, arg1 string) (*big.Int, error) {
	return _Bfs.Contract.Chunks(&_Bfs.CallOpts, arg0, arg1)
}

// Count is a free data retrieval call binding the contract method 0xba67e641.
//
// Solidity: function count(address addr, string filename) view returns(uint256)
func (_Bfs *BfsCaller) Count(opts *bind.CallOpts, addr common.Address, filename string) (*big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "count", addr, filename)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0xba67e641.
//
// Solidity: function count(address addr, string filename) view returns(uint256)
func (_Bfs *BfsSession) Count(addr common.Address, filename string) (*big.Int, error) {
	return _Bfs.Contract.Count(&_Bfs.CallOpts, addr, filename)
}

// Count is a free data retrieval call binding the contract method 0xba67e641.
//
// Solidity: function count(address addr, string filename) view returns(uint256)
func (_Bfs *BfsCallerSession) Count(addr common.Address, filename string) (*big.Int, error) {
	return _Bfs.Contract.Count(&_Bfs.CallOpts, addr, filename)
}

// DataStorage is a free data retrieval call binding the contract method 0xd92e741e.
//
// Solidity: function dataStorage(address , string , uint256 ) view returns(bytes)
func (_Bfs *BfsCaller) DataStorage(opts *bind.CallOpts, arg0 common.Address, arg1 string, arg2 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "dataStorage", arg0, arg1, arg2)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DataStorage is a free data retrieval call binding the contract method 0xd92e741e.
//
// Solidity: function dataStorage(address , string , uint256 ) view returns(bytes)
func (_Bfs *BfsSession) DataStorage(arg0 common.Address, arg1 string, arg2 *big.Int) ([]byte, error) {
	return _Bfs.Contract.DataStorage(&_Bfs.CallOpts, arg0, arg1, arg2)
}

// DataStorage is a free data retrieval call binding the contract method 0xd92e741e.
//
// Solidity: function dataStorage(address , string , uint256 ) view returns(bytes)
func (_Bfs *BfsCallerSession) DataStorage(arg0 common.Address, arg1 string, arg2 *big.Int) ([]byte, error) {
	return _Bfs.Contract.DataStorage(&_Bfs.CallOpts, arg0, arg1, arg2)
}

// GetId is a free data retrieval call binding the contract method 0x8a5ca90b.
//
// Solidity: function getId(address addr, string filename) view returns(uint256)
func (_Bfs *BfsCaller) GetId(opts *bind.CallOpts, addr common.Address, filename string) (*big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "getId", addr, filename)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetId is a free data retrieval call binding the contract method 0x8a5ca90b.
//
// Solidity: function getId(address addr, string filename) view returns(uint256)
func (_Bfs *BfsSession) GetId(addr common.Address, filename string) (*big.Int, error) {
	return _Bfs.Contract.GetId(&_Bfs.CallOpts, addr, filename)
}

// GetId is a free data retrieval call binding the contract method 0x8a5ca90b.
//
// Solidity: function getId(address addr, string filename) view returns(uint256)
func (_Bfs *BfsCallerSession) GetId(addr common.Address, filename string) (*big.Int, error) {
	return _Bfs.Contract.GetId(&_Bfs.CallOpts, addr, filename)
}

// Load is a free data retrieval call binding the contract method 0xa0e56764.
//
// Solidity: function load(address addr, string filename, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsCaller) Load(opts *bind.CallOpts, addr common.Address, filename string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "load", addr, filename, chunkIndex)

	if err != nil {
		return *new([]byte), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// Load is a free data retrieval call binding the contract method 0xa0e56764.
//
// Solidity: function load(address addr, string filename, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsSession) Load(addr common.Address, filename string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	return _Bfs.Contract.Load(&_Bfs.CallOpts, addr, filename, chunkIndex)
}

// Load is a free data retrieval call binding the contract method 0xa0e56764.
//
// Solidity: function load(address addr, string filename, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsCallerSession) Load(addr common.Address, filename string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	return _Bfs.Contract.Load(&_Bfs.CallOpts, addr, filename, chunkIndex)
}

// Store is a paid mutator transaction binding the contract method 0x928a4b3e.
//
// Solidity: function store(string filename, uint256 chunkIndex, bytes _data) returns()
func (_Bfs *BfsTransactor) Store(opts *bind.TransactOpts, filename string, chunkIndex *big.Int, _data []byte) (*types.Transaction, error) {
	return _Bfs.contract.Transact(opts, "store", filename, chunkIndex, _data)
}

// Store is a paid mutator transaction binding the contract method 0x928a4b3e.
//
// Solidity: function store(string filename, uint256 chunkIndex, bytes _data) returns()
func (_Bfs *BfsSession) Store(filename string, chunkIndex *big.Int, _data []byte) (*types.Transaction, error) {
	return _Bfs.Contract.Store(&_Bfs.TransactOpts, filename, chunkIndex, _data)
}

// Store is a paid mutator transaction binding the contract method 0x928a4b3e.
//
// Solidity: function store(string filename, uint256 chunkIndex, bytes _data) returns()
func (_Bfs *BfsTransactorSession) Store(filename string, chunkIndex *big.Int, _data []byte) (*types.Transaction, error) {
	return _Bfs.Contract.Store(&_Bfs.TransactOpts, filename, chunkIndex, _data)
}
