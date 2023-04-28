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
	ABI: "[{\"inputs\":[],\"name\":\"FileExists\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidBfsResult\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidURI\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"bfsId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"FileStored\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"bfsId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"chunks\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"}],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"dataStorage\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filenames\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAddresses\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getAllFilenames\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"}],\"name\":\"getId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"}],\"name\":\"load\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"}],\"name\":\"loadWithUri\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"},{\"internalType\":\"int256\",\"name\":\"\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"filename\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"chunkIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"store\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// Filenames is a free data retrieval call binding the contract method 0xc9b09659.
//
// Solidity: function filenames(address , uint256 ) view returns(string)
func (_Bfs *BfsCaller) Filenames(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (string, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "filenames", arg0, arg1)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Filenames is a free data retrieval call binding the contract method 0xc9b09659.
//
// Solidity: function filenames(address , uint256 ) view returns(string)
func (_Bfs *BfsSession) Filenames(arg0 common.Address, arg1 *big.Int) (string, error) {
	return _Bfs.Contract.Filenames(&_Bfs.CallOpts, arg0, arg1)
}

// Filenames is a free data retrieval call binding the contract method 0xc9b09659.
//
// Solidity: function filenames(address , uint256 ) view returns(string)
func (_Bfs *BfsCallerSession) Filenames(arg0 common.Address, arg1 *big.Int) (string, error) {
	return _Bfs.Contract.Filenames(&_Bfs.CallOpts, arg0, arg1)
}

// GetAllAddresses is a free data retrieval call binding the contract method 0x9516a104.
//
// Solidity: function getAllAddresses() view returns(address[])
func (_Bfs *BfsCaller) GetAllAddresses(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "getAllAddresses")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetAllAddresses is a free data retrieval call binding the contract method 0x9516a104.
//
// Solidity: function getAllAddresses() view returns(address[])
func (_Bfs *BfsSession) GetAllAddresses() ([]common.Address, error) {
	return _Bfs.Contract.GetAllAddresses(&_Bfs.CallOpts)
}

// GetAllAddresses is a free data retrieval call binding the contract method 0x9516a104.
//
// Solidity: function getAllAddresses() view returns(address[])
func (_Bfs *BfsCallerSession) GetAllAddresses() ([]common.Address, error) {
	return _Bfs.Contract.GetAllAddresses(&_Bfs.CallOpts)
}

// GetAllFilenames is a free data retrieval call binding the contract method 0xbf1450df.
//
// Solidity: function getAllFilenames(address addr) view returns(string[])
func (_Bfs *BfsCaller) GetAllFilenames(opts *bind.CallOpts, addr common.Address) ([]string, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "getAllFilenames", addr)

	if err != nil {
		return *new([]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]string)).(*[]string)

	return out0, err

}

// GetAllFilenames is a free data retrieval call binding the contract method 0xbf1450df.
//
// Solidity: function getAllFilenames(address addr) view returns(string[])
func (_Bfs *BfsSession) GetAllFilenames(addr common.Address) ([]string, error) {
	return _Bfs.Contract.GetAllFilenames(&_Bfs.CallOpts, addr)
}

// GetAllFilenames is a free data retrieval call binding the contract method 0xbf1450df.
//
// Solidity: function getAllFilenames(address addr) view returns(string[])
func (_Bfs *BfsCallerSession) GetAllFilenames(addr common.Address) ([]string, error) {
	return _Bfs.Contract.GetAllFilenames(&_Bfs.CallOpts, addr)
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

// LoadWithUri is a free data retrieval call binding the contract method 0xabb15219.
//
// Solidity: function loadWithUri(string uri, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsCaller) LoadWithUri(opts *bind.CallOpts, uri string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	var out []interface{}
	err := _Bfs.contract.Call(opts, &out, "loadWithUri", uri, chunkIndex)

	if err != nil {
		return *new([]byte), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// LoadWithUri is a free data retrieval call binding the contract method 0xabb15219.
//
// Solidity: function loadWithUri(string uri, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsSession) LoadWithUri(uri string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	return _Bfs.Contract.LoadWithUri(&_Bfs.CallOpts, uri, chunkIndex)
}

// LoadWithUri is a free data retrieval call binding the contract method 0xabb15219.
//
// Solidity: function loadWithUri(string uri, uint256 chunkIndex) view returns(bytes, int256)
func (_Bfs *BfsCallerSession) LoadWithUri(uri string, chunkIndex *big.Int) ([]byte, *big.Int, error) {
	return _Bfs.Contract.LoadWithUri(&_Bfs.CallOpts, uri, chunkIndex)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bfs *BfsTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Bfs.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bfs *BfsSession) Initialize() (*types.Transaction, error) {
	return _Bfs.Contract.Initialize(&_Bfs.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Bfs *BfsTransactorSession) Initialize() (*types.Transaction, error) {
	return _Bfs.Contract.Initialize(&_Bfs.TransactOpts)
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

// BfsFileStoredIterator is returned from FilterFileStored and is used to iterate over the raw logs and unpacked data for FileStored events raised by the Bfs contract.
type BfsFileStoredIterator struct {
	Event *BfsFileStored // Event containing the contract specifics and raw log

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
func (it *BfsFileStoredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BfsFileStored)
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
		it.Event = new(BfsFileStored)
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
func (it *BfsFileStoredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BfsFileStoredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BfsFileStored represents a FileStored event raised by the Bfs contract.
type BfsFileStored struct {
	Addr       common.Address
	Filename   string
	ChunkIndex *big.Int
	BfsId      *big.Int
	Uri        string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterFileStored is a free log retrieval operation binding the contract event 0x3570439426bde556b3a1036679abf502ade483de5bc63fe5578a687650c96bdd.
//
// Solidity: event FileStored(address indexed addr, string filename, uint256 chunkIndex, uint256 indexed bfsId, string uri)
func (_Bfs *BfsFilterer) FilterFileStored(opts *bind.FilterOpts, addr []common.Address, bfsId []*big.Int) (*BfsFileStoredIterator, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	var bfsIdRule []interface{}
	for _, bfsIdItem := range bfsId {
		bfsIdRule = append(bfsIdRule, bfsIdItem)
	}

	logs, sub, err := _Bfs.contract.FilterLogs(opts, "FileStored", addrRule, bfsIdRule)
	if err != nil {
		return nil, err
	}
	return &BfsFileStoredIterator{contract: _Bfs.contract, event: "FileStored", logs: logs, sub: sub}, nil
}

// WatchFileStored is a free log subscription operation binding the contract event 0x3570439426bde556b3a1036679abf502ade483de5bc63fe5578a687650c96bdd.
//
// Solidity: event FileStored(address indexed addr, string filename, uint256 chunkIndex, uint256 indexed bfsId, string uri)
func (_Bfs *BfsFilterer) WatchFileStored(opts *bind.WatchOpts, sink chan<- *BfsFileStored, addr []common.Address, bfsId []*big.Int) (event.Subscription, error) {

	var addrRule []interface{}
	for _, addrItem := range addr {
		addrRule = append(addrRule, addrItem)
	}

	var bfsIdRule []interface{}
	for _, bfsIdItem := range bfsId {
		bfsIdRule = append(bfsIdRule, bfsIdItem)
	}

	logs, sub, err := _Bfs.contract.WatchLogs(opts, "FileStored", addrRule, bfsIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BfsFileStored)
				if err := _Bfs.contract.UnpackLog(event, "FileStored", log); err != nil {
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

// ParseFileStored is a log parse operation binding the contract event 0x3570439426bde556b3a1036679abf502ade483de5bc63fe5578a687650c96bdd.
//
// Solidity: event FileStored(address indexed addr, string filename, uint256 chunkIndex, uint256 indexed bfsId, string uri)
func (_Bfs *BfsFilterer) ParseFileStored(log types.Log) (*BfsFileStored, error) {
	event := new(BfsFileStored)
	if err := _Bfs.contract.UnpackLog(event, "FileStored", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BfsInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Bfs contract.
type BfsInitializedIterator struct {
	Event *BfsInitialized // Event containing the contract specifics and raw log

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
func (it *BfsInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BfsInitialized)
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
		it.Event = new(BfsInitialized)
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
func (it *BfsInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BfsInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BfsInitialized represents a Initialized event raised by the Bfs contract.
type BfsInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Bfs *BfsFilterer) FilterInitialized(opts *bind.FilterOpts) (*BfsInitializedIterator, error) {

	logs, sub, err := _Bfs.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &BfsInitializedIterator{contract: _Bfs.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Bfs *BfsFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *BfsInitialized) (event.Subscription, error) {

	logs, sub, err := _Bfs.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BfsInitialized)
				if err := _Bfs.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Bfs *BfsFilterer) ParseInitialized(log types.Log) (*BfsInitialized, error) {
	event := new(BfsInitialized)
	if err := _Bfs.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
